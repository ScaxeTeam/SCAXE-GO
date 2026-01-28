package layer

type GenLayerAddIsland struct {
	*BaseLayer
}

func NewGenLayerAddIsland(seed int64, parent GenLayer) *GenLayerAddIsland {
	l := &GenLayerAddIsland{
		BaseLayer: NewBaseLayer(seed),
	}
	l.Parent = parent
	return l
}

func (l *GenLayerAddIsland) GetInts(areaX, areaY, width, height int) []int {
	parentX := areaX - 1
	parentY := areaY - 1
	parentWidth := width + 2
	parentHeight := height + 2

	parentInts := l.Parent.GetInts(parentX, parentY, parentWidth, parentHeight)
	result := make([]int, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			k := parentWidth

			k1 := parentInts[x+0+(y+0)*k]
			l1 := parentInts[x+2+(y+0)*k]
			i2 := parentInts[x+0+(y+2)*k]
			j2 := parentInts[x+2+(y+2)*k]

			centerVal := parentInts[x+1+(y+1)*k]

			l.InitChunkSeed(int64(x+areaX), int64(y+areaY))

			if centerVal != 0 || (k1 == 0 && l1 == 0 && i2 == 0 && j2 == 0) {

				if centerVal > 0 && (k1 == 0 || l1 == 0 || i2 == 0 || j2 == 0) {

					if l.NextInt(5) == 0 {
						if centerVal == 4 {
							result[x+y*width] = 4
						} else {
							result[x+y*width] = 0
						}
					} else {
						result[x+y*width] = centerVal
					}
				} else {
					result[x+y*width] = centerVal
				}
			} else {

				selected := 1
				countLand := 1

				if k1 != 0 {
					if l.NextInt(countLand) == 0 {
						selected = k1
					}
					countLand++
				}
				if l1 != 0 {
					if l.NextInt(countLand) == 0 {
						selected = l1
					}
					countLand++
				}
				if i2 != 0 {
					if l.NextInt(countLand) == 0 {
						selected = i2
					}
					countLand++
				}
				if j2 != 0 {
					if l.NextInt(countLand) == 0 {
						selected = j2
					}
					countLand++
				}

				if l.NextInt(3) == 0 {
					result[x+y*width] = selected
				} else if selected == 4 {
					result[x+y*width] = 4
				} else {
					result[x+y*width] = 0
				}
			}
		}
	}

	return result
}
