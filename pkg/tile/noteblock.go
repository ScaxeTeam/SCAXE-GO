package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

// Noteblock 音符盒 TileEntity
// 对应 PHP class Noteblock extends Tile（注意：不是 Spawnable）
type Noteblock struct {
	BaseTile
}

// NewNoteblock 创建 Noteblock 实例
// 对应 PHP Noteblock::__construct(FullChunk $chunk, CompoundTag $nbt)
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

// GetName 返回名称
func (n *Noteblock) GetName() string {
	return "Noteblock"
}

// GetPitch 获取音高 (0-24)
// 对应 PHP Noteblock::getPitch()
func (n *Noteblock) GetPitch() int {
	return int(n.NBT.GetByte("note"))
}

// SetPitch 设置音高 (0-24)
// 对应 PHP Noteblock::setPitch(int $pitch)
func (n *Noteblock) SetPitch(pitch int) {
	if pitch < 0 || pitch > 24 {
		return
	}
	n.NBT.Set(nbt.NewByteTag("note", int8(pitch)))
}

// IsPowered 是否被红石激活
// 对应 PHP Noteblock::isPowered()
func (n *Noteblock) IsPowered() bool {
	return n.NBT.GetByte("powered") > 0
}

// SetPowered 设置红石激活状态
// 对应 PHP Noteblock::setPowered(bool $powered)
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
