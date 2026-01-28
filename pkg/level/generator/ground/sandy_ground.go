package ground

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/world"
)

type SandyGroundGenerator struct {
	*GroundGenerator
}

func NewSandyGroundGenerator() *SandyGroundGenerator {
	g := &SandyGroundGenerator{
		GroundGenerator: NewGroundGenerator(),
	}
	g.SetTopMaterial(BlockSand, 0)
	g.SetGroundMaterial(BlockSandstone, 0)
	return g
}

func (g *SandyGroundGenerator) GenerateTerrainColumn(chunk *world.Chunk, random *rand.Rand, x, z int, biome int, surfaceNoise float64) {

	g.GroundGenerator.GenerateTerrainColumn(chunk, random, x, z, biome, surfaceNoise)
}
