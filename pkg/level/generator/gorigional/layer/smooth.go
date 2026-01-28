package layer

type GenLayerSmooth struct {
	*BaseLayer
}

func NewGenLayerSmooth(baseSeed int64, parent GenLayer) *GenLayerSmooth {
	l := &GenLayerSmooth{
		BaseLayer: NewBaseLayer(baseSeed),
	}
	l.Parent = parent
	return l
}

func (l *GenLayerSmooth) GetInts(x, z, width, depth int) []int {
	parentInts := l.Parent.GetInts(x-1, z-1, width+2, depth+2)
	output := make([]int, width*depth)
	parentWidth := width + 2

	for i := 0; i < depth; i++ {
		for j := 0; j < width; j++ {
			center := parentInts[j+1+(i+1)*parentWidth]
			right := parentInts[j+2+(i+1)*parentWidth]
			left := parentInts[j+(i+1)*parentWidth]
			down := parentInts[j+1+(i+2)*parentWidth]
			up := parentInts[j+1+(i)*parentWidth]

			if left == right && up == down {
				l.InitChunkSeed(int64(x+j), int64(z+i))
				if l.NextInt(2) == 0 {
					output[j+i*width] = left
				} else {
					output[j+i*width] = up
				}
			} else {
				if left == right {
					output[j+i*width] = left
				} else if up == down {
					output[j+i*width] = up
				} else {
					output[j+i*width] = center
				}
			}
		}
	}
	return output
}
