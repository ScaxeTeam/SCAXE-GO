package protocol

const (
	MapBitflagTextureUpdate    = 0x02
	MapBitflagDecorationUpdate = 0x04
)

type MapColor struct {
	R, G, B, A byte
}

type ClientboundMapItemDataPacket struct {
	BasePacket
	MapID   int64
	Type    int32
	Scale   byte
	Width   int32
	Height  int32
	XOffset int32
	YOffset int32
	Colors  [][]MapColor
}

func NewClientboundMapItemDataPacket() *ClientboundMapItemDataPacket {
	return &ClientboundMapItemDataPacket{
		BasePacket: BasePacket{PacketID: IDClientboundMapItemData},
		Width:      128,
		Height:     128,
	}
}

func (p *ClientboundMapItemDataPacket) Name() string {
	return "ClientboundMapItemDataPacket"
}

func (p *ClientboundMapItemDataPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())

	stream.WriteLong(p.MapID)

	typeFlags := int32(0)
	if len(p.Colors) > 0 {
		typeFlags |= MapBitflagTextureUpdate
	}
	stream.WriteInt(typeFlags)

	if (typeFlags & MapBitflagTextureUpdate) != 0 {
		stream.WriteByte(p.Scale)
		stream.WriteInt(p.Width)
		stream.WriteInt(p.Height)
		stream.WriteInt(p.XOffset)
		stream.WriteInt(p.YOffset)

		for y := int32(0); y < p.Height && int(y) < len(p.Colors); y++ {
			row := p.Colors[y]
			for x := int32(0); x < p.Width && int(x) < len(row); x++ {
				color := row[x]
				stream.WriteByte(color.R)
				stream.WriteByte(color.G)
				stream.WriteByte(color.B)
				stream.WriteByte(color.A)
			}
		}
	}

	return nil
}

func (p *ClientboundMapItemDataPacket) Decode(stream *BinaryStream) error {

	return nil
}

func init() {
	RegisterPacket(IDClientboundMapItemData, func() DataPacket { return NewClientboundMapItemDataPacket() })
}
