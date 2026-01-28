package rand

import (
	"time"
)

type XorShift128 struct {
	x, y, z, w int64
}

const (
	X128 = 123456789
	Y128 = 362436069
	Z128 = 521288629
	W128 = 88675123
)

func NewXorShift128(seed int64) *XorShift128 {
	if seed == -1 {
		seed = time.Now().UnixNano()
	}
	r := &XorShift128{}
	r.SetSeed(seed)
	return r
}

func (r *XorShift128) SetSeed(seed int64) {

	r.x = (X128 ^ seed) & 0xffffffff

	term1_y := Y128 ^ (seed << 17)
	term2_y := (seed >> 15) & 0x7fffffff
	r.y = (term1_y | term2_y) & 0xffffffff

	term1_z := Z128 ^ (seed << 31)
	term2_z := (seed >> 1) & 0x7fffffff
	r.z = (term1_z | term2_z) & 0xffffffff

	term1_w := W128 ^ (seed << 18)
	term2_w := (seed >> 14) & 0x7fffffff
	r.w = (term1_w | term2_w) & 0xffffffff
}

func (r *XorShift128) NextSignedInt() int32 {
	t := (r.x ^ (r.x << 11)) & 0xffffffff

	r.x = r.y
	r.y = r.z
	r.z = r.w

	part1 := r.w ^ ((r.w >> 19) & 0x7fffffff)

	part2 := t ^ ((t >> 8) & 0x7fffffff)

	r.w = (part1 ^ part2) & 0xffffffff

	return int32(r.w)
}

func (r *XorShift128) NextInt() int32 {
	return r.NextSignedInt() & 0x7fffffff
}

func (r *XorShift128) NextFloat() float64 {
	return float64(r.NextInt()) / float64(0x7fffffff)
}

func (r *XorShift128) NextSignedFloat() float64 {
	return float64(r.NextSignedInt()) / float64(0x7fffffff)
}

func (r *XorShift128) NextBoundedInt(bound int) int {
	return int(r.NextInt()) % bound
}

func (r *XorShift128) NextRange(start, end int) int {
	return start + (int(r.NextInt()) % (end + 1 - start))
}
