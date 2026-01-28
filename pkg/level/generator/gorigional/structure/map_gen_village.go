package structure

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type MapGenVillage struct {
	worldSeed    int64
	distance     int
	rand         *rand.Random
	structureMap map[int64]*StructureStart
}

func NewMapGenVillage(seed int64) *MapGenVillage {
	return &MapGenVillage{
		worldSeed:    seed,
		distance:     32,
		rand:         rand.NewRandom(seed),
		structureMap: make(map[int64]*StructureStart),
	}
}

func (m *MapGenVillage) CanSpawnStructureAtCoords(chunkX, chunkZ int) bool {
	i := chunkX
	j := chunkZ

	if chunkX < 0 {
		chunkX -= m.distance - 1
	}
	if chunkZ < 0 {
		chunkZ -= m.distance - 1
	}

	k := chunkX / m.distance
	l := chunkZ / m.distance

	seed := int64(k)*341873128712 + int64(l)*132897987541 + m.worldSeed + 10387312
	m.rand.SetSeed(seed)

	k *= m.distance
	l *= m.distance

	k += m.rand.NextBoundedInt(m.distance - 8)
	l += m.rand.NextBoundedInt(m.distance - 8)

	return i == k && j == l
}

func (m *MapGenVillage) GetStructureStart(chunkX, chunkZ int) *StructureStart {
	key := int64(chunkX) + (int64(chunkZ) << 32)

	if start, ok := m.structureMap[key]; ok {
		return start
	}

	if m.CanSpawnStructureAtCoords(chunkX, chunkZ) {
		start := NewVillageStart(m.worldSeed, m.rand, chunkX, chunkZ, 0)
		m.structureMap[key] = start.StructureStart
		return start.StructureStart
	}
	return nil
}

func (m *MapGenVillage) GenerateStructure(w WorldAccess, chunkX, chunkZ int) bool {

	x := chunkX * 16
	z := chunkZ * 16
	chunkBox := NewBoundingBox(x, 0, z, x+15, 255, z+15)

	success := false

	for i := chunkX - 8; i <= chunkX+8; i++ {
		for j := chunkZ - 8; j <= chunkZ+8; j++ {
			if m.CanSpawnStructureAtCoords(i, j) {
				start := m.GetStructureStart(i, j)
				if start != nil && start.BoundingBox != nil && start.BoundingBox.IntersectsWith(chunkBox) {

					start.GenerateStructure(w, m.rand, chunkBox)
					success = true
				}
			}
		}
	}
	return success
}
