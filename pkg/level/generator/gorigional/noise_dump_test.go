package gorigional

import (
	"fmt"
	"testing"

	"github.com/scaxe/scaxe-go/pkg/level/generator/gorigional/noise"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

func TestDumpNoiseMain(t *testing.T) {
	seed := int64(114514)
	rnd := rand.NewRandom(seed)

	minLimit := noise.NewOctavesNoise(rnd, 16)
	maxLimit := noise.NewOctavesNoise(rnd, 16)
	mainNoise := noise.NewOctavesNoise(rnd, 8)

	_ = minLimit
	_ = maxLimit

	x := -3
	z := -3
	hx := x * 4
	hz := z * 4

	buffer := make([]float64, 825)

	f := 684.412
	f1 := 684.412
	scaleX := 80.0
	scaleY := 160.0
	scaleZ := 80.0

	buffer = mainNoise.GenerateNoiseOctaves(buffer, hx, 0, hz, 5, 33, 5, f/scaleX, f1/scaleY, f/scaleZ)

	fmt.Println("Go Noise Dump:")
	for i := 0; i < 10; i++ {
		fmt.Printf("%d: %f\n", i, buffer[i])
	}
	fmt.Printf("MainNoise [100]: %v\n", buffer[100])
	fmt.Printf("MainNoise [500]: %v\n", buffer[500])
}
