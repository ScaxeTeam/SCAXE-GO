package biome

import (
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type SunflowerPlainsBiome struct {
	*PlainsBiome
}

func NewSunflowerPlainsBiome() *SunflowerPlainsBiome {
	plains := NewPlainsBiome()
	plains.ID = 129
	plains.Name = "Sunflower Plains"
	return &SunflowerPlainsBiome{
		PlainsBiome: plains,
	}
}

func (b *SunflowerPlainsBiome) Decorate(level populator.ChunkManager, r *rand.Random, pos world.BlockPos) {

	b.PlainsBiome.Decorate(level, r, pos)

	sunflower := object.NewDoublePlant(object.DoublePlantSunflower)
	for i := 0; i < 10; i++ {
		x := pos.X() + int32(r.NextBoundedInt(16)) + 8
		z := pos.Z() + int32(r.NextBoundedInt(16)) + 8

		terrainHeight := level.GetHeight(x, z)
		yMax := int(terrainHeight) + 32
		if yMax > 0 {
			y := r.NextBoundedInt(yMax)

			sunflower.Generate(level, r, world.NewBlockPos(x, int32(y), z))
		}
	}
}

func (b *SunflowerPlainsBiome) GetID() uint8 {
	return 129
}
