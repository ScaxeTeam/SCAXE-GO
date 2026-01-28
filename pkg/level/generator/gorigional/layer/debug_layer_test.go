package layer

import (
	"fmt"
	"testing"
)

func TestDumpIntermediateLayers(t *testing.T) {
	seed := int64(114514)
	fmt.Printf("Dumping Layers for Seed: %d\n", seed)

	var l GenLayer = NewGenLayerIsland(1)
	l = NewGenLayerFuzzyZoom(2000, l)
	l = NewGenLayerAddIsland(1, l)
	l = NewGenLayerZoom(2001, l)
	l = NewGenLayerAddIsland(2, l)
	l = NewGenLayerAddIsland(50, l)
	l = NewGenLayerAddIsland(70, l)
	l = NewGenLayerRemoveTooMuchOcean(2, l)
	l = NewGenLayerAddSnow(2, l)
	dumpLayer(t, "Manual AddSnow", l, seed)

	l = NewGenLayerAddIsland(3, l)
	l = NewGenLayerEdge(2, l, CoolWarm)
	dumpLayer(t, "Manual Edge CoolWarm", l, seed)

	l = NewGenLayerZoom(2002, l)
	l = NewGenLayerZoom(2003, l)
	l = NewGenLayerAddIsland(4, l)
	l = NewGenLayerAddMushroomIsland(5, l)
	l = NewGenLayerDeepOcean(4, l)

	l = NewGenLayerBiome(200, l)

	fmt.Println("--- InitializeAll Stack ---")
	layers := InitializeAll(seed)
	riverMix := layers[0].(*GenLayerRiverMix)
	current := riverMix.biomePatternGeneratorChain

	smoothHills := current.(*GenLayerSmooth)
	current = smoothHills.Parent

	var hills *GenLayerHills
	for {
		if h, ok := current.(*GenLayerHills); ok {
			hills = h
			break
		}
		switch val := current.(type) {
		case *GenLayerZoom:
			current = val.Parent
		case *GenLayerShore:
			current = val.Parent
		case *GenLayerAddIsland:
			current = val.Parent
		case *GenLayerRareBiome:
			current = val.Parent
		case *GenLayerRiver:
			current = val.Parent
		case *GenLayerSmooth:
			current = val.Parent
		default:
			panic(fmt.Sprintf("Unknown layer in Hills chain: %T", current))
		}
	}

	current = hills.Parent
	biomeEdge := current.(*GenLayerBiomeEdge)
	current = biomeEdge.Parent
	z1 := current.(*GenLayerZoom)
	current = z1.Parent
	z2 := current.(*GenLayerZoom)
	current = z2.Parent
	l8 := current

	biomeLayer := l8.(*GenLayerBiome)
	deepOcean := biomeLayer.Parent

	if z, ok := deepOcean.(*GenLayerZoom); ok {
		fmt.Printf("InitializeAll DeepOcean IS wrapped in Zoom! %T\n", z)
		deepOcean = z.Parent
	}

	if _, ok := deepOcean.(*GenLayerDeepOcean); ok {
		dumpLayer(t, "InitializeAll DeepOcean", deepOcean, seed)
	} else {
		fmt.Printf("Expected DeepOcean parent, got %T\n", deepOcean)
	}

	dumpLayer(t, "InitializeAll Biome", l8, seed)
}

func dumpLayer(t *testing.T, name string, l GenLayer, seed int64) {
	fmt.Printf("--- %s (16x16 at 0,0) ---\n", name)
	l.InitWorldGenSeed(seed)
	data := l.GetInts(0, 0, 16, 16)
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			fmt.Printf("%4d ", data[x+y*16])
		}
		fmt.Println()
	}
}
