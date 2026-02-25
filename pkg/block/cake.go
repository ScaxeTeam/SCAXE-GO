package block

// CakeBlock 蛋糕方块
// 对应 PHP class Cake extends Transparent
//
// meta 0-6 表示已被吃掉的片数 (0=完整, 6=只剩一片)
// 每次右键吃一片，恢复 2 饥饿值

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
	CakeMaxSlices = 6 // meta 0 = 完整 (7片), meta 6 = 1片剩余
)

func (b *CakeBlock) CanBeActivated() bool {
	return true
}

// OnActivate 右键蛋糕 — 吃一片
func (b *CakeBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	// meta++ 表示吃掉一片; 如果 meta >= 6 则破坏方块
	return true
}

// GetDrops 蛋糕破坏不掉落
func (b *CakeBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

func init() {
	Registry.Register(NewCakeBlock())
}
