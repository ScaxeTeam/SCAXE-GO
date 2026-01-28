package protocol

const (
	LinkTypeRemove    byte = 0
	LinkTypeRider     byte = 1
	LinkTypePassenger byte = 2
)

type SetEntityLinkPacket struct {
	BasePacket
	From     int64
	To       int64
	LinkType byte
}

func NewSetEntityLinkPacket() *SetEntityLinkPacket {
	return &SetEntityLinkPacket{
		BasePacket: BasePacket{PacketID: IDSetEntityLink},
	}
}

func (p *SetEntityLinkPacket) Name() string {
	return "SetEntityLinkPacket"
}

func (p *SetEntityLinkPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteLong(p.From)
	stream.WriteLong(p.To)
	stream.WriteByte(p.LinkType)
	return nil
}

func (p *SetEntityLinkPacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())
	var err error
	p.From, err = stream.ReadLong()
	if err != nil {
		return err
	}
	p.To, err = stream.ReadLong()
	if err != nil {
		return err
	}
	p.LinkType, err = stream.ReadByte()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RegisterPacket(IDSetEntityLink, func() DataPacket { return NewSetEntityLinkPacket() })
}
