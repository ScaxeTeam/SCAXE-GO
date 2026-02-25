package tile

import (
	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)
const ChestSize = 27
type Chest struct {
	SpawnableBase
	ContainerBase
	NameableBase
	pairX int32
	pairZ int32
}
func NewChest(chunk *world.Chunk, nbtData *nbt.CompoundTag) *Chest {
	c := &Chest{
		pairX: -1,
		pairZ: -1,
	}

	InitSpawnableBase(&c.SpawnableBase, TypeChest, chunk, nbtData)
	InitContainerBase(&c.ContainerBase, ChestSize)
	c.ContainerBase.LoadItemsFromNBT(nbtData)
	c.NameableBase.LoadNameFromNBT(nbtData)
	if nbtData.Has("pairx") && nbtData.Has("pairz") {
		c.pairX = nbtData.GetInt("pairx")
		c.pairZ = nbtData.GetInt("pairz")
	}

	return c
}
func (c *Chest) GetSize() int {
	return ChestSize
}
func (c *Chest) GetName() string {
	if c.HasCustomName() {
		return c.GetCustomName()
	}
	return "Chest"
}
func (c *Chest) IsPaired() bool {
	return c.pairX != -1 && c.pairZ != -1
}
func (c *Chest) GetPairPosition() (x, z int32) {
	return c.pairX, c.pairZ
}
func (c *Chest) PairWith(other *Chest) bool {
	if c.IsPaired() || other.IsPaired() {
		return false
	}

	c.createPair(other)
	return true
}
func (c *Chest) Unpair(getPairFunc func(x, z int32) *Chest) bool {
	if !c.IsPaired() {
		return false
	}
	pair := getPairFunc(c.pairX, c.pairZ)
	c.pairX = -1
	c.pairZ = -1
	c.NBT.Remove("pairx")
	c.NBT.Remove("pairz")
	if pair != nil {
		pair.pairX = -1
		pair.pairZ = -1
		pair.NBT.Remove("pairx")
		pair.NBT.Remove("pairz")
	}

	return true
}
func (c *Chest) createPair(other *Chest) {
	c.pairX = other.X
	c.pairZ = other.Z
	c.NBT.Set(nbt.NewIntTag("pairx", other.X))
	c.NBT.Set(nbt.NewIntTag("pairz", other.Z))

	other.pairX = c.X
	other.pairZ = c.Z
	other.NBT.Set(nbt.NewIntTag("pairx", c.X))
	other.NBT.Set(nbt.NewIntTag("pairz", c.Z))
}
func (c *Chest) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeChest))
	compound.Set(nbt.NewIntTag("x", c.X))
	compound.Set(nbt.NewIntTag("y", c.Y))
	compound.Set(nbt.NewIntTag("z", c.Z))
	if c.IsPaired() {
		compound.Set(nbt.NewIntTag("pairx", c.pairX))
		compound.Set(nbt.NewIntTag("pairz", c.pairZ))
	}
	if c.HasCustomName() {
		compound.Set(nbt.NewStringTag("CustomName", c.GetCustomName()))
	}

	return compound
}
func (c *Chest) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	if nbtData.GetString("id") != TypeChest {
		return false
	}
	return true
}
func (c *Chest) SpawnTo(sender PacketSender) bool {
	return SpawnTo(c, sender)
}
func (c *Chest) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(c, broadcaster)
}
func (c *Chest) SaveNBT() {
	c.BaseTile.SaveNBT()
	c.ContainerBase.SaveItemsToNBT(c.NBT)
	c.NameableBase.SaveNameToNBT(c.NBT)
}
func (c *Chest) Close() {
	if c.IsClosed() {
		return
	}
	c.ClearAll()
	c.BaseTile.Close()
}
func GetContentsForDoubleChest(left, right *Chest) []item.Item {
	result := make([]item.Item, ChestSize*2)
	for i := 0; i < ChestSize; i++ {
		result[i] = left.GetItem(i)
		result[ChestSize+i] = right.GetItem(i)
	}
	return result
}

func init() {
	RegisterTile(TypeChest, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewChest(chunk, nbtData)
	})
}
