package protocol

type EntityEventPacket struct {
	BasePacket
	EntityID int64
	Event    byte
	Data     int32
}

func NewEntityEventPacket() *EntityEventPacket {
	return &EntityEventPacket{
		BasePacket: BasePacket{PacketID: IDEntityEvent},
	}
}

func (p *EntityEventPacket) ID() byte {
	return IDEntityEvent
}

func (p *EntityEventPacket) Name() string {
	return "EntityEventPacket"
}

func (p *EntityEventPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteEntityID(p.EntityID)
	stream.WriteByte(p.Event)

	return nil
}

func (p *EntityEventPacket) Decode(stream *BinaryStream) error {
	var err error
	p.EntityID, err = stream.ReadEntityID()
	if err != nil {
		return err
	}
	p.Event, err = stream.ReadByte()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	RegisterPacket(IDEntityEvent, func() DataPacket {
		return NewEntityEventPacket()
	})
}
