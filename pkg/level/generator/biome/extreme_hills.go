package biome

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type ExtremeHillsBiome struct {
	*BaseBiome
	Type int
}

const (
	HILLS_NORMAL      = 0
	HILLS_EXTRA_TREES = 1
	HILLS_MUTATED     = 2
)

func NewExtremeHillsBiome() *ExtremeHillsBiome {
	b := &ExtremeHillsBiome{
		BaseBiome: &BaseBiome{
			ID:              3,
			Name:            "Extreme Hills",
			BaseHeight:      1.0,
			HeightVariation: 0.5,
			Temperature:     0.2,
			Rainfall:        0.3,
			Decorator:       NewDecorator(),
		},
		Type: HILLS_NORMAL,
	}

	b.Decorator.TreesPerChunk = 0
	return b
}

func NewExtremeHillsPlusBiome() *ExtremeHillsBiome {
	b := &ExtremeHillsBiome{
		BaseBiome: &BaseBiome{
			ID:              34,
			Name:            "Extreme Hills+",
			BaseHeight:      1.0,
			HeightVariation: 0.5,
			Temperature:     0.2,
			Rainfall:        0.3,
			Decorator:       NewDecorator(),
		},
		Type: HILLS_EXTRA_TREES,
	}

	b.Decorator.TreesPerChunk = 3
	return b
}

func NewExtremeHillsEdgeBiome() *ExtremeHillsBiome {
	b := &ExtremeHillsBiome{
		BaseBiome: &BaseBiome{
			ID:              20,
			Name:            "Extreme Hills Edge",
			BaseHeight:      0.8,
			HeightVariation: 0.3,
			Temperature:     0.2,
			Rainfall:        0.3,
			Decorator:       NewDecorator(),
		},
		Type: HILLS_NORMAL,
	}
	b.Decorator.TreesPerChunk = 0
	return b
}

func (b *ExtremeHillsBiome) GetTreeFeature(r *rand.Random) Generator {

	if r.NextBoundedInt(3) > 0 {
		return object.NewSpruceTree()
	}
	return object.NewOakTree()
}

func (b *ExtremeHillsBiome) Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos) {

	b.BaseBiome.Decorate(level, r, pos)

	emeraldCount := r.NextBoundedInt(3) + 3 + r.NextBoundedInt(6)
	for i := 0; i < emeraldCount; i++ {
		x := int32(r.NextBoundedInt(16))
		y := int32(r.NextBoundedInt(28) + 4)
		z := int32(r.NextBoundedInt(16))

		targetPos := pos.Add(x, y, z)
		if level.GetBlockId(targetPos.X(), targetPos.Y(), targetPos.Z()) == block.STONE {
			level.SetBlock(targetPos.X(), targetPos.Y(), targetPos.Z(), block.EMERALD_ORE, 0, false)
		}
	}

	for i := 0; i < 7; i++ {
		x := int32(r.NextBoundedInt(16))
		y := int32(r.NextBoundedInt(64))
		z := int32(r.NextBoundedInt(16))

		targetPos := pos.Add(x, y, z)
		currentBlock := level.GetBlockId(targetPos.X(), targetPos.Y(), targetPos.Z())
		if currentBlock == block.STONE {

			level.SetBlock(targetPos.X(), targetPos.Y(), targetPos.Z(), block.MONSTER_EGG, 0, false)
		}
	}
}
