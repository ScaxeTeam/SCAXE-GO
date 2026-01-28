package protocol

import (
	"github.com/google/uuid"
)

type RemovePlayerPacket struct {
	BasePacket
	EntityID int64
	UUID     uuid.UUID
}

func NewRemovePlayerPacket() *RemovePlayerPacket {
	return &RemovePlayerPacket{
		BasePacket: BasePacket{PacketID: IDRemovePlayer},
	}
}

func (p *RemovePlayerPacket) Name() string {
	return "RemovePlayerPacket"
}

func (p *RemovePlayerPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteEntityID(p.EntityID)

	stream.WriteUUID(string(p.UUID[:]))
	return nil
}

func (p *RemovePlayerPacket) Decode(stream *BinaryStream) error {
	var err error
	p.EntityID, err = stream.ReadEntityID()
	if err != nil {
		return err
	}
	uuidStr, err := stream.ReadUUID()
	if err != nil {
		return err
	}

	p.UUID, err = uuid.FromBytes([]byte(uuidStr))
	return err
}

func init() {
	RegisterPacket(IDRemovePlayer, func() DataPacket {
		return NewRemovePlayerPacket()
	})
}
