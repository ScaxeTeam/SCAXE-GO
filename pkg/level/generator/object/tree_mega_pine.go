package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type MegaPineTree struct {
	HasPodzol bool
}

func NewMegaPineTree(hasPodzol bool) *MegaPineTree {
	return &MegaPineTree{HasPodzol: hasPodzol}
}

func (mpt *MegaPineTree) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {
	x, y, z := pos.X(), pos.Y(), pos.Z()

	groundY := mpt.findGround(level, x, y, z)
	if groundY < 1 || groundY > 110 {
		return false
	}

	for dx := int32(0); dx <= 1; dx++ {
		for dz := int32(0); dz <= 1; dz++ {
			belowId := level.GetBlockId(x+dx, groundY-1, z+dz)
			if belowId != byte(block.GRASS) && belowId != byte(block.DIRT) && belowId != byte(block.PODZOL) {
				return false
			}
		}
	}

	height := 13 + random.NextBoundedInt(15)

	if mpt.HasPodzol {
		mpt.placePodzol(level, random, x, groundY, z)
	}

	for ty := 0; ty < height; ty++ {
		for dx := int32(0); dx <= 1; dx++ {
			for dz := int32(0); dz <= 1; dz++ {

				level.SetBlock(x+dx, groundY+int32(ty), z+dz, block.LOG, 1, false)
			}
		}
	}

	leafStart := groundY + int32(height) - int32(random.NextBoundedInt(3)) - 3
	radius := 0

	for ly := groundY + int32(height); ly >= leafStart; ly-- {
		yFromTop := int(groundY+int32(height)) - int(ly)

		if yFromTop < 2 {
			radius = 0
		} else if yFromTop < 6 {
			radius = 1 + yFromTop/2
		} else {
			radius = 2 + random.NextBoundedInt(2)
		}

		for lx := x - int32(radius); lx <= x+1+int32(radius); lx++ {
			for lz := z - int32(radius); lz <= z+1+int32(radius); lz++ {
				xDist := lx - x
				zDist := lz - z

				if xDist < 0 {
					xDist = -xDist
				}
				if zDist < 0 {
					zDist = -zDist
				}
				if int(xDist)+int(zDist) > radius+2 {
					continue
				}

				id := level.GetBlockId(lx, ly, lz)
				if id == 0 || id == byte(block.LEAVES) {

					level.SetBlock(lx, ly, lz, block.LEAVES, 1, false)
				}
			}
		}
	}

	return true
}

func (mpt *MegaPineTree) placePodzol(level populator.ChunkManager, random *rand.Random, x, y, z int32) {
	radius := 2 + random.NextBoundedInt(2)
	for dx := -int32(radius); dx <= 1+int32(radius); dx++ {
		for dz := -int32(radius); dz <= 1+int32(radius); dz++ {

			dist := dx*dx + dz*dz
			if dist > int32((radius+1)*(radius+1)) {
				continue
			}

			checkY := y - 1
			id := level.GetBlockId(x+dx, checkY, z+dz)
			if id == byte(block.GRASS) || id == byte(block.DIRT) {

				level.SetBlock(x+dx, checkY, z+dz, block.DIRT, 2, false)
			}
		}
	}
}

func (mpt *MegaPineTree) findGround(level populator.ChunkManager, x, startY, z int32) int32 {
	for y := startY; y > 0; y-- {
		id := level.GetBlockId(x, y-1, z)
		if id == byte(block.GRASS) || id == byte(block.DIRT) || id == byte(block.PODZOL) {
			return y
		}
		if id != 0 && id != byte(block.TALL_GRASS) && id != byte(block.LEAVES) {
			return -1
		}
	}
	return -1
}
