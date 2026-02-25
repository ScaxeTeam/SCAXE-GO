package tile

// flower_pot.go — 花盆 TileEntity
// 对应 PHP FlowerPot.php
// Spawnable: 存储种植的植物 ID 和 meta

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

// 花盆允许放置的植物方块 ID
var flowerPotAllowed = map[int16]bool{
	31: true, // TALL_GRASS
	6:  true, // SAPLING
	32: true, // DEAD_BUSH
	37: true, // DANDELION
	38: true, // RED_FLOWER
	39: true, // BROWN_MUSHROOM
	40: true, // RED_MUSHROOM
	81: true, // CACTUS
}

// FlowerPot 花盆 TileEntity
type FlowerPot struct {
	SpawnableBase
	PlantID   int16
	PlantData int32
}

// NewFlowerPot 创建花盆实例
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

// GetSpawnCompound 返回客户端渲染用的 NBT
func (fp *FlowerPot) GetSpawnCompound() *nbt.CompoundTag {
	// 校验植物 ID 是否合法
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

// SaveNBT 保存花盆数据
func (fp *FlowerPot) SaveNBT() {
	fp.SpawnableBase.SaveNBT()
	fp.NBT.Set(nbt.NewShortTag("item", fp.PlantID))
	fp.NBT.Set(nbt.NewIntTag("mData", fp.PlantData))
}

// SetFlowerPotData 设置花盆中的植物
func (fp *FlowerPot) SetFlowerPotData(itemID int16, data int32) {
	fp.PlantID = itemID
	fp.PlantData = data
}

func init() {
	RegisterTile(TypeFlowerPot, NewFlowerPot)
}
