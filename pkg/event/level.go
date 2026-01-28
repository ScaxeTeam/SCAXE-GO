package event

type LevelEvent struct {
	*BaseEvent
	LevelName string
}

func NewLevelEvent(name, levelName string) *LevelEvent {
	return &LevelEvent{
		BaseEvent: NewBaseEvent(name),
		LevelName: levelName,
	}
}

func (e *LevelEvent) GetLevelName() string {
	return e.LevelName
}

type LevelLoadEvent struct {
	*LevelEvent
}

var levelLoadHandlers = NewHandlerList()

func NewLevelLoadEvent(levelName string) *LevelLoadEvent {
	return &LevelLoadEvent{
		LevelEvent: NewLevelEvent("LevelLoadEvent", levelName),
	}
}

func (e *LevelLoadEvent) GetHandlers() *HandlerList {
	return levelLoadHandlers
}

type LevelUnloadEvent struct {
	*LevelEvent
}

var levelUnloadHandlers = NewHandlerList()

func NewLevelUnloadEvent(levelName string) *LevelUnloadEvent {
	return &LevelUnloadEvent{
		LevelEvent: NewLevelEvent("LevelUnloadEvent", levelName),
	}
}

func (e *LevelUnloadEvent) GetHandlers() *HandlerList {
	return levelUnloadHandlers
}

type ChunkLoadEvent struct {
	*LevelEvent
	ChunkX int
	ChunkZ int
	IsNew  bool
}

var chunkLoadHandlers = NewHandlerList()

func NewChunkLoadEvent(levelName string, chunkX, chunkZ int, isNew bool) *ChunkLoadEvent {
	return &ChunkLoadEvent{
		LevelEvent: NewLevelEvent("ChunkLoadEvent", levelName),
		ChunkX:     chunkX,
		ChunkZ:     chunkZ,
		IsNew:      isNew,
	}
}

func (e *ChunkLoadEvent) GetHandlers() *HandlerList {
	return chunkLoadHandlers
}

type ChunkUnloadEvent struct {
	*LevelEvent
	ChunkX int
	ChunkZ int
}

var chunkUnloadHandlers = NewHandlerList()

func NewChunkUnloadEvent(levelName string, chunkX, chunkZ int) *ChunkUnloadEvent {
	return &ChunkUnloadEvent{
		LevelEvent: NewLevelEvent("ChunkUnloadEvent", levelName),
		ChunkX:     chunkX,
		ChunkZ:     chunkZ,
	}
}

func (e *ChunkUnloadEvent) GetHandlers() *HandlerList {
	return chunkUnloadHandlers
}

type ChunkPopulateEvent struct {
	*LevelEvent
	ChunkX int
	ChunkZ int
}

var chunkPopulateHandlers = NewHandlerList()

func NewChunkPopulateEvent(levelName string, chunkX, chunkZ int) *ChunkPopulateEvent {
	return &ChunkPopulateEvent{
		LevelEvent: NewLevelEvent("ChunkPopulateEvent", levelName),
		ChunkX:     chunkX,
		ChunkZ:     chunkZ,
	}
}

func (e *ChunkPopulateEvent) GetHandlers() *HandlerList {
	return chunkPopulateHandlers
}

type SpawnChangeEvent struct {
	*LevelEvent
	OldX, OldY, OldZ float64
	NewX, NewY, NewZ float64
}

var spawnChangeHandlers = NewHandlerList()

func NewSpawnChangeEvent(levelName string, ox, oy, oz, nx, ny, nz float64) *SpawnChangeEvent {
	return &SpawnChangeEvent{
		LevelEvent: NewLevelEvent("SpawnChangeEvent", levelName),
		OldX:       ox, OldY: oy, OldZ: oz,
		NewX: nx, NewY: ny, NewZ: nz,
	}
}

func (e *SpawnChangeEvent) GetHandlers() *HandlerList {
	return spawnChangeHandlers
}
