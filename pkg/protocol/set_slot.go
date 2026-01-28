package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/item"
)

type ContainerSetSlotPacket struct {
	BasePacket
	WindowID   byte
	Slot       uint16
	HotbarSlot uint16
	Item       item.Item
}

func NewContainerSetSlotPacket(windowID byte, slot uint16, it item.Item) *ContainerSetSlotPacket {
	return &ContainerSetSlotPacket{
		BasePacket: BasePacket{PacketID: IDContainerSetSlot},
		WindowID:   windowID,
		Slot:       slot,
		HotbarSlot: 0,
		Item:       it,
	}
}

func (p *ContainerSetSlotPacket) Name() string {
	return "ContainerSetSlotPacket"
}

func (p *ContainerSetSlotPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteByte(p.WindowID)
	stream.WriteShort(int16(p.Slot))

	if p.ID() == IDContainerSetSlot {

		stream.WriteShort(int16(p.HotbarSlot))
	}

	stream.WriteSlot(p.Item)
	return nil
}

func (p *ContainerSetSlotPacket) Decode(stream *BinaryStream) error {
	return nil
}
