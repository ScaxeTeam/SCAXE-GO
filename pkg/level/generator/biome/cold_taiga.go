package biome

import (
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type ColdTaigaBiome struct {
	*BaseBiome
}

func NewColdTaigaBiome() *ColdTaigaBiome {
	b := &ColdTaigaBiome{
		BaseBiome: &BaseBiome{
			ID:              30,
			Name:            "Cold Taiga",
			BaseHeight:      0.2,
			HeightVariation: 0.2,
			Temperature:     -0.5,
			Rainfall:        0.4,
			Decorator:       NewDecorator(),
		},
	}

	b.Decorator.TreesPerChunk = 10
	b.Decorator.GrassPerChunk = 1
	b.Decorator.MushroomsPerChunk = 1
	return b
}

func NewColdTaigaHillsBiome() *ColdTaigaBiome {
	b := &ColdTaigaBiome{
		BaseBiome: &BaseBiome{
			ID:              31,
			Name:            "Cold Taiga Hills",
			BaseHeight:      0.45,
			HeightVariation: 0.3,
			Temperature:     -0.5,
			Rainfall:        0.4,
			Decorator:       NewDecorator(),
		},
	}
	b.Decorator.TreesPerChunk = 10
	b.Decorator.GrassPerChunk = 1
	b.Decorator.MushroomsPerChunk = 1
	return b
}

func (b *ColdTaigaBiome) GetTreeFeature(r *rand.Random) Generator {

	if r.NextBoundedInt(3) == 0 {

		return object.NewPineTree()
	}

	return object.NewSpruceTree()
}
