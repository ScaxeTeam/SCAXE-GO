package tile

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)
const FurnaceSize = 3
const (
	FurnaceSlotInput  = 0
	FurnaceSlotFuel   = 1
	FurnaceSlotOutput = 2
)
type Furnace struct {
	SpawnableBase
	ContainerBase
	NameableBase
	BurnTime  int16
	CookTime  int16
	MaxTime   int16
	BurnTicks int16
	needUpdate bool
}
func NewFurnace(chunk *world.Chunk, nbtData *nbt.CompoundTag) *Furnace {
	f := &Furnace{
		needUpdate: true,
	}
	burnTime := nbtData.GetShort("BurnTime")
	if burnTime < 0 {
		burnTime = 0
	}

	cookTime := nbtData.GetShort("CookTime")
	if cookTime < 0 || (burnTime == 0 && cookTime > 0) {
		cookTime = 0
	}

	maxTime := nbtData.GetShort("MaxTime")
	if maxTime == 0 {
		maxTime = burnTime
	}

	burnTicks := nbtData.GetShort("BurnTicks")

	f.BurnTime = burnTime
	f.CookTime = cookTime
	f.MaxTime = maxTime
	f.BurnTicks = burnTicks
	nbtData.Set(nbt.NewShortTag("BurnTime", f.BurnTime))
	nbtData.Set(nbt.NewShortTag("CookTime", f.CookTime))
	nbtData.Set(nbt.NewShortTag("MaxTime", f.MaxTime))
	nbtData.Set(nbt.NewShortTag("BurnTicks", f.BurnTicks))

	InitSpawnableBase(&f.SpawnableBase, TypeFurnace, chunk, nbtData)
	InitContainerBase(&f.ContainerBase, FurnaceSize)
	f.ContainerBase.LoadItemsFromNBT(nbtData)
	f.NameableBase.LoadNameFromNBT(nbtData)

	return f
}
func (f *Furnace) GetSize() int {
	return FurnaceSize
}
func (f *Furnace) GetSmelting() item.Item {
	return f.GetItem(FurnaceSlotInput)
}
func (f *Furnace) GetFuel() item.Item {
	return f.GetItem(FurnaceSlotFuel)
}
func (f *Furnace) GetResult() item.Item {
	return f.GetItem(FurnaceSlotOutput)
}
func (f *Furnace) SetSmelting(it item.Item) {
	f.SetItem(FurnaceSlotInput, it)
}
func (f *Furnace) SetFuel(it item.Item) {
	f.SetItem(FurnaceSlotFuel, it)
}
func (f *Furnace) SetResult(it item.Item) {
	f.SetItem(FurnaceSlotOutput, it)
}
func (f *Furnace) GetName() string {
	if f.HasCustomName() {
		return f.GetCustomName()
	}
	return "Furnace"
}
func (f *Furnace) IsBurning() bool {
	return f.BurnTime > 0
}
func (f *Furnace) SetNeedUpdate(need bool) {
	f.needUpdate = need
}
func (f *Furnace) OnUpdate() bool {
	if f.IsClosed() {
		return false
	}

	if f.BurnTime > 0 {
		f.BurnTime--
		if f.MaxTime > 0 {
			f.BurnTicks = int16(math.Ceil(float64(f.BurnTime) / float64(f.MaxTime) * 200))
		}
	}
	f.NBT.Set(nbt.NewShortTag("BurnTime", f.BurnTime))
	f.NBT.Set(nbt.NewShortTag("CookTime", f.CookTime))
	f.NBT.Set(nbt.NewShortTag("BurnTicks", f.BurnTicks))

	return true
}
func (f *Furnace) StartBurning(fuelTime int16) {
	f.MaxTime = fuelTime
	f.BurnTime = fuelTime
	f.BurnTicks = 0

	f.NBT.Set(nbt.NewShortTag("MaxTime", f.MaxTime))
	f.NBT.Set(nbt.NewShortTag("BurnTime", f.BurnTime))
	f.NBT.Set(nbt.NewShortTag("BurnTicks", f.BurnTicks))
}
func (f *Furnace) IncrementCookTime() bool {
	f.CookTime++
	if f.CookTime >= 200 {
		f.CookTime -= 200
		f.NBT.Set(nbt.NewShortTag("CookTime", f.CookTime))
		return true
	}
	f.NBT.Set(nbt.NewShortTag("CookTime", f.CookTime))
	return false
}
func (f *Furnace) ResetCookTime() {
	f.CookTime = 0
	f.NBT.Set(nbt.NewShortTag("CookTime", f.CookTime))
}
func (f *Furnace) ResetAll() {
	f.BurnTime = 0
	f.CookTime = 0
	f.BurnTicks = 0
	f.NBT.Set(nbt.NewShortTag("BurnTime", int16(0)))
	f.NBT.Set(nbt.NewShortTag("CookTime", int16(0)))
	f.NBT.Set(nbt.NewShortTag("BurnTicks", int16(0)))
}
func (f *Furnace) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeFurnace))
	compound.Set(nbt.NewIntTag("x", f.X))
	compound.Set(nbt.NewIntTag("y", f.Y))
	compound.Set(nbt.NewIntTag("z", f.Z))
	compound.Set(nbt.NewShortTag("BurnTime", f.BurnTime))
	compound.Set(nbt.NewShortTag("CookTime", f.CookTime))

	if f.HasCustomName() {
		compound.Set(nbt.NewStringTag("CustomName", f.GetCustomName()))
	}

	return compound
}
func (f *Furnace) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	if nbtData.GetString("id") != TypeFurnace {
		return false
	}
	return true
}
func (f *Furnace) SpawnTo(sender PacketSender) bool {
	return SpawnTo(f, sender)
}
func (f *Furnace) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(f, broadcaster)
}
func (f *Furnace) SaveNBT() {
	f.BaseTile.SaveNBT()
	f.ContainerBase.SaveItemsToNBT(f.NBT)
	f.NameableBase.SaveNameToNBT(f.NBT)
	f.NBT.Set(nbt.NewShortTag("BurnTime", f.BurnTime))
	f.NBT.Set(nbt.NewShortTag("CookTime", f.CookTime))
	f.NBT.Set(nbt.NewShortTag("MaxTime", f.MaxTime))
	f.NBT.Set(nbt.NewShortTag("BurnTicks", f.BurnTicks))
}
func (f *Furnace) Close() {
	if f.IsClosed() {
		return
	}
	f.ClearAll()
	f.BaseTile.Close()
}

func init() {
	RegisterTile(TypeFurnace, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewFurnace(chunk, nbtData)
	})
}
