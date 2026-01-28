package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Waterlily struct{}

func NewWaterlily() *Waterlily {
	return &Waterlily{}
}

func (wl *Waterlily) Generate(w populator.ChunkManager, r *rand.Random, pos world.BlockPos) bool {
	for i := 0; i < 10; i++ {
		target := pos.Add(
			int32(r.NextBoundedInt(8)-r.NextBoundedInt(8)),
			int32(r.NextBoundedInt(4)-r.NextBoundedInt(4)),
			int32(r.NextBoundedInt(8)-r.NextBoundedInt(8)),
		)

		if target.Y() < 0 || target.Y() >= 256 {
			continue
		}

		if w.GetBlockId(target.X(), target.Y(), target.Z()) == 0 &&
			(w.GetBlockId(target.X(), target.Y()-1, target.Z()) == block.WATER || w.GetBlockId(target.X(), target.Y()-1, target.Z()) == block.STILL_WATER) {
			w.SetBlock(target.X(), target.Y(), target.Z(), block.WATER_LILY, 0, false)
		}
	}
	return true
}
