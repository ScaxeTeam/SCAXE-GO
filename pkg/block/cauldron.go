package block

type CauldronBlock struct {
	SolidBase
}

func NewCauldronBlock() *CauldronBlock {
	return &CauldronBlock{
		SolidBase: SolidBase{
			BlockID:       CAULDRON_BLOCK,
			BlockName:     "Cauldron",
			BlockHardness: 2,
			BlockToolType: ToolTypePickaxe,
		},
	}
}

const (
	CauldronEmpty = 0
	CauldronFull  = 3
)

func (b *CauldronBlock) CanBeActivated() bool {
	return true
}

func (b *CauldronBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

func (b *CauldronBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: 380, Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewCauldronBlock())
}
