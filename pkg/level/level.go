package level

import (
	"math"
	"sync"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/entity"
	"github.com/scaxe/scaxe-go/pkg/level/generator"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/tile"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	YMax = 128
	YMin = 0

	BlockUpdateNormal    = 1
	BlockUpdateRandom    = 2
	BlockUpdateScheduled = 3
	BlockUpdateWeak      = 4
	BlockUpdateTouch     = 5
	BlockUpdateRedstone  = 6

	TimeDay     = 0
	TimeSunset  = 12000
	TimeNight   = 14000
	TimeSunrise = 23000
	TimeFull    = 24000

	DimensionNormal = 0
	DimensionNether = 1
	DimensionEnd    = 2
)

type Level struct {
	mu sync.RWMutex

	ID   int
	Name string
	Path string

	Provider Provider

	Generator generator.Generator
	Seed      int64

	Chunks map[int64]*world.Chunk

	Entities map[int64]entity.IEntity

	Time     int64
	StopTime bool

	Dimension int

	Closed bool

	tickState *TickState
	Tiles     *tile.TileManager
}

var levelCounter int = 1

func NewLevel(name string, path string, provider Provider, generatorName string) *Level {
	l := &Level{
		ID:        levelCounter,
		Name:      name,
		Path:      path,
		Provider:  provider,
		Chunks:    make(map[int64]*world.Chunk),
		Entities:  make(map[int64]entity.IEntity),
		Time:      0,
		StopTime:  false,
		Dimension: DimensionNormal,
		Closed:    false,
		Seed:      12345,
		tickState: NewTickState(),
		Tiles:     tile.NewTileManager(),
	}
	levelCounter++

	if generatorName == "" || generatorName == "DEFAULT" {
		generatorName = "gorigional"
	}

	l.Generator = generator.GetGenerator(generatorName, nil)
	if l.Generator != nil {
		l.Generator.Init(l, l.Seed)
		logger.Info("Level generator initialized", "name", l.Generator.GetName())
	} else {
		logger.Warn("Unknown generator, falling back to gorigional", "name", generatorName)
		l.Generator = generator.GetGenerator("gorigional", nil)
		if l.Generator != nil {
			l.Generator.Init(l, l.Seed)
		} else {
			logger.Error("Failed to load any generator, world generation will not work")
		}
	}

	logger.Info("Level created", "name", name, "id", l.ID, "provider", provider.GetName())
	return l
}

func (l *Level) GetChunk(x, z int32, generate bool) *world.Chunk {
	hash := world.ChunkHash(x, z)
	l.mu.RLock()
	chunk, exists := l.Chunks[hash]
	l.mu.RUnlock()

	if exists {
		return chunk
	}

	if l.Provider != nil {
		c, err := l.Provider.LoadChunk(x, z)
		if err == nil && c != nil {
			l.mu.Lock()

			if existing, ok := l.Chunks[hash]; ok {
				l.mu.Unlock()
				return existing
			}
			l.Chunks[hash] = c
			l.mu.Unlock()

			// 从区块 NBT 恢复 Tile 实体（告示牌/箱子/熔炉等）
			l.loadTilesFromChunk(c)

			return c
		}
	}

	if !generate {
		return nil
	}

	if l.Generator != nil {
		l.Generator.GenerateChunk(x, z)
		l.Generator.PopulateChunk(x, z)
	}

	l.mu.RLock()
	chunk = l.Chunks[hash]
	l.mu.RUnlock()
	return chunk
}

func (l *Level) GetSeed() int64 {
	return l.Seed
}

func (l *Level) SetChunk(x, z int32, chunk *world.Chunk) {
	hash := world.ChunkHash(x, z)
	l.mu.Lock()
	l.Chunks[hash] = chunk
	l.mu.Unlock()
}

func (l *Level) IsChunkLoaded(x, z int32) bool {
	hash := world.ChunkHash(x, z)
	l.mu.RLock()
	_, exists := l.Chunks[hash]
	l.mu.RUnlock()
	return exists
}

func (l *Level) UnloadChunk(x, z int32, safe bool, save bool) bool {
	hash := world.ChunkHash(x, z)
	l.mu.Lock()
	defer l.mu.Unlock()

	chunk, exists := l.Chunks[hash]
	if !exists {
		return true
	}

	if save && chunk.HasChanged() {
		if l.Provider != nil {
			err := l.Provider.SaveChunk(chunk)
			if err != nil {
				logger.Error("Failed to save chunk on unload", "x", x, "z", z, "error", err)
			} else {
				chunk.SetChanged(false)
			}
		}
	}

	delete(l.Chunks, hash)
	return true
}

