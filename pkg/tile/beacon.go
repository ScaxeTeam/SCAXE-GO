package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	TypeBeacon = "Beacon"
)
type Beacon struct {
	SpawnableBase
	NameableBase
	Primary   int32
	Secondary int32
}
func NewBeacon(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
	b := &Beacon{}
	InitSpawnableBase(&b.SpawnableBase, TypeBeacon, chunk, nbtData)
	b.NameableBase.LoadNameFromNBT(nbtData)

	b.Primary = nbtData.GetInt("primary")
	b.Secondary = nbtData.GetInt("secondary")

	return b
}

func (b *Beacon) GetName() string {
	if b.HasCustomName() {
		return b.GetCustomName()
	}
	return "Beacon"
}
func (b *Beacon) OnUpdate() bool {
	return false
}
func (b *Beacon) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeBeacon))
	compound.Set(nbt.NewIntTag("x", b.X))
	compound.Set(nbt.NewIntTag("y", b.Y))
	compound.Set(nbt.NewIntTag("z", b.Z))
	compound.Set(nbt.NewIntTag("primary", b.Primary))
	compound.Set(nbt.NewIntTag("secondary", b.Secondary))

	if b.HasCustomName() {
		compound.Set(nbt.NewStringTag("CustomName", b.GetCustomName()))
	}
	return compound
}

func (b *Beacon) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	b.Primary = nbtData.GetInt("primary")
	b.Secondary = nbtData.GetInt("secondary")
	return true
}

func (b *Beacon) SpawnTo(sender PacketSender) bool {
	return SpawnTo(b, sender)
}

func (b *Beacon) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(b, broadcaster)
}

func (b *Beacon) SaveNBT() {
	b.SpawnableBase.SaveNBT()
	b.NameableBase.SaveNameToNBT(b.NBT)
	b.NBT.Set(nbt.NewIntTag("primary", b.Primary))
	b.NBT.Set(nbt.NewIntTag("secondary", b.Secondary))
}

func init() {
	RegisterTile(TypeBeacon, NewBeacon)
}
