package noise

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type SimplexOctaveGenerator struct {
	BaseOctaveGenerator
	generators []*Simplex
}

func NewSimplexOctaveGenerator(seed int64, octaves int, sizeX, sizeZ int) *SimplexOctaveGenerator {
	g := &SimplexOctaveGenerator{
		BaseOctaveGenerator: BaseOctaveGenerator{
			Octaves: octaves,
			SizeX:   sizeX,
			SizeY:   1,
			SizeZ:   sizeZ,
			XScale:  1.0,
			YScale:  1.0,
			ZScale:  1.0,
		},
		generators: make([]*Simplex, octaves),
	}

	r := rand.NewRandom(seed)
	for i := 0; i < octaves; i++ {
		g.generators[i] = NewSimplex(r, 1, 0.5, 2.0)
	}

	return g
}

func (g *SimplexOctaveGenerator) GetFractalBrownianMotion(x, z float64, lacunarity, persistence float64) []float64 {
	noise := make([]float64, g.SizeX*g.SizeZ)

	freq := 1.0
	amp := 1.0

	for _, gen := range g.generators {
		for i := 0; i < g.SizeX; i++ {
			for j := 0; j < g.SizeZ; j++ {
				idx := i + j*g.SizeX
				nx := (x + float64(i)) * g.XScale * freq
				nz := (z + float64(j)) * g.ZScale * freq
				noise[idx] += gen.GetNoise2D(nx, nz) * amp
			}
		}
		freq *= lacunarity
		amp *= persistence
	}

	return noise
}

func (g *SimplexOctaveGenerator) Noise2D(x, z float64) float64 {
	result := 0.0
	freq := 1.0
	amp := 1.0

	for _, gen := range g.generators {
		result += gen.GetNoise2D(x*g.XScale*freq, z*g.ZScale*freq) * amp
		freq *= 2.0
		amp *= 0.5
	}

	return result
}

func (g *SimplexOctaveGenerator) Noise3D(x, y, z float64) float64 {
	result := 0.0
	freq := 1.0
	amp := 1.0

	for _, gen := range g.generators {
		result += gen.GetNoise3D(x*g.XScale*freq, y*g.YScale*freq, z*g.ZScale*freq) * amp
		freq *= 2.0
		amp *= 0.5
	}

	return result
}
