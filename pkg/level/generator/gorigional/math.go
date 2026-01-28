package gorigional

type MathHelper struct {
	CoordinateScale float64
	HeightScale     float64
	StretchY        float64
	DepthNoiseScale float64
}

func NewMathHelper() *MathHelper {
	return &MathHelper{

		CoordinateScale: 684.412,
		HeightScale:     684.412,

		StretchY:        12.0 * 2.0,
		DepthNoiseScale: 200.0,
	}
}

func (h *MathHelper) CalculateDensity(yIndex int, baseHeight float64, volatility float64, minLimit, maxLimit, mainNoise float64) float64 {

	d1 := (float64(yIndex) - baseHeight) * h.StretchY * 0.5 / volatility

	if d1 < 0.0 {
		d1 *= 4.0
	}

	d2 := minLimit / 512.0
	d3 := maxLimit / 512.0
	d4 := (mainNoise/10.0 + 1.0) / 2.0

	d5 := clampedLerp(d2, d3, d4) - d1
	return d5
}

func clampedLerp(lowerB, upperB, t float64) float64 {
	if t < 0.0 {
		return lowerB
	} else if t > 1.0 {
		return upperB
	}
	return lowerB + (upperB-lowerB)*t
}

func GetSquashedBaseHeight(baseSize float64, heightNoise float64) float64 {

	d0 := baseSize + heightNoise*4.0
	return d0 / 2.0
}
