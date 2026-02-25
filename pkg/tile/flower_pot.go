package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)
var flowerPotAllowed = map[int16]bool{
	31: true,
	6:  true,
	32: true,
	37: true,
	38: true,
	39: true,
	40: true,
	81: true,
}
type FlowerPot struct {
	SpawnableBase
	PlantID   int16
	PlantData int32
}
func NewFlowerPot(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
	fp := &FlowerPot{}
	InitSpawnableBase(&fp.SpawnableBase, TypeFlowerPot, chunk, nbtData)

	fp.PlantID = nbtData.GetShort("item")
	fp.PlantData = nbtData.GetInt("mData")

	return fp
}

func (fp *FlowerPot) GetName() string {
	return "Flower Pot"
}
func (fp *FlowerPot) GetSpawnCompound() *nbt.CompoundTag {
	plantID := fp.PlantID
	if _, ok := flowerPotAllowed[plantID]; !ok {
		plantID = 0
	}

	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeFlowerPot))
	compound.Set(nbt.NewIntTag("x", fp.X))
	compound.Set(nbt.NewIntTag("y", fp.Y))
	compound.Set(nbt.NewIntTag("z", fp.Z))
	compound.Set(nbt.NewShortTag("item", plantID))
	compound.Set(nbt.NewIntTag("mData", fp.PlantData))
	return compound
}

func (fp *FlowerPot) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return DefaultUpdateCompoundTag(nbtData)
}

func (fp *FlowerPot) SpawnTo(sender PacketSender) bool {
	return SpawnTo(fp, sender)
}

func (fp *FlowerPot) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(fp, broadcaster)
}
func (fp *FlowerPot) SaveNBT() {
	fp.SpawnableBase.SaveNBT()
	fp.NBT.Set(nbt.NewShortTag("item", fp.PlantID))
	fp.NBT.Set(nbt.NewIntTag("mData", fp.PlantData))
}
func (fp *FlowerPot) SetFlowerPotData(itemID int16, data int32) {
	fp.PlantID = itemID
	fp.PlantData = data
}

func init() {
	RegisterTile(TypeFlowerPot, NewFlowerPot)
}
