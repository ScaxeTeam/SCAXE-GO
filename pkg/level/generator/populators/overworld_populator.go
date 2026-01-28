package populators

import (
	"github.com/scaxe/scaxe-go/pkg/level/generator/biome"
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type OverworldPopulator struct {
	Populators          []populator.Populator
	CavePopulator       *CavePopulator
	OrePopulator        *OrePopulator
	TreePopulator       *TreePopulator
	VegetationPopulator *VegetationPopulator
	SnowPopulator       *SnowPopulator
}

func NewOverworldPopulator() *OverworldPopulator {
	op := &OverworldPopulator{
		Populators:          make([]populator.Populator, 0),
		CavePopulator:       NewCavePopulator(),
		OrePopulator:        NewOrePopulator(make([]object.OreType, 0)),
		TreePopulator:       NewTreePopulator(0),
		VegetationPopulator: NewVegetationPopulator(),
		SnowPopulator:       NewSnowPopulator(),
	}

	return op
}

func (p *OverworldPopulator) Populate(level populator.ChunkManager, chunk *world.Chunk, x, z int32, random *rand.Random) {

	biomeID := chunk.GetBiomeID(8, 8)
	b := biome.GetBiome(biomeID)
	b.Decorate(level, random, world.NewBlockPos(int32(x*16), 0, int32(z*16)))

	p.SnowPopulator.Populate(level, chunk, x, z, random)
}
