package block

// natural.go — 自然方块（泥土/草方块/沙子/砂砾）
// 对应 PHP: Dirt, Grass, Sand, Gravel
//
// 这些方块都存在于 blocks.go 的旧注册中（dirtBlock, grassBlock 等），
// 此文件提供增强的方块行为，包括交互、随机 tick、特殊掉落等。

// ============ Dirt 泥土 ============

// DirtBlock 泥土方块
// 对应 PHP class Dirt extends Solid
type DirtBlock struct {
	SolidBase
}

func NewDirtBlock() *DirtBlock {
	return &DirtBlock{
		SolidBase: SolidBase{
			BlockID:       DIRT,
			BlockName:     "Dirt",
			BlockHardness: 0.5,
			BlockToolType: ToolTypeShovel,
		},
	}
}

// CanBeActivated 泥土可以被锄头激活（→耕地）
func (b *DirtBlock) CanBeActivated() bool {
	return true
}

// OnActivateResult 激活结果
type OnActivateResult struct {
	Handled      bool  // 是否处理了
	ReplaceBlock uint8 // 替换为什么方块（0=不替换）
	UseTool      bool  // 是否消耗工具耐久
}

// OnDirtActivate 泥土被激活时的处理
// 对应 PHP Dirt::onActivate()
// 锄头右键 → 变成耕地
func OnDirtActivate(isHoe bool) OnActivateResult {
	if isHoe {
		return OnActivateResult{Handled: true, ReplaceBlock: FARMLAND, UseTool: true}
	}
	return OnActivateResult{Handled: false}
}

// ============ Grass 草方块 ============

// GrassBlock 草方块
// 对应 PHP class Grass extends Solid
type GrassBlock struct {
	SolidBase
}

func NewGrassBlock() *GrassBlock {
	return &GrassBlock{
		SolidBase: SolidBase{
			BlockID:       GRASS,
			BlockName:     "Grass",
			BlockHardness: 0.6,
			BlockToolType: ToolTypeShovel,
		},
	}
}

// CanBeActivated 草方块可以被激活（骨粉/锄头/铲）
func (b *GrassBlock) CanBeActivated() bool {
	return true
}

// GetDrops 草方块默认掉落泥土（精准采集掉落自身，由 Level 层处理）
func (b *GrassBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(DIRT), Meta: 0, Count: 1}}
}

// ---------- 草方块交互 ----------

// GrassActivateType 草方块激活类型
type GrassActivateType uint8

const (
	GrassActivateNone     GrassActivateType = iota
	GrassActivateBoneMeal                   // 骨粉催草
	GrassActivateHoe                        // 锄头→耕地
	GrassActivateShovel                     // 铲→草径
)

// OnGrassActivate 草方块被激活时的处理
// 对应 PHP Grass::onActivate()
func OnGrassActivate(isBoneMeal bool, isHoe bool, isShovel bool, topBlockIsAir bool) OnActivateResult {
	if isBoneMeal {
		// 催草由 Level 层执行（TallGrassObject::growGrass）
		return OnActivateResult{Handled: true, UseTool: true}
	}
	if isHoe {
		return OnActivateResult{Handled: true, ReplaceBlock: FARMLAND, UseTool: true}
	}
	if isShovel && topBlockIsAir {
		return OnActivateResult{Handled: true, ReplaceBlock: GRASS_PATH, UseTool: true}
	}
	return OnActivateResult{Handled: false}
}

// ---------- 草方块随机 Tick ----------

// GrassRandomTickResult 随机 tick 结果
type GrassRandomTickResult uint8

const (
	GrassTickNoChange GrassRandomTickResult = iota
	GrassTickDie                            // 草死亡 → 泥土
	GrassTickSpread                         // 草蔓延（尝试将附近泥土变成草）
)

// CheckGrassRandomTick 检查草方块随机 tick 应该做什么
// 对应 PHP Grass::onUpdate(BLOCK_UPDATE_RANDOM)
//
// 参数:
//   lightAbove: 上方完整光照等级
//   lightFilterAbove: 上方方块的光过滤值
//
// 规则:
//   - 如果上方光照 < 4 且上方光过滤 >= 3 → 草死亡（变泥土）
//   - 如果上方光照 >= 9 → 尝试蔓延到附近泥土
func CheckGrassRandomTick(lightAbove int, lightFilterAbove int) GrassRandomTickResult {
	if lightAbove < 4 && lightFilterAbove >= 3 {
		return GrassTickDie
	}
	if lightAbove >= 9 {
		return GrassTickSpread
	}
	return GrassTickNoChange
}

// CanSpreadTo 判断草是否可以蔓延到指定方块
// 对应 PHP Grass::onUpdate() 中的蔓延条件检查
//
// 条件: 目标是泥土(非粗泥) + 上方光照 >= 4 + 上方光过滤 < 3 + 上方是空气
func CanGrassSpreadTo(targetID, targetMeta uint8, lightAboveTarget int, filterAboveTarget int, aboveTargetIsAir bool) bool {
	return targetID == DIRT &&
		targetMeta != 1 && // meta 1 = 粗泥，不能长草
		lightAboveTarget >= 4 &&
		filterAboveTarget < 3 &&
		aboveTargetIsAir
}

// ============ Sand 沙子 ============

// SandBlock 沙子方块（受重力影响）
// 对应 PHP class Sand extends Fallable
type SandBlock struct {
	FallableBase
}

func NewSandBlock() *SandBlock {
	return &SandBlock{
		FallableBase: FallableBase{
			SolidBase: SolidBase{
				BlockID:       SAND,
				BlockName:     "Sand",
				BlockHardness: 0.5,
				BlockToolType: ToolTypeShovel,
			},
		},
	}
}

// ============ Gravel 砂砾 ============

// GravelBlock 砂砾方块（受重力影响，10%掉落燧石）
// 对应 PHP class Gravel extends Fallable
type GravelBlock struct {
	FallableBase
}

func NewGravelBlock() *GravelBlock {
	return &GravelBlock{
		FallableBase: FallableBase{
			SolidBase: SolidBase{
				BlockID:       GRAVEL,
				BlockName:     "Gravel",
				BlockHardness: 0.6,
				BlockToolType: ToolTypeShovel,
			},
		},
	}
}

// ItemFlint 燧石物品ID
const ItemFlint = 318

// GetDrops 砂砾 10% 概率掉落燧石（精准采集总是掉落自身，由 Level 层处理）
// 对应 PHP Gravel::getDrops()
// 注意: 随机判断由 Level 层执行，此处提供两种掉落方法
func (b *GravelBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(GRAVEL), Meta: 0, Count: 1}}
}

// GetGravelFlintDrop 获取砂砾的燧石掉落（10%概率）
func GetGravelFlintDrop() Drop {
	return Drop{ID: ItemFlint, Meta: 0, Count: 1}
}

// GravelFlintChance 砂砾掉落燧石的概率（1/10）
const GravelFlintChance = 10

// ============ 注册 ============

func init() {
	Registry.Register(NewDirtBlock())
	Registry.Register(NewGrassBlock())
	Registry.Register(NewSandBlock())
	Registry.Register(NewGravelBlock())
}
