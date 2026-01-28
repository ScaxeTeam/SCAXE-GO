package noise

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

var GrassColorNoise *PerlinNoiseGenerator

func init() {

	GrassColorNoise = NewGrassColorNoiseGenerator()
}

func NewGrassColorNoiseGenerator() *PerlinNoiseGenerator {
	r := rand.NewRandom(2345)
	return NewPerlinNoiseGeneratorWithRandom(r, 1)
}
