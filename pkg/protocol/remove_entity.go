package protocol

type RemoveEntityPacket struct {
	BasePacket
	EntityID int64
}

func NewRemoveEntityPacket() *RemoveEntityPacket {
	return &RemoveEntityPacket{
		BasePacket: BasePacket{PacketID: IDRemoveEntity},
	}
}

func (p *RemoveEntityPacket) Name() string {
	return "RemoveEntityPacket"
}

func (p *RemoveEntityPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteEntityID(p.EntityID)
	return nil
}

func (p *RemoveEntityPacket) Decode(stream *BinaryStream) error {
	var err error
	p.EntityID, err = stream.ReadEntityID()
	return err
}

func init() {
	RegisterPacket(IDRemoveEntity, func() DataPacket {
		return NewRemoveEntityPacket()
	})
}
