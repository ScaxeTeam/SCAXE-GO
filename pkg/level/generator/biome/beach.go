package biome

import (
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type BeachBiome struct {
	*BaseBiome
}

func NewBeachBiome() *BeachBiome {
	b := &BeachBiome{
		BaseBiome: &BaseBiome{
			ID:              16,
			Name:            "Beach",
			BaseHeight:      0.0,
			HeightVariation: 0.025,
			Temperature:     0.8,
			Rainfall:        0.4,
			Decorator:       NewDecorator(),
		},
	}

	b.Decorator.TreesPerChunk = -999
	b.Decorator.DeadBushPerChunk = 0
	b.Decorator.ReedsPerChunk = 0
	b.Decorator.CactiPerChunk = 0
	return b
}

func (b *BeachBiome) Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos) {

	b.BaseBiome.Decorate(level, r, pos)
}

type StoneBeachBiome struct {
	*BaseBiome
}

func NewStoneBeachBiome() *StoneBeachBiome {
	b := &StoneBeachBiome{
		BaseBiome: &BaseBiome{
			ID:              25,
			Name:            "Stone Beach",
			BaseHeight:      0.1,
			HeightVariation: 0.8,
			Temperature:     0.2,
			Rainfall:        0.3,
			Decorator:       NewDecorator(),
		},
	}

	b.Decorator.TreesPerChunk = -999
	return b
}

func (b *StoneBeachBiome) Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos) {

	b.BaseBiome.Decorate(level, r, pos)
}

type ColdBeachBiome struct {
	*BaseBiome
}

func NewColdBeachBiome() *ColdBeachBiome {
	b := &ColdBeachBiome{
		BaseBiome: &BaseBiome{
			ID:              26,
			Name:            "Cold Beach",
			BaseHeight:      0.0,
			HeightVariation: 0.025,
			Temperature:     0.05,
			Rainfall:        0.3,
			Decorator:       NewDecorator(),
		},
	}
	b.Decorator.TreesPerChunk = -999
	return b
}

func (b *ColdBeachBiome) Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos) {
	b.BaseBiome.Decorate(level, r, pos)
}