func (l *Level) RequestChunk(x, z int32, loader ChunkLoader) {

	chunk := l.GetChunk(x, z, false)
	if chunk != nil {
		loader.OnChunkLoaded(chunk)
		return
	}

	l.AsyncLoadChunk(x, z, func(c *world.Chunk) {
		if c != nil {

			loader.OnChunkLoaded(c)
		}
	})
}

func (l *Level) AsyncLoadChunk(x, z int32, callback func(*world.Chunk)) {

	if chunk := l.GetChunk(x, z, false); chunk != nil {
		if callback != nil {
			callback(chunk)
		}
		return
	}

	go func() {
		var chunk *world.Chunk
		var err error

		if l.Provider != nil {
			chunk, err = l.Provider.LoadChunk(x, z)
			if err != nil {
				logger.Error("AsyncLoadChunk failed", "x", x, "z", z, "error", err)
			}
		}

		if chunk == nil {

			if l.Generator != nil {
				l.Generator.GenerateChunk(x, z)
				l.Generator.PopulateChunk(x, z)
			}

			chunk = l.GetChunk(x, z, false)
		} else {

			l.SetChunk(x, z, chunk)
		}

		if callback != nil {
			callback(chunk)
		}
	}()
}

func (l *Level) GetBlock(x, y, z int32) block.BlockState {
	if y < YMin || y >= YMax {
		return block.NewBlockState(block.AIR, 0)
	}
	chunkX := x >> 4
	chunkZ := z >> 4
	chunk := l.GetChunk(chunkX, chunkZ, false)
	if chunk == nil {
		return block.NewBlockState(block.AIR, 0)
	}
	localX := x & 0x0F
	localZ := z & 0x0F
	id := chunk.GetBlockId(int(localX), int(y), int(localZ))
	meta := chunk.GetBlockData(int(localX), int(y), int(localZ))
	return block.NewBlockState(id, meta)
}

func (l *Level) SetBlock(x, y, z int32, id, meta byte, update bool) bool {
	if y < YMin || y >= YMax {
		return false
	}
	chunkX := x >> 4
	chunkZ := z >> 4
	chunk := l.GetChunk(chunkX, chunkZ, true)
	if chunk == nil {
		return false
	}
	localX := x & 0x0F
	localZ := z & 0x0F
	chunk.SetBlock(int(localX), int(y), int(localZ), id, meta)

	if update {
		l.UpdateBlockLight(x, y, z, -1)

		l.UpdateBlockLight(x+1, y, z, -1)
		l.UpdateBlockLight(x-1, y, z, -1)
		l.UpdateBlockLight(x, y+1, z, -1)
		l.UpdateBlockLight(x, y-1, z, -1)
		l.UpdateBlockLight(x, y, z+1, -1)
		l.UpdateBlockLight(x, y, z-1, -1)

		// 通知6个相邻方块发生更新（对应 PHP Level::setBlock 中调用 updateAround）
		l.UpdateAround(x, y, z)
	}

	return true
}

func (l *Level) GetBlockId(x, y, z int32) byte {
	if y < YMin || y >= YMax {
		return block.AIR
	}
	chunkX := x >> 4
	chunkZ := z >> 4
	chunk := l.GetChunk(chunkX, chunkZ, false)
	if chunk == nil {
		return block.AIR
	}
	return chunk.GetBlockId(int(x&0x0F), int(y), int(z&0x0F))
}

func (l *Level) GetBlockData(x, y, z int32) byte {
	if y < YMin || y >= YMax {
		return 0
	}
	chunkX := x >> 4
	chunkZ := z >> 4
	chunk := l.GetChunk(chunkX, chunkZ, false)
	if chunk == nil {
		return 0
	}
	return chunk.GetBlockData(int(x&0x0F), int(y), int(z&0x0F))
}

func (l *Level) GetHeight(x, z int32) int32 {
	chunkX := x >> 4
	chunkZ := z >> 4
	chunk := l.GetChunk(chunkX, chunkZ, false)
	if chunk == nil {
		return 0
	}
	localX := int(x & 0x0F)
	localZ := int(z & 0x0F)

	for y := YMax - 1; y >= 0; y-- {
		if chunk.GetBlockId(localX, y, localZ) != 0 {
			return int32(y + 1)
		}
	}
	return 0
}

func (l *Level) GetTime() int64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.Time
}

func (l *Level) SetTime(time int64) {
	l.mu.Lock()
	l.Time = time % TimeFull
	l.mu.Unlock()
}

func (l *Level) AddEntity(e entity.IEntity) {
	l.mu.Lock()
	l.Entities[e.GetID()] = e
	l.mu.Unlock()
}

func (l *Level) RemoveEntity(e entity.IEntity) {
	l.mu.Lock()
	delete(l.Entities, e.GetID())
	l.mu.Unlock()
}

