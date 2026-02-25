package item
type SwordInfo struct {
	ID           int
	Name         string
	Tier         int
	AttackDamage float64
	Durability   int
}
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
func IsSword(id int) bool {
	_, ok := swords[id]
	return ok
}
func GetSwordInfo(id int) *SwordInfo {
	info, ok := swords[id]
	if !ok {
		return nil
	}
	return &info
}
func GetAttackDamage(id int) float64 {
	if info, ok := swords[id]; ok {
		return info.AttackDamage
	}
	if info := GetToolInfo(id); info != nil {
		return info.BaseDamage
	}
	return 1
}
func GetSwordTier(id int) int {
	if info, ok := swords[id]; ok {
		return info.Tier
	}
	return TierNone
}
func AllSwordIDs() []int {
	return []int{WOODEN_SWORD, STONE_SWORD, IRON_SWORD, GOLD_SWORD, DIAMOND_SWORD}
}
