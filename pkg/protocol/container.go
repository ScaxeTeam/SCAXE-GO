package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/logger"
)

type ContainerSetContentPacket struct {
	BasePacket
	WindowID    byte
	Items       []item.Item
	HotbarTypes []int32
}

func NewContainerSetContentPacket(windowID byte, items []item.Item) *ContainerSetContentPacket {
	return &ContainerSetContentPacket{
		BasePacket:  BasePacket{PacketID: IDContainerSetContent},
		WindowID:    windowID,
		Items:       items,
		HotbarTypes: nil,
	}
}

func (p *ContainerSetContentPacket) Name() string {
	return "ContainerSetContentPacket"
}

func (p *ContainerSetContentPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteByte(p.WindowID)

	stream.WriteShort(int16(len(p.Items)))

	for _, it := range p.Items {
		stream.WriteSlot(it)
	}

	if p.WindowID == 0 && len(p.HotbarTypes) > 0 {
		stream.WriteShort(int16(len(p.HotbarTypes)))
		for _, slot := range p.HotbarTypes {
			stream.WriteInt(slot)
		}
	} else {

		stream.WriteShort(0)
	}

	logger.Debug("ContainerSetContent",
		"windowID", p.WindowID,
		"itemCount", len(p.Items),
		"hotbarCount", len(p.HotbarTypes),
		"totalBytes", stream.Len())

	return nil
}

func (p *ContainerSetContentPacket) Decode(stream *BinaryStream) error {
	return nil
}
