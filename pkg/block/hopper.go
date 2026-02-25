package block

// HopperBlock 漏斗方块
// 对应 PHP class Hopper extends Transparent
//
// meta 低 3 位 = 输出朝向 (0=下, 2=北, 3=南, 4=西, 5=东; 注意: 1=下 也有效)
// meta bit 3 (0x08) = 是否被红石信号禁用

type HopperBlock struct {
	TransparentBase
}

func NewHopperBlock() *HopperBlock {
	return &HopperBlock{
		TransparentBase: TransparentBase{
			BlockID:       HOPPER_BLOCK,
			BlockName:     "Hopper",
			BlockHardness: 3,
			BlockToolType: ToolTypePickaxe,
		},
	}
}

func (b *HopperBlock) CanBeActivated() bool {
	return true
}

// OnActivate 右键打开漏斗背包
func (b *HopperBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

// GetPlacementMeta 漏斗放置 — 根据放置面决定输出方向
// 默认朝下(0)，放置在侧面则指向对应方向
func (b *HopperBlock) GetPlacementMeta(playerDirection int) uint8 {
	// 水平: 0南→3, 1西→4, 2北→2, 3东→5
	switch playerDirection {
	case 0:
		return 3
	case 1:
		return 4
	case 2:
		return 2
	case 3:
		return 5
	default:
		return 0 // 朝下
	}
}

// HopperGetFacing 从 meta 获取输出朝向
func HopperGetFacing(meta uint8) int {
	return int(meta & 0x07)
}

// HopperIsDisabled 漏斗是否被红石禁用
func HopperIsDisabled(meta uint8) bool {
	return meta&0x08 != 0
}

func (b *HopperBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(HOPPER_BLOCK), Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewHopperBlock())
}
