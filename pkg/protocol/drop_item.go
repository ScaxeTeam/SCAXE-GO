package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/item"
)

type DropItemPacket struct {
	BasePacket
	Type byte
	Item item.Item
}

func NewDropItemPacket() *DropItemPacket {
	return &DropItemPacket{
		BasePacket: BasePacket{PacketID: IDDropItem},
	}
}

func (p *DropItemPacket) Name() string {
	return "DropItemPacket"
}

func (p *DropItemPacket) Encode(stream *BinaryStream) error {

	return nil
}

func (p *DropItemPacket) Decode(stream *BinaryStream) error {
	var err error
	p.Type, err = stream.ReadByte()
	if err != nil {
		return err
	}
	p.Item, err = stream.ReadSlot()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RegisterPacket(IDDropItem, func() DataPacket { return NewDropItemPacket() })
}
