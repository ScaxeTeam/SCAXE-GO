package noise

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

const (
	F2      = 0.5 * (1.7320508075688772935274463415059 - 1.0)
	G2      = (3.0 - 1.7320508075688772935274463415059) / 6.0
	G22     = G2*2.0 - 1.0
	F3      = 1.0 / 3.0
	G3      = 1.0 / 6.0
	GradMax = 12
)

var grad3 = [12][3]float64{
	{1, 1, 0}, {-1, 1, 0}, {1, -1, 0}, {-1, -1, 0},
	{1, 0, 1}, {-1, 0, 1}, {1, 0, -1}, {-1, 0, -1},
	{0, 1, 1}, {0, -1, 1}, {0, 1, -1}, {0, -1, -1},
}

type Simplex struct {
	*Noise
	perm [512]uint8
}

var _ Generator = (*Simplex)(nil)

func NewSimplex(r *rand.Random, octaves int, persistence, expansion float64) *Simplex {
	s := &Simplex{
		Noise: NewNoise(octaves, persistence, expansion),
	}

	s.OffsetX = r.NextFloat() * 256
	s.OffsetY = r.NextFloat() * 256
	s.OffsetZ = r.NextFloat() * 256

	for i := 0; i < 512; i++ {
		s.perm[i] = 0
	}

	for i := 0; i < 256; i++ {
		s.perm[i] = uint8(r.NextBoundedInt(256))
	}

	for i := 0; i < 256; i++ {
		pos := r.NextBoundedInt(256-i) + i
		old := s.perm[i]

		s.perm[i] = s.perm[pos]
		s.perm[pos] = old
		s.perm[i+256] = s.perm[i]
	}

	r.NextInt()

	return s
}

func fastFloor(value float64) int {
	if value > 0 {
		return int(value)
	}
	return int(value) - 1
}

func dot(g [3]float64, x, y float64) float64 {
	return g[0]*x + g[1]*y
}

func (s *Simplex) Add(noiseArray []float64, x, z float64, sizeX, sizeZ int, scaleX, scaleZ, scaleFactor float64) {
	index := 0
	for j := 0; j < sizeZ; j++ {
		d0 := (z+float64(j))*scaleZ + s.OffsetY

		for k := 0; k < sizeX; k++ {
			d1 := (x+float64(k))*scaleX + s.OffsetX

			d5 := (d1 + d0) * F2
			l := fastFloor(d1 + d5)
			i1 := fastFloor(d0 + d5)
			d6 := float64(l+i1) * G2
			d7 := float64(l) - d6
			d8 := float64(i1) - d6
			d9 := d1 - d7
			d10 := d0 - d8

			var j1, k1 int
			if d9 > d10 {
				j1, k1 = 1, 0
			} else {
				j1, k1 = 0, 1
			}

			d11 := d9 - float64(j1) + G2
			d12 := d10 - float64(k1) + G2
			d13 := d9 - 1.0 + 2.0*G2
			d14 := d10 - 1.0 + 2.0*G2

			l1 := l & 255
			i2 := i1 & 255

			j2 := s.perm[l1+int(s.perm[i2])] % 12
			k2 := s.perm[l1+j1+int(s.perm[i2+k1])] % 12
			l2 := s.perm[l1+1+int(s.perm[i2+1])] % 12

			d15 := 0.5 - d9*d9 - d10*d10
			var d2 float64
			if d15 < 0.0 {
				d2 = 0.0
			} else {
				d15 *= d15
				d2 = d15 * d15 * dot(grad3[j2], d9, d10)
			}

			d16 := 0.5 - d11*d11 - d12*d12
			var d3 float64
			if d16 < 0.0 {
				d3 = 0.0
			} else {
				d16 *= d16
				d3 = d16 * d16 * dot(grad3[k2], d11, d12)
			}

			d17 := 0.5 - d13*d13 - d14*d14
			var d4 float64
			if d17 < 0.0 {
				d4 = 0.0
			} else {
				d17 *= d17
				d4 = d17 * d17 * dot(grad3[l2], d13, d14)
			}

			noiseArray[index] += 70.0 * (d2 + d3 + d4) * scaleFactor
			index++
		}
	}
}

