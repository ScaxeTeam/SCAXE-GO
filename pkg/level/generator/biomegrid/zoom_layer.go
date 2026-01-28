package biomegrid

type ZoomType int

const (
	ZoomNormal ZoomType = iota
	ZoomBlurry
)

type ZoomMapLayer struct {
	*BaseMapLayer
	ZoomType ZoomType
}

func NewZoomMapLayer(seed int64, parent MapLayer, zoomType ZoomType) *ZoomMapLayer {
	return &ZoomMapLayer{
		BaseMapLayer: NewBaseMapLayer(seed, parent),
		ZoomType:     zoomType,
	}
}

func (l *ZoomMapLayer) GenerateValues(x, z, sizeX, sizeZ int) []int {

	parentX := x >> 1
	parentZ := z >> 1
	parentSizeX := (sizeX >> 1) + 2
	parentSizeZ := (sizeZ >> 1) + 2
	parentValues := l.Parent.GenerateValues(parentX, parentZ, parentSizeX, parentSizeZ)

	zoomedSizeX := (parentSizeX - 1) << 1
	zoomedSizeZ := (parentSizeZ - 1) << 1
	zoomedValues := make([]int, zoomedSizeX*zoomedSizeZ)

	for j := 0; j < parentSizeZ-1; j++ {
		for i := 0; i < parentSizeX-1; i++ {
			l.SetCoordsSeed((parentX+i)<<1, (parentZ+j)<<1)

			val00 := parentValues[i+j*parentSizeX]
			val10 := parentValues[i+1+j*parentSizeX]
			val01 := parentValues[i+(j+1)*parentSizeX]
			val11 := parentValues[i+1+(j+1)*parentSizeX]

			outI := i << 1
			outJ := j << 1

			zoomedValues[outI+outJ*zoomedSizeX] = val00

			zoomedValues[outI+1+outJ*zoomedSizeX] = l.selectRandom(val00, val10)

			zoomedValues[outI+(outJ+1)*zoomedSizeX] = l.selectRandom(val00, val01)

			zoomedValues[outI+1+(outJ+1)*zoomedSizeX] = l.selectModeOrRandom(val00, val10, val01, val11)
		}
	}

	result := make([]int, sizeX*sizeZ)
	offsetX := x & 1
	offsetZ := z & 1
	for j := 0; j < sizeZ; j++ {
		for i := 0; i < sizeX; i++ {
			result[i+j*sizeX] = zoomedValues[i+offsetX+(j+offsetZ)*zoomedSizeX]
		}
	}

	return result
}

func (l *ZoomMapLayer) selectRandom(a, b int) int {
	if l.NextInt(2) == 0 {
		return a
	}
	return b
}

func (l *ZoomMapLayer) selectModeOrRandom(a, b, c, d int) int {
	if l.ZoomType == ZoomBlurry {

		switch l.NextInt(4) {
		case 0:
			return a
		case 1:
			return b
		case 2:
			return c
		default:
			return d
		}
	}

	if a == b && b == c {
		return a
	}
	if a == b && b == d {
		return a
	}
	if a == c && c == d {
		return a
	}
	if b == c && c == d {
		return b
	}
	if a == b {
		return a
	}
	if a == c {
		return a
	}
	if a == d {
		return a
	}
	if b == c {
		return b
	}
	if b == d {
		return b
	}
	if c == d {
		return c
	}

	switch l.NextInt(4) {
	case 0:
		return a
	case 1:
		return b
	case 2:
		return c
	default:
		return d
	}
}
