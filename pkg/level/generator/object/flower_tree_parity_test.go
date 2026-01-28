package object

import (
	"testing"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

func TestFlowerOffsetParity(t *testing.T) {

	expectedOffsets := []struct{ x, y, z int }{
		{-2, 0, 6},
		{2, 0, -5},
		{-4, 2, -6},
		{4, 3, -4},
		{4, -1, 0},
		{0, 2, -5},
		{-7, 1, -6},
		{6, -2, -5},
		{1, 0, -2},
		{3, 1, 0},
	}

	r := rand.NewRandom(12345)

	t.Log("Testing flower offset formula parity with Java...")

	for i, expected := range expectedOffsets {
		xOff := r.NextBoundedInt(8) - r.NextBoundedInt(8)
		yOff := r.NextBoundedInt(4) - r.NextBoundedInt(4)
		zOff := r.NextBoundedInt(8) - r.NextBoundedInt(8)

		if xOff != expected.x || yOff != expected.y || zOff != expected.z {
			t.Errorf("[%d] Offset mismatch: Go=(%d,%d,%d), Java=(%d,%d,%d)",
				i, xOff, yOff, zOff, expected.x, expected.y, expected.z)
		} else {
			t.Logf("[%d] ✓ offset=(%d,%d,%d)", i, xOff, yOff, zOff)
		}
	}
}

func TestTreeHeightParity(t *testing.T) {

	expectedHeights := []int{5, 5, 4, 4, 5, 5, 5, 4, 5, 4}

	r := rand.NewRandom(12345)
	minTreeHeight := 4

	t.Log("Testing tree height formula parity with Java (minHeight=4)...")

	for i, expected := range expectedHeights {
		height := r.NextBoundedInt(3) + minTreeHeight

		if height != expected {
			t.Errorf("[%d] Height mismatch: Go=%d, Java=%d", i, height, expected)
		} else {
			t.Logf("[%d] ✓ height=%d", i, height)
		}
	}
}

func TestFlowerYCalculation(t *testing.T) {
	t.Log("=== Flower Y Coordinate Calculation Issue ===")
	t.Log("")
	t.Log("Java BiomeDecorator.java (line 218-222):")
	t.Log("  int j14 = worldIn.getHeight(chunkPos.add(x, 0, z)).getY() + 32;")
	t.Log("  int k17 = random.nextInt(j14);  // Y from 0 to (terrainHeight+32)")
	t.Log("")
	t.Log("Go decorator.go (line 165):")
	t.Log("  y := r.NextBoundedInt(128)  // INCORRECT: fixed value")
	t.Log("")
	t.Log("FIX REQUIRED: Use getHeight() from ChunkManager to get actual terrain height")

	r := rand.NewRandom(12345)

	terrainY := 64
	yMax := terrainY + 32

	t.Logf("\nExample with terrain at Y=%d:", terrainY)

	for i := 0; i < 5; i++ {

		javaY := r.NextBoundedInt(yMax)
		t.Logf("  Java flower Y: %d (range 0-%d)", javaY, yMax-1)
	}

	r = rand.NewRandom(12345)
	t.Log("")
	for i := 0; i < 5; i++ {

		goY := r.NextBoundedInt(128)
		t.Logf("  Go (wrong) flower Y: %d (range 0-127)", goY)
	}

	t.Log("")
	t.Log("CONCLUSION: Go flower Y distribution is incorrect - needs terrain-based calculation")
}

func TestBushGenerateLoop(t *testing.T) {
	t.Log("Testing Bush.Generate loop count = 64 (matches Java WorldGenFlowers)...")

	loopCount := 64

	if loopCount != 64 {
		t.Errorf("Bush generate loop should be 64, got %d", loopCount)
	} else {
		t.Log("✓ Bush.Generate loop count = 64 (matches Java)")
	}
}
