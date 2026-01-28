package noise

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type SimplexNoise struct {
	xo, yo, zo float64
	p          [512]int
}

var grad3 = [12][3]int{
	{1, 1, 0}, {-1, 1, 0}, {1, -1, 0}, {-1, -1, 0},
	{1, 0, 1}, {-1, 0, 1}, {1, 0, -1}, {-1, 0, -1},
	{0, 1, 1}, {0, -1, 1}, {0, 1, -1}, {0, -1, -1},
}

var sqrt3 = math.Sqrt(3.0)
var f2 = 0.5 * (sqrt3 - 1.0)
var g2 = (3.0 - sqrt3) / 6.0

func NewSimplexNoise(r *rand.Random) *SimplexNoise {
	s := &SimplexNoise{}
	s.xo = r.NextDouble() * 256.0
	s.yo = r.NextDouble() * 256.0
	s.zo = r.NextDouble() * 256.0

	for i := 0; i < 256; i++ {
		s.p[i] = i
	}

	for i := 0; i < 256; i++ {
		j := r.NextBoundedInt(256-i) + i
		k := s.p[i]
		s.p[i] = s.p[j]
		s.p[j] = k
		s.p[i+256] = s.p[i]
	}

	return s
}

func fastFloor(x float64) int {
	if x > 0 {
		return int(x)
	}
	return int(x) - 1
}

func dot(g [3]int, x, y float64) float64 {
	return float64(g[0])*x + float64(g[1])*y
}

func (s *SimplexNoise) Add(out []float64, x, z float64, width, height int, scaleX, scaleZ, scaleMod float64) {
	idx := 0
	for i := 0; i < height; i++ {
		d0 := (z+float64(i))*scaleZ + s.yo
		for j := 0; j < width; j++ {
			d1 := (x+float64(j))*scaleX + s.xo

			d5 := (d1 + d0) * f2
			l := fastFloor(d1 + d5)
			i1 := fastFloor(d0 + d5)
			d6 := float64(l+i1) * g2
			d7 := float64(l) - d6
			d8 := float64(i1) - d6
			d9 := d1 - d7
			d10 := d0 - d8

			j1 := 0
			k1 := 1
			if d9 > d10 {
				j1 = 1
				k1 = 0
			}

			d11 := d9 - float64(j1) + g2
			d12 := d10 - float64(k1) + g2
			d13 := d9 - 1.0 + 2.0*g2
			d14 := d10 - 1.0 + 2.0*g2

			l1 := l & 255
			i2 := i1 & 255

			j2 := s.p[l1+s.p[i2]] % 12
			k2 := s.p[l1+j1+s.p[i2+k1]] % 12
			l2 := s.p[l1+1+s.p[i2+1]] % 12

			d15 := 0.5 - d9*d9 - d10*d10
			d2 := 0.0
			if d15 >= 0 {
				d15 *= d15
				d2 = d15 * d15 * dot(grad3[j2], d9, d10)
			}

			d16 := 0.5 - d11*d11 - d12*d12
			d3 := 0.0
			if d16 >= 0 {
				d16 *= d16
				d3 = d16 * d16 * dot(grad3[k2], d11, d12)
			}

			d17 := 0.5 - d13*d13 - d14*d14
			d4 := 0.0
			if d17 >= 0 {
				d17 *= d17
				d4 = d17 * d17 * dot(grad3[l2], d13, d14)
			}

			out[idx] += 70.0 * (d2 + d3 + d4) * scaleMod
			idx++
		}
	}
}

func (s *SimplexNoise) GetValue(pX, pZ float64) float64 {

	d3 := 0.5 * (sqrt3 - 1.0)
	d4 := (pX + pZ) * d3
	i := fastFloor(pX + d4)
	j := fastFloor(pZ + d4)
	d5 := (3.0 - sqrt3) / 6.0
	d6 := float64(i+j) * d5
	d7 := float64(i) - d6
	d8 := float64(j) - d6
	d9 := pX - d7
	d10 := pZ - d8

	k := 0
	l := 0
	if d9 > d10 {
		k = 1
		l = 0
	} else {
		k = 0
		l = 1
	}

	d11 := d9 - float64(k) + d5
	d12 := d10 - float64(l) + d5
	d13 := d9 - 1.0 + 2.0*d5
	d14 := d10 - 1.0 + 2.0*d5

	i1 := i & 255
	j1 := j & 255
	k1 := s.p[i1+s.p[j1]] % 12
	l1 := s.p[i1+k+s.p[j1+l]] % 12
	i2 := s.p[i1+1+s.p[j1+1]] % 12

	d15 := 0.5 - d9*d9 - d10*d10
	d0 := 0.0
	if d15 >= 0 {
		d15 *= d15
		d0 = d15 * d15 * dot(grad3[k1], d9, d10)
	}

	d16 := 0.5 - d11*d11 - d12*d12
	d1 := 0.0
	if d16 >= 0 {
		d16 *= d16
		d1 = d16 * d16 * dot(grad3[l1], d11, d12)
	}

	d17 := 0.5 - d13*d13 - d14*d14
	d2 := 0.0
	if d17 >= 0 {
		d17 *= d17
		d2 = d17 * d17 * dot(grad3[i2], d13, d14)
	}

	return 70.0 * (d0 + d1 + d2)
}