func (l *Level) GetEntities() []entity.IEntity {
	l.mu.RLock()
	defer l.mu.RUnlock()
	entities := make([]entity.IEntity, 0, len(l.Entities))
	for _, e := range l.Entities {
		entities = append(entities, e)
	}
	return entities
}

// GetEntityByID 通过 ID 查找实体
func (l *Level) GetEntityByID(id int64) entity.IEntity {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.Entities[id]
}

func (l *Level) GetNearbyEntities(bb *entity.AxisAlignedBB, except entity.IEntity) []entity.IEntity {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var nearby []entity.IEntity

	entityCount := len(l.Entities)
	if entityCount > 0 {
		for _, e := range l.Entities {
			if e == except {
				continue
			}
			eBB := e.GetBoundingBox()
			if eBB == nil {
				continue
			}
			if eBB.IntersectsWith(bb) {
				nearby = append(nearby, e)
			}
		}
	}
	return nearby
}

func (l *Level) GetCollisionCubes(e entity.IEntity, bb *entity.AxisAlignedBB, includeEntities bool) []*entity.AxisAlignedBB {
	minX := int32(math.Floor(bb.MinX))
	minY := int32(math.Floor(bb.MinY))
	minZ := int32(math.Floor(bb.MinZ))
	maxX := int32(math.Floor(bb.MaxX + 1.0))
	maxY := int32(math.Floor(bb.MaxY + 1.0))
	maxZ := int32(math.Floor(bb.MaxZ + 1.0))

	var collisions []*entity.AxisAlignedBB

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			for z := minZ; z <= maxZ; z++ {
				blk := l.GetBlock(x, y, z)

				if blk.ID != block.AIR && blk.ID != block.TALL_GRASS && blk.ID != block.DANDELION && blk.ID != block.RED_FLOWER {

					blockBB := entity.NewAxisAlignedBB(
						float64(x), float64(y), float64(z),
						float64(x)+1, float64(y)+1, float64(z)+1,
					)
					if blockBB.IntersectsWith(bb) {
						collisions = append(collisions, blockBB)
					}
				}
			}
		}
	}

	if includeEntities {
		entities := l.GetNearbyEntities(bb, e)
		for _, ent := range entities {
			entBB := ent.GetBoundingBox()
			if entBB != nil && entBB.IntersectsWith(bb) {
				collisions = append(collisions, entBB)
			}
		}
	}

	return collisions
}

func (l *Level) Tick() {
	l.mu.Lock()
	if !l.StopTime {
		l.Time++
		if l.Time >= TimeFull {
			l.Time = 0
		}
	}
	l.tickState.currentTick++

	entities := make([]entity.IEntity, 0, len(l.Entities))
	for _, e := range l.Entities {
		entities = append(entities, e)
	}
	l.mu.Unlock()

	// 实体 tick
	for _, e := range entities {
		if !e.Tick(l.Time) {
			l.RemoveEntity(e)
		}
	}

	// 计划方块更新（红石/液体等延迟更新）
	l.processScheduledUpdates()

	// 随机方块 tick（作物生长/草蔓延等）
	l.tickChunks()

	// Tile Entity tick（熔炉烧炼、酿造台等）
	// 对应 PHP Level::actuallyDoTick() L738-751
	l.Tiles.TickUpdates()
}

func (l *Level) Save() {
	l.mu.RLock()
	defer l.mu.RUnlock()

	savedCount := 0
	for _, chunk := range l.Chunks {
		if chunk.HasChanged() && l.Provider != nil {
			if err := l.Provider.SaveChunk(chunk); err != nil {
				logger.Error("Failed to save chunk", "x", chunk.X, "z", chunk.Z, "error", err)
			} else {
				chunk.SetChanged(false)
				savedCount++
			}
		}
	}
	if savedCount > 0 {
		logger.Debug("Level saved", "name", l.Name, "chunks", savedCount)
	}
}

func (l *Level) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for hash, chunk := range l.Chunks {
		if chunk.HasChanged() && l.Provider != nil {
			if err := l.Provider.SaveChunk(chunk); err != nil {
				logger.Error("Failed to save chunk on close", "x", chunk.X, "z", chunk.Z, "error", err)
			}
		}
		delete(l.Chunks, hash)
	}

	for _, e := range l.Entities {
		e.Close()
	}
	l.Entities = make(map[int64]entity.IEntity)

	if l.Provider != nil {
		l.Provider.Close()
	}

	l.Closed = true
	logger.Info("Level closed", "name", l.Name)
}

func (l *Level) GetLoadedChunks() []*world.Chunk {
	l.mu.RLock()
	defer l.mu.RUnlock()
	chunks := make([]*world.Chunk, 0, len(l.Chunks))
	for _, chunk := range l.Chunks {
		chunks = append(chunks, chunk)
	}
	return chunks
}

