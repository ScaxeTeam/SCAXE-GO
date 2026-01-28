package layer

import (
	"testing"
)

func TestGenLayerLCG(t *testing.T) {

	layer := NewBaseLayer(1)
	layer.InitWorldGenSeed(12345)
	layer.InitChunkSeed(0, 0)

	val1 := layer.NextInt(10)
	val2 := layer.NextInt(10)
	val3 := layer.NextInt(10)

	layer2 := NewBaseLayer(1)
	layer2.InitWorldGenSeed(12345)
	layer2.InitChunkSeed(0, 0)

	if v := layer2.NextInt(10); v != val1 {
		t.Errorf("Mismatch val1: %d vs %d", v, val1)
	}
	if v := layer2.NextInt(10); v != val2 {
		t.Errorf("Mismatch val2: %d vs %d", v, val2)
	}
	if v := layer2.NextInt(10); v != val3 {
		t.Errorf("Mismatch val3: %d vs %d", v, val3)
	}

	t.Logf("Generated sequence: %d, %d, %d", val1, val2, val3)
}

func TestMixSeed(t *testing.T) {
	l := &BaseLayer{}
	res := l.mixSeed(1, 1)

	if res == 0 {
		t.Error("mixSeed returned 0, unlikely")
	}
}

func TestLayerStack(t *testing.T) {

	island := NewGenLayerIsland(1)
	fuzzy := NewGenLayerFuzzyZoom(2000, island)
	addIsland := NewGenLayerAddIsland(1, fuzzy)
	zoom := NewGenLayerZoom(2001, addIsland)
	removeOcean := NewGenLayerRemoveTooMuchOcean(2, zoom)

	removeOcean.InitWorldGenSeed(12345)

	ints := removeOcean.GetInts(0, 0, 16, 16)

	if len(ints) != 256 {
		t.Fatalf("Output size mismatch: %d", len(ints))
	}

	landCount := 0
	for _, v := range ints {
		if v == 1 {
			landCount++
		}
	}
	t.Logf("Stack Output Land Count: %d / 256", landCount)

	removeOcean2 := NewGenLayerRemoveTooMuchOcean(2, NewGenLayerZoom(2001, NewGenLayerAddIsland(1, NewGenLayerFuzzyZoom(2000, NewGenLayerIsland(1)))))
	removeOcean2.InitWorldGenSeed(12345)
	ints2 := removeOcean2.GetInts(0, 0, 16, 16)

	for i, v := range ints {
		if ints2[i] != v {
			t.Errorf("Stack determinism check failed at %d", i)
		}
	}
}

func TestLayerStackBatch2(t *testing.T) {

	island := NewGenLayerIsland(1)
	zoom := NewGenLayerZoom(2000, island)
	addIsland := NewGenLayerAddIsland(1, zoom)
	removeOcean := NewGenLayerRemoveTooMuchOcean(2, addIsland)
	addSnow := NewGenLayerAddSnow(3, removeOcean)
	deepOcean := NewGenLayerDeepOcean(4, addSnow)

	deepOcean.InitWorldGenSeed(12345)

	w, h := 64, 64
	ints := deepOcean.GetInts(0, 0, w, h)

	snowCount := 0
	coolCount := 0
	deepOceanCount := 0

	for _, v := range ints {
		if v == 4 {
			snowCount++
		}
		if v == 3 {
			coolCount++
		}
		if v == 24 {
			deepOceanCount++
		}
	}

	t.Logf("Batch 2 Stats: Snow=%d, Cool=%d, DeepOcean=%d / %d", snowCount, coolCount, deepOceanCount, w*h)

}

func TestLayerStackBiome(t *testing.T) {

	mock := &MockLayer{

		Vals: []int{0, 1, 2, 3, 4, 1, 2, 3, 4},
		W:    3, H: 3,
	}

	biomeLayer := NewGenLayerBiome(100, mock)
	biomeLayer.InitWorldGenSeed(12345)

	ints := biomeLayer.GetInts(0, 0, 3, 3)

	t.Logf("Biome Output: %v", ints)

	foundForest := false
	foundDesert := false
	foundIce := false
	foundTaiga := false
	foundJungle := false

	for _, id := range ints {
		if id == 4 {
			foundForest = true
		}
		if id == 2 {
			foundDesert = true
		}
		if id == 12 || id == 30 {
			foundIce = true
		}
		if id == 5 {
			foundTaiga = true
		}
		if id == 21 {
			foundJungle = true
		}
	}

	t.Logf("Found Biomes: Forest=%v, Desert=%v, Ice=%v, Taiga=%v, Jungle=%v",
		foundForest, foundDesert, foundIce, foundTaiga, foundJungle)

	if len(ints) != 9 {
		t.Errorf("Expected 9 ints, got %d", len(ints))
	}
}

type MockLayer struct {
	Vals       []int
	W, H       int
	*BaseLayer
}

func (m *MockLayer) GetInts(x, y, w, h int) []int {

	res := make([]int, w*h)
	for i := 0; i < w*h; i++ {
		res[i] = m.Vals[i%len(m.Vals)]
	}
	return res
}
func (m *MockLayer) InitWorldGenSeed(s int64) {}
