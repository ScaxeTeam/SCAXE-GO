package object_test

import (
	"testing"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type MockChunkManager struct {
	Blocks map[world.BlockPos]uint8
	Metas  map[world.BlockPos]uint8
}

func NewMockChunkManager() *MockChunkManager {
	return &MockChunkManager{
		Blocks: make(map[world.BlockPos]uint8),
		Metas:  make(map[world.BlockPos]uint8),
	}
}

func (m *MockChunkManager) GetChunk(x, z int32, create bool) *world.Chunk { return nil }
func (m *MockChunkManager) SetChunk(x, z int32, chunk *world.Chunk)       {}
func (m *MockChunkManager) GetSeed() int64                                { return 0 }
func (m *MockChunkManager) GetBlockId(x, y, z int32) byte {
	pos := world.NewBlockPos(x, y, z)
	if id, ok := m.Blocks[pos]; ok {
		return id
	}
	return block.AIR
}
func (m *MockChunkManager) SetBlock(x, y, z int32, id, meta byte, update bool) bool {
	pos := world.NewBlockPos(x, y, z)
	m.Blocks[pos] = id
	m.Metas[pos] = meta
	return true
}
func (m *MockChunkManager) GetHeight(x, z int32) int32 {

	for y := int32(255); y >= 0; y-- {
		if m.GetBlockId(x, y, z) != block.AIR {
			return y + 1
		}
	}
	return 0
}

func TestBigMushroomGeneration(t *testing.T) {
	r := rand.NewRandom(12345)

	gen := object.NewBigMushroom(block.RED_MUSHROOM_BLOCK)
	cm := NewMockChunkManager()

	center := world.NewBlockPos(0, 64, 0)
	cm.SetBlock(center.X(), center.Y()-1, center.Z(), block.DIRT, 0, false)

	if !gen.Generate(cm, r, center) {
		t.Error("Failed to generate red mushroom on dirt")
	}

	stemBase := cm.GetBlockId(center.X(), center.Y(), center.Z())
	if stemBase != block.RED_MUSHROOM_BLOCK {
		t.Errorf("Expected stem base, got %d", stemBase)
	}

	genBrown := object.NewBigMushroom(block.BROWN_MUSHROOM_BLOCK)
	cmBrown := NewMockChunkManager()
	centerBrown := world.NewBlockPos(20, 64, 20)
	cmBrown.SetBlock(centerBrown.X(), centerBrown.Y()-1, centerBrown.Z(), block.GRASS, 0, false)

	if !genBrown.Generate(cmBrown, r, centerBrown) {
		t.Error("Failed to generate brown mushroom on grass")
	}
}
