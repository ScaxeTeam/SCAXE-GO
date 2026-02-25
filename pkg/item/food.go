package item

// food.go — 食物系统（FoodSource 接口 + 饥饿系统）
// 对应 PHP: item/FoodSource.php (接口), item/Food.php (基类)
//
// PHP FoodSource 接口定义: getFoodRestore/getSaturationRestore/getResidue/getAdditionalEffects
// PHP Food 基类提供: canBeConsumed/canBeConsumedBy/onConsume
//
// Go 端 properties.go 已有 IsFood/GetFoodRestore/GetSaturation 基础查询。
// 此文件提供完整的食物系统：详细属性表、残留物、附加效果、食用结果计算。

// ============ 食物属性 ============

// FoodInfo 食物的完整属性
type FoodInfo struct {
	ID             int     // 物品ID
	Name           string  // 显示名称
	FoodRestore    int     // 饱食度恢复量
	Saturation     float32 // 饱和度恢复量
	ResidueID      int     // 食用后残留物品ID（0=无）
	ResidueMeta    int     // 残留物品 meta
	EffectChance   float32 // 附加效果触发概率 (0.0~1.0, 0=无效果)
	EffectID       int     // 附加效果ID（药水效果ID）
	EffectDuration int     // 附加效果持续时间（tick）
	EffectLevel    int     // 附加效果等级
	AlwaysEdible   bool    // 是否始终可食用（如金苹果，不管饱食度）
}

// foods 全部食物属性表
var foods = map[int]FoodInfo{
	// 基础食物
	APPLE:        {ID: APPLE, Name: "Apple", FoodRestore: 4, Saturation: 2.4},
	BREAD:        {ID: BREAD, Name: "Bread", FoodRestore: 5, Saturation: 6.0},
	COOKIE:       {ID: COOKIE, Name: "Cookie", FoodRestore: 2, Saturation: 0.4},
	MELON:        {ID: MELON, Name: "Melon Slice", FoodRestore: 2, Saturation: 1.2},
	CARROT:       {ID: CARROT, Name: "Carrot", FoodRestore: 3, Saturation: 3.6},
	POTATO:       {ID: POTATO, Name: "Potato", FoodRestore: 1, Saturation: 0.6},
	BAKED_POTATO: {ID: BAKED_POTATO, Name: "Baked Potato", FoodRestore: 5, Saturation: 6.0},
	BEETROOT:     {ID: BEETROOT, Name: "Beetroot", FoodRestore: 1, Saturation: 1.2},
	PUMPKIN_PIE:  {ID: PUMPKIN_PIE, Name: "Pumpkin Pie", FoodRestore: 8, Saturation: 4.8},

	// 肉类（生）
	RAW_PORKCHOP: {ID: RAW_PORKCHOP, Name: "Raw Porkchop", FoodRestore: 3, Saturation: 1.8},
	RAW_BEEF:     {ID: RAW_BEEF, Name: "Raw Beef", FoodRestore: 3, Saturation: 1.8},
	RAW_CHICKEN: {ID: RAW_CHICKEN, Name: "Raw Chicken", FoodRestore: 2, Saturation: 1.2,
		EffectChance: 0.3, EffectID: 17, EffectDuration: 600, EffectLevel: 0}, // 30% 饥饿效果
	RAW_RABBIT: {ID: RAW_RABBIT, Name: "Raw Rabbit", FoodRestore: 3, Saturation: 1.8},
	RAW_FISH:   {ID: RAW_FISH, Name: "Raw Fish", FoodRestore: 2, Saturation: 0.4},

	// 肉类（熟）
	COOKED_PORKCHOP: {ID: COOKED_PORKCHOP, Name: "Cooked Porkchop", FoodRestore: 8, Saturation: 12.8},
	STEAK:           {ID: STEAK, Name: "Steak", FoodRestore: 8, Saturation: 12.8},
	COOKED_CHICKEN:  {ID: COOKED_CHICKEN, Name: "Cooked Chicken", FoodRestore: 6, Saturation: 7.2},
	COOKED_RABBIT:   {ID: COOKED_RABBIT, Name: "Cooked Rabbit", FoodRestore: 5, Saturation: 6.0},
	COOKED_FISH:     {ID: COOKED_FISH, Name: "Cooked Fish", FoodRestore: 5, Saturation: 6.0},

	// 特殊食物
	GOLDEN_APPLE:  {ID: GOLDEN_APPLE, Name: "Golden Apple", FoodRestore: 4, Saturation: 9.6, AlwaysEdible: true},
	GOLDEN_CARROT: {ID: GOLDEN_CARROT, Name: "Golden Carrot", FoodRestore: 6, Saturation: 14.4},
	ROTTEN_FLESH: {ID: ROTTEN_FLESH, Name: "Rotten Flesh", FoodRestore: 4, Saturation: 0.8,
		EffectChance: 0.8, EffectID: 17, EffectDuration: 600, EffectLevel: 0}, // 80% 饥饿效果
	POISONOUS_POTATO: {ID: POISONOUS_POTATO, Name: "Poisonous Potato", FoodRestore: 2, Saturation: 1.2,
		EffectChance: 0.6, EffectID: 19, EffectDuration: 100, EffectLevel: 0}, // 60% 中毒

	// 容器食物（有残留物品）
	MUSHROOM_STEW: {ID: MUSHROOM_STEW, Name: "Mushroom Stew", FoodRestore: 6, Saturation: 7.2,
		ResidueID: BOWL},
	RABBIT_STEW: {ID: RABBIT_STEW, Name: "Rabbit Stew", FoodRestore: 10, Saturation: 12.0,
		ResidueID: BOWL},
	BEETROOT_SOUP: {ID: BEETROOT_SOUP, Name: "Beetroot Soup", FoodRestore: 6, Saturation: 7.2,
		ResidueID: BOWL},

	// 鱼类（额外变种）
	RAW_SALMON:    {ID: RAW_SALMON, Name: "Raw Salmon", FoodRestore: 2, Saturation: 0.4},
	COOKED_SALMON: {ID: COOKED_SALMON, Name: "Cooked Salmon", FoodRestore: 6, Saturation: 9.6},
	CLOWN_FISH:    {ID: CLOWN_FISH, Name: "Clownfish", FoodRestore: 1, Saturation: 0.2},
	PUFFER_FISH: {ID: PUFFER_FISH, Name: "Pufferfish", FoodRestore: 1, Saturation: 0.2,
		EffectChance: 1.0, EffectID: 19, EffectDuration: 1200, EffectLevel: 3}, // 100% 中毒IV 60秒

	// 其他可食用物品
	SPIDER_EYE: {ID: SPIDER_EYE, Name: "Spider Eye", FoodRestore: 2, Saturation: 3.2,
		EffectChance: 1.0, EffectID: 19, EffectDuration: 100, EffectLevel: 0}, // 100% 中毒
	ENCHANTED_GOLDEN_APPLE: {ID: ENCHANTED_GOLDEN_APPLE, Name: "Enchanted Golden Apple",
		FoodRestore: 4, Saturation: 9.6, AlwaysEdible: true},
}

