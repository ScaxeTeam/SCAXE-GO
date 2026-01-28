package object

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type BigOakTree struct {
	basePos     world.BlockPos
	heightLimit int
	height      int

	heightAttenuation float64
	branchSlope       float64
	scaleWidth        float64
	leafDensity       float64

	trunkSize         int
	heightLimitLimit  int
	leafDistanceLimit int

	foliageCoords []foliageCoordinates
	level         populator.ChunkManager
	rand          *rand.Random
}

type foliageCoordinates struct {
	world.BlockPos
	branchBase int
}

func NewBigOakTree() *BigOakTree {
	return &BigOakTree{
		heightAttenuation: 0.618,
		branchSlope:       0.381,
		scaleWidth:        1.0,
		leafDensity:       1.0,
		trunkSize:         1,
		heightLimitLimit:  12,
		leafDistanceLimit: 4,
	}
}

func (t *BigOakTree) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {
	t.level = level
	t.basePos = pos

	t.rand = rand.NewRandom(int64(random.NextInt()))

	t.heightLimit = 5 + t.rand.NextBoundedInt(t.heightLimitLimit)

	if !t.validTreeLocation() {
		return false
	}

	t.generateLeafNodeList()
	t.generateLeaves()
	t.generateTrunk()
	t.generateLeafNodeBases()
	return true
}

func (t *BigOakTree) validTreeLocation() bool {
	down := t.basePos.Down()
	id := t.level.GetBlockId(down.X(), down.Y(), down.Z())

	if id != byte(block.DIRT) && id != byte(block.GRASS) && id != byte(block.FARMLAND) {
		return false
	}

	limit := t.checkBlockLine(t.basePos, t.basePos.Up(int32(t.heightLimit-1)))
	if limit == -1 {
		return true
	} else if limit < 6 {
		return false
	} else {
		t.heightLimit = limit
		return true
	}
}

func (t *BigOakTree) generateLeafNodeList() {
	t.height = int(float64(t.heightLimit) * t.heightAttenuation)
	if t.height >= t.heightLimit {
		t.height = t.heightLimit - 1
	}

	i := int(1.382 + math.Pow(t.leafDensity*float64(t.heightLimit)/13.0, 2.0))
	if i < 1 {
		i = 1
	}

	j := t.basePos.Y() + int32(t.height)
	k := t.heightLimit - t.leafDistanceLimit

	t.foliageCoords = []foliageCoordinates{}
	t.foliageCoords = append(t.foliageCoords, foliageCoordinates{
		BlockPos:   t.basePos.Up(int32(k)),
		branchBase: int(j),
	})

	for ; k >= 0; k-- {
		f := t.layerSize(k)
		if f >= 0.0 {
			for l := 0; l < i; l++ {
				d0 := t.scaleWidth * float64(f) * (float64(t.rand.NextFloat()) + 0.328)
				d1 := float64(t.rand.NextFloat()*2.0) * math.Pi
				d2 := d0*math.Sin(d1) + 0.5
				d3 := d0*math.Cos(d1) + 0.5

				xVal := float64(t.basePos.X()) + d2
				yVal := float64(t.basePos.Y()) + float64(k-1)
				zVal := float64(t.basePos.Z()) + d3

				blockpos := world.NewBlockPos(
					int32(math.Floor(xVal)),
					int32(math.Floor(yVal)),
					int32(math.Floor(zVal)),
				)
				blockpos1 := blockpos.Up(int32(t.leafDistanceLimit))

				if t.checkBlockLine(blockpos, blockpos1) == -1 {
					i1 := t.basePos.X() - blockpos.X()
					j1 := t.basePos.Z() - blockpos.Z()
					d4 := float64(blockpos.Y()) - math.Sqrt(float64(i1*i1+j1*j1))*t.branchSlope

					k1 := int(d4)
					if d4 > float64(j) {
						k1 = int(j)
					}

					blockpos2 := world.NewBlockPos(t.basePos.X(), int32(k1), t.basePos.Z())

					if t.checkBlockLine(blockpos2, blockpos) == -1 {
						t.foliageCoords = append(t.foliageCoords, foliageCoordinates{
							BlockPos:   blockpos,
							branchBase: int(blockpos2.Y()),
						})
					}
				}
			}
		}
	}
}

func (t *BigOakTree) layerSize(y int) float32 {
	if float32(y) < float32(t.heightLimit)*0.3 {
		return -1.0
	}
	f := float32(t.heightLimit) / 2.0
	f1 := f - float32(y)
	f2 := float32(math.Sqrt(float64(f*f - f1*f1)))
	if f1 == 0.0 {
		f2 = f
	} else if math.Abs(float64(f1)) >= float64(f) {
		return 0.0
	}
	return f2 * 0.5
}

func (t *BigOakTree) generateLeaves() {
	for _, coord := range t.foliageCoords {
		t.generateLeafNode(coord.BlockPos)
	}
}

func (t *BigOakTree) generateLeafNode(pos world.BlockPos) {
	for i := 0; i < t.leafDistanceLimit; i++ {
		t.crosSection(pos.Up(int32(i)), t.leafSize(i))
	}
}

