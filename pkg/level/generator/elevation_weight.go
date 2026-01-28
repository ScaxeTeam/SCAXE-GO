package generator

import "math"

var ElevationWeight [5][5]float64

func init() {

	for x := 0; x < 5; x++ {
		for z := 0; z < 5; z++ {
			sqX := float64((x - 2) * (x - 2))
			sqZ := float64((z - 2) * (z - 2))
			ElevationWeight[x][z] = 10.0 / math.Sqrt(sqX+sqZ+0.2)
		}
	}
}

func GetElevationWeight(x, z int) float64 {
	if x < 0 || x > 4 || z < 0 || z > 4 {
		return 0
	}
	return ElevationWeight[x][z]
}
