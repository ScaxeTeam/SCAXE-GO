package protocol

import (
	"encoding/binary"
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/logger"
)

type AddPlayerPacket struct {
	BasePacket
	UUID     string
	RawUUID  []byte
	Username string
	EntityID int64
	X        float32
	Y        float32
	Z        float32
	SpeedX   float32
	SpeedY   float32
	SpeedZ   float32
	Yaw      float32
	Pitch    float32
	Metadata map[int]interface{}
}

func NewAddPlayerPacket() *AddPlayerPacket {
	return &AddPlayerPacket{
		BasePacket: BasePacket{PacketID: IDAddPlayer},
		Metadata:   make(map[int]interface{}),
	}
}

func (p *AddPlayerPacket) ID() byte {
	return IDAddPlayer
}

func (p *AddPlayerPacket) Name() string {
	return "AddPlayerPacket"
}

func (p *AddPlayerPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())

	if len(p.RawUUID) == 16 {
		stream.WriteBytes(p.RawUUID)
		logger.Debug("AddPlayerPacket Encode", "RawUUID", fmt.Sprintf("%x", p.RawUUID))
	} else {
		rawUUID := parseUUIDToBytes(p.UUID)
		stream.WriteBytes(rawUUID)
		logger.Debug("AddPlayerPacket Encode", "UUID", p.UUID, "Raw", fmt.Sprintf("%x", rawUUID))
	}

	stream.WriteString16(p.Username)

	stream.WriteEntityID(p.EntityID)

	stream.WriteFloat(p.X)
	stream.WriteFloat(p.Y)
	stream.WriteFloat(p.Z)

	stream.WriteFloat(p.SpeedX)
	stream.WriteFloat(p.SpeedY)
	stream.WriteFloat(p.SpeedZ)

	stream.WriteFloat(p.Yaw)
	stream.WriteFloat(p.Yaw)
	stream.WriteFloat(p.Pitch)

	stream.WriteShort(0)

	metaBytes := encodeMetadata(p.Username)
	stream.WriteBytes(metaBytes)

	return nil
}

func encodeMetadata(username string) []byte {
	var meta []byte

	meta = append(meta, 0x00)
	meta = append(meta, 0)

	meta = append(meta, 0x21)

	airBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(airBytes, 300)
	meta = append(meta, airBytes...)

	meta = append(meta, 0x82)

	nameBytes := []byte(username)
	lenBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(lenBytes, uint16(len(nameBytes)))
	meta = append(meta, lenBytes...)
	meta = append(meta, nameBytes...)

	meta = append(meta, 0x03)
	meta = append(meta, 1)

	meta = append(meta, 0x04)
	meta = append(meta, 0)

	meta = append(meta, 0x10)
	meta = append(meta, 0)

	meta = append(meta, 0xD1)

	posBytes := make([]byte, 12)

	meta = append(meta, posBytes...)

	meta = append(meta, 0xF7)

	longBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(longBytes, 0xFFFFFFFFFFFFFFFF)
	meta = append(meta, longBytes...)

	meta = append(meta, 0x18)
	meta = append(meta, 0)

	meta = append(meta, 0x7f)

	return meta
}

func parseUUIDToBytes(uuid string) []byte {
	result := make([]byte, 16)
	if len(uuid) == 0 {
		return result
	}

	hexStr := ""
	for _, c := range uuid {
		if c != '-' {
			hexStr += string(c)
		}
	}

	for i := 0; i < 16 && i*2+1 < len(hexStr); i++ {
		var b byte
		hexPair := hexStr[i*2 : i*2+2]
		for j := 0; j < 2; j++ {
			c := hexPair[j]
			b <<= 4
			if c >= '0' && c <= '9' {
				b |= c - '0'
			} else if c >= 'a' && c <= 'f' {
				b |= c - 'a' + 10
			} else if c >= 'A' && c <= 'F' {
				b |= c - 'A' + 10
			}
		}
		result[i] = b
	}

	return result
}

func (p *AddPlayerPacket) Decode(stream *BinaryStream) error {

	return nil
}

func init() {
	RegisterPacket(IDAddPlayer, func() DataPacket {
		return NewAddPlayerPacket()
	})
}
