package block

// CauldronBlock 炼药锅方块
// 对应 PHP class Cauldron extends Solid
//
// meta 0-3 表示水量 (0=空, 3=满)

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

// OnActivate 右键炼药锅 — 装水/取水/清洗皮革
func (b *CauldronBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

// GetDrops 需要镐
func (b *CauldronBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: 380, Meta: 0, Count: 1}} // CAULDRON item
}

func init() {
	Registry.Register(NewCauldronBlock())
}
