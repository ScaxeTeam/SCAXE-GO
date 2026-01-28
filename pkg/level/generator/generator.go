package generator

import (
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Generator interface {
	Init(level ChunkManager, seed int64)

	GenerateChunk(cx, cz int32)

	PopulateChunk(cx, cz int32)

	GetSpawn() *world.Vector3

	GetName() string

	GetSettings() map[string]interface{}
}

type ChunkManager interface {
	GetChunk(x, z int32, create bool) *world.Chunk

	SetChunk(x, z int32, chunk *world.Chunk)

	GetSeed() int64

	GetBlockId(x, y, z int32) byte

	SetBlock(x, y, z int32, id, meta byte, update bool) bool

	GetHeight(x, z int32) int32
}

var RegisteredGenerators = make(map[string]func(map[string]interface{}) Generator)

func RegisterGenerator(name string, factory func(map[string]interface{}) Generator) {
	RegisteredGenerators[name] = factory
}

func GetGenerator(name string, settings map[string]interface{}) Generator {
	if factory, ok := RegisteredGenerators[name]; ok {
		return factory(settings)
	}
	return nil
}
