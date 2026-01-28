package biome

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

var GrassColorNoise *PerlinNoiseGen

func init() {
	GrassColorNoise = NewGrassColorNoiseGen()
}

func NewGrassColorNoiseGen() *PerlinNoiseGen {
	r := rand.NewRandom(2345)
	return &PerlinNoiseGen{
		simplex: newSimplexNoiseLocal(r),
	}
}

type PerlinNoiseGen struct {
	simplex *simplexNoiseLocal
}

func (p *PerlinNoiseGen) GetValue(x, z float64) float64 {
	return p.simplex.getValue(x, z)
}

type simplexNoiseLocal struct {
	xo, yo, zo float64
	p          [512]int
}

var grad3Local = [12][3]int{
	{1, 1, 0}, {-1, 1, 0}, {1, -1, 0}, {-1, -1, 0},
	{1, 0, 1}, {-1, 0, 1}, {1, 0, -1}, {-1, 0, -1},
	{0, 1, 1}, {0, -1, 1}, {0, 1, -1}, {0, -1, -1},
}

var sqrt3Local = math.Sqrt(3.0)

func newSimplexNoiseLocal(r *rand.Random) *simplexNoiseLocal {
	s := &simplexNoiseLocal{}
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

func fastFloorLocal(x float64) int {
	xi := int(x)
	if x < float64(xi) {
		return xi - 1
	}
	return xi
}

func dotLocal(g [3]int, x, y float64) float64 {
	return float64(g[0])*x + float64(g[1])*y
}

func (s *simplexNoiseLocal) getValue(pX, pZ float64) float64 {
	d3 := 0.5 * (sqrt3Local - 1.0)
	d4 := (pX + pZ) * d3
	i := fastFloorLocal(pX + d4)
	j := fastFloorLocal(pZ + d4)
	d5 := (3.0 - sqrt3Local) / 6.0
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
		d0 = d15 * d15 * dotLocal(grad3Local[k1], d9, d10)
	}

	d16 := 0.5 - d11*d11 - d12*d12
	d1 := 0.0
	if d16 >= 0 {
		d16 *= d16
		d1 = d16 * d16 * dotLocal(grad3Local[l1], d11, d12)
	}

	d17 := 0.5 - d13*d13 - d14*d14
	d2 := 0.0
	if d17 >= 0 {
		d17 *= d17
		d2 = d17 * d17 * dotLocal(grad3Local[i2], d13, d14)
	}

	return 70.0 * (d0 + d1 + d2)
}
