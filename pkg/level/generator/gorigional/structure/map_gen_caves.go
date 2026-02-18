package structure

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/level/generator"
	"github.com/scaxe/scaxe-go/pkg/math/java"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type MapGenCaves struct {
	*MapGenBase
	MaxHeight int
}

func NewMapGenCaves(seed int64) *MapGenCaves {
	return &MapGenCaves{
		MapGenBase: NewMapGenBase(seed),
		MaxHeight:  256,
	}
}

func (c *MapGenCaves) GenerateChunk(cx, cz int32, chunk *world.Chunk) {
	c.Generate(int(cx), int(cz), chunk, c)
}

func (c *MapGenCaves) RecursiveGenerate(chunkX, chunkZ, originX, originZ int, chunk BlockSetter) {

	nodes := c.Rand.NextBoundedInt(c.Rand.NextBoundedInt(c.Rand.NextBoundedInt(15)+1) + 1)
	if c.Rand.NextBoundedInt(7) != 0 {
		nodes = 0
	}

	for j := 0; j < nodes; j++ {
		rx := float64(originX*16 + c.Rand.NextBoundedInt(16))
		ry := float64(c.Rand.NextBoundedInt(c.Rand.NextBoundedInt(120) + 8))
		rz := float64(originZ*16 + c.Rand.NextBoundedInt(16))

		count := 1
		if c.Rand.NextBoundedInt(4) == 0 {
			c.AddRoom(c.Rand.NextLong(), chunkX, chunkZ, chunk, rx, ry, rz)
			count += c.Rand.NextBoundedInt(4)
		}

		for k := 0; k < count; k++ {
			f := c.Rand.NextFloat() * float64(math.Pi) * 2.0
			f1 := (c.Rand.NextFloat() - 0.5) * 2.0 / 8.0
			size := c.Rand.NextFloat()*2.0 + c.Rand.NextFloat()

			if c.Rand.NextBoundedInt(10) == 0 {
				size *= (c.Rand.NextFloat()*c.Rand.NextFloat()*3.0 + 1.0)
			}

			c.AddTunnel(c.Rand.NextLong(), chunkX, chunkZ, chunk, rx, ry, rz, float32(size), float32(f), float32(f1), 0, 0, 1.0)
		}
	}
}

func (c *MapGenCaves) AddRoom(seed int64, chunkX, chunkZ int, chunk BlockSetter, x, y, z float64) {
	c.AddTunnel(seed, chunkX, chunkZ, chunk, x, y, z, 1.0+float32(c.Rand.NextFloat())*6.0, 0.0, 0.0, -1, -1, 0.5)
}

func (c *MapGenCaves) AddTunnel(seed int64, chunkX, chunkZ int, chunk BlockSetter, x, y, z float64, size, yaw, pitch float32, currentStep, tunnelLength int, heightScale float64) {
	cx := float64(chunkX*16 + 8)
	cz := float64(chunkZ*16 + 8)

	f := float32(0.0)
	f1 := float32(0.0)

	rnd := rand.NewRandom(seed)

	if tunnelLength <= 0 {
		i := c.Range*16 - 16
		tunnelLength = i - rnd.NextBoundedInt(i/4)
	}

	flag := false
	if currentStep == -1 {
		currentStep = tunnelLength / 2
		flag = true
	}

	j := rnd.NextBoundedInt(tunnelLength/2) + tunnelLength/4
	flag1 := rnd.NextBoundedInt(6) == 0

	for ; currentStep < tunnelLength; currentStep++ {
		d2 := 1.5 + float64(java.Sin(float32(currentStep)*float32(math.Pi)/float32(tunnelLength))*size)
		d3 := d2 * heightScale
		f2 := java.Cos(pitch)
		f3 := java.Sin(pitch)

		x += float64(java.Cos(yaw) * f2)
		y += float64(f3)
		z += float64(java.Sin(yaw) * f2)

		if flag1 {
			pitch *= 0.92
		} else {
			pitch *= 0.7
		}

		pitch += f1 * 0.1
		yaw += f * 0.1
		f1 *= 0.9
		f *= 0.75
		f1 += (float32(rnd.NextFloat()) - float32(rnd.NextFloat())) * float32(rnd.NextFloat()) * 2.0
		f += (float32(rnd.NextFloat()) - float32(rnd.NextFloat())) * float32(rnd.NextFloat()) * 4.0

		if !flag && currentStep == j && size > 1.0 && tunnelLength > 0 {
			c.AddTunnel(rnd.NextLong(), chunkX, chunkZ, chunk, x, y, z, float32(rnd.NextFloat())*0.5+0.5, yaw-float32(math.Pi)/2.0, pitch/3.0, currentStep, tunnelLength, 1.0)
			c.AddTunnel(rnd.NextLong(), chunkX, chunkZ, chunk, x, y, z, float32(rnd.NextFloat())*0.5+0.5, yaw+float32(math.Pi)/2.0, pitch/3.0, currentStep, tunnelLength, 1.0)
			return
		}

		if flag || rnd.NextBoundedInt(4) != 0 {
			d4 := x - cx
			d5 := z - cz
			d6 := float64(tunnelLength - currentStep)
			d7 := float64(size + 2.0 + 16.0)

			if d4*d4+d5*d5-d6*d6 > d7*d7 {
				return
			}

			if x >= cx-16.0-d2*2.0 && z >= cz-16.0-d2*2.0 && x <= cx+16.0+d2*2.0 && z <= cz+16.0+d2*2.0 {
				k := int(x-d2) - chunkX*16 - 1
				l := int(x+d2) - chunkX*16 + 1
				i1 := int(y-d3) - 1
				j1 := int(y+d3) + 1
				k1 := int(z-d2) - chunkZ*16 - 1
				l1 := int(z+d2) - chunkZ*16 + 1

				if k < 0 {
					k = 0
				}
				if l > 16 {
					l = 16
				}
				if i1 < 1 {
					i1 = 1
				}
				if j1 > c.MaxHeight-8 {
					j1 = c.MaxHeight - 8
				}
				if k1 < 0 {
					k1 = 0
				}
				if l1 > 16 {
					l1 = 16
				}

				for i2 := k; i2 < l; i2++ {
					d8 := (float64(i2+chunkX*16) + 0.5 - x) / d2
					for k2 := k1; k2 < l1; k2++ {
						d9 := (float64(k2+chunkZ*16) + 0.5 - z) / d2
						if d8*d8+d9*d9 < 1.0 {
							for j2 := j1; j2 > i1; j2-- {
								d10 := (float64(j2-1) + 0.5 - y) / d3
								if d10 > -0.7 && d8*d8+d10*d10+d9*d9 < 1.0 {

									b, _ := chunk.GetBlock(i2, j2, k2)
									if b == 1 || b == 3 || b == 2 {
										aboveID, _ := chunk.GetBlock(i2, j2+1, k2)
										if aboveID == 8 || aboveID == 9 {
											continue
										}
										if j2 < 10 {
											chunk.SetBlock(i2, j2, k2, 10, 0)
										} else {
											chunk.SetBlock(i2, j2, k2, 0, 0)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

var _ generator.Generator = nil
