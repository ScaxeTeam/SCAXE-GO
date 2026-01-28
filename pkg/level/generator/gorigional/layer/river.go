package layer

type GenLayerRiver struct {
	*BaseLayer
}

func NewGenLayerRiver(baseSeed int64, parent GenLayer) *GenLayerRiver {
	l := &GenLayerRiver{
		BaseLayer: NewBaseLayer(baseSeed),
	}
	l.Parent = parent
	return l
}

func (l *GenLayerRiver) GetInts(x, z, width, depth int) []int {

	parentInts := l.Parent.GetInts(x-1, z-1, width+2, depth+2)
	output := make([]int, width*depth)

	parentWidth := width + 2

	for i := 0; i < depth; i++ {
		for j := 0; j < width; j++ {

			center := l.riverFilter(parentInts[j+1+(i+1)*parentWidth])

			left := l.riverFilter(parentInts[j+(i+1)*parentWidth])
			right := l.riverFilter(parentInts[j+2+(i+1)*parentWidth])
			up := l.riverFilter(parentInts[j+1+(i)*parentWidth])
			down := l.riverFilter(parentInts[j+1+(i+2)*parentWidth])

			if center == left && center == right && center == up && center == down {
				output[j+i*width] = -1
			} else {
				output[j+i*width] = 7
			}
		}
	}
	return output
}

func (l *GenLayerRiver) riverFilter(val int) int {
	if val >= 2 {
		return 2 + (val & 1)
	}
	return val
}
