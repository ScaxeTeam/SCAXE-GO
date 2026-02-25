package item

// armor.go — 盔甲基类（防御值 + 耐久消耗 + 自定义颜色）
// 对应 PHP: item/Armor.php
//
// PHP 中 Armor 继承 Item，提供:
//   - 盔甲等级（皮革/金/锁链/铁/钻石）
//   - 盔甲类型（头盔/胸甲/护腿/靴子）
//   - 防御值
//   - 耐久消耗（与工具不同的概率表）
//   - 自定义颜色（皮革甲专用）
//
// Go 端 properties.go 已有 IsArmor/GetArmorType/GetArmorDefense 基础函数。
// 此文件增加完整的盔甲系统。

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

// ============ 盔甲等级 ============

const (
	ArmorTierLeather = 1
	ArmorTierGold    = 2
	ArmorTierChain   = 3
	ArmorTierIron    = 4
	ArmorTierDiamond = 5
)

// ============ 盔甲槽位 ============

const (
	ArmorSlotHelmet     = 0
	ArmorSlotChestplate = 1
	ArmorSlotLeggings   = 2
	ArmorSlotBoots      = 3
)

// ============ 盔甲详细属性表 ============

// ArmorPieceInfo 单件盔甲的完整属性
type ArmorPieceInfo struct {
	ID         int
	Name       string
	Tier       int // ArmorTier*
	Slot       int // ArmorSlot*
	Defense    int // 防御值
	Durability int // 最大耐久
}

var armorPieces = map[int]ArmorPieceInfo{
	// 皮革
	LEATHER_CAP:   {ID: LEATHER_CAP, Name: "Leather Cap", Tier: ArmorTierLeather, Slot: ArmorSlotHelmet, Defense: 1, Durability: 56},
	LEATHER_TUNIC: {ID: LEATHER_TUNIC, Name: "Leather Tunic", Tier: ArmorTierLeather, Slot: ArmorSlotChestplate, Defense: 3, Durability: 81},
	LEATHER_PANTS: {ID: LEATHER_PANTS, Name: "Leather Pants", Tier: ArmorTierLeather, Slot: ArmorSlotLeggings, Defense: 2, Durability: 76},
	LEATHER_BOOTS: {ID: LEATHER_BOOTS, Name: "Leather Boots", Tier: ArmorTierLeather, Slot: ArmorSlotBoots, Defense: 1, Durability: 66},

	// 锁链
	CHAIN_HELMET:     {ID: CHAIN_HELMET, Name: "Chain Helmet", Tier: ArmorTierChain, Slot: ArmorSlotHelmet, Defense: 2, Durability: 166},
	CHAIN_CHESTPLATE: {ID: CHAIN_CHESTPLATE, Name: "Chain Chestplate", Tier: ArmorTierChain, Slot: ArmorSlotChestplate, Defense: 5, Durability: 241},
	CHAIN_LEGGINGS:   {ID: CHAIN_LEGGINGS, Name: "Chain Leggings", Tier: ArmorTierChain, Slot: ArmorSlotLeggings, Defense: 4, Durability: 226},
	CHAIN_BOOTS:      {ID: CHAIN_BOOTS, Name: "Chain Boots", Tier: ArmorTierChain, Slot: ArmorSlotBoots, Defense: 1, Durability: 196},

	// 铁
	IRON_HELMET:     {ID: IRON_HELMET, Name: "Iron Helmet", Tier: ArmorTierIron, Slot: ArmorSlotHelmet, Defense: 2, Durability: 166},
	IRON_CHESTPLATE: {ID: IRON_CHESTPLATE, Name: "Iron Chestplate", Tier: ArmorTierIron, Slot: ArmorSlotChestplate, Defense: 6, Durability: 241},
	IRON_LEGGINGS:   {ID: IRON_LEGGINGS, Name: "Iron Leggings", Tier: ArmorTierIron, Slot: ArmorSlotLeggings, Defense: 5, Durability: 226},
	IRON_BOOTS:      {ID: IRON_BOOTS, Name: "Iron Boots", Tier: ArmorTierIron, Slot: ArmorSlotBoots, Defense: 2, Durability: 196},

	// 钻石
	DIAMOND_HELMET:     {ID: DIAMOND_HELMET, Name: "Diamond Helmet", Tier: ArmorTierDiamond, Slot: ArmorSlotHelmet, Defense: 3, Durability: 364},
	DIAMOND_CHESTPLATE: {ID: DIAMOND_CHESTPLATE, Name: "Diamond Chestplate", Tier: ArmorTierDiamond, Slot: ArmorSlotChestplate, Defense: 8, Durability: 529},
	DIAMOND_LEGGINGS:   {ID: DIAMOND_LEGGINGS, Name: "Diamond Leggings", Tier: ArmorTierDiamond, Slot: ArmorSlotLeggings, Defense: 6, Durability: 496},
	DIAMOND_BOOTS:      {ID: DIAMOND_BOOTS, Name: "Diamond Boots", Tier: ArmorTierDiamond, Slot: ArmorSlotBoots, Defense: 3, Durability: 430},

	// 金
	GOLD_HELMET:     {ID: GOLD_HELMET, Name: "Golden Helmet", Tier: ArmorTierGold, Slot: ArmorSlotHelmet, Defense: 2, Durability: 78},
	GOLD_CHESTPLATE: {ID: GOLD_CHESTPLATE, Name: "Golden Chestplate", Tier: ArmorTierGold, Slot: ArmorSlotChestplate, Defense: 5, Durability: 113},
	GOLD_LEGGINGS:   {ID: GOLD_LEGGINGS, Name: "Golden Leggings", Tier: ArmorTierGold, Slot: ArmorSlotLeggings, Defense: 3, Durability: 106},
	GOLD_BOOTS:      {ID: GOLD_BOOTS, Name: "Golden Boots", Tier: ArmorTierGold, Slot: ArmorSlotBoots, Defense: 1, Durability: 92},
}

