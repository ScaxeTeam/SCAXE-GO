package item

// sword.go — 剑类物品（5种材质）
// 对应 PHP: item/WoodenSword.php, StoneSword.php, IronSword.php, GoldSword.php, DiamondSword.php
//
// PHP 中每把剑继承 Tool，仅提供 ID/Name/AttackPoints/isSword()==tier。
// Go 端 tool.go 的 toolData 表已包含所有数据（ToolType/Tier/BaseDamage）。
// 此文件提供剑类专用的便捷查询函数和攻击属性。

// ============ 剑描述数据 ============

// SwordInfo 描述一把剑的详细属性
type SwordInfo struct {
	ID           int     // 物品ID
	Name         string  // 显示名称
	Tier         int     // 材质等级
	AttackDamage float64 // 攻击伤害（= BaseDamage from toolData）
	Durability   int     // 最大耐久
}

// swords 全部5把剑的属性表
var swords = map[int]SwordInfo{
	WOODEN_SWORD: {
		ID:           WOODEN_SWORD,
		Name:         "Wooden Sword",
		Tier:         TierWooden,
		AttackDamage: 4,
		Durability:   60,
	},
	STONE_SWORD: {
		ID:           STONE_SWORD,
		Name:         "Stone Sword",
		Tier:         TierStone,
		AttackDamage: 5,
		Durability:   132,
	},
	IRON_SWORD: {
		ID:           IRON_SWORD,
		Name:         "Iron Sword",
		Tier:         TierIron,
		AttackDamage: 6,
		Durability:   251,
	},
	GOLD_SWORD: {
		ID:           GOLD_SWORD,
		Name:         "Golden Sword",
		Tier:         TierGolden,
		AttackDamage: 4,
		Durability:   33,
	},
	DIAMOND_SWORD: {
		ID:           DIAMOND_SWORD,
		Name:         "Diamond Sword",
		Tier:         TierDiamond,
		AttackDamage: 7,
		Durability:   1562,
	},
}

// ============ 查询函数 ============

// IsSword 判断物品是否为剑
func IsSword(id int) bool {
	_, ok := swords[id]
	return ok
}

// GetSwordInfo 获取剑的详细属性，非剑返回 nil
func GetSwordInfo(id int) *SwordInfo {
	info, ok := swords[id]
	if !ok {
		return nil
	}
	return &info
}

// GetAttackDamage 获取物品的攻击伤害点数
// 对应 PHP Tool::getAttackPoints()
// 剑返回其 AttackDamage，非剑工具使用 toolData 中的 BaseDamage，其他物品返回 1
func GetAttackDamage(id int) float64 {
	if info, ok := swords[id]; ok {
		return info.AttackDamage
	}
	if info := GetToolInfo(id); info != nil {
		return info.BaseDamage
	}
	return 1 // 非工具物品默认 1 点伤害
}

// GetSwordTier 获取剑的材质等级，非剑返回 TierNone
// 对应 PHP Sword::isSword() 返回的 tier
func GetSwordTier(id int) int {
	if info, ok := swords[id]; ok {
		return info.Tier
	}
	return TierNone
}

// AllSwordIDs 返回所有剑的物品ID列表
func AllSwordIDs() []int {
	return []int{WOODEN_SWORD, STONE_SWORD, IRON_SWORD, GOLD_SWORD, DIAMOND_SWORD}
}
