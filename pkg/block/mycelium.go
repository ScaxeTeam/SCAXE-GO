package block

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
func (b *MyceliumBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(DIRT), Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewMyceliumBlock())
}
