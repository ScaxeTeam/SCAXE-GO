package protocol

type SetDifficultyPacket struct {
	BasePacket
	Difficulty int32
}

func NewSetDifficultyPacket() *SetDifficultyPacket {
	return &SetDifficultyPacket{
		BasePacket: BasePacket{PacketID: IDSetDifficulty},
		Difficulty: 1,
	}
}

func (p *SetDifficultyPacket) Name() string {
	return "SetDifficultyPacket"
}

func (p *SetDifficultyPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteInt(p.Difficulty)
	return nil
}

func (p *SetDifficultyPacket) Decode(stream *BinaryStream) error {
	var err error
	p.Difficulty, err = stream.ReadInt()
	return err
}

func init() {
	RegisterPacket(IDSetDifficulty, func() DataPacket { return NewSetDifficultyPacket() })
}
