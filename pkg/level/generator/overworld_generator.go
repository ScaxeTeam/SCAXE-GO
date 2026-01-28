package generator

import (
	"math"
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/math/noise"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type OverworldGenerator struct {
	Seed int64

	heightOctave     *noise.PerlinOctaveGenerator
	roughnessOctave  *noise.PerlinOctaveGenerator
	roughness2Octave *noise.PerlinOctaveGenerator
	detailOctave     *noise.PerlinOctaveGenerator
	surfaceOctave    *noise.SimplexOctaveGenerator

	density [5][5][17]float64

	coordinateScale   float64
	heightScale       float64
	heightNoiseScaleX float64
	heightNoiseScaleZ float64
	detailNoiseScaleX float64
	detailNoiseScaleY float64
	detailNoiseScaleZ float64
	surfaceScale      float64
	baseSize          float64
	stretchY          float64
	biomeHeightOffset float64
	biomeHeightWeight float64
	biomeScaleOffset  float64
	biomeScaleWeight  float64
}

func NewOverworldGenerator(seed int64) *OverworldGenerator {
	r := rand.New(rand.NewSource(seed))

	g := &OverworldGenerator{
		Seed: seed,

		coordinateScale:   684.412,
		heightScale:       684.412,
		heightNoiseScaleX: 200.0,
		heightNoiseScaleZ: 200.0,
		detailNoiseScaleX: 80.0,
		detailNoiseScaleY: 160.0,
		detailNoiseScaleZ: 80.0,
		surfaceScale:      0.0625,
		baseSize:          8.5,
		stretchY:          12.0,
		biomeHeightOffset: 0.0,
		biomeHeightWeight: 1.0,
		biomeScaleOffset:  0.0,
		biomeScaleWeight:  1.0,
	}

	g.heightOctave = noise.NewPerlinOctaveGenerator2D(r, 16, 5, 5)
	g.heightOctave.SetXScale(g.heightNoiseScaleX)
	g.heightOctave.SetZScale(g.heightNoiseScaleZ)

	g.roughnessOctave = noise.NewPerlinOctaveGenerator(r, 16, 5, 17, 5)
	g.roughnessOctave.SetXScale(g.coordinateScale)
	g.roughnessOctave.SetYScale(g.heightScale)
	g.roughnessOctave.SetZScale(g.coordinateScale)

	g.roughness2Octave = noise.NewPerlinOctaveGenerator(r, 16, 5, 17, 5)
	g.roughness2Octave.SetXScale(g.coordinateScale)
	g.roughness2Octave.SetYScale(g.heightScale)
	g.roughness2Octave.SetZScale(g.coordinateScale)

	g.detailOctave = noise.NewPerlinOctaveGenerator(r, 8, 5, 17, 5)
	g.detailOctave.SetXScale(g.coordinateScale / g.detailNoiseScaleX)
	g.detailOctave.SetYScale(g.heightScale / g.detailNoiseScaleY)
	g.detailOctave.SetZScale(g.coordinateScale / g.detailNoiseScaleZ)

	g.surfaceOctave = noise.NewSimplexOctaveGenerator(seed, 4, 16, 16)
	g.surfaceOctave.SetScale(g.surfaceScale)

	return g
}

func (g *OverworldGenerator) GetName() string {
	return "overworld"
}

func (g *OverworldGenerator) GenerateChunk(chunk *world.Chunk) {
	g.generateRawTerrain(chunk, nil)
}

func (g *OverworldGenerator) GenerateChunkWithBiomes(chunk *world.Chunk, biomeGrid []int) {
	g.generateRawTerrain(chunk, biomeGrid)
}

func (g *OverworldGenerator) generateRawTerrain(chunk *world.Chunk, biomeGrid []int) {
	chunkX := chunk.X
	chunkZ := chunk.Z

	g.generateTerrainDensityWithBiomes(chunkX, chunkZ, biomeGrid)

	seaLevel := 64

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 16; k++ {

				d1 := g.density[i][j][k]
				d2 := g.density[i+1][j][k]
				d3 := g.density[i][j+1][k]
				d4 := g.density[i+1][j+1][k]

				d5 := (g.density[i][j][k+1] - d1) / 8.0
				d6 := (g.density[i+1][j][k+1] - d2) / 8.0
				d7 := (g.density[i][j+1][k+1] - d3) / 8.0
				d8 := (g.density[i+1][j+1][k+1] - d4) / 8.0

				for l := 0; l < 8; l++ {
					d9 := d1
					d10 := d3

					for m := 0; m < 4; m++ {
						dens := d9

						for n := 0; n < 4; n++ {
							bx := m + (i << 2)
							by := l + (k << 3)
							bz := n + (j << 2)

							if dens > 0 {
								chunk.SetBlock(bx, by, bz, 1, 0)
							} else if by < seaLevel {
								chunk.SetBlock(bx, by, bz, 9, 0)
							}

							dens += (d10 - d9) / 4.0
						}

						d9 += (d2 - d1) / 4.0
						d10 += (d4 - d3) / 4.0
					}

					d1 += d5
					d3 += d7
					d2 += d6
					d4 += d8
				}
			}
		}
	}
}

