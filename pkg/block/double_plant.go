package block

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
func DoublePlantGetType(meta uint8) int {
	return int(meta & 0x07)
}
func DoublePlantIsTop(meta uint8) bool {
	return meta&0x08 != 0
}

func (b *DoublePlantBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

func init() {
	Registry.Register(NewDoublePlantBlock())
}
