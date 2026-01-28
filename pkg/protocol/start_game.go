package protocol

type StartGamePacket struct {
	BasePacket
	Seed      int32
	Dimension byte
	Generator int32
	Gamemode  int32
	EntityID  int64
	RuntimeID int64
	SpawnX    int32
	SpawnY    int32
	SpawnZ    int32
	X         float32
	Y         float32
	Z         float32
	LevelID string
}

func NewStartGamePacket() *StartGamePacket {
	return &StartGamePacket{
		BasePacket: BasePacket{PacketID: IDStartGame},
		Dimension:  0,
		Generator:  1,
		Gamemode:   0,
	}
}

func (p *StartGamePacket) Name() string {
	return "StartGamePacket"
}

func (p *StartGamePacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())

	stream.WriteInt(p.Seed)
	stream.WriteByte(p.Dimension)
	stream.WriteInt(p.Generator)
	stream.WriteInt(p.Gamemode)
	stream.WriteLong(p.EntityID)
	stream.WriteUnsignedVarLong(uint64(p.RuntimeID))
	stream.WriteInt(p.SpawnX)
	stream.WriteInt(p.SpawnY)
	stream.WriteInt(p.SpawnZ)
	stream.WriteFloat(p.X)
	stream.WriteFloat(p.Y)
	stream.WriteFloat(p.Z)

	stream.WriteByte(1)
	stream.WriteByte(1)
	stream.WriteByte(0)

	stream.WriteString16(p.LevelID)

	return nil
}

func (p *StartGamePacket) Decode(stream *BinaryStream) error {

	return nil
}

func init() {
	RegisterPacket(IDStartGame, func() DataPacket { return NewStartGamePacket() })
}
