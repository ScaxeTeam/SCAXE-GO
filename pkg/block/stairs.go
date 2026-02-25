package block

// stairs.go — 楼梯方块基类 + 所有变种
// 对应 PHP abstract class Stair extends Transparent + 13个变种
//
// Meta 位编码:
//   bit0-1: 朝向 (0-3)
//   bit2:   倒置 (0x04)
//
// 碰撞箱: 由两个 AABB 组成 — 底半(或顶半) + 台阶部分

// StairBlock 楼梯方块
type StairBlock struct {
	TransparentBase
}

// Stair meta 掩码
const (
	StairMaskDirection  = 0x03
	StairMaskUpsideDown = 0x04
)

func newStair(blockID uint8, name string, toolType int) *StairBlock {
	return &StairBlock{
		TransparentBase: TransparentBase{
			BlockID:         blockID,
			BlockName:       name,
			BlockHardness:   2,
			BlockResistance: 15,
			BlockToolType:   toolType,
			BlockCanPlace:   true,
		},
	}
}

// GetDrops 石质楼梯需镐，木质直接掉落
func (b *StairBlock) GetDrops(toolType, toolTier int) []Drop {
	if b.BlockToolType == ToolTypePickaxe && toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(b.BlockID), Meta: 0, Count: 1}}
}

// ---------- Meta 工具函数 ----------

func StairIsUpsideDown(meta uint8) bool  { return meta&StairMaskUpsideDown != 0 }
func StairGetDirection(meta uint8) uint8 { return meta & StairMaskDirection }

// StairDirectionToMeta 放置朝向映射
// 对应 PHP Stair::place() 中的 $faces
var StairDirectionToMeta = [4]uint8{0, 2, 1, 3}

// GetStairPlacementMeta 获取放置楼梯的 meta
func GetStairPlacementMeta(playerDirection int, clickY float64, face int) uint8 {
	meta := StairDirectionToMeta[playerDirection&0x03]
	if (clickY > 0.5 && face != 1) || face == 0 {
		meta |= StairMaskUpsideDown
	}
	return meta
}

// ---------- 碰撞箱 ----------

// StairBoundingBox 楼梯碰撞箱（由两个 AABB 组成）
type StairBoundingBox struct {
	// Slab 底半（或顶半）— 全宽半高
	SlabMinX, SlabMinY, SlabMinZ float64
	SlabMaxX, SlabMaxY, SlabMaxZ float64
	// Step 台阶部分 — 半宽半高
	StepMinX, StepMinY, StepMinZ float64
	StepMaxX, StepMaxY, StepMaxZ float64
}

// GetStairBoundingBoxes 计算楼梯碰撞箱
// 对应 PHP Stair::recalculateCollisionBoxes()
func GetStairBoundingBoxes(x, y, z int, meta uint8) StairBoundingBox {
	fx, fy, fz := float64(x), float64(y), float64(z)

	upsideDown := StairIsUpsideDown(meta)
	dir := StairGetDirection(meta)

	// 底半(Slab)
	var slabMinY, slabMaxY float64
	if upsideDown {
		slabMinY = 0.5
	}
	slabMaxY = slabMinY + 0.5

	// 台阶部分(Step) — 另一半高度
	var stepMinY, stepMaxY float64
	if upsideDown {
		stepMinY = 0
	} else {
		stepMinY = 0.5
	}
	stepMaxY = stepMinY + 0.5

	// 台阶方向
	stepMinX, stepMinZ := 0.0, 0.0
	stepMaxX, stepMaxZ := 1.0, 1.0
	switch dir {
	case 0:
		stepMinX = 0.5
	case 1:
		stepMaxX = 0.5
	case 2:
		stepMinZ = 0.5
	case 3:
		stepMaxZ = 0.5
	}

	return StairBoundingBox{
		SlabMinX: fx, SlabMinY: fy + slabMinY, SlabMinZ: fz,
		SlabMaxX: fx + 1, SlabMaxY: fy + slabMaxY, SlabMaxZ: fz + 1,
		StepMinX: fx + stepMinX, StepMinY: fy + stepMinY, StepMinZ: fz + stepMinZ,
		StepMaxX: fx + stepMaxX, StepMaxY: fy + stepMaxY, StepMaxZ: fz + stepMaxZ,
	}
}

// ---------- 注册 ----------

func init() {
	// 石质楼梯（需镐）
	Registry.Register(newStair(COBBLESTONE_STAIRS, "Cobblestone Stairs", ToolTypePickaxe))
	Registry.Register(newStair(BRICK_STAIRS, "Brick Stairs", ToolTypePickaxe))
	Registry.Register(newStair(STONE_BRICK_STAIRS, "Stone Brick Stairs", ToolTypePickaxe))
	Registry.Register(newStair(NETHER_BRICKS_STAIRS, "Nether Brick Stairs", ToolTypePickaxe))
	Registry.Register(newStair(SANDSTONE_STAIRS, "Sandstone Stairs", ToolTypePickaxe))
	Registry.Register(newStair(QUARTZ_STAIRS, "Quartz Stairs", ToolTypePickaxe))
	Registry.Register(newStair(RED_SANDSTONE_STAIRS, "Red Sandstone Stairs", ToolTypePickaxe))

	// 木质楼梯（斧）
	Registry.Register(newStair(WOOD_STAIRS, "Oak Wood Stairs", ToolTypeAxe))
	Registry.Register(newStair(SPRUCE_WOOD_STAIRS, "Spruce Wood Stairs", ToolTypeAxe))
	Registry.Register(newStair(BIRCH_WOOD_STAIRS, "Birch Wood Stairs", ToolTypeAxe))
	Registry.Register(newStair(JUNGLE_WOOD_STAIRS, "Jungle Wood Stairs", ToolTypeAxe))
	Registry.Register(newStair(ACACIA_WOOD_STAIRS, "Acacia Wood Stairs", ToolTypeAxe))
	Registry.Register(newStair(DARK_OAK_WOOD_STAIRS, "Dark Oak Wood Stairs", ToolTypeAxe))
}