func (s *Simplex) GetNoise2D(x, y float64) float64 {
	x += s.OffsetX
	y += s.OffsetY

	skew := (x + y) * F2

	i := int(x + skew)
	j := int(y + skew)
	t := float64(i+j) * G2

	x0 := x - (float64(i) - t)
	y0 := y - (float64(j) - t)

	var i1, j1 int
	if x0 > y0 {
		i1, j1 = 1, 0
	} else {
		i1, j1 = 0, 1
	}

	x1 := x0 - float64(i1) + G2
	y1 := y0 - float64(j1) + G2
	x2 := x0 - 1.0 + 2.0*G2
	y2 := y0 - 1.0 + 2.0*G2

	ii := i & 255
	jj := j & 255

	var n0, n1, n2 float64

	t0 := 0.5 - x0*x0 - y0*y0
	if t0 > 0 {
		t0 *= t0
		t0 *= t0

		gi0 := s.perm[ii+int(s.perm[jj])] % 12
		g := grad3[gi0]
		n0 = t0 * (g[0]*x0 + g[1]*y0)
	}

	t1 := 0.5 - x1*x1 - y1*y1
	if t1 > 0 {
		t1 *= t1
		t1 *= t1
		gi1 := s.perm[ii+i1+int(s.perm[jj+j1])] % 12
		g := grad3[gi1]
		n1 = t1 * (g[0]*x1 + g[1]*y1)
	}

	t2 := 0.5 - x2*x2 - y2*y2
	if t2 > 0 {
		t2 *= t2
		t2 *= t2
		gi2 := s.perm[ii+1+int(s.perm[jj+1])] % 12
		g := grad3[gi2]
		n2 = t2 * (g[0]*x2 + g[1]*y2)
	}

	return 70.0 * (n0 + n1 + n2)
}

func (s *Simplex) GetNoise3D(x, y, z float64) float64 {
	x += s.OffsetX
	y += s.OffsetY
	z += s.OffsetZ

	sk := (x + y + z) * F3
	i := int(x + sk)
	j := int(y + sk)
	k := int(z + sk)
	t := float64(i+j+k) * G3

	x0 := x - (float64(i) - t)
	y0 := y - (float64(j) - t)
	z0 := z - (float64(k) - t)

	var i1, j1, k1, i2, j2, k2 int

	if x0 >= y0 {
		if y0 >= z0 {
			i1, j1, k1 = 1, 0, 0
			i2, j2, k2 = 1, 1, 0
		} else if x0 >= z0 {
			i1, j1, k1 = 1, 0, 0
			i2, j2, k2 = 1, 0, 1
		} else {
			i1, j1, k1 = 0, 0, 1
			i2, j2, k2 = 1, 0, 1
		}
	} else {
		if y0 < z0 {
			i1, j1, k1 = 0, 0, 1
			i2, j2, k2 = 0, 1, 1
		} else if x0 < z0 {
			i1, j1, k1 = 0, 1, 0
			i2, j2, k2 = 0, 1, 1
		} else {
			i1, j1, k1 = 0, 1, 0
			i2, j2, k2 = 1, 1, 0
		}
	}

	x1 := x0 - float64(i1) + G3
	y1 := y0 - float64(j1) + G3
	z1 := z0 - float64(k1) + G3
	x2 := x0 - float64(i2) + 2.0*G3
	y2 := y0 - float64(j2) + 2.0*G3
	z2 := z0 - float64(k2) + 2.0*G3
	x3 := x0 - 1.0 + 3.0*G3
	y3 := y0 - 1.0 + 3.0*G3
	z3 := z0 - 1.0 + 3.0*G3

	ii := i & 255
	jj := j & 255
	kk := k & 255

	var n0, n1, n2, n3 float64

	t0 := 0.6 - x0*x0 - y0*y0 - z0*z0
	if t0 > 0 {
		t0 *= t0
		t0 *= t0
		gi0 := s.perm[ii+int(s.perm[jj+int(s.perm[kk])])] % 12
		g := grad3[gi0]
		n0 = t0 * (g[0]*x0 + g[1]*y0 + g[2]*z0)
	}

	t1 := 0.6 - x1*x1 - y1*y1 - z1*z1
	if t1 > 0 {
		t1 *= t1
		t1 *= t1
		gi1 := s.perm[ii+i1+int(s.perm[jj+j1+int(s.perm[kk+k1])])] % 12
		g := grad3[gi1]
		n1 = t1 * (g[0]*x1 + g[1]*y1 + g[2]*z1)
	}

	t2 := 0.6 - x2*x2 - y2*y2 - z2*z2
	if t2 > 0 {
		t2 *= t2
		t2 *= t2
		gi2 := s.perm[ii+i2+int(s.perm[jj+j2+int(s.perm[kk+k2])])] % 12
		g := grad3[gi2]
		n2 = t2 * (g[0]*x2 + g[1]*y2 + g[2]*z2)
	}

	t3 := 0.6 - x3*x3 - y3*y3 - z3*z3
	if t3 > 0 {
		t3 *= t3
		t3 *= t3
		gi3 := s.perm[ii+1+int(s.perm[jj+1+int(s.perm[kk+1])])] % 12
		g := grad3[gi3]
		n3 = t3 * (g[0]*x3 + g[1]*y3 + g[2]*z3)
	}

	return 32.0 * (n0 + n1 + n2 + n3)
}

