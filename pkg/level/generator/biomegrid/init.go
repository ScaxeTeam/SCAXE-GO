package biomegrid

func Initialize(seed int64) (MapLayer, MapLayer) {
	zoom := 2

	layer := MapLayer(NewNoiseMapLayer(seed))

	layer = NewWhittakerMapLayer(seed+1, layer, ClimateWarmWet)
	layer = NewWhittakerMapLayer(seed+1, layer, ClimateColdDry)
	layer = NewWhittakerMapLayer(seed+2, layer, ClimateLargerBiomes)

	for i := 0; i < 2; i++ {
		layer = NewZoomMapLayer(seed+int64(100+i), layer, ZoomBlurry)
	}

	for i := 0; i < 2; i++ {
		layer = NewErosionMapLayer(seed+int64(3+i), layer)
	}

	layer = NewDeepOceanMapLayer(seed+4, layer)

	layerMountains := MapLayer(layer)
	for i := 0; i < 2; i++ {
		layerMountains = NewZoomMapLayer(seed+int64(200+i), layerMountains, ZoomNormal)
	}

	layer = NewBiomeMapLayer(seed+5, layer)
	for i := 0; i < 2; i++ {
		layer = NewZoomMapLayer(seed+int64(200+i), layer, ZoomNormal)
	}

	layer = NewShoreMapLayer(seed+7, layer)

	for i := 0; i < zoom; i++ {
		layer = NewZoomMapLayer(seed+int64(500+i), layer, ZoomNormal)
	}

	layerRiver := MapLayer(layerMountains)
	layerRiver = NewZoomMapLayer(seed+300, layerRiver, ZoomNormal)
	layerRiver = NewZoomMapLayer(seed+400, layerRiver, ZoomNormal)
	for i := 0; i < zoom; i++ {
		layerRiver = NewZoomMapLayer(seed+int64(500+i), layerRiver, ZoomNormal)
	}
	layerRiver = NewRiverMapLayer(seed+10, layerRiver)

	layer = NewRiverMapLayerMerged(seed+1000, layerRiver, layer)

	layerLowerRes := layer

	for i := 0; i < 2; i++ {
		layer = NewZoomMapLayer(seed+int64(2000+i), layer, ZoomNormal)
	}

	layer = NewSmoothMapLayer(seed+1001, layer)

	return layer, layerLowerRes
}
