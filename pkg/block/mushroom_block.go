package block

// HugeMushroomBlock 巨型蘑菇方块
// MCPE 方块 ID 99 (棕色) / 100 (红色)
//
// meta 0-15 控制蘑菇方块纹理:
//   0 = 全孔隙面  10 = 全蘑菇面 (柄)  14 = 全蘑菇顶面
//   1-9 = 各角/边/顶面组合  15 = 全蘑菇面

type HugeMushroomBlock struct {
	SolidBase
}

func NewBrownMushroomBlock() *HugeMushroomBlock {
	return &HugeMushroomBlock{
		SolidBase: SolidBase{
			BlockID:       BROWN_MUSHROOM_BLOCK,
			BlockName:     "Brown Mushroom Block",
			BlockHardness: 0.2,
			BlockToolType: ToolTypeAxe,
		},
	}
}

func NewRedMushroomBlock() *HugeMushroomBlock {
	return &HugeMushroomBlock{
		SolidBase: SolidBase{
			BlockID:       RED_MUSHROOM_BLOCK,
			BlockName:     "Red Mushroom Block",
			BlockHardness: 0.2,
			BlockToolType: ToolTypeAxe,
		},
	}
}

func (b *HugeMushroomBlock) GetDrops(toolType, toolTier int) []Drop {
	// 掉落 0-2 个对应小蘑菇
	if b.BlockID == BROWN_MUSHROOM_BLOCK {
		return []Drop{{ID: int(BROWN_MUSHROOM), Meta: 0, Count: 1}}
	}
	return []Drop{{ID: int(RED_MUSHROOM), Meta: 0, Count: 1}}
}

// GetFuelTime 蘑菇方块可作燃料 (300 tick)
func (b *HugeMushroomBlock) GetFuelTime() int {
	return 300
}

func init() {
	Registry.Register(NewBrownMushroomBlock())
	Registry.Register(NewRedMushroomBlock())
}
