package block

// EndPortalBlock 末地传送门方块
// MCPE 方块 ID 119
// 不可破坏，不可放置，由末影之眼填满传送门框架后生成

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

// EndPortalFrameBlock 末地传送门框架
// MCPE 方块 ID 120
// meta bit 2 (0x04) = 是否已放入末影之眼
// meta 低 2 位 = 朝向 (0=南 1=西 2=北 3=东)

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

// OnActivate 右键放入末影之眼
func (b *EndPortalFrameBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

// EndPortalFrameHasEye 是否已有末影之眼
func EndPortalFrameHasEye(meta uint8) bool {
	return meta&0x04 != 0
}

// EndPortalFrameGetDirection 获取朝向
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
	return nil // 不可破坏/不掉落
}

func init() {
	Registry.Register(NewEndPortalBlock())
	Registry.Register(NewEndPortalFrameBlock())
}
