package protocol

type PlayerActionPacket struct {
	BasePacket
	EntityID int64
	Action   int32
	X        int32
	Y        int32
	Z        int32
	Face     int32
}

const (
	ActionStartBreak         int32 = 0
	ActionAbortBreak         int32 = 1
	ActionStopBreak          int32 = 2
	ActionReleaseItem        int32 = 5
	ActionStopSleeping       int32 = 6
	ActionRespawn            int32 = 7
	ActionJump               int32 = 8
	ActionStartSprint        int32 = 9
	ActionStopSprint         int32 = 10
	ActionStartSneak         int32 = 11
	ActionStopSneak          int32 = 12
	ActionDimensionChange    int32 = 13
	ActionDimensionChangeAck int32 = 14
)

func NewPlayerActionPacket() *PlayerActionPacket {
	return &PlayerActionPacket{
		BasePacket: BasePacket{PacketID: IDPlayerAction},
	}
}

func (p *PlayerActionPacket) ID() byte {
	return IDPlayerAction
}

func (p *PlayerActionPacket) Name() string {
	return "PlayerActionPacket"
}

func (p *PlayerActionPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteEntityID(p.EntityID)
	stream.WriteInt(p.Action)
	stream.WriteInt(p.X)
	stream.WriteInt(p.Y)
	stream.WriteInt(p.Z)
	stream.WriteInt(p.Face)
	return nil
}

func (p *PlayerActionPacket) Decode(stream *BinaryStream) error {
	var err error
	p.EntityID, err = stream.ReadEntityID()
	if err != nil {
		return err
	}
	p.Action, err = stream.ReadInt()
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
	p.Face, err = stream.ReadInt()
	return err
}

func init() {
	RegisterPacket(IDPlayerAction, func() DataPacket {
		return NewPlayerActionPacket()
	})
}
