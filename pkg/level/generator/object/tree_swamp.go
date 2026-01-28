package object

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type SwampTree struct {
	*Tree
}

func NewSwampTree() *SwampTree {

	t := NewBaseTree(block.LOG, block.LEAVES, 0)
	return &SwampTree{Tree: t}
}

func (t *SwampTree) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {

	height := random.NextBoundedInt(4) + 5

	for level.GetBlockId(pos.X(), pos.Y()-1, pos.Z()) == byte(block.STILL_WATER) || level.GetBlockId(pos.X(), pos.Y()-1, pos.Z()) == byte(block.WATER) {
		pos = pos.Down()
	}

	x, y, z := pos.X(), pos.Y(), pos.Z()

	if y < 1 || y+int32(height)+1 > 256 {
		return false
	}

	for yy := y; yy <= y+1+int32(height); yy++ {
		radius := 1
		if yy == y {
			radius = 0
		}
		if yy >= y+1+int32(height)-2 {
			radius = 3
		}

		for xx := x - int32(radius); xx <= x+int32(radius); xx++ {
			for zz := z - int32(radius); zz <= z+int32(radius); zz++ {
				if yy >= 0 && yy < 256 {
					id := level.GetBlockId(xx, yy, zz)

					if id != 0 && id != byte(block.LEAVES) {
						if id != byte(block.STILL_WATER) && id != byte(block.WATER) {
							return false
						} else if yy > y {
							return false
						}
					}
				} else {
					return false
				}
			}
		}
	}

	soil := level.GetBlockId(x, y-1, z)
	if (soil == byte(block.GRASS) || soil == byte(block.DIRT)) && y < 256-int32(height)-1 {

		level.SetBlock(x, y-1, z, byte(block.DIRT), 0, false)

		for yy := y - 3 + int32(height); yy <= y+int32(height); yy++ {
			yOff := yy - (y + int32(height))

			mid := 2 - int(yOff)/2

			for xx := x - int32(mid); xx <= x+int32(mid); xx++ {
				xOff := xx - x
				for zz := z - int32(mid); zz <= z+int32(mid); zz++ {
					zOff := zz - z

					if math.Abs(float64(xOff)) != float64(mid) || math.Abs(float64(zOff)) != float64(mid) || (random.NextBoundedInt(2) != 0 && yOff != 0) {

						level.SetBlock(xx, yy, zz, byte(block.LEAVES), 0, false)
					}
				}
			}
		}

		for yy := 0; yy < height; yy++ {
			id := level.GetBlockId(x, y+int32(yy), z)
			if id == 0 || id == byte(block.LEAVES) || id == byte(block.STILL_WATER) || id == byte(block.WATER) {
				level.SetBlock(x, y+int32(yy), z, byte(block.LOG), 0, false)
			}
		}

		for yy := y - 3 + int32(height); yy <= y+int32(height); yy++ {
			yOff := yy - (y + int32(height))
			mid := 2 - int(yOff)/2

			for xx := x - int32(mid); xx <= x+int32(mid); xx++ {
				for zz := z - int32(mid); zz <= z+int32(mid); zz++ {
					if level.GetBlockId(xx, yy, zz) == byte(block.LEAVES) {

						if random.NextBoundedInt(4) == 0 && level.GetBlockId(xx-1, yy, zz) == 0 {
							t.addVine(level, xx-1, yy, zz, 8)

						}
						if random.NextBoundedInt(4) == 0 && level.GetBlockId(xx+1, yy, zz) == 0 {
							t.addVine(level, xx+1, yy, zz, 2)
						}

						if random.NextBoundedInt(4) == 0 && level.GetBlockId(xx, yy, zz-1) == 0 {
							t.addVine(level, xx, yy, zz-1, 1)
						}

						if random.NextBoundedInt(4) == 0 && level.GetBlockId(xx, yy, zz+1) == 0 {
							t.addVine(level, xx, yy, zz+1, 4)
						}
					}
				}
			}
		}

		return true
	}
	return false
}

func (t *SwampTree) addVine(level populator.ChunkManager, x, y, z int32, meta int) {
	level.SetBlock(x, y, z, byte(block.VINE), byte(meta), false)

	for i := 0; i < 4; i++ {
		y--
		if level.GetBlockId(x, y, z) == 0 {
			level.SetBlock(x, y, z, byte(block.VINE), byte(meta), false)
		} else {
			break
		}
	}
}
