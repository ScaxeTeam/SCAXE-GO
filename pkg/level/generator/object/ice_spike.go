package object

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type IceSpike struct{}

func NewIceSpike() *IceSpike {
	return &IceSpike{}
}

func (is *IceSpike) Generate(level populator.ChunkManager, random *rand.Random, pos world.BlockPos) bool {
	x, y, z := pos.X(), pos.Y(), pos.Z()

	for level.GetBlockId(x, y, z) == 0 && y > 2 {
		y--
	}

	if level.GetBlockId(x, y, z) != byte(block.SNOW_BLOCK) {
		return false
	}

	y += int32(random.NextBoundedInt(4))

	height := random.NextBoundedInt(4) + 7
	baseRadius := height/4 + random.NextBoundedInt(2)

	if baseRadius > 1 && random.NextBoundedInt(60) == 0 {
		y += int32(10 + random.NextBoundedInt(30))
	}

	for k := 0; k < height; k++ {

		f := (1.0 - float64(k)/float64(height)) * float64(baseRadius)
		radius := int(math.Ceil(f))

		for dx := -radius; dx <= radius; dx++ {
			f1 := float64(absInt(dx)) - 0.25

			for dz := -radius; dz <= radius; dz++ {
				f2 := float64(absInt(dz)) - 0.25

				if (dx == 0 && dz == 0) ||
					(f1*f1+f2*f2 <= f*f &&
						(dx != -radius && dx != radius && dz != -radius && dz != radius ||
							random.NextFloat() <= 0.75)) {

					bx, by, bz := x+int32(dx), y+int32(k), z+int32(dz)
					blockId := level.GetBlockId(bx, by, bz)

					if blockId == 0 || blockId == byte(block.DIRT) ||
						blockId == byte(block.SNOW_BLOCK) || blockId == byte(block.ICE) {
						level.SetBlock(bx, by, bz, byte(block.PACKED_ICE), 0, false)
					}

					if k != 0 && radius > 1 {
						bx, by, bz = x+int32(dx), y-int32(k), z+int32(dz)
						blockId = level.GetBlockId(bx, by, bz)

						if blockId == 0 || blockId == byte(block.DIRT) ||
							blockId == byte(block.SNOW_BLOCK) || blockId == byte(block.ICE) {
							level.SetBlock(bx, by, bz, byte(block.PACKED_ICE), 0, false)
						}
					}
				}
			}
		}
	}

	rootRadius := baseRadius - 1
	if rootRadius < 0 {
		rootRadius = 0
	} else if rootRadius > 1 {
		rootRadius = 1
	}

	for dx := -rootRadius; dx <= rootRadius; dx++ {
		for dz := -rootRadius; dz <= rootRadius; dz++ {
			rootY := y - 1
			depth := 50

			if absInt(dx) == 1 && absInt(dz) == 1 {
				depth = random.NextBoundedInt(5)
			}

			for rootY > 50 {
				blockId := level.GetBlockId(x+int32(dx), rootY, z+int32(dz))

				if blockId != 0 && blockId != byte(block.DIRT) &&
					blockId != byte(block.SNOW_BLOCK) && blockId != byte(block.ICE) &&
					blockId != byte(block.PACKED_ICE) {
					break
				}

				level.SetBlock(x+int32(dx), rootY, z+int32(dz), byte(block.PACKED_ICE), 0, false)
				rootY--
				depth--

				if depth <= 0 {
					rootY -= int32(random.NextBoundedInt(5) + 1)
					depth = random.NextBoundedInt(5)
				}
			}
		}
	}

	return true
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
