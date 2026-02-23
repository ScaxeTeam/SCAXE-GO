package block

type PistonBlock struct {
	SolidBase
}

func NewPistonBlock() *PistonBlock {
	return &PistonBlock{
		SolidBase: SolidBase{
			BlockID:       PISTON,
			BlockName:     "Piston",
			BlockHardness: 0.5,
			BlockToolType: ToolTypeNone,
		},
	}
}

func (b *PistonBlock) GetPlacementMeta(playerDirection int) uint8 {

	switch playerDirection {
	case 0:
		return 3
	case 1:
		return 4
	case 2:
		return 2
	case 3:
		return 5
	default:
		return 1
	}
}

func PistonGetFacing(meta uint8) int {
	return int(meta & 0x07)
}

func PistonIsExtended(meta uint8) bool {
	return meta&0x08 != 0
}

func (b *PistonBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(PISTON), Meta: 0, Count: 1}}
}

type StickyPistonBlock struct {
	SolidBase
}

func NewStickyPistonBlock() *StickyPistonBlock {
	return &StickyPistonBlock{
		SolidBase: SolidBase{
			BlockID:       STICKY_PISTON,
			BlockName:     "Sticky Piston",
			BlockHardness: 0.5,
			BlockToolType: ToolTypeNone,
		},
	}
}

func (b *StickyPistonBlock) GetPlacementMeta(playerDirection int) uint8 {
	switch playerDirection {
	case 0:
		return 3
	case 1:
		return 4
	case 2:
		return 2
	case 3:
		return 5
	default:
		return 1
	}
}

func (b *StickyPistonBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(STICKY_PISTON), Meta: 0, Count: 1}}
}

type PistonHeadBlock struct {
	TransparentBase
}

func NewPistonHeadBlock() *PistonHeadBlock {
	return &PistonHeadBlock{
		TransparentBase: TransparentBase{
			BlockID:       PISTON_HEAD,
			BlockName:     "Piston Head",
			BlockHardness: 0.5,
			BlockToolType: ToolTypeNone,
			BlockCanPlace: false,
		},
	}
}

func PistonHeadIsSticky(meta uint8) bool {
	return meta&0x08 != 0
}

func (b *PistonHeadBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

const (
	PistonFacingDown  = 0
	PistonFacingUp    = 1
	PistonFacingNorth = 2
	PistonFacingSouth = 3
	PistonFacingWest  = 4
	PistonFacingEast  = 5

	PistonMaxPushDistance = 12
)

func PistonFacingOffset(facing int) (dx, dy, dz int) {
	switch facing {
	case PistonFacingDown:
		return 0, -1, 0
	case PistonFacingUp:
		return 0, 1, 0
	case PistonFacingNorth:
		return 0, 0, -1
	case PistonFacingSouth:
		return 0, 0, 1
	case PistonFacingWest:
		return -1, 0, 0
	case PistonFacingEast:
		return 1, 0, 0
	default:
		return 0, 0, 0
	}
}

func init() {
	Registry.Register(NewPistonBlock())
	Registry.Register(NewStickyPistonBlock())
	Registry.Register(NewPistonHeadBlock())
}
