package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/logger"
)

type BinaryStream struct {
	buf    *bytes.Buffer
	offset int
}

func NewBinaryStream() *BinaryStream {
	logger.Debug("NewBinaryStream", "action", "create empty stream")
	return &BinaryStream{
		buf:    new(bytes.Buffer),
		offset: 0,
	}
}

func NewBinaryStreamFromBytes(data []byte) *BinaryStream {
	logger.Debug("NewBinaryStreamFromBytes", "size", len(data))
	return &BinaryStream{
		buf:    bytes.NewBuffer(data),
		offset: 0,
	}
}

func (b *BinaryStream) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *BinaryStream) Len() int {
	return b.buf.Len()
}

func (b *BinaryStream) Reset() {
	b.buf.Reset()
	b.offset = 0
}

func (b *BinaryStream) Feof() bool {
	return b.buf.Len() == 0
}

func (b *BinaryStream) ReadByte() (byte, error) {
	val, err := b.buf.ReadByte()
	if err != nil {
		logger.Error("BinaryStream.ReadByte", "error", err)
		return 0, err
	}
	b.offset++
	return val, nil
}

func (b *BinaryStream) ReadBool() (bool, error) {
	val, err := b.ReadByte()
	if err != nil {
		return false, err
	}
	result := val != 0x00
	logger.Debug("BinaryStream.ReadBool", "value", result, "offset", b.offset)
	return result, nil
}

func (b *BinaryStream) ReadBytes(n int) ([]byte, error) {
	data := make([]byte, n)
	read, err := b.buf.Read(data)
	if err != nil {
		logger.Error("BinaryStream.ReadBytes", "error", err, "requested", n)
		return nil, err
	}
	if read != n {
		return nil, fmt.Errorf("expected %d bytes, got %d", n, read)
	}
	b.offset += n
	return data, nil
}

func (b *BinaryStream) ReadShort() (int16, error) {
	data, err := b.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	val := int16(binary.BigEndian.Uint16(data))
	logger.Debug("BinaryStream.ReadShort", "value", val, "offset", b.offset)
	return val, nil
}

func (b *BinaryStream) ReadLShort() (int16, error) {
	data, err := b.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	val := int16(binary.LittleEndian.Uint16(data))
	logger.Debug("BinaryStream.ReadLShort", "value", val, "offset", b.offset)
	return val, nil
}

func (b *BinaryStream) ReadInt() (int32, error) {
	data, err := b.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	val := int32(binary.BigEndian.Uint32(data))
	logger.Debug("BinaryStream.ReadInt", "value", val, "offset", b.offset)
	return val, nil
}

func (b *BinaryStream) ReadLInt() (int32, error) {
	data, err := b.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	val := int32(binary.LittleEndian.Uint32(data))
	logger.Debug("BinaryStream.ReadLInt", "value", val, "offset", b.offset)
	return val, nil
}

func (b *BinaryStream) ReadLong() (int64, error) {
	data, err := b.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	val := int64(binary.BigEndian.Uint64(data))
	logger.Debug("BinaryStream.ReadLong", "value", val, "offset", b.offset)
	return val, nil
}

func (b *BinaryStream) ReadLLong() (int64, error) {
	data, err := b.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	val := int64(binary.LittleEndian.Uint64(data))
	logger.Debug("BinaryStream.ReadLLong", "value", val, "offset", b.offset)
	return val, nil
}

func (b *BinaryStream) ReadFloat() (float32, error) {
	data, err := b.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	bits := binary.BigEndian.Uint32(data)
	val := math.Float32frombits(bits)
	logger.Debug("BinaryStream.ReadFloat", "value", val, "offset", b.offset)
	return val, nil
}

func (b *BinaryStream) ReadLFloat() (float32, error) {
	data, err := b.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint32(data)
	val := math.Float32frombits(bits)
	logger.Debug("BinaryStream.ReadLFloat", "value", val, "offset", b.offset)
	return val, nil
}

func (b *BinaryStream) ReadDouble() (float64, error) {
	data, err := b.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	bits := binary.BigEndian.Uint64(data)
	val := math.Float64frombits(bits)
	logger.Debug("BinaryStream.ReadDouble", "value", val, "offset", b.offset)
	return val, nil
}

func (b *BinaryStream) ReadLDouble() (float64, error) {
	data, err := b.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint64(data)
	val := math.Float64frombits(bits)
	logger.Debug("BinaryStream.ReadLDouble", "value", val, "offset", b.offset)
	return val, nil
}

func (b *BinaryStream) ReadVarInt() (int32, error) {
	raw, err := b.ReadUnsignedVarInt()
	if err != nil {
		return 0, err
	}

	val := int32((raw >> 1) ^ -(raw & 1))
	logger.Debug("BinaryStream.ReadVarInt", "value", val, "offset", b.offset)
	return val, nil
}

