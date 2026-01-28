package protocol

type HurtArmorPacket struct {
	BasePacket
	Health int32
}

func NewHurtArmorPacket() *HurtArmorPacket {
	return &HurtArmorPacket{
		BasePacket: BasePacket{PacketID: IDHurtArmor},
	}
}

func (p *HurtArmorPacket) Name() string {
	return "HurtArmorPacket"
}

func (p *HurtArmorPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteInt(p.Health)
	return nil
}

func (p *HurtArmorPacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())
	var err error
	p.Health, err = stream.ReadInt()
	return err
}

func init() {
	RegisterPacket(IDHurtArmor, func() DataPacket { return NewHurtArmorPacket() })
}
