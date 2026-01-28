package biomegrid

const (
	BiomeOcean          = 0
	BiomePlains         = 1
	BiomeDesert         = 2
	BiomeExtremeHills   = 3
	BiomeForest         = 4
	BiomeTaiga          = 5
	BiomeSwampland      = 6
	BiomeRiver          = 7
	BiomeFrozenOcean    = 10
	BiomeFrozenRiver    = 11
	BiomeIcePlains      = 12
	BiomeIceMountains   = 13
	BiomeMushroomIsland = 14
	BiomeBeach          = 16
	BiomeDesertHills    = 17
	BiomeForestHills    = 18
	BiomeTaigaHills     = 19
	BiomeJungle         = 21
	BiomeJungleHills    = 22
	BiomeJungleEdge     = 23
	BiomeDeepOcean      = 24
	BiomeSavanna        = 35
	BiomeSavannaPlateau = 36
	BiomeMesa           = 37
)

type BiomeMapLayer struct {
	*BaseMapLayer
}

func NewBiomeMapLayer(seed int64, parent MapLayer) *BiomeMapLayer {
	return &BiomeMapLayer{
		BaseMapLayer: NewBaseMapLayer(seed, parent),
	}
}

var (
	warmBiomes    = []int{BiomeDesert, BiomeSavanna, BiomePlains}
	wetBiomes     = []int{BiomeForest, BiomeSwampland, BiomePlains}
	coldBiomes    = []int{BiomeTaiga, BiomeIcePlains, BiomeExtremeHills}
	dryBiomes     = []int{BiomeDesert, BiomeSavanna, BiomeMesa}
	defaultBiomes = []int{BiomePlains, BiomeForest, BiomeTaiga, BiomeExtremeHills}
)

func (l *BiomeMapLayer) GenerateValues(x, z, sizeX, sizeZ int) []int {
	parentValues := l.Parent.GenerateValues(x, z, sizeX, sizeZ)
	values := make([]int, sizeX*sizeZ)

	for j := 0; j < sizeZ; j++ {
		for i := 0; i < sizeX; i++ {
			idx := i + j*sizeX
			val := parentValues[idx]
			l.SetCoordsSeed(x+i, z+j)

			switch val {
			case 0:
				values[idx] = BiomeOcean
			case 1:
				values[idx] = defaultBiomes[l.NextInt(len(defaultBiomes))]
			case 2:
				values[idx] = warmBiomes[l.NextInt(len(warmBiomes))]
			case 3:
				values[idx] = wetBiomes[l.NextInt(len(wetBiomes))]
			case 4:
				values[idx] = coldBiomes[l.NextInt(len(coldBiomes))]
			case 5:
				values[idx] = dryBiomes[l.NextInt(len(dryBiomes))]
			case 24:
				values[idx] = BiomeDeepOcean
			default:

				if val > 10 {
					values[idx] = val
				} else {
					values[idx] = defaultBiomes[l.NextInt(len(defaultBiomes))]
				}
			}
		}
	}

	return values
}
