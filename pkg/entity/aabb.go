package entity

type AxisAlignedBB struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}

func NewAxisAlignedBB(minX, minY, minZ, maxX, maxY, maxZ float64) *AxisAlignedBB {
	if minX > maxX {
		minX, maxX = maxX, minX
	}
	if minY > maxY {
		minY, maxY = maxY, minY
	}
	if minZ > maxZ {
		minZ, maxZ = maxZ, minZ
	}
	return &AxisAlignedBB{
		MinX: minX, MinY: minY, MinZ: minZ,
		MaxX: maxX, MaxY: maxY, MaxZ: maxZ,
	}
}

func (bb *AxisAlignedBB) SetBounds(minX, minY, minZ, maxX, maxY, maxZ float64) {
	bb.MinX = minX
	bb.MinY = minY
	bb.MinZ = minZ
	bb.MaxX = maxX
	bb.MaxY = maxY
	bb.MaxZ = maxZ
}

func (bb *AxisAlignedBB) Offset(x, y, z float64) *AxisAlignedBB {
	return NewAxisAlignedBB(
		bb.MinX+x, bb.MinY+y, bb.MinZ+z,
		bb.MaxX+x, bb.MaxY+y, bb.MaxZ+z,
	)
}

func (bb *AxisAlignedBB) Expand(x, y, z float64) *AxisAlignedBB {
	minX, minY, minZ := bb.MinX, bb.MinY, bb.MinZ
	maxX, maxY, maxZ := bb.MaxX, bb.MaxY, bb.MaxZ

	if x < 0 {
		minX += x
	} else {
		maxX += x
	}

	if y < 0 {
		minY += y
	} else {
		maxY += y
	}

	if z < 0 {
		minZ += z
	} else {
		maxZ += z
	}

	return NewAxisAlignedBB(minX, minY, minZ, maxX, maxY, maxZ)
}

func (bb *AxisAlignedBB) IntersectsWith(other *AxisAlignedBB) bool {
	return other.MaxX > bb.MinX && other.MinX < bb.MaxX &&
		other.MaxY > bb.MinY && other.MinY < bb.MaxY &&
		other.MaxZ > bb.MinZ && other.MinZ < bb.MaxZ
}

func (bb *AxisAlignedBB) CalculateIntercept(pos, vec *Vector3) *Vector3 {

	return nil
}

func (bb *AxisAlignedBB) AddCoord(x, y, z float64) *AxisAlignedBB {
	minX, minY, minZ := bb.MinX, bb.MinY, bb.MinZ
	maxX, maxY, maxZ := bb.MaxX, bb.MaxY, bb.MaxZ

	if x < 0 {
		minX += x
	} else if x > 0 {
		maxX += x
	}

	if y < 0 {
		minY += y
	} else if y > 0 {
		maxY += y
	}

	if z < 0 {
		minZ += z
	} else if z > 0 {
		maxZ += z
	}

	return NewAxisAlignedBB(minX, minY, minZ, maxX, maxY, maxZ)
}

func (bb *AxisAlignedBB) CalculateXOffset(other *AxisAlignedBB, offsetX float64) float64 {
	if other.MaxY > bb.MinY && other.MinY < bb.MaxY && other.MaxZ > bb.MinZ && other.MinZ < bb.MaxZ {
		if offsetX > 0 && other.MaxX <= bb.MinX {
			delta := bb.MinX - other.MaxX
			if delta < offsetX {
				offsetX = delta
			}
		} else if offsetX < 0 && other.MinX >= bb.MaxX {
			delta := bb.MaxX - other.MinX
			if delta > offsetX {
				offsetX = delta
			}
		}
	}
	return offsetX
}

func (bb *AxisAlignedBB) CalculateYOffset(other *AxisAlignedBB, offsetY float64) float64 {
	if other.MaxX > bb.MinX && other.MinX < bb.MaxX && other.MaxZ > bb.MinZ && other.MinZ < bb.MaxZ {
		if offsetY > 0 && other.MaxY <= bb.MinY {
			delta := bb.MinY - other.MaxY
			if delta < offsetY {
				offsetY = delta
			}
		} else if offsetY < 0 && other.MinY >= bb.MaxY {
			delta := bb.MaxY - other.MinY
			if delta > offsetY {
				offsetY = delta
			}
		}
	}
	return offsetY
}

func (bb *AxisAlignedBB) CalculateZOffset(other *AxisAlignedBB, offsetZ float64) float64 {
	if other.MaxX > bb.MinX && other.MinX < bb.MaxX && other.MaxY > bb.MinY && other.MinY < bb.MaxY {
		if offsetZ > 0 && other.MaxZ <= bb.MinZ {
			delta := bb.MinZ - other.MaxZ
			if delta < offsetZ {
				offsetZ = delta
			}
		} else if offsetZ < 0 && other.MinZ >= bb.MaxZ {
			delta := bb.MaxZ - other.MinZ
			if delta > offsetZ {
				offsetZ = delta
			}
		}
	}
	return offsetZ
}

func (bb *AxisAlignedBB) Clone() *AxisAlignedBB {
	return NewAxisAlignedBB(bb.MinX, bb.MinY, bb.MinZ, bb.MaxX, bb.MaxY, bb.MaxZ)
}

func (bb *AxisAlignedBB) Grow(x, y, z float64) *AxisAlignedBB {
	return NewAxisAlignedBB(
		bb.MinX-x, bb.MinY-y, bb.MinZ-z,
		bb.MaxX+x, bb.MaxY+y, bb.MaxZ+z,
	)
}
