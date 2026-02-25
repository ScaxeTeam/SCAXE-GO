package tile

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

// FurnaceSize 熔炉槽位数
const FurnaceSize = 3

// 熔炉槽位索引
const (
	FurnaceSlotInput  = 0 // 原料槽
	FurnaceSlotFuel   = 1 // 燃料槽
	FurnaceSlotOutput = 2 // 产出槽
)

// Furnace 熔炉 TileEntity
// 对应 PHP class Furnace extends Spawnable implements InventoryHolder, Container, Nameable
// 实现 Spawnable + Container + Nameable，支持熔炼逻辑和 tick 更新
type Furnace struct {
	SpawnableBase
	ContainerBase
	NameableBase

	// 熔炼状态
	BurnTime  int16 // 剩余燃烧时间（tick）
	CookTime  int16 // 已烹饪时间（tick），到 200 完成一次熔炼
	MaxTime   int16 // 当前燃料最大燃烧时间（用于客户端进度条）
	BurnTicks int16 // 客户端显示用的燃烧进度

	// 是否需要 tick 更新
	needUpdate bool
}

// NewFurnace 创建 Furnace 实例
// 对应 PHP Furnace::__construct(FullChunk $chunk, CompoundTag $nbt)
func NewFurnace(chunk *world.Chunk, nbtData *nbt.CompoundTag) *Furnace {
	f := &Furnace{
		needUpdate: true,
	}

	// 初始化熔炼计时器（对应 PHP 构造函数中的默认值验证）
	burnTime := nbtData.GetShort("BurnTime")
	if burnTime < 0 {
		burnTime = 0
	}

	cookTime := nbtData.GetShort("CookTime")
	if cookTime < 0 || (burnTime == 0 && cookTime > 0) {
		cookTime = 0
	}

	maxTime := nbtData.GetShort("MaxTime")
	if maxTime == 0 {
		maxTime = burnTime
	}

	burnTicks := nbtData.GetShort("BurnTicks")

	f.BurnTime = burnTime
	f.CookTime = cookTime
	f.MaxTime = maxTime
	f.BurnTicks = burnTicks

	// 同步到 NBT
	nbtData.Set(nbt.NewShortTag("BurnTime", f.BurnTime))
	nbtData.Set(nbt.NewShortTag("CookTime", f.CookTime))
	nbtData.Set(nbt.NewShortTag("MaxTime", f.MaxTime))
	nbtData.Set(nbt.NewShortTag("BurnTicks", f.BurnTicks))

	InitSpawnableBase(&f.SpawnableBase, TypeFurnace, chunk, nbtData)
	InitContainerBase(&f.ContainerBase, FurnaceSize)

	// 从 NBT 加载物品
	f.ContainerBase.LoadItemsFromNBT(nbtData)

	// 从 NBT 加载自定义名称
	f.NameableBase.LoadNameFromNBT(nbtData)

	return f
}

// ---------- Container 接口 ----------

// GetSize 返回熔炉容量 (3)
// 对应 PHP Furnace::getSize() { return 3; }
func (f *Furnace) GetSize() int {
	return FurnaceSize
}

// ---------- 便捷方法（对应 FurnaceInventory 的方法） ----------

// GetSmelting 获取原料槽物品
func (f *Furnace) GetSmelting() item.Item {
	return f.GetItem(FurnaceSlotInput)
}

// GetFuel 获取燃料槽物品
func (f *Furnace) GetFuel() item.Item {
	return f.GetItem(FurnaceSlotFuel)
}

// GetResult 获取产出槽物品
func (f *Furnace) GetResult() item.Item {
	return f.GetItem(FurnaceSlotOutput)
}

// SetSmelting 设置原料槽物品
func (f *Furnace) SetSmelting(it item.Item) {
	f.SetItem(FurnaceSlotInput, it)
}

// SetFuel 设置燃料槽物品
func (f *Furnace) SetFuel(it item.Item) {
	f.SetItem(FurnaceSlotFuel, it)
}

// SetResult 设置产出槽物品
func (f *Furnace) SetResult(it item.Item) {
	f.SetItem(FurnaceSlotOutput, it)
}

// ---------- Nameable 接口实现 ----------

// GetName 返回熔炉名称
// 对应 PHP Furnace::getName() 优先返回 CustomName，否则返回 "Furnace"
func (f *Furnace) GetName() string {
	if f.HasCustomName() {
		return f.GetCustomName()
	}
	return "Furnace"
}

// ---------- 熔炼逻辑 ----------

// IsBurning 返回熔炉是否正在燃烧
func (f *Furnace) IsBurning() bool {
	return f.BurnTime > 0
}

// SetNeedUpdate 设置是否需要在下一 tick 更新
// 对应 PHP Furnace::setNeedUpdate(bool $needUpdate)
func (f *Furnace) SetNeedUpdate(need bool) {
	f.needUpdate = need
}

