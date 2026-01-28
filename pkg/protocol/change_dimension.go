package protocol

const (
	DimensionOverworld byte = 0
	DimensionNether    byte = 1
	DimensionEnd       byte = 2
)

type ChangeDimensionPacket struct {
	BasePacket
	Dimension byte
	X         float32
	Y         float32
	Z         float32
	Unknown   bool
}

func NewChangeDimensionPacket() *ChangeDimensionPacket {
	return &ChangeDimensionPacket{
		BasePacket: BasePacket{PacketID: IDChangeDimension},
	}
}

func (p *ChangeDimensionPacket) Name() string {
	return "ChangeDimensionPacket"
}

func (p *ChangeDimensionPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteInt(int32(p.Dimension))
	stream.WriteFloat(p.X)
	stream.WriteFloat(p.Y)
	stream.WriteFloat(p.Z)
	if p.Unknown {
		stream.WriteByte(1)
	} else {
		stream.WriteByte(0)
	}
	return nil
}

func (p *ChangeDimensionPacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())
	var err error
	dim, err := stream.ReadInt()
	if err != nil {
		return err
	}
	p.Dimension = byte(dim)
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
	unknown, err := stream.ReadByte()
	if err != nil {
		return err
	}
	p.Unknown = unknown != 0
	return nil
}

func init() {
	RegisterPacket(IDChangeDimension, func() DataPacket { return NewChangeDimensionPacket() })
}
