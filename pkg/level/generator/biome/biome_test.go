package biome

import (
	"testing"

	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

func TestJungleTreeDistribution(t *testing.T) {
	jungle := NewJungleBiome()
	r := rand.NewRandom(12345)

	counts := map[string]int{
		"BigOakTree":      0,
		"JungleBush":      0,
		"MegaJungleTree":  0,
		"JungleSmallTree": 0,
	}

	iterations := 10000
	for i := 0; i < iterations; i++ {
		tree := jungle.GetTreeFeature(r)
		switch tree.(type) {
		case *object.BigOakTree:
			counts["BigOakTree"]++
		case *object.JungleBush:
			counts["JungleBush"]++
		case *object.MegaJungleTree:
			counts["MegaJungleTree"]++
		case *object.JungleSmallTree:
			counts["JungleSmallTree"]++
		}
	}

	total := float64(iterations)
	bigOakPct := float64(counts["BigOakTree"]) / total * 100
	shrubPct := float64(counts["JungleBush"]) / total * 100
	megaPct := float64(counts["MegaJungleTree"]) / total * 100
	smallPct := float64(counts["JungleSmallTree"]) / total * 100

	t.Logf("Tree Distribution (n=%d):", iterations)
	t.Logf("  BigOak:       %.1f%% (expected ~10%%)", bigOakPct)
	t.Logf("  JungleShrub:  %.1f%% (expected ~45%%)", shrubPct)
	t.Logf("  MegaJungle:   %.1f%% (expected ~15%%)", megaPct)
	t.Logf("  SmallJungle:  %.1f%% (expected ~30%%)", smallPct)

	if bigOakPct < 5 || bigOakPct > 15 {
		t.Errorf("BigOak percentage %.1f%% outside expected range [5-15%%]", bigOakPct)
	}
	if shrubPct < 35 || shrubPct > 55 {
		t.Errorf("JungleShrub percentage %.1f%% outside expected range [35-55%%]", shrubPct)
	}
}

func TestMesaClayBandsDeterminism(t *testing.T) {
	mesa1 := NewMesaBiome()
	mesa2 := NewMesaBiome()

	for i := 0; i < 64; i++ {
		if mesa1.clayBands[i] != mesa2.clayBands[i] {
			t.Errorf("Clay band mismatch at layer %d: %d vs %d", i, mesa1.clayBands[i], mesa2.clayBands[i])
		}
	}

	t.Logf("Mesa Clay Bands (first 16): %v", mesa1.clayBands[:16])
}

func TestMegaTaigaBiomeProperties(t *testing.T) {
	mega := NewMegaTaigaBiome()

	if mega.Decorator.TreesPerChunk != 10 {
		t.Errorf("TreesPerChunk = %d, expected 10", mega.Decorator.TreesPerChunk)
	}
	if mega.Decorator.GrassPerChunk != 7 {
		t.Errorf("GrassPerChunk = %d, expected 7", mega.Decorator.GrassPerChunk)
	}
	if mega.Decorator.MushroomsPerChunk != 3 {
		t.Errorf("MushroomsPerChunk = %d, expected 3", mega.Decorator.MushroomsPerChunk)
	}
	if mega.Decorator.GrassGen == nil {
		t.Error("GrassGen is nil")
	} else if mega.Decorator.GrassGen.Type != 2 {
		t.Errorf("GrassGen.Type = %d, expected 2 (Fern)", mega.Decorator.GrassGen.Type)
	}

	t.Logf("MegaTaigaBiome: Trees=%d, Grass=%d, Mushrooms=%d",
		mega.Decorator.TreesPerChunk, mega.Decorator.GrassPerChunk, mega.Decorator.MushroomsPerChunk)
}

func TestDecoratorGrassGenExists(t *testing.T) {
	d := NewDecorator()
	if d.GrassGen == nil {
		t.Error("GrassGen is nil in new Decorator")
	}
}
