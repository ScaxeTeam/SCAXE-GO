package block

// FurnaceBlock 熔炉方块（方块层面的放置/破坏/交互逻辑）
// 对应 PHP class BurningFurnace extends Solid 和 class Furnace extends BurningFurnace
//
// Go 中统一为一个结构体，通过 BlockID 区分熄灭/燃烧状态：
//   FURNACE (61)         — 未燃烧的熔炉
//   BURNING_FURNACE (62) — 正在燃烧的熔炉（发光）
type FurnaceBlock struct {
	SolidBase
}

// NewFurnaceBlock 创建未燃烧的熔炉方块
// 对应 PHP class Furnace extends BurningFurnace
func NewFurnaceBlock() *FurnaceBlock {
	return &FurnaceBlock{
		SolidBase: SolidBase{
			BlockID:         FURNACE,
			BlockName:       "Furnace",
			BlockHardness:   3.5,
			BlockLightLevel: 0,
			BlockToolType:   ToolTypePickaxe,
		},
	}
}

// NewBurningFurnaceBlock 创建正在燃烧的熔炉方块（发光）
// 对应 PHP class BurningFurnace extends Solid
func NewBurningFurnaceBlock() *FurnaceBlock {
	return &FurnaceBlock{
		SolidBase: SolidBase{
			BlockID:         BURNING_FURNACE,
			BlockName:       "Burning Furnace",
			BlockHardness:   3.5,
			BlockLightLevel: 13,
			BlockToolType:   ToolTypePickaxe,
		},
	}
}

// CanBeActivated 熔炉可以被右键激活（打开背包）
// 对应 PHP BurningFurnace::canBeActivated() { return true; }
func (b *FurnaceBlock) CanBeActivated() bool {
	return true
}

// OnActivate 熔炉右键交互 — 打开熔炉背包
// 实际背包创建由服务器层通过 FurnaceOnActivate() 执行
func (b *FurnaceBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

// GetDrops 熔炉需要镐才能掉落，且始终掉落未燃烧的熔炉
// 对应 PHP BurningFurnace::getDrops()
func (b *FurnaceBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	// 不论是否在燃烧，掉落物都是未燃烧的熔炉
	return []Drop{{ID: int(FURNACE), Meta: 0, Count: 1}}
}

// ---------- 放置朝向 ----------

// FurnaceDirectionToMeta 与箱子相同的朝向映射
// 对应 PHP BurningFurnace::place() 中的 $faces 映射
var FurnaceDirectionToMeta = [4]uint8{4, 2, 5, 3}

// GetPlacementMeta 根据玩家朝向返回放置的 meta 值
func (b *FurnaceBlock) GetPlacementMeta(playerDirection int) uint8 {
	if playerDirection < 0 || playerDirection > 3 {
		playerDirection = 0
	}
	return FurnaceDirectionToMeta[playerDirection]
}

// ---------- 注册 ----------

func init() {
	Registry.Register(NewFurnaceBlock())
	Registry.Register(NewBurningFurnaceBlock())
}
