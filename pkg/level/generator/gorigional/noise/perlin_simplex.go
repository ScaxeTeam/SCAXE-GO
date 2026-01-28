package noise

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type PerlinSimplexGenerator struct {
	levels      int
	noiseLevels []*SimplexNoise
}

func NewPerlinSimplexGenerator(r *rand.Random, levels int) *PerlinSimplexGenerator {
	g := &PerlinSimplexGenerator{
		levels:      levels,
		noiseLevels: make([]*SimplexNoise, levels),
	}

	for i := 0; i < levels; i++ {
		g.noiseLevels[i] = NewSimplexNoise(r)
	}

	return g
}

func (g *PerlinSimplexGenerator) GetValue(x, y float64) float64 {
	d0 := 0.0
	d1 := 1.0

	for i := 0; i < g.levels; i++ {

		d0 += g.noiseLevels[i].GetValue(x*d1, y*d1) / d1
		d1 /= 2.0
	}

	return d0
}

func (g *PerlinSimplexGenerator) GetRegion(noiseArray []float64, x, z float64, sizeX, sizeZ int, scaleX, scaleZ, scaleExp float64) []float64 {
	return g.GetRegionWithDivide(noiseArray, x, z, sizeX, sizeZ, scaleX, scaleZ, scaleExp, 0.5)
}

func (g *PerlinSimplexGenerator) GetRegionWithDivide(noiseArray []float64, x, z float64, sizeX, sizeZ int, scaleX, scaleZ, scaleExp, val float64) []float64 {
	if noiseArray != nil && len(noiseArray) >= sizeX*sizeZ {
		for i := range noiseArray {
			noiseArray[i] = 0.0
		}
	} else {
		noiseArray = make([]float64, sizeX*sizeZ)
	}

	d1 := 1.0
	d0 := 1.0

	for j := 0; j < g.levels; j++ {

		g.noiseLevels[j].Add(noiseArray, x, z, sizeX, sizeZ, scaleX*d0*d1, scaleZ*d0*d1, 0.55/d1)
		d0 *= scaleExp
		d1 *= val
	}

	return noiseArray
}
