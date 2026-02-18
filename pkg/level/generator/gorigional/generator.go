package gorigional

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/level/generator"
	"github.com/scaxe/scaxe-go/pkg/level/generator/biome"
	"github.com/scaxe/scaxe-go/pkg/level/generator/gorigional/layer"
	"github.com/scaxe/scaxe-go/pkg/level/generator/gorigional/noise"
	"github.com/scaxe/scaxe-go/pkg/level/generator/gorigional/structure"
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type ChunkGeneratorOverworld struct {
	seed int64

	minLimitPerlinNoise *noise.OctavesNoise
	maxLimitPerlinNoise *noise.OctavesNoise
	mainPerlinNoise     *noise.OctavesNoise
	surfaceNoise        *noise.PerlinSimplexGenerator
	scaleNoise          *noise.OctavesNoise
	depthNoise          *noise.OctavesNoise
	forestNoise         *noise.OctavesNoise
	flowerNoise         *noise.PerlinNoiseGenerator

	genLayer layer.GenLayer
	biomeGen layer.GenLayer

	caves            *structure.MapGenCaves
	ravines          *structure.MapGenRavine
	village          *structure.MapGenVillage
	stronghold       *structure.MapGenStronghold
	mineshaft        *structure.MapGenMineshaft
	scatteredFeature *structure.MapGenScatteredFeature

	biomeWeights []float64

	settings *ChunkGeneratorSettings
}

type ChunkGeneratorSettings struct {
	coordinateScale         float64
	heightScale             float64
	lowerLimitScale         float64
	upperLimitScale         float64
	depthNoiseScaleX        float64
	depthNoiseScaleZ        float64
	depthNoiseScaleExponent float64
	mainNoiseScaleX         float64
	mainNoiseScaleY         float64
	mainNoiseScaleZ         float64
	baseSize                float64
	stretchY                float64
	biomeDepthOffSet        float64
	biomeDepthWeight        float64
	biomeScaleOffset        float64
	biomeScaleWeight        float64
	seaLevel                int
	ChunkManager            generator.ChunkManager

	useCaves        bool
	useDungeons     bool
	dungeonChance   int
	useStrongholds  bool
	useVillages     bool
	useMineShafts   bool
	useTemples      bool
	useMonuments    bool
	useMansions     bool
	useRavines      bool
	useWaterLakes   bool
	waterLakeChance int
	useLavaLakes    bool
	lavaLakeChance  int
	useLavaOceans   bool
	fixedBiome      int
	biomeSize       int
	riverSize       int
	MaxHeight       int
}

func DefaultSettings() *ChunkGeneratorSettings {
	return &ChunkGeneratorSettings{
		coordinateScale:         684.412,
		heightScale:             684.412,
		lowerLimitScale:         512.0,
		upperLimitScale:         512.0,
		depthNoiseScaleX:        200.0,
		depthNoiseScaleZ:        200.0,
		depthNoiseScaleExponent: 0.5,
		mainNoiseScaleX:         80.0,
		mainNoiseScaleY:         160.0,
		mainNoiseScaleZ:         80.0,
		baseSize:                8.5,
		stretchY:                12.0,
		biomeDepthOffSet:        0.0,
		biomeDepthWeight:        1.0,
		biomeScaleOffset:        0.0,
		biomeScaleWeight:        1.0,
		seaLevel:                63,
		useCaves:                true,
		useDungeons:             true,
		dungeonChance:           8,
		useStrongholds:          true,
		useVillages:             true,
		useMineShafts:           true,
		useTemples:              true,
		useMonuments:            true,
		useMansions:             true,
		useRavines:              true,
		useWaterLakes:           true,
		waterLakeChance:         80,
		useLavaLakes:            true,
		lavaLakeChance:          80,
		useLavaOceans:           false,
		fixedBiome:              -1,
		biomeSize:               4,
		riverSize:               4,
		MaxHeight:               128,
	}
}

