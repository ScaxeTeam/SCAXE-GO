package block

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
func NetherWartGetAge(meta uint8) int {
	return int(meta & 0x03)
}
func NetherWartIsMature(meta uint8) bool {
	return meta&0x03 >= NetherWartMaxAge
}

func (b *NetherWartBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 372, Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewNetherWartBlock())
}
