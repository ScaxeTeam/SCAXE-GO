package level

import "github.com/scaxe/scaxe-go/pkg/world"

type ChunkLoader interface {
	GetLoaderId() int64

	OnChunkLoaded(chunk *world.Chunk)

	OnChunkUnloaded(chunk *world.Chunk)
}
