package block

// trapdoor.go — 活板门方块（木质 + 铁质）
// 对应 PHP class Trapdoor extends Transparent 和 class IronTrapdoor extends Trapdoor
//
// Meta 位编码:
//   bit0-1: 朝向 (0-3)
//   bit2:   上半 (MASK_UPPER = 0x04)
//   bit3:   打开 (MASK_OPENED = 0x08)

// TrapdoorBlock 活板门方块
type TrapdoorBlock struct {
	TransparentBase
}

// Trapdoor meta 掩码
const (
	TrapdoorMaskDirection = 0x03
	TrapdoorMaskUpper     = 0x04
	TrapdoorMaskOpened    = 0x08
)

// NewTrapdoorBlock 创建木质活板门
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

// NewIronTrapdoorBlock 创建铁质活板门
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

// CanBeActivated 活板门可以被右键激活
func (b *TrapdoorBlock) CanBeActivated() bool {
	return true
}

// OnActivate 活板门右键交互 — 开/关切换
// 实际 meta 变更由服务器层通过 TrapdoorOnActivate() 执行
func (b *TrapdoorBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	// 铁活板门不可手动开（仅红石）
	if b.BlockID == IRON_TRAPDOOR {
		return false
	}
	return true
}

// IsSolid 活板门不是完整实心方块
func (b *TrapdoorBlock) IsSolid() bool {
	return false
}

// GetFuelTime 木质活板门可作燃料，铁质不行
func (b *TrapdoorBlock) GetFuelTime() int {
	if b.BlockID == TRAPDOOR {
		return 300
	}
	return 0
}

// GetDrops 活板门掉落自身
func (b *TrapdoorBlock) GetDrops(toolType, toolTier int) []Drop {
	if b.BlockID == IRON_TRAPDOOR && toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(b.BlockID), Meta: 0, Count: 1}}
}

// ---------- Meta 工具函数 ----------

// TrapdoorIsOpen 判断活板门是否打开
func TrapdoorIsOpen(meta uint8) bool {
	return meta&TrapdoorMaskOpened != 0
}

// TrapdoorIsUpper 判断活板门是否在上半
func TrapdoorIsUpper(meta uint8) bool {
	return meta&TrapdoorMaskUpper != 0
}

// TrapdoorGetDirection 获取朝向 (0-3)
func TrapdoorGetDirection(meta uint8) uint8 {
	return meta & TrapdoorMaskDirection
}

// TrapdoorToggleOpen 切换开关，返回新 meta
func TrapdoorToggleOpen(meta uint8) uint8 {
	return meta ^ TrapdoorMaskOpened
}

// ---------- 放置 ----------

// TrapdoorDirectionToMeta 玩家朝向 → meta 映射
// 对应 PHP Trapdoor::place() 中的 $directions
var TrapdoorDirectionToMeta = [4]uint8{1, 3, 0, 2}

// GetTrapdoorPlacementMeta 获取放置活板门的 meta
// playerDirection: 玩家朝向(0-3), clickY: 点击的Y坐标(0-1), face: 放置面
// 如果 clickY > 0.5 且不是上面放置，或者从下面放置，则为上半
func GetTrapdoorPlacementMeta(playerDirection int, clickY float64, face int) uint8 {
	meta := TrapdoorDirectionToMeta[playerDirection&0x03]
	if (clickY > 0.5 && face != 1) || face == 0 { // face 1=上, 0=下
		meta |= TrapdoorMaskUpper
	}
	return meta
}

// ---------- 碰撞箱 ----------

// TrapdoorBoundingBox 活板门碰撞箱
type TrapdoorBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}

// GetTrapdoorBoundingBox 计算活板门碰撞箱
// 对应 PHP Trapdoor::recalculateBoundingBox()
func GetTrapdoorBoundingBox(x, y, z int, meta uint8) TrapdoorBoundingBox {
	fx, fy, fz := float64(x), float64(y), float64(z)
	f := 0.1875 // 3/16 厚度

	isOpen := TrapdoorIsOpen(meta)
	isUpper := TrapdoorIsUpper(meta)
	dir := TrapdoorGetDirection(meta)

	// 关闭时：水平薄片
	if !isOpen {
		if isUpper {
			return TrapdoorBoundingBox{fx, fy + 1 - f, fz, fx + 1, fy + 1, fz + 1}
		}
		return TrapdoorBoundingBox{fx, fy, fz, fx + 1, fy + f, fz + 1}
	}

	// 打开时：垂直薄片（按朝向）
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

// ---------- 注册 ----------

func init() {
	Registry.Register(NewTrapdoorBlock())
	Registry.Register(NewIronTrapdoorBlock())
}
