package noise

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type PerlinNoiseGenerator struct {
	levels []*SimplexNoise
}

func NewPerlinNoiseGenerator(rnd *rand.Random, octaves int) *PerlinNoiseGenerator {
	p := &PerlinNoiseGenerator{
		levels: make([]*SimplexNoise, octaves),
	}
	for i := 0; i < octaves; i++ {
		p.levels[i] = NewSimplexNoise(rnd)
	}
	return p
}

func NewPerlinNoiseGeneratorWithRandom(rnd *rand.Random, octaves int) *PerlinNoiseGenerator {
	return NewPerlinNoiseGenerator(rnd, octaves)
}

func (p *PerlinNoiseGenerator) GetRegion(buffer []float64, x, z float64, xSize, zSize int, xScale, zScale, scaleMod float64) []float64 {
	if buffer == nil || len(buffer) < xSize*zSize {
		buffer = make([]float64, xSize*zSize)
	} else {
		for i := range buffer {
			buffer[i] = 0
		}
	}

	d1 := 1.0
	d0 := 1.0

	for _, s := range p.levels {
		s.Add(buffer, x, z, xSize, zSize, xScale*d0*d1, zScale*d0*d1, 0.55/d1)
		d0 *= scaleMod
		d1 *= 0.5
	}

	return buffer
}

func (p *PerlinNoiseGenerator) GetValue(x, z float64) float64 {
	d0 := 0.0
	d1 := 1.0
	for _, s := range p.levels {
		d0 += s.GetValue(x*d1, z*d1) / d1
		d1 /= 2.0
	}
	return d0
}
