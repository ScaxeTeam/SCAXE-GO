package protocol

type SetTimePacket struct {
	BasePacket
	Time    int32
	Started bool
}

func NewSetTimePacket() *SetTimePacket {
	return &SetTimePacket{
		BasePacket: BasePacket{PacketID: IDSetTime},
		Time:       0,
		Started:    true,
	}
}

func (p *SetTimePacket) Name() string {
	return "SetTimePacket"
}

func (p *SetTimePacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteInt(p.Time)
	stream.WriteBool(p.Started)
	return nil
}

func (p *SetTimePacket) Decode(stream *BinaryStream) error {
	var err error
	p.Time, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.Started, err = stream.ReadBool()
	return err
}

func init() {
	RegisterPacket(IDSetTime, func() DataPacket { return NewSetTimePacket() })
}
