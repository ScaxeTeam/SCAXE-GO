package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/item"
)

type CraftingEventPacket struct {
	BasePacket
	WindowID byte
	Type     int32
	UUID1    int64
	UUID2    int64
	Input    []item.Item
	Output   []item.Item
}

func NewCraftingEventPacket() *CraftingEventPacket {
	return &CraftingEventPacket{
		BasePacket: BasePacket{PacketID: IDCraftingEvent},
		Input:      make([]item.Item, 0),
		Output:     make([]item.Item, 0),
	}
}

func (p *CraftingEventPacket) Name() string {
	return "CraftingEventPacket"
}

func (p *CraftingEventPacket) Encode(stream *BinaryStream) error {

	return nil
}

func (p *CraftingEventPacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())
	var err error

	p.WindowID, err = stream.ReadByte()
	if err != nil {
		return err
	}

	p.Type, err = stream.ReadInt()
	if err != nil {
		return err
	}

	p.UUID1, err = stream.ReadLong()
	if err != nil {
		return err
	}
	p.UUID2, err = stream.ReadLong()
	if err != nil {
		return err
	}

	inputCount, err := stream.ReadInt()
	if err != nil {
		return err
	}
	p.Input = make([]item.Item, inputCount)
	for i := int32(0); i < inputCount; i++ {
		p.Input[i], err = readSlot(stream)
		if err != nil {
			return err
		}
	}

	outputCount, err := stream.ReadInt()
	if err != nil {
		return err
	}
	p.Output = make([]item.Item, outputCount)
	for i := int32(0); i < outputCount; i++ {
		p.Output[i], err = readSlot(stream)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	RegisterPacket(IDCraftingEvent, func() DataPacket { return NewCraftingEventPacket() })
}
