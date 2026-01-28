package object

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type SavannaTree struct {
	*Tree
}

func NewSavannaTree() *SavannaTree {

	t := NewBaseTree(block.WOOD2, block.LEAVES2, 0)
	return &SavannaTree{Tree: t}
}

func (t *SavannaTree) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {

	height := random.NextBoundedInt(3) + random.NextBoundedInt(3) + 5
	flag := true

	x, y, z := pos.X(), pos.Y(), pos.Z()

	if y >= 1 && y+int32(height)+1 <= 256 {
		for j := y; j <= y+1+int32(height); j++ {
			k := int32(1)
			if j == y {
				k = 0
			}
			if j >= y+1+int32(height)-2 {
				k = 2
			}

			for l := x - k; l <= x+k && flag; l++ {
				for i1 := z - k; i1 <= z+k && flag; i1++ {
					if j >= 0 && j < 256 {
						id := level.GetBlockId(l, j, i1)

						if t.Overrides[int(id)] == false {

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

		soil := level.GetBlockId(x, y-1, z)
		if (soil == byte(block.GRASS) || soil == byte(block.DIRT)) && y < 256-int32(height)-1 {
			level.SetBlock(x, y-1, z, byte(block.DIRT), 0, false)

			dir := random.NextBoundedInt(4)

			xOffset := int32(0)
			zOffset := int32(0)
			switch dir {
			case 0:
				zOffset = -1
			case 1:
				zOffset = 1
			case 2:
				xOffset = -1
			case 3:
				xOffset = 1
			}

			k2 := int32(height) - int32(random.NextBoundedInt(4)) - 1
			l2 := 3 - random.NextBoundedInt(3)
			i3, k1 := x, z
			topY := int32(0)

			for l1 := int32(0); l1 < int32(height); l1++ {
				i2 := y + l1
				if l1 >= k2 && l2 > 0 {
					i3 += xOffset
					k1 += zOffset
					l2--
				}

				level.SetBlock(i3, i2, k1, byte(block.WOOD2), 0, false)
				topY = i2
			}

			pos2 := world.NewBlockPos(i3, topY, k1)
			t.generateCanopy(level, pos2)

			pos2 = pos2.Up(1)

			dir1 := random.NextBoundedInt(4)
			if dir1 != dir {

				xOffset1 := int32(0)
				zOffset1 := int32(0)
				switch dir1 {
				case 0:
					zOffset1 = -1
				case 1:
					zOffset1 = 1
				case 2:
					xOffset1 = -1
				case 3:
					xOffset1 = 1
				}

				l3 := k2 - int32(random.NextBoundedInt(2)) - 1
				k4 := 1 + random.NextBoundedInt(3)
				topYBranch := int32(0)
				i3Branch, k1Branch := x, z

				for l4 := l3; l4 < int32(height) && k4 > 0; k4-- {
					if l4 >= 1 {
						j2 := y + l4
						i3Branch += xOffset1
						k1Branch += zOffset1

						level.SetBlock(i3Branch, j2, k1Branch, byte(block.WOOD2), 0, false)
						topYBranch = j2
					}
					l4++
				}

				if topYBranch > 0 {
					pos3 := world.NewBlockPos(i3Branch, topYBranch, k1Branch)

					t.generateSmallCanopy(level, pos3)
				}
			}

			return true
		}
	}
	return false
}

func (t *SavannaTree) generateCanopy(level populator.ChunkManager, pos world.BlockPos) {

	for j3 := int32(-3); j3 <= 3; j3++ {
		for i4 := int32(-3); i4 <= 3; i4++ {
			if math.Abs(float64(j3)) != 3 || math.Abs(float64(i4)) != 3 {
				t.placeLeafAt(level, pos.Add(j3, 0, i4))
			}
		}
	}

	pos = pos.Up(1)
	for k3 := int32(-1); k3 <= 1; k3++ {
		for j4 := int32(-1); j4 <= 1; j4++ {
			t.placeLeafAt(level, pos.Add(k3, 0, j4))
		}
	}

	t.placeLeafAt(level, pos.East().East())
	t.placeLeafAt(level, pos.West().West())
	t.placeLeafAt(level, pos.South().South())
	t.placeLeafAt(level, pos.North().North())
}

func (t *SavannaTree) generateSmallCanopy(level populator.ChunkManager, pos world.BlockPos) {

	for i5 := int32(-2); i5 <= 2; i5++ {
		for k5 := int32(-2); k5 <= 2; k5++ {
			if math.Abs(float64(i5)) != 2 || math.Abs(float64(k5)) != 2 {
				t.placeLeafAt(level, pos.Add(i5, 0, k5))
			}
		}
	}

	pos = pos.Up(1)
	for j5 := int32(-1); j5 <= 1; j5++ {
		for l5 := int32(-1); l5 <= 1; l5++ {
			t.placeLeafAt(level, pos.Add(j5, 0, l5))
		}
	}
}

func (t *SavannaTree) placeLeafAt(level populator.ChunkManager, pos world.BlockPos) {
	id := level.GetBlockId(pos.X(), pos.Y(), pos.Z())
	if id == 0 || id == byte(block.LEAVES) || id == byte(block.LEAVES2) {
		level.SetBlock(pos.X(), pos.Y(), pos.Z(), byte(block.LEAVES2), 0, false)
	}
}
