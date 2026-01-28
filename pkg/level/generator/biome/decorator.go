package biome

import (
	"github.com/scaxe/scaxe-go/pkg/block"

	"github.com/scaxe/scaxe-go/pkg/level/generator/gorigional/noise"
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type NoiseGenerator interface {
	GetValue(x, z float64) float64
}

type Decorator struct {
	TreesPerChunk         int
	ExtraTreeChance       float64
	FlowersPerChunk       int
	GrassPerChunk         int
	DeadBushPerChunk      int
	MushroomsPerChunk     int
	ReedsPerChunk         int
	CactiPerChunk         int
	WaterlilyPerChunk     int
	GravelPatchesPerChunk int
	SandPatchesPerChunk   int
	ClayPerChunk          int
	BigMushroomsPerChunk  int
	GenerateFalls         bool

	DirtGen *object.Ore

	FlowerNoise *noise.PerlinNoiseGenerator

	GravelOreGen *object.Ore
	GraniteGen   *object.Ore
	DioriteGen   *object.Ore
	AndesiteGen  *object.Ore
	CoalGen      *object.Ore
	IronGen      *object.Ore
	GoldGen      *object.Ore
	RedstoneGen  *object.Ore
	DiamondGen   *object.Ore
	LapisGen     *object.Ore

	SandGen        *object.Sand
	GravelGen      *object.Sand
	ClayGen        *object.Clay
	MushroomBrGen  *object.Bush
	MushroomRdGen  *object.Bush
	FlowerYGen     *object.Bush
	FlowerRGen     *object.Bush
	GrassGen       *object.Grass
	ReedGen        *object.Reed
	CactusGen      *object.Cactus
	WaterlilyGen   *object.Waterlily
	PumpkinGen     *object.Pumpkin
	DeadBushGen    *object.DeadBush
	WaterSpringGen *object.Spring
	LavaSpringGen  *object.Spring
}

type BiomeProvider interface {
	GetTreeFeature(r *rand.Random) Generator
	GetFlowerType(r *rand.Random, pos world.BlockPos) Generator
}

func NewDecorator() *Decorator {
	d := &Decorator{
		TreesPerChunk:         0,
		ExtraTreeChance:       0.1,
		FlowersPerChunk:       2,
		GrassPerChunk:         1,
		DeadBushPerChunk:      0,
		MushroomsPerChunk:     0,
		ReedsPerChunk:         0,
		CactiPerChunk:         0,
		WaterlilyPerChunk:     0,
		GravelPatchesPerChunk: 1,
		SandPatchesPerChunk:   3,
		ClayPerChunk:          1,
		BigMushroomsPerChunk:  0,
		GenerateFalls:         true,
	}
	d.InitOres()

	d.SandGen = object.NewSand(block.SAND, 7)
	d.GravelGen = object.NewSand(block.GRAVEL, 6)
	d.ClayGen = object.NewClay(4)

	d.MushroomBrGen = object.NewBush(block.BROWN_MUSHROOM, 0)
	d.MushroomRdGen = object.NewBush(block.RED_MUSHROOM, 0)
	d.FlowerYGen = object.NewBush(block.DANDELION, 0)
	d.FlowerRGen = object.NewBush(block.RED_FLOWER, 0)
	d.GrassGen = object.NewGrass(1)
	d.ReedGen = object.NewReed()
	d.CactusGen = object.NewCactus()
	d.WaterlilyGen = object.NewWaterlily()
	d.PumpkinGen = object.NewPumpkin()
	d.DeadBushGen = object.NewDeadBush()
	d.WaterSpringGen = object.NewSpring(block.WATER)
	d.LavaSpringGen = object.NewSpring(block.LAVA)

	return d
}

func (d *Decorator) InitOres() {

	d.DirtGen = object.NewOre(block.DIRT, 0, 33)
	d.GravelOreGen = object.NewOre(block.GRAVEL, 0, 33)
	d.GraniteGen = object.NewOre(block.STONE, 1, 33)
	d.DioriteGen = object.NewOre(block.STONE, 3, 33)
	d.AndesiteGen = object.NewOre(block.STONE, 5, 33)
	d.CoalGen = object.NewOre(block.COAL_ORE, 0, 17)
	d.IronGen = object.NewOre(block.IRON_ORE, 0, 9)
	d.GoldGen = object.NewOre(block.GOLD_ORE, 0, 9)
	d.RedstoneGen = object.NewOre(block.REDSTONE_ORE, 0, 8)
	d.DiamondGen = object.NewOre(block.DIAMOND_ORE, 0, 8)
	d.LapisGen = object.NewOre(block.LAPIS_ORE, 0, 7)
}

func (d *Decorator) Decorate(w populator.ChunkManager, r *rand.Random, b BiomeProvider, pos world.BlockPos) {

	d.GenStandardOre1(w, r, 10, d.DirtGen, 0, 256, pos)
	d.GenStandardOre1(w, r, 8, d.GravelOreGen, 0, 256, pos)
	d.GenStandardOre1(w, r, 10, d.DioriteGen, 0, 80, pos)
	d.GenStandardOre1(w, r, 10, d.GraniteGen, 0, 80, pos)
	d.GenStandardOre1(w, r, 10, d.AndesiteGen, 0, 80, pos)
	d.GenStandardOre1(w, r, 20, d.CoalGen, 0, 128, pos)
	d.GenStandardOre1(w, r, 20, d.IronGen, 0, 64, pos)
	d.GenStandardOre1(w, r, 2, d.GoldGen, 0, 32, pos)
	d.GenStandardOre1(w, r, 8, d.RedstoneGen, 0, 16, pos)
	d.GenStandardOre1(w, r, 1, d.DiamondGen, 0, 16, pos)
	d.GenStandardOre2(w, r, 1, d.LapisGen, 16, 16, pos)

	for i := 0; i < d.SandPatchesPerChunk; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		d.SandGen.Generate(w, r, pos.Add(int32(x), terrainHeight, int32(z)))
	}

	for i := 0; i < d.ClayPerChunk; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		d.ClayGen.Generate(w, r, pos.Add(int32(x), terrainHeight, int32(z)))
	}

	for i := 0; i < d.GravelPatchesPerChunk; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		d.GravelGen.Generate(w, r, pos.Add(int32(x), terrainHeight, int32(z)))
	}

	trees := d.TreesPerChunk
	if r.NextFloat() < d.ExtraTreeChance {
		trees++
	}
	for i := 0; i < trees; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8

		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))

		feature := b.GetTreeFeature(r)
		if feature != nil {
			feature.Generate(w, r, pos.Add(int32(x), terrainHeight, int32(z)))
		}
	}

	for i := 0; i < d.BigMushroomsPerChunk; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		bigMushroom := object.NewBigMushroom(0)
		bigMushroom.Generate(w, r, pos.Add(int32(x), terrainHeight, int32(z)))
	}

	const flowerNoiseScale = 41.66666667

	for i := 0; i < d.FlowersPerChunk; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		yMax := int(terrainHeight) + 32

		if yMax > 0 {
			y := r.NextBoundedInt(yMax)
			blockPos := pos.Add(int32(x), int32(y), int32(z))

			gen := b.GetFlowerType(r, blockPos)
			if gen != nil {
				gen.Generate(w, r, blockPos)
			}
		}
	}

	for i := 0; i < d.GrassPerChunk; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8

		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		yMax := int(terrainHeight) * 2
		if yMax > 0 {
			y := r.NextBoundedInt(yMax)

			isDouble := false
			if d.FlowerNoise != nil {

				nv := d.FlowerNoise.GetValue(float64(pos.X()+int32(x))/80.0, float64(pos.Z()+int32(z))/80.0)
				if nv > 0.4 {
					isDouble = true
				}
			}

			if isDouble {

				gen := object.NewDoublePlant(object.DoublePlantGrass)
				gen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
			} else {

				d.GrassGen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
			}
		}
	}

	for i := 0; i < d.DeadBushPerChunk; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		yMax := int(terrainHeight) * 2
		if yMax > 0 {
			y := r.NextBoundedInt(yMax)
			d.DeadBushGen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
		}
	}

	for i := 0; i < d.WaterlilyPerChunk; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		yMax := int(terrainHeight) * 2
		if yMax > 0 {
			y := r.NextBoundedInt(yMax)
			d.WaterlilyGen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
		}
	}

	for i := 0; i < d.MushroomsPerChunk; i++ {
		if r.NextBoundedInt(4) == 0 {
			x := r.NextBoundedInt(16) + 8
			z := r.NextBoundedInt(16) + 8
			terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
			d.MushroomBrGen.Generate(w, r, pos.Add(int32(x), terrainHeight, int32(z)))
		}
		if r.NextBoundedInt(8) == 0 {
			x := r.NextBoundedInt(16) + 8
			z := r.NextBoundedInt(16) + 8
			terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
			yMax := int(terrainHeight) * 2
			if yMax > 0 {
				y := r.NextBoundedInt(yMax)
				d.MushroomRdGen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
			}
		}
	}

	if r.NextBoundedInt(4) == 0 {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		yMax := int(terrainHeight) * 2
		if yMax > 0 {
			y := r.NextBoundedInt(yMax)
			d.MushroomBrGen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
		}
	}
	if r.NextBoundedInt(8) == 0 {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		yMax := int(terrainHeight) * 2
		if yMax > 0 {
			y := r.NextBoundedInt(yMax)
			d.MushroomRdGen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
		}
	}

	for i := 0; i < d.ReedsPerChunk; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		yMax := int(terrainHeight) * 2
		if yMax > 0 {
			y := r.NextBoundedInt(yMax)
			d.ReedGen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
		}
	}

	for i := 0; i < 10; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		yMax := int(terrainHeight) * 2
		if yMax > 0 {
			y := r.NextBoundedInt(yMax)
			d.ReedGen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
		}
	}

	if r.NextBoundedInt(32) == 0 {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		yMax := int(terrainHeight) * 2
		if yMax > 0 {
			y := r.NextBoundedInt(yMax)
			d.PumpkinGen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
		}
	}

	for i := 0; i < d.CactiPerChunk; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8
		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		yMax := int(terrainHeight) * 2
		if yMax > 0 {
			y := r.NextBoundedInt(yMax)
			d.CactusGen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
		}
	}

	if d.GenerateFalls {
		for i := 0; i < 50; i++ {
			p := pos.Add(int32(r.NextBoundedInt(16)+8), int32(r.NextBoundedInt(r.NextBoundedInt(248)+8)), int32(r.NextBoundedInt(16)+8))
			d.WaterSpringGen.Generate(w, r, p)
		}
		for i := 0; i < 20; i++ {
			p := pos.Add(int32(r.NextBoundedInt(16)+8), int32(r.NextBoundedInt(r.NextBoundedInt(r.NextBoundedInt(240)+8)+8)), int32(r.NextBoundedInt(16)+8))
			d.LavaSpringGen.Generate(w, r, p)
		}
	}
}

