package ground

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	BlockAir       = 0
	BlockStone     = 1
	BlockGrass     = 2
	BlockDirt      = 3
	BlockSand      = 12
	BlockGravel    = 13
	BlockSandstone = 24
	BlockSnow      = 80
	BlockIce       = 79
	BlockMycelium  = 110
	BlockWater     = 9
)

type Generator interface {
	GenerateTerrainColumn(chunk *world.Chunk, random *rand.Rand, x, z int, biome int, surfaceNoise float64)
}

type GroundGenerator struct {
	TopMaterial    int
	TopMeta        byte
	GroundMaterial int
	GroundMeta     byte
}

func NewGroundGenerator() *GroundGenerator {
	return &GroundGenerator{
		TopMaterial:    BlockGrass,
		TopMeta:        0,
		GroundMaterial: BlockDirt,
		GroundMeta:     0,
	}
}

func (g *GroundGenerator) GenerateTerrainColumn(chunk *world.Chunk, random *rand.Rand, x, z int, biome int, surfaceNoise float64) {
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
			} else {
				chunk.SetBlock(x, y, z, byte(g.TopMaterial), g.TopMeta)
			}
			depth = surfaceDepth
		} else if depth > 0 {

			depth--
			chunk.SetBlock(x, y, z, byte(g.GroundMaterial), g.GroundMeta)
		}
	}
}

func (g *GroundGenerator) SetTopMaterial(id int, meta byte) {
	g.TopMaterial = id
	g.TopMeta = meta
}

func (g *GroundGenerator) SetGroundMaterial(id int, meta byte) {
	g.GroundMaterial = id
	g.GroundMeta = meta
}
