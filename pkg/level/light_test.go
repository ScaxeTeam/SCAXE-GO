package level

import (
	"testing"

	"github.com/scaxe/scaxe-go/pkg/world"
)

type MockProvider struct{}

func (m *MockProvider) GetName() string                            { return "mock" }
func (m *MockProvider) Close() error                               { return nil }
func (m *MockProvider) LoadChunk(x, z int32) (*world.Chunk, error) { return nil, nil }
func (m *MockProvider) SaveChunk(chunk *world.Chunk) error         { return nil }
func (m *MockProvider) GetSpawn() *world.Vector3                   { return world.NewVector3(0, 0, 0) }
func (m *MockProvider) SetSpawn(v *world.Vector3)                  {}
func (m *MockProvider) IsChunkLoaded(x, z int32) bool              { return false }

func TestBlockLightPropagation(t *testing.T) {

	lvl := NewLevel("test", "test", &MockProvider{}, "flat")

	lvl.GetChunk(0, 0, true)

	lvl.SetBlock(8, 50, 8, 50, 0, true)

	light := lvl.GetBlockLightAt(8, 50, 8)
	if light != 14 {
		t.Errorf("Expected torch light 14, got %d", light)
	}

	adjLight := lvl.GetBlockLightAt(9, 50, 8)
	if adjLight != 13 {
		t.Errorf("Expected adjacent light 13, got %d", adjLight)
	}

	diagLight := lvl.GetBlockLightAt(9, 50, 9)
	if diagLight != 12 {
		t.Errorf("Expected diagonal light 12, got %d", diagLight)
	}

	lvl.SetBlock(8, 50, 7, 1, 0, true)

	wrappedLight := lvl.GetBlockLightAt(8, 50, 6)
	if wrappedLight != 10 {
		t.Errorf("Expected wrapped light 10, got %d", wrappedLight)
	}

	stoneLight := lvl.GetBlockLightAt(8, 50, 7)
	if stoneLight != 0 {
		t.Errorf("Expected stone light 0, got %d", stoneLight)
	}

	lvl.SetBlock(8, 50, 8, 0, 0, true)

	if l := lvl.GetBlockLightAt(8, 50, 8); l != 0 {
		t.Errorf("Expected light 0 after removal, got %d", l)
	}

	if l := lvl.GetBlockLightAt(9, 50, 8); l != 0 {
		t.Errorf("Expected neighbor light 0 after removal, got %d", l)
	}
}
