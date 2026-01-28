package protocol

type SetEntityDataPacket struct {
	BasePacket
	EntityID int64
	Metadata []byte
}

func NewSetEntityDataPacket() *SetEntityDataPacket {
	return &SetEntityDataPacket{
		BasePacket: BasePacket{PacketID: IDSetEntityData},
		Metadata:   []byte{0x7f},
	}
}

func (p *SetEntityDataPacket) ID() byte {
	return IDSetEntityData
}

func (p *SetEntityDataPacket) Name() string {
	return "SetEntityDataPacket"
}

func (p *SetEntityDataPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteEntityID(p.EntityID)
	stream.WriteBytes(p.Metadata)
	return nil
}

func (p *SetEntityDataPacket) Decode(stream *BinaryStream) error {
	var err error
	p.EntityID, err = stream.ReadEntityID()

	return err
}

func init() {
	RegisterPacket(IDSetEntityData, func() DataPacket {
		return NewSetEntityDataPacket()
	})
}
