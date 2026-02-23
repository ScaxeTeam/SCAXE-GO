package item

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

func IsAxe(id int) bool {
	_, ok := axes[id]
	return ok
}

func GetAxeInfo(id int) *AxeInfo {
	info, ok := axes[id]
	if !ok {
		return nil
	}
	return &info
}

func GetAxeTier(id int) int {
	if info, ok := axes[id]; ok {
		return info.Tier
	}
	return TierNone
}

func AllAxeIDs() []int {
	return []int{WOODEN_AXE, STONE_AXE, IRON_AXE, GOLD_AXE, DIAMOND_AXE}
}

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

func IsShovel(id int) bool {
	_, ok := shovels[id]
	return ok
}

func GetShovelInfo(id int) *ShovelInfo {
	info, ok := shovels[id]
	if !ok {
		return nil
	}
	return &info
}

func GetShovelTier(id int) int {
	if info, ok := shovels[id]; ok {
		return info.Tier
	}
	return TierNone
}

func AllShovelIDs() []int {
	return []int{WOODEN_SHOVEL, STONE_SHOVEL, IRON_SHOVEL, GOLD_SHOVEL, DIAMOND_SHOVEL}
}

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

func IsHoe(id int) bool {
	_, ok := hoes[id]
	return ok
}

func GetHoeInfo(id int) *HoeInfo {
	info, ok := hoes[id]
	if !ok {
		return nil
	}
	return &info
}

func GetHoeTier(id int) int {
	if info, ok := hoes[id]; ok {
		return info.Tier
	}
	return TierNone
}

func AllHoeIDs() []int {
	return []int{WOODEN_HOE, STONE_HOE, IRON_HOE, GOLD_HOE, DIAMOND_HOE}
}
