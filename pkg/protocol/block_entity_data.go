package protocol

type BlockEntityDataPacket struct {
	BasePacket
	X       int32
	Y       int32
	Z       int32
	NBTData []byte
}

func NewBlockEntityDataPacket() *BlockEntityDataPacket {
	return &BlockEntityDataPacket{
		BasePacket: BasePacket{PacketID: IDBlockEntityData},
	}
}

func (p *BlockEntityDataPacket) Name() string {
	return "BlockEntityDataPacket"
}

func (p *BlockEntityDataPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteInt(p.X)
	stream.WriteInt(p.Y)
	stream.WriteInt(p.Z)
	stream.WriteBytes(p.NBTData)
	return nil
}

func (p *BlockEntityDataPacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())
	var err error
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

	p.NBTData, _ = stream.ReadRemaining()
	return nil
}

func init() {
	RegisterPacket(IDBlockEntityData, func() DataPacket { return NewBlockEntityDataPacket() })
}
