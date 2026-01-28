package protocol

type AnimatePacket struct {
	BasePacket
	Action   byte
	EntityID int64
	Float    float32
}

const (
	AnimateActionSwingArm    byte = 1
	AnimateActionStopSleep   byte = 3
	AnimateActionCriticalHit byte = 4
	AnimateActionRowRight    byte = 128
	AnimateActionRowLeft     byte = 129
)

func NewAnimatePacket() *AnimatePacket {
	return &AnimatePacket{
		BasePacket: BasePacket{PacketID: IDAnimate},
	}
}

func (p *AnimatePacket) ID() byte {
	return IDAnimate
}

func (p *AnimatePacket) Name() string {
	return "AnimatePacket"
}

func (p *AnimatePacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteByte(p.Action)
	stream.WriteEntityID(p.EntityID)
	if (p.Action & 0x80) != 0 {
		stream.WriteFloat(p.Float)
	}
	return nil
}

func (p *AnimatePacket) Decode(stream *BinaryStream) error {
	var err error
	p.Action, err = stream.ReadByte()
	if err != nil {
		return err
	}
	p.EntityID, err = stream.ReadEntityID()
	if err != nil {
		return err
	}
	if (p.Action & 0x80) != 0 {
		p.Float, err = stream.ReadFloat()
	}
	return err
}

func init() {
	RegisterPacket(IDAnimate, func() DataPacket {
		return NewAnimatePacket()
	})
}
