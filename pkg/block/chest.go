package block

// ChestBlock 箱子方块（方块层面的放置/破坏/交互逻辑）
// 对应 PHP class Chest extends Transparent
//
// 关键行为:
//   - 放置时根据玩家朝向设置 meta（朝向）
//   - 放置时搜索相邻同向箱子进行大箱子配对
//   - 破坏时解除配对
//   - 右键激活时打开箱子背包（检查上方方块是否透明）
//   - 碰撞箱比完整方块略小（每侧内缩 0.0625）
type ChestBlock struct {
	TransparentBase
}

// NewChestBlock 创建箱子方块行为
func NewChestBlock() *ChestBlock {
	return &ChestBlock{
		TransparentBase: TransparentBase{
			BlockID:       CHEST,
			BlockName:     "Chest",
			BlockHardness: 2.5,
			BlockToolType: ToolTypeAxe,
			BlockCanPlace: true,
		},
	}
}

// CanBeActivated 箱子可以被右键激活（打开背包）
// 对应 PHP Chest::canBeActivated() { return true; }
func (b *ChestBlock) CanBeActivated() bool {
	return true
}

// OnActivate 箱子右键交互 — 打开背包
// 实际背包创建由服务器层通过 ChestOnActivate() 执行
func (b *ChestBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

// GetFuelTime 箱子可以作燃料（300 tick = 15秒）
// 对应 PHP Chest::getFuelTime() { return 300; }
func (b *ChestBlock) GetFuelTime() int {
	return 300
}

// GetDrops 箱子破坏后掉落自身
// 对应 PHP Chest::getDrops()
func (b *ChestBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(CHEST), Meta: 0, Count: 1}}
}

// ---------- 放置朝向 ----------

// DirectionToMeta 将玩家朝向（0-3）转换为箱子方块 meta
// 对应 PHP Chest::place() 中的 $faces 映射
//
//	方向 0 (南) → meta 4
//	方向 1 (西) → meta 2
//	方向 2 (北) → meta 5
//	方向 3 (东) → meta 3
var ChestDirectionToMeta = [4]uint8{4, 2, 5, 3}

// GetPlacementMeta 根据玩家朝向返回放置的 meta 值
func (b *ChestBlock) GetPlacementMeta(playerDirection int) uint8 {
	if playerDirection < 0 || playerDirection > 3 {
		playerDirection = 0
	}
	return ChestDirectionToMeta[playerDirection]
}

// ---------- 大箱子配对搜索 ----------

// FindPairableSides 返回可能配对的方向列表（排除与箱子朝向相同轴的方向）
// 对应 PHP Chest::place() 中 side 2~5 的过滤逻辑
//
//	meta 4 或 5 (东西朝向) → 只检查南(2)北(3)
//	meta 2 或 3 (南北朝向) → 只检查西(4)东(5)
func GetPairSearchSides(meta uint8) []int {
	switch meta {
	case 4, 5: // 东西朝向，只搜南北
		return []int{2, 3}
	case 2, 3: // 南北朝向，只搜东西
		return []int{4, 5}
	default:
		return []int{2, 3, 4, 5}
	}
}

// ---------- 碰撞箱 ----------

// BoundingBox 箱子的碰撞箱（比完整方块每侧内缩 1/16）
// 对应 PHP Chest::recalculateBoundingBox()
type ChestBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}

// GetChestBoundingBox 返回指定坐标的箱子碰撞箱
func GetChestBoundingBox(x, y, z int) ChestBoundingBox {
	return ChestBoundingBox{
		MinX: float64(x) + 0.0625,
		MinY: float64(y),
		MinZ: float64(z) + 0.0625,
		MaxX: float64(x) + 0.9375,
		MaxY: float64(y) + 0.9475, // 与 PHP 原始值一致
		MaxZ: float64(z) + 0.9375,
	}
}

func init() {
	Registry.Register(NewChestBlock())
}
