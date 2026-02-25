package block

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
	if b.BlockID == BROWN_MUSHROOM_BLOCK {
		return []Drop{{ID: int(BROWN_MUSHROOM), Meta: 0, Count: 1}}
	}
	return []Drop{{ID: int(RED_MUSHROOM), Meta: 0, Count: 1}}
}
func (b *HugeMushroomBlock) GetFuelTime() int {
	return 300
}

func init() {
	Registry.Register(NewBrownMushroomBlock())
	Registry.Register(NewRedMushroomBlock())
}
