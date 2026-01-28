package protocol

import "github.com/scaxe/scaxe-go/pkg/item"

type AddItemEntityPacket struct {
	BasePacket
	EntityID int64
	Item     item.Item
	X        float32
	Y        float32
	Z        float32
	SpeedX   float32
	SpeedY   float32
	SpeedZ   float32
	Metadata []byte

}

func NewAddItemEntityPacket() *AddItemEntityPacket {
	return &AddItemEntityPacket{
		BasePacket: BasePacket{PacketID: IDAddItemEntity},
	}
}

func (p *AddItemEntityPacket) ID() byte {
	return IDAddItemEntity
}

func (p *AddItemEntityPacket) Name() string {
	return "AddItemEntityPacket"
}

func (p *AddItemEntityPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteEntityID(p.EntityID)

	stream.WriteSlot(p.Item)
	stream.WriteFloat(p.X)
	stream.WriteFloat(p.Y)
	stream.WriteFloat(p.Z)
	stream.WriteFloat(p.SpeedX)
	stream.WriteFloat(p.SpeedY)
	stream.WriteFloat(p.SpeedZ)

	return nil
}

func (p *AddItemEntityPacket) Decode(stream *BinaryStream) error {
	var err error
	p.EntityID, err = stream.ReadEntityID()
	if err != nil {
		return err
	}
	p.Item, err = stream.ReadSlot()
	if err != nil {
		return err
	}
	p.X, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.Y, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.Z, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.SpeedX, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.SpeedY, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.SpeedZ, err = stream.ReadFloat()
	if err != nil {
		return err
	}

	return err
}

func init() {
	RegisterPacket(IDAddItemEntity, func() DataPacket {
		return NewAddItemEntityPacket()
	})
}
