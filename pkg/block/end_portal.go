package block

type EndPortalBlock struct {
	TransparentBase
}

func NewEndPortalBlock() *EndPortalBlock {
	return &EndPortalBlock{
		TransparentBase: TransparentBase{
			BlockID:         END_PORTAL,
			BlockName:       "End Portal",
			BlockHardness:   -1,
			BlockLightLevel: 15,
			BlockToolType:   ToolTypeNone,
			BlockCanPlace:   false,
		},
	}
}

func (b *EndPortalBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

type EndPortalFrameBlock struct {
	SolidBase
}

func NewEndPortalFrameBlock() *EndPortalFrameBlock {
	return &EndPortalFrameBlock{
		SolidBase: SolidBase{
			BlockID:         END_PORTAL_FRAME,
			BlockName:       "End Portal Frame",
			BlockHardness:   -1,
			BlockLightLevel: 1,
			BlockToolType:   ToolTypeNone,
		},
	}
}

func (b *EndPortalFrameBlock) CanBeActivated() bool {
	return true
}
func (b *EndPortalFrameBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}
func EndPortalFrameHasEye(meta uint8) bool {
	return meta&0x04 != 0
}
func EndPortalFrameGetDirection(meta uint8) int {
	return int(meta & 0x03)
}

func (b *EndPortalFrameBlock) GetPlacementMeta(playerDirection int) uint8 {
	if playerDirection < 0 || playerDirection > 3 {
		return 0
	}
	return uint8(playerDirection)
}

func (b *EndPortalFrameBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

func init() {
	Registry.Register(NewEndPortalBlock())
	Registry.Register(NewEndPortalFrameBlock())
}
