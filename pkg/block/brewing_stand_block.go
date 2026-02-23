package block

type BrewingStandBlockType struct {
	TransparentBase
}

func NewBrewingStandBlockType() *BrewingStandBlockType {
	return &BrewingStandBlockType{
		TransparentBase: TransparentBase{
			BlockID:         BREWING_STAND_BLOCK,
			BlockName:       "Brewing Stand",
			BlockHardness:   0.5,
			BlockLightLevel: 1,
			BlockToolType:   ToolTypePickaxe,
		},
	}
}

func (b *BrewingStandBlockType) CanBeActivated() bool {
	return true
}

func (b *BrewingStandBlockType) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

func (b *BrewingStandBlockType) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: 379, Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewBrewingStandBlockType())
}
