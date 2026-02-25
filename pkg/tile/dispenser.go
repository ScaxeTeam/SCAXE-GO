package tile

// dispenser.go — 发射器 TileEntity
// 对应 PHP Dispenser.php
// Spawnable + Container(9格) + Nameable
// activate() 发射物品需要 Entity 创建系统，此处实现结构框架

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	// DispenserSlots 发射器槽位数
	DispenserSlots = 9
)

// Dispenser 发射器 TileEntity
type Dispenser struct {
	SpawnableBase
	ContainerBase
	NameableBase
}

// NewDispenser 创建发射器实例
func NewDispenser(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
	d := &Dispenser{}
	InitSpawnableBase(&d.SpawnableBase, TypeDispenser, chunk, nbtData)
	InitContainerBase(&d.ContainerBase, DispenserSlots)
	d.NameableBase.LoadNameFromNBT(nbtData)

	// 从 NBT 加载物品
	d.ContainerBase.LoadItemsFromNBT(nbtData)

	return d
}

func (d *Dispenser) GetName() string {
	if d.HasCustomName() {
		return d.GetCustomName()
	}
	return "Dispenser"
}

// Activate 发射物品
// 对应 PHP Dispenser::activate()
// 完整实现需要 Entity.CreateEntity + Level.DropItem，此处预留接口
func (d *Dispenser) Activate() {
	// TODO: 随机选择一个非空槽位
	// TODO: 根据物品类型创建对应投射物 (Arrow/Snowball/Egg/SplashPotion/ExpBottle)
	// TODO: 或者将物品作为掉落物弹出
	// TODO: 播放烟雾粒子效果
}

// GetSpawnCompound 客户端渲染数据
func (d *Dispenser) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeDispenser))
	compound.Set(nbt.NewIntTag("x", d.X))
	compound.Set(nbt.NewIntTag("y", d.Y))
	compound.Set(nbt.NewIntTag("z", d.Z))

	if d.HasCustomName() {
		compound.Set(nbt.NewStringTag("CustomName", d.GetCustomName()))
	}
	return compound
}

func (d *Dispenser) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return DefaultUpdateCompoundTag(nbtData)
}

func (d *Dispenser) SpawnTo(sender PacketSender) bool {
	return SpawnTo(d, sender)
}

func (d *Dispenser) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(d, broadcaster)
}

func (d *Dispenser) SaveNBT() {
	d.SpawnableBase.SaveNBT()
	d.ContainerBase.SaveItemsToNBT(d.NBT)
	d.NameableBase.SaveNameToNBT(d.NBT)
}

func init() {
	RegisterTile(TypeDispenser, NewDispenser)
}
