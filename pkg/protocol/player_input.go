package protocol

type PlayerInputPacket struct {
	BasePacket
	MotionX  float32
	MotionY  float32
	Jumping  bool
	Sneaking bool
}

func NewPlayerInputPacket() *PlayerInputPacket {
	return &PlayerInputPacket{
		BasePacket: BasePacket{PacketID: IDPlayerInput},
	}
}

func (p *PlayerInputPacket) Name() string {
	return "PlayerInputPacket"
}

func (p *PlayerInputPacket) Encode(stream *BinaryStream) error {

	return nil
}

func (p *PlayerInputPacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())
	var err error
	p.MotionX, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.MotionY, err = stream.ReadFloat()
	if err != nil {
		return err
	}

	jumping, err := stream.ReadByte()
	if err != nil {
		return err
	}
	p.Jumping = jumping != 0

	sneaking, err := stream.ReadByte()
	if err != nil {
		return err
	}
	p.Sneaking = sneaking != 0

	return nil
}

func init() {
	RegisterPacket(IDPlayerInput, func() DataPacket { return NewPlayerInputPacket() })
}
