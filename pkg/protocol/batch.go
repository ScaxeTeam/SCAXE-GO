package protocol

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/logger"
)

type BatchPacket struct {
	BasePacket
	Payload []byte
}

func NewBatchPacket() *BatchPacket {
	return &BatchPacket{
		BasePacket: BasePacket{PacketID: IDBatch},
	}
}

func (p *BatchPacket) Name() string {
	return "BatchPacket"
}

func (p *BatchPacket) Encode(stream *BinaryStream) error {

	EncodeHeader(stream, p.ID())
	stream.WriteInt(int32(len(p.Payload)))
	stream.WriteBytes(p.Payload)
	return nil
}

func (p *BatchPacket) Decode(stream *BinaryStream) error {

	length, err := stream.ReadInt()
	if err != nil {
		return fmt.Errorf("failed to read batch length: %w", err)
	}

	if length <= 0 {
		return nil
	}

	data, err := stream.ReadBytes(int(length))
	if err != nil {
		return fmt.Errorf("failed to read batch payload (len=%d): %w", length, err)
	}
	p.Payload = data
	return nil
}

func CreateBatch(packets []DataPacket) ([]byte, error) {
	var rawBuf bytes.Buffer

	for _, pkt := range packets {

		stream := NewBinaryStream()

		err := pkt.Encode(stream)
		if err != nil {
			logger.Error("CreateBatch", "error", err, "packet", pkt.Name())
			continue
		}
		pktData := stream.Bytes()

		binary.Write(&rawBuf, binary.BigEndian, uint32(len(pktData)))

		rawBuf.Write(pktData)
	}

	var compressedBuf bytes.Buffer

	w := zlib.NewWriter(&compressedBuf)
	_, err := w.Write(rawBuf.Bytes())
	if err != nil {
		w.Close()
		return nil, err
	}
	w.Close()

	return compressedBuf.Bytes(), nil
}

func (p *BatchPacket) Compress(data []byte) error {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	if _, err := w.Write(data); err != nil {
		w.Close()
		return err
	}
	w.Close()
	p.Payload = b.Bytes()
	return nil
}

func (p *BatchPacket) Decompress() ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(p.Payload))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var b bytes.Buffer
	if _, err := b.ReadFrom(r); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func DecodePackets(data []byte) ([]DataPacket, error) {
	var packets []DataPacket
	buf := bytes.NewBuffer(data)

	for buf.Len() > 0 {

		if buf.Len() < 4 {
			break
		}
		var length uint32
		if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
			return nil, err
		}

		if int64(length) > int64(buf.Len()) {
			break
		}

		pktData := make([]byte, length)
		if _, err := buf.Read(pktData); err != nil {
			return nil, err
		}

		if len(pktData) == 0 {
			continue
		}
		pid := pktData[0]

		pkt := GetPacket(pid)
		if pkt == nil {

			continue
		}

		stream := NewBinaryStreamFromBytes(pktData[1:])
		if err := pkt.Decode(stream); err != nil {
			logger.Error("DecodePackets", "error", err, "packetID", pid)
			continue
		}

		packets = append(packets, pkt)
	}

	return packets, nil
}

func init() {
	RegisterPacket(IDBatch, func() DataPacket { return NewBatchPacket() })
}
