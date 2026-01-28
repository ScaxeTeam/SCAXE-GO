package biome

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type MesaBiome struct {
	*BaseBiome
	HasTrees  bool
	IsPlateau bool
	clayBands []byte
}

func NewMesaBiome() *MesaBiome {
	b := &MesaBiome{
		BaseBiome: &BaseBiome{
			ID:              MESA,
			Name:            "Mesa",
			BaseHeight:      0.1,
			HeightVariation: 0.2,
			Temperature:     2.0,
			Rainfall:        0.0,
			GroundCover:     []block.BlockState{block.NewBlockState(block.RED_SANDSTONE, 0), block.NewBlockState(block.STAINED_CLAY, 1)},
			Decorator:       NewDecorator(),
		},
		HasTrees:  false,
		IsPlateau: false,
	}
	b.Decorator.DeadBushPerChunk = 20
	b.Decorator.CactiPerChunk = 5
	b.Decorator.FlowersPerChunk = 0
	b.Decorator.TreesPerChunk = -999
	b.initClayBands()
	return b
}

func NewMesaPlateauBiome() *MesaBiome {
	b := &MesaBiome{
		BaseBiome: &BaseBiome{
			ID:              MESA_PLATEAU,
			Name:            "Mesa Plateau",
			BaseHeight:      1.5,
			HeightVariation: 0.025,
			Temperature:     2.0,
			Rainfall:        0.0,
			GroundCover:     []block.BlockState{block.NewBlockState(block.RED_SANDSTONE, 0), block.NewBlockState(block.STAINED_CLAY, 1)},
			Decorator:       NewDecorator(),
		},
		HasTrees:  false,
		IsPlateau: true,
	}
	b.Decorator.DeadBushPerChunk = 20
	b.Decorator.CactiPerChunk = 5
	b.initClayBands()
	return b
}

func NewMesaPlateauFBiome() *MesaBiome {
	b := &MesaBiome{
		BaseBiome: &BaseBiome{
			ID:              MESA_PLATEAU_F,
			Name:            "Mesa Plateau F",
			BaseHeight:      1.5,
			HeightVariation: 0.025,
			Temperature:     2.0,
			Rainfall:        0.0,
			GroundCover:     []block.BlockState{block.NewBlockState(block.RED_SANDSTONE, 0), block.NewBlockState(block.STAINED_CLAY, 1)},
			Decorator:       NewDecorator(),
		},
		HasTrees:  true,
		IsPlateau: true,
	}
	b.Decorator.DeadBushPerChunk = 20
	b.Decorator.CactiPerChunk = 5
	b.Decorator.TreesPerChunk = 5
	b.initClayBands()
	return b
}

func (b *MesaBiome) initClayBands() {

	b.initClayBandsWithSeed(12345)
}

func (b *MesaBiome) initClayBandsWithSeed(seed int64) {
	b.clayBands = make([]byte, 64)
	r := rand.NewRandom(seed)

	for i := range b.clayBands {
		b.clayBands[i] = 0
	}

	for l := 0; l < 64; {
		l += r.NextBoundedInt(5) + 1
		if l < 64 {
			b.clayBands[l] = 1
		}
	}

	yellowCount := r.NextBoundedInt(4) + 2
	for i := 0; i < yellowCount; i++ {
		width := r.NextBoundedInt(3) + 1
		start := r.NextBoundedInt(64)
		for j := 0; j < width && start+j < 64; j++ {
			b.clayBands[start+j] = 4
		}
	}

	brownCount := r.NextBoundedInt(4) + 2
	for i := 0; i < brownCount; i++ {
		width := r.NextBoundedInt(3) + 2
		start := r.NextBoundedInt(64)
		for j := 0; j < width && start+j < 64; j++ {
			b.clayBands[start+j] = 12
		}
	}

	redCount := r.NextBoundedInt(4) + 2
	for i := 0; i < redCount; i++ {
		width := r.NextBoundedInt(3) + 1
		start := r.NextBoundedInt(64)
		for j := 0; j < width && start+j < 64; j++ {
			b.clayBands[start+j] = 14
		}
	}

	whiteCount := r.NextBoundedInt(3) + 3
	pos := 0
	for i := 0; i < whiteCount; i++ {
		pos += r.NextBoundedInt(16) + 4
		if pos >= 64 {
			break
		}
		b.clayBands[pos] = 0

		b.clayBands[pos] = 8

		if pos > 0 && r.NextBoundedInt(2) == 0 {
			b.clayBands[pos-1] = 7
		}
		if pos < 63 && r.NextBoundedInt(2) == 0 {
			b.clayBands[pos+1] = 7
		}
	}
}

func (b *MesaBiome) GenTerrainBlocks(chunk *world.Chunk, r *rand.Random, x, z int, noiseVal float64) {

	seaLevel := 63
	chunkX := x & 15
	chunkZ := z & 15

	topBlock := byte(block.RED_SANDSTONE)
	fillerBlock := byte(block.STAINED_CLAY)
	stone := byte(1)
	air := byte(0)

	_ = topBlock
	_ = fillerBlock

	run := -1
	clayBandIndex := 0

	for y := 255; y >= 0; y-- {
		id, _ := chunk.GetBlock(chunkX, y, chunkZ)

		if id == air {
			run = -1
		} else if id == stone {
			if run == -1 {
				run = 3 + int(math.Abs(noiseVal))

				if y >= seaLevel-1 {
					if y >= seaLevel+2+int(noiseVal*3) {

						chunk.SetBlock(chunkX, y, chunkZ, block.SAND, 1)
					} else {

						chunk.SetBlock(chunkX, y, chunkZ, block.HARDENED_CLAY, 0)
					}
				} else {
					chunk.SetBlock(chunkX, y, chunkZ, block.STAINED_CLAY, 1)
				}
			} else if run > 0 {
				run--

				clayBandIndex = y % 64
				clayMeta := b.clayBands[clayBandIndex]
				chunk.SetBlock(chunkX, y, chunkZ, block.STAINED_CLAY, clayMeta)
			}
		}
	}
}

func (b *MesaBiome) GetTreeFeature(r *rand.Random) Generator {

	if b.HasTrees {
		return nil
	}
	return nil
}
