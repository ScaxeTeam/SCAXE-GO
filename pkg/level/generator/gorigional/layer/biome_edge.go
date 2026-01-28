package layer

type GenLayerBiomeEdge struct {
	*BaseLayer
}

func NewGenLayerBiomeEdge(baseSeed int64, parent GenLayer) *GenLayerBiomeEdge {
	l := &GenLayerBiomeEdge{
		BaseLayer: NewBaseLayer(baseSeed),
	}
	l.Parent = parent
	return l
}

func (l *GenLayerBiomeEdge) GetInts(x, z, width, depth int) []int {
	parentInts := l.Parent.GetInts(x-1, z-1, width+2, depth+2)
	output := make([]int, width*depth)
	parentWidth := width + 2

	for i := 0; i < depth; i++ {
		for j := 0; j < width; j++ {
			l.InitChunkSeed(int64(x+j), int64(z+i))
			k := parentInts[j+1+(i+1)*parentWidth]

			replaced := l.replaceBiomeEdgeIfNecessary(parentInts, output, j, i, width, k, BiomeExtremeHills, BiomeExtremeHillsEdge)
			if !replaced {

				replaced = l.replaceBiomeEdge(parentInts, output, j, i, width, k, BiomeMesaPlateau, BiomeMesa)
			}
			if !replaced {

				replaced = l.replaceBiomeEdge(parentInts, output, j, i, width, k, BiomeMesaPlateauF, BiomeMesa)
			}
			if !replaced {

				replaced = l.replaceBiomeEdge(parentInts, output, j, i, width, k, BiomeMegaTaiga, BiomeTaiga)
			}

			if !replaced {
				if k == BiomeDesert {

					l1 := parentInts[j+1+(i)*parentWidth]
					i2 := parentInts[j+1+(i+2)*parentWidth]
					j2 := parentInts[j+(i+1)*parentWidth]
					k2 := parentInts[j+2+(i+1)*parentWidth]

					if l1 != BiomeIcePlains && i2 != BiomeIcePlains && j2 != BiomeIcePlains && k2 != BiomeIcePlains {
						output[j+i*width] = k
					} else {
						output[j+i*width] = BiomeExtremeHillsPlus
					}

				} else if k == BiomeSwampland {

					l1 := parentInts[j+1+(i)*parentWidth]
					i2 := parentInts[j+1+(i+2)*parentWidth]
					j2 := parentInts[j+(i+1)*parentWidth]
					k2 := parentInts[j+2+(i+1)*parentWidth]

					isDesert := l1 == BiomeDesert || i2 == BiomeDesert || j2 == BiomeDesert || k2 == BiomeDesert
					isColdTaiga := l1 == BiomeColdTaiga || i2 == BiomeColdTaiga || j2 == BiomeColdTaiga || k2 == BiomeColdTaiga
					isIcePlains := l1 == BiomeIcePlains || i2 == BiomeIcePlains || j2 == BiomeIcePlains || k2 == BiomeIcePlains

					if !isDesert && !isColdTaiga && !isIcePlains {
						isJungle := l1 == BiomeJungle || i2 == BiomeJungle || j2 == BiomeJungle || k2 == BiomeJungle
						if !isJungle {
							output[j+i*width] = k
						} else {
							output[j+i*width] = BiomeJungleEdge
						}
					} else {
						output[j+i*width] = BiomePlains
					}
				} else {
					output[j+i*width] = k
				}
			}
		}
	}
	return output
}

func (l *GenLayerBiomeEdge) replaceBiomeEdgeIfNecessary(parentInts []int, output []int, x, z, width, currentBiome, targetBiome, edgeBiome int) bool {
	if !biomesEqualOrMesaPlateau(currentBiome, targetBiome) {
		return false
	}

	parentWidth := width + 2
	i := parentInts[x+1+(z)*parentWidth]
	j := parentInts[x+1+(z+2)*parentWidth]
	k := parentInts[x+(z+1)*parentWidth]
	m := parentInts[x+2+(z+1)*parentWidth]

	if l.canBiomesBeNeighbors(i, targetBiome) && l.canBiomesBeNeighbors(j, targetBiome) && l.canBiomesBeNeighbors(k, targetBiome) && l.canBiomesBeNeighbors(m, targetBiome) {
		output[x+z*width] = currentBiome
	} else {
		output[x+z*width] = edgeBiome
	}
	return true
}

