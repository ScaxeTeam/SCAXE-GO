package item

// pickaxe.go — 镐类物品（5种材质）
// 对应 PHP: item/WoodenPickaxe.php, StonePickaxe.php, IronPickaxe.php, GoldPickaxe.php, DiamondPickaxe.php
//
// PHP 中每把镐继承 Tool，仅提供 ID/Name/AttackPoints/isPickaxe()==tier。
// Go 端 tool.go 的 toolData 表已包含所有数据。
// 此文件提供镐类专用的便捷查询函数。

// ============ 镐描述数据 ============

// PickaxeInfo 描述一把镐的详细属性
type PickaxeInfo struct {
	ID           int     // 物品ID
	Name         string  // 显示名称
	Tier         int     // 材质等级
	AttackDamage float64 // 攻击伤害
	Durability   int     // 最大耐久
}

// pickaxes 全部5把镐的属性表
var pickaxes = map[int]PickaxeInfo{
	WOODEN_PICKAXE: {
		ID:           WOODEN_PICKAXE,
		Name:         "Wooden Pickaxe",
		Tier:         TierWooden,
		AttackDamage: 2,
		Durability:   60,
	},
	STONE_PICKAXE: {
		ID:           STONE_PICKAXE,
		Name:         "Stone Pickaxe",
		Tier:         TierStone,
		AttackDamage: 3,
		Durability:   132,
	},
	IRON_PICKAXE: {
		ID:           IRON_PICKAXE,
		Name:         "Iron Pickaxe",
		Tier:         TierIron,
		AttackDamage: 4,
		Durability:   251,
	},
	GOLD_PICKAXE: {
		ID:           GOLD_PICKAXE,
		Name:         "Golden Pickaxe",
		Tier:         TierGolden,
		AttackDamage: 2,
		Durability:   33,
	},
	DIAMOND_PICKAXE: {
		ID:           DIAMOND_PICKAXE,
		Name:         "Diamond Pickaxe",
		Tier:         TierDiamond,
		AttackDamage: 5,
		Durability:   1562,
	},
}

// ============ 查询函数 ============

// IsPickaxe 判断物品是否为镐
func IsPickaxe(id int) bool {
	_, ok := pickaxes[id]
	return ok
}

// GetPickaxeInfo 获取镐的详细属性，非镐返回 nil
func GetPickaxeInfo(id int) *PickaxeInfo {
	info, ok := pickaxes[id]
	if !ok {
		return nil
	}
	return &info
}

// GetPickaxeTier 获取镐的材质等级，非镐返回 TierNone
// 对应 PHP Pickaxe::isPickaxe() 返回的 tier
func GetPickaxeTier(id int) int {
	if info, ok := pickaxes[id]; ok {
		return info.Tier
	}
	return TierNone
}

// CanMine 判断镐是否能挖掘指定等级的方块
// 例如：铁镐（TierIron=4）可以挖掘需要 TierStone（3）的钻石矿
func CanMine(pickaxeID int, requiredTier int) bool {
	tier := GetPickaxeTier(pickaxeID)
	if tier == TierNone {
		return false
	}
	return tier >= requiredTier
}

// AllPickaxeIDs 返回所有镐的物品ID列表
func AllPickaxeIDs() []int {
	return []int{WOODEN_PICKAXE, STONE_PICKAXE, IRON_PICKAXE, GOLD_PICKAXE, DIAMOND_PICKAXE}
}
