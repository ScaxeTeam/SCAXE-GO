package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)
type EnchantTable struct {
	SpawnableBase
	NameableBase
}
func NewEnchantTable(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
	et := &EnchantTable{}
	InitSpawnableBase(&et.SpawnableBase, TypeEnchantTable, chunk, nbtData)
	et.NameableBase.LoadNameFromNBT(nbtData)
	return et
}

func (et *EnchantTable) GetName() string {
	if et.HasCustomName() {
		return et.GetCustomName()
	}
	return "Enchanting Table"
}
func (et *EnchantTable) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeEnchantTable))
	compound.Set(nbt.NewIntTag("x", et.X))
	compound.Set(nbt.NewIntTag("y", et.Y))
	compound.Set(nbt.NewIntTag("z", et.Z))

	if et.HasCustomName() {
		compound.Set(nbt.NewStringTag("CustomName", et.GetCustomName()))
	}
	return compound
}

func (et *EnchantTable) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return DefaultUpdateCompoundTag(nbtData)
}

func (et *EnchantTable) SpawnTo(sender PacketSender) bool {
	return SpawnTo(et, sender)
}

func (et *EnchantTable) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(et, broadcaster)
}

func (et *EnchantTable) SaveNBT() {
	et.SpawnableBase.SaveNBT()
	et.NameableBase.SaveNameToNBT(et.NBT)
}

func init() {
	RegisterTile(TypeEnchantTable, NewEnchantTable)
}
