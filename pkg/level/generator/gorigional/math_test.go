package gorigional

import (
	"math"
	"testing"
)

func TestSquashLogic(t *testing.T) {

	vanillaBaseHeight := 10.0

	volatility := 1.0

	h := NewMathHelper()

	gorigionalBaseHeight := vanillaBaseHeight / 2.0

	d_at_5 := h.CalculateDensity(5, gorigionalBaseHeight, volatility, 0, 0, -100)

	if math.Abs(d_at_5) > 0.001 {
		t.Errorf("Expected surface at index 5, got density %f", d_at_5)
	}

	d_at_6 := h.CalculateDensity(6, gorigionalBaseHeight, volatility, 0, 0, -100)
	expected_d := -12.0

	if math.Abs(d_at_6-expected_d) > 0.001 {
		t.Errorf("Expected density -12.0 at index 6, got %f", d_at_6)
	}
}
