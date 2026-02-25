package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	HopperSlots = 5
	HopperCooldownTicks = 8
)
type Hopper struct {
	SpawnableBase
	ContainerBase
	NameableBase

	Cooldown   int
	NeedUpdate bool
}
func NewHopper(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
	h := &Hopper{}
	InitSpawnableBase(&h.SpawnableBase, TypeHopper, chunk, nbtData)
	InitContainerBase(&h.ContainerBase, HopperSlots)
	h.NameableBase.LoadNameFromNBT(nbtData)
	h.ContainerBase.LoadItemsFromNBT(nbtData)
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
func (h *Hopper) OnUpdate() bool {
	if h.IsClosed() {
		return false
	}
	if h.Cooldown > 0 {
		h.Cooldown--
		return true
	}

	h.Cooldown = HopperCooldownTicks
	return true
}
func (h *Hopper) ResetCooldown(ticks int) {
	h.Cooldown = ticks
}
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
