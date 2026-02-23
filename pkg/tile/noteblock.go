package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Noteblock struct {
	BaseTile
}

func NewNoteblock(chunk *world.Chunk, nbtData *nbt.CompoundTag) *Noteblock {
	n := &Noteblock{}

	if nbtData.Get("note") == nil {
		nbtData.Set(nbt.NewByteTag("note", 0))
	}
	if nbtData.Get("powered") == nil {
		nbtData.Set(nbt.NewByteTag("powered", 0))
	}

	InitBaseTile(&n.BaseTile, TypeNoteblock, chunk, nbtData)
	return n
}

func (n *Noteblock) GetName() string {
	return "Noteblock"
}

func (n *Noteblock) GetPitch() int {
	return int(n.NBT.GetByte("note"))
}

func (n *Noteblock) SetPitch(pitch int) {
	if pitch < 0 || pitch > 24 {
		return
	}
	n.NBT.Set(nbt.NewByteTag("note", int8(pitch)))
}

func (n *Noteblock) IsPowered() bool {
	return n.NBT.GetByte("powered") > 0
}

func (n *Noteblock) SetPowered(powered bool) {
	val := int8(0)
	if powered {
		val = 1
	}
	n.NBT.Set(nbt.NewByteTag("powered", val))
}

func init() {
	RegisterTile(TypeNoteblock, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewNoteblock(chunk, nbtData)
	})
}
