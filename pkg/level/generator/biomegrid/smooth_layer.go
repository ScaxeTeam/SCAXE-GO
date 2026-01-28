package biomegrid

type SmoothMapLayer struct {
	*BaseMapLayer
}

func NewSmoothMapLayer(seed int64, parent MapLayer) *SmoothMapLayer {
	return &SmoothMapLayer{
		BaseMapLayer: NewBaseMapLayer(seed, parent),
	}
}

func (l *SmoothMapLayer) GenerateValues(x, z, sizeX, sizeZ int) []int {
	parentX := x - 1
	parentZ := z - 1
	parentSizeX := sizeX + 2
	parentSizeZ := sizeZ + 2
	parentValues := l.Parent.GenerateValues(parentX, parentZ, parentSizeX, parentSizeZ)

	values := make([]int, sizeX*sizeZ)

	for j := 0; j < sizeZ; j++ {
		for i := 0; i < sizeX; i++ {
			center := parentValues[(i+1)+(j+1)*parentSizeX]
			north := parentValues[(i+1)+j*parentSizeX]
			south := parentValues[(i+1)+(j+2)*parentSizeX]
			east := parentValues[(i+2)+(j+1)*parentSizeX]
			west := parentValues[i+(j+1)*parentSizeX]

			idx := i + j*sizeX
			l.SetCoordsSeed(x+i, z+j)

			if north == south && west == east {
				if l.NextInt(2) == 0 {
					values[idx] = north
				} else {
					values[idx] = west
				}
			} else if north == south {
				values[idx] = north
			} else if west == east {
				values[idx] = west
			} else {
				values[idx] = center
			}
		}
	}

	return values
}
