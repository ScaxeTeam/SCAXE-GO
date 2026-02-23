package block

type FenceGateBlock struct {
	TransparentBase
}

const (
	FenceGateMaskDirection = 0x03
	FenceGateMaskOpen      = 0x04
)

func newFenceGate(blockID uint8, name string) *FenceGateBlock {
	return &FenceGateBlock{
		TransparentBase: TransparentBase{
			BlockID:         blockID,
			BlockName:       name,
			BlockHardness:   2,
			BlockResistance: 10,
			BlockToolType:   ToolTypeAxe,
			BlockCanPlace:   true,
		},
	}
}

func (b *FenceGateBlock) CanBeActivated() bool {
	return true
}

func (b *FenceGateBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

func (b *FenceGateBlock) GetFuelTime() int {
	return 300
}

func (b *FenceGateBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(b.BlockID), Meta: 0, Count: 1}}
}

func FenceGateIsOpen(meta uint8) bool {
	return meta&FenceGateMaskOpen != 0
}

func FenceGateGetDirection(meta uint8) uint8 {
	return meta & FenceGateMaskDirection
}

func FenceGateToggleOpen(meta uint8) uint8 {
	return meta ^ FenceGateMaskOpen
}

var FenceGateDirectionToMeta = [4]uint8{3, 0, 1, 2}

func GetFenceGatePlacementMeta(playerDirection int) uint8 {
	return FenceGateDirectionToMeta[playerDirection&0x03]
}

type FenceGateBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
	HasCollision     bool
}

func GetFenceGateBoundingBox(x, y, z int, meta uint8) FenceGateBoundingBox {
	fx, fy, fz := float64(x), float64(y), float64(z)

	if FenceGateIsOpen(meta) {
		return FenceGateBoundingBox{HasCollision: false}
	}

	dir := FenceGateGetDirection(meta)

	if dir == 0 || dir == 2 {
		return FenceGateBoundingBox{
			MinX: fx, MinY: fy, MinZ: fz + 0.375,
			MaxX: fx + 1, MaxY: fy + 1.5, MaxZ: fz + 0.625,
			HasCollision: true,
		}
	}
	return FenceGateBoundingBox{
		MinX: fx + 0.375, MinY: fy, MinZ: fz,
		MaxX: fx + 0.625, MaxY: fy + 1.5, MaxZ: fz + 1,
		HasCollision: true,
	}
}

func init() {
	Registry.Register(newFenceGate(FENCE_GATE, "Oak Fence Gate"))
	Registry.Register(newFenceGate(FENCE_GATE_SPRUCE, "Spruce Fence Gate"))
	Registry.Register(newFenceGate(FENCE_GATE_BIRCH, "Birch Fence Gate"))
	Registry.Register(newFenceGate(FENCE_GATE_JUNGLE, "Jungle Fence Gate"))
	Registry.Register(newFenceGate(FENCE_GATE_DARK_OAK, "Dark Oak Fence Gate"))
	Registry.Register(newFenceGate(FENCE_GATE_ACACIA, "Acacia Fence Gate"))
}
