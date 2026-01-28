package structure

import (
	"fmt"
	"testing"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type MockWorld struct {
	Blocks map[string]struct{ ID, Meta byte }
}

func NewMockWorld() *MockWorld {
	return &MockWorld{
		Blocks: make(map[string]struct{ ID, Meta byte }),
	}
}

func (m *MockWorld) SetBlock(x, y, z int, id, meta byte) {
	key := fmt.Sprintf("%d,%d,%d", x, y, z)
	m.Blocks[key] = struct{ ID, Meta byte }{id, meta}
}

func (m *MockWorld) GetBlock(x, y, z int) (byte, byte) {
	key := fmt.Sprintf("%d,%d,%d", x, y, z)
	if b, ok := m.Blocks[key]; ok {
		return b.ID, b.Meta
	}
	return 0, 0
}

func (m *MockWorld) AssertBlock(t *testing.T, x, y, z int, expectedID, expectedMeta byte, context string) {
	id, meta := m.GetBlock(x, y, z)
	if id != expectedID || meta != expectedMeta {
		t.Errorf("[%s] Expected block at %d,%d,%d to be %d:%d, but got %d:%d", context, x, y, z, expectedID, expectedMeta, id, meta)
	}
}

func TestDesertPyramid_Rotation(t *testing.T) {

	testRotation(t, 0, 3, "SOUTH")

	testRotation(t, 2, 2, "NORTH")

	testRotation(t, 1, 0, "WEST")

	testRotation(t, 3, 1, "EAST")
}

func testRotation(t *testing.T, mode int, expectedMeta byte, modeName string) {
	w := NewMockWorld()
	rnd := rand.NewRandom()

	dp := NewDesertPyramid(rnd, 0, 0)
	dp.CoordBaseMode = mode
	dp.BoundingBox = NewBoundingBox(0, 0, 0, 20, 20, 20)

	SANDSTONE_STAIRS := byte(128)
	NORTH := byte(3)

	dp.SetBlockState(w, SANDSTONE_STAIRS, NORTH, 2, 10, 0, dp.BoundingBox)

	x := dp.GetXWithOffset(2, 0)
	y := dp.GetYWithOffset(10)
	z := dp.GetZWithOffset(2, 0)

	w.AssertBlock(t, x, y, z, SANDSTONE_STAIRS, expectedMeta, modeName)
}

func TestDesertPyramid_FloorHeight(t *testing.T) {
	w := NewMockWorld()
	rnd := rand.NewRandom()
	dp := NewDesertPyramid(rnd, 0, 0)
	dp.CoordBaseMode = 0
	dp.BoundingBox = NewBoundingBox(0, 0, 0, 30, 30, 30)

}

func TestSwampHut_Stairs(t *testing.T) {

	sh := NewSwampHut(rand.NewRandom(), 0, 0)

}
