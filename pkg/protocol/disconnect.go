package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/logger"
)

type DisconnectPacket struct {
	BasePacket
	HideDisconnectionScreen bool
	Message                 string
}

func NewDisconnectPacket() *DisconnectPacket {
	return &DisconnectPacket{
		BasePacket: BasePacket{PacketID: IDDisconnect},
	}
}

func (p *DisconnectPacket) Name() string {
	return "DisconnectPacket"
}

func (p *DisconnectPacket) Encode(stream *BinaryStream) error {
	logger.Debug("DisconnectPacket.Encode", "message", p.Message, "hideScreen", p.HideDisconnectionScreen)
	EncodeHeader(stream, p.ID())
	stream.WriteBool(p.HideDisconnectionScreen)
	if !p.HideDisconnectionScreen {
		stream.WriteString16(p.Message)
	}
	return nil
}

func (p *DisconnectPacket) Decode(stream *BinaryStream) error {
	var err error

	p.HideDisconnectionScreen, err = stream.ReadBool()
	if err != nil {
		logger.Error("DisconnectPacket.Decode", "error", err)
		return err
	}
	if !p.HideDisconnectionScreen {
		p.Message, err = stream.ReadString16()
		if err != nil {
			logger.Error("DisconnectPacket.Decode", "error", err)
			return err
		}
	}
	logger.Debug("DisconnectPacket.Decode", "message", p.Message)
	return nil
}

func init() {
	RegisterPacket(IDDisconnect, func() DataPacket { return NewDisconnectPacket() })
}
