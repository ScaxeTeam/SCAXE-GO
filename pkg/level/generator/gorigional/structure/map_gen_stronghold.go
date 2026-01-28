package structure

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type MapGenStronghold struct {
	structureCoords []world.ChunkPos
	worldSeed       int64
	ran             *rand.Random
	distance        float64
	spread          int
	structureMap    map[int64]*StructureStart
}

func NewMapGenStronghold(seed int64) *MapGenStronghold {
	m := &MapGenStronghold{
		worldSeed:       seed,
		ran:             rand.NewRandom(seed),
		distance:        32.0,
		spread:          3,
		structureCoords: make([]world.ChunkPos, 0, 128),
		structureMap:    make(map[int64]*StructureStart),
	}
	m.generatePositions()
	return m
}

func (m *MapGenStronghold) generatePositions() {
	m.ran.SetSeed(m.worldSeed)

	generated := 0
	ring := 0

	for generated < 128 {
		d1 := (4.0*m.distance + m.distance*float64(ring)*6.0) + (m.ran.NextDouble()-0.5)*m.distance*2.5

		countInRing := 0
		if ring == 0 {
			countInRing = m.spread
		} else {
			countInRing = ring*m.spread + m.spread
		}

		if generated+countInRing > 128 {
			countInRing = 128 - generated
		}

		d2 := m.ran.NextDouble() * math.Pi * 2.0
		d3 := math.Pi * 2.0 / float64(countInRing)

		for l := 0; l < countInRing; l++ {
			d4 := (4.0*m.distance + m.distance*float64(ring)*6.0) + (m.ran.NextDouble()-0.5)*m.distance*2.5
			_ = d4

			angle := d2 + float64(l)*d3

			cx := int(math.Round(math.Cos(angle) * d1))
			cz := int(math.Round(math.Sin(angle) * d1))

			m.structureCoords = append(m.structureCoords, world.ChunkPos{X: int32(cx), Z: int32(cz)})
			generated++
		}
		ring++
	}
}

func (m *MapGenStronghold) CanSpawnStructureAtCoords(chunkX, chunkZ int) bool {
	for _, pos := range m.structureCoords {
		if int(pos.X) == chunkX && int(pos.Z) == chunkZ {
			return true
		}
	}
	return false
}

func (m *MapGenStronghold) GetStructureStart(chunkX, chunkZ int) *StructureStart {
	key := int64(chunkX) + (int64(chunkZ) << 32)

	if start, ok := m.structureMap[key]; ok {
		return start
	}

	if m.CanSpawnStructureAtCoords(chunkX, chunkZ) {
		start := NewStrongholdStart(m.worldSeed, m.ran, chunkX, chunkZ).StructureStart
		m.structureMap[key] = start
		return start
	}
	return nil
}

func (m *MapGenStronghold) GenerateStructure(w WorldAccess, chunkX, chunkZ int) bool {
	x := chunkX * 16
	z := chunkZ * 16
	chunkBox := NewBoundingBox(x, 0, z, x+15, 255, z+15)

	success := false

	for i := chunkX - 8; i <= chunkX+8; i++ {
		for j := chunkZ - 8; j <= chunkZ+8; j++ {
			if m.CanSpawnStructureAtCoords(i, j) {
				start := m.GetStructureStart(i, j)
				if start != nil && start.BoundingBox != nil && start.BoundingBox.IntersectsWith(chunkBox) {
					start.GenerateStructure(w, m.ran, chunkBox)
					success = true
				}
			}
		}
	}
	return success
}

func (m *MapGenStronghold) GetStructureCoords() []world.ChunkPos {
	return m.structureCoords
}
