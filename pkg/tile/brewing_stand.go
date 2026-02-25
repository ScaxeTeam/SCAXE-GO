package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	MaxBrewTime = 400
	BrewingSlots = 4
)
var brewingIngredients = map[int16]bool{
	372: true,
	348: true,
	331: true,
	376: true,
	378: true,
	353: true,
	382: true,
	375: true,
	370: true,
	377: true,
	396: true,
	462: true,
	414: true,
	289: true,
}
type BrewingStand struct {
	SpawnableBase
	ContainerBase
	NameableBase

	CookTime   int16
	NeedUpdate bool
}
func NewBrewingStand(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
	bs := &BrewingStand{}
	InitSpawnableBase(&bs.SpawnableBase, TypeBrewingStand, chunk, nbtData)
	InitContainerBase(&bs.ContainerBase, BrewingSlots)
	bs.NameableBase.LoadNameFromNBT(nbtData)
	bs.ContainerBase.LoadItemsFromNBT(nbtData)
	bs.CookTime = nbtData.GetShort("CookTime")
	if bs.CookTime == 0 {
		bs.CookTime = MaxBrewTime
	}

	bs.NeedUpdate = true

	return bs
}

func (bs *BrewingStand) GetName() string {
	if bs.HasCustomName() {
		return bs.GetCustomName()
	}
	return "Brewing Stand"
}
func (bs *BrewingStand) OnUpdate() bool {
	if bs.IsClosed() {
		return false
	}
	if !bs.NeedUpdate {
		return true
	}
	ingredient := bs.GetItem(0)
	canBrew := false
	for i := 1; i <= 3; i++ {
		it := bs.GetItem(i)
		if it.ID == 373 || it.ID == 438 {
			canBrew = true
			break
		}
	}
	if canBrew && ingredient.ID != 0 && ingredient.Count > 0 {
		if _, ok := brewingIngredients[int16(ingredient.ID)]; !ok {
			canBrew = false
		}
	} else {
		canBrew = false
	}

	if canBrew {
		bs.CookTime--
		if bs.CookTime <= 0 {
			bs.CookTime = MaxBrewTime
		}
	} else {
		bs.CookTime = MaxBrewTime
	}

	bs.NeedUpdate = canBrew
	return true
}
func (bs *BrewingStand) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeBrewingStand))
	compound.Set(nbt.NewIntTag("x", bs.X))
	compound.Set(nbt.NewIntTag("y", bs.Y))
	compound.Set(nbt.NewIntTag("z", bs.Z))
	compound.Set(nbt.NewShortTag("CookTime", MaxBrewTime))

	if bs.HasCustomName() {
		compound.Set(nbt.NewStringTag("CustomName", bs.GetCustomName()))
	}
	return compound
}

func (bs *BrewingStand) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return DefaultUpdateCompoundTag(nbtData)
}

func (bs *BrewingStand) SpawnTo(sender PacketSender) bool {
	return SpawnTo(bs, sender)
}

func (bs *BrewingStand) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(bs, broadcaster)
}

func (bs *BrewingStand) SaveNBT() {
	bs.SpawnableBase.SaveNBT()
	bs.ContainerBase.SaveItemsToNBT(bs.NBT)
	bs.NameableBase.SaveNameToNBT(bs.NBT)
	bs.NBT.Set(nbt.NewShortTag("CookTime", bs.CookTime))
}
func CheckIngredient(itemID int16) bool {
	_, ok := brewingIngredients[itemID]
	return ok
}

func init() {
	RegisterTile(TypeBrewingStand, NewBrewingStand)
}
