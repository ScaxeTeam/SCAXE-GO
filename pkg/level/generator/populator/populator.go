package populator

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type ChunkManager interface {
	GetChunk(x, z int32, create bool) *world.Chunk
	SetChunk(x, z int32, chunk *world.Chunk)
	GetSeed() int64
	GetBlockId(x, y, z int32) byte
	SetBlock(x, y, z int32, id, meta byte, update bool) bool
	GetHeight(x, z int32) int32
}

type Populator interface {
	Populate(level ChunkManager, chunk *world.Chunk, x, z int32, random *rand.Random)
}
