package populators

import (
	"github.com/scaxe/scaxe-go/pkg/level/generator/biome"
	"github.com/scaxe/scaxe-go/pkg/level/generator/biomegrid"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	BlockSnowLayer = 78
	BlockIce       = 79
)

type SnowPopulator struct{}

func NewSnowPopulator() *SnowPopulator {
	return &SnowPopulator{}
}

func (s *SnowPopulator) Populate(level populator.ChunkManager, chunk *world.Chunk, x, z int32, random *rand.Random) {
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			biomeID := chunk.GetBiomeID(x, z)
			if biome.GetBiome(biomeID).GetTemperature() < 0.15 {
				continue
			}

			for y := 127; y > 0; y-- {
				blockID, _ := chunk.GetBlock(x, y, z)

				if blockID == 0 || blockID == 18 || blockID == 161 {
					continue
				}

				if blockID == 9 || blockID == 8 {

					chunk.SetBlock(x, y, z, BlockIce, 0)
					break
				}

				blockAbove, _ := chunk.GetBlock(x, y+1, z)
				if blockAbove == 0 && canSnowOn(blockID) {
					chunk.SetBlock(x, y+1, z, BlockSnowLayer, 0)
				}
				break
			}
		}
	}
}

func isColdBiome(biome int) bool {
	switch biome {
	case biomegrid.BiomeIcePlains,
		biomegrid.BiomeIceMountains,
		biomegrid.BiomeFrozenOcean,
		biomegrid.BiomeFrozenRiver,
		30,
		31:
		return true
	}
	return false
}

func canSnowOn(blockID byte) bool {

	switch blockID {
	case 1, 2, 3, 4, 12, 13, 17, 18, 24, 48, 98, 161, 162:
		return true
	}
	return false
}
