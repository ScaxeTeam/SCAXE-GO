package protocol

type ContainerClosePacket struct {
	BasePacket
	WindowID byte
}

func NewContainerClosePacket() *ContainerClosePacket {
	return &ContainerClosePacket{
		BasePacket: BasePacket{PacketID: IDContainerClose},
	}
}

func (p *ContainerClosePacket) ID() byte {
	return IDContainerClose
}

func (p *ContainerClosePacket) Name() string {
	return "ContainerClosePacket"
}

func (p *ContainerClosePacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteByte(p.WindowID)
	return nil
}

func (p *ContainerClosePacket) Decode(stream *BinaryStream) error {
	var err error
	p.WindowID, err = stream.ReadByte()
	return err
}

func init() {
	RegisterPacket(IDContainerClose, func() DataPacket {
		return NewContainerClosePacket()
	})
}
