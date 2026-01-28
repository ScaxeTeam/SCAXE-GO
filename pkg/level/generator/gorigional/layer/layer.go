package layer

type GenLayer interface {
	GetInts(x, z, width, depth int) []int

	InitWorldGenSeed(seed int64)
}

type BaseLayer struct {
	Parent       GenLayer
	BaseSeed     int64
	WorldGenSeed int64
	ChunkSeed    int64
}

func NewBaseLayer(baseSeed int64) *BaseLayer {

	l := &BaseLayer{BaseSeed: baseSeed}
	l.BaseSeed = l.mixSeed(l.BaseSeed, baseSeed)
	l.BaseSeed = l.mixSeed(l.BaseSeed, baseSeed)
	l.BaseSeed = l.mixSeed(l.BaseSeed, baseSeed)
	return l
}

func (l *BaseLayer) mixSeed(current int64, add int64) int64 {
	current *= current*6364136223846793005 + 1442695040888963407
	current += add
	return current
}

func (l *BaseLayer) InitWorldGenSeed(seed int64) {
	l.WorldGenSeed = seed
	if l.Parent != nil {
		l.Parent.InitWorldGenSeed(seed)
	}

	l.WorldGenSeed = l.mixSeed(l.WorldGenSeed, l.BaseSeed)
	l.WorldGenSeed = l.mixSeed(l.WorldGenSeed, l.BaseSeed)
	l.WorldGenSeed = l.mixSeed(l.WorldGenSeed, l.BaseSeed)
}

func (l *BaseLayer) InitChunkSeed(x, z int64) {
	l.ChunkSeed = l.WorldGenSeed
	l.ChunkSeed = l.mixSeed(l.ChunkSeed, x)
	l.ChunkSeed = l.mixSeed(l.ChunkSeed, z)
	l.ChunkSeed = l.mixSeed(l.ChunkSeed, x)
	l.ChunkSeed = l.mixSeed(l.ChunkSeed, z)
}

func (l *BaseLayer) NextInt(max int) int {

	val := int((l.ChunkSeed >> 24) % int64(max))
	if val < 0 {
		val += max
	}

	l.ChunkSeed *= l.ChunkSeed*6364136223846793005 + 1442695040888963407
	l.ChunkSeed += l.WorldGenSeed

	return val
}

func (l *BaseLayer) SelectRandom(choices ...int) int {
	return choices[l.NextInt(len(choices))]
}
