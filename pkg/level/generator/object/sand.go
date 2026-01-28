package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Sand struct {
	BlockID uint8
	Radius  int
}

func NewSand(blockID uint8, radius int) *Sand {
	return &Sand{
		BlockID: blockID,
		Radius:  radius,
	}
}

func (s *Sand) Generate(w populator.ChunkManager, r *rand.Random, pos world.BlockPos) bool {
	if w.GetBlockId(pos.X(), pos.Y(), pos.Z()) != block.WATER && w.GetBlockId(pos.X(), pos.Y(), pos.Z()) != block.STILL_WATER {
		return false
	}

	i := r.NextBoundedInt(s.Radius-2) + 2
	j := 2

	for k := pos.X() - int32(i); k <= pos.X()+int32(i); k++ {
		for l := pos.Z() - int32(i); l <= pos.Z()+int32(i); l++ {
			dx := float64(k - pos.X())
			dz := float64(l - pos.Z())

			if dx*dx+dz*dz <= float64(i*i) {
				for m := pos.Y() - int32(j); m <= pos.Y()+int32(j); m++ {
					id := w.GetBlockId(k, m, l)

					if id == block.DIRT || id == block.GRASS {
						w.SetBlock(k, m, l, s.BlockID, 0, false)
					}
				}
			}
		}
	}
	return true
}
