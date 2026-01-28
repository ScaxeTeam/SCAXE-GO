package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	HUGE_MUSHROOM_ALL_INSIDE  = 0
	HUGE_MUSHROOM_NORTH_WEST  = 1
	HUGE_MUSHROOM_NORTH       = 2
	HUGE_MUSHROOM_NORTH_EAST  = 3
	HUGE_MUSHROOM_WEST        = 4
	HUGE_MUSHROOM_CENTER      = 5
	HUGE_MUSHROOM_EAST        = 6
	HUGE_MUSHROOM_SOUTH_WEST  = 7
	HUGE_MUSHROOM_SOUTH       = 8
	HUGE_MUSHROOM_SOUTH_EAST  = 9
	HUGE_MUSHROOM_STEM        = 10
	HUGE_MUSHROOM_ALL_OUTSIDE = 14
	HUGE_MUSHROOM_ALL_STEM    = 15
)

type BigMushroom struct {
	MushroomType uint8
}

var mushroomOverridable = map[uint8]bool{
	block.AIR:        true,
	block.SAPLING:    true,
	block.LOG:        true,
	block.LEAVES:     true,
	block.SNOW_LAYER: true,
	block.WOOD2:      true,
	block.LEAVES2:    true,
}

func NewBigMushroom(mushroomType uint8) *BigMushroom {
	return &BigMushroom{MushroomType: mushroomType}
}

func (g *BigMushroom) Generate(w populator.ChunkManager, r *rand.Random, pos world.BlockPos) bool {
	blockID := g.MushroomType
	if blockID == 0 {
		if r.NextBoolean() {
			blockID = block.BROWN_MUSHROOM_BLOCK
		} else {
			blockID = block.RED_MUSHROOM_BLOCK
		}
	}

	i := r.NextBoundedInt(3) + 4

	if r.NextBoundedInt(12) == 0 {
		i *= 2
	}

	flag := true

	if pos.Y() >= 1 && pos.Y()+int32(i)+1 < 256 {
		for y := pos.Y(); y <= pos.Y()+int32(i)+1; y++ {
			k := 3
			if y <= pos.Y()+3 {
				k = 0
			}

			for x := pos.X() - int32(k); x <= pos.X()+int32(k) && flag; x++ {
				for z := pos.Z() - int32(k); z <= pos.Z()+int32(k) && flag; z++ {
					if y >= 0 && y < 256 {
						material := w.GetBlockId(int32(x), int32(y), int32(z))

						if material != block.AIR && material != block.LEAVES && material != block.LEAVES2 {
							flag = false
						}
					} else {
						flag = false
					}
				}
			}
		}

		if !flag {
			return false
		}

		soil := w.GetBlockId(pos.X(), pos.Y()-1, pos.Z())
		if soil != block.DIRT && soil != block.GRASS && soil != block.MYCELIUM {
			return false
		}

		k2 := pos.Y() + int32(i)
		if blockID == block.RED_MUSHROOM_BLOCK {
			k2 = pos.Y() + int32(i) - 3
		}

		for l2 := k2; l2 <= pos.Y()+int32(i); l2++ {
			j3 := 1
			if l2 < pos.Y()+int32(i) {
				j3++
			}
			if blockID == block.BROWN_MUSHROOM_BLOCK {
				j3 = 3
			}

			k3 := pos.X() - int32(j3)
			l3 := pos.X() + int32(j3)
			j1 := pos.Z() - int32(j3)
			k1 := pos.Z() + int32(j3)

			for l1 := k3; l1 <= l3; l1++ {
				for i2 := j1; i2 <= k1; i2++ {
					j2 := 5

					if l1 == k3 {
						j2--
					} else if l1 == l3 {
						j2++
					}

					if i2 == j1 {
						j2 -= 3
					} else if i2 == k1 {
						j2 += 3
					}

					meta := j2

					if blockID == block.BROWN_MUSHROOM_BLOCK || l2 < pos.Y()+int32(i) {

						if (l1 == k3 || l1 == l3) && (i2 == j1 || i2 == k1) {
							continue
						}

						if l1 == pos.X()-(int32(j3)-1) && i2 == j1 {
							meta = HUGE_MUSHROOM_NORTH_WEST
						}
						if l1 == k3 && i2 == pos.Z()-(int32(j3)-1) {
							meta = HUGE_MUSHROOM_NORTH_WEST
						}
						if l1 == pos.X()+(int32(j3)-1) && i2 == j1 {
							meta = HUGE_MUSHROOM_NORTH_EAST
						}
						if l1 == l3 && i2 == pos.Z()-(int32(j3)-1) {
							meta = HUGE_MUSHROOM_NORTH_EAST
						}
						if l1 == pos.X()-(int32(j3)-1) && i2 == k1 {
							meta = HUGE_MUSHROOM_SOUTH_WEST
						}
						if l1 == k3 && i2 == pos.Z()+(int32(j3)-1) {
							meta = HUGE_MUSHROOM_SOUTH_WEST
						}
						if l1 == pos.X()+(int32(j3)-1) && i2 == k1 {
							meta = HUGE_MUSHROOM_SOUTH_EAST
						}
						if l1 == l3 && i2 == pos.Z()+(int32(j3)-1) {
							meta = HUGE_MUSHROOM_SOUTH_EAST
						}
					}

					if meta == HUGE_MUSHROOM_CENTER && l2 < pos.Y()+int32(i) {
						meta = HUGE_MUSHROOM_ALL_INSIDE
					}

					if pos.Y() >= pos.Y()+int32(i)-1 || meta != HUGE_MUSHROOM_ALL_INSIDE {

						existingBlock := w.GetBlockId(int32(l1), int32(l2), int32(i2))
						if mushroomOverridable[existingBlock] {
							w.SetBlock(int32(l1), int32(l2), int32(i2), blockID, uint8(meta), false)
						}
					}
				}
			}
		}

		for i3 := int32(0); i3 < int32(i); i3++ {
			pY := pos.Y() + i3
			pBlock := w.GetBlockId(pos.X(), pY, pos.Z())

			if mushroomOverridable[pBlock] {
				w.SetBlock(pos.X(), pY, pos.Z(), blockID, HUGE_MUSHROOM_STEM, false)
			}
		}

		return true
	}
	return false
}
