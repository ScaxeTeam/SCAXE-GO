package biome

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

func getHeightAtForBiome(level populator.ChunkManager, x, z int32) int {
	for y := 255; y >= 0; y-- {
		bid := level.GetBlockId(x, int32(y), z)

		if bid != 0 && bid != 8 && bid != 9 && bid != 10 && bid != 11 {
			return y + 1
		}
	}
	return 64
}

type PlainsBiome struct {
	*BaseBiome
}

func NewPlainsBiome() *PlainsBiome {
	b := &PlainsBiome{
		BaseBiome: &BaseBiome{
			ID:              PLAINS,
			Name:            "Plains",
			BaseHeight:      0.125,
			HeightVariation: 0.05,
			Temperature:     0.8,
			Rainfall:        0.4,
			GroundCover:     []block.BlockState{block.NewBlockState(block.GRASS, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
	}
	b.Decorator.TreesPerChunk = 0
	b.Decorator.ExtraTreeChance = 0.05
	b.Decorator.FlowersPerChunk = 4
	b.Decorator.GrassPerChunk = 10
	return b
}

func (b *PlainsBiome) GetTreeFeature(r *rand.Random) Generator {
	return object.NewOakTree()
}

func (b *PlainsBiome) GetFlowerType(r *rand.Random, pos world.BlockPos) Generator {

	noise := GrassColorNoise.GetValue(float64(pos.X())/200.0, float64(pos.Z())/200.0)

	if noise < -0.8 {

		switch r.NextBoundedInt(4) {
		case 0:
			return object.NewBush(block.RED_FLOWER, object.FlowerOrangeTulip)
		case 1:
			return object.NewBush(block.RED_FLOWER, object.FlowerRedTulip)
		case 2:
			return object.NewBush(block.RED_FLOWER, object.FlowerPinkTulip)
		default:
			return object.NewBush(block.RED_FLOWER, object.FlowerWhiteTulip)
		}
	} else if r.NextBoundedInt(3) > 0 {

		switch r.NextBoundedInt(3) {
		case 0:
			return object.NewBush(block.RED_FLOWER, object.FlowerPoppy)
		case 1:
			return object.NewBush(block.RED_FLOWER, object.FlowerAzureBluet)
		default:
			return object.NewBush(block.RED_FLOWER, object.FlowerOxeyeDaisy)
		}
	}

	return object.NewBush(block.DANDELION, 0)
}

type DesertBiome struct {
	*BaseBiome
}

func NewDesertBiome() *DesertBiome {
	b := &DesertBiome{
		BaseBiome: &BaseBiome{
			ID:              DESERT,
			Name:            "Desert",
			BaseHeight:      0.125,
			HeightVariation: 0.05,
			Temperature:     2.0,
			Rainfall:        0.0,
			GroundCover:     []block.BlockState{block.NewBlockState(block.SAND, 0), block.NewBlockState(block.SAND, 0), block.NewBlockState(block.SANDSTONE, 0)},
			Decorator:       NewDecorator(),
		},
	}
	b.Decorator.TreesPerChunk = -999
	b.Decorator.DeadBushPerChunk = 2
	b.Decorator.ReedsPerChunk = 50
	b.Decorator.CactiPerChunk = 10
	return b
}

func (b *DesertBiome) Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos) {
	b.BaseBiome.Decorate(level, r, pos)

}

type ForestBiome struct {
	*BaseBiome
	Type int
}

const (
	FOREST_NORMAL = 0
	FOREST_FLOWER = 1
	FOREST_BIRCH  = 2
	FOREST_ROOFED = 3
)

func NewForestBiome(t int) *ForestBiome {

	id := uint8(FOREST)
	name := "Forest"
	if t == FOREST_BIRCH {
		id = BIRCH_FOREST
		name = "Birch Forest"
	} else if t == FOREST_ROOFED {
		id = ROOFED_FOREST
		name = "Roofed Forest"
	}

	b := &ForestBiome{
		BaseBiome: &BaseBiome{
			ID:              id,
			Name:            name,
			BaseHeight:      0.1,
			HeightVariation: 0.2,
			Temperature:     0.7,
			Rainfall:        0.8,
			GroundCover:     []block.BlockState{block.NewBlockState(block.GRASS, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
		Type: t,
	}

	b.Decorator.TreesPerChunk = 10
	b.Decorator.GrassPerChunk = 2
	b.Decorator.FlowersPerChunk = 4

	if t == FOREST_FLOWER {
		b.Decorator.TreesPerChunk = 6
		b.Decorator.FlowersPerChunk = 100
		b.Decorator.GrassPerChunk = 1

	}

	if t == FOREST_ROOFED {
		b.Decorator.TreesPerChunk = -999

	}

	return b
}

func (b *ForestBiome) GetTreeFeature(r *rand.Random) Generator {

	if b.Type == FOREST_ROOFED && r.NextBoundedInt(3) > 0 {
		return object.NewDarkOakTree()
	}

	if b.Type == FOREST_BIRCH || r.NextBoundedInt(5) == 0 {
		return object.NewBirchTree(false)
	}

	if r.NextBoundedInt(10) == 0 {
		return object.NewBigOakTree()
	}
	return object.NewOakTree()
}

func (b *ForestBiome) Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos) {

	doublePlantCount := r.NextBoundedInt(5) - 3
	if b.Type == FOREST_FLOWER {
		doublePlantCount += 2
	}

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

			terrainHeight := getHeightAtForBiome(level, pos.X()+int32(px), pos.Z()+int32(pz))
			py := r.NextBoundedInt(terrainHeight + 32)
			if plant.Generate(level, r, pos.Add(int32(px), int32(py), int32(pz))) {
				break
			}
		}
	}

	b.BaseBiome.Decorate(level, r, pos)
}

type TaigaBiome struct {
	*BaseBiome
}

func NewTaigaBiome() *TaigaBiome {
	b := &TaigaBiome{
		BaseBiome: &BaseBiome{
			ID:              TAIGA,
			Name:            "Taiga",
			BaseHeight:      0.2,
			HeightVariation: 0.2,
			Temperature:     0.25,
			Rainfall:        0.8,
			GroundCover:     []block.BlockState{block.NewBlockState(block.GRASS, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
	}
	b.Decorator.TreesPerChunk = 10
	b.Decorator.GrassPerChunk = 1
	b.Decorator.MushroomsPerChunk = 1
	return b
}

func (b *TaigaBiome) GetTreeFeature(r *rand.Random) Generator {
	if r.NextBoundedInt(3) == 0 {
		return object.NewPineTree()
	}
	return object.NewSpruceTree()
}

type MegaTaigaBiome struct {
	*BaseBiome
	IsSpruce bool
}

const (
	MEGA_TAIGA_NORMAL = 0
	MEGA_TAIGA_SPRUCE = 1
)

func NewMegaTaigaBiome() *MegaTaigaBiome {
	b := &MegaTaigaBiome{
		BaseBiome: &BaseBiome{
			ID:              MEGA_TAIGA,
			Name:            "Mega Taiga",
			BaseHeight:      0.2,
			HeightVariation: 0.2,
			Temperature:     0.3,
			Rainfall:        0.8,
			GroundCover:     []block.BlockState{block.NewBlockState(block.GRASS, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
		IsSpruce: false,
	}

	b.Decorator.TreesPerChunk = 10
	b.Decorator.GrassPerChunk = 7
	b.Decorator.DeadBushPerChunk = 1
	b.Decorator.MushroomsPerChunk = 3

	b.Decorator.GrassGen = object.NewGrass(2)
	return b
}

func NewMegaTaigaHillsBiome() *MegaTaigaBiome {
	b := &MegaTaigaBiome{
		BaseBiome: &BaseBiome{
			ID:              MEGA_TAIGA_HILLS,
			Name:            "Mega Taiga Hills",
			BaseHeight:      0.45,
			HeightVariation: 0.3,
			Temperature:     0.3,
			Rainfall:        0.8,
			GroundCover:     []block.BlockState{block.NewBlockState(block.GRASS, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
		IsSpruce: true,
	}
	b.Decorator.TreesPerChunk = 10
	b.Decorator.GrassPerChunk = 7
	b.Decorator.DeadBushPerChunk = 1
	b.Decorator.MushroomsPerChunk = 3
	b.Decorator.GrassGen = object.NewGrass(2)
	return b
}

func (b *MegaTaigaBiome) GetTreeFeature(r *rand.Random) Generator {

	if r.NextBoundedInt(3) == 0 {

		if b.IsSpruce || r.NextBoundedInt(13) == 0 {
			return object.NewMegaPineTree(true)
		}
		return object.NewMegaPineTree(false)
	}

	if r.NextBoundedInt(3) == 0 {
		return object.NewPineTree()
	}
	return object.NewSpruceTree()
}

func (b *MegaTaigaBiome) Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos) {

	boulderCount := r.NextBoundedInt(3)
	for i := 0; i < boulderCount; i++ {
		bx := r.NextBoundedInt(16) + 8
		bz := r.NextBoundedInt(16) + 8
		boulder := object.NewMossyCobblestoneBlob()

		boulder.Generate(level, r, pos.Add(int32(bx), 128, int32(bz)))
	}

	for i := 0; i < 7; i++ {
		fx := r.NextBoundedInt(16) + 8
		fz := r.NextBoundedInt(16) + 8
		yMax := 128 + 32
		if yMax > 0 {
			fy := r.NextBoundedInt(yMax)
			fern := object.NewDoublePlant(object.DoublePlantFern)
			fern.Generate(level, r, pos.Add(int32(fx), int32(fy), int32(fz)))
		}
	}

	b.BaseBiome.Decorate(level, r, pos)
}

type SwampBiome struct {
	*BaseBiome
}

func NewSwampBiome() *SwampBiome {
	b := &SwampBiome{
		BaseBiome: &BaseBiome{
			ID:              SWAMP,
			Name:            "Swamp",
			BaseHeight:      -0.2,
			HeightVariation: 0.1,
			Temperature:     0.8,
			Rainfall:        0.9,
			GroundCover:     []block.BlockState{block.NewBlockState(block.GRASS, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
	}
	b.Decorator.TreesPerChunk = 2
	b.Decorator.FlowersPerChunk = 1
	b.Decorator.DeadBushPerChunk = 1
	b.Decorator.MushroomsPerChunk = 8
	b.Decorator.ReedsPerChunk = 10
	b.Decorator.ClayPerChunk = 1
	b.Decorator.WaterlilyPerChunk = 4
	b.Decorator.SandPatchesPerChunk = 0
	b.Decorator.GravelPatchesPerChunk = 0
	b.Decorator.GrassPerChunk = 5
	return b
}

func (b *SwampBiome) GetTreeFeature(r *rand.Random) Generator {
	return object.NewSwampTree()
}

func (b *SwampBiome) GetFlowerType(r *rand.Random, pos world.BlockPos) Generator {
	return object.NewBush(block.RED_FLOWER, object.FlowerBlueOrchid)
}

type IcePlainsBiome struct {
	*BaseBiome
}

func NewIcePlainsBiome() *IcePlainsBiome {
	b := &IcePlainsBiome{
		BaseBiome: &BaseBiome{
			ID:              ICE_PLAINS,
			Name:            "Ice Plains",
			BaseHeight:      0.125,
			HeightVariation: 0.05,
			Temperature:     0.0,
			Rainfall:        0.5,
			GroundCover:     []block.BlockState{block.NewBlockState(block.SNOW_BLOCK, 0), block.NewBlockState(block.DIRT, 0)},

			Decorator: NewDecorator(),
		},
	}
	b.Decorator.TreesPerChunk = 0
	b.Decorator.ExtraTreeChance = 0.05
	b.Decorator.FlowersPerChunk = 0
	b.Decorator.GrassPerChunk = 0
	return b
}

func (b *IcePlainsBiome) GetTreeFeature(r *rand.Random) Generator {
	return object.NewSpruceTree()
}

type IcePlainsSpikesBiome struct {
	*BaseBiome
}

func NewIcePlainsSpikesBiome() *IcePlainsSpikesBiome {
	b := &IcePlainsSpikesBiome{
		BaseBiome: &BaseBiome{
			ID:              ICE_PLAINS_SPIKES,
			Name:            "Ice Plains Spikes",
			BaseHeight:      0.425,
			HeightVariation: 0.45,
			Temperature:     0.0,
			Rainfall:        0.5,
			GroundCover:     []block.BlockState{block.NewBlockState(block.SNOW_BLOCK, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
	}
	b.Decorator.TreesPerChunk = 0
	b.Decorator.ExtraTreeChance = 0.05
	b.Decorator.FlowersPerChunk = 0
	b.Decorator.GrassPerChunk = 0
	return b
}

func (b *IcePlainsSpikesBiome) GetTreeFeature(r *rand.Random) Generator {
	return object.NewSpruceTree()
}

func (b *IcePlainsSpikesBiome) Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos) {

	b.BaseBiome.Decorate(level, r, pos)

	icePath := object.NewIcePath(4)
	for i := 0; i < 3; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		icePath.Generate(level, r, pos.Add(int32(x), 128, int32(z)))
	}

	iceSpike := object.NewIceSpike()
	for i := 0; i < 2; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		iceSpike.Generate(level, r, pos.Add(int32(x), 128, int32(z)))
	}
}

type MushroomIslandBiome struct {
	*BaseBiome
}

func NewMushroomIslandBiome(id uint8) *MushroomIslandBiome {
	name := "Mushroom Island"
	height := 0.2
	variation := 0.3
	if id == MUSHROOM_ISLAND_SHORE {
		name = "Mushroom Island Shore"
		height = 0.0
		variation = 0.025
	}

	b := &MushroomIslandBiome{
		BaseBiome: &BaseBiome{
			ID:              id,
			Name:            name,
			BaseHeight:      height,
			HeightVariation: variation,
			Temperature:     0.9,
			Rainfall:        1.0,
			GroundCover:     []block.BlockState{block.NewBlockState(block.MYCELIUM, 0), block.NewBlockState(block.DIRT, 0)},
			Decorator:       NewDecorator(),
		},
	}
	b.Decorator.TreesPerChunk = 0
	b.Decorator.ExtraTreeChance = 0
	b.Decorator.FlowersPerChunk = 0
	b.Decorator.GrassPerChunk = 0
	b.Decorator.MushroomsPerChunk = 1
	b.Decorator.BigMushroomsPerChunk = 1
	return b
}
