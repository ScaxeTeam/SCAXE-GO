package biomegrid

type ClimateType int

const (
	ClimateWarmWet ClimateType = iota
	ClimateColdDry
	ClimateLargerBiomes
)

type WhittakerMapLayer struct {
	*BaseMapLayer
	Climate ClimateType
}

func NewWhittakerMapLayer(seed int64, parent MapLayer, climate ClimateType) *WhittakerMapLayer {
	return &WhittakerMapLayer{
		BaseMapLayer: NewBaseMapLayer(seed, parent),
		Climate:      climate,
	}
}

func (l *WhittakerMapLayer) GenerateValues(x, z, sizeX, sizeZ int) []int {
	parentValues := l.Parent.GenerateValues(x, z, sizeX, sizeZ)
	values := make([]int, sizeX*sizeZ)

	for j := 0; j < sizeZ; j++ {
		for i := 0; i < sizeX; i++ {
			idx := i + j*sizeX
			val := parentValues[idx]
			l.SetCoordsSeed(x+i, z+j)

			if val != 0 {
				switch l.Climate {
				case ClimateWarmWet:

					if l.NextInt(6) == 0 {
						val = 2 + l.NextInt(2)
					}
				case ClimateColdDry:

					if l.NextInt(6) == 0 {
						val = 4 + l.NextInt(2)
					}
				case ClimateLargerBiomes:

					if l.NextInt(3) == 0 && val == 1 {
						val = 1
					}
				}
			}

			values[idx] = val
		}
	}

	return values
}
