package protocol

type TelemetryEventPacket struct {
	BasePacket
	EntityID  int64
	EventID   int32
	EventData byte
}

func NewTelemetryEventPacket() *TelemetryEventPacket {
	return &TelemetryEventPacket{
		BasePacket: BasePacket{PacketID: IDTelemetryEvent},
	}
}

func (p *TelemetryEventPacket) Name() string {
	return "TelemetryEventPacket"
}

func (p *TelemetryEventPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteLong(p.EntityID)
	stream.WriteInt(p.EventID)
	stream.WriteByte(p.EventData)
	return nil
}

func (p *TelemetryEventPacket) Decode(stream *BinaryStream) error {
	var err error
	p.EntityID, err = stream.ReadLong()
	if err != nil {
		return err
	}
	p.EventID, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.EventData, err = stream.ReadByte()
	return err
}

func init() {
	RegisterPacket(IDTelemetryEvent, func() DataPacket { return NewTelemetryEventPacket() })
}
