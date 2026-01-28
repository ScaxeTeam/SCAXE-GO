package object

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type SpruceTree struct {
	*Tree
}

func NewSpruceTree() *SpruceTree {
	t := NewBaseTree(block.LOG, block.LEAVES, 1)

	return &SpruceTree{Tree: t}
}

func (st *SpruceTree) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {
	x, y, z := pos.X(), pos.Y(), pos.Z()

	soil := level.GetBlockId(x, y-1, z)
	if soil != byte(block.GRASS) && soil != byte(block.DIRT) && soil != byte(block.PODZOL) {
		return false
	}

	st.TreeHeight = random.NextBoundedInt(4) + 6

	topSize := st.TreeHeight - (1 + random.NextBoundedInt(2))
	lRadius := 2 + random.NextBoundedInt(2)

	st.PlaceTrunk(level, x, y, z, random, st.TreeHeight-random.NextBoundedInt(3))

	radius := random.NextBoundedInt(2)
	maxR := 1
	minR := 0

	for yy := 0; yy <= topSize; yy++ {
		yyy := int(y) + st.TreeHeight - yy

		for xx := int(x) - radius; xx <= int(x)+radius; xx++ {
			xOff := int(math.Abs(float64(xx - int(x))))
			for zz := int(z) - radius; zz <= int(z)+radius; zz++ {
				zOff := int(math.Abs(float64(zz - int(z))))

				if xOff == radius && zOff == radius && radius > 0 {
					continue
				}

				id := level.GetBlockId(int32(xx), int32(yyy), int32(zz))
				if id == 0 || id == byte(block.LEAVES) || id == byte(block.SNOW_LAYER) {
					level.SetBlock(int32(xx), int32(yyy), int32(zz), byte(st.LeafBlock), byte(st.Type), false)
				}
			}
		}

		if radius >= maxR {
			radius = minR
			minR = 1
			maxR++
			if maxR > lRadius {
				maxR = lRadius
			}
		} else {
			radius++
		}
	}
	return true
}
