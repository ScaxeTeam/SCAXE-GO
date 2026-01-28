package biomegrid

type ShoreMapLayer struct {
	*BaseMapLayer
}

func NewShoreMapLayer(seed int64, parent MapLayer) *ShoreMapLayer {
	return &ShoreMapLayer{
		BaseMapLayer: NewBaseMapLayer(seed, parent),
	}
}

func (l *ShoreMapLayer) GenerateValues(x, z, sizeX, sizeZ int) []int {

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

			if !isOcean(center) {
				north := parentValues[(i+1)+j*parentSizeX]
				south := parentValues[(i+1)+(j+2)*parentSizeX]
				east := parentValues[(i+2)+(j+1)*parentSizeX]
				west := parentValues[i+(j+1)*parentSizeX]

				if isOcean(north) || isOcean(south) || isOcean(east) || isOcean(west) {

					values[idx] = getShoreFor(center)
				} else {
					values[idx] = center
				}
			} else {
				values[idx] = center
			}
		}
	}

	return values
}

func isOcean(biome int) bool {
	return biome == BiomeOcean || biome == BiomeDeepOcean || biome == BiomeFrozenOcean
}

func getShoreFor(biome int) int {
	switch biome {
	case BiomeExtremeHills, BiomeIceMountains:
		return 25
	case BiomeIcePlains, BiomeTaiga:
		return 26
	case BiomeMushroomIsland:
		return 15
	case BiomeMesa, BiomeSavanna:
		return biome
	default:
		return BiomeBeach
	}
}
