package gaussian

import (
	"math"
)

type Gaussian struct {
	Size        int
	Kernel1D    []float64
	WeightSum1D float64
	Kernel2D    [][]float64
	WeightSum   float64
}

func NewGaussian(smoothSize int) *Gaussian {
	g := &Gaussian{
		Size:        smoothSize,
		Kernel1D:    make([]float64, smoothSize*2+1),
		Kernel2D:    make([][]float64, smoothSize*2+1),
		WeightSum1D: 0,
		WeightSum:   0,
	}

	bellSize := 1.0 / float64(smoothSize)
	bellHeight := 2.0 * float64(smoothSize)

	for sx := -smoothSize; sx <= smoothSize; sx++ {
		bx := bellSize * float64(sx)

		val := math.Sqrt(bellHeight) * math.Exp(-(bx*bx)/2.0)
		g.Kernel1D[sx+smoothSize] = val

		g.Kernel2D[sx+smoothSize] = make([]float64, smoothSize*2+1)

		for sz := -smoothSize; sz <= smoothSize; sz++ {
			bz := bellSize * float64(sz)

			val2D := bellHeight * math.Exp(-(bx*bx+bz*bz)/2.0)
			g.Kernel2D[sx+smoothSize][sz+smoothSize] = val2D
		}
	}

	sum1D := 0.0
	for _, v := range g.Kernel1D {
		sum1D += v
	}
	g.WeightSum1D = sum1D
	g.WeightSum = sum1D * sum1D

	return g
}
