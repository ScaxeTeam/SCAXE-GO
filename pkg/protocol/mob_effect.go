package protocol

type MobEffectPacket struct {
	BasePacket
	EntityID  int64
	EventID   byte
	EffectID  byte
	Amplifier byte
	Particles bool
	Duration  int32
}

func NewMobEffectPacket() *MobEffectPacket {
	return &MobEffectPacket{
		BasePacket: BasePacket{PacketID: IDMobEffect},
		Particles:  true,
	}
}

func (p *MobEffectPacket) ID() byte {
	return IDMobEffect
}

func (p *MobEffectPacket) Name() string {
	return "MobEffectPacket"
}

func (p *MobEffectPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteEntityID(p.EntityID)
	stream.WriteByte(p.EventID)
	stream.WriteByte(p.EffectID)
	stream.WriteByte(p.Amplifier)
	stream.WriteBool(p.Particles)
	stream.WriteInt(p.Duration)
	return nil
}

func (p *MobEffectPacket) Decode(stream *BinaryStream) error {
	var err error
	p.EntityID, err = stream.ReadEntityID()
	if err != nil {
		return err
	}
	p.EventID, err = stream.ReadByte()
	if err != nil {
		return err
	}
	p.EffectID, err = stream.ReadByte()
	if err != nil {
		return err
	}
	p.Amplifier, err = stream.ReadByte()
	if err != nil {
		return err
	}

	particles, err := stream.ReadByte()
	if err != nil {
		return err
	}
	p.Particles = particles != 0

	p.Duration, err = stream.ReadInt()
	return err
}

func init() {
	RegisterPacket(IDMobEffect, func() DataPacket {
		return NewMobEffectPacket()
	})
}
