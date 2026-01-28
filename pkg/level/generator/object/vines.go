package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Vines struct{}

func NewVines() *Vines {
	return &Vines{}
}

func (v *Vines) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {

	for i := 0; i < 64; i++ {
		target := pos.Add(
			int32(random.NextBoundedInt(8)-random.NextBoundedInt(8)),
			int32(random.NextBoundedInt(4)-random.NextBoundedInt(4)),
			int32(random.NextBoundedInt(8)-random.NextBoundedInt(8)),
		)

		if target.Y() < 0 || target.Y() >= 256 {
			continue
		}

		if level.GetBlockId(target.X(), target.Y(), target.Z()) == 0 {

			meta := 0

			if v.isSolid(level, target.North()) {
				meta |= 1
			}
			if v.isSolid(level, target.South()) {
				meta |= 4
			}
			if v.isSolid(level, target.West()) {
				meta |= 8
			}
			if v.isSolid(level, target.East()) {
				meta |= 2
			}

			if meta != 0 {
				level.SetBlock(target.X(), target.Y(), target.Z(), byte(block.VINE), byte(meta), false)
			}
		}
	}
	return true
}

func (v *Vines) isSolid(level populator.ChunkManager, pos world.BlockPos) bool {
	id := level.GetBlockId(pos.X(), pos.Y(), pos.Z())

	return id != 0 && id != byte(block.VINE) && id != byte(block.TALL_GRASS) && id != byte(block.DEAD_BUSH)
}
