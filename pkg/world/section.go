package world

type ChunkSection struct {
	Y          byte
	Blocks     []byte
	Data       []byte
	BlockLight []byte
	SkyLight   []byte
}

func NewChunkSection(y byte) *ChunkSection {
	s := &ChunkSection{
		Y:          y,
		Blocks:     make([]byte, 4096),
		Data:       make([]byte, 2048),
		BlockLight: make([]byte, 2048),
		SkyLight:   make([]byte, 2048),
	}

	for i := range s.SkyLight {
		s.SkyLight[i] = 0xFF
	}
	return s
}

func (s *ChunkSection) getBlockIndex(x, y, z int) int {
	return (y << 8) | (z << 4) | x
}

func (s *ChunkSection) getNibbleIndex(x, y, z int) (int, int) {
	index := (y << 7) | (z << 3) | (x >> 1)
	shift := (x & 1) << 2
	return index, shift
}

func (s *ChunkSection) GetBlockId(x, y, z int) byte {
	return s.Blocks[s.getBlockIndex(x, y, z)]
}

func (s *ChunkSection) SetBlockId(x, y, z int, id byte) {
	s.Blocks[s.getBlockIndex(x, y, z)] = id
}

func (s *ChunkSection) GetBlockData(x, y, z int) byte {
	index, shift := s.getNibbleIndex(x, y, z)
	return (s.Data[index] >> shift) & 0x0f
}

func (s *ChunkSection) SetBlockData(x, y, z int, data byte) {
	index, shift := s.getNibbleIndex(x, y, z)
	s.Data[index] &= ^(0x0f << shift)
	s.Data[index] |= (data & 0x0f) << shift
}

func (s *ChunkSection) GetBlockLight(x, y, z int) byte {
	index, shift := s.getNibbleIndex(x, y, z)
	return (s.BlockLight[index] >> shift) & 0x0f
}

func (s *ChunkSection) SetBlockLight(x, y, z int, level byte) {
	index, shift := s.getNibbleIndex(x, y, z)
	s.BlockLight[index] &= ^(0x0f << shift)
	s.BlockLight[index] |= (level & 0x0f) << shift
}

func (s *ChunkSection) GetSkyLight(x, y, z int) byte {
	index, shift := s.getNibbleIndex(x, y, z)
	return (s.SkyLight[index] >> shift) & 0x0f
}

func (s *ChunkSection) SetSkyLight(x, y, z int, level byte) {
	index, shift := s.getNibbleIndex(x, y, z)
	s.SkyLight[index] &= ^(0x0f << shift)
	s.SkyLight[index] |= (level & 0x0f) << shift
}

func (s *ChunkSection) IsEmpty() bool {

	for _, b := range s.Blocks {
		if b != 0 {
			return false
		}
	}
	return true
}

func (s *ChunkSection) FillSkyLight(value byte) {
	packed := (value & 0x0f) | ((value & 0x0f) << 4)
	for i := range s.SkyLight {
		s.SkyLight[i] = packed
	}
}

func (s *ChunkSection) Copy() *ChunkSection {
	c := NewChunkSection(s.Y)
	copy(c.Blocks, s.Blocks)
	copy(c.Data, s.Data)
	copy(c.BlockLight, s.BlockLight)
	copy(c.SkyLight, s.SkyLight)
	return c
}
