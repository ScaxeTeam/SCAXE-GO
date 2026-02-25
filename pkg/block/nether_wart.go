package block

// NetherWartBlock 地狱疣方块 (种植)
// MCPE 方块 ID 115
//
// meta 0-3 = 生长阶段 (0=刚种下, 3=成熟)
// 只能种在灵魂沙上

type NetherWartBlock struct {
	TransparentBase
}

func NewNetherWartBlock() *NetherWartBlock {
	return &NetherWartBlock{
		TransparentBase: TransparentBase{
			BlockID:       NETHER_WART_BLOCK,
			BlockName:     "Nether Wart",
			BlockHardness: 0,
			BlockToolType: ToolTypeNone,
		},
	}
}

const NetherWartMaxAge = 3

// NetherWartGetAge 获取生长阶段 (0-3)
func NetherWartGetAge(meta uint8) int {
	return int(meta & 0x03)
}

// NetherWartIsMature 是否成熟
func NetherWartIsMature(meta uint8) bool {
	return meta&0x03 >= NetherWartMaxAge
}

func (b *NetherWartBlock) GetDrops(toolType, toolTier int) []Drop {
	// 成熟: 2-4个; 未成熟: 1个
	return []Drop{{ID: 372, Meta: 0, Count: 1}} // NETHER_WART item
}

func init() {
	Registry.Register(NewNetherWartBlock())
}