func NewChunkGeneratorOverworld(seed int64) *ChunkGeneratorOverworld {
	rnd := rand.NewRandom(seed)

	layers := layer.InitializeAll(seed)
	genLayer := layers[0]
	biomeGen := layers[1]
	biomeSource := &genLayerAdapter{layer: genLayer}

	g := &ChunkGeneratorOverworld{
		seed:     seed,
		settings: DefaultSettings(),

		minLimitPerlinNoise: noise.NewOctavesNoise(rnd, 16),
		maxLimitPerlinNoise: noise.NewOctavesNoise(rnd, 16),
		mainPerlinNoise:     noise.NewOctavesNoise(rnd, 8),

		surfaceNoise: noise.NewPerlinSimplexGenerator(rnd, 4),
		scaleNoise:   noise.NewOctavesNoise(rnd, 10),

		depthNoise:  noise.NewOctavesNoise(rnd, 16),
		forestNoise: noise.NewOctavesNoise(rnd, 8),
		flowerNoise: noise.NewPerlinNoiseGenerator(rnd, 1),

		biomeWeights: make([]float64, 25),

		genLayer: genLayer,
		biomeGen: biomeGen,

		caves:            structure.NewMapGenCaves(seed),
		ravines:          structure.NewMapGenRavine(seed),
		village:          structure.NewMapGenVillage(seed),
		stronghold:       structure.NewMapGenStronghold(seed),
		mineshaft:        structure.NewMapGenMineshaft(seed),
		scatteredFeature: structure.NewMapGenScatteredFeature(seed, biomeSource),
	}

	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			f := 10.0 / float64(math.Sqrt(float64(i*i+j*j)+0.2))
			g.biomeWeights[i+2+(j+2)*5] = f
		}
	}

	g.caves.MaxHeight = g.settings.MaxHeight
	g.ravines.MaxHeight = g.settings.MaxHeight

	return g
}

func (g *ChunkGeneratorOverworld) Init(level generator.ChunkManager, seed int64) {
	g.seed = seed
	rnd := rand.NewRandom(seed)
	g.settings.ChunkManager = level

	layers := layer.InitializeAll(seed)
	g.genLayer = layers[0]
	g.biomeGen = layers[1]

	biomeSource := &genLayerAdapter{layer: g.genLayer}
	g.scatteredFeature = structure.NewMapGenScatteredFeature(seed, biomeSource)

	g.minLimitPerlinNoise = noise.NewOctavesNoise(rnd, 16)
	g.maxLimitPerlinNoise = noise.NewOctavesNoise(rnd, 16)
	g.mainPerlinNoise = noise.NewOctavesNoise(rnd, 8)
	g.surfaceNoise = noise.NewPerlinSimplexGenerator(rnd, 4)
	g.scaleNoise = noise.NewOctavesNoise(rnd, 10)
	g.depthNoise = noise.NewOctavesNoise(rnd, 16)
	g.forestNoise = noise.NewOctavesNoise(rnd, 8)
	g.flowerNoise = noise.NewPerlinNoiseGenerator(rnd, 1)
}

func (g *ChunkGeneratorOverworld) GetName() string {
	return "gorigional"
}

func (g *ChunkGeneratorOverworld) GetSettings() map[string]interface{} {
	return nil
}

func (g *ChunkGeneratorOverworld) GetSpawn() *world.Vector3 {

	spawnX := int32(0)
	spawnZ := int32(0)

	if g.settings != nil && g.settings.ChunkManager != nil {

		chunk := g.settings.ChunkManager.GetChunk(0, 0, true)
		if chunk != nil {

			for y := 255; y >= 0; y-- {
				blockId := chunk.GetBlockId(0, y, 0)

				if blockId == 0 || blockId == 8 || blockId == 9 || blockId == 10 || blockId == 11 || blockId == 18 || blockId == 161 {
					continue
				}

				hasSpace := true
				for check := 1; check <= 2; check++ {
					aboveId := chunk.GetBlockId(0, y+check, 0)

					if aboveId != 0 {
						hasSpace = false
						break
					}
				}
				if hasSpace {
					return world.NewVector3(float64(spawnX), float64(y+1), float64(spawnZ))
				}

			}
		}
	}

	return world.NewVector3(0, 65, 0)
}

