package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type IcePath struct {
	BasePathWidth int
}

func NewIcePath(baseWidth int) *IcePath {
	return &IcePath{BasePathWidth: baseWidth}
}

func (ip *IcePath) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {
	x, y, z := pos.X(), pos.Y(), pos.Z()

	for level.GetBlockId(x, y, z) == 0 && y > 2 {
		y--
	}

	if level.GetBlockId(x, y, z) != byte(block.SNOW_BLOCK) {
		return false
	}

	radius := 2
	if ip.BasePathWidth > 2 {
		radius += random.NextBoundedInt(ip.BasePathWidth - 2)
	}

	for dx := -int32(radius); dx <= int32(radius); dx++ {
		for dz := -int32(radius); dz <= int32(radius); dz++ {

			distSq := dx*dx + dz*dz
			if distSq > int32(radius*radius) {
				continue
			}

			for dy := int32(-1); dy <= 1; dy++ {
				bx, by, bz := x+dx, y+dy, z+dz
				blockId := level.GetBlockId(bx, by, bz)

				if blockId == byte(block.DIRT) ||
					blockId == byte(block.SNOW_BLOCK) ||
					blockId == byte(block.ICE) {
					level.SetBlock(bx, by, bz, byte(block.PACKED_ICE), 0, false)
				}
			}
		}
	}

	return true
}