func (t *BigOakTree) leafSize(y int) float32 {
	if y >= 0 && y < t.leafDistanceLimit {
		if y != 0 && y != t.leafDistanceLimit-1 {
			return 3.0
		}
		return 2.0
	}
	return -1.0
}

func (t *BigOakTree) crosSection(pos world.BlockPos, radius float32) {
	i := int(float64(radius) + 0.618)

	for j := -i; j <= i; j++ {
		for k := -i; k <= i; k++ {
			if math.Pow(math.Abs(float64(j))+0.5, 2.0)+math.Pow(math.Abs(float64(k))+0.5, 2.0) <= float64(radius*radius) {
				blockpos := pos.Add(int32(j), 0, int32(k))
				mat := t.level.GetBlockId(blockpos.X(), blockpos.Y(), blockpos.Z())
				if mat == 0 || mat == byte(block.LEAVES) {

					t.level.SetBlock(blockpos.X(), blockpos.Y(), blockpos.Z(), byte(block.LEAVES), 0, false)
				}
			}
		}
	}
}

func (t *BigOakTree) generateTrunk() {
	blockpos := t.basePos
	blockpos1 := t.basePos.Up(int32(t.height))

	t.limb(blockpos, blockpos1)

	if t.trunkSize == 2 {
		t.limb(blockpos.East(), blockpos1.East())
		t.limb(blockpos.East().South(), blockpos1.East().South())
		t.limb(blockpos.South(), blockpos1.South())
	}
}

func (t *BigOakTree) generateLeafNodeBases() {
	for _, coord := range t.foliageCoords {
		i := coord.branchBase
		blockpos := world.NewBlockPos(t.basePos.X(), int32(i), t.basePos.Z())

		if (blockpos != coord.BlockPos) && t.leafNodeNeedsBase(i-int(t.basePos.Y())) {
			t.limb(blockpos, coord.BlockPos)
		}
	}
}

func (t *BigOakTree) leafNodeNeedsBase(p_76493_1_ int) bool {
	return float64(p_76493_1_) >= float64(t.heightLimit)*0.2
}

func (t *BigOakTree) limb(start, end world.BlockPos) {

	blockpos := end.Add(-start.X(), -start.Y(), -start.Z())
	i := t.getGreatestDistance(blockpos)

	if i == 0 {
		return
	}

	f := float32(blockpos.X()) / float32(i)
	f1 := float32(blockpos.Y()) / float32(i)
	f2 := float32(blockpos.Z()) / float32(i)

	for j := 0; j <= i; j++ {

		xVal := float64(start.X()) + 0.5 + float64(float32(j)*f)
		yVal := float64(start.Y()) + 0.5 + float64(float32(j)*f1)
		zVal := float64(start.Z()) + 0.5 + float64(float32(j)*f2)

		blockpos1 := world.NewBlockPos(
			int32(math.Floor(xVal)),
			int32(math.Floor(yVal)),
			int32(math.Floor(zVal)),
		)

		axis := t.getLogAxis(start, blockpos1)

		meta := 0
		switch axis {
		case 0:

		case 1:
			meta |= 4
		case 2:
			meta |= 8
		}

		t.level.SetBlock(blockpos1.X(), blockpos1.Y(), blockpos1.Z(), byte(block.LOG), byte(meta), false)
	}
}

func (t *BigOakTree) getLogAxis(start, end world.BlockPos) int {
	axis := 0
	i := int(math.Abs(float64(end.X() - start.X())))
	j := int(math.Abs(float64(end.Z() - start.Z())))
	k := int(math.Max(float64(i), float64(j)))

	if k > 0 {
		if i == k {
			axis = 1
		} else if j == k {
			axis = 2
		}
	}
	return axis
}

func (t *BigOakTree) getGreatestDistance(pos world.BlockPos) int {
	i := int(math.Abs(float64(pos.X())))
	j := int(math.Abs(float64(pos.Y())))
	k := int(math.Abs(float64(pos.Z())))

	if k > i && k > j {
		return k
	}
	if j > i {
		return j
	}
	return i
}

func (t *BigOakTree) checkBlockLine(start, end world.BlockPos) int {
	blockpos := end.Add(-start.X(), -start.Y(), -start.Z())
	i := t.getGreatestDistance(blockpos)

	if i == 0 {
		return -1
	}

	f := float32(blockpos.X()) / float32(i)
	f1 := float32(blockpos.Y()) / float32(i)
	f2 := float32(blockpos.Z()) / float32(i)

	for j := 0; j <= i; j++ {
		xVal := float64(start.X()) + 0.5 + float64(float32(j)*f)
		yVal := float64(start.Y()) + 0.5 + float64(float32(j)*f1)
		zVal := float64(start.Z()) + 0.5 + float64(float32(j)*f2)

		blockpos1 := world.NewBlockPos(
			int32(math.Floor(xVal)),
			int32(math.Floor(yVal)),
			int32(math.Floor(zVal)),
		)

		id := t.level.GetBlockId(blockpos1.X(), blockpos1.Y(), blockpos1.Z())

		if id != 0 && id != byte(block.LEAVES) && id != byte(block.SAPLING) && id != byte(block.VINE) && id != byte(block.LOG) {
			return j
		}
	}
	return -1
}
