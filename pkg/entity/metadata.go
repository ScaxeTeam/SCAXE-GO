package entity

import (
	"bytes"
	"encoding/binary"
	_ "math"
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
	DataFlags              = 0
	DataAir                = 1
	DataNameTag            = 2
	DataShowNameTag        = 3
	DataSilent             = 4
	DataPotionColor        = 7
	DataPotionAmbient      = 8
	DataZombieIsBaby       = 12
	DataAgeableFlags       = 14
	DataAnimalFlags        = 14
	DataIsBaby             = 14
	DataNoAI               = 15
	DataProfessionID       = 16
	DataPotionID           = 16
	DataSlimeSize          = 16
	DataColorInfo          = 16
	DataIsResting          = 16
	DataCharge             = 16
	DataPlayerFlags        = 16
	DataHurtTime           = 17
	DataPlayerBedPos       = 17
	DataShooterID          = 17
	DataOwnerEID           = 17
	DataHurtDirection      = 18
	DataRabbitType         = 18
	DataCatType            = 18
	DataHurtDamage         = 19
	DataBlockInfo          = 20
	DataWoodID             = 20
	DataMinecartBlock      = 20
	DataInLove             = 21
	DataMinecartOffset     = 21
	DataMinecartHasDisplay = 22
	DataLeadHolder         = 23
	DataLead               = 24
)

const (
	DataFlagOnFire    = 0
	DataFlagSneaking  = 1
	DataFlagRiding    = 2
	DataFlagSprinting = 3
	DataFlagAction    = 4
	DataFlagInvisible = 5
	DataFlagTempted   = 6
)

const (
	DataAnimalFlagIsBaby     = 0
	DataAnimalFlagSitting    = 1
	DataAnimalFlagAngry      = 2
	DataAnimalFlagInterested = 3
)

const (
	DataPlayerFlagSleep = 1
	DataPlayerFlagDead  = 2
)

type MetadataProperty struct {
	Type  int
	Value interface{}
}

type MetadataStore struct {
	Properties map[int]*MetadataProperty
}

func NewMetadataStore() *MetadataStore {
	m := &MetadataStore{
		Properties: make(map[int]*MetadataProperty),
	}

	m.SetByte(DataFlags, 0)
	m.SetShort(DataAir, 300)
	m.SetString(DataNameTag, "")
	m.SetByte(DataShowNameTag, 1)
	m.SetByte(DataSilent, 0)
	m.SetByte(DataNoAI, 0)
	m.SetLong(DataLeadHolder, -1)
	m.SetByte(DataLead, 0)
	return m
}

func (m *MetadataStore) SetByte(key int, value byte) {
	m.Properties[key] = &MetadataProperty{Type: DataTypeByte, Value: value}
}

func (m *MetadataStore) SetShort(key int, value int16) {
	m.Properties[key] = &MetadataProperty{Type: DataTypeShort, Value: value}
}

func (m *MetadataStore) SetInt(key int, value int32) {
	m.Properties[key] = &MetadataProperty{Type: DataTypeInt, Value: value}
}

func (m *MetadataStore) SetFloat(key int, value float32) {
	m.Properties[key] = &MetadataProperty{Type: DataTypeFloat, Value: value}
}

func (m *MetadataStore) SetString(key int, value string) {
	m.Properties[key] = &MetadataProperty{Type: DataTypeString, Value: value}
}

func (m *MetadataStore) SetLong(key int, value int64) {
	m.Properties[key] = &MetadataProperty{Type: DataTypeLong, Value: value}
}

func (m *MetadataStore) GetByte(key int) byte {
	if prop, ok := m.Properties[key]; ok {
		if v, ok := prop.Value.(byte); ok {
			return v
		}
	}
	return 0
}

func (m *MetadataStore) GetShort(key int) int16 {
	if prop, ok := m.Properties[key]; ok {
		if v, ok := prop.Value.(int16); ok {
			return v
		}
	}
	return 0
}

func (m *MetadataStore) GetInt(key int) int32 {
	if prop, ok := m.Properties[key]; ok {
		if v, ok := prop.Value.(int32); ok {
			return v
		}
	}
	return 0
}

