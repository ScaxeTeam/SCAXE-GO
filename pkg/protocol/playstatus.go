package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/logger"
)

const (
	PlayStatusLoginSuccess               int32 = 0
	PlayStatusLoginFailedClient          int32 = 1
	PlayStatusLoginFailedServer          int32 = 2
	PlayStatusPlayerSpawn                int32 = 3
	PlayStatusLoginFailedInvalidTenant   int32 = 4
	PlayStatusLoginFailedEditionMismatch int32 = 5
)

type PlayStatusPacket struct {
	BasePacket
	Status int32
}

func NewPlayStatusPacket() *PlayStatusPacket {
	return &PlayStatusPacket{
		BasePacket: BasePacket{PacketID: IDPlayStatus},
	}
}

func (p *PlayStatusPacket) Name() string {
	return "PlayStatusPacket"
}

func (p *PlayStatusPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	logger.Debug("PlayStatusPacket.Encode", "status", p.Status)
	stream.WriteInt(p.Status)
	return nil
}

func (p *PlayStatusPacket) Decode(stream *BinaryStream) error {
	var err error
	p.Status, err = stream.ReadInt()
	if err != nil {
		logger.Error("PlayStatusPacket.Decode", "error", err)
		return err
	}
	logger.Debug("PlayStatusPacket.Decode", "status", p.Status)
	return nil
}

func init() {
	RegisterPacket(IDPlayStatus, func() DataPacket { return NewPlayStatusPacket() })
}
