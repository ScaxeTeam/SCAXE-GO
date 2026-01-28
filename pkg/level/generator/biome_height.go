package generator

type BiomeHeight struct {
	Height float64
	Scale  float64
}

var (
	BiomeHeightDefault        = BiomeHeight{Height: 0.1, Scale: 0.2}
	BiomeHeightFlatShore      = BiomeHeight{Height: 0.0, Scale: 0.025}
	BiomeHeightHighPlateau    = BiomeHeight{Height: 1.5, Scale: 0.025}
	BiomeHeightFlatlands      = BiomeHeight{Height: 0.125, Scale: 0.05}
	BiomeHeightSwampland      = BiomeHeight{Height: -0.2, Scale: 0.1}
	BiomeHeightMidPlains      = BiomeHeight{Height: 0.2, Scale: 0.2}
	BiomeHeightFlatlandsHills = BiomeHeight{Height: 0.45, Scale: 0.3}
	BiomeHeightSwamplandHills = BiomeHeight{Height: -0.1, Scale: 0.3}
	BiomeHeightLowHills       = BiomeHeight{Height: 0.2, Scale: 0.3}
	BiomeHeightHills          = BiomeHeight{Height: 0.45, Scale: 0.3}
	BiomeHeightMidHills2      = BiomeHeight{Height: 0.3, Scale: 0.4}
	BiomeHeightDefaultHills   = BiomeHeight{Height: 0.2, Scale: 0.4}
	BiomeHeightMidHills       = BiomeHeight{Height: 0.25, Scale: 0.4}
	BiomeHeightBigHills       = BiomeHeight{Height: 0.55, Scale: 0.5}
	BiomeHeightBigHills2      = BiomeHeight{Height: 0.5, Scale: 0.5}
	BiomeHeightExtremeHills   = BiomeHeight{Height: 1.0, Scale: 0.5}
	BiomeHeightRockyShore     = BiomeHeight{Height: 0.1, Scale: 0.8}
	BiomeHeightLowSpikes      = BiomeHeight{Height: 0.3625, Scale: 1.225}
	BiomeHeightHighSpikes     = BiomeHeight{Height: 1.05, Scale: 1.2125}

	BiomeHeightRiver     = BiomeHeight{Height: -0.5, Scale: 0.0}
	BiomeHeightOcean     = BiomeHeight{Height: -1.0, Scale: 0.1}
	BiomeHeightDeepOcean = BiomeHeight{Height: -1.8, Scale: 0.1}
)

var BiomeHeightMap = map[int]BiomeHeight{

	0:  BiomeHeightOcean,
	24: BiomeHeightDeepOcean,
	10: BiomeHeightOcean,

	7:  BiomeHeightRiver,
	11: BiomeHeightRiver,

	16: BiomeHeightFlatShore,
	26: BiomeHeightFlatShore,
	25: BiomeHeightRockyShore,

	1:  BiomeHeightDefault,
	2:  BiomeHeightFlatlands,
	12: BiomeHeightFlatlands,
	35: BiomeHeightFlatlands,

	4:  BiomeHeightDefault,
	5:  BiomeHeightMidPlains,
	27: BiomeHeightDefault,
	29: BiomeHeightDefault,
	30: BiomeHeightMidPlains,
	32: BiomeHeightMidPlains,

	3:  BiomeHeightHills,
	13: BiomeHeightHills,
	17: BiomeHeightHills,
	18: BiomeHeightHills,
	19: BiomeHeightHills,
	28: BiomeHeightHills,
	31: BiomeHeightHills,
	33: BiomeHeightHills,
	22: BiomeHeightHills,

	34: BiomeHeightExtremeHills,

	6: BiomeHeightSwampland,

	21: BiomeHeightDefault,
	23: BiomeHeightDefault,

	14: BiomeHeightLowHills,
	15: BiomeHeightFlatShore,

	36: BiomeHeightHighPlateau,

	37: BiomeHeightHighPlateau,
	38: BiomeHeightHighPlateau,
	39: BiomeHeightHighPlateau,
}

func GetBiomeHeight(biomeID int) BiomeHeight {
	if h, ok := BiomeHeightMap[biomeID]; ok {
		return h
	}
	return BiomeHeightDefault
}
