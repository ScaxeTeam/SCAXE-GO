package biome

import (
	"github.com/scaxe/scaxe-go/pkg/math/noise"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type BiomeLookup func(temperature, rainfall, river, ocean, hills float64) uint8

type BiomeSelector struct {
	Fallback Biome

	Temperature *noise.Simplex
	Rainfall    *noise.Simplex
	River       *noise.Simplex
	Ocean       *noise.Simplex
	Hills       *noise.Simplex

	Biomes map[uint8]Biome
	Lookup BiomeLookup
}

func NewBiomeSelector(r *rand.Random, lookup BiomeLookup, fallback Biome) *BiomeSelector {
	return &BiomeSelector{
		Fallback: fallback,
		Lookup:   lookup,
		Biomes:   make(map[uint8]Biome),

		Temperature: noise.NewSimplex(r, 2, 1.0/8.0, 1.0/2048.0),
		Rainfall:    noise.NewSimplex(r, 2, 1.0/8.0, 1.0/2048.0),
		River:       noise.NewSimplex(r, 6, 1.0/2.0, 1.0/1024.0),
		Ocean:       noise.NewSimplex(r, 6, 1.0/2.0, 1.0/2048.0),
		Hills:       noise.NewSimplex(r, 2, 1.0/2.0, 1.0/2048.0),
	}
}

func (s *BiomeSelector) AddBiome(b Biome) {
	s.Biomes[b.GetID()] = b
}

func (s *BiomeSelector) GetTemperature(x, z float64) float64 {
	return s.Temperature.Noise2D(nil, x, z, true)

}

func (s *BiomeSelector) GetRainfall(x, z float64) float64 {
	return s.Rainfall.Noise2D(s.Rainfall, x, z, true)
}

func (s *BiomeSelector) GetRiver(x, z float64) float64 {
	return s.River.Noise2D(s.River, x, z, true)
}

func (s *BiomeSelector) GetOcean(x, z float64) float64 {
	return s.Ocean.Noise2D(s.Ocean, x, z, true)
}

func (s *BiomeSelector) GetHills(x, z float64) float64 {
	return s.Hills.Noise2D(s.Hills, x, z, true)
}

func (s *BiomeSelector) PickBiome(x, z float64) Biome {
	temp := s.Temperature.Noise2D(s.Temperature, x, z, true)
	rain := s.Rainfall.Noise2D(s.Rainfall, x, z, true)
	river := s.River.Noise2D(s.River, x, z, true)
	ocean := s.Ocean.Noise2D(s.Ocean, x, z, true)
	hills := s.Hills.Noise2D(s.Hills, x, z, true)

	biomeID := s.Lookup(temp, rain, river, ocean, hills)

	if b, ok := s.Biomes[biomeID]; ok {
		return b
	}
	return s.Fallback
}
