package level

// tick.go — Level 方块 tick 系统
// 对应 PHP Level.php 中的 doTick/tickChunks/scheduleUpdate/updateAround
//
// 三大子系统：
//   1. 计划更新队列 (Scheduled Block Updates) — 红石/液体/活板门等延迟更新
//   2. 随机方块 Tick — 作物生长、草蔓延、冰融化等
//   3. 邻居方块更新 — 放/破方块后通知6个相邻方块

import (
	"container/heap"
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/world"
)

// ── 配置常量 ──────────────────────────────────────

const (
	// RandomTickSpeed 每个 chunk section 每 tick 随机选择的方块数
	// 对应 PHP $server->randomTickSpeed（默认 = 3）
	RandomTickSpeed = 3

	// MaxScheduledUpdatesPerTick 每 tick 最多处理的计划更新数
	// 对应 PHP $server->maxBlockUpdatesPreTick
	MaxScheduledUpdatesPerTick = 512
)

// ── 随机 Tick 方块注册表 ──────────────────────────────

// randomTickBlocks 需要随机 tick 的方块 ID 集合
// 对应 PHP Level::$randomTickBlocks
var randomTickBlocks = map[byte]bool{
	block.GRASS:        true, // 草蔓延/死亡
	block.FARMLAND:     true, // 耕地干燥
	block.WHEAT_BLOCK:  true, // 小麦生长
	block.CARROT_BLOCK: true, // 胡萝卜生长
	block.POTATO_BLOCK: true, // 马铃薯生长

	block.BEETROOT_BLOCK:  true, // 甜菜生长
	block.SUGARCANE_BLOCK: true, // 甘蔗生长
	block.CACTUS:          true, // 仙人掌生长
	block.PUMPKIN_STEM:    true, // 南瓜梗
	block.MELON_STEM:      true, // 西瓜梗

	block.SAPLING: true, // 树苗生长
	block.LEAVES:  true, // 树叶腐烂
	block.LEAVES2: true,

	block.SNOW_LAYER: true, // 雪层融化
	block.ICE:        true, // 冰融化
	block.MYCELIUM:   true, // 菌丝蔓延

	block.VINE: true, // 藤蔓生长
	block.FIRE: true, // 火焰蔓延/熄灭
}

// RegisterRandomTickBlock 注册一个需要随机 tick 的方块 ID
func RegisterRandomTickBlock(blockID byte) {
	randomTickBlocks[blockID] = true
}

// ── 计划更新优先队列 ──────────────────────────────────

// ScheduledUpdate 计划更新条目
type ScheduledUpdate struct {
	X, Y, Z  int32 // 方块世界坐标
	Priority int64 // 执行 tick（越小越先执行）
	index    int   // heap 内部索引
}

// scheduledUpdateQueue 优先队列（最小堆）
type scheduledUpdateQueue []*ScheduledUpdate

func (q scheduledUpdateQueue) Len() int           { return len(q) }
func (q scheduledUpdateQueue) Less(i, j int) bool { return q[i].Priority < q[j].Priority }
func (q scheduledUpdateQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *scheduledUpdateQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(*ScheduledUpdate)
	item.index = n
	*q = append(*q, item)
}

func (q *scheduledUpdateQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*q = old[:n-1]
	return item
}

// ── Level tick 字段 ──────────────────────────────────

// TickState 存储 Level 的 tick 相关状态
type TickState struct {
	// 计划更新队列
	updateQueue      scheduledUpdateQueue
	updateQueueIndex map[int64]int64 // blockHash → 已注册的 delay

	// 当前全局 tick 计数器
	currentTick int64
}

// NewTickState 创建 tick 状态
func NewTickState() *TickState {
	ts := &TickState{
		updateQueue:      make(scheduledUpdateQueue, 0),
		updateQueueIndex: make(map[int64]int64),
	}
	heap.Init(&ts.updateQueue)
	return ts
}

// GetCurrentTick 返回当前 tick 数
func (l *Level) GetCurrentTick() int64 {
	return l.tickState.currentTick
}

// ── 计划更新 ──────────────────────────────────────

// ScheduleUpdate 安排一个方块的延迟更新
// 对应 PHP Level::scheduleUpdate(Vector3 $pos, int $delay)
//
// delay 是延迟的 tick 数（例如红石 = 2, 液体 = 5/30）
func (l *Level) ScheduleUpdate(x, y, z int32, delay int) {
	hash := blockHash(x, y, z)

	// 如果已有更短延迟的更新，跳过
	if existingDelay, ok := l.tickState.updateQueueIndex[hash]; ok {
		if existingDelay <= int64(delay) {
			return
		}
	}

	l.tickState.updateQueueIndex[hash] = int64(delay)
	heap.Push(&l.tickState.updateQueue, &ScheduledUpdate{
		X:        x,
		Y:        y,
		Z:        z,
		Priority: l.tickState.currentTick + int64(delay),
	})
}

