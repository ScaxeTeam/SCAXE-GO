package block

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
			BlockCanPlace: false,
		},
	}
}

func (b *BedBlock) CanBeActivated() bool {
	return true
}
func (b *BedBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}
func IsHeadPart(meta uint8) bool {
	return meta&0x08 != 0
}
func GetBedDirection(meta uint8) int {
	return int(meta & 0x03)
}

func (b *BedBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 355, Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewBedBlock())
}
