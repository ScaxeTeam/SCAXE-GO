package layer

type GenLayerDeepOcean struct {
	*BaseLayer
}

func NewGenLayerDeepOcean(seed int64, parent GenLayer) *GenLayerDeepOcean {
	l := &GenLayerDeepOcean{
		BaseLayer: NewBaseLayer(seed),
	}
	l.Parent = parent
	return l
}

func (l *GenLayerDeepOcean) GetInts(areaX, areaY, width, height int) []int {
	parentX := areaX - 1
	parentY := areaY - 1
	parentWidth := width + 2
	parentHeight := height + 2

	parentInts := l.Parent.GetInts(parentX, parentY, parentWidth, parentHeight)
	result := make([]int, width*height)

	k := parentWidth

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			center := parentInts[x+1+(y+1)*k]

			up := parentInts[x+1+(y+0)*k]
			right := parentInts[x+2+(y+1)*k]
			left := parentInts[x+0+(y+1)*k]
			down := parentInts[x+1+(y+2)*k]

			count := 0
			if up == 0 {
				count++
			}
			if right == 0 {
				count++
			}
			if left == 0 {
				count++
			}
			if down == 0 {
				count++
			}

			if center == 0 && count == 4 {

				result[x+y*width] = 24
			} else {
				result[x+y*width] = center
			}
		}
	}

	return result
}
