package block
type DoorBase struct {
	TransparentBase
}
func (b *DoorBase) IsSolid() bool {
	return false
}
func (b *DoorBase) CanBeActivated() bool {
	return true
}
func (b *DoorBase) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}
const (
	DoorMetaDirection  = 0x03
	DoorMetaOpen       = 0x04
	DoorMetaTop        = 0x08
	DoorMetaHingeRight = 0x01
	DoorMetaPowered    = 0x02
)
func DoorIsTopHalf(meta uint8) bool {
	return meta&DoorMetaTop != 0
}
func DoorIsOpen(bottomMeta uint8) bool {
	return bottomMeta&DoorMetaOpen != 0
}
func DoorGetDirection(bottomMeta uint8) uint8 {
	return bottomMeta & DoorMetaDirection
}
func DoorToggleOpen(bottomMeta uint8) uint8 {
	return bottomMeta ^ DoorMetaOpen
}
func DoorIsHingeRight(topMeta uint8) bool {
	return topMeta&DoorMetaHingeRight != 0
}
var DoorPlacementFaces = [4]int{3, 4, 2, 5}
func GetDoorPlacementMeta(playerDirection int) uint8 {
	return uint8(playerDirection) & DoorMetaDirection
}
func GetDoorTopMeta(hingeRight bool) uint8 {
	meta := uint8(DoorMetaTop)
	if hingeRight {
		meta |= DoorMetaHingeRight
	}
	return meta
}
func ShouldHingeRight(sameBlockOnRight bool, leftTransparent bool, rightTransparent bool) bool {
	return sameBlockOnRight || (!leftTransparent && rightTransparent)
}
const DoorThickness = 0.1875
type DoorBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}
func GetDoorBoundingBox(x, y, z int, direction uint8, isOpen bool, isRight bool) DoorBoundingBox {
	fx, fy, fz := float64(x), float64(y), float64(z)
	f := DoorThickness
	bb := DoorBoundingBox{fx, fy, fz, fx + 1, fy + 1, fz + 1}

	switch direction & 0x03 {
	case 0:
		if isOpen {
			if !isRight {
				bb = DoorBoundingBox{fx, fy, fz, fx + 1, fy + 1, fz + f}
			} else {
				bb = DoorBoundingBox{fx, fy, fz + 1 - f, fx + 1, fy + 1, fz + 1}
			}
		} else {
			bb = DoorBoundingBox{fx, fy, fz, fx + f, fy + 1, fz + 1}
		}
	case 1:
		if isOpen {
			if !isRight {
				bb = DoorBoundingBox{fx + 1 - f, fy, fz, fx + 1, fy + 1, fz + 1}
			} else {
				bb = DoorBoundingBox{fx, fy, fz, fx + f, fy + 1, fz + 1}
			}
		} else {
			bb = DoorBoundingBox{fx, fy, fz, fx + 1, fy + 1, fz + f}
		}
	case 2:
		if isOpen {
			if !isRight {
				bb = DoorBoundingBox{fx, fy, fz + 1 - f, fx + 1, fy + 1, fz + 1}
			} else {
				bb = DoorBoundingBox{fx, fy, fz, fx + 1, fy + 1, fz + f}
			}
		} else {
			bb = DoorBoundingBox{fx + 1 - f, fy, fz, fx + 1, fy + 1, fz + 1}
		}
	case 3:
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
type DoorFullState struct {
	Direction  uint8
	IsOpen     bool
	IsTopHalf  bool
	HingeRight bool
}
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
