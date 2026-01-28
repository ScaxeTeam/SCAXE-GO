package biomegrid

type ConstantBiomeMapLayer struct {
	*BaseMapLayer
	BiomeID int
}

func NewConstantBiomeMapLayer(seed int64, biomeID int) *ConstantBiomeMapLayer {
	return &ConstantBiomeMapLayer{
		BaseMapLayer: NewBaseMapLayer(seed, nil),
		BiomeID:      biomeID,
	}
}

func (l *ConstantBiomeMapLayer) GenerateValues(x, z, sizeX, sizeZ int) []int {
	values := make([]int, sizeX*sizeZ)
	for i := range values {
		values[i] = l.BiomeID
	}
	return values
}
