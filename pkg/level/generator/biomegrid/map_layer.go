package biomegrid

import (
	"math/rand"
)

type MapLayer interface {
	GenerateValues(x, z, sizeX, sizeZ int) []int
}

type BaseMapLayer struct {
	Seed   int64
	Parent MapLayer
	random *rand.Rand
}

func NewBaseMapLayer(seed int64, parent MapLayer) *BaseMapLayer {
	return &BaseMapLayer{
		Seed:   seed,
		Parent: parent,
		random: rand.New(rand.NewSource(seed)),
	}
}

func (l *BaseMapLayer) SetCoordsSeed(x, z int) {
	l.random.Seed(l.Seed)
	seedX := l.random.Int63()
	seedZ := l.random.Int63()
	l.random.Seed(int64(x)*seedX + int64(z)*seedZ ^ l.Seed)
}

func (l *BaseMapLayer) NextInt(max int) int {
	if max <= 0 {
		return 0
	}
	return l.random.Intn(max)
}

func (l *BaseMapLayer) GetParent() MapLayer {
	return l.Parent
}
