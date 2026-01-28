package rand

import (
	"time"
)

type Random struct {
	seed int64
}

func NewRandom(seed int64) *Random {
	r := &Random{}
	r.SetSeed(seed)
	return r
}

func New() *Random {
	return NewRandom(time.Now().UnixNano())
}

func (r *Random) SetSeed(seed int64) {
	r.seed = (seed ^ 0x5DEECE66D) & ((1 << 48) - 1)
}

func (r *Random) next(bits int) int32 {
	r.seed = (r.seed*0x5DEECE66D + 0xB) & ((1 << 48) - 1)
	return int32(r.seed >> (48 - bits))
}

func (r *Random) NextInt() int32 {
	return r.next(32)
}

func (r *Random) NextBoundedInt(n int) int {
	if n <= 0 {
		return 0
	}

	if (n & -n) == n {
		return int((int64(n) * int64(r.next(31))) >> 31)
	}

	var bits, val int32
	for {
		bits = r.next(31)
		val = bits % int32(n)
		if bits-val+(int32(n)-1) >= 0 {
			return int(val)
		}
	}
}

func (r *Random) NextFloat() float64 {
	return float64(r.next(24)) / (float64(1 << 24))
}

func (r *Random) NextDouble() float64 {
	return (float64(r.next(26))*float64(1<<27) + float64(r.next(27))) / (float64(1 << 53))
}

func (r *Random) NextRange(start, end int) int {
	return start + r.NextBoundedInt(end-start+1)
}

func (r *Random) NextBoolean() bool {
	return r.next(1) != 0
}

func (r *Random) NextLong() int64 {
	return (int64(r.next(32)) << 32) + int64(r.next(32))
}