// processScheduledUpdates 处理到期的计划更新
// 对应 PHP Level::actuallyDoTick() L708-724
func (l *Level) processScheduledUpdates() {
	processed := 0

	for l.tickState.updateQueue.Len() > 0 && processed < MaxScheduledUpdatesPerTick {
		// peek
		top := l.tickState.updateQueue[0]
		if top.Priority > l.tickState.currentTick {
			break
		}

		// extract
		item := heap.Pop(&l.tickState.updateQueue).(*ScheduledUpdate)
		delete(l.tickState.updateQueueIndex, blockHash(item.X, item.Y, item.Z))

		// 检查区块是否加载
		chunkX := item.X >> 4
		chunkZ := item.Z >> 4
		if !l.IsChunkLoaded(chunkX, chunkZ) {
			continue
		}

		// 获取方块行为并调用 OnUpdate
		bs := l.GetBlock(item.X, item.Y, item.Z)
		behavior := block.Registry.GetBehavior(bs.ID)
		if behavior != nil {
			ctx := &block.BlockContext{
				X: int(item.X), Y: int(item.Y), Z: int(item.Z),
				Meta: bs.Meta,
			}
			behavior.OnUpdate(ctx, BlockUpdateScheduled)
		}
		processed++
	}
}

// ── 随机方块 Tick ──────────────────────────────────

// tickChunks 在已加载的区块中执行随机方块 tick
// 对应 PHP Level::tickChunks()
//
// 算法: 对每个 chunk 的每个非空 section，随机选择 RandomTickSpeed 个方块，
// 如果该方块注册了随机 tick，则调用 OnUpdate(BLOCK_UPDATE_RANDOM)
func (l *Level) tickChunks() {
	l.mu.RLock()
	chunks := make([]*world.Chunk, 0, len(l.Chunks))
	for _, c := range l.Chunks {
		chunks = append(chunks, c)
	}
	l.mu.RUnlock()

	for _, chunk := range chunks {
		l.tickChunk(chunk)
	}
}

// tickChunk 对单个区块执行随机 tick
func (l *Level) tickChunk(chunk *world.Chunk) {
	if chunk == nil {
		return
	}

	chunkX := int32(chunk.X)
	chunkZ := int32(chunk.Z)

	// 遍历 8 个 section (Y=0..7, 对应 y=0..127)
	for sectionY := 0; sectionY < 8; sectionY++ {
		section := chunk.Sections[sectionY]
		if section == nil || section.IsEmpty() {
			continue
		}

		// 每个 section 随机选择 RandomTickSpeed 个方块
		// 对应 PHP L962-980
		k := rand.Int63()
		for i := 0; i < RandomTickSpeed; i++ {
			x := int(k & 0x0f)
			y := int((k >> 4) & 0x0f)
			z := int((k >> 8) & 0x0f)
			k >>= 12

			blockID := section.GetBlockId(x, y, z)
			if randomTickBlocks[blockID] {
				worldY := (sectionY << 4) + y
				meta := section.GetBlockData(x, y, z)

				behavior := block.Registry.GetBehavior(blockID)
				if behavior != nil {
					ctx := &block.BlockContext{
						X:    int(chunkX)*16 + x,
						Y:    worldY,
						Z:    int(chunkZ)*16 + z,
						Meta: meta,
					}
					behavior.OnUpdate(ctx, BlockUpdateRandom)
				}
			}
		}
	}
}

// ── 邻居方块更新 ──────────────────────────────────

// UpdateAround 通知6个相邻方块发生了 BLOCK_UPDATE_NORMAL
// 对应 PHP Level::updateAround()
func (l *Level) UpdateAround(x, y, z int32) {
	offsets := [6][3]int32{
		{0, -1, 0}, {0, 1, 0},
		{-1, 0, 0}, {1, 0, 0},
		{0, 0, -1}, {0, 0, 1},
	}

	for _, off := range offsets {
		nx, ny, nz := x+off[0], y+off[1], z+off[2]
		if ny < YMin || ny >= YMax {
			continue
		}

		bs := l.GetBlock(nx, ny, nz)
		behavior := block.Registry.GetBehavior(bs.ID)
		if behavior != nil {
			ctx := &block.BlockContext{
				X: int(nx), Y: int(ny), Z: int(nz),
				Meta: bs.Meta,
			}
			behavior.OnUpdate(ctx, BlockUpdateNormal)
		}
	}
}

// ── blockHash ──────────────────────────────────────

// blockHash 方块坐标哈希
func blockHash(x, y, z int32) int64 {
	return (int64(x) << 32) | (int64(z) & 0xFFFFFFFF) ^ (int64(y) << 48)
}
