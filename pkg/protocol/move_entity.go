package protocol

type MoveEntityPacket struct {
	BasePacket
	EntityID int64
	X        float32
	Y        float32
	Z        float32
	Pitch    float32
	HeadYaw  float32
	Yaw      float32
}

func NewMoveEntityPacket() *MoveEntityPacket {
	return &MoveEntityPacket{
		BasePacket: BasePacket{PacketID: IDMoveEntity},
	}
}

func (p *MoveEntityPacket) ID() byte {
	return IDMoveEntity
}

func (p *MoveEntityPacket) Name() string {
	return "MoveEntityPacket"
}

func (p *MoveEntityPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteEntityID(p.EntityID)
	stream.WriteFloat(p.X)
	stream.WriteFloat(p.Y)
	stream.WriteFloat(p.Z)

	stream.WriteByte(byte(p.Pitch / 360.0 * 256.0))
	stream.WriteByte(byte(p.HeadYaw / 360.0 * 256.0))
	stream.WriteByte(byte(p.Yaw / 360.0 * 256.0))

	return nil
}

func (p *MoveEntityPacket) Decode(stream *BinaryStream) error {
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

	pitch, err := stream.ReadByte()
	if err != nil {
		return err
	}
	p.Pitch = float32(pitch) / 256.0 * 360.0

	headYaw, err := stream.ReadByte()
	if err != nil {
		return err
	}
	p.HeadYaw = float32(headYaw) / 256.0 * 360.0

	yaw, err := stream.ReadByte()
	if err != nil {
		return err
	}
	p.Yaw = float32(yaw) / 256.0 * 360.0

	return nil
}

func init() {
	RegisterPacket(IDMoveEntity, func() DataPacket {
		return NewMoveEntityPacket()
	})
}
