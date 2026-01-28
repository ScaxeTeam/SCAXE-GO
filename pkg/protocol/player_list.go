package protocol

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/logger"
)

type PlayerListPacket struct {
	BasePacket
	Type    byte
	Entries []PlayerListEntry
}

const (
	PlayerListTypeAdd    byte = 0
	PlayerListTypeRemove byte = 1
)

type PlayerListEntry struct {
	UUID     string
	EntityID int64
	Username string
	SkinName string
	SkinData string
}

func NewPlayerListPacket() *PlayerListPacket {
	return &PlayerListPacket{
		BasePacket: BasePacket{PacketID: IDPlayerList},
		Entries:    make([]PlayerListEntry, 0),
	}
}

func (p *PlayerListPacket) ID() byte {
	return IDPlayerList
}

func (p *PlayerListPacket) Name() string {
	return "PlayerListPacket"
}

func (p *PlayerListPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteByte(p.Type)
	stream.WriteInt(int32(len(p.Entries)))

	for _, entry := range p.Entries {

		rawUUID := parseUUIDToBytes(entry.UUID)
		logger.Debug("PlayerListPacket Encode", "UUID", entry.UUID, "Raw", fmt.Sprintf("%x", rawUUID), "Type", p.Type)
		stream.WriteBytes(rawUUID)

		if p.Type == PlayerListTypeAdd {

			stream.WriteEntityID(entry.EntityID)

			stream.WriteString16(entry.Username)

			stream.WriteString16(entry.SkinName)

			stream.WriteString16(entry.SkinData)
		}

	}
	return nil
}

func (p *PlayerListPacket) Decode(stream *BinaryStream) error {
	var err error
	p.Type, err = stream.ReadByte()
	if err != nil {
		return err
	}

	count, err := stream.ReadInt()
	if err != nil {
		return err
	}

	p.Entries = make([]PlayerListEntry, count)
	for i := int32(0); i < count; i++ {
		entry := PlayerListEntry{}
		entry.UUID, err = stream.ReadUUID()
		if err != nil {
			return err
		}

		if p.Type == PlayerListTypeAdd {
			entry.EntityID, err = stream.ReadEntityID()
			if err != nil {
				return err
			}
			entry.Username, err = stream.ReadString()
			if err != nil {
				return err
			}
			entry.SkinName, err = stream.ReadString()
			if err != nil {
				return err
			}
			entry.SkinData, err = stream.ReadString()
			if err != nil {
				return err
			}
		}
		p.Entries[i] = entry
	}
	return nil
}

func init() {
	RegisterPacket(IDPlayerList, func() DataPacket {
		return NewPlayerListPacket()
	})
}
