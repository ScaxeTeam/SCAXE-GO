package biome

var RegisteredBiomes = make(map[uint8]Biome)

func Register(id uint8, b Biome) {
	RegisteredBiomes[id] = b
}

func GetBiome(id uint8) Biome {
	if b, ok := RegisteredBiomes[id]; ok {
		return b
	}
	return RegisteredBiomes[0]
}

func InitBiomes() {

	Register(0, &BaseBiome{
		ID: 0, Name: "Ocean",
		BaseHeight: -1.0, HeightVariation: 0.1,
		Temperature: 0.5, Rainfall: 0.5,
	})

	Register(1, NewPlainsBiome())

	Register(2, NewDesertBiome())

	Register(3, NewExtremeHillsBiome())

	Register(4, NewForestBiome(FOREST_NORMAL))

	Register(5, NewTaigaBiome())

	Register(6, NewSwampBiome())

	Register(7, &BaseBiome{
		ID: 7, Name: "River",
		BaseHeight: -0.5, HeightVariation: 0.0,
		Temperature: 0.5, Rainfall: 0.5,
		Decorator: NewDecorator(),
	})

	Register(8, &BaseBiome{
		ID: 8, Name: "Hell",
		BaseHeight: 0.1, HeightVariation: 0.2,
		Temperature: 2.0, Rainfall: 0.0,
		Decorator: NewDecorator(),
	})

	Register(9, &BaseBiome{
		ID: 9, Name: "The End",
		BaseHeight: 0.1, HeightVariation: 0.2,
		Temperature: 0.5, Rainfall: 0.5,
		Decorator: NewDecorator(),
	})

	Register(10, &BaseBiome{
		ID: 10, Name: "Frozen Ocean",
		BaseHeight: -1.0, HeightVariation: 0.1,
		Temperature: 0.0, Rainfall: 0.5,
		Decorator: NewDecorator(),
	})

	Register(11, &BaseBiome{
		ID: 11, Name: "Frozen River",
		BaseHeight: -0.5, HeightVariation: 0.0,
		Temperature: 0.0, Rainfall: 0.5,
		Decorator: NewDecorator(),
	})

	Register(12, NewIcePlainsBiome())

	Register(13, &BaseBiome{
		ID: 13, Name: "Ice Mountains",
		BaseHeight: 0.45, HeightVariation: 0.3,
		Temperature: 0.0, Rainfall: 0.5,
		Decorator: NewDecorator(),
	})

	Register(14, NewMushroomIslandBiome(MUSHROOM_ISLAND))

	Register(15, NewMushroomIslandBiome(MUSHROOM_ISLAND_SHORE))

	Register(16, NewBeachBiome())

	Register(17, &BaseBiome{
		ID: 17, Name: "Desert Hills",
		BaseHeight: 0.45, HeightVariation: 0.3,
		Temperature: 2.0, Rainfall: 0.0,
		Decorator: NewDecorator(),
	})

	Register(18, &BaseBiome{
		ID: 18, Name: "Forest Hills",
		BaseHeight: 0.45, HeightVariation: 0.3,
		Temperature: 0.7, Rainfall: 0.8,
		Decorator: NewDecorator(),
	})

	Register(19, &BaseBiome{
		ID: 19, Name: "Taiga Hills",
		BaseHeight: 0.45, HeightVariation: 0.3,
		Temperature: 0.25, Rainfall: 0.8,
		Decorator: NewDecorator(),
	})

	Register(20, NewExtremeHillsEdgeBiome())

	Register(21, NewJungleBiome())

	Register(22, NewJungleHillsBiome())

	Register(23, NewJungleEdgeBiome())

	Register(24, &BaseBiome{
		ID: 24, Name: "Deep Ocean",
		BaseHeight: -1.8, HeightVariation: 0.1,
		Temperature: 0.5, Rainfall: 0.5,
		Decorator: NewDecorator(),
	})

	Register(25, NewStoneBeachBiome())

	Register(26, NewColdBeachBiome())

	Register(27, NewForestBiome(FOREST_BIRCH))

	Register(28, &BaseBiome{
		ID: 28, Name: "Birch Forest Hills",
		BaseHeight: 0.45, HeightVariation: 0.3,
		Temperature: 0.6, Rainfall: 0.6,
		Decorator: NewDecorator(),
	})

	Register(29, NewRoofedForestBiome())

	Register(30, NewColdTaigaBiome())

	Register(31, NewColdTaigaHillsBiome())

	Register(32, NewMegaTaigaBiome())

	Register(33, NewMegaTaigaHillsBiome())

	Register(34, NewExtremeHillsPlusBiome())

	Register(35, NewSavannaBiome())

	Register(36, NewSavannaPlateauBiome())

	Register(37, NewMesaBiome())

	Register(38, NewMesaPlateauFBiome())

	Register(39, NewMesaPlateauBiome())

	Register(129, NewSunflowerPlainsBiome())

	Register(130, &BaseBiome{
		ID: 130, Name: "Desert M",
		BaseHeight: 0.225, HeightVariation: 0.25,
		Temperature: 2.0, Rainfall: 0.0,
	})

	Register(140, NewIcePlainsSpikesBiome())
}
