package noise

import (
	"math"
	"testing"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

func TestImprovedNoise_Coords_JavaParity(t *testing.T) {

	expectedXCoord := 92.62159543308078000
	expectedYCoord := 238.84633223386650000
	expectedZCoord := 213.27138533658206000

	r := rand.NewRandom(12345)
	n := NewImprovedNoise(r)

	tolerance := 1e-14

	if math.Abs(n.GetXCoord()-expectedXCoord) > tolerance {
		t.Errorf("xCoord mismatch: got %.17f, expected %.17f", n.GetXCoord(), expectedXCoord)
	}
	if math.Abs(n.GetYCoord()-expectedYCoord) > tolerance {
		t.Errorf("yCoord mismatch: got %.17f, expected %.17f", n.GetYCoord(), expectedYCoord)
	}
	if math.Abs(n.GetZCoord()-expectedZCoord) > tolerance {
		t.Errorf("zCoord mismatch: got %.17f, expected %.17f", n.GetZCoord(), expectedZCoord)
	}

	t.Logf("✓ ImprovedNoise coords match Java exactly:")
	t.Logf("  xCoord: %.17f", n.GetXCoord())
	t.Logf("  yCoord: %.17f", n.GetYCoord())
	t.Logf("  zCoord: %.17f", n.GetZCoord())
}

func TestImprovedNoise_Permutations_JavaParity(t *testing.T) {

	expectedFirst20 := []int{83, 88, 161, 90, 117, 220, 146, 221, 68, 86, 213, 124, 192, 112, 203, 19, 248, 82, 184, 97}
	expectedLast16 := []int{17, 43, 137, 70, 25, 118, 123, 16, 233, 51, 215, 193, 129, 198, 202, 14}

	r := rand.NewRandom(12345)
	n := NewImprovedNoise(r)
	perm := n.GetPermutations()

	allMatch := true
	for i := 0; i < 20; i++ {
		if perm[i] != expectedFirst20[i] {
			t.Errorf("permutations[%d] mismatch: got %d, expected %d", i, perm[i], expectedFirst20[i])
			allMatch = false
		}
	}

	for i := 0; i < 16; i++ {
		if perm[240+i] != expectedLast16[i] {
			t.Errorf("permutations[%d] mismatch: got %d, expected %d", 240+i, perm[240+i], expectedLast16[i])
			allMatch = false
		}
	}

	if allMatch {
		t.Log("✓ ImprovedNoise permutations: 36/36 checked values match Java exactly")
	}
}

func TestImprovedNoise_NoiseOutput_JavaParity(t *testing.T) {

	testCases := []struct {
		x, y, z float64
		noise   float64
	}{
		{0.0, 0.0, 0.0, -0.37230011176761096},
		{1.0, 0.0, 0.0, -0.49162915131929685},
		{0.0, 1.0, 0.0, -0.04668359815933333},
		{0.0, 0.0, 1.0, 0.14988113475352210},
		{0.5, 0.5, 0.5, -0.51448999811635290},
		{10.0, 20.0, 30.0, -0.25374550838340176},
		{-5.0, -10.0, -15.0, 0.41655778095012796},
	}

	r := rand.NewRandom(12345)
	n := NewImprovedNoise(r)

	tolerance := 1e-14
	allMatch := true

	for _, tc := range testCases {

		got := computeNoise3D(n, tc.x, tc.y, tc.z)

		if math.Abs(got-tc.noise) > tolerance {
			t.Errorf("noise(%.1f, %.1f, %.1f) mismatch: got %.17f, expected %.17f",
				tc.x, tc.y, tc.z, got, tc.noise)
			allMatch = false
		}
	}

	if allMatch {
		t.Log("✓ ImprovedNoise 3D output: 7/7 test coordinates match Java exactly")
	}
}

func computeNoise3D(n *ImprovedNoise, x, y, z float64) float64 {

	x = x + n.GetXCoord()
	y = y + n.GetYCoord()
	z = z + n.GetZCoord()

	i := int(math.Floor(x))
	j := int(math.Floor(y))
	k := int(math.Floor(z))

	x = x - float64(i)
	y = y - float64(j)
	z = z - float64(k)

	i = i & 255
	j = j & 255
	k = k & 255

	u := fade(x)
	v := fade(y)
	w := fade(z)

	perm := n.GetPermutations()

	A := perm[i] + j
	AA := perm[A] + k
	AB := perm[A+1] + k
	B := perm[i+1] + j
	BA := perm[B] + k
	BB := perm[B+1] + k

	return lerpNoise(w,
		lerpNoise(v,
			lerpNoise(u, grad3DTest(perm[AA], x, y, z), grad3DTest(perm[BA], x-1, y, z)),
			lerpNoise(u, grad3DTest(perm[AB], x, y-1, z), grad3DTest(perm[BB], x-1, y-1, z))),
		lerpNoise(v,
			lerpNoise(u, grad3DTest(perm[AA+1], x, y, z-1), grad3DTest(perm[BA+1], x-1, y, z-1)),
			lerpNoise(u, grad3DTest(perm[AB+1], x, y-1, z-1), grad3DTest(perm[BB+1], x-1, y-1, z-1))))
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerpNoise(t, a, b float64) float64 {
	return a + t*(b-a)
}

var gX = []float64{1, -1, 1, -1, 1, -1, 1, -1, 0, 0, 0, 0, 1, 0, -1, 0}
var gY = []float64{1, 1, -1, -1, 0, 0, 0, 0, 1, -1, 1, -1, 1, -1, 1, -1}
var gZ = []float64{0, 0, 0, 0, 1, 1, -1, -1, 1, 1, -1, -1, 0, 1, 0, -1}

func grad3DTest(hash int, x, y, z float64) float64 {
	i := hash & 15
	return gX[i]*x + gY[i]*y + gZ[i]*z
}
