package protocol

import (
	"encoding/binary"
	"math"
)

const (
	DataTypeByte   = 0
	DataTypeShort  = 1
	DataTypeInt    = 2
	DataTypeFloat  = 3
	DataTypeString = 4
	DataTypeSlot   = 5
	DataTypePos    = 6
	DataTypeLong   = 7
)

const (
	DataFlags         = 0
	DataAir           = 1
	DataNametag       = 2
	DataShowNametag   = 3
	DataNoAI          = 15
	DataLeadHolder    = 23
	DataLeadHolderEID = 23
)

const (
	DataFlagOnFire    = 0
	DataFlagSneaking  = 1
	DataFlagRiding    = 2
	DataFlagSprinting = 3
	DataFlagAction    = 4
	DataFlagInvisible = 5
)

type MetadataMap map[int]interface{}

func NewMetadataMap() MetadataMap {
	return make(MetadataMap)
}

func (m MetadataMap) PutByte(key int, value byte) {
	m[key] = value
}

func (m MetadataMap) PutShort(key int, value int16) {
	m[key] = value
}

func (m MetadataMap) PutInt(key int, value int32) {
	m[key] = value
}

func (m MetadataMap) PutFloat(key int, value float32) {
	m[key] = value
}

func (m MetadataMap) PutString(key int, value string) {
	m[key] = value
}

func (m MetadataMap) PutLong(key int, value int64) {
	m[key] = value
}

func EncodeMetadata(m MetadataMap) []byte {
	var buf []byte

	for key, value := range m {

		if key >= 32 {
			continue
		}

		var typeID byte
		switch value.(type) {
		case byte:
			typeID = DataTypeByte
		case int16:
			typeID = DataTypeShort
		case int32:
			typeID = DataTypeInt
		case float32:
			typeID = DataTypeFloat
		case string:
			typeID = DataTypeString
		case int64:
			typeID = DataTypeLong

		default:
			continue
		}

		header := (typeID << 5) | byte(key&0x1F)
		buf = append(buf, header)

		switch v := value.(type) {
		case byte:
			buf = append(buf, v)
		case int16:
			b := make([]byte, 2)
			binary.LittleEndian.PutUint16(b, uint16(v))
			buf = append(buf, b...)
		case int32:
			b := make([]byte, 4)
			binary.LittleEndian.PutUint32(b, uint32(v))
			buf = append(buf, b...)
		case float32:
			b := make([]byte, 4)
			binary.LittleEndian.PutUint32(b, math.Float32bits(v))
			buf = append(buf, b...)
		case string:
			strBytes := []byte(v)
			lenBytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(lenBytes, uint16(len(strBytes)))
			buf = append(buf, lenBytes...)
			buf = append(buf, strBytes...)
		case int64:
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(v))
			buf = append(buf, b...)
		}
	}

	buf = append(buf, 0x7f)

	return buf
}
