package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type ItemFrame struct {
	SpawnableBase
}

func NewItemFrame(chunk *world.Chunk, nbtData *nbt.CompoundTag) *ItemFrame {
	f := &ItemFrame{}

	if nbtData.Get("Item") == nil {

		itemTag := nbt.NewCompoundTag("Item")
		itemTag.Set(nbt.NewShortTag("id", 0))
		itemTag.Set(nbt.NewShortTag("Damage", 0))
		itemTag.Set(nbt.NewByteTag("Count", 0))
		nbtData.Set(itemTag)
	}

	if nbtData.Get("ItemRotation") == nil {
		nbtData.Set(nbt.NewByteTag("ItemRotation", 0))
	}

	if nbtData.Get("ItemDropChance") == nil {
		nbtData.Set(nbt.NewFloatTag("ItemDropChance", 1.0))
	}

	InitSpawnableBase(&f.SpawnableBase, TypeItemFrame, chunk, nbtData)
	return f
}

func (f *ItemFrame) GetName() string {
	return "Item Frame"
}

func (f *ItemFrame) GetItemRotation() int8 {
	return f.NBT.GetByte("ItemRotation")
}

func (f *ItemFrame) SetItemRotation(rotation int8) {
	f.NBT.Set(nbt.NewByteTag("ItemRotation", rotation))
}

func (f *ItemFrame) GetItemDropChance() float32 {
	return f.NBT.GetFloat("ItemDropChance")
}

func (f *ItemFrame) SetItemDropChance(chance float32) {
	f.NBT.Set(nbt.NewFloatTag("ItemDropChance", chance))
}

func (f *ItemFrame) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeItemFrame))
	compound.Set(nbt.NewIntTag("x", f.X))
	compound.Set(nbt.NewIntTag("y", f.Y))
	compound.Set(nbt.NewIntTag("z", f.Z))
	compound.Set(nbt.NewByteTag("ItemRotation", f.GetItemRotation()))
	compound.Set(nbt.NewFloatTag("ItemDropChance", f.GetItemDropChance()))

	if itemTag := f.NBT.Get("Item"); itemTag != nil {
		if ct, ok := itemTag.(*nbt.CompoundTag); ok {

			if ct.GetShort("id") != 0 {
				clone := ct.Clone().(*nbt.CompoundTag)
				clone.SetName("Item")
				compound.Set(clone)
			}
		}
	}

	return compound
}

func (f *ItemFrame) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return false
}

func (f *ItemFrame) SpawnTo(sender PacketSender) bool {
	return SpawnTo(f, sender)
}

func (f *ItemFrame) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(f, broadcaster)
}

func init() {
	RegisterTile(TypeItemFrame, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewItemFrame(chunk, nbtData)
	})
}
