package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Lake struct {
	BlockID uint8
}

func NewLake(blockID uint8) *Lake {
	return &Lake{BlockID: blockID}
}

func (l *Lake) Generate(w populator.ChunkManager, r *rand.Random, pos world.BlockPos) bool {

	x := pos.X() - 8
	y := pos.Y()
	z := pos.Z() - 8

	for y > 5 && w.GetBlockId(x, y, z) == 0 {
		y--
	}

	if y <= 4 {
		return false
	}

	y -= 4

	aboolean := make([]bool, 2048)
	i := r.NextBoundedInt(4) + 4

	for j := 0; j < i; j++ {
		d0 := r.NextDouble()*6.0 + 3.0
		d1 := r.NextDouble()*4.0 + 2.0
		d2 := r.NextDouble()*6.0 + 3.0
		d3 := r.NextDouble()*(16.0-d0-2.0) + 1.0 + d0/2.0
		d4 := r.NextDouble()*(8.0-d1-4.0) + 2.0 + d1/2.0
		d5 := r.NextDouble()*(16.0-d2-2.0) + 1.0 + d2/2.0

		for l := 1; l < 15; l++ {
			for i1 := 1; i1 < 15; i1++ {
				for j1 := 1; j1 < 7; j1++ {
					d6 := (float64(l) - d3) / (d0 / 2.0)
					d7 := (float64(j1) - d4) / (d1 / 2.0)
					d8 := (float64(i1) - d5) / (d2 / 2.0)
					d9 := d6*d6 + d7*d7 + d8*d8

					if d9 < 1.0 {
						aboolean[(l*16+i1)*8+j1] = true
					}
				}
			}
		}
	}

	for k1 := 0; k1 < 16; k1++ {
		for l2 := 0; l2 < 16; l2++ {
			for k := 0; k < 8; k++ {
				flag := !aboolean[(k1*16+l2)*8+k] && ((k1 < 15 && aboolean[((k1+1)*16+l2)*8+k]) ||
					(k1 > 0 && aboolean[((k1-1)*16+l2)*8+k]) ||
					(l2 < 15 && aboolean[(k1*16+l2+1)*8+k]) ||
					(l2 > 0 && aboolean[(k1*16+(l2-1))*8+k]) ||
					(k < 7 && aboolean[(k1*16+l2)*8+k+1]) ||
					(k > 0 && aboolean[(k1*16+l2)*8+(k-1)]))

				if flag {
					mat := w.GetBlockId(x+int32(k1), y+int32(k), z+int32(l2))

					if k >= 4 && (mat == block.WATER || mat == block.STILL_WATER || mat == block.LAVA || mat == block.STILL_LAVA) {
						return false
					}

					if k < 4 && mat != 0 && mat != l.BlockID {

						if mat == 0 {
							return false
						}

					}
				}
			}
		}
	}

	for l1 := 0; l1 < 16; l1++ {
		for i3 := 0; i3 < 16; i3++ {
			for i4 := 0; i4 < 8; i4++ {
				if aboolean[(l1*16+i3)*8+i4] {
					bx := x + int32(l1)
					by := y + int32(i4)
					bz := z + int32(i3)

					existing := w.GetBlockId(bx, by, bz)
					if existing == block.WATER || existing == block.STILL_WATER {
						continue
					}

					above := w.GetBlockId(bx, by+1, bz)
					if above == block.WATER || above == block.STILL_WATER {
						continue
					}

					targetBlock := l.BlockID
					if i4 >= 4 {
						targetBlock = 0
					}
					w.SetBlock(bx, by, bz, targetBlock, 0, false)
				}
			}
		}
	}

	return true
}
