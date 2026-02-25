package block

// fence.go — 围栏/栅栏 + 石墙方块
// 对应 PHP: Fence, NetherBrickFence, StoneWall
//
// 围栏/石墙共同特点:
//   - 高1.5格（实体无法跳过）
//   - 根据相邻方块动态连接（围栏/围门/实心不透明方块）
//   - 碰撞箱随连接方向变化

// ============ 围栏 ============

// FenceBlock 木质围栏（6种木材变种，通过 meta 区分）
// 对应 PHP class Fence extends Transparent
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

// IsSolid 围栏视为实心（用于碰撞）
func (b *FenceBlock) IsSolid() bool { return true }

// GetFuelTime 木围栏可燃
func (b *FenceBlock) GetFuelTime() int { return 300 }

// GetDrops 围栏掉落自身
func (b *FenceBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(FENCE), Meta: 0, Count: 1}}
}

// 围栏变种名称 (meta 0-5)
var FenceVariantNames = [6]string{
	"Oak Fence", "Spruce Fence", "Birch Fence",
	"Jungle Fence", "Acacia Fence", "Dark Oak Fence",
}

// ============ 下界砖围栏 ============

// NetherBrickFenceBlock 下界砖围栏
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

// ============ 石墙 ============

// StoneWallBlock 石墙（圆石墙/苔石墙）
// 对应 PHP class StoneWall extends Transparent
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

// ============ 连接判断 ============

// FenceConnectable 围栏连接方向
type FenceConnections struct {
	North, South, West, East bool
}

// CanFenceConnect 判断围栏是否可以连接目标方块
// 对应 PHP Fence::canConnect()
// 围栏连接: 同类围栏、栅栏门、实心不透明方块
func CanFenceConnect(targetID uint8, targetIsSolid, targetIsTransparent bool) bool {
	if targetID == FENCE || targetID == NETHER_BRICK_FENCE {
		return true
	}
	if IsFenceGate(targetID) {
		return true
	}
	return targetIsSolid && !targetIsTransparent
}

// CanWallConnect 判断石墙是否可以连接目标方块
// 对应 PHP StoneWall::canConnect()
func CanWallConnect(targetID uint8, targetIsSolid, targetIsTransparent bool) bool {
	if targetID == STONE_WALL || IsFenceGate(targetID) {
		return true
	}
	return targetIsSolid && !targetIsTransparent
}

// IsFenceGate 判断方块是否为栅栏门
func IsFenceGate(blockID uint8) bool {
	switch blockID {
	case FENCE_GATE, FENCE_GATE_SPRUCE, FENCE_GATE_BIRCH,
		FENCE_GATE_JUNGLE, FENCE_GATE_DARK_OAK, FENCE_GATE_ACACIA:
		return true
	}
	return false
}

// ============ 碰撞箱 ============

const (
	FenceThickness = 0.25 // 围栏柱厚度
	FenceHeight    = 1.5  // 围栏高度
	WallThickness  = 0.5  // 石墙柱厚度
)

// FenceBoundingBox 围栏碰撞箱
type FenceBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}

// GetFenceBoundingBox 计算围栏碰撞箱（简化版，单 AABB）
// 对应 PHP Fence::recalculateBoundingBox()
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

// GetWallBoundingBox 计算石墙碰撞箱
// 对应 PHP StoneWall::recalculateBoundingBox()
// 石墙特殊: 只有南北连接或只有东西连接时，变窄
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

	// 直线连接时变窄
	if conn.North && conn.South && !conn.West && !conn.East {
		w = 0.3125
		e = 0.6875
	} else if !conn.North && !conn.South && conn.West && conn.East {
		n = 0.3125
		s = 0.6875
	}

	return FenceBoundingBox{fx + w, fy, fz + n, fx + e, fy + FenceHeight, fz + s}
}

// ============ 注册 ============

func init() {
	Registry.Register(NewFenceBlock())
	Registry.Register(NewNetherBrickFenceBlock())
	Registry.Register(NewStoneWallBlock())
}
