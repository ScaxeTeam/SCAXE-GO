package noise

import (
	"math"
)

type Generator interface {
	GetNoise2D(x, z float64) float64

	GetNoise3D(x, y, z float64) float64

	SetOffset(x, y, z float64)
}

type Noise struct {
	Octaves     int
	Persistence float64
	Expansion   float64
	OffsetX     float64
	OffsetY     float64
	OffsetZ     float64
}

func NewNoise(octaves int, persistence, expansion float64) *Noise {
	return &Noise{
		Octaves:     octaves,
		Persistence: persistence,
		Expansion:   expansion,
	}
}

func (n *Noise) SetOffset(x, y, z float64) {
	n.OffsetX = x
	n.OffsetY = y
	n.OffsetZ = z
}

func LinearLerp(x, x1, x2, q0, q1 float64) float64 {
	return ((x2-x)/(x2-x1))*q0 + ((x-x1)/(x2-x1))*q1
}

func BilinearLerp(x, y, q00, q01, q10, q11, x1, x2, y1, y2 float64) float64 {
	dx1 := (x2 - x) / (x2 - x1)
	dx2 := (x - x1) / (x2 - x1)

	return ((y2-y)/(y2-y1))*(dx1*q00+dx2*q10) + ((y-y1)/(y2-y1))*(dx1*q01+dx2*q11)
}

func TrilinearLerp(x, y, z, q000, q001, q010, q011, q100, q101, q110, q111, x1, x2, y1, y2, z1, z2 float64) float64 {
	dx1 := (x2 - x) / (x2 - x1)
	dx2 := (x - x1) / (x2 - x1)
	dy1 := (y2 - y) / (y2 - y1)
	dy2 := (y - y1) / (y2 - y1)

	return ((z2-z)/(z2-z1))*(dy1*(dx1*q000+dx2*q100)+dy2*(dx1*q001+dx2*q101)) +
		((z-z1)/(z2-z1))*(dy1*(dx1*q010+dx2*q110)+dy2*(dx1*q011+dx2*q111))
}

func (n *Noise) Mix1D(gen Generator, x float64, normalized bool) float64 {

	return 0
}

func (n *Noise) Noise2D(gen Generator, x, z float64, normalized bool) float64 {
	result := 0.0
	amp := 1.0
	freq := 1.0
	maxVal := 0.0

	x *= n.Expansion
	z *= n.Expansion

	for i := 0; i < n.Octaves; i++ {
		result += gen.GetNoise2D(x*freq, z*freq) * amp
		maxVal += amp
		freq *= 2
		amp *= n.Persistence
	}

	if normalized {
		result /= maxVal
	}

	return result
}

func (n *Noise) Noise3D(gen Generator, x, y, z float64, normalized bool) float64 {
	result := 0.0
	amp := 1.0
	freq := 1.0
	maxVal := 0.0

	x *= n.Expansion
	y *= n.Expansion
	z *= n.Expansion

	for i := 0; i < n.Octaves; i++ {
		result += gen.GetNoise3D(x*freq, y*freq, z*freq) * amp
		maxVal += amp
		freq *= 2
		amp *= n.Persistence
	}

	if normalized {
		result /= maxVal
	}

	return result
}

func RidgedNoise(gen Generator, x, y, z float64) float64 {
	val := gen.GetNoise3D(x, y, z)
	return 1.0 - math.Abs(val)
}

func DomainWarp(gen Generator, x, y, z, strength float64) float64 {
	q := gen.GetNoise3D(x, y, z)
	return gen.GetNoise3D(x+strength*q, y+strength*q, z+strength*q)
}

func DomainWarp2D(gen Generator, x, z, strength float64) float64 {
	q := gen.GetNoise2D(x, z)
	return gen.GetNoise2D(x+strength*q, z+strength*q)
}
