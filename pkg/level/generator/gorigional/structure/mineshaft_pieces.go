package structure

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type MineshaftPiece struct {
	*StructureComponentBase
}

func NewMineshaftPiece(componentType int, rnd *rand.Random, box *BoundingBox, facing int) *MineshaftPiece {
	return &MineshaftPiece{
		StructureComponentBase: &StructureComponentBase{
			ComponentType: componentType,
			BoundingBox:   box,
			CoordBaseMode: facing,
		},
	}
}

func CreateRandomShaftPiece(components *[]StructureComponent, rnd *rand.Random, x, y, z int, facing int, typeId int) StructureComponent {
	i := rnd.NextBoundedInt(100)
	if i >= 80 {
		box := GetMineshaftCrossingBoundingBox(*components, rnd, x, y, z, facing)
		if box != nil {
			return NewMineshaftCrossing(nil, rnd, box, facing)
		}
	} else if i >= 70 {
		box := GetMineshaftStairsBoundingBox(*components, rnd, x, y, z, facing)
		if box != nil {
			return NewMineshaftStairs(nil, rnd, box, facing)
		}
	} else {
		box := GetMineshaftCorridorBoundingBox(*components, rnd, x, y, z, facing)
		if box != nil {
			return NewMineshaftCorridor(nil, rnd, box, facing)
		}
	}
	return nil
}

func GetNextMineshaftComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random, x, y, z int, facing int, typeId int) StructureComponent {
	if typeId > 8 {
		return nil
	}
	newPiece := CreateRandomShaftPiece(components, rnd, x, y, z, facing, typeId+1)
	return newPiece
}

func GetMineshaftCorridorBoundingBox(components []StructureComponent, rnd *rand.Random, x, y, z int, facing int) *BoundingBox {
	length := 20
	box := GetComponentToAddBoundingBox(x, y, z, -1, 0, 0, 3, 3, length, facing)
	if FindIntersecting(components, box) != nil {
		return nil
	}
	return box
}

func GetMineshaftCrossingBoundingBox(components []StructureComponent, rnd *rand.Random, x, y, z int, facing int) *BoundingBox {
	box := GetComponentToAddBoundingBox(x, y, z, -2, 0, 0, 5, 3, 5, facing)
	if FindIntersecting(components, box) != nil {
		return nil
	}
	return box
}

func GetMineshaftStairsBoundingBox(components []StructureComponent, rnd *rand.Random, x, y, z int, facing int) *BoundingBox {
	box := GetComponentToAddBoundingBox(x, y, z, -1, 0, 0, 5, 5, 5, facing)
	if FindIntersecting(components, box) != nil {
		return nil
	}
	return box
}

func FindIntersecting(list []StructureComponent, box *BoundingBox) StructureComponent {
	for _, c := range list {
		if c.GetBoundingBox() != nil && c.GetBoundingBox().IntersectsWith(box) {
			return c
		}
	}
	return nil
}

type MineshaftRoom struct {
	*MineshaftPiece
}

func NewMineshaftRoom(start *StructureStart, rnd *rand.Random, x, z int) *MineshaftRoom {

	l := 10
	h := 4
	w := 10
	box := NewBoundingBox(x, 50, z, x+l-1, 50+h-1, z+w-1)

	m := &MineshaftRoom{
		MineshaftPiece: NewMineshaftPiece(0, rnd, box, 0),
	}
	return m
}

func (m *MineshaftRoom) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {

	if c := GetNextMineshaftComponent(m, components, rnd, m.BoundingBox.MinX+5, m.BoundingBox.MinY, m.BoundingBox.MinZ-1, 2, 0); c != nil {
		*components = append(*components, c)
		c.BuildComponent(m, components, rnd)
	}

	if c := GetNextMineshaftComponent(m, components, rnd, m.BoundingBox.MinX+5, m.BoundingBox.MinY, m.BoundingBox.MaxZ+1, 0, 0); c != nil {
		*components = append(*components, c)
		c.BuildComponent(m, components, rnd)
	}

	if c := GetNextMineshaftComponent(m, components, rnd, m.BoundingBox.MaxX+1, m.BoundingBox.MinY, m.BoundingBox.MinZ+5, 3, 0); c != nil {
		*components = append(*components, c)
		c.BuildComponent(m, components, rnd)
	}

	if c := GetNextMineshaftComponent(m, components, rnd, m.BoundingBox.MinX-1, m.BoundingBox.MinY, m.BoundingBox.MinZ+5, 1, 0); c != nil {
		*components = append(*components, c)
		c.BuildComponent(m, components, rnd)
	}
}

