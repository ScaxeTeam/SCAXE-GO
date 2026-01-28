package java

import "math"

var SinTable [65536]float32

func init() {
	for i := 0; i < 65536; i++ {
		SinTable[i] = float32(math.Sin(float64(i) * math.Pi * 2.0 / 65536.0))
	}
}

func Sin(value float32) float32 {
	return SinTable[int(value*float32(10430.378))&65535]
}

func Cos(value float32) float32 {
	return SinTable[int(value*float32(10430.378)+16384.0)&65535]
}

func Floor(value float64) int {
	i := int(value)
	if value < float64(i) {
		return i - 1
	}
	return i
}

func FloorFloat(value float32) int {
	i := int(value)
	if value < float32(i) {
		return i - 1
	}
	return i
}

func Ceil(value float64) int {
	i := int(value)
	if value > float64(i) {
		return i + 1
	}
	return i
}

func Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
