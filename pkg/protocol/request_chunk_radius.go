package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/logger"
)

type RequestChunkRadiusPacket struct {
	BasePacket
	Radius int32
}

func NewRequestChunkRadiusPacket() *RequestChunkRadiusPacket {
	return &RequestChunkRadiusPacket{
		BasePacket: BasePacket{PacketID: IDRequestChunkRadius},
	}
}

func (p *RequestChunkRadiusPacket) Name() string {
	return "RequestChunkRadiusPacket"
}

func (p *RequestChunkRadiusPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteInt(p.Radius)
	return nil
}

func (p *RequestChunkRadiusPacket) Decode(stream *BinaryStream) error {
	var err error
	p.Radius, err = stream.ReadInt()
	if err != nil {
		logger.Error("RequestChunkRadiusPacket.Decode", "error", err)
		return err
	}
	logger.Debug("RequestChunkRadiusPacket.Decode", "radius", p.Radius)
	return nil
}

func init() {
	RegisterPacket(IDRequestChunkRadius, func() DataPacket { return NewRequestChunkRadiusPacket() })
}
