package block

type CakeBlock struct {
	TransparentBase
}

func NewCakeBlock() *CakeBlock {
	return &CakeBlock{
		TransparentBase: TransparentBase{
			BlockID:       CAKE_BLOCK,
			BlockName:     "Cake Block",
			BlockHardness: 0.5,
			BlockToolType: ToolTypeNone,
			BlockCanPlace: false,
		},
	}
}

const (
	CakeMaxSlices = 6
)

func (b *CakeBlock) CanBeActivated() bool {
	return true
}
func (b *CakeBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}
func (b *CakeBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

func init() {
	Registry.Register(NewCakeBlock())
}