func (s *Simplex) GetFastNoise3D(xSize, ySize, zSize, xSampling, ySampling, zSampling, x, y, z int) []float64 {
	xArraySize := xSize + 1
	zArraySize := zSize + 1
	yArraySize := ySize + 1

	noiseArray := make([]float64, xArraySize*zArraySize*yArraySize)

	idx := func(ix, iz, iy int) int {
		return (ix*zArraySize+iz)*yArraySize + iy
	}

	for xx := 0; xx <= xSize; xx += xSampling {
		for zz := 0; zz <= zSize; zz += zSampling {
			for yy := 0; yy <= ySize; yy += ySampling {

				val := s.Noise3D(s, float64(x+xx), float64(y+yy), float64(z+zz), true)
				noiseArray[idx(xx, zz, yy)] = val
			}
		}
	}

	xLerpStep := 1.0 / float64(xSampling)
	yLerpStep := 1.0 / float64(ySampling)
	zLerpStep := 1.0 / float64(zSampling)

	for leftX := 0; leftX < xSize; leftX += xSampling {
		rightX := leftX + xSampling
		for leftZ := 0; leftZ < zSize; leftZ += zSampling {
			rightZ := leftZ + zSampling
			for leftY := 0; leftY < ySize; leftY += ySampling {
				rightY := leftY + ySampling

				c000 := noiseArray[idx(leftX, leftZ, leftY)]
				c100 := noiseArray[idx(rightX, leftZ, leftY)]
				c001 := noiseArray[idx(leftX, leftZ, rightY)]
				c101 := noiseArray[idx(rightX, leftZ, rightY)]
				c010 := noiseArray[idx(leftX, rightZ, leftY)]
				c110 := noiseArray[idx(rightX, rightZ, leftY)]
				c011 := noiseArray[idx(leftX, rightZ, rightY)]
				c111 := noiseArray[idx(rightX, rightZ, rightY)]

				for xStep := 0; xStep < xSampling; xStep++ {
					xx := leftX + xStep
					dx2 := float64(xStep) * xLerpStep
					dx1 := 1.0 - dx2

					x00 := c000*dx1 + c100*dx2
					x01 := c001*dx1 + c101*dx2
					x10 := c010*dx1 + c110*dx2
					x11 := c011*dx1 + c111*dx2

					for zStep := 0; zStep < zSampling; zStep++ {
						zz := leftZ + zStep
						dz2 := float64(zStep) * zLerpStep
						dz1 := 1.0 - dz2

						z0 := x00*dz1 + x10*dz2
						z1 := x01*dz1 + x11*dz2

						yStart := 0
						if xStep == 0 && zStep == 0 {
							yStart = 1
						}

						for yStep := yStart; yStep < ySampling; yStep++ {
							yy := leftY + yStep
							dy2 := float64(yStep) * yLerpStep
							dy1 := 1.0 - dy2

							noiseArray[idx(xx, zz, yy)] = dy1*z0 + dy2*z1
						}
					}
				}
			}
		}
	}
	return noiseArray
}
