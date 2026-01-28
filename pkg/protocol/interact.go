package protocol

const (
	ActionRightClick   byte = 1
	ActionLeftClick    byte = 2
	ActionLeaveVehicle byte = 3
)

type InteractPacket struct {
	BasePacket
	Action byte
	Target int64
}

func NewInteractPacket() *InteractPacket {
	return &InteractPacket{
		BasePacket: BasePacket{PacketID: IDInteract},
	}
}

func (p *InteractPacket) Name() string {
	return "InteractPacket"
}

func (p *InteractPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteByte(p.Action)
	stream.WriteLong(p.Target)
	return nil
}

func (p *InteractPacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())
	var err error
	p.Action, err = stream.ReadByte()
	if err != nil {
		return err
	}
	p.Target, err = stream.ReadLong()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RegisterPacket(IDInteract, func() DataPacket { return NewInteractPacket() })
}