func (l *GenLayerBiomeEdge) replaceBiomeEdge(parentInts []int, output []int, x, z, width, currentBiome, targetBiome, edgeBiome int) bool {
	if currentBiome != targetBiome {
		return false
	}

	parentWidth := width + 2
	i := parentInts[x+1+(z)*parentWidth]
	j := parentInts[x+1+(z+2)*parentWidth]
	k := parentInts[x+(z+1)*parentWidth]
	m := parentInts[x+2+(z+1)*parentWidth]

	if biomesEqualOrMesaPlateau(i, targetBiome) && biomesEqualOrMesaPlateau(j, targetBiome) && biomesEqualOrMesaPlateau(k, targetBiome) && biomesEqualOrMesaPlateau(m, targetBiome) {
		output[x+z*width] = currentBiome
	} else {
		output[x+z*width] = edgeBiome
	}
	return true
}

func (l *GenLayerBiomeEdge) canBiomesBeNeighbors(id1, id2 int) bool {
	if biomesEqualOrMesaPlateau(id1, id2) {
		return true
	}

	t1 := getTempCategory(id1)
	t2 := getTempCategory(id2)

	return t1 == t2 || t1 == TempCategoryMedium || t2 == TempCategoryMedium
}

func biomesEqualOrMesaPlateau(a, b int) bool {
	if a == b {
		return true
	}

	isMesaVariantA := a == BiomeMesaPlateau || a == BiomeMesaPlateauF
	isMesaVariantB := b == BiomeMesaPlateau || b == BiomeMesaPlateauF

	if !isMesaVariantA {

		return getBiomeClass(a) == getBiomeClass(b)
	} else {
		return isMesaVariantB
	}
}

type TempCategory int

const (
	TempCategoryOcean  TempCategory = 0
	TempCategoryCold   TempCategory = 1
	TempCategoryMedium TempCategory = 2
	TempCategoryWarm   TempCategory = 3
)

func getTempCategory(id int) TempCategory {

	switch id {
	case BiomeOcean, BiomeDeepOcean, BiomeRiver, BiomeSwampland, BiomeMushroomIsland, BiomeMushroomIslandShore, BiomeBeach:
		return TempCategoryMedium
	case BiomeFrozenOcean, BiomeFrozenRiver, BiomeIcePlains, BiomeIceMountains, BiomeColdBeach, BiomeColdTaiga, BiomeColdTaigaHills:
		return TempCategoryCold
	case BiomePlains, BiomeForest, BiomeExtremeHills, BiomeTaiga, BiomeExtremeHillsEdge, BiomeForestHills, BiomeTaigaHills, BiomeExtremeHillsPlus, BiomeBirchForest, BiomeBirchForestHills, BiomeRoofedForest, BiomeMegaTaiga, BiomeMegaTaigaHills, BiomeStoneBeach:
		return TempCategoryMedium
	case BiomeDesert, BiomeDesertHills, BiomeJungle, BiomeJungleHills, BiomeJungleEdge, BiomeSavanna, BiomeSavannaPlateau, BiomeMesa, BiomeMesaPlateauF, BiomeMesaPlateau, BiomeHell:
		return TempCategoryWarm
	}
	return TempCategoryMedium
}

func getBiomeClass(id int) int {

	switch id {
	case BiomeMesa, BiomeMesaPlateauF, BiomeMesaPlateau:
		return 1001
	case BiomeJungle, BiomeJungleHills, BiomeJungleEdge:
		return 1002
	case BiomeMegaTaiga, BiomeMegaTaigaHills, BiomeTaiga, BiomeTaigaHills, BiomeColdTaiga, BiomeColdTaigaHills:
		return 1003
	case BiomeOcean, BiomeDeepOcean, BiomeFrozenOcean:
		return 1004
	}
	return id
}
