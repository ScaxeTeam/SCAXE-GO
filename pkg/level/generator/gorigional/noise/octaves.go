package noise

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type OctavesNoise struct {
	generators []*ImprovedNoise
	octaves    int
}

func NewOctavesNoise(rnd *rand.Random, octaves int) *OctavesNoise {
	o := &OctavesNoise{
		octaves:    octaves,
		generators: make([]*ImprovedNoise, octaves),
	}

	for i := 0; i < octaves; i++ {
		o.generators[i] = NewImprovedNoise(rnd)
	}
	return o
}

func (o *OctavesNoise) GenerateNoiseOctaves(noiseArray []float64, xOffset, yOffset, zOffset int, xSize, ySize, zSize int, xScale, yScale, zScale float64) []float64 {
	if noiseArray == nil || len(noiseArray) < xSize*ySize*zSize {
		noiseArray = make([]float64, xSize*ySize*zSize)
	} else {

		clear(noiseArray)
	}

	amp := 1.0

	for j := 0; j < o.octaves; j++ {
		d0 := float64(xOffset) * amp * xScale
		d1 := float64(yOffset) * amp * yScale
		d2 := float64(zOffset) * amp * zScale

		k := int64(math.Floor(d0))
		l := int64(math.Floor(d2))

		d0 -= float64(k)
		d2 -= float64(l)

		k %= 16777216
		l %= 16777216

		d0 += float64(k)
		d2 += float64(l)

		o.generators[j].PopulateNoiseArray(noiseArray, d0, d1, d2, xSize, ySize, zSize, xScale*amp, yScale*amp, zScale*amp, amp)
		amp /= 2.0
	}

	return noiseArray
}

func (o *OctavesNoise) GenerateNoiseOctaves2D(noiseArray []float64, xOffset, zOffset int, xSize, zSize int, xScale, zScale float64) []float64 {
	return o.GenerateNoiseOctaves(noiseArray, xOffset, 10, zOffset, xSize, 1, zSize, xScale, 1.0, zScale)
}
