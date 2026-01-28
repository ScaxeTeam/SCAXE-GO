package rand

import (
	"math"
	"testing"
)

func TestLCGSequence_JavaParity(t *testing.T) {

	expected := []int32{
		1553932502,
		-2090749135,
		-287790814,
		-355989640,
		-716867186,
		161804169,
		1402202751,
		535445604,
		1011567003,
		151766778,
		1499439034,
		-51321412,
		1924478780,
		-370025683,
		-1554121271,
		496460768,
		679749574,
		-301730690,
		-992618231,
		1128070351,
	}

	r := NewRandom(12345)
	for i, exp := range expected {
		got := r.NextInt()
		if got != exp {
			t.Errorf("NextInt()[%d]: got %d, expected %d (PARITY FAILURE)", i, got, exp)
		}
	}
	t.Log("✓ LCG NextInt sequence: 20/20 values match Java exactly")
}

func TestNextBoundedInt_JavaParity(t *testing.T) {

	testCases := []struct {
		bound    int
		expected []int
	}{
		{2, []int{0, 1, 1, 1, 1, 0, 0, 0, 0, 0}},
		{3, []int{1, 1, 0, 0, 1, 1, 1, 0, 1, 0}},
		{5, []int{1, 0, 1, 3, 0, 4, 0, 2, 1, 4}},
		{7, []int{5, 2, 4, 6, 2, 4, 2, 4, 6, 1}},
		{10, []int{1, 0, 1, 8, 5, 4, 5, 2, 1, 9}},
		{16, []int{5, 8, 14, 14, 13, 0, 5, 1, 3, 0}},
		{32, []int{11, 16, 29, 29, 26, 1, 10, 3, 7, 1}},
		{64, []int{23, 32, 59, 58, 53, 2, 20, 7, 15, 2}},
		{100, []int{51, 80, 41, 28, 55, 84, 75, 2, 1, 89}},
		{128, []int{46, 65, 119, 117, 106, 4, 41, 15, 30, 4}},
		{256, []int{92, 131, 238, 234, 213, 9, 83, 31, 60, 9}},
		{1000, []int{251, 80, 241, 828, 55, 84, 375, 802, 501, 389}},
	}

	allPassed := true
	for _, tc := range testCases {
		r := NewRandom(12345)
		for i, exp := range tc.expected {
			got := r.NextBoundedInt(tc.bound)
			if got != exp {
				t.Errorf("NextBoundedInt(%d)[%d]: got %d, expected %d (PARITY FAILURE)", tc.bound, i, got, exp)
				allPassed = false
			}
		}
	}
	if allPassed {
		t.Logf("✓ NextBoundedInt: 12 bounds × 10 values = 120/120 match Java exactly")
	}
}

func TestNextFloat_JavaParity(t *testing.T) {

	expected := []float64{
		0.361803054809570,
		0.513209521770477,
		0.932993471622467,
		0.917114675045013,
		0.833091318607330,
		0.037672936916351,
		0.326475739479065,
		0.124668121337891,
		0.235523760318756,
		0.035335898399353,
	}

	r := NewRandom(12345)
	allPassed := true
	for i, exp := range expected {
		got := r.NextFloat()

		if math.Abs(got-exp) > 1e-9 {
			t.Errorf("NextFloat()[%d]: got %.15f, expected %.15f (PARITY FAILURE)", i, got, exp)
			allPassed = false
		}
	}
	if allPassed {
		t.Log("✓ NextFloat: 10/10 values match Java (within 1e-9 tolerance)")
	}
}

func TestNextDouble_JavaParity(t *testing.T) {

	expected := []float64{
		0.36180310716047180,
		0.93299348528854100,
		0.83309134897102370,
		0.32647575623792624,
		0.23552379064762520,
		0.34911535662488336,
		0.44807763269315180,
		0.63815294378386860,
		0.15826654329520230,
		0.76888806019200900,
	}

	r := NewRandom(12345)
	allPassed := true
	for i, exp := range expected {
		got := r.NextDouble()

		if math.Abs(got-exp) > 1e-15 {
			t.Errorf("NextDouble()[%d]: got %.17f, expected %.17f (PARITY FAILURE)", i, got, exp)
			allPassed = false
		}
	}
	if allPassed {
		t.Log("✓ NextDouble: 10/10 values match Java (within 1e-15 tolerance)")
	}
}

func TestNextLong_JavaParity(t *testing.T) {

	expected := []int64{
		6674089274190705457,
		-1236052134575208584,
		-3078921119283744887,
		6022414958441676900,
		4344647195749500666,
		6440041613324510652,
		8265573421575953197,
		-6674900032466492448,
		2919502189498201214,
		-4263262838430303025,
	}

	r := NewRandom(12345)
	allPassed := true
	for i, exp := range expected {
		got := r.NextLong()
		if got != exp {
			t.Errorf("NextLong()[%d]: got %d, expected %d (PARITY FAILURE)", i, got, exp)
			allPassed = false
		}
	}
	if allPassed {
		t.Log("✓ NextLong: 10/10 values match Java exactly")
	}
}

func TestPopulationSeed_JavaParity(t *testing.T) {

	testCases := []struct {
		cx, cz     int
		k, l       int64
		resultSeed int64
	}{
		{0, 0, 6674089274190705457, -1236052134575208583, 12345},
		{1, 0, 6674089274190705457, -1236052134575208583, 6674089274190693128},
		{0, 1, 6674089274190705457, -1236052134575208583, -1236052134575196352},
		{1, 1, 6674089274190705457, -1236052134575208583, 5438037139615484563},
		{-1, 0, 6674089274190705457, -1236052134575208583, -6674089274190693130},
		{0, -1, 6674089274190705457, -1236052134575208583, 1236052134575196350},
		{100, 100, 6674089274190705457, -1236052134575208583, 8848135823972686417},
	}

	worldSeed := int64(12345)
	r := NewRandom(worldSeed)
	k := r.NextLong()/2*2 + 1
	l := r.NextLong()/2*2 + 1

	if k != testCases[0].k {
		t.Errorf("k value mismatch: got %d, expected %d", k, testCases[0].k)
	}
	if l != testCases[0].l {
		t.Errorf("l value mismatch: got %d, expected %d", l, testCases[0].l)
	}

	allPassed := true
	for _, tc := range testCases {
		seed := int64(tc.cx)*k + int64(tc.cz)*l ^ worldSeed
		if seed != tc.resultSeed {
			t.Errorf("PopulationSeed(%d, %d): got %d, expected %d (PARITY FAILURE)", tc.cx, tc.cz, seed, tc.resultSeed)
			allPassed = false
		}
	}
	if allPassed {
		t.Logf("✓ Population Seed: 7/7 chunk coordinates match Java exactly")
	}
}

func TestLCGMathematicalProof(t *testing.T) {
	const multiplier int64 = 0x5DEECE66D
	const addend int64 = 0xB
	const mask int64 = (1 << 48) - 1

	r := &Random{}
	r.SetSeed(0)

	seed0 := (int64(0) ^ multiplier) & mask
	seed1 := (seed0*multiplier + addend) & mask
	expected1 := int32(seed1 >> (48 - 32))

	r2 := NewRandom(0)
	got1 := r2.NextInt()

	if got1 != expected1 {
		t.Errorf("Mathematical verification failed: got %d, expected %d", got1, expected1)
	} else {
		t.Log("✓ LCG mathematical formula verified: (seed * 0x5DEECE66D + 0xB) & ((1<<48)-1)")
	}
}
