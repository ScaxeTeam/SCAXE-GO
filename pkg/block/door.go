package block

// door.go — 门方块基类
// 对应 PHP abstract class Door extends Transparent implements IRedstone
//
// 门是双方块（占据上下两格），meta 位编码：
//   下半: bit0-1=朝向(0-3), bit2=是否打开, bit3=0
//   上半: bit0=铰链侧(0=左,1=右), bit1=红石供电, bit3=1
//
// 具体门类型（木门/铁门等）嵌入 DoorBase，只需设置 BlockID 和基础属性。

// DoorBase 门方块基类
// 对应 PHP abstract class Door extends Transparent
type DoorBase struct {
	TransparentBase
}

// IsSolid 门不是完整实心方块
// 对应 PHP Door::isSolid() { return false; }
func (b *DoorBase) IsSolid() bool {
	return false
}

// CanBeActivated 门可以被右键激活（开/关）
// 对应 PHP Door::canBeActivated() { return true; }
func (b *DoorBase) CanBeActivated() bool {
	return true
}

// OnActivate 门的右键交互 — 开/关切换
// 实际 meta 变更由服务器层通过 DoorOnActivate() 执行
func (b *DoorBase) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

// ---------- Meta 位操作 ----------

// DoorMeta 位掩码常量
const (
	DoorMetaDirection  = 0x03 // bit0-1: 方向 (0-3)
	DoorMetaOpen       = 0x04 // bit2: 是否打开（下半）
	DoorMetaTop        = 0x08 // bit3: 是否为上半
	DoorMetaHingeRight = 0x01 // bit0: 铰链在右侧（上半）
	DoorMetaPowered    = 0x02 // bit1: 红石供电（上半）
)

// IsTopHalf 判断 meta 是否为上半门
func DoorIsTopHalf(meta uint8) bool {
	return meta&DoorMetaTop != 0
}

// IsOpen 判断门是否打开（需传入下半的 meta）
func DoorIsOpen(bottomMeta uint8) bool {
	return bottomMeta&DoorMetaOpen != 0
}

// GetDirection 获取门的朝向（从下半 meta）
func DoorGetDirection(bottomMeta uint8) uint8 {
	return bottomMeta & DoorMetaDirection
}

// ToggleOpen 切换下半门的开关状态，返回新 meta
func DoorToggleOpen(bottomMeta uint8) uint8 {
	return bottomMeta ^ DoorMetaOpen
}

// IsHingeRight 判断铰链是否在右侧（从上半 meta）
func DoorIsHingeRight(topMeta uint8) bool {
	return topMeta&DoorMetaHingeRight != 0
}

// ---------- 放置 ----------

// DoorPlacementFaces 放置时玩家朝向 → 检查侧面的映射
// 对应 PHP Door::place() 中的 $face 映射
var DoorPlacementFaces = [4]int{3, 4, 2, 5}

// GetDoorPlacementMeta 获取放置门时下半的 meta
// playerDirection: 玩家朝向 (0-3)
func GetDoorPlacementMeta(playerDirection int) uint8 {
	return uint8(playerDirection) & DoorMetaDirection
}

// GetDoorTopMeta 计算上半门的 meta
// hingeRight: 铰链是否在右侧
func GetDoorTopMeta(hingeRight bool) uint8 {
	meta := uint8(DoorMetaTop)
	if hingeRight {
		meta |= DoorMetaHingeRight
	}
	return meta
}

// ShouldHingeRight 判断铰链应该在哪一侧
// 对应 PHP Door::place() 中铰链侧的判断逻辑:
//   如果右侧有同类门，或（左侧不透明且右侧透明），铰链在右侧
//
// 参数:
//   sameBlockOnRight: 右侧(direction+2)%4方向是否有同类门
//   leftTransparent: 左侧(direction方向)是否透明
//   rightTransparent: 右侧是否透明
func ShouldHingeRight(sameBlockOnRight bool, leftTransparent bool, rightTransparent bool) bool {
	return sameBlockOnRight || (!leftTransparent && rightTransparent)
}

