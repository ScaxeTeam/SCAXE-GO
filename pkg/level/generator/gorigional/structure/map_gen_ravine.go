package structure

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type MapGenRavine struct {
	rangeR    int
	worldSeed int64
	rand      *rand.Random
	rs        []float64
	MaxHeight int
}

func NewMapGenRavine(seed int64) *MapGenRavine {
	return &MapGenRavine{
		rangeR:    8,
		worldSeed: seed,
		rand:      rand.NewRandom(seed),
		MaxHeight: 256,
	}
}

func (m *MapGenRavine) GenerateChunk(chunkX, chunkZ int32, chunk *world.Chunk) {
	m.rand.SetSeed(m.worldSeed)
	r1 := m.rand.NextLong()
	r2 := m.rand.NextLong()

	for x := int(chunkX) - m.rangeR; x <= int(chunkX)+m.rangeR; x++ {
		for z := int(chunkZ) - m.rangeR; z <= int(chunkZ)+m.rangeR; z++ {
			seed := int64(x)*r1 ^ int64(z)*r2 ^ m.worldSeed
			m.rand.SetSeed(seed)

			if m.rand.NextFloat() < 0.02 {
				m.recursiveGenerate(int(chunkX), int(chunkZ), x, z, chunk)
			}
		}
	}
}

func (m *MapGenRavine) recursiveGenerate(chunkX, chunkZ, x, z int, chunk *world.Chunk) {

	d0 := float64(int64(x)*16 + int64(m.rand.NextBoundedInt(16)))
	d1 := float64(m.rand.NextBoundedInt(m.rand.NextBoundedInt(40)+8) + 20)
	d2 := float64(int64(z)*16 + int64(m.rand.NextBoundedInt(16)))

	f := m.rand.NextFloat() * (math.Pi * 2.0)
	f1 := (m.rand.NextFloat() - 0.5) * 2.0 / 8.0
	f2 := (m.rand.NextFloat()*2.0 + m.rand.NextFloat()) * 2.0

	steps := 112

	steps = m.rand.NextBoundedInt(m.rand.NextBoundedInt(140) + 8)

	m.addTunnel(int64(m.rand.NextLong()), chunkX, chunkZ, chunk, d0, d1, d2, float32(f), float32(f1), float32(f2), 0, steps, 3.0)
}

func (m *MapGenRavine) addTunnel(seed int64, chunkX, chunkZ int, chunk *world.Chunk, x, y, z float64, yaw, pitch, scale float32, startStep, endStep int, heightMod float64) {
	r := rand.NewRandom(seed)

	cxMin := float64(chunkX * 16)
	czMin := float64(chunkZ * 16)

	for i := startStep; i < endStep; i++ {

		d0 := 1.5 + float64(math.Sin(float64(i)*math.Pi/float64(endStep))*float64(scale))
		d1 := d0 * heightMod

		d0 *= float64(r.NextFloat()*0.25 + 0.75)
		d1 *= float64(r.NextFloat()*0.25 + 0.75)

		x += float64(math.Cos(float64(yaw)))
		z += float64(math.Sin(float64(yaw)))
		y += float64(math.Sin(float64(pitch)))

		pitch *= 0.7
		pitch += float32(r.NextFloat()-r.NextFloat()) * 0.05
		yaw += float32(r.NextFloat()-r.NextFloat()) * 0.05

		if r.NextFloat() < 0.25 {
			r.NextFloat()
			r.NextFloat()

		}

		if x < cxMin-16.0-d0*2.0 || z < czMin-16.0-d0*2.0 || x > cxMin+16.0+d0*2.0 || z > czMin+16.0+d0*2.0 {
			continue
		}

		minX := int(math.Floor(x - d0))
		maxX := int(math.Floor(x + d0))
		minY := int(math.Floor(y - d1))
		maxY := int(math.Floor(y + d1))
		minZ := int(math.Floor(z - d0))
		maxZ := int(math.Floor(z + d0))

		if minX < int(chunkX)*16 {
			minX = int(chunkX) * 16
		}
		if maxX > int(chunkX)*16+15 {
			maxX = int(chunkX)*16 + 15
		}
		if minZ < int(chunkZ)*16 {
			minZ = int(chunkZ) * 16
		}
		if maxZ > int(chunkZ)*16+15 {
			maxZ = int(chunkZ)*16 + 15
		}
		if minY < 1 {
			minY = 1
		}
		if maxY > m.MaxHeight-8 {
			maxY = m.MaxHeight - 8
		}

		for ix := minX; ix <= maxX; ix++ {
			relX := float64(ix) + 0.5 - x
			for iz := minZ; iz <= maxZ; iz++ {
				relZ := float64(iz) + 0.5 - z

				if relX*relX+relZ*relZ < d0*d0 {
					for iy := minY; iy <= maxY; iy++ {
						relY := float64(iy) + 0.5 - y
						if relX*relX+relZ*relZ < d0*d0 && (relX*relX+relZ*relZ)*heightMod+relY*relY < d0*d0*heightMod {

							lx := ix & 15
							lz := iz & 15

							id, _ := chunk.GetBlock(lx, iy, lz)

							if id == 1 || id == 2 || id == 3 {

								chunk.SetBlock(lx, iy, lz, 0, 0)

								if iy < 10 {
									chunk.SetBlock(lx, iy, lz, 10, 0)
								}
							}
						}
					}
				}
			}
		}

		if i%2 == 0 {
			continue
		}
	}
}
