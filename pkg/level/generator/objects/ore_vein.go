package objects

import (
	"math"
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/world"
)

type OreVein struct {
	OreType OreType
}

func NewOreVein(oreType OreType) *OreVein {
	return &OreVein{OreType: oreType}
}

func (v *OreVein) Generate(chunk *world.Chunk, random *rand.Rand, sourceX, sourceY, sourceZ int) bool {
	ore := v.OreType

	if sourceY < ore.MinY || sourceY > ore.MaxY {
		return false
	}

	chunkX := int(chunk.X) * 16
	chunkZ := int(chunk.Z) * 16

	count := ore.Size
	if count <= 0 {
		return false
	}

	angle := random.Float64() * math.Pi

	dx := math.Sin(angle) * float64(count) / 8.0
	dz := math.Cos(angle) * float64(count) / 8.0

	x1 := float64(sourceX) + dx
	x2 := float64(sourceX) - dx
	z1 := float64(sourceZ) + dz
	z2 := float64(sourceZ) - dz

	y1 := float64(sourceY) + float64(random.Intn(3)) - 2
	y2 := float64(sourceY) + float64(random.Intn(3)) - 2

	placed := 0

	for i := 0; i < count; i++ {
		t := float64(i) / float64(count)

		centerX := x1 + (x2-x1)*t
		centerY := y1 + (y2-y1)*t
		centerZ := z1 + (z2-z1)*t

		radius := (math.Sin(t*math.Pi) + 1) * random.Float64() * float64(count) / 16.0

		minX := int(math.Floor(centerX - radius))
		maxX := int(math.Floor(centerX + radius))
		minY := int(math.Floor(centerY - radius))
		maxY := int(math.Floor(centerY + radius))
		minZ := int(math.Floor(centerZ - radius))
		maxZ := int(math.Floor(centerZ + radius))

		for x := minX; x <= maxX; x++ {
			dx := (float64(x) + 0.5 - centerX) / radius
			if dx*dx >= 1 {
				continue
			}

			for y := minY; y <= maxY; y++ {
				if y < 0 || y >= 128 {
					continue
				}

				dy := (float64(y) + 0.5 - centerY) / radius
				if dx*dx+dy*dy >= 1 {
					continue
				}

				for z := minZ; z <= maxZ; z++ {
					dz := (float64(z) + 0.5 - centerZ) / radius
					if dx*dx+dy*dy+dz*dz >= 1 {
						continue
					}

					localX := x - chunkX
					localZ := z - chunkZ
					if localX < 0 || localX >= 16 || localZ < 0 || localZ >= 16 {
						continue
					}

					blockID, _ := chunk.GetBlock(localX, y, localZ)
					if blockID == byte(ore.TargetType) {
						chunk.SetBlock(localX, y, localZ, byte(ore.Material), ore.Meta)
						placed++
					}
				}
			}
		}
	}

	return placed > 0
}
