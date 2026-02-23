package item

type PickaxeInfo struct {
	ID           int
	Name         string
	Tier         int
	AttackDamage float64
	Durability   int
}

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

func IsPickaxe(id int) bool {
	_, ok := pickaxes[id]
	return ok
}

func GetPickaxeInfo(id int) *PickaxeInfo {
	info, ok := pickaxes[id]
	if !ok {
		return nil
	}
	return &info
}

func GetPickaxeTier(id int) int {
	if info, ok := pickaxes[id]; ok {
		return info.Tier
	}
	return TierNone
}

func CanMine(pickaxeID int, requiredTier int) bool {
	tier := GetPickaxeTier(pickaxeID)
	if tier == TierNone {
		return false
	}
	return tier >= requiredTier
}

func AllPickaxeIDs() []int {
	return []int{WOODEN_PICKAXE, STONE_PICKAXE, IRON_PICKAXE, GOLD_PICKAXE, DIAMOND_PICKAXE}
}
