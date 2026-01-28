package protocol

type ContainerSetDataPacket struct {
	BasePacket
	WindowID byte
	Property uint16
	Value    uint16
}

func NewContainerSetDataPacket() *ContainerSetDataPacket {
	return &ContainerSetDataPacket{
		BasePacket: BasePacket{PacketID: IDContainerSetData},
	}
}

func (p *ContainerSetDataPacket) Name() string {
	return "ContainerSetDataPacket"
}

func (p *ContainerSetDataPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteByte(p.WindowID)
	stream.WriteBEUShort(p.Property)
	stream.WriteBEUShort(p.Value)
	return nil
}

func (p *ContainerSetDataPacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())
	var err error
	p.WindowID, err = stream.ReadByte()
	if err != nil {
		return err
	}
	p.Property, err = stream.ReadBEUShort()
	if err != nil {
		return err
	}
	p.Value, err = stream.ReadBEUShort()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RegisterPacket(IDContainerSetData, func() DataPacket { return NewContainerSetDataPacket() })
}
