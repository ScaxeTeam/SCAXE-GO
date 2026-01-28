package populators

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/biome"
	"github.com/scaxe/scaxe-go/pkg/level/generator/biomegrid"
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type VegetationPopulator struct{}

func NewVegetationPopulator() *VegetationPopulator {
	return &VegetationPopulator{}
}

func (vp *VegetationPopulator) Populate(level populator.ChunkManager, chunk *world.Chunk, cx, cz int32, random *rand.Random) {

	centerBiome := int(chunk.GetBiomeID(8, 8))

	baseX := cx * 16
	baseZ := cz * 16

	grassCount := getGrassCount(centerBiome)
	grass := object.NewGrass(1)
	for i := 0; i < grassCount; i++ {
		tryX := random.NextBoundedInt(16)
		tryZ := random.NextBoundedInt(16)
		y := findGround(chunk, int(tryX), int(tryZ))
		if y > 0 {
			grass.Generate(level, random, world.NewBlockPos(int32(baseX)+int32(tryX), int32(y), int32(baseZ)+int32(tryZ)))
		}
	}

	flowerCount := getFlowerCount(centerBiome)
	for i := 0; i < flowerCount; i++ {
		tryX := random.NextBoundedInt(16)
		tryZ := random.NextBoundedInt(16)
		y := findGround(chunk, int(tryX), int(tryZ))
		if y > 0 {
			var flower *object.Bush
			if random.NextBoundedInt(2) == 0 {
				flower = object.NewBush(block.DANDELION, 0)
			} else {
				flower = object.NewBush(block.RED_FLOWER, 0)
			}
			flower.Generate(level, random, world.NewBlockPos(int32(baseX)+int32(tryX), int32(y), int32(baseZ)+int32(tryZ)))
		}
	}

	if centerBiome == biome.DESERT || centerBiome == biome.DESERT_HILLS {
		cactus := object.NewCactus()
		for i := 0; i < 10; i++ {
			tryX := random.NextBoundedInt(16)
			tryZ := random.NextBoundedInt(16)
			y := findGround(chunk, int(tryX), int(tryZ))
			if y > 0 {
				cactus.Generate(level, random, world.NewBlockPos(int32(baseX)+int32(tryX), int32(y), int32(baseZ)+int32(tryZ)))
			}
		}
	}

	if centerBiome == biome.SWAMP || centerBiome == biome.RIVER {
		sugarcane := object.NewReed()
		for i := 0; i < 10; i++ {
			tryX := random.NextBoundedInt(16)
			tryZ := random.NextBoundedInt(16)
			y := findGround(chunk, int(tryX), int(tryZ))
			if y > 0 {
				sugarcane.Generate(level, random, world.NewBlockPos(int32(baseX)+int32(tryX), int32(y), int32(baseZ)+int32(tryZ)))
			}
		}
	}

	if centerBiome == biome.DESERT || centerBiome == biome.MESA {
		deadbush := object.NewDeadBush()
		for i := 0; i < 4; i++ {
			tryX := random.NextBoundedInt(16)
			tryZ := random.NextBoundedInt(16)
			y := findGround(chunk, int(tryX), int(tryZ))
			if y > 0 {
				deadbush.Generate(level, random, world.NewBlockPos(int32(baseX)+int32(tryX), int32(y), int32(baseZ)+int32(tryZ)))
			}
		}
	}
}

func findGround(chunk *world.Chunk, x, z int) int {
	for y := 127; y > 0; y-- {
		blockID, _ := chunk.GetBlock(x, y, z)
		if blockID != 0 && blockID != 9 && blockID != 18 && blockID != 161 {
			return y + 1
		}
	}
	return -1
}

func getGrassCount(biomeID int) int {
	switch biomeID {
	case biome.PLAINS:
		return 20
	case biomegrid.BiomeTaiga, biomegrid.BiomeTaigaHills:
		return 10
	case biomegrid.BiomeJungle, biomegrid.BiomeJungleHills:
		return 50
	case biomegrid.BiomeSavanna:
		return 30
	case biomegrid.BiomeSwampland:
		return 20
	default:
		return 5
	}
}

func getFlowerCount(biomeID int) int {
	switch biomeID {
	case biomegrid.BiomePlains:
		return 4
	case biomegrid.BiomeForest:
		return 3
	case biomegrid.BiomeSwampland:
		return 1
	default:
		return 1
	}
}
