package structure

import "fmt"

type BoundingBox struct {
	MinX, MinY, MinZ int
	MaxX, MaxY, MaxZ int
}

func NewBoundingBox(x1, y1, z1, x2, y2, z2 int) *BoundingBox {
	b := &BoundingBox{
		MinX: x1, MinY: y1, MinZ: z1,
		MaxX: x2, MaxY: y2, MaxZ: z2,
	}

	if b.MinX > b.MaxX {
		b.MinX, b.MaxX = b.MaxX, b.MinX
	}
	if b.MinY > b.MaxY {
		b.MinY, b.MaxY = b.MaxY, b.MinY
	}
	if b.MinZ > b.MaxZ {
		b.MinZ, b.MaxZ = b.MaxZ, b.MinZ
	}
	return b
}

func (b *BoundingBox) IntersectsWith(other *BoundingBox) bool {
	return b.MaxX >= other.MinX && b.MinX <= other.MaxX &&
		b.MaxZ >= other.MinZ && b.MinZ <= other.MaxZ &&
		b.MaxY >= other.MinY && b.MinY <= other.MaxY
}

func (b *BoundingBox) String() string {
	return fmt.Sprintf("(%d,%d,%d -> %d,%d,%d)", b.MinX, b.MinY, b.MinZ, b.MaxX, b.MaxY, b.MaxZ)
}

func (b *BoundingBox) ResultIsInside(x, y, z int) bool {
	return x >= b.MinX && x <= b.MaxX && z >= b.MinZ && z <= b.MaxZ && y >= b.MinY && y <= b.MaxY
}

func (b *BoundingBox) Offset(x, y, z int) {
	b.MinX += x
	b.MinY += y
	b.MinZ += z
	b.MaxX += x
	b.MaxY += y
	b.MaxZ += z
}

func GetComponentToAddBoundingBox(structureMinX, structureMinY, structureMinZ, xMin, yMin, zMin, xMax, yMax, zMax int, facing int) *BoundingBox {
	switch facing {
	case 0:

		return NewBoundingBox(structureMinX+xMin, structureMinY+yMin, structureMinZ+zMin, structureMinX+xMax-1+xMin, structureMinY+yMax-1+yMin, structureMinZ+zMax-1+zMin)
	case 2:
		return NewBoundingBox(structureMinX+xMin, structureMinY+yMin, structureMinZ-zMax+1+zMin, structureMinX+xMax-1+xMin, structureMinY+yMax-1+yMin, structureMinZ+zMin)
	case 1:
		return NewBoundingBox(structureMinX-zMax+1+zMin, structureMinY+yMin, structureMinZ+xMin, structureMinX+zMin, structureMinY+yMax-1+yMin, structureMinZ+xMax-1+xMin)
	case 3:
		return NewBoundingBox(structureMinX+zMin, structureMinY+yMin, structureMinZ+xMin, structureMinX+zMax-1+zMin, structureMinY+yMax-1+yMin, structureMinZ+xMax-1+xMin)
	}

	return NewBoundingBox(structureMinX+xMin, structureMinY+yMin, structureMinZ+zMin, structureMinX+xMax-1+xMin, structureMinY+yMax-1+yMin, structureMinZ+zMax-1+zMin)
}
