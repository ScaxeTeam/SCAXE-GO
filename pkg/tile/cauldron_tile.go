package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Cauldron struct {
	SpawnableBase
}

func NewCauldron(chunk *world.Chunk, nbtData *nbt.CompoundTag) *Cauldron {
	c := &Cauldron{}

	if nbtData.Get("PotionId") == nil {
		nbtData.Set(nbt.NewShortTag("PotionId", -1))
	}
	if nbtData.Get("SplashPotion") == nil {
		nbtData.Set(nbt.NewByteTag("SplashPotion", 0))
	}
	if nbtData.Get("Items") == nil {
		nbtData.Set(nbt.NewListTag("Items", nbt.TagCompound))
	}

	InitSpawnableBase(&c.SpawnableBase, TypeCauldron, chunk, nbtData)
	return c
}

func (c *Cauldron) GetName() string {
	return "Cauldron"
}

func (c *Cauldron) GetPotionId() int16 {
	return c.NBT.GetShort("PotionId")
}

func (c *Cauldron) SetPotionId(potionId int16) {
	c.NBT.Set(nbt.NewShortTag("PotionId", potionId))
}

func (c *Cauldron) HasPotion() bool {
	return c.GetPotionId() != -1
}

func (c *Cauldron) GetSplashPotion() bool {
	return c.NBT.GetByte("SplashPotion") != 0
}

func (c *Cauldron) SetSplashPotion(splash bool) {
	val := int8(0)
	if splash {
		val = 1
	}
	c.NBT.Set(nbt.NewByteTag("SplashPotion", val))
}

func (c *Cauldron) IsCustomColor() bool {
	return c.NBT.Get("CustomColor") != nil
}

func (c *Cauldron) GetCustomColor() (r, g, b int) {
	if !c.IsCustomColor() {
		return 0, 0, 0
	}
	color := c.NBT.GetInt("CustomColor")
	r = int((color >> 16) & 0xFF)
	g = int((color >> 8) & 0xFF)
	b = int(color & 0xFF)
	return
}

func (c *Cauldron) SetCustomColor(r, g, b int) {
	color := int32((r&0xFF)<<16 | (g&0xFF)<<8 | (b & 0xFF))
	c.NBT.Set(nbt.NewIntTag("CustomColor", color))
}

func (c *Cauldron) ClearCustomColor() {
	c.NBT.Remove("CustomColor")
}

func (c *Cauldron) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeCauldron))
	compound.Set(nbt.NewIntTag("x", c.X))
	compound.Set(nbt.NewIntTag("y", c.Y))
	compound.Set(nbt.NewIntTag("z", c.Z))
	compound.Set(nbt.NewShortTag("PotionId", c.GetPotionId()))

	splashVal := int8(0)
	if c.GetSplashPotion() {
		splashVal = 1
	}
	compound.Set(nbt.NewByteTag("SplashPotion", splashVal))

	if items := c.NBT.Get("Items"); items != nil {
		compound.Set(items.Clone())
	} else {
		compound.Set(nbt.NewListTag("Items", nbt.TagCompound))
	}

	if c.GetPotionId() == -1 && c.IsCustomColor() {
		compound.Set(c.NBT.Get("CustomColor").Clone())
	}

	return compound
}

func (c *Cauldron) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return false
}

func (c *Cauldron) SpawnTo(sender PacketSender) bool {
	return SpawnTo(c, sender)
}

func (c *Cauldron) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(c, broadcaster)
}

func init() {
	RegisterTile(TypeCauldron, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewCauldron(chunk, nbtData)
	})
}
