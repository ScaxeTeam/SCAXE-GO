package anvil

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type AnvilProvider struct {
	path    string
	loaders map[uint64]*RegionLoader
	mu      sync.Mutex
}

var _ level.Provider = (*AnvilProvider)(nil)

func NewAnvilProvider(path string) (*AnvilProvider, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, err
		}
	}

	regionDir := filepath.Join(path, "region")
	if _, err := os.Stat(regionDir); os.IsNotExist(err) {
		if err := os.MkdirAll(regionDir, 0755); err != nil {
			return nil, err
		}
	}

	return &AnvilProvider{
		path:    path,
		loaders: make(map[uint64]*RegionLoader),
	}, nil
}

func (p *AnvilProvider) GetName() string {
	return "anvil"
}

func regionHash(x, z int32) uint64 {
	return (uint64(x) << 32) | (uint64(z) & 0xFFFFFFFF)
}

func (p *AnvilProvider) getRegionLoader(rx, rz int32) (*RegionLoader, error) {
	hash := regionHash(rx, rz)

	p.mu.Lock()
	defer p.mu.Unlock()

	if loader, ok := p.loaders[hash]; ok {
		return loader, nil
	}

	loader, err := NewRegionLoader(p.path, rx, rz)
	if err != nil {
		return nil, err
	}

	p.loaders[hash] = loader
	return loader, nil
}

func (p *AnvilProvider) LoadChunk(x, z int32) (*world.Chunk, error) {
	rx := x >> 5
	rz := z >> 5

	loader, err := p.getRegionLoader(rx, rz)
	if err != nil {
		return nil, err
	}

	return loader.ReadChunk(x, z)
}

func (p *AnvilProvider) SaveChunk(chunk *world.Chunk) error {
	rx := chunk.X >> 5
	rz := chunk.Z >> 5

	loader, err := p.getRegionLoader(rx, rz)
	if err != nil {
		return err
	}

	return loader.WriteChunk(chunk)
}

func (p *AnvilProvider) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var lastErr error
	for _, loader := range p.loaders {
		if err := loader.Close(); err != nil {
			logger.Error("Failed to close region loader", "error", err)
			lastErr = err
		}
	}

	p.loaders = make(map[uint64]*RegionLoader)
	return lastErr
}