func (g *OverworldGenerator) generateTerrainDensityWithBiomes(chunkX, chunkZ int32, biomeGrid []int) {

	x := int(chunkX) << 2
	z := int(chunkZ) << 2

	heightNoise := g.heightOctave.GetFractalBrownianMotion(float64(x), float64(z), 0.5, 2.0)
	roughnessNoise := g.roughnessOctave.GetFractalBrownianMotion3D(float64(x), 0, float64(z), 0.5, 2.0)
	roughness2Noise := g.roughness2Octave.GetFractalBrownianMotion3D(float64(x), 0, float64(z), 0.5, 2.0)
	detailNoise := g.detailOctave.GetFractalBrownianMotion3D(float64(x), 0, float64(z), 0.5, 2.0)

	index := 0
	indexHeight := 0

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {

			avgHeightScale := 0.0
			avgHeightBase := 0.0
			totalWeight := 0.0

			var biomeHeight BiomeHeight
			if biomeGrid != nil && len(biomeGrid) >= 100 {

				biomeID := biomeGrid[(i+2)+(j+2)*10]
				biomeHeight = GetBiomeHeight(biomeID)
			} else {
				biomeHeight = BiomeHeightDefault
			}

			for m := 0; m < 5; m++ {
				for n := 0; n < 5; n++ {
					var nearBiomeHeight BiomeHeight
					if biomeGrid != nil && len(biomeGrid) >= 100 {
						nearBiomeID := biomeGrid[(i+m)+(j+n)*10]
						nearBiomeHeight = GetBiomeHeight(nearBiomeID)
					} else {
						nearBiomeHeight = biomeHeight
					}

					heightBase := g.biomeHeightOffset + nearBiomeHeight.Height*g.biomeHeightWeight
					heightScale := g.biomeScaleOffset + nearBiomeHeight.Scale*g.biomeScaleWeight

					weight := ElevationWeight[m][n] / (heightBase + 2.0)
					if nearBiomeHeight.Height > biomeHeight.Height {
						weight *= 0.5
					}

					avgHeightScale += heightScale * weight
					avgHeightBase += heightBase * weight
					totalWeight += weight
				}
			}

			avgHeightScale /= totalWeight
			avgHeightBase /= totalWeight
			avgHeightScale = avgHeightScale*0.9 + 0.1
			avgHeightBase = (avgHeightBase*4.0 - 1.0) / 8.0

			noiseH := heightNoise[indexHeight] / 8000.0
			indexHeight++

			if noiseH < 0 {
				noiseH = math.Abs(noiseH) * 0.3
			}
			noiseH = noiseH*3.0 - 2.0
			if noiseH < 0 {
				noiseH = math.Max(noiseH*0.5, -1) / 1.4 * 0.5
			} else {
				noiseH = math.Min(noiseH, 1) / 8.0
			}

			noiseH = (noiseH*0.2+avgHeightBase)*g.baseSize/8.0*4.0 + g.baseSize

			for k := 0; k < 17; k++ {

				nh := (float64(k) - noiseH) * g.stretchY * 128.0 / 256.0 / avgHeightScale
				if nh < 0 {
					nh *= 4.0
				}

				noiseR := roughnessNoise[index] / 512.0
				noiseR2 := roughness2Noise[index] / 512.0
				noiseD := (detailNoise[index]/10.0 + 1.0) / 2.0

				var dens float64
				if noiseD < 0 {
					dens = noiseR
				} else if noiseD > 1 {
					dens = noiseR2
				} else {
					dens = noiseR + (noiseR2-noiseR)*noiseD
				}
				dens -= nh
				index++

				if k > 13 {
					lowering := float64(k-13) / 3.0
					dens = dens*(1.0-lowering) + -10.0*lowering
				}

				g.density[i][j][k] = dens
			}
		}
	}
}
