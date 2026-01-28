package layer

import "testing"

func TestGenLayerIsland(t *testing.T) {
	seed := int64(1)
	island := NewGenLayerIsland(seed)
	island.InitWorldGenSeed(234234)

	x, z := -5, -5
	w, d := 10, 10
	ints := island.GetInts(x, z, w, d)

	if len(ints) != w*d {
		t.Fatalf("Expected %d ints, got %d", w*d, len(ints))
	}

	centerIdx := -x + -z*w
	if ints[centerIdx] != 1 {
		t.Errorf("Center (0,0) at index %d must be LAND (1), got %d", centerIdx, ints[centerIdx])
	}

	ints2 := island.GetInts(x, z, w, d)
	for i, v := range ints {
		if ints2[i] != v {
			t.Errorf("Determinism failure at index %d", i)
		}
	}

	count := 0
	for _, v := range ints {
		if v == 1 {
			count++
		}
	}
	t.Logf("Land cells in 100 blocks: %d (Expected approx 10 + 1 forced)", count)
}
