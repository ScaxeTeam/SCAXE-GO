package layer

type ModeSelector func(l *GenLayerZoom, v1, v2, v3, v4 int) int

type GenLayerZoom struct {
	*BaseLayer
	Selector ModeSelector
}

func NewGenLayerZoom(seed int64, parent GenLayer) *GenLayerZoom {
	z := &GenLayerZoom{
		BaseLayer: NewBaseLayer(seed),
		Selector:  selectModeOrRandom,
	}
	z.Parent = parent
	return z
}

func NewGenLayerFuzzyZoom(seed int64, parent GenLayer) *GenLayerZoom {
	z := &GenLayerZoom{
		BaseLayer: NewBaseLayer(seed),
		Selector:  selectRandom,
	}
	z.Parent = parent
	return z
}

func Magnify(seed int64, parent GenLayer, times int) GenLayer {
	layer := parent
	for i := 0; i < times; i++ {
		layer = NewGenLayerZoom(seed+int64(i), layer)
	}
	return layer
}

func (l *GenLayerZoom) GetInts(areaX, areaY, width, height int) []int {

	parentX := areaX >> 1
	parentY := areaY >> 1
	parentWidth := (width >> 1) + 2
	parentHeight := (height >> 1) + 2

	parentInts := l.Parent.GetInts(parentX, parentY, parentWidth, parentHeight)

	outWidth := (parentWidth - 1) << 1
	outHeight := (parentHeight - 1) << 1

	tempInts := make([]int, outWidth*outHeight)

	for y := 0; y < parentHeight-1; y++ {
		index := (y << 1) * outWidth

		parentIndex := y * parentWidth

		var val0, val1, val2, val3 int

		val0 = parentInts[parentIndex]
		parentIndex++
		val1 = parentInts[parentIndex]

		for x := 0; x < parentWidth-1; x++ {

			val2 = parentInts[parentIndex+parentWidth-1]
			val3 = parentInts[parentIndex+parentWidth]

			l.InitChunkSeed(int64((x+parentX)<<1), int64((y+parentY)<<1))

			tempInts[index] = val0
			tempInts[index+outWidth] = l.SelectRandom(val0, val2)

			index++
			tempInts[index] = l.SelectRandom(val0, val1)

			tempInts[index+outWidth] = l.Selector(l, val0, val1, val2, val3)
			index++

			val0 = val1
			parentIndex++
			if x < parentWidth-2 {
				val1 = parentInts[parentIndex]
			}
		}
	}

	result := make([]int, width*height)
	for y := 0; y < height; y++ {

		srcY := y + (areaY & 1)
		srcOffset := srcY*outWidth + (areaX & 1)

		copy(result[y*width:], tempInts[srcOffset:srcOffset+width])
	}

	return result
}

func selectRandom(l *GenLayerZoom, v1, v2, v3, v4 int) int {
	return l.SelectRandom(v1, v2, v3, v4)
}

func selectModeOrRandom(l *GenLayerZoom, v1, v2, v3, v4 int) int {
	if v2 == v3 && v3 == v4 {
		return v2
	}
	if v1 == v2 && v1 == v3 {
		return v1
	}
	if v1 == v2 && v1 == v4 {
		return v1
	}
	if v1 == v3 && v1 == v4 {
		return v1
	}
	if v1 == v2 && v3 != v4 {
		return v1
	}
	if v1 == v3 && v2 != v4 {
		return v1
	}
	if v1 == v4 && v2 != v3 {
		return v1
	}
	if v2 == v3 && v1 != v4 {
		return v2
	}
	if v2 == v4 && v1 != v3 {
		return v2
	}
	if v3 == v4 && v1 != v2 {
		return v3
	}

	return l.SelectRandom(v1, v2, v3, v4)
}
