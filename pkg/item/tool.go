package item

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

// tool.go — 工具基类增强（耐久系统）
// 对应 PHP: item/Tool.php
//
// 核心功能：
//   - UseOn()         — 破方块/攻击实体时的耐久消耗
//   - UseOnActivate() — 激活（锄头翻地/打火石/铲）时的耐久消耗
//   - IsUnbreakable() — NBT "Unbreakable" 标签检查
//   - GetBlockToolType() — 获取工具对应的方块工具类型
//   - GetMiningEfficiencyFor() — 根据方块类型计算采掘效率

// ============ 常量 ============

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

// UseType 耐久消耗使用类型
const (
	UseTypeBreak    = 1 // 破坏方块 / 攻击实体
	UseTypeActivate = 2 // 激活（锄头/打火石/铲）
)

// BlockIDCobweb 蜘蛛网方块ID（剑对蜘蛛网有加速效果）
const BlockIDCobweb = 30

// ============ 等级数据表 ============

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

// ============ ToolInfo 工具信息 ============

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
	// 打火石和钓鱼竿也算工具
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

// GetBlockToolType 获取工具对应的方块工具类型
// 对应 PHP Tool::getBlockToolType()
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

// ============ 采掘效率 ============

// GetMiningEfficiency 获取工具的基础采掘效率（不考虑方块匹配）
func GetMiningEfficiency(id int) float64 {
	if info := GetToolInfo(id); info != nil {
		if eff, ok := TierEfficiency[info.Tier]; ok {
			return eff
		}
	}
	return 1.0
}

// GetMiningEfficiencyFor 根据方块类型计算实际采掘效率
// 对应 PHP Tool::getMiningEfficiency()
//
// 参数:
//   - toolID: 工具物品ID
//   - blockToolType: 方块所需的工具类型（block.GetToolType()）
//   - blockID: 方块ID（用于特殊检查如蜘蛛网）
//   - enchantEfficiency: 效率附魔等级（0=无附魔）
func GetMiningEfficiencyFor(toolID int, blockToolType int, blockID int, enchantEfficiency int) float64 {
	efficiency := 1.0

	myToolType := GetBlockToolType(toolID)

	// 工具类型匹配时应用基础效率
	if blockToolType != ToolTypeNone && myToolType == blockToolType {
		efficiency = GetMiningEfficiency(toolID)

		// 效率附魔加成
		if enchantEfficiency > 0 {
			efficiency += float64(enchantEfficiency*enchantEfficiency+1) * 2 // HACK: 0.14 附魔判断需要 *2
		}
	}

	// 剑的特殊加速（对应 PHP 的 isSword 分支）
	if GetToolType(toolID) == ToolTypeSword {
		efficiency *= 1.5
		// 剑破蜘蛛网 15x 加速
		if blockID == BlockIDCobweb {
			efficiency *= 10
		}
	}

	return efficiency
}

// ============ 耐久消耗系统 ============

// IsUnbreakable 检查物品是否有 "Unbreakable" NBT 标签
// 对应 PHP Tool::isUnbreakable()
func IsUnbreakable(nbtData *nbt.CompoundTag) bool {
	if nbtData == nil {
		return false
	}
	tag := nbtData.GetCompound("") // root
	if tag == nil {
		tag = nbtData
	}
	val := tag.GetByte("Unbreakable")
	return val > 0
}

// applyUnbreaking 应用耐久附魔概率检查
// 返回 true 表示跳过此次耐久消耗
// 对应 PHP Tool::useOn() 中的 Enchantment::TYPE_MINING_DURABILITY 检查
func applyUnbreaking(enchantUnbreaking int) bool {
	if enchantUnbreaking <= 0 {
		return false
	}
	if enchantUnbreaking > 3 {
		enchantUnbreaking = 3
	}
	// 1/(level+1) 的概率消耗耐久
	return rand.Intn(enchantUnbreaking+1) != 0
}

// UseOnResult 使用工具后的结果
type UseOnResult struct {
	DamageIncrease int  // 耐久增加量（meta += 此值）
	IsBroken       bool // 工具是否已损坏
}

