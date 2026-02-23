package tile

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	TypeSign         = "Sign"
	TypeChest        = "Chest"
	TypeFurnace      = "Furnace"
	TypeFlowerPot    = "FlowerPot"
	TypeMobSpawner   = "MobSpawner"
	TypeSkull        = "Skull"
	TypeBrewingStand = "BrewingStand"
	TypeEnchantTable = "EnchantTable"
	TypeItemFrame    = "ItemFrame"
	TypeDispenser    = "Dispenser"
	TypeDropper      = "Dropper"
	TypeDLDetector   = "DLDetector"
	TypeCauldron     = "Cauldron"
	TypeHopper       = "Hopper"
	TypeComparator   = "Comparator"
	TypeNoteblock    = "Noteblock"
)

var tileCounter int64

func nextTileID() int64 {
	return atomic.AddInt64(&tileCounter, 1)
}

type TileFactory func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile

var (
	registryMu sync.RWMutex
	knownTiles = map[string]TileFactory{}
	shortNames = map[string]string{}
)

func RegisterTile(typeName string, factory TileFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()
	knownTiles[typeName] = factory
	shortNames[typeName] = typeName
}

func CreateTile(typeName string, chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
	registryMu.RLock()
	factory, ok := knownTiles[typeName]
	registryMu.RUnlock()
	if !ok {
		return nil
	}
	return factory(chunk, nbtData)
}

type Tile interface {
	GetID() int64

	GetSaveID() string

	GetName() string

	GetPosition() (x, y, z int32)

	GetChunk() *world.Chunk

	GetNBT() *nbt.CompoundTag

	SaveNBT()

	OnUpdate() bool

	IsClosed() bool

	Close()
}

type BaseTile struct {
	id     int64
	saveID string
	name   string

	X, Y, Z int32

	Chunk *world.Chunk
	NBT   *nbt.CompoundTag

	closed bool
}

func InitBaseTile(t *BaseTile, saveID string, chunk *world.Chunk, nbtData *nbt.CompoundTag) {
	t.id = nextTileID()
	t.saveID = saveID
	t.name = ""
	t.Chunk = chunk
	t.NBT = nbtData

	t.X = nbtData.GetInt("x")
	t.Y = nbtData.GetInt("y")
	t.Z = nbtData.GetInt("z")
}

func (t *BaseTile) GetID() int64 {
	return t.id
}

func (t *BaseTile) GetSaveID() string {
	return t.saveID
}

func (t *BaseTile) GetName() string {
	return t.name
}

func (t *BaseTile) GetPosition() (x, y, z int32) {
	return t.X, t.Y, t.Z
}

func (t *BaseTile) GetChunk() *world.Chunk {
	return t.Chunk
}

func (t *BaseTile) GetNBT() *nbt.CompoundTag {
	return t.NBT
}

func (t *BaseTile) SaveNBT() {
	t.NBT.Set(nbt.NewStringTag("id", t.saveID))
	t.NBT.Set(nbt.NewIntTag("x", t.X))
	t.NBT.Set(nbt.NewIntTag("y", t.Y))
	t.NBT.Set(nbt.NewIntTag("z", t.Z))
}

func (t *BaseTile) OnUpdate() bool {
	return false
}

func (t *BaseTile) IsClosed() bool {
	return t.closed
}

func (t *BaseTile) Close() {
	if t.closed {
		return
	}
	t.closed = true
}

func (t *BaseTile) String() string {
	return fmt.Sprintf("Tile(%s, id=%d, pos=[%d,%d,%d])", t.saveID, t.id, t.X, t.Y, t.Z)
}

type TileManager struct {
	mu          sync.RWMutex
	tiles       map[int64]Tile
	tilesByPos  map[int64]Tile
	updateTiles map[int64]Tile
}

func NewTileManager() *TileManager {
	return &TileManager{
		tiles:       make(map[int64]Tile),
		tilesByPos:  make(map[int64]Tile),
		updateTiles: make(map[int64]Tile),
	}
}

func posHash(x, y, z int32) int64 {
	return (int64(x) & 0x3FFFFFF) | ((int64(z) & 0x3FFFFFF) << 26) | (int64(y) << 52)
}

func (m *TileManager) AddTile(t Tile) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tiles[t.GetID()] = t
	x, y, z := t.GetPosition()
	m.tilesByPos[posHash(x, y, z)] = t
}

func (m *TileManager) RemoveTile(t Tile) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.tiles, t.GetID())
	x, y, z := t.GetPosition()
	delete(m.tilesByPos, posHash(x, y, z))
	delete(m.updateTiles, t.GetID())
}

func (m *TileManager) GetTileAt(x, y, z int32) Tile {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tilesByPos[posHash(x, y, z)]
}

func (m *TileManager) GetTileByID(id int64) Tile {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tiles[id]
}

func (m *TileManager) ScheduleUpdate(t Tile) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateTiles[t.GetID()] = t
}

func (m *TileManager) TickUpdates() {
	m.mu.Lock()

	toUpdate := make([]Tile, 0, len(m.updateTiles))
	for _, t := range m.updateTiles {
		toUpdate = append(toUpdate, t)
	}
	m.mu.Unlock()

	for _, t := range toUpdate {
		if t.IsClosed() {
			m.mu.Lock()
			delete(m.updateTiles, t.GetID())
			m.mu.Unlock()
			continue
		}
		if !t.OnUpdate() {
			m.mu.Lock()
			delete(m.updateTiles, t.GetID())
			m.mu.Unlock()
		}
	}
}

func (m *TileManager) GetAllTiles() []Tile {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]Tile, 0, len(m.tiles))
	for _, t := range m.tiles {
		result = append(result, t)
	}
	return result
}

func (m *TileManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.tiles)
}
