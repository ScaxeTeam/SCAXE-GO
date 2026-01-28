package protocol

func init() {
	RegisterPacket(IDRemoveBlock, func() DataPacket { return NewRemoveBlockPacket() })
}

type RemoveBlockPacket struct {
	BasePacket
	EntityID int64
	X        int32
	Z        int32
	Y        byte
}

func NewRemoveBlockPacket() *RemoveBlockPacket {
	return &RemoveBlockPacket{
		BasePacket: BasePacket{PacketID: IDRemoveBlock},
	}
}

func (p *RemoveBlockPacket) ID() byte {
	return IDRemoveBlock
}

func (p *RemoveBlockPacket) Name() string {
	return "RemoveBlockPacket"
}

func (p *RemoveBlockPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteBELong(p.EntityID)
	stream.WriteBEInt(p.X)
	stream.WriteBEInt(p.Z)
	stream.WriteByte(p.Y)
	return nil
}

func (p *RemoveBlockPacket) Decode(stream *BinaryStream) error {
	var err error

	p.EntityID, err = stream.ReadBELong()
	if err != nil {
		return err
	}

	p.X, err = stream.ReadBEInt()
	if err != nil {
		return err
	}

	p.Z, err = stream.ReadBEInt()
	if err != nil {
		return err
	}

	p.Y, err = stream.ReadByte()
	return err
}
