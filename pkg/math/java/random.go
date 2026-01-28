package java

import (
	"sync"
)

const (
	multiplier = 0x5DEECE66D
	addend     = 0xB
	mask       = (1 << 48) - 1
)

type Random struct {
	seed int64
	mu sync.Mutex
}

func NewRandom(seed int64) *Random {
	r := &Random{}
	r.SetSeed(seed)
	return r
}

func (r *Random) SetSeed(seed int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seed = (seed ^ multiplier) & mask
}

func (r *Random) next(bits int) int32 {
	r.mu.Lock()
	defer r.mu.Unlock()

	nextseed := (r.seed*multiplier + addend) & mask
	r.seed = nextseed

	return int32(nextseed >> (48 - bits))
}

func (r *Random) NextInt() int32 {
	return r.next(32)
}

func (r *Random) NextIntn(bound int32) int32 {
	if bound <= 0 {
		panic("bound must be positive")
	}

	if (bound & -bound) == bound {
		return int32((int64(bound) * int64(r.next(31))) >> 31)
	}

	var bits, val int32
	for {
		bits = r.next(31)
		val = bits % bound
		if bits-val+(bound-1) >= 0 {
			return val
		}
	}
}

func (r *Random) NextLong() int64 {
	return (int64(r.next(32)) << 32) + int64(r.next(32))
}

func (r *Random) NextBoolean() bool {
	return r.next(1) != 0
}

func (r *Random) NextFloat() float32 {
	return float32(r.next(24)) / (float32(1 << 24))
}

func (r *Random) NextDouble() float64 {
	high := int64(r.next(26)) << 27
	low := int64(r.next(27))
	return float64(high+low) / (float64(1 << 53))
}

type GaussianRandom struct {
	*Random
	nextNextGaussian     float64
	haveNextNextGaussian bool
}
