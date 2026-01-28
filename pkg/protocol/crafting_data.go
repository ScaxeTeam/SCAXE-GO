package protocol

import (
	"github.com/google/uuid"
	"github.com/scaxe/scaxe-go/pkg/item"
)

const (
	EntryShapeless   int32 = 0
	EntryShaped      int32 = 1
	EntryFurnace     int32 = 2
	EntryFurnaceData int32 = 3
	EntryEnchantList int32 = 4
)

type ShapelessRecipe struct {
	UUID   uuid.UUID
	Inputs []item.Item
	Output item.Item
}

type ShapedRecipe struct {
	UUID   uuid.UUID
	Width  int32
	Height int32
	Inputs []item.Item
	Output item.Item
}

type FurnaceRecipe struct {
	InputID   int32
	InputMeta int32
	Output    item.Item
}

type EnchantmentEntry struct {
	Cost         int32
	Enchantments []EnchantData
	RandomName   string
}

type EnchantData struct {
	ID    int32
	Level int32
}

type EnchantmentList struct {
	Entries []EnchantmentEntry
}

type CraftingEntry struct {
	Type int32
	Data []byte
}

type CraftingDataPacket struct {
	BasePacket
	Entries      []CraftingEntry
	CleanRecipes bool
}

func NewCraftingDataPacket() *CraftingDataPacket {
	return &CraftingDataPacket{
		BasePacket: BasePacket{PacketID: IDCraftingData},
		Entries:    make([]CraftingEntry, 0),
	}
}

func (p *CraftingDataPacket) Name() string {
	return "CraftingDataPacket"
}

func (p *CraftingDataPacket) AddShapelessRecipe(recipe ShapelessRecipe) {
	stream := NewBinaryStream()

	stream.WriteInt(int32(len(recipe.Inputs)))
	for _, inp := range recipe.Inputs {
		stream.WriteSlot(inp)
	}

	stream.WriteInt(1)
	stream.WriteSlot(recipe.Output)

	stream.WriteUUID(recipe.UUID.String())

	p.Entries = append(p.Entries, CraftingEntry{
		Type: EntryShapeless,
		Data: stream.Bytes(),
	})
}

func (p *CraftingDataPacket) AddShapedRecipe(recipe ShapedRecipe) {
	stream := NewBinaryStream()

	stream.WriteInt(recipe.Width)
	stream.WriteInt(recipe.Height)

	for _, inp := range recipe.Inputs {
		stream.WriteSlot(inp)
	}

	stream.WriteInt(1)
	stream.WriteSlot(recipe.Output)

	stream.WriteUUID(recipe.UUID.String())

	p.Entries = append(p.Entries, CraftingEntry{
		Type: EntryShaped,
		Data: stream.Bytes(),
	})
}

func (p *CraftingDataPacket) AddFurnaceRecipe(recipe FurnaceRecipe) {
	stream := NewBinaryStream()

	entryType := EntryFurnace
	if recipe.InputMeta != -1 {

		stream.WriteInt((recipe.InputID << 16) | recipe.InputMeta)
		entryType = EntryFurnaceData
	} else {
		stream.WriteInt(recipe.InputID)
	}

	stream.WriteSlot(recipe.Output)

	p.Entries = append(p.Entries, CraftingEntry{
		Type: entryType,
		Data: stream.Bytes(),
	})
}

func (p *CraftingDataPacket) AddEnchantList(list EnchantmentList) {
	stream := NewBinaryStream()

	stream.WriteByte(byte(len(list.Entries)))
	for _, entry := range list.Entries {
		stream.WriteInt(entry.Cost)
		stream.WriteByte(byte(len(entry.Enchantments)))
		for _, ench := range entry.Enchantments {
			stream.WriteInt(ench.ID)
			stream.WriteInt(ench.Level)
		}
		stream.WriteString(entry.RandomName)
	}

	p.Entries = append(p.Entries, CraftingEntry{
		Type: EntryEnchantList,
		Data: stream.Bytes(),
	})
}

func (p *CraftingDataPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())

	stream.WriteInt(int32(len(p.Entries)))

	for _, entry := range p.Entries {
		stream.WriteInt(entry.Type)
		stream.WriteInt(int32(len(entry.Data)))
		stream.WriteBytes(entry.Data)
	}

	if p.CleanRecipes {
		stream.WriteByte(1)
	} else {
		stream.WriteByte(0)
	}
	return nil
}

func (p *CraftingDataPacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())

	return nil
}

func init() {
	RegisterPacket(IDCraftingData, func() DataPacket { return NewCraftingDataPacket() })
}
