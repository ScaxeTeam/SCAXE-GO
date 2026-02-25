package block

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
func (b *DropperBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}
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
