package protocol

type AddEntityPacket struct {
	BasePacket
	EntityID int64
	Type     int32
	X        float32
	Y        float32
	Z        float32
	SpeedX   float32
	SpeedY   float32
	SpeedZ   float32
	Yaw      float32
	Pitch    float32
	Metadata []byte

}

func NewAddEntityPacket() *AddEntityPacket {
	return &AddEntityPacket{
		BasePacket: BasePacket{PacketID: IDAddEntity},
		Metadata:   []byte{0x7f},
	}
}

func (p *AddEntityPacket) ID() byte {
	return IDAddEntity
}

func (p *AddEntityPacket) Name() string {
	return "AddEntityPacket"
}

func (p *AddEntityPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteEntityID(p.EntityID)
	stream.WriteInt(p.Type)
	stream.WriteFloat(p.X)
	stream.WriteFloat(p.Y)
	stream.WriteFloat(p.Z)
	stream.WriteFloat(p.SpeedX)
	stream.WriteFloat(p.SpeedY)
	stream.WriteFloat(p.SpeedZ)
	stream.WriteFloat(p.Yaw)
	stream.WriteFloat(p.Pitch)
	stream.WriteBytes(p.Metadata)

	stream.WriteShort(0)

	return nil
}

func (p *AddEntityPacket) Decode(stream *BinaryStream) error {
	var err error
	p.EntityID, err = stream.ReadEntityID()
	if err != nil {
		return err
	}
	p.Type, err = stream.ReadInt()
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
	p.SpeedX, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.SpeedY, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.SpeedZ, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.Yaw, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.Pitch, err = stream.ReadFloat()
	if err != nil {
		return err
	}

	return err
}

func init() {
	RegisterPacket(IDAddEntity, func() DataPacket {
		return NewAddEntityPacket()
	})
}