func (g *ChunkGeneratorOverworld) generateHeightmap(x, y, z int, heightMap, depthRegion, mainNoiseRegion, minLimitRegion, maxLimitRegion []float64, biomesForGeneration []int) {
	const horizontalPoints = 5

	verticalPoints := g.settings.MaxHeight/8 + 1

	depthRegion = g.depthNoise.GenerateNoiseOctaves2D(depthRegion, x, z, horizontalPoints, horizontalPoints, g.settings.depthNoiseScaleX, g.settings.depthNoiseScaleZ)

	f := g.settings.coordinateScale
	f1 := g.settings.heightScale

	mainNoiseRegion = g.mainPerlinNoise.GenerateNoiseOctaves(mainNoiseRegion, x, y, z, horizontalPoints, verticalPoints, horizontalPoints, f/g.settings.mainNoiseScaleX, f1/g.settings.mainNoiseScaleY, f/g.settings.mainNoiseScaleZ)
	minLimitRegion = g.minLimitPerlinNoise.GenerateNoiseOctaves(minLimitRegion, x, y, z, horizontalPoints, verticalPoints, horizontalPoints, f, f1, f)
	maxLimitRegion = g.maxLimitPerlinNoise.GenerateNoiseOctaves(maxLimitRegion, x, y, z, horizontalPoints, verticalPoints, horizontalPoints, f, f1, f)

	i := 0
	j := 0

	for k := 0; k < horizontalPoints; k++ {
		for l := 0; l < horizontalPoints; l++ {
			f2 := 0.0
			f3 := 0.0
			f4 := 0.0

			centerBiomeID := biomesForGeneration[k+2+(l+2)*10]
			centerBiome := biome.GetBiome(uint8(centerBiomeID))
			centerBaseHeight := centerBiome.GetMinElevation()

			for j1 := -2; j1 <= 2; j1++ {
				for k1 := -2; k1 <= 2; k1++ {

					biomeID := biomesForGeneration[k+j1+2+(l+k1+2)*10]
					neighborBiome := biome.GetBiome(uint8(biomeID))
					baseHeight := neighborBiome.GetMinElevation()
					heightVar := neighborBiome.GetMaxElevation()

					f5 := g.settings.biomeDepthOffSet + baseHeight*g.settings.biomeDepthWeight
					f6 := g.settings.biomeScaleOffset + heightVar*g.settings.biomeScaleWeight

					f7 := g.biomeWeights[j1+2+(k1+2)*5] / (f5 + 2.0)
					if baseHeight > centerBaseHeight {
						f7 /= 2.0
					}

					f2 += f6 * f7
					f3 += f5 * f7
					f4 += f7
				}
			}

			f2 /= f4
			f3 /= f4
			f2 = f2*0.9 + 0.1
			f3 = (f3*4.0 - 1.0) / 8.0

			d7 := depthRegion[j] / 8000.0
			if d7 < 0.0 {
				d7 = -d7 * 0.3
			}
			d7 = d7*3.0 - 2.0
			if d7 < 0.0 {
				d7 /= 2.0
				if d7 < -1.0 {
					d7 = -1.0
				}
				d7 /= 1.4
				d7 /= 2.0
			} else {
				if d7 > 1.0 {
					d7 = 1.0
				}
				d7 /= 8.0
			}
			j++

			d8 := f3
			d9 := f2
			d8 += d7 * 0.2
			d8 = d8 * g.settings.baseSize / 8.0
			d0 := g.settings.baseSize + d8*4.0

			for l1 := 0; l1 < verticalPoints; l1++ {

				d1 := (float64(l1) - d0) * g.settings.stretchY * 128.0 / 256.0 / d9

				if d1 < 0.0 {
					d1 *= 4.0
				}

				minVal := minLimitRegion[i] / g.settings.lowerLimitScale
				maxVal := maxLimitRegion[i] / g.settings.upperLimitScale
				mainVal := (mainNoiseRegion[i]/10.0 + 1.0) / 2.0

				d5 := clampedLerp(minVal, maxVal, mainVal) - d1

				if l1 > verticalPoints-4 {

					d6 := float64(l1-(verticalPoints-4)) / 3.0
					d5 = d5*(1.0-d6) + -10.0*d6
				}

				heightMap[i] = d5
				i++
			}
		}
	}
}

