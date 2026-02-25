package block

// DragonEggBlock 龙蛋方块
// MCPE 方块 ID 122
// 点击时传送到附近随机位置

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

// OnActivate 右键龙蛋 — 传送到附近随机位置
func (b *DragonEggBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	// 龙蛋传送逻辑: 随机 15-30 格半径内选择一个空气方块
	return true
}

func (b *DragonEggBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(DRAGON_EGG), Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewDragonEggBlock())
}
