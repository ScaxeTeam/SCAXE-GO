package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type DoublePlant struct {
	Type int
}

func NewDoublePlant(t int) *DoublePlant {
	return &DoublePlant{Type: t}
}

func (d *DoublePlant) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {

	placed := false
	for i := 0; i < 64; i++ {

		x := pos.X() + int32(random.NextBoundedInt(8)) - int32(random.NextBoundedInt(8))
		y := pos.Y() + int32(random.NextBoundedInt(4)) - int32(random.NextBoundedInt(4))
		z := pos.Z() + int32(random.NextBoundedInt(8)) - int32(random.NextBoundedInt(8))

		if y >= 254 || y < 0 {
			continue
		}

		if level.GetBlockId(x, y, z) == 0 && level.GetBlockId(x, y+1, z) == 0 {
			soil := level.GetBlockId(x, y-1, z)
			if soil == byte(block.GRASS) || soil == byte(block.DIRT) || soil == byte(block.FARMLAND) {

				level.SetBlock(x, y, z, byte(block.DOUBLE_PLANT), byte(d.Type), false)

				level.SetBlock(x, y+1, z, byte(block.DOUBLE_PLANT), 8|byte(d.Type&7), false)
				placed = true
			}
		}
	}
	return placed
}
