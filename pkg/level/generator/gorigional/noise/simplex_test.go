package noise

import (
	"fmt"
	"testing"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

func TestFastFloor(t *testing.T) {
	cases := []struct {
		in  float64
		out int
	}{
		{0.0, 0},
		{0.5, 0},
		{0.9, 0},
		{1.0, 1},
		{-0.1, -1},
		{-0.5, -1},
		{-0.9, -1},
		{-1.0, -1},
		{-1.1, -2},
	}

	for _, c := range cases {
		if res := fastFloor(c.in); res != c.out {
			t.Errorf("fastFloor(%f) = %d, want %d", c.in, res, c.out)
		}
	}
}

func TestSimplexNoise(t *testing.T) {

	rnd := rand.NewRandom(1)
	s := NewSimplexNoise(rnd)

	buffer := make([]float64, 256)

	s.Add(buffer, 0, 0, 16, 16, 0.0625, 0.0625, 1.0)

	val1 := buffer[0]

	for i := range buffer {
		buffer[i] = 0
	}

	s.Add(buffer, 0, 0, 16, 16, 0.0625, 0.0625, 1.0)
	if buffer[0] != val1 {
		t.Errorf("Simplex noise not deterministic! %f vs %f", val1, buffer[0])
	}

	fmt.Printf("Simplex[0] = %f\n", buffer[0])
	fmt.Printf("Simplex[1] = %f\n", buffer[1])
}
