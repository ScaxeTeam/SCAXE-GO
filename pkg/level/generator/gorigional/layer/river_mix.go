package layer

import "github.com/scaxe/scaxe-go/pkg/level/generator/biome"

type GenLayerRiverMix struct {
	*BaseLayer
	biomePatternGeneratorChain GenLayer
	riverPatternGeneratorChain GenLayer
}

func NewGenLayerRiverMix(baseSeed int64, biomeChain, riverChain GenLayer) *GenLayerRiverMix {
	l := &GenLayerRiverMix{
		BaseLayer: NewBaseLayer(baseSeed),
	}
	l.biomePatternGeneratorChain = biomeChain
	l.riverPatternGeneratorChain = riverChain
	return l
}

func (l *GenLayerRiverMix) InitWorldGenSeed(seed int64) {
	l.BaseLayer.InitWorldGenSeed(seed)
	l.biomePatternGeneratorChain.InitWorldGenSeed(seed)
	l.riverPatternGeneratorChain.InitWorldGenSeed(seed)
}

func (l *GenLayerRiverMix) GetInts(x, z, width, depth int) []int {
	biomeInts := l.biomePatternGeneratorChain.GetInts(x, z, width, depth)
	riverInts := l.riverPatternGeneratorChain.GetInts(x, z, width, depth)
	output := make([]int, width*depth)

	for i := 0; i < width*depth; i++ {
		b := biomeInts[i]
		r := riverInts[i]

		if b != int(biome.OCEAN) && b != int(biome.DEEP_OCEAN) {

			if r == int(biome.RIVER) {
				if b == int(biome.ICE_PLAINS) {
					output[i] = int(biome.FROZEN_RIVER)
				} else if b != int(biome.MUSHROOM_ISLAND) && b != int(biome.MUSHROOM_ISLAND_SHORE) {
					output[i] = r & 255
				} else {
					output[i] = int(biome.MUSHROOM_ISLAND_SHORE)
				}
			} else {
				output[i] = b
			}
		} else {
			output[i] = b
		}
	}
	return output
}
