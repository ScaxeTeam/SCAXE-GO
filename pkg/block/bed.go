package block

// BedBlock 床方块
// 对应 PHP class Bed extends Transparent
//
// meta 低 3 位 = 朝向 (0:南 1:西 2:北 3:东)
// meta bit 3 (0x08) = 是否为床头部分 (1=头, 0=脚)

type BedBlock struct {
	TransparentBase
}

func NewBedBlock() *BedBlock {
	return &BedBlock{
		TransparentBase: TransparentBase{
			BlockID:       BED_BLOCK,
			BlockName:     "Bed Block",
			BlockHardness: 0.2,
			BlockToolType: ToolTypeNone,
			BlockCanPlace: false, // 通过物品放置
		},
	}
}

func (b *BedBlock) CanBeActivated() bool {
	return true
}

// OnActivate 右键床 — 设置重生点 / 入睡
func (b *BedBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

// IsHeadPart 判断当前 meta 是否为床头
func IsHeadPart(meta uint8) bool {
	return meta&0x08 != 0
}

// GetBedDirection 从 meta 获取床朝向 (0-3)
func GetBedDirection(meta uint8) int {
	return int(meta & 0x03)
}

func (b *BedBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 355, Meta: 0, Count: 1}} // BED item
}

func init() {
	Registry.Register(NewBedBlock())
}
