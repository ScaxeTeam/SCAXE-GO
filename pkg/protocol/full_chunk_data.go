package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/logger"
)

const (
	ChunkOrderColumns byte = 0
	ChunkOrderLayered byte = 1
)

type FullChunkDataPacket struct {
	BasePacket
	ChunkX int32
	ChunkZ int32
	Order  byte
	Data   []byte
}

func NewFullChunkDataPacket() *FullChunkDataPacket {
	return &FullChunkDataPacket{
		BasePacket: BasePacket{PacketID: IDFullChunkData},
		Order:      ChunkOrderColumns,
		Data:       make([]byte, 0),
	}
}

func (p *FullChunkDataPacket) Name() string {
	return "FullChunkDataPacket"
}

func (p *FullChunkDataPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	logger.Debug("FullChunkDataPacket.Encode", "x", p.ChunkX, "z", p.ChunkZ, "order", p.Order, "len", len(p.Data))
	stream.WriteInt(p.ChunkX)
	stream.WriteInt(p.ChunkZ)
	stream.WriteByte(p.Order)
	stream.WriteInt(int32(len(p.Data)))
	stream.WriteBytes(p.Data)
	return nil
}

func (p *FullChunkDataPacket) Decode(stream *BinaryStream) error {
	var err error
	p.ChunkX, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.ChunkZ, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.Order, err = stream.ReadByte()
	if err != nil {
		return err
	}
	length, err := stream.ReadInt()
	if err != nil {
		return err
	}
	p.Data, err = stream.ReadBytes(int(length))
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RegisterPacket(IDFullChunkData, func() DataPacket { return NewFullChunkDataPacket() })
}
