package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)
type Dropper struct {
	SpawnableBase
}
const DropperSize = 9
func NewDropper(chunk *world.Chunk, nbtData *nbt.CompoundTag) *Dropper {
	d := &Dropper{}
	if nbtData.Get("Items") == nil {
		nbtData.Set(nbt.NewListTag("Items", nbt.TagCompound))
	}

	InitSpawnableBase(&d.SpawnableBase, TypeDropper, chunk, nbtData)
	return d
}
func (d *Dropper) GetName() string {
	if customName := d.NBT.GetString("CustomName"); customName != "" {
		return customName
	}
	return "Dropper"
}
func (d *Dropper) HasName() bool {
	return d.NBT.Get("CustomName") != nil && d.NBT.GetString("CustomName") != ""
}
func (d *Dropper) SetName(name string) {
	if name == "" {
		d.NBT.Remove("CustomName")
		return
	}
	d.NBT.Set(nbt.NewStringTag("CustomName", name))
}
func (d *Dropper) GetSize() int {
	return DropperSize
}
func GetDropperMotion(meta int) (x, y, z int) {
	switch meta {
	case 0:
		return 0, -1, 0
	case 1:
		return 0, 1, 0
	case 2:
		return 0, 0, -1
	case 3:
		return 0, 0, 1
	case 4:
		return -1, 0, 0
	case 5:
		return 1, 0, 0
	default:
		return 0, 0, 0
	}
}
func (d *Dropper) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeDropper))
	compound.Set(nbt.NewIntTag("x", d.X))
	compound.Set(nbt.NewIntTag("y", d.Y))
	compound.Set(nbt.NewIntTag("z", d.Z))

	if d.HasName() {
		compound.Set(nbt.NewStringTag("CustomName", d.GetName()))
	}

	return compound
}
func (d *Dropper) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return false
}
func (d *Dropper) SpawnTo(sender PacketSender) bool {
	return SpawnTo(d, sender)
}
func (d *Dropper) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(d, broadcaster)
}

func init() {
	RegisterTile(TypeDropper, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewDropper(chunk, nbtData)
	})
}