func (d *Decorator) GenStandardBush(w populator.ChunkManager, r *rand.Random, count int, gen Generator, pos world.BlockPos) {
	for i := 0; i < count; i++ {
		x := r.NextBoundedInt(16) + 8
		z := r.NextBoundedInt(16) + 8

		terrainHeight := w.GetHeight(pos.X()+int32(x), pos.Z()+int32(z))
		yMax := int(terrainHeight) + 32
		if yMax > 0 {
			y := r.NextBoundedInt(yMax)
			gen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
		}
	}
}

func (d *Decorator) GenSand(w populator.ChunkManager, r *rand.Random, count int, gen *object.Sand, maxH int, pos world.BlockPos) {
	for i := 0; i < count; i++ {
		x := r.NextBoundedInt(16)
		y := r.NextBoundedInt(maxH)
		z := r.NextBoundedInt(16)
		gen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
	}
}

func (d *Decorator) GenClay(w populator.ChunkManager, r *rand.Random, count int, gen *object.Clay, pos world.BlockPos) {
	for i := 0; i < count; i++ {
		x := r.NextBoundedInt(16)
		y := r.NextBoundedInt(128)
		z := r.NextBoundedInt(16)
		gen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
	}
}

func (d *Decorator) GenStandardOre1(w populator.ChunkManager, r *rand.Random, count int, gen *object.Ore, minH, maxH int, pos world.BlockPos) {
	if maxH < minH {
		minH, maxH = maxH, minH
	} else if maxH == minH {
		if minH < 255 {
			maxH++
		} else {
			minH--
		}
	}

	for i := 0; i < count; i++ {

		x := r.NextBoundedInt(16)
		y := r.NextBoundedInt(maxH-minH) + minH
		z := r.NextBoundedInt(16)

		gen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
	}
}

func (d *Decorator) GenStandardOre2(w populator.ChunkManager, r *rand.Random, count int, gen *object.Ore, centerH, spread int, pos world.BlockPos) {
	for i := 0; i < count; i++ {
		x := r.NextBoundedInt(16)
		y := r.NextBoundedInt(spread) + r.NextBoundedInt(spread) + centerH - spread
		z := r.NextBoundedInt(16)
		gen.Generate(w, r, pos.Add(int32(x), int32(y), int32(z)))
	}
}
