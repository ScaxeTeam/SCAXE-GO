package level

import (
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Provider interface {
	GetName() string

	LoadChunk(x, z int32) (*world.Chunk, error)

	SaveChunk(chunk *world.Chunk) error

	Close() error
}
