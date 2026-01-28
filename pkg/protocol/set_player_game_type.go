package protocol

func init() {
	RegisterPacket(IDSetPlayerGameType, func() DataPacket { return NewSetPlayerGameTypePacket() })
}

type SetPlayerGameTypePacket struct {
	BasePacket
	Gamemode int32
}

func NewSetPlayerGameTypePacket() *SetPlayerGameTypePacket {
	return &SetPlayerGameTypePacket{
		BasePacket: BasePacket{PacketID: IDSetPlayerGameType},
	}
}

func (p *SetPlayerGameTypePacket) ID() byte {
	return IDSetPlayerGameType
}

func (p *SetPlayerGameTypePacket) Name() string {
	return "SetPlayerGameTypePacket"
}

func (p *SetPlayerGameTypePacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteBEInt(p.Gamemode)
	return nil
}

func (p *SetPlayerGameTypePacket) Decode(stream *BinaryStream) error {
	var err error
	p.Gamemode, err = stream.ReadBEInt()
	return err
}
