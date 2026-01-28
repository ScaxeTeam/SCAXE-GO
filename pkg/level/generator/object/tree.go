package object

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Tree struct {
	TrunkBlock int
	LeafBlock  int
	Type       int
	TreeHeight int
	Overrides  map[int]bool
}

func NewBaseTree(trunk, leaf, typeData int) *Tree {
	return &Tree{
		TrunkBlock: trunk,
		LeafBlock:  leaf,
		Type:       typeData,
		Overrides: map[int]bool{
			0:                true,
			block.SAPLING:    true,
			block.LOG:        true,
			block.LEAVES:     true,
			block.SNOW_LAYER: true,
			block.LEAVES2:    true,
			block.WOOD2:      true,
		},
	}
}

func (t *Tree) CanPlaceObject(level populator.ChunkManager, x, y, z int32, random *rand.Random) bool {
	radiusToCheck := 0

	for yy := int32(0); yy < int32(t.TreeHeight)+3; yy++ {

		if yy == 1 || yy == int32(t.TreeHeight) {
			radiusToCheck++
		}

		for xx := int32(-radiusToCheck); xx < int32(radiusToCheck)+1; xx++ {
			for zz := int32(-radiusToCheck); zz < int32(radiusToCheck)+1; zz++ {
				checkY := y + yy
				if checkY >= 0 && checkY < 256 {
					id := level.GetBlockId(x+xx, checkY, z+zz)

					if !t.Overrides[int(id)] {
						return false
					}
				} else {
					return false
				}
			}
		}
	}
	return true
}

func (t *Tree) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {
	x, y, z := pos.X(), pos.Y(), pos.Z()

	soilBlock := level.GetBlockId(x, y-1, z)
	if soilBlock != byte(block.GRASS) && soilBlock != byte(block.DIRT) && soilBlock != byte(block.FARMLAND) {
		return false
	}

	level.SetBlock(x, y-1, z, block.DIRT, 0, false)

	for i3 := int(y) - 3 + t.TreeHeight; i3 <= int(y)+t.TreeHeight; i3++ {

		i4 := i3 - (int(y) + t.TreeHeight)

		j1 := 1 - i4/2

		for k1 := int(x) - j1; k1 <= int(x)+j1; k1++ {

			l1 := k1 - int(x)

			for i2 := int(z) - j1; i2 <= int(z)+j1; i2++ {

				j2 := i2 - int(z)

				absL1 := int(math.Abs(float64(l1)))
				absJ2 := int(math.Abs(float64(j2)))

				if absL1 != j1 || absJ2 != j1 || (random.NextBoundedInt(2) != 0 && i4 != 0) {

					blockId := level.GetBlockId(int32(k1), int32(i3), int32(i2))
					if blockId == 0 || blockId == byte(block.LEAVES) || blockId == byte(block.VINE) {
						level.SetBlock(int32(k1), int32(i3), int32(i2), byte(t.LeafBlock), byte(t.Type), false)
					}
				}
			}
		}
	}

	for j3 := 0; j3 < t.TreeHeight; j3++ {

		blockId := level.GetBlockId(x, y+int32(j3), z)

		if blockId == 0 || blockId == byte(block.LEAVES) || blockId == byte(block.VINE) || blockId == byte(block.LEAVES2) {
			level.SetBlock(x, y+int32(j3), z, byte(t.TrunkBlock), byte(t.Type), false)
		}
	}

	return true
}

func (t *Tree) PlaceTrunk(level populator.ChunkManager, x, y, z int32, random *rand.Random, trunkHeight int) {

	level.SetBlock(x, y-1, z, block.DIRT, 0, false)

	for yy := 0; yy < trunkHeight; yy++ {
		blockId := level.GetBlockId(x, y+int32(yy), z)

		if t.Overrides[int(blockId)] {
			level.SetBlock(x, y+int32(yy), z, byte(t.TrunkBlock), byte(t.Type), false)
		}
	}
}

type OakTree struct {
	*Tree
}

func NewOakTree() *OakTree {
	t := NewBaseTree(block.LOG, block.LEAVES, 0)
	return &OakTree{Tree: t}
}

func (ot *OakTree) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {
	ot.TreeHeight = random.NextBoundedInt(3) + 4

	if !ot.Tree.CanPlaceObject(level, pos.X(), pos.Y(), pos.Z(), random) {
		return false
	}
	return ot.Tree.Generate(level, random, pos)
}

type BirchTree struct {
	*Tree
	SuperBirch bool
}

func NewBirchTree(super bool) *BirchTree {
	t := NewBaseTree(block.LOG, block.LEAVES, 2)
	return &BirchTree{Tree: t, SuperBirch: super}
}

func (bt *BirchTree) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {
	bt.TreeHeight = random.NextBoundedInt(3) + 5
	if bt.SuperBirch {
		bt.TreeHeight += random.NextBoundedInt(7)
	}

	if !bt.Tree.CanPlaceObject(level, pos.X(), pos.Y(), pos.Z(), random) {
		return false
	}
	return bt.Tree.Generate(level, random, pos)
}
