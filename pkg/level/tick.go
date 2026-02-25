package level

import (
	"container/heap"
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	RandomTickSpeed = 3
	MaxScheduledUpdatesPerTick = 512
)
var randomTickBlocks = map[byte]bool{
	block.GRASS:        true,
	block.FARMLAND:     true,
	block.WHEAT_BLOCK:  true,
	block.CARROT_BLOCK: true,
	block.POTATO_BLOCK: true,

	block.BEETROOT_BLOCK:  true,
	block.SUGARCANE_BLOCK: true,
	block.CACTUS:          true,
	block.PUMPKIN_STEM:    true,
	block.MELON_STEM:      true,

	block.SAPLING: true,
	block.LEAVES:  true,
	block.LEAVES2: true,

	block.SNOW_LAYER: true,
	block.ICE:        true,
	block.MYCELIUM:   true,

	block.VINE: true,
	block.FIRE: true,
}
func RegisterRandomTickBlock(blockID byte) {
	randomTickBlocks[blockID] = true
}
type ScheduledUpdate struct {
	X, Y, Z  int32
	Priority int64
	index    int
}
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
type TickState struct {
	updateQueue      scheduledUpdateQueue
	updateQueueIndex map[int64]int64
	currentTick int64
}
func NewTickState() *TickState {
	ts := &TickState{
		updateQueue:      make(scheduledUpdateQueue, 0),
		updateQueueIndex: make(map[int64]int64),
	}
	heap.Init(&ts.updateQueue)
	return ts
}
func (l *Level) GetCurrentTick() int64 {
	return l.tickState.currentTick
}
func (l *Level) ScheduleUpdate(x, y, z int32, delay int) {
	hash := blockHash(x, y, z)
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
func (l *Level) processScheduledUpdates() {
	processed := 0

	for l.tickState.updateQueue.Len() > 0 && processed < MaxScheduledUpdatesPerTick {
		top := l.tickState.updateQueue[0]
		if top.Priority > l.tickState.currentTick {
			break
		}

		item := heap.Pop(&l.tickState.updateQueue).(*ScheduledUpdate)
		delete(l.tickState.updateQueueIndex, blockHash(item.X, item.Y, item.Z))

		chunkX := item.X >> 4
		chunkZ := item.Z >> 4
		if !l.IsChunkLoaded(chunkX, chunkZ) {
			continue
		}

		bs := l.GetBlock(item.X, item.Y, item.Z)
		behavior := block.Registry.GetBehavior(bs.ID)
		if behavior != nil {
			ctx := &block.BlockContext{
				X: int(item.X), Y: int(item.Y), Z: int(item.Z),
				Meta:    bs.Meta,
				Powered: l.getBlockPowered(bs.ID, bs.Meta, item.X, item.Y, item.Z),
			}
			behavior.OnUpdate(ctx, BlockUpdateScheduled)
			l.applyBlockContextResult(ctx, item.X, item.Y, item.Z)
		}
		processed++
	}
}
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
func (l *Level) tickChunk(chunk *world.Chunk) {
	if chunk == nil {
		return
	}

	chunkX := int32(chunk.X)
	chunkZ := int32(chunk.Z)
	for sectionY := 0; sectionY < 8; sectionY++ {
		section := chunk.Sections[sectionY]
		if section == nil || section.IsEmpty() {
			continue
		}
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
		if bs.ID == block.REDSTONE_WIRE {
			l.UpdateRedstoneWire(nx, ny, nz)
			continue
		}
		behavior := block.Registry.GetBehavior(bs.ID)
		if behavior != nil {
			ctx := &block.BlockContext{
				X: int(nx), Y: int(ny), Z: int(nz),
				Meta:    bs.Meta,
				Powered: l.getBlockPowered(bs.ID, bs.Meta, nx, ny, nz),
			}
			behavior.OnUpdate(ctx, BlockUpdateNormal)
			l.applyBlockContextResult(ctx, nx, ny, nz)
		}
	}
}

func (l *Level) applyBlockContextResult(ctx *block.BlockContext, x, y, z int32) {
	if ctx.ReplaceBlockID != 0 {
		l.SetBlock(x, y, z, ctx.ReplaceBlockID, ctx.ReplaceBlockMeta, false)
		l.PendingBlockUpdates = append(l.PendingBlockUpdates, PendingBlockUpdate{
			X: x, Y: y, Z: z,
			ID: ctx.ReplaceBlockID, Meta: ctx.ReplaceBlockMeta,
		})
		l.UpdateAround(x, y, z)
	}
	if ctx.ScheduleDelay > 0 {
		l.ScheduleUpdate(x, y, z, ctx.ScheduleDelay)
	}
}
func blockHash(x, y, z int32) int64 {
	return (int64(x) << 32) | (int64(z) & 0xFFFFFFFF) ^ (int64(y) << 48)
}

func (l *Level) getBlockPowered(bid byte, meta byte, x, y, z int32) bool {
	switch bid {
	case block.REDSTONE_TORCH, block.UNLIT_REDSTONE_TORCH:
		ax, ay, az := attachmentOffset(meta, x, y, z)
		return l.IsBlockPowered(ax, ay, az)
	default:
		return l.IsBlockPowered(x, y, z)
	}
}

func attachmentOffset(meta byte, x, y, z int32) (int32, int32, int32) {
	switch meta & 0x07 {
	case 1:
		return x - 1, y, z
	case 2:
		return x + 1, y, z
	case 3:
		return x, y, z - 1
	case 4:
		return x, y, z + 1
	case 5:
		return x, y - 1, z
	default:
		return x, y, z
	}
}
