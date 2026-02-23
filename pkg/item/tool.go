package item

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

const (
	ToolTypeNone    = 0
	ToolTypeSword   = 1
	ToolTypeShovel  = 2
	ToolTypePickaxe = 3
	ToolTypeAxe     = 4
	ToolTypeShears  = 5
	ToolTypeHoe     = 6
)

const (
	TierNone    = 0
	TierWooden  = 1
	TierGolden  = 2
	TierStone   = 3
	TierIron    = 4
	TierDiamond = 5
)

const (
	UseTypeBreak    = 1
	UseTypeActivate = 2
)

const BlockIDCobweb = 30

var TierDurability = map[int]int{
	TierWooden:  60,
	TierGolden:  33,
	TierStone:   132,
	TierIron:    251,
	TierDiamond: 1562,
}

var TierEfficiency = map[int]float64{
	TierWooden:  2.0,
	TierGolden:  12.0,
	TierStone:   4.0,
	TierIron:    6.0,
	TierDiamond: 8.0,
}

type ToolInfo struct {
	ToolType   int
	Tier       int
	BaseDamage float64
}

func GetToolInfo(id int) *ToolInfo {
	info, ok := toolData[id]
	if !ok {
		return nil
	}
	return &info
}

func IsTool(id int) bool {
	_, ok := toolData[id]
	if ok {
		return true
	}

	return id == FLINT_AND_STEEL || id == FISHING_ROD
}

func GetToolType(id int) int {
	if info := GetToolInfo(id); info != nil {
		return info.ToolType
	}
	return ToolTypeNone
}

func GetToolTier(id int) int {
	if info := GetToolInfo(id); info != nil {
		return info.Tier
	}
	return TierNone
}

func GetMaxDurability(id int) int {
	if info := GetToolInfo(id); info != nil {
		if dur, ok := TierDurability[info.Tier]; ok {
			return dur
		}
	}

	switch id {
	case FLINT_AND_STEEL:
		return 65
	case SHEARS:
		return 239
	case BOW:
		return 385
	case FISHING_ROD:
		return 65
	}
	return 0
}

func GetBlockToolType(id int) int {
	info := GetToolInfo(id)
	if info == nil {
		return ToolTypeNone
	}
	switch info.ToolType {
	case ToolTypePickaxe:
		return ToolTypePickaxe
	case ToolTypeAxe:
		return ToolTypeAxe
	case ToolTypeSword:
		return ToolTypeSword
	case ToolTypeShovel:
		return ToolTypeShovel
	case ToolTypeShears:
		return ToolTypeShears
	default:
		return ToolTypeNone
	}
}

func GetMiningEfficiency(id int) float64 {
	if info := GetToolInfo(id); info != nil {
		if eff, ok := TierEfficiency[info.Tier]; ok {
			return eff
		}
	}
	return 1.0
}

func GetMiningEfficiencyFor(toolID int, blockToolType int, blockID int, enchantEfficiency int) float64 {
	efficiency := 1.0

	myToolType := GetBlockToolType(toolID)

	if blockToolType != ToolTypeNone && myToolType == blockToolType {
		efficiency = GetMiningEfficiency(toolID)

		if enchantEfficiency > 0 {
			efficiency += float64(enchantEfficiency*enchantEfficiency+1) * 2
		}
	}

	if GetToolType(toolID) == ToolTypeSword {
		efficiency *= 1.5

		if blockID == BlockIDCobweb {
			efficiency *= 10
		}
	}

	return efficiency
}

func IsUnbreakable(nbtData *nbt.CompoundTag) bool {
	if nbtData == nil {
		return false
	}
	tag := nbtData.GetCompound("")
	if tag == nil {
		tag = nbtData
	}
	val := tag.GetByte("Unbreakable")
	return val > 0
}

func applyUnbreaking(enchantUnbreaking int) bool {
	if enchantUnbreaking <= 0 {
		return false
	}
	if enchantUnbreaking > 3 {
		enchantUnbreaking = 3
	}

	return rand.Intn(enchantUnbreaking+1) != 0
}

type UseOnResult struct {
	DamageIncrease int
	IsBroken       bool
}

