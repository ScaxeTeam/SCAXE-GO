package block

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
func (b *HopperBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}
func (b *HopperBlock) GetPlacementMeta(playerDirection int) uint8 {
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
		return 0
	}
}
func HopperGetFacing(meta uint8) int {
	return int(meta & 0x07)
}
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
