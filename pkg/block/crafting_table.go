package block

type CraftingTableBlock struct {
	SolidBase
}

func NewCraftingTableBlock() *CraftingTableBlock {
	return &CraftingTableBlock{
		SolidBase: SolidBase{
			BlockID:       WORKBENCH,
			BlockName:     "Crafting Table",
			BlockHardness: 2.5,
			BlockToolType: ToolTypeAxe,
		},
	}
}

func (b *CraftingTableBlock) CanBeActivated() bool {
	return true
}

func (b *CraftingTableBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

func (b *CraftingTableBlock) GetFuelTime() int {
	return 300
}

func (b *CraftingTableBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(WORKBENCH), Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewCraftingTableBlock())
}
