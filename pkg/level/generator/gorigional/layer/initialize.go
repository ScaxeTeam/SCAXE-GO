package layer

func InitializeAll(seed int64) []GenLayer {

	l := GenLayer(NewGenLayerIsland(1))
	l = NewGenLayerFuzzyZoom(2000, l)
	l = NewGenLayerAddIsland(1, l)
	l = NewGenLayerZoom(2001, l)
	l = NewGenLayerAddIsland(2, l)
	l = NewGenLayerAddIsland(50, l)
	l = NewGenLayerAddIsland(70, l)
	l = NewGenLayerRemoveTooMuchOcean(2, l)
	l = NewGenLayerAddSnow(2, l)
	l = NewGenLayerAddIsland(3, l)
	l = NewGenLayerEdge(2, l, CoolWarm)
	l = NewGenLayerEdge(2, l, HeatIce)
	l = NewGenLayerEdge(3, l, Special)
	l = NewGenLayerZoom(2002, l)
	l = NewGenLayerZoom(2003, l)
	l = NewGenLayerAddIsland(4, l)

	l = NewGenLayerAddMushroomIsland(5, l)

	l = NewGenLayerDeepOcean(4, l)

	l4 := Magnify(1000, l, 0)

	l7 := Magnify(1000, l4, 0)

	riverInit := GenLayer(NewGenLayerRiverInit(100, l7))

	l8 := GenLayer(NewGenLayerBiome(200, l4))

	l6 := Magnify(1000, l8, 2)

	biomeEdge := GenLayer(NewGenLayerBiomeEdge(1000, l6))

	l9 := Magnify(1000, riverInit, 2)

	hills := GenLayer(NewGenLayerHills(1000, biomeEdge, l9))

	l5 := Magnify(1000, riverInit, 2)

	l5 = Magnify(1000, l5, 4)

	river := GenLayer(NewGenLayerRiver(1, l5))
	smoothRiver := GenLayer(NewGenLayerSmooth(1000, river))

	hills = NewGenLayerRareBiome(1001, hills)

	for k := 0; k < 4; k++ {
		hills = NewGenLayerZoom(int64(1000+k), hills)

		if k == 0 {
			hills = NewGenLayerAddIsland(3, hills)
		}

		if k == 1 {
			hills = NewGenLayerShore(1000, hills)
		}
	}

	smoothHills := GenLayer(NewGenLayerSmooth(1000, hills))

	riverMix := GenLayer(NewGenLayerRiverMix(100, smoothHills, smoothRiver))

	voronoi := GenLayer(NewGenLayerVoronoiZoom(10, riverMix))

	riverMix.InitWorldGenSeed(seed)
	voronoi.InitWorldGenSeed(seed)

	return []GenLayer{riverMix, voronoi, riverMix}
}
