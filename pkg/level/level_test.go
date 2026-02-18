package level

import (
	"testing"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator"
	"github.com/scaxe/scaxe-go/pkg/world"
)

func init() {
	block.Registry.Init()
}

func makeTestLevel(spawnX, spawnY, spawnZ float64) *Level {
	l := &Level{
		ID:     1,
		Name:   "test",
		Chunks: make(map[int64]*world.Chunk),
	}
	l.Generator = &mockGen{spawn: world.NewVector3(spawnX, spawnY, spawnZ)}
	return l
}

func TestGetSafeSpawn_NormalTerrain(t *testing.T) {
	l := makeTestLevel(8, 64, 8)

	chunk := world.NewChunk(0, 0)
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			chunk.SetBlock(x, 0, z, block.BEDROCK, 0)
			for y := 1; y <= 62; y++ {
				chunk.SetBlock(x, y, z, block.STONE, 0)
			}
			chunk.SetBlock(x, 63, z, block.GRASS, 0)
		}
	}
	l.Chunks[world.ChunkHash(0, 0)] = chunk

	safe := l.GetSafeSpawn()

	if safe.Y < 64 {
		t.Errorf("Expected Y >= 64 (above grass at 63), got Y=%v", safe.Y)
	}
	if safe.Y > 127 {
		t.Errorf("Expected Y <= 127, got Y=%v", safe.Y)
	}
	t.Logf("Safe spawn Y=%v (terrain top=63)", safe.Y)
}

func TestGetSafeSpawn_SpawnInsideTerrain(t *testing.T) {
	l := makeTestLevel(8, 40, 8)

	chunk := world.NewChunk(0, 0)
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			for y := 0; y <= 80; y++ {
				chunk.SetBlock(x, y, z, block.STONE, 0)
			}
		}
	}
	l.Chunks[world.ChunkHash(0, 0)] = chunk

	safe := l.GetSafeSpawn()

	if safe.Y <= 80 {
		t.Errorf("Expected Y > 80 (above stone terrain), got Y=%v", safe.Y)
	}
	t.Logf("Safe spawn Y=%v (terrain top=80)", safe.Y)
}

func TestGetSafeSpawn_FlatBedrock(t *testing.T) {
	l := makeTestLevel(8, 64, 8)

	chunk := world.NewChunk(0, 0)
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			chunk.SetBlock(x, 0, z, block.BEDROCK, 0)
		}
	}
	l.Chunks[world.ChunkHash(0, 0)] = chunk

	safe := l.GetSafeSpawn()

	if safe.Y < 1 {
		t.Errorf("Expected Y >= 1 (above bedrock at 0), got Y=%v", safe.Y)
	}
	t.Logf("Safe spawn Y=%v (only bedrock at y=0)", safe.Y)
}

func TestGetSafeSpawn_EmptyChunk(t *testing.T) {
	l := makeTestLevel(8, 64, 8)

	chunk := world.NewChunk(0, 0)
	l.Chunks[world.ChunkHash(0, 0)] = chunk

	safe := l.GetSafeSpawn()

	if safe.Y < 0 || safe.Y > 128 {
		t.Errorf("Expected Y in [0, 128], got Y=%v", safe.Y)
	}
	t.Logf("Safe spawn Y=%v (empty chunk)", safe.Y)
}

func TestGetSafeSpawn_HighTerrain(t *testing.T) {
	l := makeTestLevel(8, 64, 8)

	chunk := world.NewChunk(0, 0)
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			for y := 0; y <= 120; y++ {
				chunk.SetBlock(x, y, z, block.STONE, 0)
			}
		}
	}
	l.Chunks[world.ChunkHash(0, 0)] = chunk

	safe := l.GetSafeSpawn()

	if safe.Y <= 120 {
		t.Errorf("Expected Y > 120 (above terrain), got Y=%v", safe.Y)
	}
	t.Logf("Safe spawn Y=%v (terrain to y=120)", safe.Y)
}

type mockGen struct {
	spawn *world.Vector3
}

var _ generator.Generator = (*mockGen)(nil)

func (g *mockGen) GetName() string                               { return "mock" }
func (g *mockGen) Init(level generator.ChunkManager, seed int64) {}
func (g *mockGen) GenerateChunk(cx, cz int32)                    {}
func (g *mockGen) PopulateChunk(cx, cz int32)                    {}
func (g *mockGen) GetSpawn() *world.Vector3                      { return g.spawn }
func (g *mockGen) GetSettings() map[string]interface{}           { return nil }
