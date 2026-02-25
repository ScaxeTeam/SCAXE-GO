package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)
type Skull struct {
	SpawnableBase
}
func NewSkull(chunk *world.Chunk, nbtData *nbt.CompoundTag) *Skull {
	s := &Skull{}

	if nbtData.Get("SkullType") == nil {
		nbtData.Set(nbt.NewByteTag("SkullType", 0))
	}
	if nbtData.Get("Rot") == nil {
		nbtData.Set(nbt.NewByteTag("Rot", 0))
	}

	InitSpawnableBase(&s.SpawnableBase, TypeSkull, chunk, nbtData)
	return s
}
func (s *Skull) GetName() string {
	return "Skull"
}
func (s *Skull) GetSkullType() int8 {
	return s.NBT.GetByte("SkullType")
}
func (s *Skull) SetSkullType(skullType int8) {
	s.NBT.Set(nbt.NewByteTag("SkullType", skullType))
}
func (s *Skull) GetRot() int8 {
	return s.NBT.GetByte("Rot")
}
func (s *Skull) SaveNBT() {
	s.BaseTile.SaveNBT()
	s.NBT.Remove("Creator")
}
func (s *Skull) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeSkull))
	compound.Set(nbt.NewIntTag("x", s.X))
	compound.Set(nbt.NewIntTag("y", s.Y))
	compound.Set(nbt.NewIntTag("z", s.Z))
	compound.Set(nbt.NewByteTag("SkullType", s.GetSkullType()))
	compound.Set(nbt.NewByteTag("Rot", s.GetRot()))
	return compound
}
func (s *Skull) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return false
}
func (s *Skull) SpawnTo(sender PacketSender) bool {
	return SpawnTo(s, sender)
}
func (s *Skull) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(s, broadcaster)
}

func init() {
	RegisterTile(TypeSkull, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewSkull(chunk, nbtData)
	})
}