// ============ 查询函数 ============

// GetArmorPieceInfo 获取盔甲件的完整属性，非盔甲返回 nil
func GetArmorPieceInfo(id int) *ArmorPieceInfo {
	info, ok := armorPieces[id]
	if !ok {
		return nil
	}
	return &info
}

// GetArmorTier 获取盔甲的材质等级
func GetArmorTier(id int) int {
	if info, ok := armorPieces[id]; ok {
		return info.Tier
	}
	return 0
}

// GetArmorSlot 获取盔甲的装备槽位（0=头/1=胸/2=腿/3=脚），非盔甲返回 -1
func GetArmorSlot(id int) int {
	if info, ok := armorPieces[id]; ok {
		return info.Slot
	}
	return -1
}

// IsHelmet 判断是否为头盔
func IsHelmet(id int) bool {
	return GetArmorSlot(id) == ArmorSlotHelmet
}

// IsChestplate 判断是否为胸甲
func IsChestplate(id int) bool {
	return GetArmorSlot(id) == ArmorSlotChestplate
}

// IsLeggings 判断是否为护腿
func IsLeggings(id int) bool {
	return GetArmorSlot(id) == ArmorSlotLeggings
}

// IsBoots 判断是否为靴子
func IsBoots(id int) bool {
	return GetArmorSlot(id) == ArmorSlotBoots
}

// IsLeatherArmor 判断是否为皮革甲（支持自定义颜色）
func IsLeatherArmor(id int) bool {
	return GetArmorTier(id) == ArmorTierLeather
}

// ============ 盔甲耐久消耗 ============

// armorUnbreakingChance 盔甲耐久附魔的消耗概率表
// 与工具不同，盔甲有独立的概率表
// 对应 PHP Armor::useOn() 中的 $unbreakings 数组
var armorUnbreakingChance = [4]int{100, 80, 73, 70}

