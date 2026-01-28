package layer

type EdgeMode int

const (
	CoolWarm EdgeMode = iota
	HeatIce
	Special
)

type GenLayerEdge struct {
	*BaseLayer
	Mode EdgeMode
}

func NewGenLayerEdge(seed int64, parent GenLayer, mode EdgeMode) *GenLayerEdge {
	l := &GenLayerEdge{
		BaseLayer: NewBaseLayer(seed),
		Mode:      mode,
	}
	l.Parent = parent
	return l
}

func (l *GenLayerEdge) GetInts(areaX, areaY, width, height int) []int {
	switch l.Mode {
	case CoolWarm:
		return l.getIntsCoolWarm(areaX, areaY, width, height)
	case HeatIce:
		return l.getIntsHeatIce(areaX, areaY, width, height)
	case Special:
		return l.getIntsSpecial(areaX, areaY, width, height)
	default:
		return l.getIntsCoolWarm(areaX, areaY, width, height)
	}
}

func (l *GenLayerEdge) getIntsCoolWarm(areaX, areaY, width, height int) []int {
	parentX := areaX - 1
	parentY := areaY - 1
	parentWidth := width + 2
	parentHeight := height + 2

	parentInts := l.Parent.GetInts(parentX, parentY, parentWidth, parentHeight)
	result := make([]int, width*height)

	k := parentWidth

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			l.InitChunkSeed(int64(x+areaX), int64(y+areaY))
			k1 := parentInts[x+1+(y+1)*k]

			if k1 == 1 {
				l1 := parentInts[x+1+(y+0)*k]
				i2 := parentInts[x+2+(y+1)*k]
				j2 := parentInts[x+0+(y+1)*k]
				k2 := parentInts[x+1+(y+2)*k]

				flag := (l1 == 3 || i2 == 3 || j2 == 3 || k2 == 3)
				flag1 := (l1 == 4 || i2 == 4 || j2 == 4 || k2 == 4)

				if flag || flag1 {
					k1 = 2

				}
			}
			result[x+y*width] = k1
		}
	}
	return result
}

func (l *GenLayerEdge) getIntsHeatIce(areaX, areaY, width, height int) []int {

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
			if k1 == 4 {
				l1 := parentInts[x+1+(y+0)*k]
				i2 := parentInts[x+2+(y+1)*k]
				j2 := parentInts[x+0+(y+1)*k]
				k2 := parentInts[x+1+(y+2)*k]

				flag := (l1 == 2 || i2 == 2 || j2 == 2 || k2 == 2)
				flag1 := (l1 == 1 || i2 == 1 || j2 == 1 || k2 == 1)

				if flag || flag1 {
					k1 = 3
				}
			}
			result[x+y*width] = k1
		}
	}
	return result
}

func (l *GenLayerEdge) getIntsSpecial(areaX, areaY, width, height int) []int {

	parentInts := l.Parent.GetInts(areaX, areaY, width, height)
	result := make([]int, width*height)

	for i := 0; i < len(parentInts); i++ {
		l.InitChunkSeed(int64(i%width+areaX), int64(i/width+areaY))
		k := parentInts[i]

		if k != 0 && l.NextInt(13) == 0 {
			k |= (1 + l.NextInt(15)) << 8 & 3840
		}
		result[i] = k
	}
	return result
}
