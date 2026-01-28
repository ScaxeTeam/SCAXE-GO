package ground

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/world"
)

type GravelPatchGroundGenerator struct {
	*GroundGenerator
}

func NewGravelPatchGroundGenerator() *GravelPatchGroundGenerator {
	return &GravelPatchGroundGenerator{
		GroundGenerator: NewGroundGenerator(),
	}
}

func (g *GravelPatchGroundGenerator) GenerateTerrainColumn(chunk *world.Chunk, random *rand.Rand, x, z int, biome int, surfaceNoise float64) {
	if surfaceNoise > 1.0 {
		g.SetTopMaterial(BlockGravel, 0)
		g.SetGroundMaterial(BlockGravel, 0)
	} else {
		g.SetTopMaterial(BlockGrass, 0)
		g.SetGroundMaterial(BlockDirt, 0)
	}
	g.GroundGenerator.GenerateTerrainColumn(chunk, random, x, z, biome, surfaceNoise)
}

type DirtPatchGroundGenerator struct {
	*GroundGenerator
}

func NewDirtPatchGroundGenerator() *DirtPatchGroundGenerator {
	return &DirtPatchGroundGenerator{
		GroundGenerator: NewGroundGenerator(),
	}
}

func (g *DirtPatchGroundGenerator) GenerateTerrainColumn(chunk *world.Chunk, random *rand.Rand, x, z int, biome int, surfaceNoise float64) {
	if surfaceNoise > 1.75 {
		g.SetTopMaterial(BlockDirt, 1)
		g.SetGroundMaterial(BlockDirt, 0)
	} else if surfaceNoise > -0.95 {
		g.SetTopMaterial(BlockDirt, 2)
		g.SetGroundMaterial(BlockDirt, 0)
	} else {
		g.SetTopMaterial(BlockGrass, 0)
		g.SetGroundMaterial(BlockDirt, 0)
	}
	g.GroundGenerator.GenerateTerrainColumn(chunk, random, x, z, biome, surfaceNoise)
}

type StonePatchGroundGenerator struct {
	*GroundGenerator
}

func NewStonePatchGroundGenerator() *StonePatchGroundGenerator {
	return &StonePatchGroundGenerator{
		GroundGenerator: NewGroundGenerator(),
	}
}

func (g *StonePatchGroundGenerator) GenerateTerrainColumn(chunk *world.Chunk, random *rand.Rand, x, z int, biome int, surfaceNoise float64) {
	if surfaceNoise > 1.0 {
		g.SetTopMaterial(BlockStone, 0)
		g.SetGroundMaterial(BlockStone, 0)
	} else {
		g.SetTopMaterial(BlockGrass, 0)
		g.SetGroundMaterial(BlockDirt, 0)
	}
	g.GroundGenerator.GenerateTerrainColumn(chunk, random, x, z, biome, surfaceNoise)
}

type DirtAndStonePatchGroundGenerator struct {
	*GroundGenerator
}

func NewDirtAndStonePatchGroundGenerator() *DirtAndStonePatchGroundGenerator {
	return &DirtAndStonePatchGroundGenerator{
		GroundGenerator: NewGroundGenerator(),
	}
}

func (g *DirtAndStonePatchGroundGenerator) GenerateTerrainColumn(chunk *world.Chunk, random *rand.Rand, x, z int, biome int, surfaceNoise float64) {
	if surfaceNoise > 1.75 {
		g.SetTopMaterial(BlockStone, 0)
		g.SetGroundMaterial(BlockStone, 0)
	} else if surfaceNoise > -0.5 {
		g.SetTopMaterial(BlockDirt, 1)
		g.SetGroundMaterial(BlockDirt, 0)
	} else {
		g.SetTopMaterial(BlockGrass, 0)
		g.SetGroundMaterial(BlockDirt, 0)
	}
	g.GroundGenerator.GenerateTerrainColumn(chunk, random, x, z, biome, surfaceNoise)
}
