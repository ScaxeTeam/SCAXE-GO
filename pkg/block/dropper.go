package block

// DropperBlock 投掷器方块
// 对应 PHP class Dropper extends Solid
//
// meta 0-5 表示朝向 (与 Dispenser 相同)
// 类似发射器，但只会弹出物品实体而非发射投射物

type DropperBlock struct {
	SolidBase
}

func NewDropperBlock() *DropperBlock {
	return &DropperBlock{
		SolidBase: SolidBase{
			BlockID:       DROPPER,
			BlockName:     "Dropper",
			BlockHardness: 3.5,
			BlockToolType: ToolTypePickaxe,
		},
	}
}

func (b *DropperBlock) CanBeActivated() bool {
	return true
}

// OnActivate 右键打开投掷器背包
func (b *DropperBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

// GetPlacementMeta 与发射器相同的朝向映射
func (b *DropperBlock) GetPlacementMeta(playerDirection int) uint8 {
	if playerDirection < 0 || playerDirection > 3 {
		playerDirection = 0
	}
	return DispenserDirectionToMeta[playerDirection]
}

func (b *DropperBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(DROPPER), Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewDropperBlock())
}
