package protocol

type BlockEventPacket struct {
	BasePacket
	X     int32
	Y     int32
	Z     int32
	Case1 int32
	Case2 int32
}

func NewBlockEventPacket() *BlockEventPacket {
	return &BlockEventPacket{
		BasePacket: BasePacket{PacketID: IDBlockEvent},
	}
}

func (p *BlockEventPacket) Name() string {
	return "BlockEventPacket"
}

func (p *BlockEventPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteInt(p.X)
	stream.WriteInt(p.Y)
	stream.WriteInt(p.Z)
	stream.WriteInt(p.Case1)
	stream.WriteInt(p.Case2)
	return nil
}

func (p *BlockEventPacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())
	var err error
	p.X, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.Y, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.Z, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.Case1, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.Case2, err = stream.ReadInt()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RegisterPacket(IDBlockEvent, func() DataPacket { return NewBlockEventPacket() })
}
