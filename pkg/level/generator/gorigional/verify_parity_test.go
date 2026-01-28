package gorigional_test

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/scaxe/scaxe-go/pkg/level/generator/gorigional"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type ChunkData struct {
	Cx, Cz      int32
	Blocks      []byte
	Meta        []byte
	Biomes      []byte
	BiomeColors []int32
}

type MockChunkManager struct {
	chunks map[int64]*world.Chunk
	seed   int64
}

func NewMockChunkManager(seed int64) *MockChunkManager {
	return &MockChunkManager{
		chunks: make(map[int64]*world.Chunk),
		seed:   seed,
	}
}

func (m *MockChunkManager) GetChunk(x, z int32, generate bool) *world.Chunk {
	hash := world.ChunkHash(x, z)
	if c, ok := m.chunks[hash]; ok {
		return c
	}
	if generate {
		c := world.NewChunk(x, z)
		m.chunks[hash] = c
		return c
	}
	return nil
}

func (m *MockChunkManager) SetChunk(x, z int32, c *world.Chunk) {
	hash := world.ChunkHash(x, z)
	m.chunks[hash] = c
}

func (m *MockChunkManager) GetSeed() int64 {
	return m.seed
}

func (m *MockChunkManager) GetBlockId(x, y, z int32) byte {

	cx := x >> 4
	cz := z >> 4
	hash := world.ChunkHash(cx, cz)
	if c, ok := m.chunks[hash]; ok {
		return c.GetBlockId(int(x&15), int(y), int(z&15))
	}
	return 0
}

func (m *MockChunkManager) SetBlock(x, y, z int32, id, meta byte, update bool) bool {

	cx := x >> 4
	cz := z >> 4
	hash := world.ChunkHash(cx, cz)
	if c, ok := m.chunks[hash]; ok {
		return c.SetBlock(int(x&15), int(y), int(z&15), id, meta)
	}
	return true
}

func (m *MockChunkManager) GetHeight(x, z int32) int32 {

	return 64
}

