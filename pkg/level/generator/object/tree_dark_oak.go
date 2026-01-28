package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type DarkOakTree struct {
	TrunkBlock int
	LeafBlock  int
	Type       int
}

func NewDarkOakTree() *DarkOakTree {
	return &DarkOakTree{
		TrunkBlock: block.WOOD2,
		LeafBlock:  block.LEAVES2,
		Type:       1,
	}
}

func (t *DarkOakTree) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {

	height := random.NextBoundedInt(3) + random.NextBoundedInt(2) + 6

	j := pos.X()
	k := pos.Y()
	l := pos.Z()

	if k < 1 || k+int32(height)+1 >= 256 {
		return false
	}

	blockDown := pos.Down()
	soil := level.GetBlockId(blockDown.X(), blockDown.Y(), blockDown.Z())
	if soil != byte(block.GRASS) && soil != byte(block.DIRT) {
		return false
	}

	if !t.placeTreeOfHeight(level, pos, height) {
		return false
	}

	t.setDirtAt(level, blockDown)
	t.setDirtAt(level, blockDown.East())
	t.setDirtAt(level, blockDown.South())
	t.setDirtAt(level, blockDown.South().East())

	enumfacing := random.NextBoundedInt(4)
	xOffset, zOffset := t.getDirectionOffsets(enumfacing)

	i1 := height - random.NextBoundedInt(4)
	j1 := 2 - random.NextBoundedInt(3)

	k1 := j
	l1 := l
	i2 := k + int32(height) - 1

	for j2 := 0; j2 < height; j2++ {
		if j2 >= i1 && j1 > 0 {
			k1 += xOffset
			l1 += zOffset
			j1--
		}

		k2 := k + int32(j2)
		blockpos1 := world.NewBlockPos(k1, k2, l1)

		mat := level.GetBlockId(blockpos1.X(), blockpos1.Y(), blockpos1.Z())
		if mat == 0 || mat == byte(block.LEAVES) || mat == byte(block.LEAVES2) {

			t.placeLogAt(level, blockpos1)
			t.placeLogAt(level, blockpos1.East())
			t.placeLogAt(level, blockpos1.South())
			t.placeLogAt(level, blockpos1.East().South())
		}
	}

	for i3 := int32(-2); i3 <= 0; i3++ {
		for l3 := int32(-2); l3 <= 0; l3++ {
			k4 := int32(-1)
			t.placeLeafAt(level, k1+i3, i2+k4, l1+l3)
			t.placeLeafAt(level, 1+k1-i3, i2+k4, l1+l3)
			t.placeLeafAt(level, k1+i3, i2+k4, 1+l1-l3)
			t.placeLeafAt(level, 1+k1-i3, i2+k4, 1+l1-l3)

			if (i3 > -2 || l3 > -1) && (i3 != -1 || l3 != -2) {
				k4 = 1
				t.placeLeafAt(level, k1+i3, i2+k4, l1+l3)
				t.placeLeafAt(level, 1+k1-i3, i2+k4, l1+l3)
				t.placeLeafAt(level, k1+i3, i2+k4, 1+l1-l3)
				t.placeLeafAt(level, 1+k1-i3, i2+k4, 1+l1-l3)
			}
		}
	}

	if random.NextBoolean() {
		t.placeLeafAt(level, k1, i2+2, l1)
		t.placeLeafAt(level, k1+1, i2+2, l1)
		t.placeLeafAt(level, k1+1, i2+2, l1+1)
		t.placeLeafAt(level, k1, i2+2, l1+1)
	}

	for j3 := int32(-3); j3 <= 4; j3++ {
		for i4 := int32(-3); i4 <= 4; i4++ {
			if (j3 != -3 || i4 != -3) && (j3 != -3 || i4 != 4) &&
				(j3 != 4 || i4 != -3) && (j3 != 4 || i4 != 4) &&
				(abs32(j3) < 3 || abs32(i4) < 3) {
				t.placeLeafAt(level, k1+j3, i2, l1+i4)
			}
		}
	}

	for k3 := int32(-1); k3 <= 2; k3++ {
		for j4 := int32(-1); j4 <= 2; j4++ {

			if (k3 < 0 || k3 > 1 || j4 < 0 || j4 > 1) && random.NextBoundedInt(3) <= 0 {

				l4 := random.NextBoundedInt(3) + 2

				for i5 := 0; i5 < l4; i5++ {
					t.placeLogAt(level, world.NewBlockPos(j+k3, i2-int32(i5)-1, l+j4))
				}

				for j5 := int32(-1); j5 <= 1; j5++ {
					for l2 := int32(-1); l2 <= 1; l2++ {
						t.placeLeafAt(level, k1+k3+j5, i2, l1+j4+l2)
					}
				}

				for k5 := int32(-2); k5 <= 2; k5++ {
					for l5 := int32(-2); l5 <= 2; l5++ {
						if abs32(k5) != 2 || abs32(l5) != 2 {
							t.placeLeafAt(level, k1+k3+k5, i2-1, l1+j4+l5)
						}
					}
				}
			}
		}
	}

	return true
}

func (t *DarkOakTree) placeTreeOfHeight(level populator.ChunkManager, pos world.BlockPos, height int) bool {
	i := pos.X()
	j := pos.Y()
	k := pos.Z()

	for l := 0; l <= height+1; l++ {
		i1 := 1
		if l == 0 {
			i1 = 0
		}
		if l >= height-1 {
			i1 = 2
		}

		for j1 := -int32(i1); j1 <= int32(i1); j1++ {
			for k1 := -int32(i1); k1 <= int32(i1); k1++ {
				blockId := level.GetBlockId(i+j1, j+int32(l), k+k1)
				if !t.canGrowInto(blockId) {
					return false
				}
			}
		}
	}
	return true
}

func (t *DarkOakTree) canGrowInto(blockId byte) bool {
	switch blockId {
	case 0,
		block.LEAVES, block.LEAVES2,
		block.LOG, block.WOOD2,
		block.SAPLING,
		block.VINE:
		return true
	}
	return false
}

func (t *DarkOakTree) setDirtAt(level populator.ChunkManager, pos world.BlockPos) {
	blockId := level.GetBlockId(pos.X(), pos.Y(), pos.Z())
	if blockId == byte(block.GRASS) || blockId == byte(block.FARMLAND) {
		level.SetBlock(pos.X(), pos.Y(), pos.Z(), byte(block.DIRT), 0, false)
	}
}

func (t *DarkOakTree) placeLogAt(level populator.ChunkManager, pos world.BlockPos) {
	blockId := level.GetBlockId(pos.X(), pos.Y(), pos.Z())
	if t.canGrowInto(blockId) {
		level.SetBlock(pos.X(), pos.Y(), pos.Z(), byte(t.TrunkBlock), byte(t.Type), false)
	}
}

func (t *DarkOakTree) placeLeafAt(level populator.ChunkManager, x, y, z int32) {
	blockId := level.GetBlockId(x, y, z)
	if blockId == 0 {
		level.SetBlock(x, y, z, byte(t.LeafBlock), byte(t.Type), false)
	}
}

func (t *DarkOakTree) getDirectionOffsets(dir int) (int32, int32) {

	switch dir {
	case 0:
		return 0, -1
	case 1:
		return 0, 1
	case 2:
		return -1, 0
	case 3:
		return 1, 0
	}
	return 0, 0
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
