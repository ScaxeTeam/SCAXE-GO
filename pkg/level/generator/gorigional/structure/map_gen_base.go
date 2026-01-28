package structure

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type BlockSetter interface {
	GetBlock(x, y, z int) (byte, byte)
	SetBlock(x, y, z int, id byte, meta byte) bool
}

type MapGenBase struct {
	Range     int
	Rand      *rand.Random
	WorldSeed int64
}

func NewMapGenBase(seed int64) *MapGenBase {
	return &MapGenBase{
		Range:     8,
		Rand:      rand.NewRandom(seed),
		WorldSeed: seed,
	}
}

func (m *MapGenBase) Generate(chunkX, chunkZ int, chunk BlockSetter, impl RecursiveGenerator) {
	m.Rand.SetSeed(m.WorldSeed)
	r1 := m.Rand.NextLong()
	r2 := m.Rand.NextLong()

	for x := chunkX - m.Range; x <= chunkX+m.Range; x++ {
		for z := chunkZ - m.Range; z <= chunkZ+m.Range; z++ {
			rX := int64(x) * r1
			rZ := int64(z) * r2
			seed := rX ^ rZ ^ m.WorldSeed
			m.Rand.SetSeed(seed)

			impl.RecursiveGenerate(chunkX, chunkZ, x, z, chunk)
		}
	}
}

type RecursiveGenerator interface {
	RecursiveGenerate(chunkX, chunkZ, originX, originZ int, chunk BlockSetter)
}
