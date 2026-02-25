package tile

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

// Tile 类型常量（对应 PHP Tile::SIGN 等）
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

// tileCounter 全局 Tile ID 计数器（原子操作，线程安全）
var tileCounter int64

func nextTileID() int64 {
	return atomic.AddInt64(&tileCounter, 1)
}

// ---------- 注册表 ----------

// TileFactory 是创建 Tile 实例的工厂函数
type TileFactory func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile

var (
	registryMu sync.RWMutex
	knownTiles = map[string]TileFactory{} // typeName -> factory
	shortNames = map[string]string{}      // 实现类名 -> shortName（Go 中用注册名）
)

// RegisterTile 注册一个 Tile 类型的工厂函数
func RegisterTile(typeName string, factory TileFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()
	knownTiles[typeName] = factory
	shortNames[typeName] = typeName
}

// CreateTile 通过类型名创建 Tile 实例
func CreateTile(typeName string, chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
	registryMu.RLock()
	factory, ok := knownTiles[typeName]
	registryMu.RUnlock()
	if !ok {
		return nil
	}
	return factory(chunk, nbtData)
}

// ---------- Tile 接口 ----------

// Tile 接口定义了所有 TileEntity 必须实现的方法
type Tile interface {
	// GetID 返回此 Tile 的唯一ID
	GetID() int64

	// GetSaveID 返回保存用的类型名（如 "Chest", "Sign"）
	GetSaveID() string

	// GetName 返回此 Tile 的名称
	GetName() string

	// GetPosition 返回 Tile 的世界坐标
	GetPosition() (x, y, z int32)

	// GetChunk 返回此 Tile 所在的区块
	GetChunk() *world.Chunk

	// GetNBT 返回此 Tile 的 NBT 数据
	GetNBT() *nbt.CompoundTag

	// SaveNBT 将当前状态写入 NBT
	SaveNBT()

	// OnUpdate 每 tick 更新，返回 true 表示需要继续更新
	OnUpdate() bool

	// IsClosed 返回此 Tile 是否已关闭
	IsClosed() bool

	// Close 关闭并移除此 Tile
	Close()
}

// ---------- BaseTile 基类实现 ----------

// BaseTile 是所有 Tile 的基类，对应 PHP abstract class Tile extends Position
type BaseTile struct {
	id     int64
	saveID string // 注册的类型名
	name   string

	X, Y, Z int32

	Chunk *world.Chunk
	NBT   *nbt.CompoundTag

	closed bool
}

// InitBaseTile 初始化 BaseTile 的字段，应在具体 Tile 的构造函数中调用
// 对应 PHP Tile::__construct(FullChunk $chunk, CompoundTag $nbt)
func InitBaseTile(t *BaseTile, saveID string, chunk *world.Chunk, nbtData *nbt.CompoundTag) {
	t.id = nextTileID()
	t.saveID = saveID
	t.name = ""
	t.Chunk = chunk
	t.NBT = nbtData

	// 从 NBT 读取坐标（对应 PHP: $this->x = (int) $this->namedtag["x"]）
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

// SaveNBT 将坐标和类型ID写入 NBT
// 对应 PHP: Tile::saveNBT()
func (t *BaseTile) SaveNBT() {
	t.NBT.Set(nbt.NewStringTag("id", t.saveID))
	t.NBT.Set(nbt.NewIntTag("x", t.X))
	t.NBT.Set(nbt.NewIntTag("y", t.Y))
	t.NBT.Set(nbt.NewIntTag("z", t.Z))
}

// OnUpdate 默认不需要更新
// 对应 PHP: Tile::onUpdate() { return false; }
func (t *BaseTile) OnUpdate() bool {
	return false
}

func (t *BaseTile) IsClosed() bool {
	return t.closed
}

// Close 关闭此 Tile
// 对应 PHP: Tile::close()
func (t *BaseTile) Close() {
	if t.closed {
		return
	}
	t.closed = true
}

func (t *BaseTile) String() string {
	return fmt.Sprintf("Tile(%s, id=%d, pos=[%d,%d,%d])", t.saveID, t.id, t.X, t.Y, t.Z)
}

// ---------- TileManager 管理器（集成到 Level 使用） ----------

// TileManager 管理一个世界中所有活跃的 Tile
type TileManager struct {
	mu          sync.RWMutex
	tiles       map[int64]Tile // id -> Tile
	tilesByPos  map[int64]Tile // posHash -> Tile
	updateTiles map[int64]Tile // 需要 tick 更新的 Tile
}

// NewTileManager 创建一个新的 TileManager
func NewTileManager() *TileManager {
	return &TileManager{
		tiles:       make(map[int64]Tile),
		tilesByPos:  make(map[int64]Tile),
		updateTiles: make(map[int64]Tile),
	}
}

// posHash 计算位置哈希
func posHash(x, y, z int32) int64 {
	return (int64(x) & 0x3FFFFFF) | ((int64(z) & 0x3FFFFFF) << 26) | (int64(y) << 52)
}

// AddTile 添加一个 Tile 到管理器
func (m *TileManager) AddTile(t Tile) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tiles[t.GetID()] = t
	x, y, z := t.GetPosition()
	m.tilesByPos[posHash(x, y, z)] = t
}

// RemoveTile 从管理器中移除一个 Tile
func (m *TileManager) RemoveTile(t Tile) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.tiles, t.GetID())
	x, y, z := t.GetPosition()
	delete(m.tilesByPos, posHash(x, y, z))
	delete(m.updateTiles, t.GetID())
}

// GetTileAt 获取指定位置的 Tile
func (m *TileManager) GetTileAt(x, y, z int32) Tile {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tilesByPos[posHash(x, y, z)]
}

// GetTileByID 通过 ID 获取 Tile
func (m *TileManager) GetTileByID(id int64) Tile {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tiles[id]
}

// ScheduleUpdate 标记一个 Tile 需要每 tick 更新
// 对应 PHP: Tile::scheduleUpdate() { $this->level->updateTiles[$this->id] = $this; }
func (m *TileManager) ScheduleUpdate(t Tile) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateTiles[t.GetID()] = t
}

// TickUpdates 处理所有需要更新的 Tile
func (m *TileManager) TickUpdates() {
	m.mu.Lock()
	// 复制一份避免在迭代中修改
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

// GetAllTiles 返回所有活跃的 Tile
func (m *TileManager) GetAllTiles() []Tile {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]Tile, 0, len(m.tiles))
	for _, t := range m.tiles {
		result = append(result, t)
	}
	return result
}

// Count 返回 Tile 数量
func (m *TileManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.tiles)
}
