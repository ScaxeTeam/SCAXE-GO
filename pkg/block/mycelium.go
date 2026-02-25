package block

// MyceliumBlock 菌丝方块
// MCPE 方块 ID 110
//
// 类似草方块，可以向周围泥土蔓延
// 在菌丝上可以放置蘑菇（不受光照限制）

type MyceliumBlock struct {
	SolidBase
}

func NewMyceliumBlock() *MyceliumBlock {
	return &MyceliumBlock{
		SolidBase: SolidBase{
			BlockID:       MYCELIUM,
			BlockName:     "Mycelium",
			BlockHardness: 0.6,
			BlockToolType: ToolTypeShovel,
		},
	}
}

// GetDrops 菌丝需要精准采集否则掉落泥土
func (b *MyceliumBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(DIRT), Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewMyceliumBlock())
}
