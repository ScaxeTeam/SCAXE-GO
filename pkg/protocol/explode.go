package protocol

type ExplodePacket struct {
	BasePacket
	X       float32
	Y       float32
	Z       float32
	Radius  float32
	Records []ExplodeRecord
}

type ExplodeRecord struct {
	X int8
	Y int8
	Z int8
}

func NewExplodePacket() *ExplodePacket {
	return &ExplodePacket{
		BasePacket: BasePacket{PacketID: IDExplode},
		Records:    make([]ExplodeRecord, 0),
	}
}

func (p *ExplodePacket) Name() string {
	return "ExplodePacket"
}

func (p *ExplodePacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteFloat(p.X)
	stream.WriteFloat(p.Y)
	stream.WriteFloat(p.Z)
	stream.WriteFloat(p.Radius)
	stream.WriteInt(int32(len(p.Records)))
	for _, record := range p.Records {
		stream.WriteByte(byte(record.X))
		stream.WriteByte(byte(record.Y))
		stream.WriteByte(byte(record.Z))
	}
	return nil
}

func (p *ExplodePacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())
	var err error
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
	p.Radius, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	count, err := stream.ReadInt()
	if err != nil {
		return err
	}
	p.Records = make([]ExplodeRecord, count)
	for i := int32(0); i < count; i++ {
		x, err := stream.ReadByte()
		if err != nil {
			return err
		}
		y, err := stream.ReadByte()
		if err != nil {
			return err
		}
		z, err := stream.ReadByte()
		if err != nil {
			return err
		}
		p.Records[i] = ExplodeRecord{X: int8(x), Y: int8(y), Z: int8(z)}
	}
	return nil
}

func init() {
	RegisterPacket(IDExplode, func() DataPacket { return NewExplodePacket() })
}
