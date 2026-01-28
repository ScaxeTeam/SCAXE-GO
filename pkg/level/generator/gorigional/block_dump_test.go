package gorigional

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/scaxe/scaxe-go/pkg/level/generator/biome"
	"github.com/scaxe/scaxe-go/pkg/world"
)

func TestDumpBlocksForParity(t *testing.T) {

	biome.InitBiomes()

	const (
		seed       = int64(114514)
		chunkRange = 18
		maxHeight  = 128
		outputFile = "../../../../../go_world.bin"
	)

	t.Logf("=== Go World Generation Dump ===")
	t.Logf("Seed: %d", seed)
	t.Logf("Range: -%d to +%d (%d chunks)", chunkRange, chunkRange-1, chunkRange*2*chunkRange*2)

	gen := NewChunkGeneratorOverworld(seed)
	mockLevel := &ParityMockChunkManager{
		chunks: make(map[int64]*world.Chunk),
		seed:   seed,
	}
	gen.Init(mockLevel, seed)

	t.Log("Phase 1: Generating terrain...")
	count := 0
	for cx := int32(-chunkRange); cx < chunkRange; cx++ {
		for cz := int32(-chunkRange); cz < chunkRange; cz++ {
			gen.GenerateChunk(cx, cz)
			count++
			if count%100 == 0 {
				t.Logf("Terrain: %d/%d", count, chunkRange*2*chunkRange*2)
			}
		}
	}
	t.Logf("Terrain: %d chunks generated", count)

	t.Log("Phase 2: Populating...")
	count = 0
	for cx := int32(-chunkRange); cx < chunkRange; cx++ {
		for cz := int32(-chunkRange); cz < chunkRange; cz++ {
			gen.PopulateChunk(cx, cz)
			count++
			if count%100 == 0 {
				t.Logf("Populate: %d/%d", count, chunkRange*2*chunkRange*2)
			}
		}
	}
	t.Logf("Populate: %d chunks populated", count)

	t.Log("Phase 3: Writing output...")
	err := writeParityDump(outputFile, seed, mockLevel.chunks, maxHeight)
	if err != nil {
		t.Fatalf("Failed to write output: %v", err)
	}

	t.Logf("=== Complete: %d chunks written to %s ===", len(mockLevel.chunks), outputFile)
}

type ParityMockChunkManager struct {
	chunks map[int64]*world.Chunk
	seed   int64
}

func (m *ParityMockChunkManager) GetChunk(x, z int32, create bool) *world.Chunk {
	idx := (int64(x) << 32) | int64(uint32(z))
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

func (m *ParityMockChunkManager) SetChunk(x, z int32, c *world.Chunk) {
	idx := (int64(x) << 32) | int64(uint32(z))
	m.chunks[idx] = c
}

func (m *ParityMockChunkManager) GetSpawnLocation() *world.Vector3 { return world.NewVector3(0, 64, 0) }
func (m *ParityMockChunkManager) GetSeed() int64                   { return m.seed }

func (m *ParityMockChunkManager) GetBlockId(x, y, z int32) byte {
	cx := x >> 4
	cz := z >> 4
	c := m.GetChunk(cx, cz, false)
	if c != nil {
		id, _ := c.GetBlock(int(x&15), int(y), int(z&15))
		return id
	}
	return 0
}

func (m *ParityMockChunkManager) SetBlock(x, y, z int32, id, meta byte, update bool) bool {
	cx := x >> 4
	cz := z >> 4
	c := m.GetChunk(cx, cz, true)
	if c != nil {
		c.SetBlock(int(x&15), int(y), int(z&15), id, meta)
		return true
	}
	return false
}

func (m *ParityMockChunkManager) GetHeight(x, z int32) int32 {
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

func writeParityDump(filename string, seed int64, chunks map[int64]*world.Chunk, maxHeight int) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	binary.Write(f, binary.BigEndian, seed)
	binary.Write(f, binary.BigEndian, int32(len(chunks)))

	for _, chunk := range chunks {

		binary.Write(f, binary.BigEndian, chunk.X)
		binary.Write(f, binary.BigEndian, chunk.Z)

		blockData := make([]byte, 16*maxHeight*16)
		for x := 0; x < 16; x++ {
			for y := 0; y < maxHeight; y++ {
				for z := 0; z < 16; z++ {
					id, _ := chunk.GetBlock(x, y, z)
					blockData[(x*maxHeight+y)*16+z] = id
				}
			}
		}
		f.Write(blockData)

		metaData := make([]byte, 16*maxHeight*16/2)
		for x := 0; x < 16; x++ {
			for y := 0; y < maxHeight; y++ {
				for z := 0; z < 16; z++ {
					_, meta := chunk.GetBlock(x, y, z)
					idx := (x*maxHeight+y)*16 + z
					metaIdx := idx / 2
					if idx&1 == 0 {
						metaData[metaIdx] = (metaData[metaIdx] & 0xF0) | (meta & 0x0F)
					} else {
						metaData[metaIdx] = (metaData[metaIdx] & 0x0F) | ((meta & 0x0F) << 4)
					}
				}
			}
		}
		f.Write(metaData)

		biomeData := make([]byte, 256)
		for i := 0; i < 256; i++ {

			biomeData[i] = chunk.GetBiomeID(i&15, i>>4)
		}
		f.Write(biomeData)

		colorData := make([]byte, 256*4)
		for x := 0; x < 16; x++ {
			for z := 0; z < 16; z++ {
				idx := z*16 + x
				color := chunk.GetBiomeColor(x, z)
				colorData[idx*4+0] = byte((color >> 24) & 0xFF)
				colorData[idx*4+1] = byte((color >> 16) & 0xFF)
				colorData[idx*4+2] = byte((color >> 8) & 0xFF)
				colorData[idx*4+3] = byte(color & 0xFF)
			}
		}
		f.Write(colorData)
	}

	return nil
}
