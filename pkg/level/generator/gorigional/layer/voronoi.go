package layer

type GenLayerVoronoiZoom struct {
	*BaseLayer
}

func NewGenLayerVoronoiZoom(baseSeed int64, parent GenLayer) *GenLayerVoronoiZoom {
	l := &GenLayerVoronoiZoom{NewBaseLayer(baseSeed)}
	l.BaseLayer.Parent = parent
	return l
}

func (l *GenLayerVoronoiZoom) GetInts(x, z, width, depth int) []int {
	x -= 2
	z -= 2

	i := x >> 2
	j := z >> 2
	k := (width >> 2) + 2
	m := (depth >> 2) + 2

	parentInts := l.Parent.GetInts(i, j, k, m)

	out := make([]int, width*depth)

	for k1 := 0; k1 < m-1; k1++ {
		l1 := 0
		i2 := parentInts[l1+(k1+0)*k]

		for l1 = 0; l1 < k-1; l1++ {

			l.InitChunkSeed(int64((l1+i)<<2), int64((k1+j)<<2))
			d1 := (float64(l.NextInt(1024))/1024.0 - 0.5) * 3.6
			d2 := (float64(l.NextInt(1024))/1024.0 - 0.5) * 3.6

			l.InitChunkSeed(int64((l1+i+1)<<2), int64((k1+j)<<2))
			d3 := (float64(l.NextInt(1024))/1024.0-0.5)*3.6 + 4.0
			d4 := (float64(l.NextInt(1024))/1024.0 - 0.5) * 3.6

			l.InitChunkSeed(int64((l1+i)<<2), int64((k1+j+1)<<2))
			d5 := (float64(l.NextInt(1024))/1024.0 - 0.5) * 3.6
			d6 := (float64(l.NextInt(1024))/1024.0-0.5)*3.6 + 4.0

			l.InitChunkSeed(int64((l1+i+1)<<2), int64((k1+j+1)<<2))
			d7 := (float64(l.NextInt(1024))/1024.0-0.5)*3.6 + 4.0
			d8 := (float64(l.NextInt(1024))/1024.0-0.5)*3.6 + 4.0

			k2 := parentInts[l1+1+(k1+0)*k] & 255
			l2 := parentInts[l1+1+(k1+1)*k] & 255

			j2 := parentInts[l1+(k1+1)*k] & 255

			j2 = parentInts[l1+(k1+1)*k] & 255
			i2 = parentInts[l1+(k1+0)*k] & 255

			for i3 := 0; i3 < 4; i3++ {

				absZ := (k1+j)<<2 + i3
				outZ := absZ - z

				for k3 := 0; k3 < 4; k3++ {

					absX := (l1+i)<<2 + k3
					outX := absX - x

					if outX >= 0 && outX < width && outZ >= 0 && outZ < depth {

						d9 := (float64(i3)-d2)*(float64(i3)-d2) + (float64(k3)-d1)*(float64(k3)-d1)
						d10 := (float64(i3)-d4)*(float64(i3)-d4) + (float64(k3)-d3)*(float64(k3)-d3)
						d11 := (float64(i3)-d6)*(float64(i3)-d6) + (float64(k3)-d5)*(float64(k3)-d5)
						d12 := (float64(i3)-d8)*(float64(i3)-d8) + (float64(k3)-d7)*(float64(k3)-d7)

						idx := outX + outZ*width

						if d9 < d10 && d9 < d11 && d9 < d12 {
							out[idx] = i2
						} else if d10 < d9 && d10 < d11 && d10 < d12 {
							out[idx] = k2
						} else if d11 < d9 && d11 < d10 && d11 < d12 {
							out[idx] = j2
						} else {
							out[idx] = l2
						}
					}
				}
			}

		}
	}

	return out
}
