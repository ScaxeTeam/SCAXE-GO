package ground

import (
	"math"
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	BlockHardenedClay = 172
	BlockRedSand      = 12
	BlockRedSandstone = 179
)

var clayBands = []byte{1, 1, 4, 4, 7, 7, 4, 8, 1, 4, 1, 12, 12, 14, 14, 1}

type MesaGroundGenerator struct {
	*GroundGenerator
	hasBryce  bool
	hasForest bool
}

func NewMesaGroundGenerator() *MesaGroundGenerator {
	g := &MesaGroundGenerator{
		GroundGenerator: NewGroundGenerator(),
		hasBryce:        false,
		hasForest:       false,
	}
	g.SetTopMaterial(BlockRedSand, 1)
	g.SetGroundMaterial(BlockHardenedClay, 0)
	return g
}

func NewMesaBryceGroundGenerator() *MesaGroundGenerator {
	g := NewMesaGroundGenerator()
	g.hasBryce = true
	return g
}

func NewMesaForestGroundGenerator() *MesaGroundGenerator {
	g := NewMesaGroundGenerator()
	g.hasForest = true
	g.SetTopMaterial(BlockGrass, 0)
	return g
}

func (g *MesaGroundGenerator) GenerateTerrainColumn(chunk *world.Chunk, random *rand.Rand, x, z int, biome int, surfaceNoise float64) {
	seaLevel := 64
	surfaceDepth := int(surfaceNoise/3.0 + 3.0 + random.Float64()*0.25)

	topBlock := -1
	for y := 127; y >= 0; y-- {
		blockID, _ := chunk.GetBlock(x, y, z)
		if blockID != BlockAir && blockID != BlockWater {
			topBlock = y
			break
		}
	}

	if topBlock < 0 {
		return
	}

	depth := 0
	for y := topBlock; y >= 0; y-- {
		blockID, _ := chunk.GetBlock(x, y, z)

		if blockID == BlockAir {
			depth = 0
			continue
		}

		if blockID != BlockStone {
			continue
		}

		if depth == 0 {
			if y < seaLevel-1 {
				chunk.SetBlock(x, y, z, BlockGravel, 0)
			} else if g.hasForest && y > seaLevel+10 {
				chunk.SetBlock(x, y, z, BlockGrass, 0)
			} else {
				chunk.SetBlock(x, y, z, BlockRedSand, 1)
			}
			depth = surfaceDepth
		} else if depth > 0 {
			depth--

			bandIndex := int(math.Abs(float64(y))) % len(clayBands)
			chunk.SetBlock(x, y, z, BlockHardenedClay, clayBands[bandIndex])
		}
	}
}