func (g *ChunkGeneratorOverworld) SetBlocksInChunk(chunkX, chunkZ int, chunk *world.Chunk) {
	var biomesForGeneration []int
	if g.genLayer != nil {
		biomesForGeneration = g.genLayer.GetInts(chunkX*4-2, chunkZ*4-2, 10, 10)
	} else {

		biomesForGeneration = make([]int, 100)
		for i := range biomesForGeneration {
			biomesForGeneration[i] = 1
		}
	}

	verticalPoints := g.settings.MaxHeight/8 + 1
	horizontalPoints := 5
	bufferSize := horizontalPoints * horizontalPoints * verticalPoints

	heightMap := make([]float64, bufferSize)
	depthRegion := make([]float64, horizontalPoints*horizontalPoints)
	mainNoiseRegion := make([]float64, bufferSize)
	minLimitRegion := make([]float64, bufferSize)
	maxLimitRegion := make([]float64, bufferSize)

	g.generateHeightmap(chunkX*4, 0, chunkZ*4, heightMap, depthRegion, mainNoiseRegion, minLimitRegion, maxLimitRegion, biomesForGeneration)

	verticalSegments := g.settings.MaxHeight / 8

	for i := 0; i < 4; i++ {
		j := i * 5
		k := (i + 1) * 5

		for l := 0; l < 4; l++ {
			i1 := (j + l) * verticalPoints
			j1 := (j + l + 1) * verticalPoints
			k1 := (k + l) * verticalPoints
			l1 := (k + l + 1) * verticalPoints

			for i2 := 0; i2 < verticalSegments; i2++ {
				d0 := 0.125
				d1 := heightMap[i1+i2]
				d2 := heightMap[j1+i2]
				d3 := heightMap[k1+i2]
				d4 := heightMap[l1+i2]

				d5 := (heightMap[i1+i2+1] - d1) * d0
				d6 := (heightMap[j1+i2+1] - d2) * d0
				d7 := (heightMap[k1+i2+1] - d3) * d0
				d8 := (heightMap[l1+i2+1] - d4) * d0

				for j2 := 0; j2 < 8; j2++ {
					d9 := 0.25
					d10 := d1
					d11 := d2
					d12 := (d3 - d1) * d9
					d13 := (d4 - d2) * d9

					for k2 := 0; k2 < 4; k2++ {
						d14 := 0.25
						d16 := (d11 - d10) * d14
						lvt_45_1_ := d10 - d16

						for l2 := 0; l2 < 4; l2++ {

							lvt_45_1_ += d16

							y := i2*8 + j2
							if y >= 128 {
								continue
							}

							if lvt_45_1_ > 0.0 {

								chunk.SetBlock(i*4+k2, y, l*4+l2, 1, 0)
							} else if y < g.settings.seaLevel {

								chunk.SetBlock(i*4+k2, y, l*4+l2, 9, 0)
							}
						}
						d10 += d12
						d11 += d13
					}

					d1 += d5
					d2 += d6
					d3 += d7
					d4 += d8
				}
			}
		}
	}
}

func (g *ChunkGeneratorOverworld) GenerateChunk(cx, cz int32) {
	x, z := int(cx), int(cz)

	rnd := rand.NewRandom(int64(x)*341873128712 + int64(z)*132897987541)

	c := world.NewChunk(cx, cz)

	g.SetBlocksInChunk(x, z, c)

	g.replaceBiomeBlocks(x, z, c, rnd)
	g.caves.GenerateChunk(cx, cz, c)
	g.ravines.GenerateChunk(cx, cz, c)

	if g.settings.ChunkManager != nil {
		g.settings.ChunkManager.SetChunk(cx, cz, c)
	}
}