// ArmorUseOnResult 盔甲受击后的耐久结果
type ArmorUseOnResult struct {
	DamageIncrease int  // 耐久消耗增加量
	IsBroken       bool // 盔甲是否损坏
}

// UseOnArmor 盔甲受击时的耐久消耗
// 对应 PHP Armor::useOn()
//
// 参数:
//   - armorID: 盔甲物品ID
//   - currentDamage: 当前耐久消耗 (meta)
//   - nbtData: 物品NBT数据（检查 Unbreakable）
//   - enchantUnbreaking: 耐久附魔等级（0-3）
func UseOnArmor(armorID int, currentDamage int, nbtData *nbt.CompoundTag, enchantUnbreaking int) ArmorUseOnResult {
	// 不可破坏
	if IsUnbreakable(nbtData) {
		return ArmorUseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	// 盔甲专用耐久附魔概率
	level := enchantUnbreaking
	if level < 0 {
		level = 0
	}
	if level > 3 {
		level = 3
	}
	threshold := armorUnbreakingChance[level]
	if rand.Intn(100)+1 > threshold {
		// 概率跳过此次消耗
		return ArmorUseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	// 消耗 1 点耐久
	info := GetArmorPieceInfo(armorID)
	if info == nil {
		return ArmorUseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	newDamage := currentDamage + 1
	broken := newDamage >= info.Durability

	return ArmorUseOnResult{DamageIncrease: 1, IsBroken: broken}
}

// ============ 防御值计算 ============

// CalcTotalDefense 计算全套盔甲的总防御值
// 参数是4个盔甲槽位的物品ID（0表示空槽位）
func CalcTotalDefense(helmetID, chestplateID, leggingsID, bootsID int) int {
	total := 0
	for _, id := range []int{helmetID, chestplateID, leggingsID, bootsID} {
		if info, ok := armorPieces[id]; ok {
			total += info.Defense
		}
	}
	return total
}

// CalcDamageReduction 计算盔甲的伤害减免
// 对应 MC 公式: reducedDamage = damage * (1 - min(20, totalDefense) / 25)
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

// ============ 自定义颜色（皮革甲） ============

// ArmorColor 盔甲自定义颜色
type ArmorColor struct {
	R, G, B uint8
}

// ColorCode 获取颜色的整数编码 (0xRRGGBB)
func (c ArmorColor) ColorCode() int {
	return int(c.R)<<16 | int(c.G)<<8 | int(c.B)
}

// ColorFromCode 从整数编码创建颜色
func ColorFromCode(code int) ArmorColor {
	return ArmorColor{
		R: uint8((code >> 16) & 0xFF),
		G: uint8((code >> 8) & 0xFF),
		B: uint8(code & 0xFF),
	}
}

// SetArmorCustomColor 设置皮革甲的自定义颜色 (NBT "customColor" tag)
// 对应 PHP Armor::setCustomColor()
func SetArmorCustomColor(item *Item, color ArmorColor) {
	nbtTag := item.GetNBT()
	nbtTag.Set(nbt.NewIntTag("customColor", int32(color.ColorCode())))
}

// GetArmorCustomColor 获取皮革甲的自定义颜色，无颜色返回 nil
// 对应 PHP Armor::getCustomColor()
func GetArmorCustomColor(item *Item) *ArmorColor {
	if !item.HasNBT() {
		return nil
	}
	val := item.NBTData.GetInt("customColor")
	if val == 0 {
		// 检查是否真的存在该标签
		if !item.NBTData.Has("customColor") {
			return nil
		}
	}
	c := ColorFromCode(int(val))
	return &c
}

// ClearArmorCustomColor 清除皮革甲的自定义颜色
// 对应 PHP Armor::clearCustomColor()
func ClearArmorCustomColor(item *Item) {
	if !item.HasNBT() {
		return
	}
	item.NBTData.Remove("customColor")
}
