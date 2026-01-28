package layer

type GenLayerIsland struct {
	*BaseLayer
}

func NewGenLayerIsland(seed int64) *GenLayerIsland {
	return &GenLayerIsland{
		BaseLayer: NewBaseLayer(seed),
	}
}

func (l *GenLayerIsland) GetInts(x, z, width, depth int) []int {

	result := make([]int, width*depth)

	for i := 0; i < depth; i++ {
		for j := 0; j < width; j++ {
			l.InitChunkSeed(int64(x+j), int64(z+i))

			val := 0
			if l.NextInt(10) == 0 {
				val = 1
			}
			result[j+i*width] = val
		}
	}

	if x > -width && x <= 0 && z > -depth && z <= 0 {
		result[-x+-z*width] = 1
	}

	return result
}
