package block
type FurnaceBlock struct {
	SolidBase
}
func NewFurnaceBlock() *FurnaceBlock {
	return &FurnaceBlock{
		SolidBase: SolidBase{
			BlockID:         FURNACE,
			BlockName:       "Furnace",
			BlockHardness:   3.5,
			BlockLightLevel: 0,
			BlockToolType:   ToolTypePickaxe,
		},
	}
}
func NewBurningFurnaceBlock() *FurnaceBlock {
	return &FurnaceBlock{
		SolidBase: SolidBase{
			BlockID:         BURNING_FURNACE,
			BlockName:       "Burning Furnace",
			BlockHardness:   3.5,
			BlockLightLevel: 13,
			BlockToolType:   ToolTypePickaxe,
		},
	}
}
func (b *FurnaceBlock) CanBeActivated() bool {
	return true
}
func (b *FurnaceBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}
func (b *FurnaceBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(FURNACE), Meta: 0, Count: 1}}
}
var FurnaceDirectionToMeta = [4]uint8{4, 2, 5, 3}
func (b *FurnaceBlock) GetPlacementMeta(playerDirection int) uint8 {
	if playerDirection < 0 || playerDirection > 3 {
		playerDirection = 0
	}
	return FurnaceDirectionToMeta[playerDirection]
}

func init() {
	Registry.Register(NewFurnaceBlock())
	Registry.Register(NewBurningFurnaceBlock())
}
