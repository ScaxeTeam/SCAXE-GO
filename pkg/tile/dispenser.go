package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	DispenserSlots = 9
)

type Dispenser struct {
	SpawnableBase
	ContainerBase
	NameableBase
}

func NewDispenser(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
	d := &Dispenser{}
	InitSpawnableBase(&d.SpawnableBase, TypeDispenser, chunk, nbtData)
	InitContainerBase(&d.ContainerBase, DispenserSlots)
	d.NameableBase.LoadNameFromNBT(nbtData)

	d.ContainerBase.LoadItemsFromNBT(nbtData)

	return d
}

func (d *Dispenser) GetName() string {
	if d.HasCustomName() {
		return d.GetCustomName()
	}
	return "Dispenser"
}

func (d *Dispenser) Activate() {

}

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
