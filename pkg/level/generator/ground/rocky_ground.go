package ground

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/world"
)

type RockyGroundGenerator struct {
	*GroundGenerator
}

func NewRockyGroundGenerator() *RockyGroundGenerator {
	g := &RockyGroundGenerator{
		GroundGenerator: NewGroundGenerator(),
	}
	g.SetTopMaterial(BlockStone, 0)
	g.SetGroundMaterial(BlockStone, 0)
	return g
}

func (g *RockyGroundGenerator) GenerateTerrainColumn(chunk *world.Chunk, random *rand.Rand, x, z int, biome int, surfaceNoise float64) {
	g.GroundGenerator.GenerateTerrainColumn(chunk, random, x, z, biome, surfaceNoise)
}
