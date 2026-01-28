package structure

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type StructureStart struct {
	ChunkX      int
	ChunkZ      int
	Components  []StructureComponent
	BoundingBox *BoundingBox
}

func NewStructureStart(chunkX, chunkZ int) *StructureStart {
	return &StructureStart{
		ChunkX:      chunkX,
		ChunkZ:      chunkZ,
		Components:  make([]StructureComponent, 0),
		BoundingBox: NewBoundingBox(0, 0, 0, 0, 0, 0),
	}
}

type WorldAccess interface {
	GetBlock(x, y, z int) (byte, byte)
	SetBlock(x, y, z int, id byte, meta byte)
}

func (s *StructureStart) GenerateStructure(w WorldAccess, rnd *rand.Random, box *BoundingBox) {

	for _, c := range s.Components {
		bb := c.GetBoundingBox()
		if bb.IntersectsWith(box) {
			c.AddComponentParts(w, rnd, box)

		}
	}
}

func (s *StructureStart) UpdateBoundingBox() {
	if len(s.Components) == 0 {
		return
	}
	minX, minY, minZ := 1000000, 1000000, 1000000
	maxX, maxY, maxZ := -1000000, -1000000, -1000000

	for _, c := range s.Components {
		bb := c.GetBoundingBox()
		if bb.MinX < minX {
			minX = bb.MinX
		}
		if bb.MinY < minY {
			minY = bb.MinY
		}
		if bb.MinZ < minZ {
			minZ = bb.MinZ
		}
		if bb.MaxX > maxX {
			maxX = bb.MaxX
		}
		if bb.MaxY > maxY {
			maxY = bb.MaxY
		}
		if bb.MaxZ > maxZ {
			maxZ = bb.MaxZ
		}
	}
	s.BoundingBox = NewBoundingBox(minX, minY, minZ, maxX, maxY, maxZ)
}

type StructureComponent interface {
	BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random)
	AddComponentParts(w WorldAccess, rnd *rand.Random, box *BoundingBox) bool
	GetBoundingBox() *BoundingBox
	GetComponentType() int
}

type StructureComponentBase struct {
	BoundingBox   *BoundingBox
	CoordBaseMode int
	ComponentType int
}

func (c *StructureComponentBase) GetBoundingBox() *BoundingBox {
	return c.BoundingBox
}

func (c *StructureComponentBase) GetComponentType() int {
	return c.ComponentType
}

func (c *StructureComponentBase) GetXWithOffset(x, z int) int {
	switch c.CoordBaseMode {
	case 0:
		return c.BoundingBox.MinX + x
	case 1:
		return c.BoundingBox.MaxX - z
	case 2:
		return c.BoundingBox.MaxX - x
	case 3:
		return c.BoundingBox.MinX + z
	}
	return x
}

func (c *StructureComponentBase) GetYWithOffset(y int) int {
	if c.CoordBaseMode == -1 {
		return y
	}
	return y + c.BoundingBox.MinY
}

func (c *StructureComponentBase) GetZWithOffset(x, z int) int {
	switch c.CoordBaseMode {
	case 0:
		return c.BoundingBox.MinZ + z
	case 1:
		return c.BoundingBox.MinZ + x
	case 2:
		return c.BoundingBox.MaxZ - z
	case 3:
		return c.BoundingBox.MaxZ - x
	}
	return z
}

func (c *StructureComponentBase) SetBlockState(w WorldAccess, id byte, meta byte, x, y, z int, box *BoundingBox) {
	worldX := c.GetXWithOffset(x, z)
	worldY := c.GetYWithOffset(y)
	worldZ := c.GetZWithOffset(x, z)

	if box.ResultIsInside(worldX, worldY, worldZ) {
		meta = c.GetMetadataWithOffset(id, meta)
		w.SetBlock(worldX, worldY, worldZ, id, meta)
	}
}

