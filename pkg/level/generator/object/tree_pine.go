package object

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type PineTree struct {
	*Tree
}

func NewPineTree() *PineTree {

	t := NewBaseTree(block.LOG, block.LEAVES, 1)
	return &PineTree{Tree: t}
}

func (t *PineTree) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {

	height := random.NextBoundedInt(5) + 7

	foliageStart := height - random.NextBoundedInt(2) - 3

	if pos.Y() < 1 || pos.Y()+int32(height)+1 > 256 {
		return false
	}

	x, y, z := pos.X(), pos.Y(), pos.Z()

	soil := level.GetBlockId(x, y-1, z)
	if soil != byte(block.GRASS) && soil != byte(block.DIRT) && soil != byte(block.FARMLAND) && soil != byte(block.PODZOL) {
		return false
	}

	if soil == byte(block.GRASS) {
		level.SetBlock(x, y-1, z, byte(block.DIRT), 0, false)
	}

	radius := 0
	for yl := pos.Y() + int32(height); yl >= pos.Y()+int32(foliageStart); yl-- {

		for xx := x - int32(radius); xx <= x+int32(radius); xx++ {
			for zz := z - int32(radius); zz <= z+int32(radius); zz++ {
				dx := int(math.Abs(float64(xx - x)))
				dz := int(math.Abs(float64(zz - z)))

				if dx != radius || dz != radius || radius <= 0 {
					id := level.GetBlockId(xx, yl, zz)
					if id == 0 || id == byte(block.LEAVES) || id == byte(block.SNOW_LAYER) {
						level.SetBlock(xx, yl, zz, byte(block.LEAVES), 1, false)
					}
				}
			}
		}

		if radius >= 1 && yl == pos.Y()+int32(foliageStart)+1 {
			radius--
		} else if radius < 2 {
			radius++
		}
	}

	for i := int32(0); i < int32(height)-int32(random.NextBoundedInt(3)); i++ {
		id := level.GetBlockId(x, y+i, z)
		if id == 0 || id == byte(block.LEAVES) || id == byte(block.SNOW_LAYER) {
			level.SetBlock(x, y+i, z, byte(block.LOG), 1, false)
		}
	}

	return true
}
