package object

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type JungleBush struct {
	LogBlockID  byte
	LogMeta     byte
	LeafBlockID byte
	LeafMeta    byte
}

func NewJungleBush(logID, logMeta, leafID, leafMeta byte) *JungleBush {
	return &JungleBush{
		LogBlockID:  logID,
		LogMeta:     logMeta,
		LeafBlockID: leafID,
		LeafMeta:    leafMeta,
	}
}

func (jb *JungleBush) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {

	x, y, z := pos.X(), pos.Y(), pos.Z()

	for y > 0 && level.GetBlockId(x, y, z) == 0 {
		y--
	}
	y++

	soil := level.GetBlockId(x, y-1, z)
	if soil == byte(block.GRASS) || soil == byte(block.DIRT) {
		level.SetBlock(x, y-1, z, byte(block.DIRT), 0, false)

		level.SetBlock(x, y, z, jb.LogBlockID, jb.LogMeta, false)

		for l := y; l <= y+2; l++ {
			yOff := l - y
			radius := 2 - yOff

			for xx := x - int32(radius); xx <= x+int32(radius); xx++ {
				for zz := z - int32(radius); zz <= z+int32(radius); zz++ {
					if math.Abs(float64(xx-x)) != float64(radius) || math.Abs(float64(zz-z)) != float64(radius) || random.NextBoundedInt(2) != 0 {
						id := level.GetBlockId(xx, l, zz)
						if id == 0 || id == byte(block.LEAVES) {
							level.SetBlock(xx, l, zz, jb.LeafBlockID, jb.LeafMeta, false)
						}
					}
				}
			}
		}
		return true
	}
	return false
}