func (c *StructureComponentBase) GetMetadataWithOffset(id byte, meta byte) byte {
	if c.CoordBaseMode == -1 {
		return meta
	}

	if isStairs(id) {

		switch c.CoordBaseMode {
		case 0:
			return meta
		case 1:

			switch meta {
			case 2:
				return 1
			case 1:
				return 3
			case 3:
				return 0
			case 0:
				return 2
			}
		case 2:

			switch meta {
			case 2:
				return 3
			case 3:
				return 2
			case 0:
				return 1
			case 1:
				return 0
			}
		case 3:

			switch meta {
			case 2:
				return 0
			case 0:
				return 3
			case 3:
				return 1
			case 1:
				return 2
			}
		}
	}

	if id == 54 || id == 130 || id == 146 || id == 65 || id == 61 || id == 62 {
		switch c.CoordBaseMode {
		case 0:
			return meta
		case 1:

			switch meta {
			case 2:
				return 5
			case 3:
				return 4
			case 4:
				return 2
			case 5:
				return 3
			}
		case 2:

			switch meta {
			case 2:
				return 3
			case 3:
				return 2
			case 4:
				return 5
			case 5:
				return 4
			}
		case 3:

			switch meta {
			case 2:
				return 4
			case 3:
				return 5
			case 4:
				return 3
			case 5:
				return 2
			}
		}
	}

	if id == 17 || id == 162 {
		axis := meta & 12
		if axis == 4 || axis == 8 {
			if c.CoordBaseMode == 1 || c.CoordBaseMode == 3 {

				if axis == 4 {
					return (meta & 3) | 8
				}
				if axis == 8 {
					return (meta & 3) | 4
				}
			}
		}
	}

	return meta
}

func isStairs(id byte) bool {

	return id == 53 || id == 67 || id == 108 || id == 109 || id == 114 || id == 128 || id == 134 || id == 135 || id == 136 || id == 156 || id == 163 || id == 164 || id == 180
}

func (c *StructureComponentBase) FillWithBlocks(w WorldAccess, box *BoundingBox, minX, minY, minZ, maxX, maxY, maxZ int, outlineId, outlineMeta, insideId, insideMeta byte, keepOld bool) {
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			for z := minZ; z <= maxZ; z++ {
				if x != minX && x != maxX && z != minZ && z != maxZ && y != minY && y != maxY {
					c.SetBlockState(w, insideId, insideMeta, x, y, z, box)
				} else {
					c.SetBlockState(w, outlineId, outlineMeta, x, y, z, box)
				}
			}
		}
	}
}

func (c *StructureComponentBase) FillWithAir(w WorldAccess, box *BoundingBox, minX, minY, minZ, maxX, maxY, maxZ int) {
	c.FillWithBlocks(w, box, minX, minY, minZ, maxX, maxY, maxZ, 0, 0, 0, 0, false)
}

func (c *StructureComponentBase) ReplaceAirAndLiquidDownwards(w WorldAccess, id byte, meta byte, x, y, z int, box *BoundingBox) {
	worldX := c.GetXWithOffset(x, z)
	worldY := c.GetYWithOffset(y)
	worldZ := c.GetZWithOffset(x, z)

	if box.ResultIsInside(worldX, worldY, worldZ) {
		bId, _ := w.GetBlock(worldX, worldY, worldZ)

		if bId == 0 || (bId >= 8 && bId <= 11) {
			w.SetBlock(worldX, worldY, worldZ, id, meta)
		}
	}
}

func (c *StructureComponentBase) ClearCurrentPositionBlocksUpwards(w WorldAccess, x, y, z int, box *BoundingBox) {
	worldX := c.GetXWithOffset(x, z)
	worldY := c.GetYWithOffset(y)
	worldZ := c.GetZWithOffset(x, z)

	if box.ResultIsInside(worldX, worldY, worldZ) {

		for worldY < 128 {
			bId, _ := w.GetBlock(worldX, worldY, worldZ)
			if bId == 0 {
				break
			}
			w.SetBlock(worldX, worldY, worldZ, 0, 0)
			worldY++
		}
	}
}

type BlockSelector func(rnd *rand.Random, x, y, z int, wall bool) (byte, byte)

func (c *StructureComponentBase) FillWithRandomizedBlocks(w WorldAccess, box *BoundingBox, minX, minY, minZ, maxX, maxY, maxZ int, alwaysReplace bool, rnd *rand.Random, selector BlockSelector) {
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			for z := minZ; z <= maxZ; z++ {

				wall := (x == minX || x == maxX || z == minZ || z == maxZ || y == minY || y == maxY)
				id, meta := selector(rnd, x, y, z, wall)

				if id != 0 {
					c.SetBlockState(w, id, meta, x, y, z, box)
				}
			}
		}
	}
}
