package item

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

const (
	ArmorTierLeather = 1
	ArmorTierGold    = 2
	ArmorTierChain   = 3
	ArmorTierIron    = 4
	ArmorTierDiamond = 5
)

const (
	ArmorSlotHelmet     = 0
	ArmorSlotChestplate = 1
	ArmorSlotLeggings   = 2
	ArmorSlotBoots      = 3
)
type ArmorPieceInfo struct {
	ID         int
	Name       string
	Tier       int
	Slot       int
	Defense    int
	Durability int
}

var armorPieces = map[int]ArmorPieceInfo{
	LEATHER_CAP:   {ID: LEATHER_CAP, Name: "Leather Cap", Tier: ArmorTierLeather, Slot: ArmorSlotHelmet, Defense: 1, Durability: 56},
	LEATHER_TUNIC: {ID: LEATHER_TUNIC, Name: "Leather Tunic", Tier: ArmorTierLeather, Slot: ArmorSlotChestplate, Defense: 3, Durability: 81},
	LEATHER_PANTS: {ID: LEATHER_PANTS, Name: "Leather Pants", Tier: ArmorTierLeather, Slot: ArmorSlotLeggings, Defense: 2, Durability: 76},
	LEATHER_BOOTS: {ID: LEATHER_BOOTS, Name: "Leather Boots", Tier: ArmorTierLeather, Slot: ArmorSlotBoots, Defense: 1, Durability: 66},
	CHAIN_HELMET:     {ID: CHAIN_HELMET, Name: "Chain Helmet", Tier: ArmorTierChain, Slot: ArmorSlotHelmet, Defense: 2, Durability: 166},
	CHAIN_CHESTPLATE: {ID: CHAIN_CHESTPLATE, Name: "Chain Chestplate", Tier: ArmorTierChain, Slot: ArmorSlotChestplate, Defense: 5, Durability: 241},
	CHAIN_LEGGINGS:   {ID: CHAIN_LEGGINGS, Name: "Chain Leggings", Tier: ArmorTierChain, Slot: ArmorSlotLeggings, Defense: 4, Durability: 226},
	CHAIN_BOOTS:      {ID: CHAIN_BOOTS, Name: "Chain Boots", Tier: ArmorTierChain, Slot: ArmorSlotBoots, Defense: 1, Durability: 196},
	IRON_HELMET:     {ID: IRON_HELMET, Name: "Iron Helmet", Tier: ArmorTierIron, Slot: ArmorSlotHelmet, Defense: 2, Durability: 166},
	IRON_CHESTPLATE: {ID: IRON_CHESTPLATE, Name: "Iron Chestplate", Tier: ArmorTierIron, Slot: ArmorSlotChestplate, Defense: 6, Durability: 241},
	IRON_LEGGINGS:   {ID: IRON_LEGGINGS, Name: "Iron Leggings", Tier: ArmorTierIron, Slot: ArmorSlotLeggings, Defense: 5, Durability: 226},
	IRON_BOOTS:      {ID: IRON_BOOTS, Name: "Iron Boots", Tier: ArmorTierIron, Slot: ArmorSlotBoots, Defense: 2, Durability: 196},
	DIAMOND_HELMET:     {ID: DIAMOND_HELMET, Name: "Diamond Helmet", Tier: ArmorTierDiamond, Slot: ArmorSlotHelmet, Defense: 3, Durability: 364},
	DIAMOND_CHESTPLATE: {ID: DIAMOND_CHESTPLATE, Name: "Diamond Chestplate", Tier: ArmorTierDiamond, Slot: ArmorSlotChestplate, Defense: 8, Durability: 529},
	DIAMOND_LEGGINGS:   {ID: DIAMOND_LEGGINGS, Name: "Diamond Leggings", Tier: ArmorTierDiamond, Slot: ArmorSlotLeggings, Defense: 6, Durability: 496},
	DIAMOND_BOOTS:      {ID: DIAMOND_BOOTS, Name: "Diamond Boots", Tier: ArmorTierDiamond, Slot: ArmorSlotBoots, Defense: 3, Durability: 430},
	GOLD_HELMET:     {ID: GOLD_HELMET, Name: "Golden Helmet", Tier: ArmorTierGold, Slot: ArmorSlotHelmet, Defense: 2, Durability: 78},
	GOLD_CHESTPLATE: {ID: GOLD_CHESTPLATE, Name: "Golden Chestplate", Tier: ArmorTierGold, Slot: ArmorSlotChestplate, Defense: 5, Durability: 113},
	GOLD_LEGGINGS:   {ID: GOLD_LEGGINGS, Name: "Golden Leggings", Tier: ArmorTierGold, Slot: ArmorSlotLeggings, Defense: 3, Durability: 106},
	GOLD_BOOTS:      {ID: GOLD_BOOTS, Name: "Golden Boots", Tier: ArmorTierGold, Slot: ArmorSlotBoots, Defense: 1, Durability: 92},
}
func GetArmorPieceInfo(id int) *ArmorPieceInfo {
	info, ok := armorPieces[id]
	if !ok {
		return nil
	}
	return &info
}
func GetArmorTier(id int) int {
	if info, ok := armorPieces[id]; ok {
		return info.Tier
	}
	return 0
}
func GetArmorSlot(id int) int {
	if info, ok := armorPieces[id]; ok {
		return info.Slot
	}
	return -1
}
func IsHelmet(id int) bool {
	return GetArmorSlot(id) == ArmorSlotHelmet
}
func IsChestplate(id int) bool {
	return GetArmorSlot(id) == ArmorSlotChestplate
}
func IsLeggings(id int) bool {
	return GetArmorSlot(id) == ArmorSlotLeggings
}
func IsBoots(id int) bool {
	return GetArmorSlot(id) == ArmorSlotBoots
}
func IsLeatherArmor(id int) bool {
	return GetArmorTier(id) == ArmorTierLeather
}
var armorUnbreakingChance = [4]int{100, 80, 73, 70}
type ArmorUseOnResult struct {
	DamageIncrease int
	IsBroken       bool
}
func UseOnArmor(armorID int, currentDamage int, nbtData *nbt.CompoundTag, enchantUnbreaking int) ArmorUseOnResult {
	if IsUnbreakable(nbtData) {
		return ArmorUseOnResult{DamageIncrease: 0, IsBroken: false}
	}
	level := enchantUnbreaking
	if level < 0 {
		level = 0
	}
	if level > 3 {
		level = 3
	}
	threshold := armorUnbreakingChance[level]
	if rand.Intn(100)+1 > threshold {
		return ArmorUseOnResult{DamageIncrease: 0, IsBroken: false}
	}
	info := GetArmorPieceInfo(armorID)
	if info == nil {
		return ArmorUseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	newDamage := currentDamage + 1
	broken := newDamage >= info.Durability

	return ArmorUseOnResult{DamageIncrease: 1, IsBroken: broken}
}
func CalcTotalDefense(helmetID, chestplateID, leggingsID, bootsID int) int {
	total := 0
	for _, id := range []int{helmetID, chestplateID, leggingsID, bootsID} {
		if info, ok := armorPieces[id]; ok {
			total += info.Defense
		}
	}
	return total
}
func CalcDamageReduction(damage float64, totalDefense int) float64 {
	def := totalDefense
	if def > 20 {
		def = 20
	}
	if def <= 0 {
		return damage
	}
	return damage * (1.0 - float64(def)/25.0)
}
type ArmorColor struct {
	R, G, B uint8
}
func (c ArmorColor) ColorCode() int {
	return int(c.R)<<16 | int(c.G)<<8 | int(c.B)
}
func ColorFromCode(code int) ArmorColor {
	return ArmorColor{
		R: uint8((code >> 16) & 0xFF),
		G: uint8((code >> 8) & 0xFF),
		B: uint8(code & 0xFF),
	}
}
func SetArmorCustomColor(item *Item, color ArmorColor) {
	nbtTag := item.GetNBT()
	nbtTag.Set(nbt.NewIntTag("customColor", int32(color.ColorCode())))
}
func GetArmorCustomColor(item *Item) *ArmorColor {
	if !item.HasNBT() {
		return nil
	}
	val := item.NBTData.GetInt("customColor")
	if val == 0 {
		if !item.NBTData.Has("customColor") {
			return nil
		}
	}
	c := ColorFromCode(int(val))
	return &c
}
func ClearArmorCustomColor(item *Item) {
	if !item.HasNBT() {
		return
	}
	item.NBTData.Remove("customColor")
}
