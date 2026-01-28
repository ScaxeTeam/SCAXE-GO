package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Cactus struct{}

func NewCactus() *Cactus {
	return &Cactus{}
}

func (c *Cactus) Generate(w populator.ChunkManager, r *rand.Random, pos world.BlockPos) bool {
	for i := 0; i < 10; i++ {
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
			if below == block.SAND || below == block.CACTUS {

				h := 1 + r.NextBoundedInt(r.NextBoundedInt(3)+1)
				for j := 0; j < h; j++ {
					if w.GetBlockId(target.X(), target.Y()+int32(j), target.Z()) == 0 {
						w.SetBlock(target.X(), target.Y()+int32(j), target.Z(), block.CACTUS, 0, false)
					}
				}
			}
		}
	}
	return true
}

type Reed struct{}

func NewReed() *Reed {
	return &Reed{}
}

func (re *Reed) Generate(w populator.ChunkManager, r *rand.Random, pos world.BlockPos) bool {
	for i := 0; i < 20; i++ {
		target := pos.Add(
			int32(r.NextBoundedInt(4)-r.NextBoundedInt(4)),
			0,
			int32(r.NextBoundedInt(4)-r.NextBoundedInt(4)),
		)

		if target.Y() < 0 || target.Y() >= 256 {
			continue
		}

		if w.GetBlockId(target.X(), target.Y(), target.Z()) == 0 {
			below := w.GetBlockId(target.X(), target.Y()-1, target.Z())
			if below == block.GRASS || below == block.DIRT || below == block.SAND {

				hasWater := false
				if w.GetBlockId(target.X()-1, target.Y()-1, target.Z()) == block.WATER || w.GetBlockId(target.X()-1, target.Y()-1, target.Z()) == block.STILL_WATER {
					hasWater = true
				} else if w.GetBlockId(target.X()+1, target.Y()-1, target.Z()) == block.WATER || w.GetBlockId(target.X()+1, target.Y()-1, target.Z()) == block.STILL_WATER {
					hasWater = true
				} else if w.GetBlockId(target.X(), target.Y()-1, target.Z()-1) == block.WATER || w.GetBlockId(target.X(), target.Y()-1, target.Z()-1) == block.STILL_WATER {
					hasWater = true
				} else if w.GetBlockId(target.X(), target.Y()-1, target.Z()+1) == block.WATER || w.GetBlockId(target.X(), target.Y()-1, target.Z()+1) == block.STILL_WATER {
					hasWater = true
				}

				if hasWater {
					h := 2 + r.NextBoundedInt(r.NextBoundedInt(3)+1)
					for j := 0; j < h; j++ {
						if w.GetBlockId(target.X(), target.Y()+int32(j), target.Z()) == 0 {
							w.SetBlock(target.X(), target.Y()+int32(j), target.Z(), block.SUGARCANE_BLOCK, 0, false)
						}
					}
				}
			}
		}
	}
	return true
}

type Pumpkin struct{}

func NewPumpkin() *Pumpkin { return &Pumpkin{} }

func (p *Pumpkin) Generate(w populator.ChunkManager, r *rand.Random, pos world.BlockPos) bool {
	for i := 0; i < 64; i++ {
		target := pos.Add(
			int32(r.NextBoundedInt(8)-r.NextBoundedInt(8)),
			int32(r.NextBoundedInt(4)-r.NextBoundedInt(4)),
			int32(r.NextBoundedInt(8)-r.NextBoundedInt(8)),
		)

		if target.Y() < 0 || target.Y() >= 256 {
			continue
		}

		if w.GetBlockId(target.X(), target.Y(), target.Z()) == 0 && w.GetBlockId(target.X(), target.Y()-1, target.Z()) == block.GRASS {

			w.SetBlock(target.X(), target.Y(), target.Z(), block.PUMPKIN, uint8(r.NextBoundedInt(4)), false)
		}
	}
	return true
}
