package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/item"
)

type ItemFrameDropItemPacket struct {
	BasePacket
	X        int32
	Y        int32
	Z        int32
	DropItem item.Item
}

func NewItemFrameDropItemPacket() *ItemFrameDropItemPacket {
	return &ItemFrameDropItemPacket{
		BasePacket: BasePacket{PacketID: IDItemFrameDropItem},
	}
}

func (p *ItemFrameDropItemPacket) Name() string {
	return "ItemFrameDropItemPacket"
}

func (p *ItemFrameDropItemPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())

	stream.WriteInt(p.Z)
	stream.WriteInt(p.Y)
	stream.WriteInt(p.X)
	stream.WriteSlot(p.DropItem)
	return nil
}

func (p *ItemFrameDropItemPacket) Decode(stream *BinaryStream) error {
	var err error

	p.Z, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.Y, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.X, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.DropItem, err = stream.ReadSlot()
	return err
}

func init() {
	RegisterPacket(IDItemFrameDropItem, func() DataPacket { return NewItemFrameDropItemPacket() })
}