// ---------- 碰撞箱 ----------

// DoorThickness 门的厚度 (3/16)
const DoorThickness = 0.1875

// DoorBoundingBox 门的碰撞箱
type DoorBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}

// GetDoorBoundingBox 根据门的完整状态计算碰撞箱
// 对应 PHP Door::recalculateBoundingBox() 的完整逻辑
//
// 参数:
//   x, y, z: 方块坐标
//   direction: 朝向 (0-3)
//   isOpen: 是否打开
//   isRight: 铰链是否在右侧
func GetDoorBoundingBox(x, y, z int, direction uint8, isOpen bool, isRight bool) DoorBoundingBox {
	fx, fy, fz := float64(x), float64(y), float64(z)
	f := DoorThickness

	// 默认完整方块
	bb := DoorBoundingBox{fx, fy, fz, fx + 1, fy + 1, fz + 1}

	switch direction & 0x03 {
	case 0: // 面朝东 → 门板在西侧(x-)
		if isOpen {
			if !isRight {
				bb = DoorBoundingBox{fx, fy, fz, fx + 1, fy + 1, fz + f}
			} else {
				bb = DoorBoundingBox{fx, fy, fz + 1 - f, fx + 1, fy + 1, fz + 1}
			}
		} else {
			bb = DoorBoundingBox{fx, fy, fz, fx + f, fy + 1, fz + 1}
		}
	case 1: // 面朝南 → 门板在北侧(z-)
		if isOpen {
			if !isRight {
				bb = DoorBoundingBox{fx + 1 - f, fy, fz, fx + 1, fy + 1, fz + 1}
			} else {
				bb = DoorBoundingBox{fx, fy, fz, fx + f, fy + 1, fz + 1}
			}
		} else {
			bb = DoorBoundingBox{fx, fy, fz, fx + 1, fy + 1, fz + f}
		}
	case 2: // 面朝西 → 门板在东侧(x+)
		if isOpen {
			if !isRight {
				bb = DoorBoundingBox{fx, fy, fz + 1 - f, fx + 1, fy + 1, fz + 1}
			} else {
				bb = DoorBoundingBox{fx, fy, fz, fx + 1, fy + 1, fz + f}
			}
		} else {
			bb = DoorBoundingBox{fx + 1 - f, fy, fz, fx + 1, fy + 1, fz + 1}
		}
	case 3: // 面朝北 → 门板在南侧(z+)
		if isOpen {
			if !isRight {
				bb = DoorBoundingBox{fx, fy, fz, fx + f, fy + 1, fz + 1}
			} else {
				bb = DoorBoundingBox{fx + 1 - f, fy, fz, fx + 1, fy + 1, fz + 1}
			}
		} else {
			bb = DoorBoundingBox{fx, fy, fz + 1 - f, fx + 1, fy + 1, fz + 1}
		}
	}

	return bb
}

// ---------- 完整 damage 解析 ----------

// DoorFullState 门的完整状态（合并上下两半的 meta）
type DoorFullState struct {
	Direction  uint8 // 朝向 0-3
	IsOpen     bool  // 是否打开
	IsTopHalf  bool  // 是否为上半
	HingeRight bool  // 铰链在右侧
}

// ParseDoorState 从上下两半的 meta 解析出完整门状态
// 对应 PHP Door::getFullDamage()
func ParseDoorState(thisMeta uint8, otherHalfMeta uint8) DoorFullState {
	var topMeta, bottomMeta uint8
	isTop := DoorIsTopHalf(thisMeta)

	if isTop {
		topMeta = thisMeta
		bottomMeta = otherHalfMeta
	} else {
		topMeta = otherHalfMeta
		bottomMeta = thisMeta
	}

	return DoorFullState{
		Direction:  DoorGetDirection(bottomMeta),
		IsOpen:     DoorIsOpen(bottomMeta),
		IsTopHalf:  isTop,
		HingeRight: DoorIsHingeRight(topMeta),
	}
}
