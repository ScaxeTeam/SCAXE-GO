package biomegrid

type ErosionMapLayer struct {
	*BaseMapLayer
}

func NewErosionMapLayer(seed int64, parent MapLayer) *ErosionMapLayer {
	return &ErosionMapLayer{
		BaseMapLayer: NewBaseMapLayer(seed, parent),
	}
}

func (l *ErosionMapLayer) GenerateValues(x, z, sizeX, sizeZ int) []int {

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

			l.SetCoordsSeed(x+i, z+j)
			idx := i + j*sizeX

			if center != 0 {
				oceanCount := 0
				if north == 0 {
					oceanCount++
				}
				if south == 0 {
					oceanCount++
				}
				if east == 0 {
					oceanCount++
				}
				if west == 0 {
					oceanCount++
				}

				if oceanCount >= 3 && l.NextInt(2) == 0 {
					values[idx] = 0
				} else {
					values[idx] = center
				}
			} else {

				landCount := 0
				if north != 0 {
					landCount++
				}
				if south != 0 {
					landCount++
				}
				if east != 0 {
					landCount++
				}
				if west != 0 {
					landCount++
				}

				if landCount >= 3 && l.NextInt(3) == 0 {
					values[idx] = 1
				} else {
					values[idx] = 0
				}
			}
		}
	}

	return values
}
