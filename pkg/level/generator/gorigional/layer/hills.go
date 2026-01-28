package layer

type GenLayerHills struct {
	*BaseLayer
	riverLayer GenLayer
}

func NewGenLayerHills(baseSeed int64, parent, riverLayer GenLayer) *GenLayerHills {
	l := &GenLayerHills{
		BaseLayer: NewBaseLayer(baseSeed),
	}
	l.Parent = parent
	l.riverLayer = riverLayer
	return l
}

func (l *GenLayerHills) InitWorldGenSeed(seed int64) {
	l.BaseLayer.InitWorldGenSeed(seed)

}

func (l *GenLayerHills) GetInts(x, z, width, depth int) []int {
	parentInts := l.Parent.GetInts(x-1, z-1, width+2, depth+2)
	riverInts := l.riverLayer.GetInts(x-1, z-1, width+2, depth+2)
	output := make([]int, width*depth)
	parentWidth := width + 2

	for i := 0; i < depth; i++ {
		for j := 0; j < width; j++ {

			l.InitChunkSeed(int64(x+j), int64(z+i))

			k := parentInts[j+1+(i+1)*parentWidth]
			riverVal := riverInts[j+1+(i+1)*parentWidth]

			flag := (riverVal-2)%29 == 0

			if k != 0 && riverVal >= 2 && (riverVal-2)%29 == 1 && k < 128 {

				mutatedID := getMutationForBiome(k)
				if mutatedID != k {
					k = mutatedID
				}
				output[j+i*width] = k
			} else if l.NextInt(3) != 0 && !flag {

				output[j+i*width] = k
			} else {
				i1 := k

				if k == BiomeDesert {
					i1 = BiomeDesertHills
				} else if k == BiomeForest {
					i1 = BiomeForestHills
				} else if k == BiomeBirchForest {
					i1 = BiomeBirchForestHills
				} else if k == BiomeRoofedForest {
					i1 = BiomePlains
				} else if k == BiomeTaiga {
					i1 = BiomeTaigaHills
				} else if k == BiomeMegaTaiga {
					i1 = BiomeMegaTaigaHills
				} else if k == BiomeColdTaiga {
					i1 = BiomeColdTaigaHills
				} else if k == BiomePlains {
					if l.NextInt(3) == 0 {
						i1 = BiomeForestHills
					} else {
						i1 = BiomeForest
					}
				} else if k == BiomeIcePlains {
					i1 = BiomeIceMountains
				} else if k == BiomeJungle {
					i1 = BiomeJungleHills
				} else if k == BiomeOcean {
					i1 = BiomeDeepOcean
				} else if k == BiomeExtremeHills {
					i1 = BiomeExtremeHillsPlus
				} else if k == BiomeSavanna {
					i1 = BiomeSavannaPlateau
				} else if biomesEqualOrMesaPlateau(k, BiomeMesaPlateauF) {
					i1 = BiomeMesa
				} else if k == BiomeDeepOcean && l.NextInt(3) == 0 {

					if l.NextInt(2) == 0 {
						i1 = BiomePlains
					} else {
						i1 = BiomeForest
					}
				}

				if flag && i1 != k {
					mutatedOfVariant := getMutationForBiome(i1)
					if mutatedOfVariant != i1 {
						i1 = mutatedOfVariant
					} else {
						i1 = k
					}
				}

				if i1 == k {
					output[j+i*width] = k
				} else {

					k2 := parentInts[j+1+(i+0)*parentWidth]
					j1_neighbor := parentInts[j+2+(i+1)*parentWidth]
					k1 := parentInts[j+0+(i+1)*parentWidth]
					l1 := parentInts[j+1+(i+2)*parentWidth]

					countCompatible := 0
					if biomesEqualOrMesaPlateau(k2, k) {
						countCompatible++
					}
					if biomesEqualOrMesaPlateau(j1_neighbor, k) {
						countCompatible++
					}
					if biomesEqualOrMesaPlateau(k1, k) {
						countCompatible++
					}
					if biomesEqualOrMesaPlateau(l1, k) {
						countCompatible++
					}

					if countCompatible >= 3 {
						output[j+i*width] = i1
					} else {
						output[j+i*width] = k
					}
				}
			}

		}
	}
	return output
}

func getMutationForBiome(id int) int {

	switch id {
	case BiomePlains, BiomeDesert, BiomeExtremeHills, BiomeForest, BiomeTaiga, BiomeSwampland,
		BiomeIcePlains, BiomeJungle, BiomeJungleEdge, BiomeBirchForest, BiomeBirchForestHills,
		BiomeRoofedForest, BiomeColdTaiga, BiomeMegaTaiga, BiomeMegaTaigaHills,
		BiomeExtremeHillsPlus, BiomeSavanna, BiomeSavannaPlateau, BiomeMesa,
		BiomeMesaPlateauF, BiomeMesaPlateau:
		return id + 128
	}

	return id
}
