package generator

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Flat struct {
	level ChunkManager
	seed  int64

	layers map[int]byte

	options map[string]interface{}
}

var _ Generator = (*Flat)(nil)

func NewFlat(settings map[string]interface{}) Generator {
	f := &Flat{
		layers:  make(map[int]byte),
		options: settings,
	}

	f.layers[0] = byte(block.BEDROCK)
	f.layers[1] = byte(block.DIRT)
	f.layers[2] = byte(block.DIRT)
	f.layers[3] = byte(block.GRASS)

	return f
}

func (f *Flat) Init(level ChunkManager, seed int64) {
	f.level = level
	f.seed = seed
}

func (f *Flat) GetName() string {
	return "flat"
}

func (f *Flat) GetSettings() map[string]interface{} {
	return f.options
}

func (f *Flat) GetSpawn() *world.Vector3 {
	return world.NewVector3(128, 4, 128)
}

func (f *Flat) GenerateChunk(cx, cz int32) {
	chunk := world.NewChunk(cx, cz)

	for y, blockID := range f.layers {
		for x := 0; x < 16; x++ {
			for z := 0; z < 16; z++ {
				chunk.SetBlock(x, y, z, blockID, 0)
			}
		}
	}

	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			chunk.SetBiomeColor(x, z, 0x8DB360)
		}
	}

	chunk.Generated = true
	f.level.SetChunk(cx, cz, chunk)
}

func (f *Flat) PopulateChunk(cx, cz int32) {
	chunk := f.level.GetChunk(cx, cz, false)
	if chunk == nil {
		return
	}

	chunk.Populated = true
}

func init() {
	RegisterGenerator("flat", NewFlat)
}
