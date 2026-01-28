package biome

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type JungleBiome struct {
	*BaseBiome
	IsEdge bool
}

func NewJungleBiome() *JungleBiome {
	b := &JungleBiome{
		BaseBiome: &BaseBiome{
			ID:              JUNGLE,
			Name:            "Jungle",
			BaseHeight:      0.1,
			HeightVariation: 0.2,
			Temperature:     0.95,
			Rainfall:        0.9,
			GroundCover:     []block.BlockState{block.NewBlockState(block.GRASS, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
		IsEdge: false,
	}
	b.Decorator.TreesPerChunk = 50
	b.Decorator.GrassPerChunk = 25
	b.Decorator.FlowersPerChunk = 4

	return b
}

func NewJungleEdgeBiome() *JungleBiome {
	b := &JungleBiome{
		BaseBiome: &BaseBiome{
			ID:              JUNGLE_EDGE,
			Name:            "Jungle Edge",
			BaseHeight:      0.1,
			HeightVariation: 0.2,
			Temperature:     0.95,
			Rainfall:        0.8,
			GroundCover:     []block.BlockState{block.NewBlockState(block.GRASS, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
		IsEdge: true,
	}
	b.Decorator.TreesPerChunk = 2
	b.Decorator.GrassPerChunk = 3
	b.Decorator.FlowersPerChunk = 2
	return b
}

func NewJungleHillsBiome() *JungleBiome {
	b := &JungleBiome{
		BaseBiome: &BaseBiome{
			ID:              JUNGLE_HILLS,
			Name:            "Jungle Hills",
			BaseHeight:      0.45,
			HeightVariation: 0.3,
			Temperature:     0.95,
			Rainfall:        0.9,
			GroundCover:     []block.BlockState{block.NewBlockState(block.GRASS, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
		IsEdge: false,
	}
	b.Decorator.TreesPerChunk = 50
	b.Decorator.GrassPerChunk = 25
	b.Decorator.FlowersPerChunk = 4
	return b
}

func (b *JungleBiome) GetTreeFeature(r *rand.Random) Generator {

	if r.NextBoundedInt(10) == 0 {

		return object.NewBigOakTree()
	}

	if r.NextBoundedInt(2) == 0 {

		return object.NewJungleBush(block.LOG, 3, block.LEAVES, 3)
	}

	if !b.IsEdge && r.NextBoundedInt(3) == 0 {

		return object.NewMegaJungleTree()
	}

	return object.NewJungleSmallTree()
}

func (b *JungleBiome) Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos) {
	b.BaseBiome.Decorate(level, r, pos)

	x := r.NextBoundedInt(16) + 8
	z := r.NextBoundedInt(16) + 8
	yMax := 128 * 2
	if yMax > 0 {
		y := r.NextBoundedInt(yMax)
		melon := object.NewBlockPatch(block.MELON_BLOCK)
		melon.Generate(level, r, pos.Add(int32(x), int32(y), int32(z)))
	}

	vines := object.NewVines()
	for i := 0; i < 50; i++ {
		vx := r.NextBoundedInt(16) + 8
		vz := r.NextBoundedInt(16) + 8

		vines.Generate(level, r, pos.Add(int32(vx), 128, int32(vz)))
	}
}
