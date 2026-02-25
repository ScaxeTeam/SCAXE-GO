package tile

// hopper.go — 漏斗 TileEntity
// 对应 PHP Hopper.php
// Spawnable + Container(5格) + Nameable + OnUpdate(物品传输tick)
//
// 槽位布局: 0-4 = 5个物品槽

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	// HopperSlots 漏斗槽位数
	HopperSlots = 5

	// HopperCooldownTicks 漏斗传输冷却 (每8tick传输一次)
	HopperCooldownTicks = 8
)

// Hopper 漏斗 TileEntity
type Hopper struct {
	SpawnableBase
	ContainerBase
	NameableBase

	Cooldown   int
	NeedUpdate bool
}

// NewHopper 创建漏斗实例
func NewHopper(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
	h := &Hopper{}
	InitSpawnableBase(&h.SpawnableBase, TypeHopper, chunk, nbtData)
	InitContainerBase(&h.ContainerBase, HopperSlots)
	h.NameableBase.LoadNameFromNBT(nbtData)

	// 从 NBT 加载物品
	h.ContainerBase.LoadItemsFromNBT(nbtData)

	// 加载冷却时间
	h.Cooldown = int(nbtData.GetInt("TransferCooldown"))
	h.NeedUpdate = true

	return h
}

func (h *Hopper) GetName() string {
	if h.HasCustomName() {
		return h.GetCustomName()
	}
	return "Hopper"
}

// OnUpdate 漏斗 tick 逻辑
// 对应 PHP Hopper::onUpdate()
// 完整的物品推拉逻辑需要 Level.GetTile + InventoryHolder 接口
func (h *Hopper) OnUpdate() bool {
	if h.IsClosed() {
		return false
	}

	// 冷却计时
	if h.Cooldown > 0 {
		h.Cooldown--
		return true
	}

	// TODO: 完整的漏斗物品传输逻辑
	// 1. 拾取上方实体物品 (pickupArea)
	// 2. 从上方容器 (tileUp) 拉取物品
	// 3. 向下方容器 (tileDown) 推送物品
	// 以上需要 Level.GetTileAt() 和 方向判断，Phase 3 实现

	h.Cooldown = HopperCooldownTicks
	return true
}

// ResetCooldown 重置传输冷却
func (h *Hopper) ResetCooldown(ticks int) {
	h.Cooldown = ticks
}

// GetSpawnCompound 客户端渲染数据
func (h *Hopper) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeHopper))
	compound.Set(nbt.NewIntTag("x", h.X))
	compound.Set(nbt.NewIntTag("y", h.Y))
	compound.Set(nbt.NewIntTag("z", h.Z))

	if h.HasCustomName() {
		compound.Set(nbt.NewStringTag("CustomName", h.GetCustomName()))
	}
	return compound
}

func (h *Hopper) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return DefaultUpdateCompoundTag(nbtData)
}

func (h *Hopper) SpawnTo(sender PacketSender) bool {
	return SpawnTo(h, sender)
}

func (h *Hopper) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(h, broadcaster)
}

func (h *Hopper) SaveNBT() {
	h.SpawnableBase.SaveNBT()
	h.ContainerBase.SaveItemsToNBT(h.NBT)
	h.NameableBase.SaveNameToNBT(h.NBT)
	h.NBT.Set(nbt.NewIntTag("TransferCooldown", int32(h.Cooldown)))
}

func init() {
	RegisterTile(TypeHopper, NewHopper)
}
