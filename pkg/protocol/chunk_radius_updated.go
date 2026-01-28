package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/logger"
)

type ChunkRadiusUpdatedPacket struct {
	BasePacket
	Radius int32
}

func NewChunkRadiusUpdatedPacket() *ChunkRadiusUpdatedPacket {
	return &ChunkRadiusUpdatedPacket{
		BasePacket: BasePacket{PacketID: IDChunkRadiusUpdated},
	}
}

func (p *ChunkRadiusUpdatedPacket) Name() string {
	return "ChunkRadiusUpdatedPacket"
}

func (p *ChunkRadiusUpdatedPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	logger.Debug("ChunkRadiusUpdatedPacket.Encode", "radius", p.Radius)
	stream.WriteInt(p.Radius)
	return nil
}

func (p *ChunkRadiusUpdatedPacket) Decode(stream *BinaryStream) error {
	var err error
	p.Radius, err = stream.ReadInt()
	if err != nil {
		logger.Error("ChunkRadiusUpdatedPacket.Decode", "error", err)
		return err
	}
	logger.Debug("ChunkRadiusUpdatedPacket.Decode", "radius", p.Radius)
	return nil
}

func init() {
	RegisterPacket(IDChunkRadiusUpdated, func() DataPacket { return NewChunkRadiusUpdatedPacket() })
}
