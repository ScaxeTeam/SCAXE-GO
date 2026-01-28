package protocol

type TakeItemEntityPacket struct {
	BasePacket
	Target   int64
	EntityID int64
}

func NewTakeItemEntityPacket() *TakeItemEntityPacket {
	return &TakeItemEntityPacket{
		BasePacket: BasePacket{PacketID: IDTakeItemEntity},
	}
}

func (p *TakeItemEntityPacket) Name() string {
	return "TakeItemEntityPacket"
}

func (p *TakeItemEntityPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())

	stream.WriteLong(p.EntityID)
	stream.WriteLong(p.Target)
	return nil
}

func (p *TakeItemEntityPacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())
	var err error
	p.EntityID, err = stream.ReadLong()
	if err != nil {
		return err
	}
	p.Target, err = stream.ReadLong()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RegisterPacket(IDTakeItemEntity, func() DataPacket { return NewTakeItemEntityPacket() })
}