func (l *Level) GetLoadedChunkCount() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.Chunks)
}

func (l *Level) GetGenerator() generator.Generator {
	return l.Generator
}

func (l *Level) GetSpawnLocation() *world.Vector3 {
	if l.Generator != nil {
		return l.Generator.GetSpawn()
	}

	return world.NewVector3(128, 64, 128)
}

func (l *Level) GetSafeSpawn() *world.Vector3 {
	spawn := l.GetSpawnLocation()

	if l.Generator == nil {
		logger.Warn("GetSafeSpawn: generator is nil, using spawn location directly",
			"x", spawn.X, "y", spawn.Y, "z", spawn.Z)
		if spawn.Y < 64 {
			spawn.Y = 64
		}
		return spawn
	}

	chunkX := int32(spawn.X) >> 4
	chunkZ := int32(spawn.Z) >> 4
	chunk := l.GetChunk(chunkX, chunkZ, true)
	if chunk == nil {
		logger.Warn("GetSafeSpawn: chunk is nil, defaulting to sea level", "cx", chunkX, "cz", chunkZ)
		return world.NewVector3(spawn.X, 64, spawn.Z)
	}

	localX := int(int32(spawn.X) & 0x0F)
	localZ := int(int32(spawn.Z) & 0x0F)

	nonEmptySections := 0
	for i := 0; i < 8; i++ {
		if chunk.Sections[i] != nil && !chunk.Sections[i].IsEmpty() {
			nonEmptySections++
		}
	}
	logger.Info("GetSafeSpawn: chunk info",
		"cx", chunkX, "cz", chunkZ,
		"localX", localX, "localZ", localZ,
		"nonEmptySections", nonEmptySections,
		"generated", chunk.Generated,
		"populated", chunk.Populated,
	)

	if nonEmptySections == 0 {
		logger.Warn("GetSafeSpawn: chunk has no terrain data, defaulting to sea level")
		return world.NewVector3(spawn.X, 64, spawn.Z)
	}

	y := 127
	for y >= 0 {
		id := chunk.GetBlockId(localX, y, localZ)
		if block.Registry.IsSolid(id) {
			logger.Info("GetSafeSpawn: found solid",
				"y", y, "blockId", id, "safeY", y+1)
			break
		}
		y--
	}

	safeY := y + 1
	if safeY < 2 {
		logger.Warn("GetSafeSpawn: no solid block found above bedrock, defaulting to sea level")
		safeY = 64
	}

	logger.Info("GetSafeSpawn: result",
		"spawnX", spawn.X, "spawnY", float64(safeY), "spawnZ", spawn.Z)
	return world.NewVector3(spawn.X, float64(safeY), spawn.Z)
}

// loadTilesFromChunk 从区块的 NBT 数据中恢复 Tile 实体
// 对应 PHP Level::loadChunk() 中遍历 TileEntities 并 Tile::createTile() 的逻辑
func (l *Level) loadTilesFromChunk(chunk *world.Chunk) {
	if len(chunk.Tiles) == 0 {
		return
	}

	loaded := 0
	for _, tileNBT := range chunk.Tiles {
		tileID := tileNBT.GetString("id")
		if tileID == "" {
			continue
		}

		t := tile.CreateTile(tileID, chunk, tileNBT)
		if t == nil {
			logger.Debug("Unknown tile entity type, skipping", "id", tileID,
				"x", tileNBT.GetInt("x"), "y", tileNBT.GetInt("y"), "z", tileNBT.GetInt("z"))
			continue
		}

		l.Tiles.AddTile(t)

		// 如果 Tile 需要 tick 更新（如熔炉、刷怪笼），注册到更新队列
		if t.OnUpdate() {
			l.Tiles.ScheduleUpdate(t)
		}

		loaded++
	}

	if loaded > 0 {
		logger.Debug("Loaded tile entities from chunk",
			"cx", chunk.X, "cz", chunk.Z, "count", loaded)
	}
}

// SendChunkTiles 向玩家发送指定区块中所有 Spawnable Tile 的数据
// 在发送 FullChunkDataPacket 之后调用，让客户端渲染告示牌文字、箱子方向等
func (l *Level) SendChunkTiles(chunk *world.Chunk, sender tile.PacketSender) {
	tiles := l.Tiles.GetAllTiles()
	for _, t := range tiles {
		x, y, z := t.GetPosition()
		tileChunkX := x >> 4
		tileChunkZ := z >> 4
		_ = y

		if tileChunkX == chunk.X && tileChunkZ == chunk.Z {
			if spawnable, ok := t.(tile.Spawnable); ok {
				tile.SpawnTo(spawnable, sender)
			}
		}
	}
}
