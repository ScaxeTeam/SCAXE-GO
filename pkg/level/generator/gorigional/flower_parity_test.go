package gorigional

import (
	"fmt"
	"testing"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type MockLayerManager struct {
	blocks map[world.BlockPos]byte
	metas  map[world.BlockPos]byte
}

func NewMockLayerManager() *MockLayerManager {
	return &MockLayerManager{
		blocks: make(map[world.BlockPos]byte),
		metas:  make(map[world.BlockPos]byte),
	}
}

func (m *MockLayerManager) GetBlockId(x, y, z int32) byte {
	pos := world.NewBlockPos(x, y, z)
	if id, ok := m.blocks[pos]; ok {
		return id
	}
	return 0
}
func (m *MockLayerManager) SetBlock(x, y, z int32, id, meta byte, update bool) bool {
	pos := world.NewBlockPos(x, y, z)
	m.blocks[pos] = id
	m.metas[pos] = meta
	return true
}
func (m *MockLayerManager) GetHeight(x, z int32) int32                    { return 64 }
func (m *MockLayerManager) GetChunk(x, z int32, create bool) *world.Chunk { return nil }
func (m *MockLayerManager) SetChunk(x, z int32, c *world.Chunk)           {}
func (m *MockLayerManager) GetSeed() int64                                { return 12345 }

func TestFlowerParity(t *testing.T) {
	seed := int64(12345)
	gen := NewChunkGeneratorOverworld(seed)

	manager := NewMockLayerManager()
	gen.Init(manager, seed)

	for x := int32(0); x < 160; x++ {
		for z := int32(0); z < 160; z++ {
			manager.SetBlock(x, 63, z, block.GRASS, 0, false)
			manager.SetBlock(x, 64, z, 0, 0, false)
		}
	}

	for cx := int32(0); cx < 10; cx++ {
		for cz := int32(0); cz < 10; cz++ {
			gen.PopulateChunk(cx, cz)
		}
	}

	flowers := make(map[string]int)
	doublePlants := make(map[string]int)

	for pos, id := range manager.blocks {
		meta := manager.metas[pos]
		if id == block.DANDELION {
			flowers["Dandelion"]++
		}
		if id == block.RED_FLOWER {
			name := "Unknown"
			switch meta {
			case 0:
				name = "Poppy"
			case 1:
				name = "BlueOrchid"
			case 2:
				name = "Allium"
			case 3:
				name = "AzureBluet"
			case 4:
				name = "RedTulip"
			case 5:
				name = "OrangeTulip"
			case 6:
				name = "WhiteTulip"
			case 7:
				name = "PinkTulip"
			case 8:
				name = "OxeyeDaisy"
			}
			flowers[name]++
		}
		if id == block.DOUBLE_PLANT {
			name := "UnknownDouble"
			switch meta & 7 {
			case 0:
				name = "Sunflower"
			case 1:
				name = "Lilac"
			case 2:
				name = "DoubleGrass"
			case 3:
				name = "DoubleFern"
			case 4:
				name = "RoseBush"
			case 5:
				name = "Peony"
			}
			doublePlants[name]++
		}
		if id == block.WOOD || id == block.LOG || id == block.LEAVES {
			flowers["TREE_RELATED"]++
		}
	}

	fmt.Println("=== Flower Distribution Report ===")
	for name, count := range flowers {
		fmt.Printf("%s: %d\n", name, count)
	}
	fmt.Println("=== Double Plant Report ===")
	for name, count := range doublePlants {
		fmt.Printf("%s: %d\n", name, count)
	}

	if len(flowers) <= 1 {
		t.Errorf("Expected variety of flowers, got %v. Noise gradient might not be working.", flowers)
	}
}
