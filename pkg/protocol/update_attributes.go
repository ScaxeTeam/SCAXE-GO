package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/entity/attribute"
)

type UpdateAttributesPacket struct {
	BasePacket
	EntityID int64
	Entries  []*attribute.Attribute
}

func NewUpdateAttributesPacket() *UpdateAttributesPacket {
	return &UpdateAttributesPacket{
		BasePacket: BasePacket{PacketID: IDUpdateAttributes},
		Entries:    make([]*attribute.Attribute, 0),
	}
}

func (p *UpdateAttributesPacket) ID() byte {
	return IDUpdateAttributes
}

func (p *UpdateAttributesPacket) Name() string {
	return "UpdateAttributesPacket"
}

func (p *UpdateAttributesPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteEntityID(p.EntityID)
	stream.WriteShort(int16(len(p.Entries)))

	for _, entry := range p.Entries {
		stream.WriteFloat(entry.MinValue)
		stream.WriteFloat(entry.MaxValue)
		stream.WriteFloat(entry.CurrentValue)

		stream.WriteString16(entry.Name)
	}
	return nil
}

func (p *UpdateAttributesPacket) Decode(stream *BinaryStream) error {
	var err error
	p.EntityID, err = stream.ReadEntityID()
	if err != nil {
		return err
	}

	count, err := stream.ReadShort()
	if err != nil {
		return err
	}

	p.Entries = make([]*attribute.Attribute, count)
	for i := 0; i < int(count); i++ {
		min, _ := stream.ReadFloat()
		max, _ := stream.ReadFloat()
		val, _ := stream.ReadFloat()
		name, _ := stream.ReadString16()
		p.Entries[i] = attribute.NewAttribute(name, min, max, val)
	}
	return nil
}

func init() {
	RegisterPacket(IDUpdateAttributes, func() DataPacket {
		return NewUpdateAttributesPacket()
	})
}
