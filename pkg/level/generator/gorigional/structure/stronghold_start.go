package structure

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type StrongholdStart struct {
	*StructureStart
}

func NewStrongholdStart(worldSeed int64, rnd *rand.Random, chunkX, chunkZ int) *StrongholdStart {
	start := &StrongholdStart{
		StructureStart: NewStructureStart(chunkX, chunkZ),
	}

	rnd.SetSeed(int64(chunkX)*341873128712 + int64(chunkZ)*132897987541 + worldSeed)

	stairs := NewStrongholdStairs2(start.StructureStart, rnd, (chunkX<<4)+2, (chunkZ<<4)+2)
	start.Components = append(start.Components, stairs)

	stairs.BuildComponent(stairs, &start.Components, rnd)

	start.UpdateBoundingBox()
	return start
}
