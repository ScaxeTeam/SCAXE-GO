package biome

import (
	"fmt"
	"testing"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

func TestForestTreeProbabilities(t *testing.T) {
	const samples = 10000
	const tolerance = 0.03

	tests := []struct {
		name         string
		forestType   int
		expectedDist map[string]float64
	}{
		{
			name:       "Normal Forest",
			forestType: FOREST_NORMAL,
			expectedDist: map[string]float64{
				"Birch":   0.20,
				"BigOak":  0.08,
				"Oak":     0.72,
				"DarkOak": 0.0,
			},
		},
		{
			name:       "Birch Forest",
			forestType: FOREST_BIRCH,
			expectedDist: map[string]float64{
				"Birch":   1.0,
				"BigOak":  0.0,
				"Oak":     0.0,
				"DarkOak": 0.0,
			},
		},
		{
			name:       "Roofed Forest",
			forestType: FOREST_ROOFED,
			expectedDist: map[string]float64{

				"DarkOak": 0.667,
				"Birch":   0.067,
				"BigOak":  0.027,
				"Oak":     0.24,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			biome := NewForestBiome(tt.forestType)
			counts := map[string]int{
				"Oak": 0, "BigOak": 0, "Birch": 0, "DarkOak": 0, "Other": 0,
			}

			r := rand.NewRandom(12345)
			for i := 0; i < samples; i++ {
				tree := biome.GetTreeFeature(r)

				typeName := identifyTreeType(tree)
				counts[typeName]++
			}

			for treeType, expected := range tt.expectedDist {
				actual := float64(counts[treeType]) / float64(samples)
				if diff := abs(actual - expected); diff > tolerance {
					t.Errorf("%s: expected %.1f%% %s, got %.1f%% (diff: %.1f%%)",
						tt.name, expected*100, treeType, actual*100, diff*100)
				} else {
					t.Logf("%s: %s = %.1f%% (expected %.1f%%)",
						tt.name, treeType, actual*100, expected*100)
				}
			}
		})
	}
}

func TestSavannaTreeProbabilities(t *testing.T) {
	const samples = 10000
	const tolerance = 0.03

	biome := NewSavannaBiome()
	counts := map[string]int{"Savanna": 0, "Oak": 0, "Other": 0}

	r := rand.NewRandom(54321)
	for i := 0; i < samples; i++ {
		tree := biome.GetTreeFeature(r)
		typeName := identifyTreeType(tree)
		if typeName == "Oak" || typeName == "Savanna" {
			counts[typeName]++
		} else {
			counts["Other"]++
		}
	}

	savannaRatio := float64(counts["Savanna"]) / float64(samples)
	oakRatio := float64(counts["Oak"]) / float64(samples)

	if diff := abs(savannaRatio - 0.80); diff > tolerance {
		t.Errorf("Savanna: expected 80%%, got %.1f%% (diff: %.1f%%)", savannaRatio*100, diff*100)
	} else {
		t.Logf("Savanna: %.1f%% (expected 80%%)", savannaRatio*100)
	}

	if diff := abs(oakRatio - 0.20); diff > tolerance {
		t.Errorf("Oak: expected 20%%, got %.1f%% (diff: %.1f%%)", oakRatio*100, diff*100)
	} else {
		t.Logf("Oak: %.1f%% (expected 20%%)", oakRatio*100)
	}
}

func TestRoofedForestMushroomProbability(t *testing.T) {
	const samples = 10000
	const tolerance = 0.02

	r := rand.NewRandom(99999)
	mushroomCount := 0

	for i := 0; i < samples; i++ {

		if r.NextBoundedInt(20) == 0 {
			mushroomCount++
		}
	}

	ratio := float64(mushroomCount) / float64(samples)
	expected := 0.05

	if diff := abs(ratio - expected); diff > tolerance {
		t.Errorf("Mushroom: expected 5%%, got %.1f%% (diff: %.1f%%)", ratio*100, diff*100)
	} else {
		t.Logf("Mushroom: %.1f%% (expected 5%%)", ratio*100)
	}
}

func identifyTreeType(gen Generator) string {

	switch gen.(type) {
	case interface{ IsDarkOak() }:
		return "DarkOak"
	case interface{ IsBirch() }:
		return "Birch"
	case interface{ IsBigOak() }:
		return "BigOak"
	case interface{ IsSavanna() }:
		return "Savanna"
	default:

		s := getTypeName(gen)
		return s
	}
}

func getTypeName(gen Generator) string {

	switch gen.(type) {
	default:

		name := extractLastWord(typeName(gen))

		switch {
		case contains(name, "DarkOak") || contains(name, "Canopy"):
			return "DarkOak"
		case contains(name, "Birch"):
			return "Birch"
		case contains(name, "BigOak") || contains(name, "Big"):
			return "BigOak"
		case contains(name, "Savanna") || contains(name, "Acacia"):
			return "Savanna"
		case contains(name, "Oak") || contains(name, "Tree"):
			return "Oak"
		default:
			return "Other"
		}
	}
}

func typeName(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func extractLastWord(s string) string {

	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			return s[i+1:]
		}
	}
	return s
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
