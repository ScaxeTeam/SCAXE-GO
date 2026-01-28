package noise

import (
	"math/rand"
)

type PerlinOctaveGenerator struct {
	BaseOctaveGenerator
	generators []*PerlinNoiseGenerator
}

func NewPerlinOctaveGenerator(r *rand.Rand, octaves int, sizeX, sizeY, sizeZ int) *PerlinOctaveGenerator {
	g := &PerlinOctaveGenerator{
		BaseOctaveGenerator: BaseOctaveGenerator{
			Octaves: octaves,
			SizeX:   sizeX,
			SizeY:   sizeY,
			SizeZ:   sizeZ,
			XScale:  1.0,
			YScale:  1.0,
			ZScale:  1.0,
		},
		generators: make([]*PerlinNoiseGenerator, octaves),
	}

	for i := 0; i < octaves; i++ {
		g.generators[i] = NewPerlinNoiseGenerator(r)
	}

	return g
}

func NewPerlinOctaveGenerator2D(r *rand.Rand, octaves int, sizeX, sizeZ int) *PerlinOctaveGenerator {
	return NewPerlinOctaveGenerator(r, octaves, sizeX, 1, sizeZ)
}

func (g *PerlinOctaveGenerator) GetFractalBrownianMotion(x, z float64, lacunarity, persistence float64) []float64 {
	noise := make([]float64, g.SizeX*g.SizeZ)

	freq := 1.0
	amp := 1.0

	for _, gen := range g.generators {
		for i := 0; i < g.SizeX; i++ {
			for j := 0; j < g.SizeZ; j++ {
				idx := i + j*g.SizeX
				nx := (x + float64(i)) * g.XScale * freq
				nz := (z + float64(j)) * g.ZScale * freq
				noise[idx] += gen.Noise2D(nx, nz) * amp
			}
		}
		freq *= lacunarity
		amp *= persistence
	}

	return noise
}

func (g *PerlinOctaveGenerator) GetFractalBrownianMotion3D(x, y, z float64, lacunarity, persistence float64) []float64 {
	noise := make([]float64, g.SizeX*g.SizeY*g.SizeZ)

	freq := 1.0
	amp := 1.0

	for _, gen := range g.generators {
		for i := 0; i < g.SizeX; i++ {
			for j := 0; j < g.SizeY; j++ {
				for k := 0; k < g.SizeZ; k++ {
					idx := i + j*g.SizeX + k*g.SizeX*g.SizeY
					nx := (x + float64(i)) * g.XScale * freq
					ny := (y + float64(j)) * g.YScale * freq
					nz := (z + float64(k)) * g.ZScale * freq
					noise[idx] += gen.Noise3D(nx, ny, nz) * amp
				}
			}
		}
		freq *= lacunarity
		amp *= persistence
	}

	return noise
}

func (g *PerlinOctaveGenerator) Noise3D(x, y, z float64) float64 {
	result := 0.0
	freq := 1.0
	amp := 1.0

	for _, gen := range g.generators {
		result += gen.Noise3D(x*g.XScale*freq, y*g.YScale*freq, z*g.ZScale*freq) * amp
		freq *= 2.0
		amp *= 0.5
	}

	return result
}

func (g *PerlinOctaveGenerator) Noise2D(x, z float64) float64 {
	result := 0.0
	freq := 1.0
	amp := 1.0

	for _, gen := range g.generators {
		result += gen.Noise2D(x*g.XScale*freq, z*g.ZScale*freq) * amp
		freq *= 2.0
		amp *= 0.5
	}

	return result
}
