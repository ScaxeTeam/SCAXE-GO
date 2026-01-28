package biomegrid

type RiverMapLayer struct {
	*BaseMapLayer
	Merged MapLayer
}

func NewRiverMapLayer(seed int64, parent MapLayer) *RiverMapLayer {
	return &RiverMapLayer{
		BaseMapLayer: NewBaseMapLayer(seed, parent),
		Merged:       nil,
	}
}

func NewRiverMapLayerMerged(seed int64, riverParent, biomeParent MapLayer) *RiverMapLayer {
	return &RiverMapLayer{
		BaseMapLayer: NewBaseMapLayer(seed, riverParent),
		Merged:       biomeParent,
	}
}

func (l *RiverMapLayer) GenerateValues(x, z, sizeX, sizeZ int) []int {
	if l.Merged != nil {
		return l.generateMerged(x, z, sizeX, sizeZ)
	}
	return l.generateRivers(x, z, sizeX, sizeZ)
}

func (l *RiverMapLayer) generateRivers(x, z, sizeX, sizeZ int) []int {
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

			centerVal := riverVal(center)
			if centerVal == riverVal(north) &&
				centerVal == riverVal(south) &&
				centerVal == riverVal(east) &&
				centerVal == riverVal(west) {
				values[idx] = -1
			} else {
				values[idx] = BiomeRiver
			}
		}
	}

	return values
}

func (l *RiverMapLayer) generateMerged(x, z, sizeX, sizeZ int) []int {
	riverValues := l.Parent.GenerateValues(x, z, sizeX, sizeZ)
	biomeValues := l.Merged.GenerateValues(x, z, sizeX, sizeZ)

	values := make([]int, sizeX*sizeZ)

	for i := 0; i < sizeX*sizeZ; i++ {
		if riverValues[i] == BiomeRiver && !isOcean(biomeValues[i]) {

			if biomeValues[i] == BiomeIcePlains || biomeValues[i] == BiomeTaiga {
				values[i] = BiomeFrozenRiver
			} else if biomeValues[i] == BiomeMushroomIsland {
				values[i] = 15
			} else {
				values[i] = BiomeRiver
			}
		} else {
			values[i] = biomeValues[i]
		}
	}

	return values
}

func riverVal(biome int) int {
	if biome >= 2 {
		return 2 + (biome & 1)
	}
	return biome
}
