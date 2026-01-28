package layer

type GenLayerRemoveTooMuchOcean struct {
	*BaseLayer
}

func NewGenLayerRemoveTooMuchOcean(seed int64, parent GenLayer) *GenLayerRemoveTooMuchOcean {
	l := &GenLayerRemoveTooMuchOcean{
		BaseLayer: NewBaseLayer(seed),
	}
	l.Parent = parent
	return l
}

func (l *GenLayerRemoveTooMuchOcean) GetInts(areaX, areaY, width, height int) []int {
	parentX := areaX - 1
	parentY := areaY - 1
	parentWidth := width + 2
	parentHeight := height + 2

	parentInts := l.Parent.GetInts(parentX, parentY, parentWidth, parentHeight)
	result := make([]int, width*height)

	k := parentWidth

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			up := parentInts[x+1+(y+0)*k]
			right := parentInts[x+2+(y+1)*k]
			left := parentInts[x+0+(y+1)*k]
			down := parentInts[x+1+(y+2)*k]
			center := parentInts[x+1+(y+1)*k]

			l.InitChunkSeed(int64(x+areaX), int64(y+areaY))

			result[x+y*width] = center

			if center == 0 && up == 0 && right == 0 && left == 0 && down == 0 && l.NextInt(2) == 0 {
				result[x+y*width] = 1
			}
		}
	}

	return result
}