func UseOnBreakBlock(toolID int, currentDamage int, blockHardness float64, blockToolType int, nbtData *nbt.CompoundTag, enchantUnbreaking int) UseOnResult {

	if IsUnbreakable(nbtData) {
		return UseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	if applyUnbreaking(enchantUnbreaking) {
		return UseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	info := GetToolInfo(toolID)
	if info == nil {

		return UseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	increase := 0

	switch info.ToolType {
	case ToolTypeShears:

		if blockToolType == ToolTypeShears {
			increase = 1
		}
	case ToolTypeSword:

		if blockHardness > 0 {
			increase = 2
		}
	case ToolTypePickaxe, ToolTypeAxe, ToolTypeShovel:

		if blockHardness > 0 {
			increase = 1
		}
	case ToolTypeHoe:

		increase = 0
	}

	maxDur := GetMaxDurability(toolID)
	broken := maxDur > 0 && (currentDamage+increase) >= maxDur

	return UseOnResult{DamageIncrease: increase, IsBroken: broken}
}

func UseOnAttackEntity(toolID int, currentDamage int, nbtData *nbt.CompoundTag, enchantUnbreaking int) UseOnResult {
	if IsUnbreakable(nbtData) {
		return UseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	if applyUnbreaking(enchantUnbreaking) {
		return UseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	info := GetToolInfo(toolID)
	if info == nil {
		return UseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	increase := 0

	switch info.ToolType {
	case ToolTypeSword, ToolTypeHoe:

		increase = 1
	case ToolTypePickaxe, ToolTypeAxe, ToolTypeShovel:

		increase = 2
	case ToolTypeShears:

		increase = 1
	default:

		increase = 0
	}

	maxDur := GetMaxDurability(toolID)
	broken := maxDur > 0 && (currentDamage+increase) >= maxDur

	return UseOnResult{DamageIncrease: increase, IsBroken: broken}
}

func UseOnActivate(toolID int, currentDamage int, nbtData *nbt.CompoundTag, enchantUnbreaking int) UseOnResult {
	if IsUnbreakable(nbtData) {
		return UseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	if applyUnbreaking(enchantUnbreaking) {
		return UseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	increase := 0

	toolType := GetToolType(toolID)
	switch {
	case toolType == ToolTypeHoe:
		increase = 1
	case toolType == ToolTypeShovel:
		increase = 1
	case toolID == FLINT_AND_STEEL:
		increase = 1
	default:

		return UseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	maxDur := GetMaxDurability(toolID)
	broken := maxDur > 0 && (currentDamage+increase) >= maxDur

	return UseOnResult{DamageIncrease: increase, IsBroken: broken}
}

func (i *Item) UseOn(useType int, isEntity bool, blockHardness float64, blockToolType int, enchantUnbreaking int) bool {
	var result UseOnResult

	switch useType {
	case UseTypeActivate:
		result = UseOnActivate(i.ID, i.Meta, i.NBTData, enchantUnbreaking)
	case UseTypeBreak:
		if isEntity {
			result = UseOnAttackEntity(i.ID, i.Meta, i.NBTData, enchantUnbreaking)
		} else {
			result = UseOnBreakBlock(i.ID, i.Meta, blockHardness, blockToolType, i.NBTData, enchantUnbreaking)
		}
	default:
		return false
	}

	i.Meta += result.DamageIncrease
	return result.IsBroken
}

var toolData = map[int]ToolInfo{

	WOODEN_SWORD:   {ToolType: ToolTypeSword, Tier: TierWooden, BaseDamage: 4},
	WOODEN_SHOVEL:  {ToolType: ToolTypeShovel, Tier: TierWooden, BaseDamage: 1},
	WOODEN_PICKAXE: {ToolType: ToolTypePickaxe, Tier: TierWooden, BaseDamage: 2},
	WOODEN_AXE:     {ToolType: ToolTypeAxe, Tier: TierWooden, BaseDamage: 3},
	WOODEN_HOE:     {ToolType: ToolTypeHoe, Tier: TierWooden, BaseDamage: 1},

	STONE_SWORD:   {ToolType: ToolTypeSword, Tier: TierStone, BaseDamage: 5},
	STONE_SHOVEL:  {ToolType: ToolTypeShovel, Tier: TierStone, BaseDamage: 2},
	STONE_PICKAXE: {ToolType: ToolTypePickaxe, Tier: TierStone, BaseDamage: 3},
	STONE_AXE:     {ToolType: ToolTypeAxe, Tier: TierStone, BaseDamage: 4},
	STONE_HOE:     {ToolType: ToolTypeHoe, Tier: TierStone, BaseDamage: 1},

	IRON_SWORD:   {ToolType: ToolTypeSword, Tier: TierIron, BaseDamage: 6},
	IRON_SHOVEL:  {ToolType: ToolTypeShovel, Tier: TierIron, BaseDamage: 3},
	IRON_PICKAXE: {ToolType: ToolTypePickaxe, Tier: TierIron, BaseDamage: 4},
	IRON_AXE:     {ToolType: ToolTypeAxe, Tier: TierIron, BaseDamage: 5},
	IRON_HOE:     {ToolType: ToolTypeHoe, Tier: TierIron, BaseDamage: 1},

	GOLD_SWORD:   {ToolType: ToolTypeSword, Tier: TierGolden, BaseDamage: 4},
	GOLD_SHOVEL:  {ToolType: ToolTypeShovel, Tier: TierGolden, BaseDamage: 1},
	GOLD_PICKAXE: {ToolType: ToolTypePickaxe, Tier: TierGolden, BaseDamage: 2},
	GOLD_AXE:     {ToolType: ToolTypeAxe, Tier: TierGolden, BaseDamage: 3},
	GOLD_HOE:     {ToolType: ToolTypeHoe, Tier: TierGolden, BaseDamage: 1},

	DIAMOND_SWORD:   {ToolType: ToolTypeSword, Tier: TierDiamond, BaseDamage: 7},
	DIAMOND_SHOVEL:  {ToolType: ToolTypeShovel, Tier: TierDiamond, BaseDamage: 4},
	DIAMOND_PICKAXE: {ToolType: ToolTypePickaxe, Tier: TierDiamond, BaseDamage: 5},
	DIAMOND_AXE:     {ToolType: ToolTypeAxe, Tier: TierDiamond, BaseDamage: 6},
	DIAMOND_HOE:     {ToolType: ToolTypeHoe, Tier: TierDiamond, BaseDamage: 1},

	SHEARS: {ToolType: ToolTypeShears, Tier: TierNone, BaseDamage: 1},
	BOW:    {ToolType: ToolTypeNone, Tier: TierNone, BaseDamage: 0},
}
