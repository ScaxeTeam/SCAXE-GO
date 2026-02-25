package block

// BrewingStandBlockType 酿造台方块
// 对应 PHP class BrewingStand extends Transparent
//
// meta 低 3 位表示药水瓶槽位状态 (bit mask: bit0=东, bit1=南西, bit2=北西)

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

// OnActivate 右键打开酿造台背包
func (b *BrewingStandBlockType) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

func (b *BrewingStandBlockType) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: 379, Meta: 0, Count: 1}} // BREWING_STAND item
}

func init() {
	Registry.Register(NewBrewingStandBlockType())
}
