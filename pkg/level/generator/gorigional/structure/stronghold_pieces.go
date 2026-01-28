package structure

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type StrongholdPiece struct {
	*StructureComponentBase
}

func NewStrongholdPiece(componentType int, rnd *rand.Random, box *BoundingBox, facing int) *StrongholdPiece {
	return &StrongholdPiece{
		StructureComponentBase: &StructureComponentBase{
			ComponentType: componentType,
			BoundingBox:   box,
			CoordBaseMode: facing,
		},
	}
}

type StrongholdPieceWeight struct {
	PieceClass func(*StructureStart, *rand.Random, *BoundingBox, int) StructureComponent
	Weight     int
	Limit      int
	Instances  int
}

var StrongholdWeights = []*StrongholdPieceWeight{
	{NewStrongholdStraight, 175, 0, 0},
	{NewStrongholdPortalRoom, 20, 1, 0},
}

func GetNextStrongholdComponent(start *StructureStart, components *[]StructureComponent, rnd *rand.Random, x, y, z int, facing int, typeId int) StructureComponent {
	if typeId > 50 {
		return nil
	}

	totalWeight := 0
	for _, pw := range StrongholdWeights {
		if pw.Limit == 0 || pw.Instances < pw.Limit {
			totalWeight += pw.Weight
		}
	}

	if totalWeight == 0 {
		return nil
	}

	i := rnd.NextBoundedInt(totalWeight)
	var selected *StrongholdPieceWeight

	for _, pw := range StrongholdWeights {
		if pw.Limit == 0 || pw.Instances < pw.Limit {
			i -= pw.Weight
			if i < 0 {
				selected = pw
				break
			}
		}
	}

	if selected != nil {

		fmt.Printf("GetNext: Selected %v at %d,%d,%d Facing %d\n", selected.Weight, x, y, z, facing)
		return CreateStrongholdPiece(selected, start, components, rnd, x, y, z, facing, typeId)
	}
	return nil
}

func CreateStrongholdPiece(pw *StrongholdPieceWeight, start *StructureStart, components *[]StructureComponent, rnd *rand.Random, x, y, z int, facing int, typeId int) StructureComponent {
	var box *BoundingBox

	if pw == StrongholdWeights[0] {
		box = GetStrongholdStraightBoundingBox(*components, rnd, x, y, z, facing)
		if box != nil {
			pw.Instances++
			fmt.Printf("Created Straight at %v\n", box)
			return NewStrongholdStraight(start, rnd, box, facing)
		} else {
			fmt.Printf("Straight Collision/Null Box\n")
		}
	} else if pw == StrongholdWeights[1] {
		box = GetStrongholdPortalRoomBoundingBox(*components, rnd, x, y, z, facing)
		if box != nil {
			pw.Instances++
			return NewStrongholdPortalRoom(start, rnd, box, facing)
		}
	}

	return nil
}

func GetStrongholdStraightBoundingBox(components []StructureComponent, rnd *rand.Random, x, y, z int, facing int) *BoundingBox {
	box := GetComponentToAddBoundingBox(x, y, z, -1, -1, 0, 5, 5, 7, facing)

	if c := FindIntersectingComponents(components, box); c != nil {
		fmt.Printf("Straight Intersection at %v with %v\n", box, c.GetBoundingBox())
		return nil
	}
	return box
}

func GetStrongholdPortalRoomBoundingBox(components []StructureComponent, rnd *rand.Random, x, y, z int, facing int) *BoundingBox {
	box := GetComponentToAddBoundingBox(x, y, z, -4, -1, 0, 11, 8, 16, facing)
	if FindIntersectingComponents(components, box) != nil {
		return nil
	}
	return box
}

func FindIntersectingComponents(list []StructureComponent, box *BoundingBox) StructureComponent {
	for _, c := range list {
		if c.GetBoundingBox() != nil && c.GetBoundingBox().IntersectsWith(box) {
			return c
		}
	}
	return nil
}

type StrongholdStraight struct {
	*StrongholdPiece
}

func NewStrongholdStraight(start *StructureStart, rnd *rand.Random, box *BoundingBox, facing int) StructureComponent {
	return &StrongholdStraight{
		StrongholdPiece: NewStrongholdPiece(1, rnd, box, facing),
	}
}

func (s *StrongholdStraight) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {

	if c := GetNextStrongholdComponent(nil, components, rnd, s.BoundingBox.MinX, s.BoundingBox.MinY, s.BoundingBox.MinZ, s.CoordBaseMode, s.ComponentType+1); c != nil {
		*components = append(*components, c)
		c.BuildComponent(s, components, rnd)
	}

}

func (s *StrongholdStraight) AddComponentParts(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	s.FillWithBlocks(wld, box, 0, 0, 0, 4, 4, 6, 98, 0, 98, 0, false)
	s.FillWithAir(wld, box, 1, 1, 0, 3, 3, 6)

	return true
}

type StrongholdStairs2 struct {
	*StrongholdPiece
}

