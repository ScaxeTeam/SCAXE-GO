package protocol

type ContainerOpenPacket struct {
	BasePacket
	WindowID byte
	Type     byte
	Slots    int16
	X        int32
	Y        int32
	Z        int32
	EntityID int64
}

const (
	ContainerTypeChest     = 0
	ContainerTypeWorkbench = 1
	ContainerTypeFurnace   = 2
	ContainerTypeEnchant   = 3
	ContainerTypeBrewing   = 4
	ContainerTypeAnvil     = 5
	ContainerTypeDispenser = 6
	ContainerTypeDropper   = 7
	ContainerTypeHopper    = 8
)

func NewContainerOpenPacket() *ContainerOpenPacket {
	return &ContainerOpenPacket{
		BasePacket: BasePacket{PacketID: IDContainerOpen},
		EntityID:   -1,
	}
}

func (p *ContainerOpenPacket) ID() byte {
	return IDContainerOpen
}

func (p *ContainerOpenPacket) Name() string {
	return "ContainerOpenPacket"
}

func (p *ContainerOpenPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteByte(p.WindowID)
	stream.WriteByte(p.Type)
	stream.WriteShort(p.Slots)
	stream.WriteInt(p.X)
	stream.WriteInt(p.Y)
	stream.WriteInt(p.Z)
	stream.WriteEntityID(p.EntityID)
	return nil
}

func (p *ContainerOpenPacket) Decode(stream *BinaryStream) error {
	var err error
	p.WindowID, err = stream.ReadByte()
	if err != nil {
		return err
	}
	p.Type, err = stream.ReadByte()
	if err != nil {
		return err
	}
	p.Slots, err = stream.ReadShort()
	if err != nil {
		return err
	}
	p.X, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.Y, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.Z, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.EntityID, err = stream.ReadEntityID()
	return err
}

func init() {
	RegisterPacket(IDContainerOpen, func() DataPacket {
		return NewContainerOpenPacket()
	})
}
