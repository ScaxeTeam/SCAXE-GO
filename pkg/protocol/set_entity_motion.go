package protocol

type SetEntityMotionPacket struct {
	BasePacket
	EntityID int64
	SpeedX   float32
	SpeedY   float32
	SpeedZ   float32
}

func NewSetEntityMotionPacket() *SetEntityMotionPacket {
	return &SetEntityMotionPacket{
		BasePacket: BasePacket{PacketID: IDSetEntityMotion},
	}
}

func (p *SetEntityMotionPacket) ID() byte {
	return IDSetEntityMotion
}

func (p *SetEntityMotionPacket) Name() string {
	return "SetEntityMotionPacket"
}

func (p *SetEntityMotionPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteEntityID(p.EntityID)
	stream.WriteFloat(p.SpeedX)
	stream.WriteFloat(p.SpeedY)
	stream.WriteFloat(p.SpeedZ)
	return nil
}

func (p *SetEntityMotionPacket) Decode(stream *BinaryStream) error {
	var err error
	p.EntityID, err = stream.ReadEntityID()
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
	return nil
}

func init() {
	RegisterPacket(IDSetEntityMotion, func() DataPacket {
		return NewSetEntityMotionPacket()
	})
}
