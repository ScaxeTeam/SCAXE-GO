package layer

type GenLayerRiverInit struct {
	*BaseLayer
}

func NewGenLayerRiverInit(baseSeed int64, parent GenLayer) *GenLayerRiverInit {
	l := &GenLayerRiverInit{
		BaseLayer: NewBaseLayer(baseSeed),
	}
	l.Parent = parent
	return l
}

func (l *GenLayerRiverInit) GetInts(x, z, width, depth int) []int {
	parentInts := l.Parent.GetInts(x, z, width, depth)
	output := make([]int, width*depth)

	for i := 0; i < depth; i++ {
		for j := 0; j < width; j++ {
			original := parentInts[j+i*width]
			l.InitChunkSeed(int64(x+j), int64(z+i))

			if original > 0 {
				output[j+i*width] = l.NextInt(299999) + 2
			} else {
				output[j+i*width] = 0
			}
		}
	}
	return output
}
