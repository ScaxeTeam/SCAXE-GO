package world

import (
	"bytes"
	"encoding/binary"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

type ChunkPos struct {
	X int32
	Z int32
}

const (
	SectionCount = 16
)

func ChunkHash(x, z int32) int64 {
	return (int64(x) << 32) | (int64(z) & 0xFFFFFFFF)
}

type Chunk struct {
	X int32
	Z int32

	Sections [SectionCount]*ChunkSection

	BiomeColors [256]uint32

	BiomeIds [256]uint8

	HeightMap [256]byte

	Entities []*nbt.CompoundTag
	Tiles    []*nbt.CompoundTag

	Generated      bool
	Populated      bool
	LightPopulated bool

	dirty        bool
	cachedPacket []byte
}

func NewChunk(x, z int32) *Chunk {
	c := &Chunk{
		X:              x,
		Z:              z,
		Generated:      true,
		Populated:      true,
		LightPopulated: true,
		dirty:          true,
	}

	for i := range c.BiomeColors {
		c.BiomeColors[i] = 0xFF7CA800
	}

	return c
}

func (c *Chunk) HasChanged() bool {
	return c.dirty
}

func (c *Chunk) SetChanged(dirty bool) {
	c.dirty = dirty
}

func (c *Chunk) getSection(yIndex int) *ChunkSection {
	if yIndex < 0 || yIndex >= SectionCount {
		return nil
	}
	if c.Sections[yIndex] == nil {
		c.Sections[yIndex] = NewChunkSection(byte(yIndex))
	}
	return c.Sections[yIndex]
}

func (c *Chunk) SetBlock(x, y, z int, id byte, meta byte) bool {
	if x < 0 || x >= 16 || z < 0 || z >= 16 || y < 0 || y >= 256 {
		return false
	}

	sectionY := y >> 4
	sec := c.getSection(sectionY)

	sy := y & 0x0f
	sec.SetBlockId(x, sy, z, id)
	sec.SetBlockData(x, sy, z, meta)

	idx := (z << 4) | x
	currentH := int(c.HeightMap[idx])

	if id != 0 {

		if y >= currentH {
			c.HeightMap[idx] = byte(y + 1)
		}
	} else {

		if y == currentH-1 {

			c.RecalculateColumn(x, z)
		}
	}

	c.dirty = true
	return true
}

func (c *Chunk) GetBlock(x, y, z int) (byte, byte) {
	if x < 0 || x >= 16 || z < 0 || z >= 16 || y < 0 || y >= 256 {
		return 0, 0
	}

	sectionY := y >> 4
	sec := c.Sections[sectionY]
	if sec == nil {
		return 0, 0
	}

	sy := y & 0x0f
	return sec.GetBlockId(x, sy, z), sec.GetBlockData(x, sy, z)
}

func (c *Chunk) GetBlockId(x, y, z int) byte {
	id, _ := c.GetBlock(x, y, z)
	return id
}

func (c *Chunk) GetBlockData(x, y, z int) byte {
	_, meta := c.GetBlock(x, y, z)
	return meta
}

func (c *Chunk) GetBlockLight(x, y, z int) byte {
	if x < 0 || x >= 16 || z < 0 || z >= 16 || y < 0 || y >= 256 {
		return 0
	}
	sectionY := y >> 4
	sec := c.Sections[sectionY]
	if sec == nil {
		return 0
	}
	return sec.GetBlockLight(x, y&0x0f, z)
}

func (c *Chunk) SetBlockLight(x, y, z int, level byte) {
	if x < 0 || x >= 16 || z < 0 || z >= 16 || y < 0 || y >= 256 {
		return
	}
	sectionY := y >> 4
	sec := c.getSection(sectionY)
	sec.SetBlockLight(x, y&0x0f, z, level)
	c.dirty = true
}

func (c *Chunk) GetSkyLight(x, y, z int) byte {
	if x < 0 || x >= 16 || z < 0 || z >= 16 || y < 0 || y >= 256 {
		return 15

	}
	sectionY := y >> 4
	sec := c.Sections[sectionY]
	if sec == nil {
		return 15
	}
	return sec.GetSkyLight(x, y&0x0f, z)
}

func (c *Chunk) SetSkyLight(x, y, z int, level byte) {
	if x < 0 || x >= 16 || z < 0 || z >= 16 || y < 0 || y >= 256 {
		return
	}
	sectionY := y >> 4
	sec := c.getSection(sectionY)
	sec.SetSkyLight(x, y&0x0f, z, level)
	c.dirty = true
}

func (c *Chunk) SetBiomeColor(x, z int, color uint32) {
	if x < 0 || x >= 16 || z < 0 || z >= 16 {
		return
	}
	c.BiomeColors[(z<<4)|x] = color
	c.dirty = true
}

func (c *Chunk) GetBiomeColor(x, z int) uint32 {
	if x < 0 || x >= 16 || z < 0 || z >= 16 {
		return 0
	}
	return c.BiomeColors[(z<<4)|x]
}

func (c *Chunk) SetBiomeID(x, z int, id uint8) {
	if x < 0 || x >= 16 || z < 0 || z >= 16 {
		return
	}
	c.BiomeIds[(z<<4)|x] = id
	c.dirty = true
}

func (c *Chunk) GetBiomeID(x, z int) uint8 {
	if x < 0 || x >= 16 || z < 0 || z >= 16 {
		return 0
	}
	return c.BiomeIds[(z<<4)|x]
}

func (c *Chunk) GetPacketBytes() []byte {
	if c.cachedPacket != nil && !c.dirty {
		return c.cachedPacket
	}

	c.cachedPacket = c.ToPacketBytes()
	c.dirty = false
	return c.cachedPacket
}

func (c *Chunk) ToPacketBytes() []byte {
	buf := new(bytes.Buffer)

	emptyIDs := make([]byte, 4096)
	emptyData := make([]byte, 2048)

	emptySkyLight := make([]byte, 2048)
	for i := range emptySkyLight {
		emptySkyLight[i] = 0xFF
	}
	emptyBlockLight := make([]byte, 2048)

	const SerializeSectionCount = 8

	for i := 0; i < SerializeSectionCount; i++ {
		if s := c.Sections[i]; s != nil {
			buf.Write(s.Blocks)
		} else {
			buf.Write(emptyIDs)
		}
	}

	for i := 0; i < SerializeSectionCount; i++ {
		if s := c.Sections[i]; s != nil {
			buf.Write(s.Data)
		} else {
			buf.Write(emptyData)
		}
	}

	for i := 0; i < SerializeSectionCount; i++ {
		if s := c.Sections[i]; s != nil {
			buf.Write(s.SkyLight)
		} else {
			buf.Write(emptySkyLight)
		}
	}

	for i := 0; i < SerializeSectionCount; i++ {
		if s := c.Sections[i]; s != nil {
			buf.Write(s.BlockLight)
		} else {
			buf.Write(emptyBlockLight)
		}
	}

	buf.Write(c.HeightMap[:])

	for _, color := range c.BiomeColors {
		binary.Write(buf, binary.BigEndian, color)
	}

	binary.Write(buf, binary.LittleEndian, int32(0))

	return buf.Bytes()
}

func (c *Chunk) InitBasicLighting() {

	c.RecalculateHeightMap()

	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			h := int(c.HeightMap[(z<<4)|x])

			for y := 127; y > h; y-- {
				s := c.Sections[y>>4]
				if s == nil {

					continue
				}
				s.SetSkyLight(x, y&0xF, z, 15)
			}

			for y := h; y >= 0; y-- {
				id, _ := c.GetBlock(x, y, z)

				if isSolidForLight(id) {
					break
				}

				s := c.Sections[y>>4]
				if s == nil {
					s = c.getSection(y >> 4)
				}
				s.SetSkyLight(x, y&0xF, z, 15)
			}
		}
	}
	c.LightPopulated = true
}

func isSolidForLight(id byte) bool {

	switch id {
	case 0:
		return false
	case 6:
		return false
	case 18, 161:
		return false
	case 20:
		return false
	case 31:
		return false
	case 37, 38:
		return false
	case 39, 40:
		return false
	case 50:
		return false
	case 51:
		return false
	case 59:
		return false
	case 78:
		return false
	case 83:
		return false
	case 102:
		return false
	case 106:
		return false
	case 111:
		return false
	case 175:
		return false
	default:

		return id != 0
	}
}
