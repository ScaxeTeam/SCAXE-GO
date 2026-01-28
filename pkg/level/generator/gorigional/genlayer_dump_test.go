package gorigional

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/scaxe/scaxe-go/pkg/level/generator/gorigional/layer"
)

func TestDumpGenLayers(t *testing.T) {
	seed := int64(114514)

	var genlayer layer.GenLayer = layer.NewGenLayerIsland(1)
	genlayer = layer.NewGenLayerFuzzyZoom(2000, genlayer)
	genlayeraddisland := layer.NewGenLayerAddIsland(1, genlayer)
	genlayerzoom := layer.NewGenLayerZoom(2001, genlayeraddisland)
	var genlayeraddisland1 layer.GenLayer = layer.NewGenLayerAddIsland(2, genlayerzoom)
	genlayeraddisland1 = layer.NewGenLayerAddIsland(50, genlayeraddisland1)
	genlayeraddisland1 = layer.NewGenLayerAddIsland(70, genlayeraddisland1)
	genlayerremovetoomuchocean := layer.NewGenLayerRemoveTooMuchOcean(2, genlayeraddisland1)
	genlayeraddsnow := layer.NewGenLayerAddSnow(2, genlayerremovetoomuchocean)
	genlayeraddisland2 := layer.NewGenLayerAddIsland(3, genlayeraddsnow)
	var genlayeredge layer.GenLayer = layer.NewGenLayerEdge(2, genlayeraddisland2, layer.CoolWarm)
	genlayeredge = layer.NewGenLayerEdge(2, genlayeredge, layer.HeatIce)
	genlayeredge = layer.NewGenLayerEdge(3, genlayeredge, layer.Special)
	var genlayerzoom1 layer.GenLayer = layer.NewGenLayerZoom(2002, genlayeredge)
	genlayerzoom1 = layer.NewGenLayerZoom(2003, genlayerzoom1)
	genlayeraddisland3 := layer.NewGenLayerAddIsland(4, genlayerzoom1)
	genlayeraddmushroomisland := layer.NewGenLayerAddMushroomIsland(5, genlayeraddisland3)
	genlayerdeepocean := layer.NewGenLayerDeepOcean(4, genlayeraddmushroomisland)
	genlayer4 := layer.Magnify(1000, genlayerdeepocean, 0)

	i := 4
	j := 4

	lvt_7_1_ := layer.Magnify(1000, genlayer4, 0)
	genlayerriverinit := layer.NewGenLayerRiverInit(100, lvt_7_1_)
	lvt_8_1_ := layer.NewGenLayerBiome(200, genlayer4)
	genlayer6 := layer.Magnify(1000, lvt_8_1_, 2)
	genlayerbiomeedge := layer.NewGenLayerBiomeEdge(1000, genlayer6)
	lvt_9_1_ := layer.Magnify(1000, genlayerriverinit, 2)
	var genlayerhills layer.GenLayer = layer.NewGenLayerHills(1000, genlayerbiomeedge, lvt_9_1_)
	genlayer5 := layer.Magnify(1000, genlayerriverinit, 2)
	genlayer5 = layer.Magnify(1000, genlayer5, j)
	genlayerriver := layer.NewGenLayerRiver(1, genlayer5)
	genlayersmooth := layer.NewGenLayerSmooth(1000, genlayerriver)
	genlayerhills = layer.NewGenLayerRareBiome(1001, genlayerhills)

	for k := 0; k < i; k++ {
		genlayerhills = layer.NewGenLayerZoom(int64(1000+k), genlayerhills)
		if k == 0 {
			genlayerhills = layer.NewGenLayerAddIsland(3, genlayerhills)
		}
		if k == 1 || i == 1 {
			genlayerhills = layer.NewGenLayerShore(1000, genlayerhills)
		}
	}

	genlayersmooth1 := layer.NewGenLayerSmooth(1000, genlayerhills)
	genlayerrivermix := layer.NewGenLayerRiverMix(100, genlayersmooth1, genlayersmooth)
	genlayer3 := layer.NewGenLayerVoronoiZoom(10, genlayerrivermix)

	genlayerrivermix.InitWorldGenSeed(seed)
	genlayer3.InitWorldGenSeed(seed)

	x := 0
	z := 0
	size := 256

	dumpLayer(t, "1_Island", genlayer, x, z, 16, 16)
	dumpLayer(t, "2_AddIsland_3", genlayeraddisland2, x, z, 16, 16)
	dumpLayer(t, "3_AddSnow", genlayeraddsnow, x, z, 16, 16)
	dumpLayer(t, "3.5_Edge_Special", genlayeredge, x, z, 16, 16)
	dumpLayer(t, "3.6_Zoom_2003", genlayerzoom1, x, z, 64, 64)
	dumpLayer(t, "3.7_AddIsland_4", genlayeraddisland3, x, z, 64, 64)
	dumpLayer(t, "3.8_AddMushroom", genlayeraddmushroomisland, x, z, 64, 64)
	dumpLayer(t, "4_DeepOcean", genlayerdeepocean, x, z, size, size)
	dumpLayer(t, "5_Biome", lvt_8_1_, x, z, size, size)
	dumpLayer(t, "6_RiverInit", genlayerriverinit, x, z, size, size)
	dumpLayer(t, "6.5_BiomeZoom", genlayer6, x, z, size, size)
	dumpLayer(t, "7_BiomeEdge", genlayerbiomeedge, x, z, size, size)
	dumpLayer(t, "8_Hills", genlayerhills, x, z, size, size)
	dumpLayer(t, "9_River", genlayerriver, x, z, size, size)
	dumpLayer(t, "10_FinalRiverMix", genlayerrivermix, x, z, size, size)
	dumpLayer(t, "11_Voronoi", genlayer3, x, z, size, size)

	t.Logf("Dump Complete")
}

func dumpLayer(t *testing.T, name string, l layer.GenLayer, x, z, w, h int) {
	data := l.GetInts(x, z, w, h)

	f, err := os.Create("../../../../../dump_go_" + name + ".csv")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	var sb strings.Builder
	for i, val := range data {
		sb.WriteString(strconv.Itoa(val))
		if (i+1)%w == 0 {
			sb.WriteString("\n")
		} else {
			sb.WriteString(",")
		}
	}
	f.WriteString(sb.String())
	fmt.Printf("Dumped %s\n", name)
}
