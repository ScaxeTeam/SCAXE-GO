package world

import (
	"math/rand"
)

type Random struct {
	*rand.Rand
}

func NewRandom(seed int64) *Random {
	return &Random{
		Rand: rand.New(rand.NewSource(seed)),
	}
}

func (r *Random) NextRange(min, max int) int {
	if max <= min {
		return min
	}
	return min + r.Intn(max-min+1)
}

func (r *Random) NextFloat() float64 {
	return r.Float64()
}

func (r *Random) NextBool() bool {
	return r.Intn(2) == 1
}

func (r *Random) NextBoundedInt(bound int) int {
	if bound <= 0 {
		return 0
	}
	return r.Intn(bound)
}
