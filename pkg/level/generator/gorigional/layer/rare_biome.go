package layer

type GenLayerRareBiome struct {
	*BaseLayer
}

func NewGenLayerRareBiome(baseSeed int64, parent GenLayer) *GenLayerRareBiome {
	l := &GenLayerRareBiome{NewBaseLayer(baseSeed)}
	l.BaseLayer.Parent = parent
	return l
}

func (l *GenLayerRareBiome) GetInts(x, z, width, depth int) []int {

	xOff := x - 1
	zOff := z - 1
	wOff := width + 2
	dOff := depth + 2

	parentInts := l.Parent.GetInts(xOff, zOff, wOff, dOff)
	out := make([]int, width*depth)

	for dz := 0; dz < depth; dz++ {
		for dx := 0; dx < width; dx++ {
			l.InitChunkSeed(int64(dx+x), int64(dz+z))

			centerID := parentInts[(dx+1)+(dz+1)*wOff]

			if l.NextInt(57) == 0 {
				if centerID == 1 {
					out[dx+dz*width] = 129
				} else {
					out[dx+dz*width] = centerID
				}
			} else {
				out[dx+dz*width] = centerID
			}
		}
	}
	return out
}
