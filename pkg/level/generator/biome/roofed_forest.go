package biome

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type RoofedForestBiome struct {
	*BaseBiome
}

func getHeightAt(level populator.ChunkManager, x, z int32) int {
	for y := 255; y >= 0; y-- {
		bid := level.GetBlockId(x, int32(y), z)

		if bid != 0 && bid != 8 && bid != 9 && bid != 10 && bid != 11 {
			return y + 1
		}
	}
	return 64
}

func NewRoofedForestBiome() *RoofedForestBiome {
	b := &RoofedForestBiome{
		BaseBiome: &BaseBiome{
			ID:              ROOFED_FOREST,
			Name:            "Roofed Forest",
			BaseHeight:      0.1,
			HeightVariation: 0.2,
			Temperature:     0.7,
			Rainfall:        0.8,
			GroundCover:     []block.BlockState{block.NewBlockState(block.GRASS, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
	}

	b.Decorator.TreesPerChunk = -999
	b.Decorator.GrassPerChunk = 2
	b.Decorator.FlowersPerChunk = 0
	return b
}

func (b *RoofedForestBiome) GetTreeFeature(r *rand.Random) Generator {

	if r.NextBoundedInt(3) > 0 {
		return object.NewDarkOakTree()
	}

	if r.NextBoundedInt(5) == 0 {
		return object.NewBirchTree(false)
	}

	if r.NextBoundedInt(10) == 0 {
		return object.NewBigOakTree()
	}
	return object.NewOakTree()
}

func (b *RoofedForestBiome) Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos) {

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {

			k := i*4 + 1 + 8 + r.NextBoundedInt(3)
			l := j*4 + 1 + 8 + r.NextBoundedInt(3)

			treeX := pos.X() + int32(k)
			treeZ := pos.Z() + int32(l)
			treeY := int32(getHeightAt(level, treeX, treeZ))

			treePos := world.NewBlockPos(treeX, treeY, treeZ)

			if r.NextBoundedInt(20) == 0 {
				mushroom := object.NewBigMushroom(0)
				mushroom.Generate(level, r, treePos)
			} else {

				tree := b.GetTreeFeature(r)
				tree.Generate(level, r, treePos)
			}
		}
	}

	doublePlantCount := r.NextBoundedInt(5) - 3
	for i := 0; i < doublePlantCount; i++ {
		plantType := r.NextBoundedInt(3)
		var plant Generator
		switch plantType {
		case 0:
			plant = object.NewDoublePlant(object.DoublePlantLilac)
		case 1:
			plant = object.NewDoublePlant(object.DoublePlantRoseBush)
		case 2:
			plant = object.NewDoublePlant(object.DoublePlantPeony)
		}

		for k := 0; k < 5; k++ {
			px := r.NextBoundedInt(16) + 8
			pz := r.NextBoundedInt(16) + 8
			py := r.NextBoundedInt(128 + 32)
			if plant.Generate(level, r, pos.Add(int32(px), int32(py), int32(pz))) {
				break
			}
		}
	}

	b.BaseBiome.Decorate(level, r, pos)
}
