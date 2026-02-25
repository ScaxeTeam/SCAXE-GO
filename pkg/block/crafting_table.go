package block

// CraftingTableBlock 工作台方块
// 对应 PHP class Workbench extends Solid
//
// 右键激活时将玩家的合成模式切换为 3x3 大合成（CRAFTING_BIG）
// 无 TileEntity，不需要朝向
type CraftingTableBlock struct {
	SolidBase
}

// NewCraftingTableBlock 创建工作台方块
func NewCraftingTableBlock() *CraftingTableBlock {
	return &CraftingTableBlock{
		SolidBase: SolidBase{
			BlockID:       WORKBENCH,
			BlockName:     "Crafting Table",
			BlockHardness: 2.5,
			BlockToolType: ToolTypeAxe,
		},
	}
}

// CanBeActivated 工作台可以被右键激活（打开 3x3 合成界面）
// 对应 PHP Workbench::canBeActivated() { return true; }
func (b *CraftingTableBlock) CanBeActivated() bool {
	return true
}

// OnActivate 工作台右键交互 — 打开 3x3 合成界面
// 实际背包创建由服务器层通过 CraftingTableOnActivate() 执行
func (b *CraftingTableBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

// GetFuelTime 工作台可以作燃料（300 tick = 15秒）
// 对应 PHP Workbench::getFuelTime() { return 300; }
func (b *CraftingTableBlock) GetFuelTime() int {
	return 300
}

// GetDrops 工作台破坏后掉落自身
// 对应 PHP Workbench::getDrops()
func (b *CraftingTableBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(WORKBENCH), Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewCraftingTableBlock())
}
