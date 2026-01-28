package biome

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Generator interface {
	Generate(w populator.ChunkManager, r *rand.Random, pos world.BlockPos) bool
}

const (
	OCEAN                    = 0
	PLAINS                   = 1
	DESERT                   = 2
	MOUNTAINS                = 3
	FOREST                   = 4
	TAIGA                    = 5
	SWAMP                    = 6
	RIVER                    = 7
	HELL                     = 8
	END                      = 9
	FROZEN_OCEAN             = 10
	FROZEN_RIVER             = 11
	ICE_PLAINS               = 12
	ICE_MOUNTAINS            = 13
	MUSHROOM_ISLAND          = 14
	MUSHROOM_ISLAND_SHORE    = 15
	BEACH                    = 16
	DESERT_HILLS             = 17
	FOREST_HILLS             = 18
	TAIGA_HILLS              = 19
	SMALL_MOUNTAINS          = 20
	JUNGLE                   = 21
	JUNGLE_HILLS             = 22
	JUNGLE_EDGE              = 23
	DEEP_OCEAN               = 24
	STONE_BEACH              = 25
	COLD_BEACH               = 26
	BIRCH_FOREST             = 27
	BIRCH_FOREST_HILLS       = 28
	ROOFED_FOREST            = 29
	COLD_TAIGA               = 30
	COLD_TAIGA_HILLS         = 31
	MEGA_TAIGA               = 32
	MEGA_TAIGA_HILLS         = 33
	EXTREME_HILLS_PLUS_TREES = 34
	SAVANNA                  = 35
	SAVANNA_PLATEAU          = 36
	MESA                     = 37
	MESA_PLATEAU_F           = 38
	MESA_PLATEAU             = 39

	ICE_PLAINS_SPIKES = 140
)

type Biome interface {
	GetID() uint8
	GetName() string
	GetMinElevation() float64
	GetMaxElevation() float64
	GetGroundCover() []block.BlockState
	GetColor() int
	GetTemperature() float64
	GetRainfall() float64

	GenTerrainBlocks(chunk *world.Chunk, r *rand.Random, x, z int, noiseVal float64)

	Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos)

	GetTreeFeature(r *rand.Random) Generator
	GetFlowerType(r *rand.Random, pos world.BlockPos) Generator
	GetDecorator() *Decorator
}

type BaseBiome struct {
	ID              uint8
	Name            string
	BaseHeight      float64
	HeightVariation float64
	Rainfall        float64
	Temperature     float64
	GroundCover     []block.BlockState

	Decorator *Decorator
}

func (b *BaseBiome) GetID() uint8 {
	return b.ID
}

func (b *BaseBiome) GetName() string {
	return b.Name
}

func (b *BaseBiome) GetMinElevation() float64 {
	return b.BaseHeight
}

func (b *BaseBiome) GetMaxElevation() float64 {
	return b.HeightVariation
}

func (b *BaseBiome) GetGroundCover() []block.BlockState {
	return b.GroundCover
}

func (b *BaseBiome) GetTemperature() float64 {
	return b.Temperature
}

func (b *BaseBiome) GetRainfall() float64 {
	return b.Rainfall
}

func (b *BaseBiome) GetColor() int {
	return GenerateBiomeColor(b.Temperature, b.Rainfall)
}

func (b *BaseBiome) Populate(level populator.ChunkManager, chunk *world.Chunk, x, z int32, r *rand.Random) {

}

func (b *BaseBiome) AddPopulator(pop populator.Populator) {

}

func (b *BaseBiome) Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos) {
	if b.Decorator != nil {
		b.Decorator.Decorate(level, r, b, pos)
	}
}

func (b *BaseBiome) GetDecorator() *Decorator {
	return b.Decorator
}

func (b *BaseBiome) GetTreeFeature(r *rand.Random) Generator {
	return object.NewOakTree()
}

func (b *BaseBiome) GetFlowerType(r *rand.Random, pos world.BlockPos) Generator {
	if r.NextBoundedInt(3) > 0 {
		return b.Decorator.FlowerYGen
	}
	return b.Decorator.FlowerRGen
}

func (b *BaseBiome) GenTerrainBlocks(chunk *world.Chunk, r *rand.Random, x, z int, noiseVal float64) {

	topBlock := byte(2)
	fillerBlock := byte(3)
	stone := byte(1)
	air := byte(0)
	bedrock := byte(7)

	seaLevel := 63

	chunkX := x & 15
	chunkZ := z & 15

	depth := int(noiseVal/3.0 + 3.0 + r.NextDouble()*0.25)

	run := -1
	currentTop := topBlock
	currentFiller := fillerBlock

	for y := 255; y >= 0; y-- {

		if y <= r.NextBoundedInt(5) {

			chunk.SetBlock(chunkZ, y, chunkX, bedrock, 0)
			continue
		}

		id, _ := chunk.GetBlock(chunkX, y, chunkZ)

		if id == air {
			run = -1
		} else if id == stone {
			if run == -1 {

				if depth <= 0 {
					currentTop = air
					currentFiller = stone
				} else if y >= seaLevel-4 && y <= seaLevel+1 {
					currentTop = topBlock
					currentFiller = fillerBlock
				}

				run = depth

				if y >= seaLevel-1 {
					chunk.SetBlock(chunkX, y, chunkZ, currentTop, 0)
				} else if y < seaLevel-7-depth {

					currentTop = air
					currentFiller = stone
					chunk.SetBlock(chunkX, y, chunkZ, 13, 0)
				} else {
					chunk.SetBlock(chunkX, y, chunkZ, currentFiller, 0)
				}
			} else if run > 0 {
				run--
				chunk.SetBlock(chunkX, y, chunkZ, currentFiller, 0)

				if run == 0 && currentFiller == 12 && depth > 1 {
					extraDepth := r.NextBoundedInt(4)
					if y-63 > 0 {
						extraDepth += y - 63
					}
					run = extraDepth
				}
			}
		}
	}
}

func GenerateBiomeColor(temperature, rainfall float64) int {
	x := (1 - temperature) * 255
	z := (1 - rainfall*temperature) * 255

	c := interpolateColor(256, x, z, []int{0x47, 0xd0, 0x33}, []int{0x6c, 0xb4, 0x93}, []int{0xbf, 0xb6, 0x55}, []int{0x80, 0xb4, 0x97})
	return (0xFF << 24) | (c[0] << 16) | (c[1] << 8) | c[2]
}

func interpolateColor(size float64, x, z float64, c1, c2, c3, c4 []int) []int {
	l1 := lerpColor(c1, c2, x/size)
	l2 := lerpColor(c3, c4, x/size)
	return lerpColor(l1, l2, z/size)
}

func lerpColor(a, b []int, s float64) []int {
	invs := 1.0 - s
	return []int{
		int(float64(a[0])*invs + float64(b[0])*s),
		int(float64(a[1])*invs + float64(b[1])*s),
		int(float64(a[2])*invs + float64(b[2])*s),
	}
}

func IsFreezing(id uint8) bool {
	switch id {
	case FROZEN_OCEAN, FROZEN_RIVER, ICE_PLAINS, ICE_MOUNTAINS,
		COLD_BEACH, COLD_TAIGA, COLD_TAIGA_HILLS:
		return true
	}
	return false
}