// ============ 查询函数 ============

// GetFoodInfo 获取食物的完整属性，非食物返回 nil
func GetFoodInfo(id int) *FoodInfo {
	info, ok := foods[id]
	if !ok {
		return nil
	}
	return &info
}

// IsFoodItem 判断物品是否为食物（比 IsFood 更详细的版本）
func IsFoodItem(id int) bool {
	_, ok := foods[id]
	return ok
}

// ============ 食用逻辑 ============

// EatResult 食用食物后的结果
type EatResult struct {
	FoodRestore    int     // 饱食度恢复量
	Saturation     float32 // 饱和度恢复量
	ResidueItem    *Item   // 残留物品（nil=无）
	HasEffect      bool    // 是否触发附加效果
	EffectID       int     // 效果ID
	EffectDuration int     // 效果持续时间(tick)
	EffectLevel    int     // 效果等级
}

// CanEat 判断物品是否可以被食用
// 参数:
//   - itemID: 物品ID
//   - currentFood: 当前饱食度
//   - maxFood: 最大饱食度
func CanEat(itemID int, currentFood int, maxFood int) bool {
	info := GetFoodInfo(itemID)
	if info == nil {
		return false
	}
	if info.AlwaysEdible {
		return true
	}
	return currentFood < maxFood
}

// Eat 计算食用食物的结果
// 对应 PHP Food::onConsume()
func Eat(item Item) *EatResult {
	info := GetFoodInfo(item.ID)
	if info == nil {
		return nil
	}

	result := &EatResult{
		FoodRestore: info.FoodRestore,
		Saturation:  info.Saturation,
	}

	// 残留物品
	if info.ResidueID > 0 {
		residue := NewItem(info.ResidueID, info.ResidueMeta, 1)
		result.ResidueItem = &residue
	}

	// 附加效果概率
	if info.EffectChance > 0 && info.EffectID > 0 {
		// 使用简单随机判断
		// 注意: 实际调用方应在传入前判断随机
		result.HasEffect = true
		result.EffectID = info.EffectID
		result.EffectDuration = info.EffectDuration
		result.EffectLevel = info.EffectLevel
	}

	return result
}

// ============ 饥饿系统常量 ============

const (
	MaxFood       = 20   // 最大饱食度
	MaxSaturation = 20.0 // 最大饱和度（等于当前饱食度）
	FoodTickTimer = 80   // 饥饿消耗间隔（tick）

	// 饥饿消耗速率
	ExhaustionActionAttack = 0.1   // 攻击
	ExhaustionActionDamage = 0.1   // 受伤
	ExhaustionActionMine   = 0.005 // 挖掘方块
	ExhaustionActionSprint = 0.1   // 疾跑（每米）
	ExhaustionActionJump   = 0.05  // 跳跃
	ExhaustionActionSwim   = 0.01  // 游泳（每米）
	ExhaustionActionWalk   = 0.0   // 走路（不消耗）
	ExhaustionActionRegen  = 6.0   // 自然回血
)

// HungerResult 饥饿系统 tick 结果
type HungerResult struct {
	ShouldHeal   bool // 是否应该自然回血（饱食度>=18）
	ShouldDamage bool // 是否应该饥饿伤害（饱食度<=0）
}

// CalcHungerTick 计算饥饿 tick 效果
// 对应 PHP Player.php 中的饥饿消耗逻辑
func CalcHungerTick(currentFood int, saturation float32) HungerResult {
	return HungerResult{
		ShouldHeal:   currentFood >= 18 && saturation > 0,
		ShouldDamage: currentFood <= 0,
	}
}

// AllFoodIDs 返回所有食物的物品ID列表
func AllFoodIDs() []int {
	ids := make([]int, 0, len(foods))
	for id := range foods {
		ids = append(ids, id)
	}
	return ids
}
