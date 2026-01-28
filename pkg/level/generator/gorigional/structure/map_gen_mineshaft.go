package structure

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type MapGenMineshaft struct {
	worldSeed int64
	rand      *rand.Random
	chance    float64
}

func NewMapGenMineshaft(seed int64) *MapGenMineshaft {
	return &MapGenMineshaft{
		worldSeed: seed,
		rand:      rand.NewRandom(seed),
		chance:    0.004,
	}
}

func (m *MapGenMineshaft) CanSpawnStructureAtCoords(chunkX, chunkZ int) bool {

	m.rand.SetSeed(int64(chunkX)*341873128712 + int64(chunkZ)*132897987541 + m.worldSeed)
	return m.rand.NextDouble() < m.chance
}

func (m *MapGenMineshaft) GetStructureStart(chunkX, chunkZ int) *StructureStart {
	return NewMineshaftStart(m.worldSeed, chunkX, chunkZ)
}

func (m *MapGenMineshaft) GenerateStructure(w WorldAccess, chunkX, chunkZ int) bool {

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
