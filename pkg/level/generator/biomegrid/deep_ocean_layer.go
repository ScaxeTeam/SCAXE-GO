package biomegrid

type DeepOceanMapLayer struct {
	*BaseMapLayer
}

func NewDeepOceanMapLayer(seed int64, parent MapLayer) *DeepOceanMapLayer {
	return &DeepOceanMapLayer{
		BaseMapLayer: NewBaseMapLayer(seed, parent),
	}
}

func (l *DeepOceanMapLayer) GenerateValues(x, z, sizeX, sizeZ int) []int {

	parentX := x - 1
	parentZ := z - 1
	parentSizeX := sizeX + 2
	parentSizeZ := sizeZ + 2
	parentValues := l.Parent.GenerateValues(parentX, parentZ, parentSizeX, parentSizeZ)

	values := make([]int, sizeX*sizeZ)

	for j := 0; j < sizeZ; j++ {
		for i := 0; i < sizeX; i++ {
			center := parentValues[(i+1)+(j+1)*parentSizeX]
			idx := i + j*sizeX

			if center == 0 {

				north := parentValues[(i+1)+j*parentSizeX]
				south := parentValues[(i+1)+(j+2)*parentSizeX]
				east := parentValues[(i+2)+(j+1)*parentSizeX]
				west := parentValues[i+(j+1)*parentSizeX]

				if north == 0 && south == 0 && east == 0 && west == 0 {
					l.SetCoordsSeed(x+i, z+j)
					if l.NextInt(2) == 0 {
						values[idx] = 24
					} else {
						values[idx] = 0
					}
				} else {
					values[idx] = 0
				}
			} else {
				values[idx] = center
			}
		}
	}

	return values
}
