package layer

type GenLayerAddMushroomIsland struct {
	*BaseLayer
}

func NewGenLayerAddMushroomIsland(baseSeed int64, parent GenLayer) *GenLayerAddMushroomIsland {
	l := &GenLayerAddMushroomIsland{NewBaseLayer(baseSeed)}
	l.BaseLayer.Parent = parent
	return l
}

func (l *GenLayerAddMushroomIsland) GetInts(x, z, width, depth int) []int {

	xOff := x - 1
	zOff := z - 1
	wOff := width + 2
	dOff := depth + 2

	parentInts := l.Parent.GetInts(xOff, zOff, wOff, dOff)
	out := make([]int, width*depth)

	for dz := 0; dz < depth; dz++ {
		for dx := 0; dx < width; dx++ {

			k1 := parentInts[(dx+0)+(dz+0)*wOff]
			l1 := parentInts[(dx+2)+(dz+0)*wOff]
			i2 := parentInts[(dx+0)+(dz+2)*wOff]
			j2 := parentInts[(dx+2)+(dz+2)*wOff]
			k2 := parentInts[(dx+1)+(dz+1)*wOff]

			l.InitChunkSeed(int64(dx+x), int64(dz+z))

			if k2 == 0 && k1 == 0 && l1 == 0 && i2 == 0 && j2 == 0 && l.NextInt(100) == 0 {
				out[dx+dz*width] = 14
			} else {
				out[dx+dz*width] = k2
			}
		}
	}
	return out
}