func NewStrongholdStairs2(start *StructureStart, rnd *rand.Random, x, z int) *StrongholdStairs2 {

	facing := int(rnd.NextInt() % 4)

	width := 5
	height := 11
	depth := 5

	y := 64

	box := NewBoundingBox(x, y, z, x+width-1, y+height-1, z+depth-1)

	m := &StrongholdStairs2{
		StrongholdPiece: NewStrongholdPiece(0, rnd, box, facing),
	}
	return m
}

func (s *StrongholdStairs2) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {

	var c StructureComponent

	switch s.CoordBaseMode {
	case 0:

		if c = GetNextStrongholdComponent(nil, components, rnd, s.BoundingBox.MinX+2, s.BoundingBox.MinY, s.BoundingBox.MaxZ+1, 0, s.ComponentType+1); c != nil {
			*components = append(*components, c)
			c.BuildComponent(s, components, rnd)
		}

		if c = GetNextStrongholdComponent(nil, components, rnd, s.BoundingBox.MaxX+1, s.BoundingBox.MinY, s.BoundingBox.MinZ+2, 3, s.ComponentType+1); c != nil {
			*components = append(*components, c)
			c.BuildComponent(s, components, rnd)
		}

		if c = GetNextStrongholdComponent(nil, components, rnd, s.BoundingBox.MinX-1, s.BoundingBox.MinY, s.BoundingBox.MinZ+2, 1, s.ComponentType+1); c != nil {
			*components = append(*components, c)
			c.BuildComponent(s, components, rnd)
		}
	case 1:

		if c = GetNextStrongholdComponent(nil, components, rnd, s.BoundingBox.MinX-1, s.BoundingBox.MinY, s.BoundingBox.MinZ+2, 1, s.ComponentType+1); c != nil {
			*components = append(*components, c)
			c.BuildComponent(s, components, rnd)
		}

		if c = GetNextStrongholdComponent(nil, components, rnd, s.BoundingBox.MinX+2, s.BoundingBox.MinY, s.BoundingBox.MaxZ+1, 0, s.ComponentType+1); c != nil {
			*components = append(*components, c)
			c.BuildComponent(s, components, rnd)
		}

		if c = GetNextStrongholdComponent(nil, components, rnd, s.BoundingBox.MinX+2, s.BoundingBox.MinY, s.BoundingBox.MinZ-1, 2, s.ComponentType+1); c != nil {
			*components = append(*components, c)
			c.BuildComponent(s, components, rnd)
		}
	case 2:

		if c = GetNextStrongholdComponent(nil, components, rnd, s.BoundingBox.MinX+2, s.BoundingBox.MinY, s.BoundingBox.MinZ-1, 2, s.ComponentType+1); c != nil {
			*components = append(*components, c)
			c.BuildComponent(s, components, rnd)
		}

		if c = GetNextStrongholdComponent(nil, components, rnd, s.BoundingBox.MinX-1, s.BoundingBox.MinY, s.BoundingBox.MinZ+2, 1, s.ComponentType+1); c != nil {
			*components = append(*components, c)
			c.BuildComponent(s, components, rnd)
		}

		if c = GetNextStrongholdComponent(nil, components, rnd, s.BoundingBox.MaxX+1, s.BoundingBox.MinY, s.BoundingBox.MinZ+2, 3, s.ComponentType+1); c != nil {
			*components = append(*components, c)
			c.BuildComponent(s, components, rnd)
		}
	case 3:

		if c = GetNextStrongholdComponent(nil, components, rnd, s.BoundingBox.MaxX+1, s.BoundingBox.MinY, s.BoundingBox.MinZ+2, 3, s.ComponentType+1); c != nil {
			*components = append(*components, c)
			c.BuildComponent(s, components, rnd)
		}

		if c = GetNextStrongholdComponent(nil, components, rnd, s.BoundingBox.MinX+2, s.BoundingBox.MinY, s.BoundingBox.MinZ-1, 2, s.ComponentType+1); c != nil {
			*components = append(*components, c)
			c.BuildComponent(s, components, rnd)
		}

		if c = GetNextStrongholdComponent(nil, components, rnd, s.BoundingBox.MinX+2, s.BoundingBox.MinY, s.BoundingBox.MaxZ+1, 0, s.ComponentType+1); c != nil {
			*components = append(*components, c)
			c.BuildComponent(s, components, rnd)
		}
	}
}

func (s *StrongholdStairs2) AddComponentParts(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	s.FillWithBlocks(wld, box, 0, 0, 0, 4, 10, 4, 98, 0, 98, 0, false)
	s.FillWithAir(wld, box, 1, 1, 1, 3, 10, 3)

	return true
}

type StrongholdPortalRoom struct {
	*StrongholdPiece
	hasSpawner bool
}

func NewStrongholdPortalRoom(start *StructureStart, rnd *rand.Random, box *BoundingBox, facing int) StructureComponent {
	return &StrongholdPortalRoom{
		StrongholdPiece: NewStrongholdPiece(10, rnd, box, facing),
	}
}

func (s *StrongholdPortalRoom) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {

}

func (s *StrongholdPortalRoom) AddComponentParts(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	s.FillWithBlocks(wld, box, 0, 0, 0, 10, 7, 15, 98, 0, 98, 0, false)
	s.FillWithAir(wld, box, 1, 1, 1, 9, 6, 14)

	return true
}
