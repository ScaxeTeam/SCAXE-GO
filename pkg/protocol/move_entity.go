package protocol

type MoveEntityPacket struct {
	BasePacket
	Entities []MoveEntityEntry
}

type MoveEntityEntry struct {
	EntityID int64
	X        float32
	Y        float32
	Z        float32
	Yaw      float32
	HeadYaw  float32
	Pitch    float32
}

func NewMoveEntityPacket() *MoveEntityPacket {
	return &MoveEntityPacket{
		BasePacket: BasePacket{PacketID: IDMoveEntity},
	}
}

func (p *MoveEntityPacket) ID() byte {
	return IDMoveEntity
}

func (p *MoveEntityPacket) Name() string {
	return "MoveEntityPacket"
}

func (p *MoveEntityPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteInt(int32(len(p.Entities)))
	for _, e := range p.Entities {
		stream.WriteEntityID(e.EntityID)
		stream.WriteFloat(e.X)
		stream.WriteFloat(e.Y)
		stream.WriteFloat(e.Z)
		stream.WriteFloat(e.Yaw)
		stream.WriteFloat(e.HeadYaw)
		stream.WriteFloat(e.Pitch)
	}
	return nil
}

func (p *MoveEntityPacket) Decode(stream *BinaryStream) error {
	count, err := stream.ReadInt()
	if err != nil {
		return err
	}
	p.Entities = make([]MoveEntityEntry, 0, count)
	for i := int32(0); i < count; i++ {
		var e MoveEntityEntry
		e.EntityID, err = stream.ReadEntityID()
		if err != nil {
			return err
		}
		e.X, err = stream.ReadFloat()
		if err != nil {
			return err
		}
		e.Y, err = stream.ReadFloat()
		if err != nil {
			return err
		}
		e.Z, err = stream.ReadFloat()
		if err != nil {
			return err
		}
		e.Yaw, err = stream.ReadFloat()
		if err != nil {
			return err
		}
		e.HeadYaw, err = stream.ReadFloat()
		if err != nil {
			return err
		}
		e.Pitch, err = stream.ReadFloat()
		if err != nil {
			return err
		}
		p.Entities = append(p.Entities, e)
	}
	return nil
}

func init() {
	RegisterPacket(IDMoveEntity, func() DataPacket {
		return NewMoveEntityPacket()
	})
}
