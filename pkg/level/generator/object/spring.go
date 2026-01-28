package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Spring struct {
	BlockID uint8
}

func NewSpring(blockID uint8) *Spring {
	return &Spring{BlockID: blockID}
}

func (s *Spring) Generate(w populator.ChunkManager, r *rand.Random, pos world.BlockPos) bool {
	if w.GetBlockId(pos.X(), pos.Y()+1, pos.Z()) != block.STONE {
		return false
	}
	if w.GetBlockId(pos.X(), pos.Y()-1, pos.Z()) != block.STONE {
		return false
	}

	state := w.GetBlockId(pos.X(), pos.Y(), pos.Z())
	if state != 0 && state != block.STONE {
		return false
	}

	i := 0
	if w.GetBlockId(pos.X()-1, pos.Y(), pos.Z()) == block.STONE {
		i++
	}
	if w.GetBlockId(pos.X()+1, pos.Y(), pos.Z()) == block.STONE {
		i++
	}
	if w.GetBlockId(pos.X(), pos.Y(), pos.Z()-1) == block.STONE {
		i++
	}
	if w.GetBlockId(pos.X(), pos.Y(), pos.Z()+1) == block.STONE {
		i++
	}

	j := 0
	if w.GetBlockId(pos.X()-1, pos.Y(), pos.Z()) == 0 {
		j++
	}
	if w.GetBlockId(pos.X()+1, pos.Y(), pos.Z()) == 0 {
		j++
	}
	if w.GetBlockId(pos.X(), pos.Y(), pos.Z()-1) == 0 {
		j++
	}
	if w.GetBlockId(pos.X(), pos.Y(), pos.Z()+1) == 0 {
		j++
	}

	if i == 3 && j == 1 {
		w.SetBlock(pos.X(), pos.Y(), pos.Z(), s.BlockID, 0, true)

		return true
	}

	return false
}
