package block

// DoublePlantBlock 两格高植物
// MCPE 方块 ID 175
//
// meta 低 3 位 = 植物类型:
//   0 = 向日葵  1 = 丁香花  2 = 高草丛
//   3 = 大型蕨  4 = 玫瑰丛  5 = 牡丹花
// meta bit 3 (0x08) = 是否为上半部分

const (
	DoublePlantSunflower = 0
	DoublePlantLilac     = 1
	DoublePlantTallGrass = 2
	DoublePlantLargeFern = 3
	DoublePlantRoseBush  = 4
	DoublePlantPeony     = 5
)

type DoublePlantBlock struct {
	TransparentBase
}

func NewDoublePlantBlock() *DoublePlantBlock {
	return &DoublePlantBlock{
		TransparentBase: TransparentBase{
			BlockID:       DOUBLE_PLANT,
			BlockName:     "Double Plant",
			BlockHardness: 0,
			BlockToolType: ToolTypeNone,
		},
	}
}

// DoublePlantGetType 获取植物子类型 (0-5)
func DoublePlantGetType(meta uint8) int {
	return int(meta & 0x07)
}

// DoublePlantIsTop 是否为上半部分
func DoublePlantIsTop(meta uint8) bool {
	return meta&0x08 != 0
}

func (b *DoublePlantBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil // 大部分双格植物不掉落；向日葵/花丛由具体逻辑决定
}

func init() {
	Registry.Register(NewDoublePlantBlock())
}
