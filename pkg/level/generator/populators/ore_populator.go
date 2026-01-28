package populators

import (
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type OrePopulator struct {
	ores []object.OreType
}

func NewOrePopulator(ores []object.OreType) *OrePopulator {
	return &OrePopulator{
		ores: ores,
	}
}

func (op *OrePopulator) Populate(level populator.ChunkManager, chunk *world.Chunk, cx, cz int32, random *rand.Random) {
	chunkX := int(chunk.X) * 16
	chunkZ := int(chunk.Z) * 16

	for _, oreType := range op.ores {

		oreObject := object.NewOre(uint8(oreType.Material), uint8(oreType.Meta), oreType.ClusterSize)

		for i := 0; i < oreType.ClusterCount; i++ {

			rx := chunkX + random.NextBoundedInt(16)
			ry := oreType.MinHeight + random.NextBoundedInt(oreType.MaxHeight-oreType.MinHeight+1)
			rz := chunkZ + random.NextBoundedInt(16)

			oreObject.Generate(level, random, world.NewBlockPos(int32(rx), int32(ry), int32(rz)))
		}
	}
}
