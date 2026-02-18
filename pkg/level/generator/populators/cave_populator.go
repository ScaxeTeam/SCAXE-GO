package populators

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const BlockAir = 0

type CavePopulator struct{}

func NewCavePopulator() *CavePopulator {
	return &CavePopulator{}
}

func (cp *CavePopulator) Populate(level populator.ChunkManager, chunk *world.Chunk, cx, cz int32, random *rand.Random) {

	chunkX := int(chunk.X)
	chunkZ := int(chunk.Z)

	for dx := -8; dx <= 8; dx++ {
		for dz := -8; dz <= 8; dz++ {

			nx := chunkX + dx
			nz := chunkZ + dz
			seed := int64(nx)*341873128712 + int64(nz)*132897987541
			caveRandom := rand.NewRandom(seed)

			if caveRandom.NextBoundedInt(16) == 0 {
				cp.generateCaveSystem(chunk, caveRandom, nx, nz)
			}
		}
	}
}

func (cp *CavePopulator) generateCaveSystem(chunk *world.Chunk, random *rand.Random, originX, originZ int) {

	tunnelCount := random.NextBoundedInt(40) + 1

	for i := 0; i < tunnelCount; i++ {

		x := float64(originX*16 + random.NextBoundedInt(16))
		y := float64(random.NextBoundedInt(128))
		z := float64(originZ*16 + random.NextBoundedInt(16))

		yaw := random.NextFloat() * 2 * math.Pi
		pitch := (random.NextFloat() - 0.5) * math.Pi / 4

		cp.generateTunnel(chunk, random, x, y, z, yaw, pitch)
	}
}

func (cp *CavePopulator) generateTunnel(chunk *world.Chunk, random *rand.Random, x, y, z, yaw, pitch float64) {
	chunkX := int(chunk.X) * 16
	chunkZ := int(chunk.Z) * 16

	length := random.NextBoundedInt(200) + 50
	radius := 1.0 + random.NextFloat()*2.0

	for step := 0; step < length; step++ {

		dx := math.Cos(yaw) * math.Cos(pitch)
		dy := math.Sin(pitch)
		dz := math.Sin(yaw) * math.Cos(pitch)

		x += dx
		y += dy
		z += dz

		yaw += (random.NextFloat() - 0.5) * 0.3
		pitch += (random.NextFloat() - 0.5) * 0.2
		pitch = math.Max(-0.5, math.Min(0.5, pitch))

		currentRadius := radius * (0.5 + 0.5*math.Sin(float64(step)*0.1))

		minX := int(math.Floor(x - currentRadius))
		maxX := int(math.Ceil(x + currentRadius))
		minY := int(math.Max(1, math.Floor(y-currentRadius)))
		maxY := int(math.Min(127, math.Ceil(y+currentRadius)))
		minZ := int(math.Floor(z - currentRadius))
		maxZ := int(math.Ceil(z + currentRadius))

		for cx := minX; cx <= maxX; cx++ {
			for cy := minY; cy <= maxY; cy++ {
				for cz := minZ; cz <= maxZ; cz++ {

					distSq := math.Pow(float64(cx)-x, 2) +
						math.Pow(float64(cy)-y, 2) +
						math.Pow(float64(cz)-z, 2)
					if distSq > currentRadius*currentRadius {
						continue
					}

					localX := cx - chunkX
					localZ := cz - chunkZ
					if localX < 0 || localX >= 16 || localZ < 0 || localZ >= 16 {
						continue
					}

					if cy <= 0 {
						continue
					}

					blockID, _ := chunk.GetBlock(localX, cy, localZ)
					if blockID != 0 && blockID != 9 && blockID != 10 && blockID != 11 {
						if cy+1 < 128 {
							aboveID, _ := chunk.GetBlock(localX, cy+1, localZ)
							if aboveID == 9 || aboveID == 8 {
								continue
							}
						}
						chunk.SetBlock(localX, cy, localZ, BlockAir, 0)
					}
				}
			}
		}

		if random.NextBoundedInt(100) == 0 {
			break
		}
	}
}