func (m *MineshaftRoom) AddComponentParts(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	m.FillWithBlocks(wld, box, 0, 0, 0, 9, 3, 9, 3, 0, 0, 0, false)

	return true
}

type MineshaftCorridor struct {
	*MineshaftPiece
	hasRails      bool
	hasSpiders    bool
	spawnerPlaced bool
}

func NewMineshaftCorridor(start *StructureStart, rnd *rand.Random, box *BoundingBox, facing int) *MineshaftCorridor {
	return &MineshaftCorridor{
		MineshaftPiece: NewMineshaftPiece(1, rnd, box, facing),
	}
}

func (m *MineshaftCorridor) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {

	facing := m.CoordBaseMode

	nextX := 0
	nextZ := 0

	switch facing {
	case 0:
		nextX = m.BoundingBox.MinX + 1
		nextZ = m.BoundingBox.MaxZ + 1
	case 1:
		nextX = m.BoundingBox.MinX - 1
		nextZ = m.BoundingBox.MinZ + 1
	case 2:
		nextX = m.BoundingBox.MinX + 1
		nextZ = m.BoundingBox.MinZ - 1
	case 3:
		nextX = m.BoundingBox.MaxX + 1
		nextZ = m.BoundingBox.MinZ + 1
	}

	if c := GetNextMineshaftComponent(m, components, rnd, nextX, m.BoundingBox.MinY, nextZ, facing, m.ComponentType+1); c != nil {
		*components = append(*components, c)
		c.BuildComponent(m, components, rnd)
	}
}

func (m *MineshaftCorridor) AddComponentParts(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	length := box.MaxX - box.MinX
	if length < 10 {
		length = box.MaxZ - box.MinZ
	}

	m.FillWithBlocks(wld, box, 0, 0, 0, 2, 2, length, 0, 0, 0, 0, false)

	for i := 0; i < length; i += 4 {

		m.FillWithBlocks(wld, box, 0, 0, i, 0, 2, i, 85, 0, 85, 0, false)
		m.FillWithBlocks(wld, box, 2, 0, i, 2, 2, i, 85, 0, 85, 0, false)

		m.FillWithBlocks(wld, box, 0, 2, i, 2, 2, i, 5, 0, 5, 0, false)
	}

	if rnd.NextFloat() < 0.5 {
		m.SetBlockState(wld, 66, 0, 1, 0, 0, box)
	}

	return true
}

type MineshaftCrossing struct {
	*MineshaftPiece
}

func NewMineshaftCrossing(start *StructureStart, rnd *rand.Random, box *BoundingBox, facing int) *MineshaftCrossing {
	return &MineshaftCrossing{
		MineshaftPiece: NewMineshaftPiece(2, rnd, box, facing),
	}
}

func (m *MineshaftCrossing) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {

}

func (m *MineshaftCrossing) AddComponentParts(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	m.FillWithAir(wld, box, 0, 0, 0, 4, 3, 4)
	m.FillWithBlocks(wld, box, 0, 0, 0, 4, 0, 4, 5, 0, 5, 0, false)
	return true
}

type MineshaftStairs struct {
	*MineshaftPiece
}

func NewMineshaftStairs(start *StructureStart, rnd *rand.Random, box *BoundingBox, facing int) *MineshaftStairs {
	return &MineshaftStairs{
		MineshaftPiece: NewMineshaftPiece(3, rnd, box, facing),
	}
}

func (m *MineshaftStairs) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {
}

func (m *MineshaftStairs) AddComponentParts(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	m.FillWithAir(wld, box, 0, 0, 0, 4, 4, 4)

	for i := 0; i < 5; i++ {
		m.SetBlockState(wld, 53, 0, 1, i, i, box)
	}
	return true
}

func NewMineshaftStart(worldSeed int64, chunkX, chunkZ int) *StructureStart {
	start := NewStructureStart(chunkX, chunkZ)

	rnd := rand.NewRandom(int64(chunkX)*341873128712 + int64(chunkZ)*132897987541 + worldSeed)

	room := NewMineshaftRoom(start, rnd, (chunkX<<4)+2, (chunkZ<<4)+2)
	start.Components = append(start.Components, room)

	var parent StructureComponent = nil
	room.BuildComponent(parent, &start.Components, rnd)
	start.UpdateBoundingBox()
	return start
}
