package ground

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/world"
)

type SnowyGroundGenerator struct {
	*GroundGenerator
}

func NewSnowyGroundGenerator() *SnowyGroundGenerator {
	g := &SnowyGroundGenerator{
		GroundGenerator: NewGroundGenerator(),
	}
	g.SetTopMaterial(BlockSnow, 0)
	g.SetGroundMaterial(BlockDirt, 0)
	return g
}

func (g *SnowyGroundGenerator) GenerateTerrainColumn(chunk *world.Chunk, random *rand.Rand, x, z int, biome int, surfaceNoise float64) {
	g.GroundGenerator.GenerateTerrainColumn(chunk, random, x, z, biome, surfaceNoise)
}