// OnUpdate 每 tick 更新熔炼状态
// 对应 PHP Furnace::onUpdate()
// 简化版：不包含事件系统和方块状态切换，由外部（Server/Level）提供熔炼配方和方块操作回调
// 返回 true 表示需要继续更新
func (f *Furnace) OnUpdate() bool {
	if f.IsClosed() {
		return false
	}

	if f.BurnTime > 0 {
		f.BurnTime--
		if f.MaxTime > 0 {
			f.BurnTicks = int16(math.Ceil(float64(f.BurnTime) / float64(f.MaxTime) * 200))
		}
	}

	// 同步到 NBT
	f.NBT.Set(nbt.NewShortTag("BurnTime", f.BurnTime))
	f.NBT.Set(nbt.NewShortTag("CookTime", f.CookTime))
	f.NBT.Set(nbt.NewShortTag("BurnTicks", f.BurnTicks))

	return true
}

// StartBurning 开始燃烧（消耗燃料）
// 对应 PHP Furnace::checkFuel() 的核心逻辑
// 参数 fuelTime 为燃料的燃烧时间（tick）
func (f *Furnace) StartBurning(fuelTime int16) {
	f.MaxTime = fuelTime
	f.BurnTime = fuelTime
	f.BurnTicks = 0

	f.NBT.Set(nbt.NewShortTag("MaxTime", f.MaxTime))
	f.NBT.Set(nbt.NewShortTag("BurnTime", f.BurnTime))
	f.NBT.Set(nbt.NewShortTag("BurnTicks", f.BurnTicks))
}

// IncrementCookTime 增加烹饪时间，返回是否完成一次熔炼（到 200）
// 对应 PHP 中 CookTime 达到 200 时完成熔炼的逻辑
func (f *Furnace) IncrementCookTime() bool {
	f.CookTime++
	if f.CookTime >= 200 {
		f.CookTime -= 200
		f.NBT.Set(nbt.NewShortTag("CookTime", f.CookTime))
		return true // 完成一次熔炼
	}
	f.NBT.Set(nbt.NewShortTag("CookTime", f.CookTime))
	return false
}

// ResetCookTime 重置烹饪时间
func (f *Furnace) ResetCookTime() {
	f.CookTime = 0
	f.NBT.Set(nbt.NewShortTag("CookTime", f.CookTime))
}

// ResetAll 重置所有计时器
func (f *Furnace) ResetAll() {
	f.BurnTime = 0
	f.CookTime = 0
	f.BurnTicks = 0
	f.NBT.Set(nbt.NewShortTag("BurnTime", int16(0)))
	f.NBT.Set(nbt.NewShortTag("CookTime", int16(0)))
	f.NBT.Set(nbt.NewShortTag("BurnTicks", int16(0)))
}

// ---------- Spawnable 接口实现 ----------

// GetSpawnCompound 返回发送给客户端的 NBT 数据
// 对应 PHP Furnace::getSpawnCompound()
func (f *Furnace) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeFurnace))
	compound.Set(nbt.NewIntTag("x", f.X))
	compound.Set(nbt.NewIntTag("y", f.Y))
	compound.Set(nbt.NewIntTag("z", f.Z))
	compound.Set(nbt.NewShortTag("BurnTime", f.BurnTime))
	compound.Set(nbt.NewShortTag("CookTime", f.CookTime))

	if f.HasCustomName() {
		compound.Set(nbt.NewStringTag("CustomName", f.GetCustomName()))
	}

	return compound
}

// UpdateCompoundTag 处理玩家发来的 NBT 更新
func (f *Furnace) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	if nbtData.GetString("id") != TypeFurnace {
		return false
	}
	return true
}

// SpawnTo 向指定玩家发送此熔炉的数据包
func (f *Furnace) SpawnTo(sender PacketSender) bool {
	return SpawnTo(f, sender)
}

// SpawnToAll 向区块内所有玩家广播
func (f *Furnace) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(f, broadcaster)
}

// ---------- Tile 接口补充 ----------

// SaveNBT 将当前状态写入 NBT
// 对应 PHP Furnace::saveNBT()
func (f *Furnace) SaveNBT() {
	f.BaseTile.SaveNBT()

	// 保存物品到 NBT
	f.ContainerBase.SaveItemsToNBT(f.NBT)

	// 保存自定义名称
	f.NameableBase.SaveNameToNBT(f.NBT)

	// 保存熔炼状态
	f.NBT.Set(nbt.NewShortTag("BurnTime", f.BurnTime))
	f.NBT.Set(nbt.NewShortTag("CookTime", f.CookTime))
	f.NBT.Set(nbt.NewShortTag("MaxTime", f.MaxTime))
	f.NBT.Set(nbt.NewShortTag("BurnTicks", f.BurnTicks))
}

// Close 关闭熔炉 Tile
// 对应 PHP Furnace::close()
func (f *Furnace) Close() {
	if f.IsClosed() {
		return
	}
	f.ClearAll()
	f.BaseTile.Close()
}

// ---------- 注册 ----------

func init() {
	RegisterTile(TypeFurnace, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewFurnace(chunk, nbtData)
	})
}
