package biome

import (
	"testing"

	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
)

func TestIcePlainsSpikesProperties(t *testing.T) {
	spikeBiome := NewIcePlainsSpikesBiome()

	if spikeBiome.ID != 140 {
		t.Errorf("ID = %d, expected 140", spikeBiome.ID)
	}
	if spikeBiome.Name != "Ice Plains Spikes" {
		t.Errorf("Name = %s, expected 'Ice Plains Spikes'", spikeBiome.Name)
	}

	if spikeBiome.BaseHeight != 0.425 {
		t.Errorf("BaseHeight = %f, expected 0.425", spikeBiome.BaseHeight)
	}
	if spikeBiome.HeightVariation != 0.45 {
		t.Errorf("HeightVariation = %f, expected 0.45", spikeBiome.HeightVariation)
	}

	if spikeBiome.Decorator.TreesPerChunk != 0 {
		t.Errorf("TreesPerChunk = %d, expected 0", spikeBiome.Decorator.TreesPerChunk)
	}

	t.Logf("IcePlainsSpikesBiome: ID=%d, Name=%s, Height=%.3f",
		spikeBiome.ID, spikeBiome.Name, spikeBiome.BaseHeight)
}

func TestIceSpikeGeneration(t *testing.T) {

	spike := object.NewIceSpike()

	if spike == nil {
		t.Fatal("NewIceSpike returned nil")
	}
}

func TestIcePathWidth(t *testing.T) {
	width := 4
	path := object.NewIcePath(width)

	if path.BasePathWidth != width {
		t.Errorf("BasePathWidth = %d, expected %d", path.BasePathWidth, width)
	}
}