func TestVerifyParity(t *testing.T) {

	javaFile := "../../../../../mc_server/java_world.bin"
	chunks, seed, err := readJavaWorld(javaFile)
	if err != nil {
		t.Fatalf("Failed to read java world: %v", err)
	}
	t.Logf("Loaded %d chunks from Java world (Seed: %d)", len(chunks), seed)

	mockCM := NewMockChunkManager(seed)
	g := gorigional.NewChunkGeneratorOverworld(seed)
	g.Init(mockCM, seed)

	totalBlocks := 0
	diffBlocks := 0
	diffByTypes := make(map[string]int)

	totalBiomes := 0
	diffBiomes := 0

	startTime := time.Now()

	t.Log("Phase 1: Generating terrain for all chunks...")
	for i, javaChunk := range chunks {
		g.GenerateChunk(javaChunk.Cx, javaChunk.Cz)
		if (i+1)%100 == 0 {
			fmt.Printf("Terrain: %d/%d\n", i+1, len(chunks))
		}
	}

	t.Log("Phase 2: Populating all chunks...")
	for i, javaChunk := range chunks {
		g.PopulateChunk(javaChunk.Cx, javaChunk.Cz)
		if (i+1)%100 == 0 {
			fmt.Printf("Populate: %d/%d\n", i+1, len(chunks))
		}
	}

	t.Log("Phase 3: Comparing blocks...")
	for i, javaChunk := range chunks {
		cx := int(javaChunk.Cx)
		cz := int(javaChunk.Cz)

		goChunk := mockCM.GetChunk(javaChunk.Cx, javaChunk.Cz, false)
		if goChunk == nil {
			t.Errorf("Failed to get Go chunk %d,%d", cx, cz)
			continue
		}

		for y := 0; y < 128; y++ {
			for z := 0; z < 16; z++ {
				for x := 0; x < 16; x++ {

					goId, goMeta := goChunk.GetBlock(x, y, z)

					javaIdx := (x*128+y)*16 + z
					javaId := javaChunk.Blocks[javaIdx]

					metaIdx := javaIdx / 2
					var javaMeta byte
					if (javaIdx & 1) == 0 {
						javaMeta = javaChunk.Meta[metaIdx] & 0x0F
					} else {
						javaMeta = (javaChunk.Meta[metaIdx] >> 4) & 0x0F
					}

					totalBlocks++

					if int(goId) != int(javaId) {
						diffBlocks++
						key := fmt.Sprintf("%d->%d", javaId, goId)
						diffByTypes[key]++
					} else if int(goMeta) != int(javaMeta) {

					}
				}
			}
		}

		for z := 0; z < 16; z++ {
			for x := 0; x < 16; x++ {
				goBiome := goChunk.GetBiomeID(x, z)
				javaBiome := javaChunk.Biomes[z*16+x]

				totalBiomes++
				if int(goBiome) != int(javaBiome) {
					diffBiomes++
				}
			}
		}

		if (i+1)%50 == 0 {
			fmt.Printf("Processed %d/%d chunks...\n", i+1, len(chunks))
		}
	}

	duration := time.Since(startTime)
	t.Logf("Comparison finished in %v", duration)
	t.Logf("Total Blocks: %d", totalBlocks)
	t.Logf("Diff Blocks: %d (%.2f%%)", diffBlocks, float64(diffBlocks)/float64(totalBlocks)*100)
	t.Logf("Diff Biomes: %d (%.2f%%)", diffBiomes, float64(diffBiomes)/float64(totalBiomes)*100)

	t.Logf("Top Block Differences:")
	for k, v := range diffByTypes {
		if v > 1000 {
			t.Logf("%s: %d", k, v)
		}
	}

	biomeDiffByTypes := make(map[string]int)
	for _, javaChunk := range chunks {
		goChunk := mockCM.GetChunk(javaChunk.Cx, javaChunk.Cz, false)
		if goChunk == nil {
			continue
		}
		for z := 0; z < 16; z++ {
			for x := 0; x < 16; x++ {
				goBiome := goChunk.GetBiomeID(x, z)
				javaBiome := javaChunk.Biomes[z*16+x]
				if int(goBiome) != int(javaBiome) {
					key := fmt.Sprintf("B%d->B%d", javaBiome, goBiome)
					biomeDiffByTypes[key]++
				}
			}
		}
	}

	t.Logf("Top Biome Differences:")
	for k, v := range biomeDiffByTypes {
		if v > 100 {
			t.Logf("%s: %d", k, v)
		}
	}

	reportPath := "C:/Users/VHDJD/.gemini/antigravity/brain/9854d431-dac5-4613-8c83-aee191fe28cc/comparison_report_v3.txt"
	reportFile, err := os.Create(reportPath)
	if err != nil {
		t.Logf("Failed to create report file: %v", err)
	} else {
		defer reportFile.Close()
		fmt.Fprintf(reportFile, "Parity Check Report (Hills Fix)\n")
		fmt.Fprintf(reportFile, "Total Blocks: %d\n", totalBlocks)
		fmt.Fprintf(reportFile, "Diff Blocks: %d (%.2f%%)\n", diffBlocks, float64(diffBlocks)/float64(totalBlocks)*100)
		fmt.Fprintf(reportFile, "Diff Biomes: %d (%.2f%%)\n", diffBiomes, float64(diffBiomes)/float64(totalBiomes)*100)
		fmt.Fprintf(reportFile, "\nTop Block Differences:\n")
		for k, v := range diffByTypes {
			if v > 1000 {
				fmt.Fprintf(reportFile, "%s: %d\n", k, v)
			}
		}
		fmt.Fprintf(reportFile, "\nTop Biome Differences:\n")
		for k, v := range biomeDiffByTypes {
			if v > 100 {
				fmt.Fprintf(reportFile, "%s: %d\n", k, v)
			}
		}
		t.Logf("Report written to %s", reportPath)
	}
}

func readJavaWorld(filename string) ([]ChunkData, int64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	var seed int64
	err = binary.Read(r, binary.BigEndian, &seed)
	if err != nil {
		return nil, 0, err
	}

	var count int32
	err = binary.Read(r, binary.BigEndian, &count)
	if err != nil {
		return nil, 0, err
	}

	chunks := make([]ChunkData, count)

	for i := 0; i < int(count); i++ {
		var cx, cz int32
		binary.Read(r, binary.BigEndian, &cx)
		binary.Read(r, binary.BigEndian, &cz)

		blocks := make([]byte, 16*128*16)
		io.ReadFull(r, blocks)

		meta := make([]byte, 16*128*16/2)
		io.ReadFull(r, meta)

		biomes := make([]byte, 256)
		io.ReadFull(r, biomes)

		colors := make([]int32, 256)
		binary.Read(r, binary.BigEndian, &colors)

		chunks[i] = ChunkData{
			Cx: cx, Cz: cz,
			Blocks: blocks, Meta: meta,
			Biomes: biomes, BiomeColors: colors,
		}
	}

	return chunks, seed, nil
}