func (m *MetadataStore) GetFloat(key int) float32 {
	if prop, ok := m.Properties[key]; ok {
		if v, ok := prop.Value.(float32); ok {
			return v
		}
	}
	return 0
}

func (m *MetadataStore) GetString(key int) string {
	if prop, ok := m.Properties[key]; ok {
		if v, ok := prop.Value.(string); ok {
			return v
		}
	}
	return ""
}

func (m *MetadataStore) GetLong(key int) int64 {
	if prop, ok := m.Properties[key]; ok {
		if v, ok := prop.Value.(int64); ok {
			return v
		}
	}
	return 0
}

func (m *MetadataStore) SetFlag(propertyKey int, flagBit int, value bool) {
	current := m.GetByte(propertyKey)
	if value {
		current |= 1 << flagBit
	} else {
		current &^= 1 << flagBit
	}
	m.SetByte(propertyKey, current)
}

func (m *MetadataStore) GetFlag(propertyKey int, flagBit int) bool {
	current := m.GetByte(propertyKey)
	return (current & (1 << flagBit)) != 0
}

func (m *MetadataStore) Clone() *MetadataStore {
	clone := &MetadataStore{
		Properties: make(map[int]*MetadataProperty),
	}
	for k, v := range m.Properties {
		clone.Properties[k] = &MetadataProperty{
			Type:  v.Type,
			Value: v.Value,
		}
	}
	return clone
}

func (m *MetadataStore) Encode() []byte {
	buf := new(bytes.Buffer)

	for index, prop := range m.Properties {

		header := byte((prop.Type << 5) | (index & 0x1F))
		buf.WriteByte(header)

		m.writeValue(buf, prop.Type, prop.Value)
	}

	buf.WriteByte(0x7f)
	return buf.Bytes()
}

func (m *MetadataStore) writeValue(buf *bytes.Buffer, typeID int, value interface{}) {
	switch typeID {
	case DataTypeByte:
		if v, ok := value.(byte); ok {
			buf.WriteByte(v)
		} else if v, ok := value.(int); ok {
			buf.WriteByte(byte(v))
		} else {
			buf.WriteByte(0)
		}
	case DataTypeShort:
		var v int16
		if val, ok := value.(int16); ok {
			v = val
		} else if val, ok := value.(int); ok {
			v = int16(val)
		}
		binary.Write(buf, binary.LittleEndian, v)
	case DataTypeInt:
		var v int32
		if val, ok := value.(int32); ok {
			v = val
		} else if val, ok := value.(int); ok {
			v = int32(val)
		}
		binary.Write(buf, binary.LittleEndian, v)
	case DataTypeFloat:
		var v float32
		if val, ok := value.(float32); ok {
			v = val
		} else if val, ok := value.(float64); ok {
			v = float32(val)
		}
		binary.Write(buf, binary.LittleEndian, v)
	case DataTypeString:
		v, ok := value.(string)
		if !ok {
			v = ""
		}
		length := int16(len(v))
		binary.Write(buf, binary.LittleEndian, length)
		buf.WriteString(v)
	case DataTypeSlot:

		if parts, ok := value.([]int); ok && len(parts) >= 3 {
			binary.Write(buf, binary.LittleEndian, int16(parts[0]))
			buf.WriteByte(byte(parts[1]))
			binary.Write(buf, binary.LittleEndian, int16(parts[2]))
		} else {

			binary.Write(buf, binary.LittleEndian, int16(0))
			buf.WriteByte(0)
			binary.Write(buf, binary.LittleEndian, int16(0))
		}
	case DataTypePos:
		if parts, ok := value.([]int); ok && len(parts) >= 3 {
			binary.Write(buf, binary.LittleEndian, int32(parts[0]))
			binary.Write(buf, binary.LittleEndian, int32(parts[1]))
			binary.Write(buf, binary.LittleEndian, int32(parts[2]))
		} else {
			binary.Write(buf, binary.LittleEndian, int32(0))
			binary.Write(buf, binary.LittleEndian, int32(0))
			binary.Write(buf, binary.LittleEndian, int32(0))
		}
	case DataTypeLong:
		var v int64
		if val, ok := value.(int64); ok {
			v = val
		} else if val, ok := value.(int); ok {
			v = int64(val)
		}
		binary.Write(buf, binary.LittleEndian, v)
	}
}
