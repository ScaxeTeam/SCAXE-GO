package protocol

type BlockRecord struct {
	X         int32
	Z         int32
	Y         byte
	BlockID   byte
	BlockMeta byte
	Flags     byte
}

type UpdateBlockPacket struct {
	BasePacket
	Records []BlockRecord
}

const (
	UpdateBlockFlagNone         = 0x00
	UpdateBlockFlagNeighborhood = 0x01
	UpdateBlockFlagNetwork      = 0x02
	UpdateBlockFlagNoGraphic    = 0x04
	UpdateBlockFlagPriority     = 0x08
	UpdateBlockPacketFlagAll    = UpdateBlockFlagNetwork | UpdateBlockFlagNeighborhood
)

func NewUpdateBlockPacket(x, y, z int32, id, meta uint8) *UpdateBlockPacket {

	return &UpdateBlockPacket{
		BasePacket: BasePacket{PacketID: IDUpdateBlock},
		Records: []BlockRecord{
			{
				X:         x,
				Z:         z,
				Y:         byte(y),
				BlockID:   id,
				BlockMeta: meta,
				Flags:     UpdateBlockFlagNetwork | UpdateBlockFlagNeighborhood,
			},
		},
	}
}

func (p *UpdateBlockPacket) Name() string {
	return "UpdateBlockPacket"
}

func (p *UpdateBlockPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())

	stream.WriteInt(int32(len(p.Records)))

	for _, r := range p.Records {
		stream.WriteInt(r.X)
		stream.WriteInt(r.Z)
		stream.WriteByte(r.Y)
		stream.WriteByte(r.BlockID)

		packed := (r.Flags << 4) | (r.BlockMeta & 0x0F)
		stream.WriteByte(packed)
	}

	return nil
}

func (p *UpdateBlockPacket) Decode(stream *BinaryStream) error {
	return nil
}
