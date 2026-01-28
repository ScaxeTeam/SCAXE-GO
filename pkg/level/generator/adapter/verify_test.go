package adapter

import (
	"testing"

	"github.com/scaxe/scaxe-go/pkg/block"
)

func TestBlockParity(t *testing.T) {
	tests := []struct {
		name     string
		adapter  BlockState
		expected uint8
	}{
		{"Air", AIR, block.AIR},
		{"Stone", STONE, block.STONE},
		{"Grass", GRASS, block.GRASS},
		{"Dirt", DIRT, block.DIRT},
		{"Cobblestone", COBBLESTONE, block.COBBLESTONE},
		{"Planks", PLANKS, block.PLANKS},
		{"Bedrock", BEDROCK, block.BEDROCK},
		{"FlowingWater", FLOWING_WATER, block.WATER},
		{"Water", WATER, block.STILL_WATER},
		{"FlowingLava", FLOWING_LAVA, block.LAVA},
		{"Lava", LAVA, block.STILL_LAVA},

		{"DoubleWoodenSlab", DOUBLE_WOODEN_SLAB, block.DOUBLE_WOOD_SLAB},
		{"WoodenSlab", WOODEN_SLAB, block.WOOD_SLAB},
		{"Glass", GLASS, block.GLASS},
		{"Torch", TORCH, block.TORCH},
	}

	for _, tt := range tests {
		if tt.adapter.ID != tt.expected {
			t.Errorf("Block %s ID mismatch: Adapter=%d, Scaxe=%d", tt.name, tt.adapter.ID, tt.expected)
		}
	}

	if LAVA.ID != 11 {
		t.Logf("Warning: Adapter LAVA is %d", LAVA.ID)
	}
	if block.STILL_LAVA != 11 {
		t.Logf("Warning: Scaxe STILL_LAVA is %d", block.STILL_LAVA)
	}
}
