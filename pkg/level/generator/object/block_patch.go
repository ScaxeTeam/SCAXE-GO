package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type BlockPatch struct {
	BlockID uint8
}

func NewBlockPatch(id uint8) *BlockPatch {
	return &BlockPatch{BlockID: id}
}

func (bp *BlockPatch) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {
	for i := 0; i < 64; i++ {
		target := pos.Add(
			int32(random.NextBoundedInt(8)-random.NextBoundedInt(8)),
			int32(random.NextBoundedInt(4)-random.NextBoundedInt(4)),
			int32(random.NextBoundedInt(8)-random.NextBoundedInt(8)),
		)

		if target.Y() < 0 || target.Y() >= 256 {
			continue
		}

		if level.GetBlockId(target.X(), target.Y(), target.Z()) == 0 &&
			level.GetBlockId(target.X(), target.Y()-1, target.Z()) == byte(block.GRASS) {

			level.SetBlock(target.X(), target.Y(), target.Z(), bp.BlockID, 0, false)

			meta := byte(random.NextBoundedInt(4))
			level.SetBlock(target.X(), target.Y(), target.Z(), bp.BlockID, meta, false)
		}
	}
	return true
}
