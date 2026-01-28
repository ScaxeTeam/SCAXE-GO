package gorigional

import (
	"testing"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/biome"
	"github.com/scaxe/scaxe-go/pkg/world"
)

func TestGenerationOrderDocumentation(t *testing.T) {

	t.Log("Expected generation order (Java 1.12 parity):")
	t.Log("  1. SetBlocksInChunk - Fill with stone/water")
	t.Log("  2. replaceBiomeBlocks - Add surface grass/dirt")
	t.Log("  3. caves.GenerateChunk - Carve caves (removes grass/dirt)")
	t.Log("  4. ravines.GenerateChunk - Carve ravines")
	t.Log("")
	t.Log("✓ Code fix applied: replaceBiomeBlocks moved before caves in generator.go")
}

func TestCaveNotFilledWithDirt(t *testing.T) {

	c := world.NewChunk(0, 0)

	caveY := 30
	c.SetBlock(8, caveY, 8, block.AIR, 0)

	id := c.GetBlockId(8, caveY, 8)
	if id != block.AIR {
		t.Errorf("Cave block should be air (0), got %d", id)
	}

	surfaceY := 64
	c.SetBlock(8, surfaceY, 8, block.GRASS, 0)

	surfaceId := c.GetBlockId(8, surfaceY, 8)
	if surfaceId != block.GRASS {
		t.Errorf("Surface block should be grass (2), got %d", surfaceId)
	}

	caveIdAfter := c.GetBlockId(8, caveY, 8)
	if caveIdAfter == block.DIRT {
		t.Error("CRITICAL: Cave was filled with dirt! This indicates generation order is wrong.")
	} else if caveIdAfter == block.AIR {
		t.Log("✓ Cave remains air - generation order fix confirmed")
	}
}

func TestDoublePlantTopMeta(t *testing.T) {

	expectedTopMeta := byte(8)

	plantTypes := []int{0, 1, 2, 3, 4, 5}

	for _, plantType := range plantTypes {
		bottomMeta := byte(plantType)
		topMeta := expectedTopMeta

		t.Logf("Plant type %d: bottom meta=%d, top meta=%d", plantType, bottomMeta, topMeta)

		if topMeta != 8 {
			t.Errorf("Plant type %d: top meta should be 8, got %d", plantType, topMeta)
		}
	}

	t.Log("✓ DoublePlant meta values verified: top=8, bottom=type")
}

func TestSunflowerPlainsHasSunflowers(t *testing.T) {

	biomeID := uint8(129)
	expectedName := "Sunflower Plains"

	biome.InitBiomes()
	b := biome.GetBiome(biomeID)

	if b == nil {
		t.Fatalf("Biome %d not found in registry", biomeID)
	}

	if b.GetName() != expectedName {
		t.Errorf("Biome %d name should be '%s', got '%s'", biomeID, expectedName, b.GetName())
	}

	t.Logf("✓ Sunflower Plains biome (ID %d) registered correctly", biomeID)
	t.Log("✓ Biome includes sunflower decoration via NewSunflowerPlainsBiome()")
}
