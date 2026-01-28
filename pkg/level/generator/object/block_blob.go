package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type BlockBlob struct {
	Block int
	Size  int
}

func NewBlockBlob(blockId int, startRadius int) *BlockBlob {
	return &BlockBlob{
		Block: blockId,
		Size:  startRadius,
	}
}

func NewMossyCobblestoneBlob() *BlockBlob {
	return NewBlockBlob(block.MOSS_STONE, 0)
}

func (bb *BlockBlob) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {
	x, y, z := pos.X(), pos.Y(), pos.Z()

	for y > 3 {
		id := level.GetBlockId(x, y-1, z)
		if id != 0 && id != byte(block.LEAVES) && id != byte(block.LEAVES2) {
			break
		}
		y--
	}

	if y <= 3 {
		return false
	}

	for i := 0; i < 3; i++ {
		xRadius := bb.Size + random.NextBoundedInt(2)
		yRadius := bb.Size + random.NextBoundedInt(2)
		zRadius := bb.Size + random.NextBoundedInt(2)

		f := float64(xRadius+yRadius+zRadius)*0.333 + 0.5
		fSq := f * f

		for bx := x - int32(xRadius); bx <= x+int32(xRadius); bx++ {
			xDist := float64(bx-x) + 0.5
			for bz := z - int32(zRadius); bz <= z+int32(zRadius); bz++ {
				zDist := float64(bz-z) + 0.5
				for by := y - int32(yRadius); by <= y+int32(yRadius); by++ {
					yDist := float64(by-y) + 0.5

					distSq := xDist*xDist + yDist*yDist + zDist*zDist
					if distSq <= fSq {
						level.SetBlock(bx, by, bz, byte(bb.Block), 0, false)
					}
				}
			}
		}

		x += int32(-(1 + random.NextBoundedInt(2)) + random.NextBoundedInt(1+random.NextBoundedInt(2)))
		y -= int32(random.NextBoundedInt(2))
		z += int32(-(1 + random.NextBoundedInt(2)) + random.NextBoundedInt(1+random.NextBoundedInt(2)))
	}

	return true
}
