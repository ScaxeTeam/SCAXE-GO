package protocol

type MapInfoRequestPacket struct {
	BasePacket
	MapID int64
}

func NewMapInfoRequestPacket() *MapInfoRequestPacket {
	return &MapInfoRequestPacket{
		BasePacket: BasePacket{PacketID: IDMapInfoRequest},
	}
}

func (p *MapInfoRequestPacket) Name() string {
	return "MapInfoRequestPacket"
}

func (p *MapInfoRequestPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteLong(p.MapID)
	return nil
}

func (p *MapInfoRequestPacket) Decode(stream *BinaryStream) error {
	var err error
	p.MapID, err = stream.ReadLong()
	return err
}

func init() {
	RegisterPacket(IDMapInfoRequest, func() DataPacket { return NewMapInfoRequestPacket() })
}
