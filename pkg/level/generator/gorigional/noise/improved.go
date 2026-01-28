package noise

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type ImprovedNoise struct {
	permutations [512]int
	xCoord       float64
	yCoord       float64
	zCoord       float64
}

var gradX = []float64{1.0, -1.0, 1.0, -1.0, 1.0, -1.0, 1.0, -1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, -1.0, 0.0}
var gradY = []float64{1.0, 1.0, -1.0, -1.0, 0.0, 0.0, 0.0, 0.0, 1.0, -1.0, 1.0, -1.0, 1.0, -1.0, 1.0, -1.0}
var gradZ = []float64{0.0, 0.0, 0.0, 0.0, 1.0, 1.0, -1.0, -1.0, 1.0, 1.0, -1.0, -1.0, 0.0, 1.0, 0.0, -1.0}
var grad2X = []float64{1.0, -1.0, 1.0, -1.0, 1.0, -1.0, 1.0, -1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, -1.0, 0.0}
var grad2Z = []float64{0.0, 0.0, 0.0, 0.0, 1.0, 1.0, -1.0, -1.0, 1.0, 1.0, -1.0, -1.0, 0.0, 1.0, 0.0, -1.0}

func NewImprovedNoise(rnd *rand.Random) *ImprovedNoise {
	n := &ImprovedNoise{}

	n.xCoord = rnd.NextDouble() * 256.0
	n.yCoord = rnd.NextDouble() * 256.0
	n.zCoord = rnd.NextDouble() * 256.0

	for i := 0; i < 256; i++ {
		n.permutations[i] = i
	}

	for l := 0; l < 256; l++ {

		j := rnd.NextBoundedInt(256-l) + l
		k := n.permutations[l]
		n.permutations[l] = n.permutations[j]
		n.permutations[j] = k

		n.permutations[l+256] = n.permutations[l]
	}

	return n
}

func (n *ImprovedNoise) GetXCoord() float64     { return n.xCoord }
func (n *ImprovedNoise) GetYCoord() float64     { return n.yCoord }
func (n *ImprovedNoise) GetZCoord() float64     { return n.zCoord }
func (n *ImprovedNoise) GetPermutations() []int { return n.permutations[:] }

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func grad3D(hash int, x, y, z float64) float64 {
	i := hash & 15
	return gradX[i]*x + gradY[i]*y + gradZ[i]*z
}

func grad2D(hash int, x, z float64) float64 {
	i := hash & 15
	return grad2X[i]*x + grad2Z[i]*z
}

func (n *ImprovedNoise) PopulateNoiseArray(noiseArray []float64, xOffset, yOffset, zOffset float64, xSize, ySize, zSize int, xScale, yScale, zScale, noiseScale float64) {
	if ySize == 1 {

		invScale := 1.0 / noiseScale
		idx := 0

		for j2 := 0; j2 < xSize; j2++ {
			d17 := xOffset + float64(j2)*xScale + n.xCoord
			i6 := int(math.Floor(d17))
			if d17 < float64(i6) {
				i6--
			}
			k2 := i6 & 255
			d17 = d17 - float64(i6)
			d18 := d17 * d17 * d17 * (d17*(d17*6.0-15.0) + 10.0)

			for j6 := 0; j6 < zSize; j6++ {
				d19 := zOffset + float64(j6)*zScale + n.zCoord
				k6 := int(math.Floor(d19))
				if d19 < float64(k6) {
					k6--
				}
				l6 := k6 & 255
				d19 = d19 - float64(k6)
				d20 := d19 * d19 * d19 * (d19*(d19*6.0-15.0) + 10.0)

				i5 := n.permutations[k2] + 0
				j5 := n.permutations[i5] + l6
				j := n.permutations[k2+1] + 0
				k5 := n.permutations[j] + l6

				d14 := lerp(d18, grad2D(n.permutations[j5], d17, d19), grad3D(n.permutations[k5], d17-1.0, 0.0, d19))
				d15 := lerp(d18, grad3D(n.permutations[j5+1], d17, 0.0, d19-1.0), grad3D(n.permutations[k5+1], d17-1.0, 0.0, d19-1.0))
				d21 := lerp(d20, d14, d15)
				noiseArray[idx] += d21 * invScale
				idx++
			}
		}
		return
	}

	invScale := 1.0 / noiseScale
	k := -1

	var d1, d2, d3, d4 float64

	idx := 0

	for l2 := 0; l2 < xSize; l2++ {
		d5 := xOffset + float64(l2)*xScale + n.xCoord
		i3 := int(math.Floor(d5))
		if d5 < float64(i3) {
			i3--
		}
		j3 := i3 & 255
		d5 = d5 - float64(i3)
		d6 := d5 * d5 * d5 * (d5*(d5*6.0-15.0) + 10.0)

		for k3 := 0; k3 < zSize; k3++ {
			d7 := zOffset + float64(k3)*zScale + n.zCoord
			l3 := int(math.Floor(d7))
			if d7 < float64(l3) {
				l3--
			}
			i4 := l3 & 255
			d7 = d7 - float64(l3)
			d8 := d7 * d7 * d7 * (d7*(d7*6.0-15.0) + 10.0)

			for j4 := 0; j4 < ySize; j4++ {
				d9 := yOffset + float64(j4)*yScale + n.yCoord
				k4 := int(math.Floor(d9))
				if d9 < float64(k4) {
					k4--
				}
				l4 := k4 & 255
				d9 = d9 - float64(k4)
				d10 := d9 * d9 * d9 * (d9*(d9*6.0-15.0) + 10.0)

				if j4 == 0 || l4 != k {
					k = l4
					l := n.permutations[j3] + l4
					i1 := n.permutations[l] + i4
					j1 := n.permutations[l+1] + i4
					k1 := n.permutations[j3+1] + l4
					l1 := n.permutations[k1] + i4
					i2 := n.permutations[k1+1] + i4

					g1 := grad3D(n.permutations[i1], d5, d9, d7)
					g2 := grad3D(n.permutations[l1], d5-1.0, d9, d7)
					d1 = lerp(d6, g1, g2)

					g3 := grad3D(n.permutations[j1], d5, d9-1.0, d7)
					g4 := grad3D(n.permutations[i2], d5-1.0, d9-1.0, d7)
					d2 = lerp(d6, g3, g4)

					g5 := grad3D(n.permutations[i1+1], d5, d9, d7-1.0)
					g6 := grad3D(n.permutations[l1+1], d5-1.0, d9, d7-1.0)
					d3 = lerp(d6, g5, g6)

					g7 := grad3D(n.permutations[j1+1], d5, d9-1.0, d7-1.0)
					g8 := grad3D(n.permutations[i2+1], d5-1.0, d9-1.0, d7-1.0)
					d4 = lerp(d6, g7, g8)
				}

				d11 := lerp(d10, d1, d2)
				d12 := lerp(d10, d3, d4)
				d13 := lerp(d8, d11, d12)

				noiseArray[idx] += d13 * invScale
				idx++
			}
		}
	}
}
