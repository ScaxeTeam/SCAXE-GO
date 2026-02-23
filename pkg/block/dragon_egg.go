package block

type DragonEggBlock struct {
	TransparentBase
}

func NewDragonEggBlock() *DragonEggBlock {
	return &DragonEggBlock{
		TransparentBase: TransparentBase{
			BlockID:         DRAGON_EGG,
			BlockName:       "Dragon Egg",
			BlockHardness:   3,
			BlockLightLevel: 1,
			BlockToolType:   ToolTypeNone,
		},
	}
}

func (b *DragonEggBlock) CanBeActivated() bool {
	return true
}

func (b *DragonEggBlock) OnActivate(ctx *BlockContext, playerID int64) bool {

	return true
}

func (b *DragonEggBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(DRAGON_EGG), Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewDragonEggBlock())
}