func (g *ChunkGeneratorOverworld) PopulateChunk(cx, cz int32) {

	x := int(cx) * 16
	z := int(cz) * 16

	var biomeID int
	if g.genLayer != nil {
		ints := g.genLayer.GetInts(x+16, z+16, 1, 1)
		biomeID = ints[0]
	} else {
		biomeID = 1
	}
	b := biome.GetBiome(uint8(biomeID))

	rnd := rand.NewRandom(g.seed)
	k := rnd.NextLong()/2*2 + 1
	l := rnd.NextLong()/2*2 + 1
	seed := int64(cx)*k + int64(cz)*l ^ g.seed
	rnd.SetSeed(seed)

	adapter := &structureAdapter{cm: g.settings.ChunkManager}

	if g.settings.useMineShafts {
		g.mineshaft.GenerateStructure(adapter, int(cx), int(cz))
	}

	flag := false
	if g.settings.useVillages {
		flag = g.village.GenerateStructure(adapter, int(cx), int(cz))
	}

	if g.settings.useStrongholds {
		g.stronghold.GenerateStructure(adapter, int(cx), int(cz))
	}

	if g.settings.useTemples {
		g.scatteredFeature.GenerateStructure(adapter, int(cx), int(cz))
	}

	isOceanOrRiver := biomeID == 0 || biomeID == 10 || biomeID == 24 || biomeID == 7 || biomeID == 11

	if b.GetID() != 2 && b.GetID() != 17 && !isOceanOrRiver && g.settings.useWaterLakes && !flag && rnd.NextBoundedInt(g.settings.waterLakeChance) == 0 {
		i1 := rnd.NextBoundedInt(16) + 8
		j1 := rnd.NextBoundedInt(256)
		k1 := rnd.NextBoundedInt(16) + 8

		lake := object.NewLake(9)
		lake.Generate(g.settings.ChunkManager, rnd, world.NewBlockPos(int32(x+i1), int32(j1), int32(z+k1)))
	}

	if !flag && !isOceanOrRiver && rnd.NextBoundedInt(g.settings.lavaLakeChance/10) == 0 && g.settings.useLavaLakes {
		i2 := rnd.NextBoundedInt(16) + 8
		l2 := rnd.NextBoundedInt(rnd.NextBoundedInt(248) + 8)
		k3 := rnd.NextBoundedInt(16) + 8

		if l2 < g.settings.seaLevel || rnd.NextBoundedInt(g.settings.lavaLakeChance/8) == 0 {

			lake := object.NewLake(10)
			lake.Generate(g.settings.ChunkManager, rnd, world.NewBlockPos(int32(x+i2), int32(l2), int32(z+k3)))
		}
	}

	if g.settings.useDungeons {
		for j2 := 0; j2 < g.settings.dungeonChance; j2++ {
			i3 := rnd.NextBoundedInt(16) + 8
			l3 := rnd.NextBoundedInt(256)
			l1 := rnd.NextBoundedInt(16) + 8
			dungeon := object.NewDungeon()
			dungeon.Generate(g.settings.ChunkManager, rnd, world.NewBlockPos(int32(x+i3), int32(l3), int32(z+l1)))
		}
	}

	if d := b.GetDecorator(); d != nil {
		d.FlowerNoise = g.flowerNoise
	}
	b.Decorate(g.settings.ChunkManager, rnd, world.NewBlockPos(int32(x), 0, int32(z)))

	chunk := g.settings.ChunkManager.GetChunk(cx, cz, false)
	if chunk == nil {
		return
	}

	if chunk != nil {
		chunk.RecalculateHeightMap()
		chunk.InitBasicLighting()
		chunk.LightPopulated = true
		chunk.Populated = true
	}
}

func (g *ChunkGeneratorOverworld) replaceBiomeBlocks(x, z int, c *world.Chunk, rnd *rand.Random) {
	chunkX := x * 16
	chunkZ := z * 16

	var biomeIDs []int
	if g.biomeGen != nil {
		biomeIDs = g.biomeGen.GetInts(chunkX, chunkZ, 16, 16)
	} else {
		biomeIDs = make([]int, 256)
		for i := range biomeIDs {
			biomeIDs[i] = 1
		}
	}

	noiseArray := g.surfaceNoise.GetRegion(nil, float64(chunkX), float64(chunkZ), 16, 16, 0.0625, 0.0625, 1.0)

	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {

			idx := j*16 + i
			biomeID := uint8(biomeIDs[idx])

			b := biome.GetBiome(biomeID)

			nv := noiseArray[idx]

			gx := x*16 + i
			gz := z*16 + j

			b.GenTerrainBlocks(c, rnd, gx, gz, nv)

			biomeColor := b.GetColor()
			c.SetBiomeColor(i, j, uint32(biomeColor))

			c.SetBiomeID(i, j, biomeID)
		}
	}
}

type structureAdapter struct {
	cm generator.ChunkManager
}

func (a *structureAdapter) GetBlock(x, y, z int) (byte, byte) {

	return a.cm.GetBlockId(int32(x), int32(y), int32(z)), 0
}

func (a *structureAdapter) SetBlock(x, y, z int, id, meta byte) {
	a.cm.SetBlock(int32(x), int32(y), int32(z), id, meta, false)
}

type genLayerAdapter struct {
	layer layer.GenLayer
}

func (a *genLayerAdapter) GetBiome(x, z int) uint8 {
	if a.layer == nil {
		return 0
	}
	ints := a.layer.GetInts(x, z, 1, 1)
	if len(ints) > 0 {
		return uint8(ints[0])
	}
	return 0
}
