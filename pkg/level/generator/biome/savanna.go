package biome

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type SavannaBiome struct {
	*BaseBiome
	IsPlateau bool
}

func NewSavannaBiome() *SavannaBiome {
	b := &SavannaBiome{
		BaseBiome: &BaseBiome{
			ID:              SAVANNA,
			Name:            "Savanna",
			BaseHeight:      0.125,
			HeightVariation: 0.05,
			Temperature:     1.2,
			Rainfall:        0.0,
			GroundCover:     []block.BlockState{block.NewBlockState(block.GRASS, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
		IsPlateau: false,
	}
	b.Decorator.TreesPerChunk = 1
	b.Decorator.FlowersPerChunk = 4
	b.Decorator.GrassPerChunk = 20
	return b
}

func NewSavannaPlateauBiome() *SavannaBiome {
	b := &SavannaBiome{
		BaseBiome: &BaseBiome{
			ID:              SAVANNA_PLATEAU,
			Name:            "Savanna Plateau",
			BaseHeight:      1.5,
			HeightVariation: 0.025,
			Temperature:     1.0,
			Rainfall:        0.0,
			GroundCover:     []block.BlockState{block.NewBlockState(block.GRASS, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
		IsPlateau: true,
	}
	b.Decorator.TreesPerChunk = 1
	b.Decorator.FlowersPerChunk = 4
	b.Decorator.GrassPerChunk = 20
	return b
}

func (b *SavannaBiome) GetTreeFeature(r *rand.Random) Generator {

	if r.NextBoundedInt(5) > 0 {
		return object.NewSavannaTree()
	}
	return object.NewOakTree()
}

func (b *SavannaBiome) Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos) {

	for i := 0; i < 7; i++ {
		px := r.NextBoundedInt(16) + 8
		pz := r.NextBoundedInt(16) + 8

		terrainHeight := getHeightAtForBiome(level, pos.X()+int32(px), pos.Z()+int32(pz))
		py := r.NextBoundedInt(terrainHeight + 32)
		grass := object.NewDoublePlant(object.DoublePlantGrass)
		grass.Generate(level, r, pos.Add(int32(px), int32(py), int32(pz)))
	}

	b.BaseBiome.Decorate(level, r, pos)
}
