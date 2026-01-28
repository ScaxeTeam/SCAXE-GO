package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type JungleSmallTree struct {
	*Tree
	GenerateVines bool
}

func NewJungleSmallTree() *JungleSmallTree {

	t := NewBaseTree(block.LOG, block.LEAVES, 3)
	return &JungleSmallTree{
		Tree:          t,
		GenerateVines: true,
	}
}

func (jt *JungleSmallTree) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {

	jt.TreeHeight = 4 + random.NextBoundedInt(7)

	x, y, z := pos.X(), pos.Y(), pos.Z()

	groundY := findGround(level, x, y, z)
	if groundY < 1 || groundY > int32(127-jt.TreeHeight) {
		return false
	}

	belowId := level.GetBlockId(x, groundY-1, z)
	if belowId != byte(block.GRASS) && belowId != byte(block.DIRT) {
		return false
	}

	level.SetBlock(x, groundY-1, z, block.DIRT, 0, false)

	for ty := 0; ty < jt.TreeHeight; ty++ {
		level.SetBlock(x, groundY+int32(ty), z, byte(jt.TrunkBlock), byte(jt.Type), false)

		if jt.GenerateVines && ty > 0 {
			jt.tryPlaceVine(level, random, x-1, groundY+int32(ty), z, 8)
			jt.tryPlaceVine(level, random, x+1, groundY+int32(ty), z, 2)
			jt.tryPlaceVine(level, random, x, groundY+int32(ty), z-1, 1)
			jt.tryPlaceVine(level, random, x, groundY+int32(ty), z+1, 4)
		}
	}

	for ly := groundY + int32(jt.TreeHeight) - 3; ly <= groundY+int32(jt.TreeHeight); ly++ {
		yOff := int(ly - (groundY + int32(jt.TreeHeight)))
		radius := 1 - yOff/2

		for lx := x - int32(radius); lx <= x+int32(radius); lx++ {
			for lz := z - int32(radius); lz <= z+int32(radius); lz++ {
				xDist := int32Abs(lx - x)
				zDist := int32Abs(lz - z)

				if xDist == int32(radius) && zDist == int32(radius) && (yOff == 0 || random.NextBoundedInt(2) == 0) {
					continue
				}

				if level.GetBlockId(lx, ly, lz) == 0 {
					level.SetBlock(lx, ly, lz, byte(jt.LeafBlock), byte(jt.Type), false)

					if jt.GenerateVines && yOff < 0 {
						jt.tryGrowVineBelow(level, random, lx, ly, lz)
					}
				}
			}
		}
	}

	return true
}

func (jt *JungleSmallTree) tryPlaceVine(level populator.ChunkManager, random *rand.Random, x, y, z int32, meta byte) {
	if random.NextBoundedInt(3) > 0 && level.GetBlockId(x, y, z) == 0 {
		level.SetBlock(x, y, z, block.VINE, meta, false)
	}
}

func (jt *JungleSmallTree) tryGrowVineBelow(level populator.ChunkManager, random *rand.Random, x, y, z int32) {
	if random.NextBoundedInt(4) == 0 {

		dir := random.NextBoundedInt(4)
		var meta byte
		switch dir {
		case 0:
			meta = 1
		case 1:
			meta = 2
		case 2:
			meta = 4
		case 3:
			meta = 8
		}

		maxLen := random.NextBoundedInt(4) + 1
		for vy := y - 1; vy >= y-int32(maxLen) && vy > 0; vy-- {
			if level.GetBlockId(x, vy, z) == 0 {
				level.SetBlock(x, vy, z, block.VINE, meta, false)
			} else {
				break
			}
		}
	}
}

func findGround(level populator.ChunkManager, x, startY, z int32) int32 {
	for y := startY; y > 0; y-- {
		id := level.GetBlockId(x, y-1, z)
		if id == byte(block.GRASS) || id == byte(block.DIRT) || id == byte(block.SAND) {
			return y
		}

		if id != 0 && id != byte(block.TALL_GRASS) && id != byte(block.LEAVES) {
			return -1
		}
	}
	return -1
}

func int32Abs(v int32) int32 {
	if v < 0 {
		return -v
	}
	return v
}
