package block

// CocoaBlock 可可豆方块
// MCPE 方块 ID 127
//
// meta 低 2 位 = 朝向 (0=南 1=西 2=北 3=东) — 附着的丛林木方向
// meta bit 2-3 = 生长阶段 (0=小, 1=中, 2=大/可采摘)

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

// CocoaGetDirection 获取朝向 (0-3)
func CocoaGetDirection(meta uint8) int {
	return int(meta & 0x03)
}

// CocoaGetAge 获取生长阶段 (0-2)
func CocoaGetAge(meta uint8) int {
	return int((meta >> 2) & 0x03)
}

func (b *CocoaBlock) GetDrops(toolType, toolTier int) []Drop {
	// 成熟时掉落 2-3 个可可豆; 未成熟掉落 1 个
	return []Drop{{ID: 351, Meta: 3, Count: 1}} // DYE:3 = cocoa beans
}

func init() {
	Registry.Register(NewCocoaBlock())
}
