package object

import (
	"fmt"
	"math"
	"testing"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

func TestTreeParity(t *testing.T) {
	seed := int64(12345)

	fmt.Println("=== Go Tree Generation Parity Test ===")
	fmt.Println("Seed:", seed)
	fmt.Println()

	fmt.Println("--- Test 1: Dark Oak Height Formula ---")
	fmt.Println("Formula: NextBoundedInt(3) + NextBoundedInt(2) + 6")
	r := rand.NewRandom(seed)
	for i := 0; i < 20; i++ {
		height := r.NextBoundedInt(3) + r.NextBoundedInt(2) + 6
		fmt.Printf("Tree %d: height=%d\n", i, height)
	}
	fmt.Println()

	fmt.Println("--- Test 2: Leaf Radius (Integer Division) ---")
	fmt.Println("Formula: j1 = 1 - i4/2")
	for i4 := -3; i4 <= 0; i4++ {
		j1 := 1 - i4/2
		fmt.Printf("i4=%d -> j1=%d\n", i4, j1)
	}
	fmt.Println()

	fmt.Println("--- Test 3: Corner Rounding Logic ---")
	fmt.Println("Formula: place if (abs(l1)!=j1 || abs(j2)!=j1 || (NextBoundedInt(2)!=0 && i4!=0))")
	r = rand.NewRandom(seed)
	treeHeight := 5
	placedCount := 0
	skippedCount := 0

	for i3 := 0 - 3 + treeHeight; i3 <= 0+treeHeight; i3++ {
		i4 := i3 - (0 + treeHeight)
		j1 := 1 - i4/2

		for k1 := 0 - j1; k1 <= 0+j1; k1++ {
			l1 := k1 - 0

			for i2 := 0 - j1; i2 <= 0+j1; i2++ {
				j2 := i2 - 0

				absL1 := int(math.Abs(float64(l1)))
				absJ2 := int(math.Abs(float64(j2)))

				shouldPlace := absL1 != j1 || absJ2 != j1 ||
					(r.NextBoundedInt(2) != 0 && i4 != 0)

				if shouldPlace {
					placedCount++
				} else {
					skippedCount++
					fmt.Printf("SKIP: i4=%d j1=%d l1=%d j2=%d\n", i4, j1, l1, j2)
				}
			}
		}
	}
	fmt.Printf("Total placed: %d, skipped: %d\n", placedCount, skippedCount)
	fmt.Println()

	fmt.Println("--- Test 4: RNG Sequence (first 50 values) ---")
	r = rand.NewRandom(seed)
	fmt.Print("NextBoundedInt(100): ")
	for i := 0; i < 10; i++ {
		fmt.Print(r.NextBoundedInt(100), " ")
	}
	fmt.Println()

	r = rand.NewRandom(seed)
	fmt.Print("NextBoolean: ")
	for i := 0; i < 10; i++ {
		fmt.Print(r.NextBoolean(), " ")
	}
	fmt.Println()

	fmt.Println()
	fmt.Println("=== Test Complete ===")
	fmt.Println("Compare with Java TreeParityTest output.")

	t.Run("DarkOakHeightRange", func(t *testing.T) {
		r := rand.NewRandom(seed)
		for i := 0; i < 1000; i++ {
			height := r.NextBoundedInt(3) + r.NextBoundedInt(2) + 6
			if height < 6 || height > 10 {
				t.Errorf("Dark Oak height out of range: %d (expected 6-10)", height)
			}
		}
	})

	t.Run("IntegerDivision", func(t *testing.T) {

		expected := map[int]int{-3: 2, -2: 2, -1: 1, 0: 1}
		for i4, expectedJ1 := range expected {
			j1 := 1 - i4/2
			if j1 != expectedJ1 {
				t.Errorf("i4=%d: got j1=%d, want %d", i4, j1, expectedJ1)
			}
		}
	})
}

func TestDarkOakHeightDistribution(t *testing.T) {
	counts := make(map[int]int)
	r := rand.NewRandom(42)

	for i := 0; i < 10000; i++ {
		height := r.NextBoundedInt(3) + r.NextBoundedInt(2) + 6
		counts[height]++
	}

	fmt.Println("Dark Oak Height Distribution (10000 samples):")
	for h := 6; h <= 10; h++ {
		pct := float64(counts[h]) / 100.0
		fmt.Printf("  Height %d: %d (%.1f%%)\n", h, counts[h], pct)
	}

	if counts[7] <= counts[6] {
		t.Logf("Warning: Height 7 should be more common than 6 (got %d vs %d)", counts[7], counts[6])
	}
}
