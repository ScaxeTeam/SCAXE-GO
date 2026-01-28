package structure

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type BiomeSource interface {
	GetBiome(x, z int) uint8
}

type MapGenScatteredFeature struct {
	worldSeed                           int64
	maxDistanceBetweenScatteredFeatures int
	minDistanceBetweenScatteredFeatures int
	rand                                *rand.Random
	biomeSource                         BiomeSource
	structureMap                        map[int64]*StructureStart
}

func NewMapGenScatteredFeature(seed int64, biomeSource BiomeSource) *MapGenScatteredFeature {
	m := &MapGenScatteredFeature{
		worldSeed:                           seed,
		maxDistanceBetweenScatteredFeatures: 32,
		minDistanceBetweenScatteredFeatures: 8,
		rand:                                rand.NewRandom(seed),
		biomeSource:                         biomeSource,
		structureMap:                        make(map[int64]*StructureStart),
	}
	return m
}

func (m *MapGenScatteredFeature) CanSpawnStructureAtCoords(chunkX, chunkZ int) bool {
	i := chunkX
	j := chunkZ

	if chunkX < 0 {
		chunkX -= m.maxDistanceBetweenScatteredFeatures - 1
	}
	if chunkZ < 0 {
		chunkZ -= m.maxDistanceBetweenScatteredFeatures - 1
	}

	k := chunkX / m.maxDistanceBetweenScatteredFeatures
	l := chunkZ / m.maxDistanceBetweenScatteredFeatures

	seed := int64(k)*341873128712 + int64(l)*132897987541 + m.worldSeed + 14357617
	m.rand.SetSeed(seed)

	k *= m.maxDistanceBetweenScatteredFeatures
	l *= m.maxDistanceBetweenScatteredFeatures

	k += m.rand.NextBoundedInt(m.maxDistanceBetweenScatteredFeatures - 8)
	l += m.rand.NextBoundedInt(m.maxDistanceBetweenScatteredFeatures - 8)

	if i == k && j == l {

		if m.biomeSource != nil {

			bID := m.biomeSource.GetBiome(i*16+8, j*16+8)
			return m.isStructureBiome(bID)
		}
		return true
	}

	return false
}

func (m *MapGenScatteredFeature) isStructureBiome(id uint8) bool {

	switch id {
	case 2, 17, 21, 22, 6, 12, 30:
		return true
	}
	return false
}

func (m *MapGenScatteredFeature) GetStructureStart(chunkX, chunkZ int) *StructureStart {
	key := int64(chunkX) + (int64(chunkZ) << 32)

	if start, ok := m.structureMap[key]; ok {
		return start
	}

	if m.CanSpawnStructureAtCoords(chunkX, chunkZ) {
		var biomeID uint8 = 1
		if m.biomeSource != nil {
			biomeID = m.biomeSource.GetBiome(chunkX*16+8, chunkZ*16+8)
		}

		start := NewScatteredFeatureStart(m.worldSeed, m.rand, chunkX, chunkZ, biomeID)
		m.structureMap[key] = start.StructureStart
		return start.StructureStart
	}
	return nil
}

func (m *MapGenScatteredFeature) GenerateStructure(w WorldAccess, chunkX, chunkZ int) bool {
	x := chunkX * 16
	z := chunkZ * 16
	chunkBox := NewBoundingBox(x, 0, z, x+15, 255, z+15)

	success := false
	rangeVal := 8

	for i := chunkX - rangeVal; i <= chunkX+rangeVal; i++ {
		for j := chunkZ - rangeVal; j <= chunkZ+rangeVal; j++ {
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

type ScatteredFeatureStart struct {
	*StructureStart
}

func NewScatteredFeatureStart(seed int64, rnd *rand.Random, chunkX, chunkZ int, biomeID uint8) *ScatteredFeatureStart {
	s := &ScatteredFeatureStart{
		StructureStart: NewStructureStart(chunkX, chunkZ),
	}

	s.createComponents(rnd, chunkX, chunkZ, biomeID)
	s.UpdateBoundingBox()
	return s
}

func (s *ScatteredFeatureStart) createComponents(rnd *rand.Random, chunkX, chunkZ int, biomeID uint8) {

	x := chunkX * 16
	z := chunkZ * 16

	switch biomeID {
	case 21, 22:
		s.Components = append(s.Components, NewJunglePyramid(rnd, x, z))
	case 6:
		s.Components = append(s.Components, NewSwampHut(rnd, x, z))
	case 2, 17:
		s.Components = append(s.Components, NewDesertPyramid(rnd, x, z))
	case 12, 30:

	default:

		s.Components = append(s.Components, NewDesertPyramid(rnd, x, z))
	}
}
