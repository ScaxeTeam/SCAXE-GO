package block

// DispenserBlock 发射器方块
// 对应 PHP class Dispenser extends Solid
//
// meta 0-5 表示朝向 (0=下 1=上 2=北 3=南 4=西 5=东)
// 红石信号激活时发射物品

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

// OnActivate 右键打开发射器背包
func (b *DispenserBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

// GetPlacementMeta 发射器根据玩家朝向放置
var DispenserDirectionToMeta = [4]uint8{3, 4, 2, 5} // 南西北东

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
