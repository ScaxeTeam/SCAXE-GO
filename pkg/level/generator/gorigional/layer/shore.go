package layer

type GenLayerShore struct {
	*BaseLayer
}

func NewGenLayerShore(baseSeed int64, parent GenLayer) *GenLayerShore {
	l := &GenLayerShore{
		BaseLayer: NewBaseLayer(baseSeed),
	}
	l.Parent = parent
	return l
}

func (l *GenLayerShore) GetInts(x, z, width, depth int) []int {
	parentInts := l.Parent.GetInts(x-1, z-1, width+2, depth+2)
	output := make([]int, width*depth)
	parentWidth := width + 2

	for i := 0; i < depth; i++ {
		for j := 0; j < width; j++ {
			l.InitChunkSeed(int64(x+j), int64(z+i))
			k := parentInts[j+1+(i+1)*parentWidth]

			if k == BiomeMushroomIsland {
				j2 := parentInts[j+1+(i+0)*parentWidth]
				i3 := parentInts[j+2+(i+1)*parentWidth]
				l3 := parentInts[j+0+(i+1)*parentWidth]
				k4 := parentInts[j+1+(i+2)*parentWidth]

				if j2 != BiomeOcean && i3 != BiomeOcean && l3 != BiomeOcean && k4 != BiomeOcean {
					output[j+i*width] = k
				} else {
					output[j+i*width] = BiomeMushroomIslandShore
				}
			} else if isJungle(k) {
				i2 := parentInts[j+1+(i+0)*parentWidth]
				l2 := parentInts[j+2+(i+1)*parentWidth]
				k3 := parentInts[j+0+(i+1)*parentWidth]
				j4 := parentInts[j+1+(i+2)*parentWidth]

				if l.isJungleCompatible(i2) && l.isJungleCompatible(l2) && l.isJungleCompatible(k3) && l.isJungleCompatible(j4) {
					if !isOceanic(i2) && !isOceanic(l2) && !isOceanic(k3) && !isOceanic(j4) {
						output[j+i*width] = k
					} else {
						output[j+i*width] = BiomeBeach
					}
				} else {
					output[j+i*width] = BiomeJungleEdge
				}
			} else if k != BiomeExtremeHills && k != BiomeExtremeHillsPlus && k != BiomeExtremeHillsEdge {
				if isSnowy(k) {
					l.replaceIfNeighborOcean(parentInts, output, j, i, parentWidth, k, BiomeColdBeach, width)
				} else if k != BiomeMesa && k != BiomeMesaPlateauF {
					if k != BiomeOcean && k != BiomeDeepOcean && k != BiomeRiver && k != BiomeSwampland {
						l1 := parentInts[j+1+(i+0)*parentWidth]
						k2 := parentInts[j+2+(i+1)*parentWidth]
						j3 := parentInts[j+0+(i+1)*parentWidth]
						i4 := parentInts[j+1+(i+2)*parentWidth]

						if !isOceanic(l1) && !isOceanic(k2) && !isOceanic(j3) && !isOceanic(i4) {
							output[j+i*width] = k
						} else {
							output[j+i*width] = BiomeBeach
						}
					} else {
						output[j+i*width] = k
					}
				} else {

					l_val := parentInts[j+1+(i+0)*parentWidth]
					i1 := parentInts[j+2+(i+1)*parentWidth]
					j1 := parentInts[j+0+(i+1)*parentWidth]
					k1 := parentInts[j+1+(i+2)*parentWidth]

					if !isOceanic(l_val) && !isOceanic(i1) && !isOceanic(j1) && !isOceanic(k1) {
						if isMesa(l_val) && isMesa(i1) && isMesa(j1) && isMesa(k1) {
							output[j+i*width] = k
						} else {
							output[j+i*width] = BiomeDesert
						}
					} else {
						output[j+i*width] = k
					}
				}
			} else {
				l.replaceIfNeighborOcean(parentInts, output, j, i, parentWidth, k, BiomeStoneBeach, width)
			}
		}
	}
	return output
}

func (l *GenLayerShore) replaceIfNeighborOcean(parentInts []int, output []int, x, z, parentWidth, current, replace, outWidth int) {
	if isOceanic(current) {
		output[x+z*outWidth] = current
	} else {
		i := parentInts[x+1+(z+0)*parentWidth]
		j := parentInts[x+2+(z+1)*parentWidth]
		k := parentInts[x+0+(z+1)*parentWidth]
		m := parentInts[x+1+(z+2)*parentWidth]

		if !isOceanic(i) && !isOceanic(j) && !isOceanic(k) && !isOceanic(m) {
			output[x+z*outWidth] = current
		} else {
			output[x+z*outWidth] = replace
		}
	}
}

func (l *GenLayerShore) isJungleCompatible(id int) bool {
	if isJungle(id) {
		return true
	}
	return id == BiomeJungleEdge || id == BiomeJungle || id == BiomeJungleHills || id == BiomeForest || id == BiomeTaiga || isOceanic(id)
}

func isJungle(id int) bool {
	return id == BiomeJungle || id == BiomeJungleHills || id == BiomeJungleEdge
}

func isMesa(id int) bool {
	return id == BiomeMesa || id == BiomeMesaPlateauF || id == BiomeMesaPlateau
}

func isSnowy(id int) bool {
	return id == BiomeIcePlains || id == BiomeIceMountains || id == BiomeColdTaiga || id == BiomeColdTaigaHills
}
