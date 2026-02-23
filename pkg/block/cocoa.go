package block

type CocoaBlock struct {
	TransparentBase
}

func NewCocoaBlock() *CocoaBlock {
	return &CocoaBlock{
		TransparentBase: TransparentBase{
			BlockID:       COCOA_BLOCK,
			BlockName:     "Cocoa Block",
			BlockHardness: 0.2,
			BlockToolType: ToolTypeNone,
		},
	}
}

func CocoaGetDirection(meta uint8) int {
	return int(meta & 0x03)
}

func CocoaGetAge(meta uint8) int {
	return int((meta >> 2) & 0x03)
}

func (b *CocoaBlock) GetDrops(toolType, toolTier int) []Drop {

	return []Drop{{ID: 351, Meta: 3, Count: 1}}
}

func init() {
	Registry.Register(NewCocoaBlock())
}
