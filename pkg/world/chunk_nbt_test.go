package world

import (
	"testing"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

func TestChunk_ToNBT_FromNBT(t *testing.T) {
	c := NewChunk(10, 20)
	c.SetBlock(0, 0, 0, 1, 0)
	c.SetBlock(15, 127, 15, 7, 0)

	entityTag := nbt.NewCompoundTag("")
	entityTag.Set(nbt.NewStringTag("id", "Pig"))
	c.Entities = append(c.Entities, entityTag)

	nbtTag := c.ToNBT()
	if nbtTag.Name() != "Level" {
		t.Errorf("Expected root tag name 'Level', got '%s'", nbtTag.Name())
	}

	if x := nbtTag.GetInt("xPos"); x != 10 {
		t.Errorf("Expected xPos 10, got %d", x)
	}
	if z := nbtTag.GetInt("zPos"); z != 20 {
		t.Errorf("Expected zPos 20, got %d", z)
	}

	c2 := ChunkFromNBT(nbtTag)

	if c2.X != 10 || c2.Z != 20 {
		t.Errorf("Expected coordinates 10,20, got %d,%d", c2.X, c2.Z)
	}

	id, meta := c2.GetBlock(0, 0, 0)
	if id != 1 || meta != 0 {
		t.Errorf("Expected block 1:0 at 0,0,0, got %d:%d", id, meta)
	}

	id, meta = c2.GetBlock(15, 127, 15)
	if id != 7 || meta != 0 {
		t.Errorf("Expected block 7:0 at 15,127,15, got %d:%d", id, meta)
	}

	if len(c2.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(c2.Entities))
	} else {
		if c2.Entities[0].GetString("id") != "Pig" {
			t.Errorf("Expected entity id 'Pig', got '%s'", c2.Entities[0].GetString("id"))
		}
	}
}

func TestChunkSection_NBT(t *testing.T) {

}
