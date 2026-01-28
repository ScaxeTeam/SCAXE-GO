package protocol

type SetSpawnPositionPacket struct {
	BasePacket
	X int32
	Y int32
	Z int32
}

func NewSetSpawnPositionPacket() *SetSpawnPositionPacket {
	return &SetSpawnPositionPacket{
		BasePacket: BasePacket{PacketID: IDSetSpawnPosition},
	}
}

func (p *SetSpawnPositionPacket) Name() string {
	return "SetSpawnPositionPacket"
}

func (p *SetSpawnPositionPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteInt(p.X)
	stream.WriteInt(p.Y)
	stream.WriteInt(p.Z)
	return nil
}

func (p *SetSpawnPositionPacket) Decode(stream *BinaryStream) error {
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
	return err
}

func init() {
	RegisterPacket(IDSetSpawnPosition, func() DataPacket { return NewSetSpawnPositionPacket() })
}
