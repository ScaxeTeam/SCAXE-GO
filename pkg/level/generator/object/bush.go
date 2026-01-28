package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Bush struct {
	BlockID   uint8
	BlockMeta uint8
}

func NewBush(id, meta uint8) *Bush {
	return &Bush{BlockID: id, BlockMeta: meta}
}

func (b *Bush) Generate(w populator.ChunkManager, r *rand.Random, pos world.BlockPos) bool {
	for i := 0; i < 64; i++ {
		target := pos.Add(
			int32(r.NextBoundedInt(8)-r.NextBoundedInt(8)),
			int32(r.NextBoundedInt(4)-r.NextBoundedInt(4)),
			int32(r.NextBoundedInt(8)-r.NextBoundedInt(8)),
		)

		if target.Y() < 0 || target.Y() >= 256 {
			continue
		}

		at := w.GetBlockId(target.X(), target.Y(), target.Z())
		below := w.GetBlockId(target.X(), target.Y()-1, target.Z())

		validSoil := below == block.GRASS || below == block.DIRT || below == block.FARMLAND

		if at == 0 && validSoil {

			w.SetBlock(target.X(), target.Y(), target.Z(), b.BlockID, b.BlockMeta, false)
		} else {

		}
	}
	return true
}

type DeadBush struct {
}

func NewDeadBush() *DeadBush {
	return &DeadBush{}
}

func (d *DeadBush) Generate(w populator.ChunkManager, r *rand.Random, pos world.BlockPos) bool {
	for i := 0; i < 4; i++ {
		target := pos.Add(
			int32(r.NextBoundedInt(8)-r.NextBoundedInt(8)),
			int32(r.NextBoundedInt(4)-r.NextBoundedInt(4)),
			int32(r.NextBoundedInt(8)-r.NextBoundedInt(8)),
		)

		if target.Y() < 0 || target.Y() >= 256 {
			continue
		}

		if w.GetBlockId(target.X(), target.Y(), target.Z()) == 0 {
			below := w.GetBlockId(target.X(), target.Y()-1, target.Z())
			if below == block.SAND || below == block.DIRT || below == block.HARDENED_CLAY || below == block.STAINED_CLAY {
				w.SetBlock(target.X(), target.Y(), target.Z(), block.DEAD_BUSH, 0, false)
			}
		}
	}
	return true
}