// UseOnBreakBlock 工具破坏方块时的耐久消耗
// 对应 PHP Tool::useOn($block, 1) 中的 Block 分支
//
// 参数:
//   - toolID: 工具物品ID
//   - currentDamage: 当前耐久消耗 (meta)
//   - blockHardness: 被破坏方块的硬度
//   - blockToolType: 被破坏方块所需的工具类型
//   - nbtData: 物品NBT数据（检查 Unbreakable）
//   - enchantUnbreaking: 耐久附魔等级（0=无）
func UseOnBreakBlock(toolID int, currentDamage int, blockHardness float64, blockToolType int, nbtData *nbt.CompoundTag, enchantUnbreaking int) UseOnResult {
	// 不可破坏的工具
	if IsUnbreakable(nbtData) {
		return UseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	// 耐久附魔概率跳过
	if applyUnbreaking(enchantUnbreaking) {
		return UseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	info := GetToolInfo(toolID)
	if info == nil {
		// 非标准工具（打火石等）不因破坏方块消耗耐久
		return UseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	increase := 0

	switch info.ToolType {
	case ToolTypeShears:
		// 剪刀：只有破坏需要剪刀的方块时消耗
		if blockToolType == ToolTypeShears {
			increase = 1
		}
	case ToolTypeSword:
		// 剑：破坏方块消耗 2 耐久（如果方块有硬度）
		if blockHardness > 0 {
			increase = 2
		}
	case ToolTypePickaxe, ToolTypeAxe, ToolTypeShovel:
		// 镐/斧/铲：破坏有硬度的方块消耗 1 耐久
		if blockHardness > 0 {
			increase = 1
		}
	case ToolTypeHoe:
		// 锄头：破坏方块不消耗耐久
		increase = 0
	}

	maxDur := GetMaxDurability(toolID)
	broken := maxDur > 0 && (currentDamage+increase) >= maxDur

	return UseOnResult{DamageIncrease: increase, IsBroken: broken}
}

// UseOnAttackEntity 工具攻击实体时的耐久消耗
// 对应 PHP Tool::useOn($entity, 1) 中的 Entity 分支
//
// 参数:
//   - toolID: 工具物品ID
//   - currentDamage: 当前耐久消耗 (meta)
//   - nbtData: 物品NBT数据
//   - enchantUnbreaking: 耐久附魔等级
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
		// 剑/锄头 攻击实体消耗 1 耐久
		increase = 1
	case ToolTypePickaxe, ToolTypeAxe, ToolTypeShovel:
		// 镐/斧/铲 攻击实体消耗 2 耐久
		increase = 2
	case ToolTypeShears:
		// 剪刀 剪羊毛等消耗 1 耐久
		increase = 1
	default:
		// 其他工具攻击不消耗耐久
		increase = 0
	}

	maxDur := GetMaxDurability(toolID)
	broken := maxDur > 0 && (currentDamage+increase) >= maxDur

	return UseOnResult{DamageIncrease: increase, IsBroken: broken}
}

// UseOnActivate 工具激活使用时的耐久消耗
// 对应 PHP Tool::useOn($object, 2) — Touch/Activate
// 锄头翻地、打火石点火、铲做草径时消耗 1 耐久
//
// 参数:
//   - toolID: 工具物品ID
//   - currentDamage: 当前耐久消耗 (meta)
//   - nbtData: 物品NBT数据
//   - enchantUnbreaking: 耐久附魔等级
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
		// 其他工具激活不消耗耐久
		return UseOnResult{DamageIncrease: 0, IsBroken: false}
	}

	maxDur := GetMaxDurability(toolID)
	broken := maxDur > 0 && (currentDamage+increase) >= maxDur

	return UseOnResult{DamageIncrease: increase, IsBroken: broken}
}

// ============ Item 上的便捷方法 ============

// UseOn 在 Item 上直接调用的耐久消耗便捷方法
// 返回物品是否已损坏
//
// 参数:
//   - useType: UseTypeBreak(破坏方块) 或 UseTypeActivate(激活)
//   - isEntity: 是否是对实体使用（仅 useType=Break 时有效）
//   - blockHardness: 被破坏方块的硬度（仅 isEntity=false 时有效）
//   - blockToolType: 方块所需工具类型（仅 isEntity=false 时有效）
//   - enchantUnbreaking: 耐久附魔等级
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

// ============ 工具数据表 ============

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
