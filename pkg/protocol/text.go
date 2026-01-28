package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/logger"
)

const (
	TextTypeRaw         byte = 0
	TextTypeChat        byte = 1
	TextTypeTranslation byte = 2
	TextTypePopup       byte = 3
	TextTypeTip         byte = 4
	TextTypeSystem      byte = 5
)

type TextPacket struct {
	BasePacket
	TextType   byte
	SourceName string
	Message    string
	Parameters []string
	XUID       string
}

func NewTextPacket() *TextPacket {
	return &TextPacket{
		BasePacket: BasePacket{PacketID: IDText},
		Parameters: make([]string, 0),
	}
}

func (p *TextPacket) Name() string {
	return "TextPacket"
}

func (p *TextPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	logger.Debug("TextPacket.Encode", "type", p.TextType, "source", p.SourceName, "message", p.Message)

	stream.WriteByte(p.TextType)

	switch p.TextType {
	case TextTypeChat:
		stream.WriteString16(p.SourceName)
		fallthrough
	case TextTypeRaw, TextTypeTip, TextTypeSystem:
		stream.WriteString16(p.Message)
	case TextTypeTranslation, TextTypePopup:
		stream.WriteString16(p.Message)

		stream.WriteByte(byte(len(p.Parameters)))
		for _, param := range p.Parameters {
			stream.WriteString16(param)
		}
	}

	return nil
}

func (p *TextPacket) Decode(stream *BinaryStream) error {
	var err error

	p.TextType, err = stream.ReadByte()
	if err != nil {
		logger.Error("TextPacket.Decode", "error", "failed to read type", "err", err)
		return err
	}

	switch p.TextType {
	case TextTypeChat:
		p.SourceName, err = stream.ReadString16()
		if err != nil {
			return err
		}
		fallthrough
	case TextTypeRaw, TextTypeTip, TextTypeSystem:
		p.Message, err = stream.ReadString16()
		if err != nil {
			return err
		}
	case TextTypeTranslation, TextTypePopup:
		p.Message, err = stream.ReadString16()
		if err != nil {
			return err
		}

		countByte, err := stream.ReadByte()
		if err != nil {
			return err
		}
		count := uint32(countByte)
		p.Parameters = make([]string, count)
		for i := uint32(0); i < count; i++ {
			p.Parameters[i], err = stream.ReadString16()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func init() {
	RegisterPacket(IDText, func() DataPacket { return NewTextPacket() })
}
