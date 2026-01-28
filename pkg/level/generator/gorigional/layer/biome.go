package layer

const (
	BiomeOcean               = 0
	BiomePlains              = 1
	BiomeDesert              = 2
	BiomeExtremeHills        = 3
	BiomeForest              = 4
	BiomeTaiga               = 5
	BiomeSwampland           = 6
	BiomeRiver               = 7
	BiomeHell                = 8
	BiomeSky                 = 9
	BiomeFrozenOcean         = 10
	BiomeFrozenRiver         = 11
	BiomeIcePlains           = 12
	BiomeIceMountains        = 13
	BiomeMushroomIsland      = 14
	BiomeMushroomIslandShore = 15
	BiomeBeach               = 16
	BiomeDesertHills         = 17
	BiomeForestHills         = 18
	BiomeTaigaHills          = 19
	BiomeExtremeHillsEdge    = 20
	BiomeJungle              = 21
	BiomeJungleHills         = 22
	BiomeJungleEdge          = 23
	BiomeDeepOcean           = 24
	BiomeStoneBeach          = 25
	BiomeColdBeach           = 26
	BiomeBirchForest         = 27
	BiomeBirchForestHills    = 28
	BiomeRoofedForest        = 29
	BiomeColdTaiga           = 30
	BiomeColdTaigaHills      = 31
	BiomeMegaTaiga           = 32
	BiomeMegaTaigaHills      = 33
	BiomeExtremeHillsPlus    = 34
	BiomeSavanna             = 35
	BiomeSavannaPlateau      = 36
	BiomeMesa                = 37
	BiomeMesaPlateauF        = 38
	BiomeMesaPlateau         = 39

	BiomeMesaRock      = 37
	BiomeMesaClearRock = 37
)

type GenLayerBiome struct {
	*BaseLayer
	warmBiomes   []int
	mediumBiomes []int
	coldBiomes   []int
	iceBiomes    []int
}

func NewGenLayerBiome(seed int64, parent GenLayer) *GenLayerBiome {
	l := &GenLayerBiome{
		BaseLayer: NewBaseLayer(seed),
	}
	l.Parent = parent

	l.warmBiomes = []int{BiomeDesert, BiomeDesert, BiomeDesert, BiomeSavanna, BiomeSavanna, BiomePlains}
	l.mediumBiomes = []int{BiomeForest, BiomeRoofedForest, BiomeExtremeHills, BiomePlains, BiomeBirchForest, BiomeSwampland}
	l.coldBiomes = []int{BiomeForest, BiomeExtremeHills, BiomeTaiga, BiomePlains}
	l.iceBiomes = []int{BiomeIcePlains, BiomeIcePlains, BiomeIcePlains, BiomeColdTaiga}

	return l
}

func (l *GenLayerBiome) GetInts(areaX, areaY, width, height int) []int {
	parentInts := l.Parent.GetInts(areaX, areaY, width, height)
	result := make([]int, width*height)

	for i := 0; i < len(parentInts); i++ {
		l.InitChunkSeed(int64(i%width+areaX), int64(i/width+areaY))

		k := parentInts[i]

		val := (k & 3840) >> 8
		k = k & ^3840

		if isOceanic(k) {
			result[i] = k
		} else if k == BiomeMushroomIsland {
			result[i] = k
		} else if k == 1 {
			if val > 0 {
				if l.NextInt(3) == 0 {
					result[i] = BiomeMesaPlateau
				} else {
					result[i] = BiomeMesaPlateauF
				}
			} else {
				result[i] = l.warmBiomes[l.NextInt(len(l.warmBiomes))]
			}
		} else if k == 2 {
			if val > 0 {
				result[i] = BiomeJungle
			} else {
				result[i] = l.mediumBiomes[l.NextInt(len(l.mediumBiomes))]
			}
		} else if k == 3 {
			if val > 0 {
				result[i] = BiomeMegaTaiga
			} else {
				result[i] = l.coldBiomes[l.NextInt(len(l.coldBiomes))]
			}
		} else if k == 4 {
			result[i] = l.iceBiomes[l.NextInt(len(l.iceBiomes))]
		} else {
			result[i] = BiomeMushroomIsland
		}
	}

	return result
}

func isOceanic(id int) bool {

	return id == BiomeOcean || id == BiomeDeepOcean || id == BiomeFrozenOcean
}