func (b *BinaryStream) ReadUnsignedVarInt() (uint32, error) {
	var val uint32
	for i := 0; i < 5; i++ {
		byteVal, err := b.ReadByte()
		if err != nil {
			return 0, err
		}
		val |= uint32(byteVal&0x7F) << (i * 7)
		if (byteVal & 0x80) == 0 {
			logger.Debug("BinaryStream.ReadUnsignedVarInt", "value", val, "bytes", i+1)
			return val, nil
		}
	}
	return 0, fmt.Errorf("VarInt is too big")
}

func (b *BinaryStream) ReadVarLong() (int64, error) {
	raw, err := b.ReadUnsignedVarLong()
	if err != nil {
		return 0, err
	}

	val := int64((raw >> 1) ^ -(raw & 1))
	logger.Debug("BinaryStream.ReadVarLong", "value", val, "offset", b.offset)
	return val, nil
}

func (b *BinaryStream) ReadUnsignedVarLong() (uint64, error) {
	var val uint64
	for i := 0; i < 10; i++ {
		byteVal, err := b.ReadByte()
		if err != nil {
			return 0, err
		}
		val |= uint64(byteVal&0x7F) << (i * 7)
		if (byteVal & 0x80) == 0 {
			logger.Debug("BinaryStream.ReadUnsignedVarLong", "value", val, "bytes", i+1)
			return val, nil
		}
	}
	return 0, fmt.Errorf("VarLong is too big")
}

func (b *BinaryStream) ReadString() (string, error) {
	length, err := b.ReadUnsignedVarInt()
	if err != nil {
		return "", err
	}
	if length == 0 {
		return "", nil
	}
	data, err := b.ReadBytes(int(length))
	if err != nil {
		return "", err
	}
	val := string(data)
	logger.Debug("BinaryStream.ReadString", "length", length, "value", val)
	return val, nil
}

func (b *BinaryStream) ReadRemaining() ([]byte, error) {
	data, err := io.ReadAll(b.buf)
	if err != nil {
		return nil, err
	}
	b.offset += len(data)
	logger.Debug("BinaryStream.ReadRemaining", "size", len(data))
	return data, nil
}

func (b *BinaryStream) WriteByte(val byte) error {
	b.buf.WriteByte(val)
	b.offset++
	logger.Debug("BinaryStream.WriteByte", "value", val)
	return nil
}

func (b *BinaryStream) WriteBool(val bool) {
	if val {
		b.WriteByte(0x01)
	} else {
		b.WriteByte(0x00)
	}
	logger.Debug("BinaryStream.WriteBool", "value", val)
}

func (b *BinaryStream) WriteBytes(data []byte) {
	b.buf.Write(data)
	b.offset += len(data)
}

func (b *BinaryStream) WriteShort(val int16) {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, uint16(val))
	b.WriteBytes(data)
	logger.Debug("BinaryStream.WriteShort", "value", val)
}

func (b *BinaryStream) WriteLShort(val int16) {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(val))
	b.WriteBytes(data)
	logger.Debug("BinaryStream.WriteLShort", "value", val)
}

func (b *BinaryStream) WriteInt(val int32) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(val))
	b.WriteBytes(data)
	logger.Debug("BinaryStream.WriteInt", "value", val)
}

func (b *BinaryStream) WriteLInt(val int32) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(val))
	b.WriteBytes(data)
	logger.Debug("BinaryStream.WriteLInt", "value", val)
}

func (b *BinaryStream) WriteLong(val int64) {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(val))
	b.WriteBytes(data)
	logger.Debug("BinaryStream.WriteLong", "value", val)
}

func (b *BinaryStream) WriteLLong(val int64) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(val))
	b.WriteBytes(data)
	logger.Debug("BinaryStream.WriteLLong", "value", val)
}

func (b *BinaryStream) WriteFloat(val float32) {
	bits := math.Float32bits(val)
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, bits)
	b.WriteBytes(data)
	logger.Debug("BinaryStream.WriteFloat", "value", val)
}

func (b *BinaryStream) WriteLFloat(val float32) {
	bits := math.Float32bits(val)
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, bits)
	b.WriteBytes(data)
	logger.Debug("BinaryStream.WriteLFloat", "value", val)
}

func (b *BinaryStream) WriteDouble(val float64) {
	bits := math.Float64bits(val)
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, bits)
	b.WriteBytes(data)
	logger.Debug("BinaryStream.WriteDouble", "value", val)
}

func (b *BinaryStream) WriteLDouble(val float64) {
	bits := math.Float64bits(val)
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, bits)
	b.WriteBytes(data)
	logger.Debug("BinaryStream.WriteLDouble", "value", val)
}

func (b *BinaryStream) WriteVarInt(val int32) {

	encoded := uint32((val << 1) ^ (val >> 31))
	b.WriteUnsignedVarInt(encoded)
	logger.Debug("BinaryStream.WriteVarInt", "value", val)
}

