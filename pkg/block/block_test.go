package block

import (
	"testing"
)

func TestBlockRegistry(t *testing.T) {
	Registry.Init()

	airBehavior := Registry.GetBehavior(AIR)
	if airBehavior == nil {
		t.Fatal("air behavior is nil")
	}
	if airBehavior.GetName() != "Air" {
		t.Errorf("expected name 'Air', got %q", airBehavior.GetName())
	}
	if airBehavior.IsSolid() {
		t.Error("air should not be solid")
	}
	if !airBehavior.CanBeReplaced() {
		t.Error("air should be replaceable")
	}

	stoneBehavior := Registry.GetBehavior(STONE)
	if stoneBehavior == nil {
		t.Fatal("stone behavior is nil")
	}
	if stoneBehavior.GetHardness() != 1.5 {
		t.Errorf("expected hardness 1.5, got %f", stoneBehavior.GetHardness())
	}
	if stoneBehavior.GetToolType() != ToolTypePickaxe {
		t.Errorf("expected tool type pickaxe, got %d", stoneBehavior.GetToolType())
	}
}

func TestBlockState(t *testing.T) {
	state := NewBlockState(STONE, 3)
	if state.ID != STONE {
		t.Errorf("expected ID %d, got %d", STONE, state.ID)
	}
	if state.Meta != 3 {
		t.Errorf("expected meta 3, got %d", state.Meta)
	}

	fullID := state.FullID()
	expectedFullID := (int(STONE) << 4) | 3
	if fullID != expectedFullID {
		t.Errorf("expected fullID %d, got %d", expectedFullID, fullID)
	}
}

func TestRegistryLookups(t *testing.T) {
	Registry.Init()

	if !Registry.IsSolid(STONE) {
		t.Error("stone should be solid")
	}
	if Registry.IsSolid(AIR) {
		t.Error("air should not be solid")
	}

	if Registry.IsTransparent(STONE) {
		t.Error("stone should not be transparent")
	}
	if !Registry.IsTransparent(AIR) {
		t.Error("air should be transparent")
	}

	if Registry.GetLightLevel(TORCH) != 14 {
		t.Errorf("torch should have light level 14, got %d", Registry.GetLightLevel(TORCH))
	}
	if Registry.GetLightLevel(GLOWSTONE_BLOCK) != 15 {
		t.Errorf("glowstone should have light level 15, got %d", Registry.GetLightLevel(GLOWSTONE_BLOCK))
	}
}

func TestStoneDrops(t *testing.T) {
	Registry.Init()

	stoneBehavior := Registry.GetBehavior(STONE)

	drops := stoneBehavior.GetDrops(ToolTypePickaxe, TierWooden)
	if len(drops) != 1 {
		t.Fatalf("expected 1 drop, got %d", len(drops))
	}
	if drops[0].ID != COBBLESTONE {
		t.Errorf("expected cobblestone drop, got block ID %d", drops[0].ID)
	}

	drops = stoneBehavior.GetDrops(ToolTypeShovel, TierDiamond)
	if drops != nil {
		t.Errorf("expected nil drops without pickaxe, got %d drops", len(drops))
	}
}

func TestGrassDrops(t *testing.T) {
	Registry.Init()

	grassBehavior := Registry.GetBehavior(GRASS)

	drops := grassBehavior.GetDrops(ToolTypeShovel, TierWooden)
	if len(drops) != 1 {
		t.Fatalf("expected 1 drop, got %d", len(drops))
	}
	if drops[0].ID != DIRT {
		t.Errorf("expected dirt drop, got block ID %d", drops[0].ID)
	}
}

func TestToolTypes(t *testing.T) {

	if ToolTypeNone != 0 {
		t.Errorf("ToolTypeNone should be 0, got %d", ToolTypeNone)
	}
	if ToolTypePickaxe != 3 {
		t.Errorf("ToolTypePickaxe should be 3, got %d", ToolTypePickaxe)
	}
}

func TestBlockProperties(t *testing.T) {

	stoneProp := GetProperty(STONE)
	if stoneProp.Name != "Stone" {
		t.Errorf("expected name 'Stone', got %q", stoneProp.Name)
	}
	if stoneProp.ToolType != ToolTypePickaxe {
		t.Errorf("expected tool type pickaxe, got %d", stoneProp.ToolType)
	}

	unknownProp := GetProperty(254)
	if unknownProp.Name != "Unknown" {
		t.Errorf("expected name 'Unknown' for undefined block, got %q", unknownProp.Name)
	}
}

func TestBlockBreakInfo(t *testing.T) {
	info := BlockBreakInfo{
		Hardness:         3.0,
		BlastResistance:  15.0,
		RequiredToolType: ToolTypePickaxe,
		RequiredToolTier: TierIron,
	}

	if info.CanHarvestWith(ToolTypeShovel, TierDiamond) {
		t.Error("should not be harvestable with shovel")
	}
	if info.CanHarvestWith(ToolTypePickaxe, TierWooden) {
		t.Error("should not be harvestable with wooden pickaxe")
	}
	if !info.CanHarvestWith(ToolTypePickaxe, TierIron) {
		t.Error("should be harvestable with iron pickaxe")
	}
	if !info.CanHarvestWith(ToolTypePickaxe, TierDiamond) {
		t.Error("should be harvestable with diamond pickaxe")
	}

	breakTime := info.GetBreakTime(ToolTypePickaxe, TierIron, 1.0)
	if breakTime != 4.5 {
		t.Errorf("expected break time 4.5, got %f", breakTime)
	}

	wrongToolTime := info.GetBreakTime(ToolTypeShovel, TierDiamond, 1.0)
	if wrongToolTime != 15.0 {
		t.Errorf("expected break time 15.0 with wrong tool, got %f", wrongToolTime)
	}
}

func TestFuelBlocks(t *testing.T) {

	if !CanBeFuel(PLANKS) {
		t.Error("planks should be fuel")
	}

	if CanBeFuel(STONE) {
		t.Error("stone should not be fuel")
	}

	fuelTime := GetFuelTime(PLANKS)
	if fuelTime != 300 {
		t.Errorf("planks fuel time should be 300, got %d", fuelTime)
	}
}

func TestFlammableBlocks(t *testing.T) {

	if !IsFlammable(PLANKS) {
		t.Error("planks should be flammable")
	}

	if IsFlammable(STONE) {
		t.Error("stone should not be flammable")
	}
}
