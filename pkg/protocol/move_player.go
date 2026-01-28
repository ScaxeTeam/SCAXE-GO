package protocol

type MovePlayerPacket struct {
	BasePacket
	EntityID int64
	X        float32
	Y        float32
	Z        float32
	Yaw      float32
	BodyYaw  float32
	Pitch    float32
	Mode     byte
	OnGround bool
}

const (
	MovePlayerModeNormal   byte = 0
	MovePlayerModeReset    byte = 1
	MovePlayerModeRotation byte = 2
)

func NewMovePlayerPacket() *MovePlayerPacket {
	return &MovePlayerPacket{
		BasePacket: BasePacket{PacketID: IDMovePlayer},
		Mode:       MovePlayerModeNormal,
	}
}

func (p *MovePlayerPacket) ID() byte {
	return IDMovePlayer
}

func (p *MovePlayerPacket) Name() string {
	return "MovePlayerPacket"
}

func (p *MovePlayerPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteEntityID(p.EntityID)
	stream.WriteFloat(p.X)
	stream.WriteFloat(p.Y)
	stream.WriteFloat(p.Z)
	stream.WriteFloat(p.Yaw)
	stream.WriteFloat(p.BodyYaw)
	stream.WriteFloat(p.Pitch)
	stream.WriteByte(p.Mode)
	stream.WriteBool(p.OnGround)
	return nil
}

func (p *MovePlayerPacket) Decode(stream *BinaryStream) error {
	var err error
	p.EntityID, err = stream.ReadEntityID()
	if err != nil {
		return err
	}
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
	p.Yaw, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.BodyYaw, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.Pitch, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.Mode, err = stream.ReadByte()
	if err != nil {
		return err
	}
	p.OnGround, err = stream.ReadBool()
	return err
}

func init() {
	RegisterPacket(IDMovePlayer, func() DataPacket {
		return NewMovePlayerPacket()
	})
}
