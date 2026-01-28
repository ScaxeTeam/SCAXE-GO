package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/item"
)

type ReplaceSelectedItemPacket struct {
	BasePacket
	Slot item.Item
}

func NewReplaceSelectedItemPacket() *ReplaceSelectedItemPacket {
	return &ReplaceSelectedItemPacket{
		BasePacket: BasePacket{PacketID: IDReplaceSelectedItem},
	}
}

func (p *ReplaceSelectedItemPacket) Name() string {
	return "ReplaceSelectedItemPacket"
}

func (p *ReplaceSelectedItemPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteSlot(p.Slot)
	return nil
}

func (p *ReplaceSelectedItemPacket) Decode(stream *BinaryStream) error {
	var err error
	p.Slot, err = stream.ReadSlot()
	return err
}

func init() {
	RegisterPacket(IDReplaceSelectedItem, func() DataPacket { return NewReplaceSelectedItemPacket() })
}
