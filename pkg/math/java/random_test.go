package java

import (
	"testing"
)

func TestRandom_Values(t *testing.T) {

	r := NewRandom(0)
	val := r.NextInt()
	t.Logf("Seed 0 NextInt: %d", val)

	r.SetSeed(0)
	valLong := r.NextLong()
	t.Logf("Seed 0 NextLong: %d", valLong)

	r.SetSeed(0)
	valFloat := r.NextFloat()
	t.Logf("Seed 0 NextFloat: %f", valFloat)

	r.SetSeed(0)
	valDouble := r.NextDouble()
	t.Logf("Seed 0 NextDouble: %f", valDouble)
}
