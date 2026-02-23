package block

type TrapdoorBlock struct {
	TransparentBase
}

const (
	TrapdoorMaskDirection = 0x03
	TrapdoorMaskUpper     = 0x04
	TrapdoorMaskOpened    = 0x08
)

func NewTrapdoorBlock() *TrapdoorBlock {
	return &TrapdoorBlock{
		TransparentBase: TransparentBase{
			BlockID:         TRAPDOOR,
			BlockName:       "Wooden Trapdoor",
			BlockHardness:   3,
			BlockResistance: 15,
			BlockToolType:   ToolTypeAxe,
			BlockCanPlace:   true,
		},
	}
}

func NewIronTrapdoorBlock() *TrapdoorBlock {
	return &TrapdoorBlock{
		TransparentBase: TransparentBase{
			BlockID:         IRON_TRAPDOOR,
			BlockName:       "Iron Trapdoor",
			BlockHardness:   5,
			BlockResistance: 25,
			BlockToolType:   ToolTypePickaxe,
			BlockCanPlace:   true,
		},
	}
}

func (b *TrapdoorBlock) CanBeActivated() bool {
	return true
}

func (b *TrapdoorBlock) OnActivate(ctx *BlockContext, playerID int64) bool {

	if b.BlockID == IRON_TRAPDOOR {
		return false
	}
	return true
}

func (b *TrapdoorBlock) IsSolid() bool {
	return false
}

func (b *TrapdoorBlock) GetFuelTime() int {
	if b.BlockID == TRAPDOOR {
		return 300
	}
	return 0
}

func (b *TrapdoorBlock) GetDrops(toolType, toolTier int) []Drop {
	if b.BlockID == IRON_TRAPDOOR && toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(b.BlockID), Meta: 0, Count: 1}}
}

func TrapdoorIsOpen(meta uint8) bool {
	return meta&TrapdoorMaskOpened != 0
}

func TrapdoorIsUpper(meta uint8) bool {
	return meta&TrapdoorMaskUpper != 0
}

func TrapdoorGetDirection(meta uint8) uint8 {
	return meta & TrapdoorMaskDirection
}

func TrapdoorToggleOpen(meta uint8) uint8 {
	return meta ^ TrapdoorMaskOpened
}

var TrapdoorDirectionToMeta = [4]uint8{1, 3, 0, 2}

func GetTrapdoorPlacementMeta(playerDirection int, clickY float64, face int) uint8 {
	meta := TrapdoorDirectionToMeta[playerDirection&0x03]
	if (clickY > 0.5 && face != 1) || face == 0 {
		meta |= TrapdoorMaskUpper
	}
	return meta
}

type TrapdoorBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}

func GetTrapdoorBoundingBox(x, y, z int, meta uint8) TrapdoorBoundingBox {
	fx, fy, fz := float64(x), float64(y), float64(z)
	f := 0.1875

	isOpen := TrapdoorIsOpen(meta)
	isUpper := TrapdoorIsUpper(meta)
	dir := TrapdoorGetDirection(meta)

	if !isOpen {
		if isUpper {
			return TrapdoorBoundingBox{fx, fy + 1 - f, fz, fx + 1, fy + 1, fz + 1}
		}
		return TrapdoorBoundingBox{fx, fy, fz, fx + 1, fy + f, fz + 1}
	}

	switch dir {
	case 0:
		return TrapdoorBoundingBox{fx, fy, fz + 1 - f, fx + 1, fy + 1, fz + 1}
	case 1:
		return TrapdoorBoundingBox{fx, fy, fz, fx + 1, fy + 1, fz + f}
	case 2:
		return TrapdoorBoundingBox{fx + 1 - f, fy, fz, fx + 1, fy + 1, fz + 1}
	case 3:
		return TrapdoorBoundingBox{fx, fy, fz, fx + f, fy + 1, fz + 1}
	}

	return TrapdoorBoundingBox{fx, fy, fz, fx + 1, fy + 1, fz + 1}
}

func init() {
	Registry.Register(NewTrapdoorBlock())
	Registry.Register(NewIronTrapdoorBlock())
}
