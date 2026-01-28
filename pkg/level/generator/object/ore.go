package object

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/java"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type OreType struct {
	Material     int
	Meta         int
	ClusterCount int
	ClusterSize  int
	MinHeight    int
	MaxHeight    int
}

type Ore struct {
	BlockID    uint8
	BlockMeta  uint8
	BlockCount int
}

func NewOre(id uint8, meta uint8, count int) *Ore {
	return &Ore{
		BlockID:    id,
		BlockMeta:  meta,
		BlockCount: count,
	}
}

func (o *Ore) Generate(w populator.ChunkManager, r *rand.Random, pos world.BlockPos) bool {
	f := r.NextFloat() * math.Pi

	sinF := java.Sin(float32(f))
	cosF := java.Cos(float32(f))

	fCount := float32(o.BlockCount) / 8.0

	d0 := float64(float32(pos.X()+8) + sinF*fCount)
	d1 := float64(float32(pos.X()+8) - sinF*fCount)
	d2 := float64(float32(pos.Z()+8) + cosF*fCount)
	d3 := float64(float32(pos.Z()+8) - cosF*fCount)

	d4 := float64(pos.Y() + int32(r.NextBoundedInt(3)) - 2)
	d5 := float64(pos.Y() + int32(r.NextBoundedInt(3)) - 2)

	for i := 0; i < o.BlockCount; i++ {
		f1 := float64(i) / float64(o.BlockCount)
		d6 := d0 + (d1-d0)*f1
		d7 := d4 + (d5-d4)*f1
		d8 := d2 + (d3-d2)*f1

		d9 := r.NextDouble() * float64(o.BlockCount) / 16.0

		sinPiF1 := float64(java.Sin(float32(math.Pi * f1)))
		d10 := (sinPiF1+1.0)*d9 + 1.0
		d11 := (sinPiF1+1.0)*d9 + 1.0

		j := java.Floor(d6 - d10/2.0)
		k := java.Floor(d7 - d11/2.0)
		l := java.Floor(d8 - d10/2.0)

		i1 := java.Floor(d6 + d10/2.0)
		j1 := java.Floor(d7 + d11/2.0)
		k1 := java.Floor(d8 + d10/2.0)

		for l1 := j; l1 <= i1; l1++ {
			d12 := (float64(l1) + 0.5 - d6) / (d10 / 2.0)

			if d12*d12 < 1.0 {
				for i2 := k; i2 <= j1; i2++ {
					d13 := (float64(i2) + 0.5 - d7) / (d11 / 2.0)

					if d12*d12+d13*d13 < 1.0 {
						for j2 := l; j2 <= k1; j2++ {
							d14 := (float64(j2) + 0.5 - d8) / (d10 / 2.0)

							if d12*d12+d13*d13+d14*d14 < 1.0 {

								if i2 >= 0 && i2 <= 255 {
									currID := w.GetBlockId(int32(l1), int32(i2), int32(j2))
									if currID == block.STONE {
										w.SetBlock(int32(l1), int32(i2), int32(j2), o.BlockID, o.BlockMeta, false)

									}
								}
							}
						}
					}
				}
			}
		}
	}
	return true
}
