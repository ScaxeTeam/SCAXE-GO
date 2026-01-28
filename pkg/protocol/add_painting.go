package protocol

type AddPaintingPacket struct {
	BasePacket
	EntityID  int64
	X         int32
	Y         int32
	Z         int32
	Direction int32
	Title     string
}

func NewAddPaintingPacket() *AddPaintingPacket {
	return &AddPaintingPacket{
		BasePacket: BasePacket{PacketID: IDAddPainting},
	}
}

func (p *AddPaintingPacket) Name() string {
	return "AddPaintingPacket"
}

func (p *AddPaintingPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())

	stream.WriteLong(p.EntityID)

	stream.WriteInt(p.X)
	stream.WriteInt(p.Y)
	stream.WriteInt(p.Z)

	stream.WriteInt(p.Direction)

	stream.WriteString(p.Title)

	return nil
}

func (p *AddPaintingPacket) Decode(stream *BinaryStream) error {

	return nil
}

func init() {
	RegisterPacket(IDAddPainting, func() DataPacket { return NewAddPaintingPacket() })
}