func (b *BinaryStream) WriteUnsignedVarInt(val uint32) {
	for {
		if (val & ^uint32(0x7F)) == 0 {
			b.WriteByte(byte(val))
			break
		}
		b.WriteByte(byte((val & 0x7F) | 0x80))
		val >>= 7
	}
}

func (b *BinaryStream) WriteVarLong(val int64) {

	encoded := uint64((val << 1) ^ (val >> 63))
	b.WriteUnsignedVarLong(encoded)
	logger.Debug("BinaryStream.WriteVarLong", "value", val)
}

func (b *BinaryStream) WriteUnsignedVarLong(val uint64) {
	for {
		if (val & ^uint64(0x7F)) == 0 {
			b.WriteByte(byte(val))
			break
		}
		b.WriteByte(byte((val & 0x7F) | 0x80))
		val >>= 7
	}
}

func (b *BinaryStream) WriteString(val string) {
	b.WriteUnsignedVarInt(uint32(len(val)))
	b.WriteBytes([]byte(val))
	logger.Debug("BinaryStream.WriteString", "length", len(val), "value", val)
}

func (b *BinaryStream) ReadString16() (string, error) {
	length, err := b.ReadShort()
	if err != nil {
		return "", err
	}
	if length <= 0 {
		return "", nil
	}
	data, err := b.ReadBytes(int(length))
	if err != nil {
		return "", err
	}
	val := string(data)
	logger.Debug("BinaryStream.ReadString16", "length", length, "value", val)
	return val, nil
}

func (b *BinaryStream) WriteString16(val string) {
	b.WriteShort(int16(len(val)))
	b.WriteBytes([]byte(val))
	logger.Debug("BinaryStream.WriteString16", "length", len(val), "value", val)
}

func (b *BinaryStream) ReadEntityID() (int64, error) {
	return b.ReadLong()
}

func (b *BinaryStream) WriteEntityID(val int64) {
	b.WriteLong(val)
}

func (b *BinaryStream) ReadUUID() (string, error) {
	data, err := b.ReadBytes(16)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (b *BinaryStream) WriteUUID(uuid string) {
	data := make([]byte, 16)
	copy(data, []byte(uuid))
	b.WriteBytes(data)
}

func (b *BinaryStream) ReadBELong() (int64, error) {
	return b.ReadLong()
}

func (b *BinaryStream) WriteBELong(val int64) {
	b.WriteLong(val)
}

func (b *BinaryStream) ReadBEShort() (int16, error) {
	return b.ReadShort()
}

func (b *BinaryStream) WriteBEShort(val int16) {
	b.WriteShort(val)
}

func (b *BinaryStream) ReadBEUShort() (uint16, error) {
	data, err := b.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(data), nil
}

func (b *BinaryStream) WriteBEUShort(val uint16) {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, val)
	b.WriteBytes(data)
}

func (b *BinaryStream) ReadBEInt() (int32, error) {
	return b.ReadInt()
}

func (b *BinaryStream) WriteBEInt(val int32) {
	b.WriteInt(val)
}

func (b *BinaryStream) ReadLEShort() (int16, error) {
	return b.ReadLShort()
}

func (b *BinaryStream) WriteLEShort(val int16) {
	b.WriteLShort(val)
}

func (b *BinaryStream) Skip(n int) error {
	if n <= 0 {
		return nil
	}
	_, err := b.ReadBytes(n)
	return err
}

func (b *BinaryStream) WriteSlot(it item.Item) {

	if it.ID <= 0 {
		b.WriteBEUShort(0)
		return
	}

	b.WriteBEUShort(uint16(it.ID))

	count := it.Count
	if count <= 0 {
		count = 1
	}
	b.WriteByte(byte(count))

	b.WriteBEUShort(uint16(it.Meta))

	nbtBytes := it.GetCompoundTagBytes()
	b.WriteLEShort(int16(len(nbtBytes)))
	if len(nbtBytes) > 0 {
		b.WriteBytes(nbtBytes)
	}
}

func (b *BinaryStream) ReadSlot() (item.Item, error) {

	id, err := b.ReadShort()
	if err != nil {
		return item.Item{}, err
	}
	if id <= 0 {
		return item.Item{ID: 0}, nil
	}

	countByte, err := b.ReadByte()
	if err != nil {
		return item.Item{}, err
	}
	count := int(int8(countByte))
	if count < 0 {
		count = 0
	}

	meta, err := b.ReadBEUShort()
	if err != nil {
		return item.Item{}, err
	}

	nbtLen, err := b.ReadLEShort()
	if err != nil {
		return item.Item{}, err
	}

	var nbtData []byte
	if nbtLen > 0 {
		nbtData, err = b.ReadBytes(int(nbtLen))
		if err != nil {
			return item.Item{}, err
		}
	}
	_ = nbtData

	return item.Item{
		ID:    int(id),
		Count: count,
		Meta:  int(meta),
	}, nil
}
