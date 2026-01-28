package protocol

type SetHealthPacket struct {
	BasePacket
	Health int32

}

func NewSetHealthPacket() *SetHealthPacket {
	return &SetHealthPacket{
		BasePacket: BasePacket{PacketID: IDSetHealth},
		Health:     20,
	}
}

func (p *SetHealthPacket) Name() string {
	return "SetHealthPacket"
}

func (p *SetHealthPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())

	stream.WriteInt(p.Health)
	return nil
}

func (p *SetHealthPacket) Decode(stream *BinaryStream) error {
	var err error
	p.Health, err = stream.ReadInt()
	return err
}

func init() {
	RegisterPacket(IDSetHealth, func() DataPacket { return NewSetHealthPacket() })
}
