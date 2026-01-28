package gorigional

import (
	"testing"

	"github.com/scaxe/scaxe-go/pkg/world"
)

func TestChunkGeneratorOverworld_GenerateChunk(t *testing.T) {
	gen := NewChunkGeneratorOverworld(12345)
	mockLevel := &MockChunkManager{
		chunks: make(map[int64]*world.Chunk),
	}
	gen.Init(mockLevel, 12345)

	gen.GenerateChunk(0, 0)

	chunk := mockLevel.GetChunk(0, 0, false)

	if chunk == nil {
		t.Fatal("Chunk is nil")
	}

	stoneCount := 0
	waterCount := 0

	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			for y := 0; y < 128; y++ {
				id, _ := chunk.GetBlock(x, y, z)
				if id == 1 {
					stoneCount++
				} else if id == 9 {
					waterCount++
				}
			}
		}
	}

	t.Logf("Generated - Stone: %d, Water: %d", stoneCount, waterCount)

	if stoneCount == 0 {
		t.Error("No stone generated! Terrain logic failed.")
	}
}

type MockChunkManager struct {
	chunks map[int64]*world.Chunk
}

func (m *MockChunkManager) GetChunk(x, z int32, create bool) *world.Chunk {
	idx := (int64(x) << 32) | int64(z)
	if c, ok := m.chunks[idx]; ok {
		return c
	}
	if create {
		c := world.NewChunk(x, z)
		m.chunks[idx] = c
		return c
	}
	return nil
}
func (m *MockChunkManager) SetChunk(x, z int32, c *world.Chunk) {
	idx := (int64(x) << 32) | int64(z)
	m.chunks[idx] = c
}
func (m *MockChunkManager) GetSpawnLocation() *world.Vector3 { return world.NewVector3(0, 64, 0) }
func (m *MockChunkManager) GetSeed() int64                   { return 12345 }

func (m *MockChunkManager) GetBlockId(x, y, z int32) byte {
	cx := x >> 4
	cz := z >> 4
	c := m.GetChunk(cx, cz, false)
	if c != nil {
		id, _ := c.GetBlock(int(x&15), int(y), int(z&15))
		return id
	}
	return 0
}
func (m *MockChunkManager) SetBlock(x, y, z int32, id, meta byte, update bool) bool {
	return true
}
func (m *MockChunkManager) GetHeight(x, z int32) int32 {
	cx := x >> 4
	cz := z >> 4
	c := m.GetChunk(cx, cz, false)
	if c == nil {
		return 0
	}
	localX := int(x & 15)
	localZ := int(z & 15)

	for y := 127; y >= 0; y-- {
		id, _ := c.GetBlock(localX, y, localZ)
		if id != 0 {
			return int32(y + 1)
		}
	}
	return 0
}
