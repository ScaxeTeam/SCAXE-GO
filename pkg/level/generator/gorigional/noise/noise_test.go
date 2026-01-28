package noise

import (
	"testing"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

func TestImprovedNoisePermutations(t *testing.T) {

	rnd := rand.NewRandom(1)

	noise := NewImprovedNoise(rnd)

	t.Logf("Permutations[0..10]: %v", noise.permutations[:10])
	t.Logf("Coords: %f, %f, %f", noise.xCoord, noise.yCoord, noise.zCoord)

	rnd2 := rand.NewRandom(1)
	noise2 := NewImprovedNoise(rnd2)

	for i := 0; i < 512; i++ {
		if noise.permutations[i] != noise2.permutations[i] {
			t.Errorf("Determinism failed at perm index %d", i)
		}
	}
	if noise.xCoord != noise2.xCoord {
		t.Errorf("Determinism failed at xCoord")
	}
}

func TestOctavesNoise(t *testing.T) {
	rnd := rand.NewRandom(1234)
	oct := NewOctavesNoise(rnd, 4)

	arr := make([]float64, 5*5*5)
	oct.GenerateNoiseOctaves(arr, 0, 0, 0, 5, 5, 5, 0.5, 0.5, 0.5)

	nonZero := 0
	for _, v := range arr {
		if v != 0 {
			nonZero++
		}
	}

	if nonZero == 0 {
		t.Error("Noise output is all zeros")
	}

	t.Logf("Sample Noise Value: %f", arr[0])
}
