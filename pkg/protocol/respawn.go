package protocol

type RespawnPacket struct {
	BasePacket
	X float32
	Y float32
	Z float32
}

func NewRespawnPacket() *RespawnPacket {
	return &RespawnPacket{
		BasePacket: BasePacket{PacketID: IDRespawn},
	}
}

func (p *RespawnPacket) ID() byte {
	return IDRespawn
}

func (p *RespawnPacket) Name() string {
	return "RespawnPacket"
}

func (p *RespawnPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteFloat(p.X)
	stream.WriteFloat(p.Y)
	stream.WriteFloat(p.Z)
	return nil
}

func (p *RespawnPacket) Decode(stream *BinaryStream) error {
	var err error
	p.X, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.Y, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.Z, err = stream.ReadFloat()
	return err
}

func init() {
	RegisterPacket(IDRespawn, func() DataPacket {
		return NewRespawnPacket()
	})
}
