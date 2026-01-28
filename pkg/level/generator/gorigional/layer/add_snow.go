package layer

type GenLayerAddSnow struct {
	*BaseLayer
}

func NewGenLayerAddSnow(seed int64, parent GenLayer) *GenLayerAddSnow {
	l := &GenLayerAddSnow{
		BaseLayer: NewBaseLayer(seed),
	}
	l.Parent = parent
	return l
}

func (l *GenLayerAddSnow) GetInts(areaX, areaY, width, height int) []int {
	parentX := areaX - 1
	parentY := areaY - 1
	parentWidth := width + 2
	parentHeight := height + 2

	parentInts := l.Parent.GetInts(parentX, parentY, parentWidth, parentHeight)
	result := make([]int, width*height)

	k := parentWidth

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			k1 := parentInts[x+1+(y+1)*k]

			l.InitChunkSeed(int64(x+areaX), int64(y+areaY))

			if k1 == 0 {
				result[x+y*width] = 0
			} else {

				r := l.NextInt(6)
				if r == 0 {
					r = 4
				} else if r <= 1 {
					r = 3
				} else {
					r = 1
				}
				result[x+y*width] = r
			}
		}
	}

	return result
}
