package item

import (
	"testing"
)

func TestNewItem(t *testing.T) {
	item := NewItem(IRON_PICKAXE, 0, 1)
	if item.ID != IRON_PICKAXE {
		t.Errorf("expected ID %d, got %d", IRON_PICKAXE, item.ID)
	}
	if item.Meta != 0 {
		t.Errorf("expected Meta 0, got %d", item.Meta)
	}
	if item.Count != 1 {
		t.Errorf("expected Count 1, got %d", item.Count)
	}
}

func TestAir(t *testing.T) {
	air := Air()
	if !air.IsAir() {
		t.Error("Air() should return an air item")
	}
	if air.ID != 0 {
		t.Errorf("Air ID should be 0, got %d", air.ID)
	}
}

func TestIsAir(t *testing.T) {

	if !NewItem(0, 0, 1).IsAir() {
		t.Error("ID 0 should be air")
	}

	if !NewItem(DIAMOND, 0, 0).IsAir() {
		t.Error("Count 0 should be air")
	}

	if !NewItem(DIAMOND, 0, -1).IsAir() {
		t.Error("Negative count should be air")
	}

	if NewItem(DIAMOND, 0, 1).IsAir() {
		t.Error("Valid item should not be air")
	}
}

func TestMaxStackSize(t *testing.T) {

	if NewItem(DIAMOND_PICKAXE, 0, 1).GetMaxStackSize() != 1 {
		t.Error("tools should stack to 1")
	}

	if NewItem(IRON_CHESTPLATE, 0, 1).GetMaxStackSize() != 1 {
		t.Error("armor should stack to 1")
	}

	if NewItem(EGG, 0, 1).GetMaxStackSize() != 16 {
		t.Error("eggs should stack to 16")
	}

	if NewItem(DIAMOND, 0, 1).GetMaxStackSize() != 64 {
		t.Error("diamonds should stack to 64")
	}
}

func TestToolInfo(t *testing.T) {

	item := NewItem(IRON_PICKAXE, 0, 1)
	if !item.IsTool() {
		t.Error("iron pickaxe should be a tool")
	}
	if item.GetToolType() != ToolTypePickaxe {
		t.Errorf("expected tool type pickaxe, got %d", item.GetToolType())
	}
	if item.GetToolTier() != TierIron {
		t.Errorf("expected tier iron, got %d", item.GetToolTier())
	}
	if item.GetMaxDurability() != 251 {
		t.Errorf("expected durability 251, got %d", item.GetMaxDurability())
	}

	sword := NewItem(DIAMOND_SWORD, 0, 1)
	if sword.GetToolType() != ToolTypeSword {
		t.Errorf("expected tool type sword, got %d", sword.GetToolType())
	}
	if sword.GetToolTier() != TierDiamond {
		t.Errorf("expected tier diamond, got %d", sword.GetToolTier())
	}
}

func TestToolData(t *testing.T) {
	info := GetToolInfo(IRON_PICKAXE)
	if info == nil {
		t.Fatal("expected tool info for iron pickaxe")
	}
	if info.ToolType != ToolTypePickaxe {
		t.Errorf("expected pickaxe type, got %d", info.ToolType)
	}
	if info.Tier != TierIron {
		t.Errorf("expected iron tier, got %d", info.Tier)
	}
	if info.BaseDamage != 4 {
		t.Errorf("expected base damage 4, got %f", info.BaseDamage)
	}
}

func TestIsArmor(t *testing.T) {

	if !IsArmor(IRON_CHESTPLATE) {
		t.Error("iron chestplate should be armor")
	}
	if !IsArmor(DIAMOND_BOOTS) {
		t.Error("diamond boots should be armor")
	}

	if IsArmor(IRON_PICKAXE) {
		t.Error("iron pickaxe should not be armor")
	}
	if IsArmor(DIAMOND) {
		t.Error("diamond should not be armor")
	}
}

