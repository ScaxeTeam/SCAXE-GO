package biome

import (
	"fmt"
	"strings"
	"testing"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

func TestColdTaigaTreeProbabilities(t *testing.T) {
	const iterations = 10000
	const tolerance = 0.03

	biome := NewColdTaigaBiome()
	r := rand.NewRandom(12345)

	pineCount := 0
	spruceCount := 0

	for i := 0; i < iterations; i++ {
		tree := biome.GetTreeFeature(r)

		treeType := getTreeTypeName(tree)
		switch treeType {
		case "PineTree":
			pineCount++
		case "SpruceTree":
			spruceCount++
		default:
			t.Errorf("Unexpected tree type: %s", treeType)
		}
	}

	pineRatio := float64(pineCount) / float64(iterations)
	spruceRatio := float64(spruceCount) / float64(iterations)

	expectedPine := 1.0 / 3.0
	expectedSpruce := 2.0 / 3.0

	if diff := pineRatio - expectedPine; diff < -tolerance || diff > tolerance {
		t.Errorf("Pine ratio %.4f outside expected %.4f ± %.2f", pineRatio, expectedPine, tolerance)
	}
	if diff := spruceRatio - expectedSpruce; diff < -tolerance || diff > tolerance {
		t.Errorf("Spruce ratio %.4f outside expected %.4f ± %.2f", spruceRatio, expectedSpruce, tolerance)
	}

	t.Logf("Cold Taiga Tree Distribution (n=%d):", iterations)
	t.Logf("  Pine:   %d (%.2f%%) - expected ~33.3%%", pineCount, pineRatio*100)
	t.Logf("  Spruce: %d (%.2f%%) - expected ~66.7%%", spruceCount, spruceRatio*100)
}

func TestExtremeHillsTreeProbabilities(t *testing.T) {
	const iterations = 10000
	const tolerance = 0.03

	biome := NewExtremeHillsPlusBiome()
	r := rand.NewRandom(12345)

	spruceCount := 0
	oakCount := 0

	for i := 0; i < iterations; i++ {
		tree := biome.GetTreeFeature(r)
		treeType := getTreeTypeName(tree)
		switch treeType {
		case "SpruceTree":
			spruceCount++
		case "OakTree":
			oakCount++
		default:
			t.Errorf("Unexpected tree type: %s", treeType)
		}
	}

	spruceRatio := float64(spruceCount) / float64(iterations)
	oakRatio := float64(oakCount) / float64(iterations)

	expectedSpruce := 2.0 / 3.0
	expectedOak := 1.0 / 3.0

	if diff := spruceRatio - expectedSpruce; diff < -tolerance || diff > tolerance {
		t.Errorf("Spruce ratio %.4f outside expected %.4f ± %.2f", spruceRatio, expectedSpruce, tolerance)
	}
	if diff := oakRatio - expectedOak; diff < -tolerance || diff > tolerance {
		t.Errorf("Oak ratio %.4f outside expected %.4f ± %.2f", oakRatio, expectedOak, tolerance)
	}

	t.Logf("Extreme Hills Tree Distribution (n=%d):", iterations)
	t.Logf("  Spruce: %d (%.2f%%) - expected ~66.7%%", spruceCount, spruceRatio*100)
	t.Logf("  Oak:    %d (%.2f%%) - expected ~33.3%%", oakCount, oakRatio*100)
}

func TestEmeraldOreYRange(t *testing.T) {
	const iterations = 10000

	r := rand.NewRandom(12345)

	minY := 1000
	maxY := -1000

	for i := 0; i < iterations; i++ {

		y := r.NextBoundedInt(28) + 4
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	if minY < 4 {
		t.Errorf("Emerald Ore generated below Y=4: minY=%d", minY)
	}
	if maxY > 31 {
		t.Errorf("Emerald Ore generated above Y=31: maxY=%d", maxY)
	}

	t.Logf("Emerald Ore Y Range (n=%d): min=%d, max=%d (expected: 4-31)",
		iterations, minY, maxY)
}

func TestSilverfishStoneYRange(t *testing.T) {
	const iterations = 10000

	r := rand.NewRandom(12345)

	minY := 1000
	maxY := -1000

	for i := 0; i < iterations; i++ {

		y := r.NextBoundedInt(64)
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	if minY < 0 {
		t.Errorf("Silverfish Stone generated below Y=0: minY=%d", minY)
	}
	if maxY > 63 {
		t.Errorf("Silverfish Stone generated above Y=63: maxY=%d", maxY)
	}

	t.Logf("Silverfish Stone Y Range (n=%d): min=%d, max=%d (expected: 0-63)",
		iterations, minY, maxY)
}

func TestBeachNoTrees(t *testing.T) {
	beach := NewBeachBiome()
	if beach.Decorator.TreesPerChunk != -999 {
		t.Errorf("Beach treesPerChunk = %d, expected -999", beach.Decorator.TreesPerChunk)
	}

	stoneBeach := NewStoneBeachBiome()
	if stoneBeach.Decorator.TreesPerChunk != -999 {
		t.Errorf("Stone Beach treesPerChunk = %d, expected -999", stoneBeach.Decorator.TreesPerChunk)
	}

	coldBeach := NewColdBeachBiome()
	if coldBeach.Decorator.TreesPerChunk != -999 {
		t.Errorf("Cold Beach treesPerChunk = %d, expected -999", coldBeach.Decorator.TreesPerChunk)
	}

	t.Log("All beach biomes correctly have treesPerChunk = -999")
}

func getTreeTypeName(gen Generator) string {
	if gen == nil {
		return "nil"
	}

	typeName := fmt.Sprintf("%T", gen)

	parts := strings.Split(typeName, ".")
	return parts[len(parts)-1]
}
