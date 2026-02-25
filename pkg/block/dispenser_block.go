package block

type DispenserBlock struct {
	SolidBase
}

func NewDispenserBlock() *DispenserBlock {
	return &DispenserBlock{
		SolidBase: SolidBase{
			BlockID:       DISPENSER,
			BlockName:     "Dispenser",
			BlockHardness: 3.5,
			BlockToolType: ToolTypePickaxe,
		},
	}
}

func (b *DispenserBlock) CanBeActivated() bool {
	return true
}
func (b *DispenserBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}
var DispenserDirectionToMeta = [4]uint8{3, 4, 2, 5}

func (b *DispenserBlock) GetPlacementMeta(playerDirection int) uint8 {
	if playerDirection < 0 || playerDirection > 3 {
		playerDirection = 0
	}
	return DispenserDirectionToMeta[playerDirection]
}

func (b *DispenserBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(DISPENSER), Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewDispenserBlock())
}
