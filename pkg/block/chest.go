package block
type ChestBlock struct {
	TransparentBase
}
func NewChestBlock() *ChestBlock {
	return &ChestBlock{
		TransparentBase: TransparentBase{
			BlockID:       CHEST,
			BlockName:     "Chest",
			BlockHardness: 2.5,
			BlockToolType: ToolTypeAxe,
			BlockCanPlace: true,
		},
	}
}
func (b *ChestBlock) CanBeActivated() bool {
	return true
}
func (b *ChestBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}
func (b *ChestBlock) GetFuelTime() int {
	return 300
}
func (b *ChestBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(CHEST), Meta: 0, Count: 1}}
}
var ChestDirectionToMeta = [4]uint8{4, 2, 5, 3}
func (b *ChestBlock) GetPlacementMeta(playerDirection int) uint8 {
	if playerDirection < 0 || playerDirection > 3 {
		playerDirection = 0
	}
	return ChestDirectionToMeta[playerDirection]
}
func GetPairSearchSides(meta uint8) []int {
	switch meta {
	case 4, 5:
		return []int{2, 3}
	case 2, 3:
		return []int{4, 5}
	default:
		return []int{2, 3, 4, 5}
	}
}
type ChestBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}
func GetChestBoundingBox(x, y, z int) ChestBoundingBox {
	return ChestBoundingBox{
		MinX: float64(x) + 0.0625,
		MinY: float64(y),
		MinZ: float64(z) + 0.0625,
		MaxX: float64(x) + 0.9375,
		MaxY: float64(y) + 0.9475,
		MaxZ: float64(z) + 0.9375,
	}
}

func init() {
	Registry.Register(NewChestBlock())
}