func TestArmorType(t *testing.T) {
	if GetArmorType(IRON_HELMET) != 0 {
		t.Errorf("helmet should be slot 0, got %d", GetArmorType(IRON_HELMET))
	}
	if GetArmorType(IRON_CHESTPLATE) != 1 {
		t.Errorf("chestplate should be slot 1, got %d", GetArmorType(IRON_CHESTPLATE))
	}
	if GetArmorType(IRON_LEGGINGS) != 2 {
		t.Errorf("leggings should be slot 2, got %d", GetArmorType(IRON_LEGGINGS))
	}
	if GetArmorType(IRON_BOOTS) != 3 {
		t.Errorf("boots should be slot 3, got %d", GetArmorType(IRON_BOOTS))
	}
}

func TestArmorDefense(t *testing.T) {

	if GetArmorDefense(IRON_CHESTPLATE) != 6 {
		t.Errorf("iron chestplate should have 6 defense, got %d", GetArmorDefense(IRON_CHESTPLATE))
	}

	if GetArmorDefense(DIAMOND_HELMET) != 3 {
		t.Errorf("diamond helmet should have 3 defense, got %d", GetArmorDefense(DIAMOND_HELMET))
	}
}

func TestIsFood(t *testing.T) {
	if !IsFood(APPLE) {
		t.Error("apple should be food")
	}
	if !IsFood(STEAK) {
		t.Error("steak should be food")
	}
	if IsFood(DIAMOND) {
		t.Error("diamond should not be food")
	}
}

func TestFoodRestore(t *testing.T) {
	if GetFoodRestore(APPLE) != 4 {
		t.Errorf("apple should restore 4 hunger, got %d", GetFoodRestore(APPLE))
	}
	if GetFoodRestore(STEAK) != 8 {
		t.Errorf("steak should restore 8 hunger, got %d", GetFoodRestore(STEAK))
	}
}

func TestItemName(t *testing.T) {
	if GetItemName(DIAMOND) != "Diamond" {
		t.Errorf("expected 'Diamond', got %q", GetItemName(DIAMOND))
	}
	if GetItemName(IRON_PICKAXE) != "Iron Pickaxe" {
		t.Errorf("expected 'Iron Pickaxe', got %q", GetItemName(IRON_PICKAXE))
	}
}

func TestItemClone(t *testing.T) {
	original := NewItem(DIAMOND, 0, 64)
	clone := original.Clone()

	if clone.ID != original.ID {
		t.Error("cloned ID should match")
	}
	if clone.Count != original.Count {
		t.Error("cloned count should match")
	}

	clone.Count = 32
	if original.Count == 32 {
		t.Error("modifying clone should not affect original")
	}
}

func TestItemEquals(t *testing.T) {
	item1 := NewItem(DIAMOND, 0, 64)
	item2 := NewItem(DIAMOND, 0, 32)
	item3 := NewItem(DIAMOND, 1, 64)
	item4 := NewItem(EMERALD, 0, 64)

	if !item1.Equals(item2, false, false) {
		t.Error("should equal without meta check")
	}
	if !item1.Equals(item2, true, false) {
		t.Error("should equal with meta check (same meta)")
	}

	if !item1.Equals(item3, false, false) {
		t.Error("should equal without meta check")
	}
	if item1.Equals(item3, true, false) {
		t.Error("should not equal with meta check (different meta)")
	}

	if item1.Equals(item4, false, false) {
		t.Error("different IDs should not be equal")
	}
}

func TestMiningEfficiency(t *testing.T) {

	if GetMiningEfficiency(DIAMOND_PICKAXE) != 8.0 {
		t.Errorf("diamond pickaxe efficiency should be 8, got %f", GetMiningEfficiency(DIAMOND_PICKAXE))
	}

	if GetMiningEfficiency(GOLD_PICKAXE) != 12.0 {
		t.Errorf("gold pickaxe efficiency should be 12, got %f", GetMiningEfficiency(GOLD_PICKAXE))
	}

	if GetMiningEfficiency(DIAMOND) != 1.0 {
		t.Errorf("non-tool efficiency should be 1, got %f", GetMiningEfficiency(DIAMOND))
	}
}

func TestTierDurability(t *testing.T) {

	if TierDurability[TierWooden] != 60 {
		t.Errorf("wooden tier durability should be 60, got %d", TierDurability[TierWooden])
	}
	if TierDurability[TierDiamond] != 1562 {
		t.Errorf("diamond tier durability should be 1562, got %d", TierDurability[TierDiamond])
	}
}
