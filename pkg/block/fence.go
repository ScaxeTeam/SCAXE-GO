package block
type FenceBlock struct {
	TransparentBase
}

func NewFenceBlock() *FenceBlock {
	return &FenceBlock{
		TransparentBase: TransparentBase{
			BlockID:       FENCE,
			BlockName:     "Fence",
			BlockHardness: 2,
			BlockToolType: ToolTypeAxe,
			BlockCanPlace: true,
		},
	}
}
func (b *FenceBlock) IsSolid() bool { return true }
func (b *FenceBlock) GetFuelTime() int { return 300 }
func (b *FenceBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(FENCE), Meta: 0, Count: 1}}
}
var FenceVariantNames = [6]string{
	"Oak Fence", "Spruce Fence", "Birch Fence",
	"Jungle Fence", "Acacia Fence", "Dark Oak Fence",
}
type NetherBrickFenceBlock struct {
	TransparentBase
}

func NewNetherBrickFenceBlock() *NetherBrickFenceBlock {
	return &NetherBrickFenceBlock{
		TransparentBase: TransparentBase{
			BlockID:       NETHER_BRICK_FENCE,
			BlockName:     "Nether Brick Fence",
			BlockHardness: 2,
			BlockToolType: ToolTypePickaxe,
			BlockCanPlace: true,
		},
	}
}

func (b *NetherBrickFenceBlock) IsSolid() bool { return true }

func (b *NetherBrickFenceBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(NETHER_BRICK_FENCE), Meta: 0, Count: 1}}
}
type StoneWallBlock struct {
	TransparentBase
}

func NewStoneWallBlock() *StoneWallBlock {
	return &StoneWallBlock{
		TransparentBase: TransparentBase{
			BlockID:         STONE_WALL,
			BlockName:       "Cobblestone Wall",
			BlockHardness:   2,
			BlockResistance: 30,
			BlockToolType:   ToolTypePickaxe,
			BlockCanPlace:   true,
		},
	}
}

func (b *StoneWallBlock) IsSolid() bool { return false }

func (b *StoneWallBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(STONE_WALL), Meta: 0, Count: 1}}
}
type FenceConnections struct {
	North, South, West, East bool
}
func CanFenceConnect(targetID uint8, targetIsSolid, targetIsTransparent bool) bool {
	if targetID == FENCE || targetID == NETHER_BRICK_FENCE {
		return true
	}
	if IsFenceGate(targetID) {
		return true
	}
	return targetIsSolid && !targetIsTransparent
}
func CanWallConnect(targetID uint8, targetIsSolid, targetIsTransparent bool) bool {
	if targetID == STONE_WALL || IsFenceGate(targetID) {
		return true
	}
	return targetIsSolid && !targetIsTransparent
}
func IsFenceGate(blockID uint8) bool {
	switch blockID {
	case FENCE_GATE, FENCE_GATE_SPRUCE, FENCE_GATE_BIRCH,
		FENCE_GATE_JUNGLE, FENCE_GATE_DARK_OAK, FENCE_GATE_ACACIA:
		return true
	}
	return false
}

const (
	FenceThickness = 0.25
	FenceHeight    = 1.5
	WallThickness  = 0.5
)
type FenceBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}
func GetFenceBoundingBox(x, y, z int, conn FenceConnections) FenceBoundingBox {
	fx, fy, fz := float64(x), float64(y), float64(z)

	n := 0.375
	if conn.North {
		n = 0
	}
	s := 0.625
	if conn.South {
		s = 1
	}
	w := 0.375
	if conn.West {
		w = 0
	}
	e := 0.625
	if conn.East {
		e = 1
	}

	return FenceBoundingBox{fx + w, fy, fz + n, fx + e, fy + FenceHeight, fz + s}
}
func GetWallBoundingBox(x, y, z int, conn FenceConnections) FenceBoundingBox {
	fx, fy, fz := float64(x), float64(y), float64(z)

	n := 0.25
	if conn.North {
		n = 0
	}
	s := 0.75
	if conn.South {
		s = 1
	}
	w := 0.25
	if conn.West {
		w = 0
	}
	e := 0.75
	if conn.East {
		e = 1
	}
	if conn.North && conn.South && !conn.West && !conn.East {
		w = 0.3125
		e = 0.6875
	} else if !conn.North && !conn.South && conn.West && conn.East {
		n = 0.3125
		s = 0.6875
	}

	return FenceBoundingBox{fx + w, fy, fz + n, fx + e, fy + FenceHeight, fz + s}
}

func init() {
	Registry.Register(NewFenceBlock())
	Registry.Register(NewNetherBrickFenceBlock())
	Registry.Register(NewStoneWallBlock())
}
