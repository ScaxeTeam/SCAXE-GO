package item

// tools_others.go — 斧/铲/锄 物品（各5种材质，共15个）
// 对应 PHP: item/*Axe.php, *Shovel.php, *Hoe.php (各5个，共15个文件)
//
// PHP 中每个工具继承 Tool，仅提供 ID/Name/AttackPoints/isAxe|isShovel|isHoe()==tier。
// Go 端 tool.go 的 toolData 表已包含所有基础数据。
// 此文件提供斧/铲/锄的专用查询函数。

// ==================== 斧 Axe ====================

// AxeInfo 斧的属性
type AxeInfo struct {
	ID           int
	Name         string
	Tier         int
	AttackDamage float64
	Durability   int
}

var axes = map[int]AxeInfo{
	WOODEN_AXE: {
		ID: WOODEN_AXE, Name: "Wooden Axe",
		Tier: TierWooden, AttackDamage: 3, Durability: 60,
	},
	STONE_AXE: {
		ID: STONE_AXE, Name: "Stone Axe",
		Tier: TierStone, AttackDamage: 4, Durability: 132,
	},
	IRON_AXE: {
		ID: IRON_AXE, Name: "Iron Axe",
		Tier: TierIron, AttackDamage: 5, Durability: 251,
	},
	GOLD_AXE: {
		ID: GOLD_AXE, Name: "Golden Axe",
		Tier: TierGolden, AttackDamage: 3, Durability: 33,
	},
	DIAMOND_AXE: {
		ID: DIAMOND_AXE, Name: "Diamond Axe",
		Tier: TierDiamond, AttackDamage: 6, Durability: 1562,
	},
}

// IsAxe 判断物品是否为斧
func IsAxe(id int) bool {
	_, ok := axes[id]
	return ok
}

// GetAxeInfo 获取斧属性，非斧返回 nil
func GetAxeInfo(id int) *AxeInfo {
	info, ok := axes[id]
	if !ok {
		return nil
	}
	return &info
}

// GetAxeTier 获取斧的材质等级
func GetAxeTier(id int) int {
	if info, ok := axes[id]; ok {
		return info.Tier
	}
	return TierNone
}

// AllAxeIDs 返回所有斧的物品ID
func AllAxeIDs() []int {
	return []int{WOODEN_AXE, STONE_AXE, IRON_AXE, GOLD_AXE, DIAMOND_AXE}
}

// ==================== 铲 Shovel ====================

// ShovelInfo 铲的属性
type ShovelInfo struct {
	ID           int
	Name         string
	Tier         int
	AttackDamage float64
	Durability   int
}

var shovels = map[int]ShovelInfo{
	WOODEN_SHOVEL: {
		ID: WOODEN_SHOVEL, Name: "Wooden Shovel",
		Tier: TierWooden, AttackDamage: 1, Durability: 60,
	},
	STONE_SHOVEL: {
		ID: STONE_SHOVEL, Name: "Stone Shovel",
		Tier: TierStone, AttackDamage: 2, Durability: 132,
	},
	IRON_SHOVEL: {
		ID: IRON_SHOVEL, Name: "Iron Shovel",
		Tier: TierIron, AttackDamage: 3, Durability: 251,
	},
	GOLD_SHOVEL: {
		ID: GOLD_SHOVEL, Name: "Golden Shovel",
		Tier: TierGolden, AttackDamage: 1, Durability: 33,
	},
	DIAMOND_SHOVEL: {
		ID: DIAMOND_SHOVEL, Name: "Diamond Shovel",
		Tier: TierDiamond, AttackDamage: 4, Durability: 1562,
	},
}

// IsShovel 判断物品是否为铲
func IsShovel(id int) bool {
	_, ok := shovels[id]
	return ok
}

// GetShovelInfo 获取铲属性，非铲返回 nil
func GetShovelInfo(id int) *ShovelInfo {
	info, ok := shovels[id]
	if !ok {
		return nil
	}
	return &info
}

// GetShovelTier 获取铲的材质等级
func GetShovelTier(id int) int {
	if info, ok := shovels[id]; ok {
		return info.Tier
	}
	return TierNone
}

// AllShovelIDs 返回所有铲的物品ID
func AllShovelIDs() []int {
	return []int{WOODEN_SHOVEL, STONE_SHOVEL, IRON_SHOVEL, GOLD_SHOVEL, DIAMOND_SHOVEL}
}

// ==================== 锄 Hoe ====================

// HoeInfo 锄的属性
type HoeInfo struct {
	ID           int
	Name         string
	Tier         int
	AttackDamage float64
	Durability   int
}

var hoes = map[int]HoeInfo{
	WOODEN_HOE: {
		ID: WOODEN_HOE, Name: "Wooden Hoe",
		Tier: TierWooden, AttackDamage: 1, Durability: 60,
	},
	STONE_HOE: {
		ID: STONE_HOE, Name: "Stone Hoe",
		Tier: TierStone, AttackDamage: 1, Durability: 132,
	},
	IRON_HOE: {
		ID: IRON_HOE, Name: "Iron Hoe",
		Tier: TierIron, AttackDamage: 1, Durability: 251,
	},
	GOLD_HOE: {
		ID: GOLD_HOE, Name: "Golden Hoe",
		Tier: TierGolden, AttackDamage: 1, Durability: 33,
	},
	DIAMOND_HOE: {
		ID: DIAMOND_HOE, Name: "Diamond Hoe",
		Tier: TierDiamond, AttackDamage: 1, Durability: 1562,
	},
}

// IsHoe 判断物品是否为锄
func IsHoe(id int) bool {
	_, ok := hoes[id]
	return ok
}

// GetHoeInfo 获取锄属性，非锄返回 nil
func GetHoeInfo(id int) *HoeInfo {
	info, ok := hoes[id]
	if !ok {
		return nil
	}
	return &info
}

// GetHoeTier 获取锄的材质等级
func GetHoeTier(id int) int {
	if info, ok := hoes[id]; ok {
		return info.Tier
	}
	return TierNone
}

// AllHoeIDs 返回所有锄的物品ID
func AllHoeIDs() []int {
	return []int{WOODEN_HOE, STONE_HOE, IRON_HOE, GOLD_HOE, DIAMOND_HOE}
}
