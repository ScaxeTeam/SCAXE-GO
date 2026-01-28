package ground

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/world"
)

type MycelGroundGenerator struct {
	*GroundGenerator
}

func NewMycelGroundGenerator() *MycelGroundGenerator {
	g := &MycelGroundGenerator{
		GroundGenerator: NewGroundGenerator(),
	}
	g.SetTopMaterial(BlockMycelium, 0)
	g.SetGroundMaterial(BlockDirt, 0)
	return g
}

func (g *MycelGroundGenerator) GenerateTerrainColumn(chunk *world.Chunk, random *rand.Rand, x, z int, biome int, surfaceNoise float64) {
	g.GroundGenerator.GenerateTerrainColumn(chunk, random, x, z, biome, surfaceNoise)
}
