package biomegrid

type NoiseMapLayer struct {
	*BaseMapLayer
}

func NewNoiseMapLayer(seed int64) *NoiseMapLayer {
	return &NoiseMapLayer{
		BaseMapLayer: NewBaseMapLayer(seed, nil),
	}
}

func (l *NoiseMapLayer) GenerateValues(x, z, sizeX, sizeZ int) []int {
	values := make([]int, sizeX*sizeZ)

	for j := 0; j < sizeZ; j++ {
		for i := 0; i < sizeX; i++ {
			l.SetCoordsSeed(x+i, z+j)

			if l.NextInt(10) == 0 {
				values[i+j*sizeX] = 1
			} else {
				values[i+j*sizeX] = 0
			}
		}
	}

	return values
}
