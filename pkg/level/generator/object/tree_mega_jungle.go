package object

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type MegaJungleTree struct {
	baseHeight int
	woodMeta   int
	leafMeta   int
}

func NewMegaJungleTree() *MegaJungleTree {

	return &MegaJungleTree{
		baseHeight: 10,
		woodMeta:   3,
		leafMeta:   3,
	}
}

func (t *MegaJungleTree) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {
	height := random.NextBoundedInt(20) + t.baseHeight

	if !t.ensureGrowable(level, random, pos, height) {
		return false
	}

	t.createCrown(level, pos.Up(int32(height)), 2)

	for j := int32(height) - 2 - int32(random.NextBoundedInt(4)); j > int32(height)/2; j -= 2 + int32(random.NextBoundedInt(4)) {
		f := random.NextFloat() * math.Pi * 2.0

		k := pos.X() + int32(0.5+math.Cos(f)*4.0)
		l := pos.Z() + int32(0.5+math.Sin(f)*4.0)

		for i1 := 0; i1 < 5; i1++ {
			k = pos.X() + int32(1.5+math.Cos(f)*float64(i1))
			l = pos.Z() + int32(1.5+math.Sin(f)*float64(i1))

			level.SetBlock(k, pos.Y()+j-3+int32(i1/2), l, byte(block.LOG), byte(t.woodMeta), false)
		}

		j2 := 1 + random.NextBoundedInt(2)
		j1 := j

		for k1 := j - int32(j2); k1 <= j1; k1++ {
			l1 := k1 - j1
			t.growLeavesLayer(level, world.NewBlockPos(k, pos.Y()+k1, l), 1-int(l1))
		}
	}

	for i2 := int32(0); i2 < int32(height); i2++ {
		blockpos := pos.Up(i2)

		if t.canGrowInto(level, blockpos) {
			level.SetBlock(blockpos.X(), blockpos.Y(), blockpos.Z(), byte(block.LOG), byte(t.woodMeta), false)

			if i2 > 0 {
				t.placeVine(level, random, blockpos.West(), 8)
				t.placeVine(level, random, blockpos.North(), 1)
			}
		}

		if i2 < int32(height)-1 {
			blockpos1 := blockpos.East()
			if t.canGrowInto(level, blockpos1) {
				level.SetBlock(blockpos1.X(), blockpos1.Y(), blockpos1.Z(), byte(block.LOG), byte(t.woodMeta), false)
				if i2 > 0 {
					t.placeVine(level, random, blockpos1.East(), 2)
					t.placeVine(level, random, blockpos1.North(), 1)
				}
			}

			blockpos2 := blockpos.South().East()
			if t.canGrowInto(level, blockpos2) {
				level.SetBlock(blockpos2.X(), blockpos2.Y(), blockpos2.Z(), byte(block.LOG), byte(t.woodMeta), false)
				if i2 > 0 {
					t.placeVine(level, random, blockpos2.East(), 2)
					t.placeVine(level, random, blockpos2.South(), 4)
				}
			}

			blockpos3 := blockpos.South()
			if t.canGrowInto(level, blockpos3) {
				level.SetBlock(blockpos3.X(), blockpos3.Y(), blockpos3.Z(), byte(block.LOG), byte(t.woodMeta), false)
				if i2 > 0 {
					t.placeVine(level, random, blockpos3.West(), 8)
					t.placeVine(level, random, blockpos3.South(), 4)
				}
			}
		}
	}

	return true
}

func (t *MegaJungleTree) ensureGrowable(level populator.ChunkManager, random *rand.Random, pos world.BlockPos, height int) bool {

	if pos.Y() < 1 || pos.Y()+int32(height)+1 > 256 {
		return false
	}

	soil1 := level.GetBlockId(pos.X(), pos.Y()-1, pos.Z())
	soil2 := level.GetBlockId(pos.X()+1, pos.Y()-1, pos.Z())
	soil3 := level.GetBlockId(pos.X(), pos.Y()-1, pos.Z()+1)
	soil4 := level.GetBlockId(pos.X()+1, pos.Y()-1, pos.Z()+1)

	validSoil := func(id byte) bool { return id == byte(block.GRASS) || id == byte(block.DIRT) }

	if validSoil(soil1) && validSoil(soil2) && validSoil(soil3) && validSoil(soil4) {

		level.SetBlock(pos.X(), pos.Y()-1, pos.Z(), byte(block.DIRT), 0, false)
		level.SetBlock(pos.X()+1, pos.Y()-1, pos.Z(), byte(block.DIRT), 0, false)
		level.SetBlock(pos.X(), pos.Y()-1, pos.Z()+1, byte(block.DIRT), 0, false)
		level.SetBlock(pos.X()+1, pos.Y()-1, pos.Z()+1, byte(block.DIRT), 0, false)
		return true
	}

	return false
}

func (t *MegaJungleTree) createCrown(level populator.ChunkManager, pos world.BlockPos, width int) {
	for j := int32(-2); j <= 0; j++ {
		t.growLeavesLayerStrict(level, pos.Up(j), width+1-int(-j))
	}
}

func (t *MegaJungleTree) growLeavesLayer(level populator.ChunkManager, pos world.BlockPos, width int) {
	radius := int32(width)
	for x := -radius; x <= radius; x++ {
		for z := -radius; z <= radius; z++ {
			if math.Abs(float64(x)) != float64(radius) || math.Abs(float64(z)) != float64(radius) {
				t.placeLeafAt(level, pos.Add(x, 0, z))
			}
		}
	}
}

func (t *MegaJungleTree) growLeavesLayerStrict(level populator.ChunkManager, pos world.BlockPos, width int) {
	radius := int32(width)
	for x := -radius; x <= radius; x++ {
		for z := -radius; z <= radius; z++ {
			dSq := x*x + z*z
			if dSq <= radius*radius {
				t.placeLeafAt(level, pos.Add(x, 0, z))
			}
		}
	}
}

func (t *MegaJungleTree) placeLeafAt(level populator.ChunkManager, pos world.BlockPos) {
	id := level.GetBlockId(pos.X(), pos.Y(), pos.Z())
	if id == 0 || id == byte(block.LEAVES) {
		level.SetBlock(pos.X(), pos.Y(), pos.Z(), byte(block.LEAVES), byte(t.leafMeta), false)
	}
}

func (t *MegaJungleTree) canGrowInto(level populator.ChunkManager, pos world.BlockPos) bool {
	id := level.GetBlockId(pos.X(), pos.Y(), pos.Z())

	return id == 0 || id == byte(block.LEAVES) || id == byte(block.SAPLING) || id == byte(block.VINE)
}

func (t *MegaJungleTree) placeVine(level populator.ChunkManager, random *rand.Random, pos world.BlockPos, meta int) {
	if random.NextBoundedInt(3) > 0 && level.GetBlockId(pos.X(), pos.Y(), pos.Z()) == 0 {
		level.SetBlock(pos.X(), pos.Y(), pos.Z(), byte(block.VINE), byte(meta), false)
	}
}
