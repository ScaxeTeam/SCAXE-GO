package structure

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type PieceWeight struct {
	PieceClass       int
	Weight           int
	Limit            int
	InstancesSpawned int
}

const (
	VillagePieceHouse4Garden = 0
	VillagePieceChurch       = 1
	VillagePieceLibrary      = 2
	VillagePieceWoodHut      = 3
	VillagePieceHall         = 4
	VillagePieceField1       = 5
	VillagePieceField2       = 6
	VillagePieceHouse1       = 7
	VillagePieceHouse2       = 8
	VillagePieceHouse3       = 9
	VillagePiecePath         = 10
	VillagePieceStart        = 11
	VillagePieceWell         = 12
)

func NewPieceWeight(pieceClass int, weight int, limit int) *PieceWeight {
	return &PieceWeight{
		PieceClass: pieceClass,
		Weight:     weight,
		Limit:      limit,
	}
}

func (pw *PieceWeight) CheckLimit(limit int) bool {
	return pw.Limit == 0 || pw.InstancesSpawned < pw.Limit
}

func GetStructureVillageWeightedPieceList(rnd *rand.Random, size int) []*PieceWeight {
	list := make([]*PieceWeight, 0)

	list = append(list, NewPieceWeight(VillagePieceHouse4Garden, 4, 2+size+rnd.NextBoundedInt(4+size*2-(2+size)+1)))
	list = append(list, NewPieceWeight(VillagePieceChurch, 20, 0+size+rnd.NextBoundedInt(1+size-(0+size)+1)))
	list = append(list, NewPieceWeight(VillagePieceLibrary, 20, 0+size+rnd.NextBoundedInt(2+size-(0+size)+1)))
	list = append(list, NewPieceWeight(VillagePieceWoodHut, 3, 2+size+rnd.NextBoundedInt(5+size*3-(2+size)+1)))
	list = append(list, NewPieceWeight(VillagePieceHall, 15, 0+size+rnd.NextBoundedInt(2+size-(0+size)+1)))
	list = append(list, NewPieceWeight(VillagePieceField1, 3, 1+size+rnd.NextBoundedInt(4+size-(1+size)+1)))
	list = append(list, NewPieceWeight(VillagePieceField2, 3, 2+size+rnd.NextBoundedInt(4+size*2-(2+size)+1)))
	list = append(list, NewPieceWeight(VillagePieceHouse1, 9, 0+size+rnd.NextBoundedInt(1+size-(0+size)+1)))
	list = append(list, NewPieceWeight(VillagePieceHouse2, 9, 0+size+rnd.NextBoundedInt(3+size*2-(0+size)+1)))
	list = append(list, NewPieceWeight(VillagePieceHouse3, 6, 2+size+rnd.NextBoundedInt(4+size*2-(2+size)+1)))

	return list
}

func GetTotalWeight(list []*PieceWeight) int {
	total := 0
	for _, pw := range list {
		if pw.CheckLimit(0) {
			total += pw.Weight
		}
	}
	return total
}

type VillagePiece struct {
	*StructureComponentBase
	AvgGroundLevel int
}

func (v *VillagePiece) GetAverageGroundLevel(wld WorldAccess, box *BoundingBox) int {

	sumY := 0
	count := 0

	for z := v.BoundingBox.MinZ; z <= v.BoundingBox.MaxZ; z++ {
		for x := v.BoundingBox.MinX; x <= v.BoundingBox.MaxX; x++ {

			if !box.ResultIsInside(x, 64, z) {
				continue
			}

			topY := getTopSolidOrLiquidBlock(wld, x, z)
			if topY < 62 {
				topY = 62
			}
			sumY += topY
			count++
		}
	}

	if count == 0 {
		return -1
	}
	return sumY / count
}

func getTopSolidOrLiquidBlock(wld WorldAccess, x, z int) int {

	for y := 127; y >= 0; y-- {
		id, _ := wld.GetBlock(x, y, z)

		if id == 0 {
			continue
		}

		if !blocksMovement(id) {
			continue
		}
		return y + 1
	}
	return 0
}

func blocksMovement(blockID byte) bool {
	switch blockID {
	case 0:
		return false
	case 18, 161:
		return false
	case 31, 32:
		return false
	case 37, 38:
		return false
	case 6:
		return false
	case 39, 40:
		return false
	case 50, 75, 76:
		return false
	case 51:
		return false
	case 55:
		return false
	case 59, 141, 142, 207:
		return false
	case 63, 68:
		return false
	case 66, 27, 28:
		return false
	case 78:
		return false
	case 83:
		return false
	case 104, 105:
		return false
	case 106:
		return false
	case 111:
		return true
	case 115:
		return false
	case 175:
		return false
	default:

		return true
	}
}

type VillageStartPiece struct {
	*VillageWell
	StructureVillageWeightedPieceList []*PieceWeight
	PendingRoads                      []StructureComponent
	PendingHouses                     []StructureComponent
	TerrainType                       int
	LastPlaced                        *PieceWeight
}

func NewVillageStartPiece(_ interface{}, _ int, rnd *rand.Random, x, z int, list []*PieceWeight, size int) *VillageStartPiece {
	dir := rnd.NextBoundedInt(4)

	well := NewVillageWell(0, rnd, x, z)
	well.CoordBaseMode = dir

	piece := &VillageStartPiece{
		VillageWell:                       well,
		StructureVillageWeightedPieceList: list,
		PendingRoads:                      make([]StructureComponent, 0),
		PendingHouses:                     make([]StructureComponent, 0),
		TerrainType:                       size,
	}
	piece.ComponentType = VillagePieceStart
	return piece
}

func (s *VillageStartPiece) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {

	bb := s.VillageWell.BoundingBox

	road1 := GenerateAndAddRoadPiece(s, components, rnd, bb.MinX-1, bb.MaxY-4+1, bb.MinZ+1, 1, VillagePiecePath)
	road2 := GenerateAndAddRoadPiece(s, components, rnd, bb.MaxX+1, bb.MaxY-4+1, bb.MinZ+1, 3, VillagePiecePath)
	road3 := GenerateAndAddRoadPiece(s, components, rnd, bb.MinX+1, bb.MaxY-4+1, bb.MinZ-1, 2, VillagePiecePath)
	road4 := GenerateAndAddRoadPiece(s, components, rnd, bb.MinX+1, bb.MaxY-4+1, bb.MaxZ+1, 0, VillagePiecePath)

	if road1 != nil {
		s.PendingRoads = append(s.PendingRoads, road1)
	}
	if road2 != nil {
		s.PendingRoads = append(s.PendingRoads, road2)
	}
	if road3 != nil {
		s.PendingRoads = append(s.PendingRoads, road3)
	}
	if road4 != nil {
		s.PendingRoads = append(s.PendingRoads, road4)
	}
}

type VillageWell struct {
	*VillagePiece
}

func NewVillageWell(typeInt int, rnd *rand.Random, x, z int) *VillageWell {
	bb := NewBoundingBox(x, 64, z, x+6-1, 78, z+6-1)
	return &VillageWell{
		VillagePiece: &VillagePiece{
			StructureComponentBase: &StructureComponentBase{
				BoundingBox:   bb,
				ComponentType: VillagePieceWell,
			},
			AvgGroundLevel: -1,
		},
	}
}

func (w *VillageWell) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {
}

func (w *VillageWell) AddComponentParts(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {
	if w.AvgGroundLevel < 0 {
		w.AvgGroundLevel = w.GetAverageGroundLevel(wld, box)
		if w.AvgGroundLevel < 0 {
			return true
		}

		w.BoundingBox.Offset(0, w.AvgGroundLevel-w.BoundingBox.MaxY+3, 0)
	}

	cobble := byte(4)
	fence := byte(85)

	w.FillWithBlocks(wld, box, 1, 0, 1, 4, 12, 4, cobble, 0, 8, 0, false)

	w.SetBlockState(wld, 0, 0, 2, 12, 2, box)
	w.SetBlockState(wld, 0, 0, 3, 12, 2, box)
	w.SetBlockState(wld, 0, 0, 2, 12, 3, box)
	w.SetBlockState(wld, 0, 0, 3, 12, 3, box)

	w.SetBlockState(wld, fence, 0, 1, 13, 1, box)
	w.SetBlockState(wld, fence, 0, 1, 14, 1, box)
	w.SetBlockState(wld, fence, 0, 4, 13, 1, box)
	w.SetBlockState(wld, fence, 0, 4, 14, 1, box)
	w.SetBlockState(wld, fence, 0, 1, 13, 4, box)
	w.SetBlockState(wld, fence, 0, 1, 14, 4, box)
	w.SetBlockState(wld, fence, 0, 4, 13, 4, box)
	w.SetBlockState(wld, fence, 0, 4, 14, 4, box)

	w.FillWithBlocks(wld, box, 1, 15, 1, 4, 15, 4, cobble, 0, cobble, 0, false)

	for i := 0; i <= 5; i++ {
		for j := 0; j <= 5; j++ {
			if j == 0 || j == 5 || i == 0 || i == 5 {
				w.SetBlockState(wld, cobble, 0, j, 11, i, box)

				w.SetBlockState(wld, 0, 0, j, 12, i, box)
			}
		}
	}

	return true
}

type VillagePath struct {
	*VillagePiece
	Length int
}

func NewVillagePath(start *VillageStartPiece, typeInt int, rnd *rand.Random, bb *BoundingBox, facing int) *VillagePath {
	lenX := bb.MaxX - bb.MinX + 1
	lenZ := bb.MaxZ - bb.MinZ + 1
	length := lenX
	if lenZ > lenX {
		length = lenZ
	}

	return &VillagePath{
		VillagePiece: &VillagePiece{
			StructureComponentBase: &StructureComponentBase{
				BoundingBox:   bb,
				CoordBaseMode: facing,
				ComponentType: VillagePiecePath,
			},
			AvgGroundLevel: -1,
		},
		Length: length,
	}
}

func (p *VillagePath) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {
	start, ok := component.(*VillageStartPiece)
	if !ok {
		return
	}

	bb := p.BoundingBox

	switch p.CoordBaseMode {
	case 0:

		FindAndCreateComponentFactory(start, start.StructureVillageWeightedPieceList, components, rnd,
			bb.MinX-1, bb.MinY, bb.MinZ+p.Length/2, 1, 1)

		FindAndCreateComponentFactory(start, start.StructureVillageWeightedPieceList, components, rnd,
			bb.MaxX+1, bb.MinY, bb.MinZ+p.Length/2, 3, 1)
	case 1:

		FindAndCreateComponentFactory(start, start.StructureVillageWeightedPieceList, components, rnd,
			bb.MinX+p.Length/2, bb.MinY, bb.MaxZ+1, 0, 1)

		FindAndCreateComponentFactory(start, start.StructureVillageWeightedPieceList, components, rnd,
			bb.MinX+p.Length/2, bb.MinY, bb.MinZ-1, 2, 1)
	case 2:

		FindAndCreateComponentFactory(start, start.StructureVillageWeightedPieceList, components, rnd,
			bb.MaxX+1, bb.MinY, bb.MaxZ-p.Length/2, 3, 1)

		FindAndCreateComponentFactory(start, start.StructureVillageWeightedPieceList, components, rnd,
			bb.MinX-1, bb.MinY, bb.MaxZ-p.Length/2, 1, 1)
	case 3:

		FindAndCreateComponentFactory(start, start.StructureVillageWeightedPieceList, components, rnd,
			bb.MaxX-p.Length/2, bb.MinY, bb.MinZ-1, 2, 1)

		FindAndCreateComponentFactory(start, start.StructureVillageWeightedPieceList, components, rnd,
			bb.MaxX-p.Length/2, bb.MinY, bb.MaxZ+1, 0, 1)
	}
}

func (p *VillagePath) AddComponentParts(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	const (
		grassPathId = byte(198)
		planksId    = byte(5)
		gravelId    = byte(13)
		cobbleId    = byte(4)
		grassId     = byte(2)
		sandId      = byte(12)
		sandstoneId = byte(24)
		redSandId   = byte(179)
	)

	seaLevel := 63

	for x := p.BoundingBox.MinX; x <= p.BoundingBox.MaxX; x++ {
		for z := p.BoundingBox.MinZ; z <= p.BoundingBox.MaxZ; z++ {

			if !box.ResultIsInside(x, 64, z) {
				continue
			}

			topY := getTopSolidOrLiquidBlock(wld, x, z) - 1

			if topY < seaLevel {
				topY = seaLevel - 1
			}

			for topY >= seaLevel-1 {
				id, _ := wld.GetBlock(x, topY, z)

				if id == grassId {
					aboveId, _ := wld.GetBlock(x, topY+1, z)
					if aboveId == 0 {
						wld.SetBlock(x, topY, z, grassPathId, 0)
						break
					}
				}

				if id >= 8 && id <= 11 {
					wld.SetBlock(x, topY, z, planksId, 0)
					break
				}

				if id == sandId || id == sandstoneId || id == redSandId {
					wld.SetBlock(x, topY, z, gravelId, 0)
					wld.SetBlock(x, topY-1, z, cobbleId, 0)
					break
				}

				topY--
			}
		}
	}
	return true
}

type VillageHouse struct {
	*VillagePiece
}

func NewVillageHouse(start *VillageStartPiece, typeInt int, rnd *rand.Random, bb *BoundingBox, facing int) *VillageHouse {
	return &VillageHouse{
		VillagePiece: &VillagePiece{
			StructureComponentBase: &StructureComponentBase{
				BoundingBox:   bb,
				CoordBaseMode: facing,
				ComponentType: typeInt,
			},
			AvgGroundLevel: -1,
		},
	}
}

func (h *VillageHouse) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {
}
func (h *VillageHouse) truncate(wld WorldAccess, rnd *rand.Random, box *BoundingBox) {

}

func (h *VillageHouse) AddComponentParts(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {
	if h.AvgGroundLevel < 0 {
		h.AvgGroundLevel = h.GetAverageGroundLevel(wld, box)
		if h.AvgGroundLevel < 0 {
			return true
		}

		height := h.getBuildingHeight()
		h.BoundingBox.Offset(0, h.AvgGroundLevel-h.BoundingBox.MaxY+height-1, 0)
	}

	switch h.ComponentType {
	case VillagePieceChurch:
		return h.addComponentPartsChurch(wld, rnd, box)
	case VillagePieceLibrary:
		return h.addComponentPartsLibrary(wld, rnd, box)
	case VillagePieceWoodHut:
		return h.addComponentPartsWoodHut(wld, rnd, box)
	case VillagePieceHall:
		return h.addComponentPartsHall(wld, rnd, box)
	case VillagePieceField1:
		return h.addComponentPartsField1(wld, rnd, box)
	case VillagePieceField2:
		return h.addComponentPartsField2(wld, rnd, box)
	case VillagePieceHouse1:
		return h.addComponentPartsHouse1(wld, rnd, box)
	case VillagePieceHouse2:
		return h.addComponentPartsHouse2(wld, rnd, box)
	case VillagePieceHouse3:
		return h.addComponentPartsHouse3(wld, rnd, box)
	case VillagePieceHouse4Garden:
		return h.addComponentPartsHouse4Garden(wld, rnd, box)
	default:

		return true
	}
}

func (h *VillageHouse) getBuildingHeight() int {
	switch h.ComponentType {
	case VillagePieceChurch:
		return 12
	case VillagePieceLibrary:
		return 9
	case VillagePieceWoodHut:
		return 6
	case VillagePieceHall:
		return 7
	case VillagePieceField1, VillagePieceField2:
		return 4
	case VillagePieceHouse1:
		return 9
	case VillagePieceHouse2:
		return 6
	case VillagePieceHouse3:
		return 7
	case VillagePieceHouse4Garden:
		return 6
	default:
		return 6
	}
}

func (h *VillageHouse) addComponentPartsLibrary(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	cobble := byte(4)
	planks := byte(5)
	oakStairs := byte(53)
	stoneStairs := byte(67)
	glassPane := byte(102)
	fence := byte(85)
	pressurePlate := byte(72)
	bookshelf := byte(47)
	craftingTable := byte(58)

	stairsNorth := byte(2)
	stairsSouth := byte(3)
	stairsEast := byte(1)

	h.FillWithBlocks(wld, box, 1, 1, 1, 7, 5, 4, 0, 0, 0, 0, false)

	h.FillWithBlocks(wld, box, 0, 0, 0, 8, 0, 5, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 0, 5, 0, 8, 5, 5, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 0, 6, 1, 8, 6, 4, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 0, 7, 2, 8, 7, 3, cobble, 0, cobble, 0, false)

	for i := -1; i <= 2; i++ {
		for j := 0; j <= 8; j++ {
			h.SetBlockState(wld, oakStairs, stairsNorth, j, 6+i, i, box)
			h.SetBlockState(wld, oakStairs, stairsSouth, j, 6+i, 5-i, box)
		}
	}

	h.FillWithBlocks(wld, box, 0, 1, 0, 0, 1, 5, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 1, 1, 5, 8, 1, 5, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 8, 1, 0, 8, 1, 4, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 2, 1, 0, 7, 1, 0, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 0, 2, 0, 0, 4, 0, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 0, 2, 5, 0, 4, 5, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 8, 2, 5, 8, 4, 5, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 8, 2, 0, 8, 4, 0, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 0, 2, 1, 0, 4, 4, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 1, 2, 5, 7, 4, 5, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 8, 2, 1, 8, 4, 4, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 1, 2, 0, 7, 4, 0, planks, 0, planks, 0, false)

	h.SetBlockState(wld, glassPane, 0, 4, 2, 0, box)
	h.SetBlockState(wld, glassPane, 0, 5, 2, 0, box)
	h.SetBlockState(wld, glassPane, 0, 6, 2, 0, box)
	h.SetBlockState(wld, glassPane, 0, 4, 3, 0, box)
	h.SetBlockState(wld, glassPane, 0, 5, 3, 0, box)
	h.SetBlockState(wld, glassPane, 0, 6, 3, 0, box)
	h.SetBlockState(wld, glassPane, 0, 0, 2, 2, box)
	h.SetBlockState(wld, glassPane, 0, 0, 2, 3, box)
	h.SetBlockState(wld, glassPane, 0, 0, 3, 2, box)
	h.SetBlockState(wld, glassPane, 0, 0, 3, 3, box)
	h.SetBlockState(wld, glassPane, 0, 8, 2, 2, box)
	h.SetBlockState(wld, glassPane, 0, 8, 2, 3, box)
	h.SetBlockState(wld, glassPane, 0, 8, 3, 2, box)
	h.SetBlockState(wld, glassPane, 0, 8, 3, 3, box)
	h.SetBlockState(wld, glassPane, 0, 2, 2, 5, box)
	h.SetBlockState(wld, glassPane, 0, 3, 2, 5, box)
	h.SetBlockState(wld, glassPane, 0, 5, 2, 5, box)
	h.SetBlockState(wld, glassPane, 0, 6, 2, 5, box)

	h.FillWithBlocks(wld, box, 1, 4, 1, 7, 4, 1, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 1, 4, 4, 7, 4, 4, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 1, 3, 4, 7, 3, 4, bookshelf, 0, bookshelf, 0, false)

	h.SetBlockState(wld, planks, 0, 7, 1, 4, box)
	h.SetBlockState(wld, oakStairs, stairsEast, 7, 1, 3, box)
	h.SetBlockState(wld, oakStairs, stairsNorth, 6, 1, 4, box)
	h.SetBlockState(wld, oakStairs, stairsNorth, 5, 1, 4, box)
	h.SetBlockState(wld, oakStairs, stairsNorth, 4, 1, 4, box)
	h.SetBlockState(wld, oakStairs, stairsNorth, 3, 1, 4, box)
	h.SetBlockState(wld, fence, 0, 6, 1, 3, box)
	h.SetBlockState(wld, pressurePlate, 0, 6, 2, 3, box)
	h.SetBlockState(wld, fence, 0, 4, 1, 3, box)
	h.SetBlockState(wld, pressurePlate, 0, 4, 2, 3, box)
	h.SetBlockState(wld, craftingTable, 0, 7, 1, 1, box)

	h.SetBlockState(wld, 0, 0, 1, 1, 0, box)
	h.SetBlockState(wld, 0, 0, 1, 2, 0, box)
	h.SetBlockState(wld, 64, 1, 1, 1, 0, box)
	h.SetBlockState(wld, 64, 8, 1, 2, 0, box)

	h.SetBlockState(wld, stoneStairs, stairsNorth, 1, 0, -1, box)

	for l := 0; l < 6; l++ {
		for k := 0; k < 9; k++ {
			h.ClearCurrentPositionBlocksUpwards(wld, k, 9, l, box)
			h.ReplaceAirAndLiquidDownwards(wld, cobble, 0, k, -1, l, box)
		}
	}

	return true
}

func (h *VillageHouse) addComponentPartsField1(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	log := byte(17)
	farmland := byte(60)
	water := byte(9)
	dirt := byte(3)
	wheat := byte(59)
	carrots := byte(141)
	potatoes := byte(142)
	beetroot := byte(207)

	h.FillWithBlocks(wld, box, 0, 1, 0, 12, 4, 8, 0, 0, 0, 0, false)

	h.FillWithBlocks(wld, box, 1, 0, 1, 2, 0, 7, farmland, 0, farmland, 0, false)
	h.FillWithBlocks(wld, box, 4, 0, 1, 5, 0, 7, farmland, 0, farmland, 0, false)
	h.FillWithBlocks(wld, box, 7, 0, 1, 8, 0, 7, farmland, 0, farmland, 0, false)
	h.FillWithBlocks(wld, box, 10, 0, 1, 11, 0, 7, farmland, 0, farmland, 0, false)

	h.FillWithBlocks(wld, box, 0, 0, 0, 0, 0, 8, log, 0, log, 0, false)
	h.FillWithBlocks(wld, box, 6, 0, 0, 6, 0, 8, log, 0, log, 0, false)
	h.FillWithBlocks(wld, box, 12, 0, 0, 12, 0, 8, log, 0, log, 0, false)
	h.FillWithBlocks(wld, box, 1, 0, 0, 11, 0, 0, log, 0, log, 0, false)
	h.FillWithBlocks(wld, box, 1, 0, 8, 11, 0, 8, log, 0, log, 0, false)

	h.FillWithBlocks(wld, box, 3, 0, 1, 3, 0, 7, water, 0, water, 0, false)
	h.FillWithBlocks(wld, box, 9, 0, 1, 9, 0, 7, water, 0, water, 0, false)

	cropTypes := []byte{wheat, carrots, potatoes, beetroot}
	cropA := cropTypes[rnd.NextBoundedInt(4)]
	cropB := cropTypes[rnd.NextBoundedInt(4)]
	cropC := cropTypes[rnd.NextBoundedInt(4)]
	cropD := cropTypes[rnd.NextBoundedInt(4)]

	for i := 1; i <= 7; i++ {

		metaA := byte(rnd.NextBoundedInt(8))
		h.SetBlockState(wld, cropA, metaA, 1, 1, i, box)
		h.SetBlockState(wld, cropA, metaA, 2, 1, i, box)

		metaB := byte(rnd.NextBoundedInt(8))
		h.SetBlockState(wld, cropB, metaB, 4, 1, i, box)
		h.SetBlockState(wld, cropB, metaB, 5, 1, i, box)

		metaC := byte(rnd.NextBoundedInt(8))
		h.SetBlockState(wld, cropC, metaC, 7, 1, i, box)
		h.SetBlockState(wld, cropC, metaC, 8, 1, i, box)

		metaD := byte(rnd.NextBoundedInt(8))
		h.SetBlockState(wld, cropD, metaD, 10, 1, i, box)
		h.SetBlockState(wld, cropD, metaD, 11, 1, i, box)
	}

	for j2 := 0; j2 < 9; j2++ {
		for k2 := 0; k2 < 13; k2++ {
			h.ClearCurrentPositionBlocksUpwards(wld, k2, 4, j2, box)
			h.ReplaceAirAndLiquidDownwards(wld, dirt, 0, k2, -1, j2, box)
		}
	}

	return true
}

func (h *VillageHouse) addComponentPartsHouse1(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	h.FillWithAir(wld, box, 1, 1, 1, 7, 4, 4)
	h.FillWithAir(wld, box, 2, 1, 6, 8, 4, 10)
	h.FillWithBlocks(wld, box, 2, 0, 6, 8, 0, 10, 3, 0, 3, 0, false)

	h.FillWithBlocks(wld, box, 0, 0, 0, 8, 4, 8, 4, 0, 0, 0, false)
	return true
}

func (h *VillageHouse) addComponentPartsHouse2(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	cobble := byte(4)
	planks := byte(5)
	log := byte(17)
	fence := byte(85)
	stoneSlab := byte(44)
	doubleStoneSlab := byte(43)
	glassPane := byte(102)
	ironBars := byte(101)
	furnace := byte(61)
	flowingLava := byte(11)
	pressurePlate := byte(72)
	oakStairsNorth := byte(53)
	oakStairsWest := byte(53)
	stoneStairsNorth := byte(67)

	stairsNorthMeta := byte(2)
	stairsWestMeta := byte(0)

	h.FillWithBlocks(wld, box, 0, 1, 0, 9, 4, 6, 0, 0, 0, 0, false)

	h.FillWithBlocks(wld, box, 0, 0, 0, 9, 0, 6, cobble, 0, cobble, 0, false)

	h.FillWithBlocks(wld, box, 0, 4, 0, 9, 4, 6, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 0, 5, 0, 9, 5, 6, stoneSlab, 0, stoneSlab, 0, false)
	h.FillWithBlocks(wld, box, 1, 5, 1, 8, 5, 5, 0, 0, 0, 0, false)

	h.FillWithBlocks(wld, box, 1, 1, 0, 2, 3, 0, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 0, 1, 0, 0, 4, 0, log, 0, log, 0, false)
	h.FillWithBlocks(wld, box, 3, 1, 0, 3, 4, 0, log, 0, log, 0, false)
	h.FillWithBlocks(wld, box, 0, 1, 6, 0, 4, 6, log, 0, log, 0, false)
	h.SetBlockState(wld, planks, 0, 3, 3, 1, box)
	h.FillWithBlocks(wld, box, 3, 1, 2, 3, 3, 2, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 4, 1, 3, 5, 3, 3, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 0, 1, 1, 0, 3, 5, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 1, 1, 6, 5, 3, 6, planks, 0, planks, 0, false)

	h.FillWithBlocks(wld, box, 5, 1, 0, 5, 3, 0, fence, 0, fence, 0, false)
	h.FillWithBlocks(wld, box, 9, 1, 0, 9, 3, 0, fence, 0, fence, 0, false)

	h.FillWithBlocks(wld, box, 6, 1, 4, 9, 4, 6, cobble, 0, cobble, 0, false)

	h.SetBlockState(wld, flowingLava, 0, 7, 1, 5, box)
	h.SetBlockState(wld, flowingLava, 0, 8, 1, 5, box)

	h.SetBlockState(wld, ironBars, 0, 9, 2, 5, box)
	h.SetBlockState(wld, ironBars, 0, 9, 2, 4, box)

	h.FillWithBlocks(wld, box, 7, 2, 4, 8, 2, 5, 0, 0, 0, 0, false)

	h.SetBlockState(wld, cobble, 0, 6, 1, 3, box)
	h.SetBlockState(wld, furnace, 2, 6, 2, 3, box)
	h.SetBlockState(wld, furnace, 2, 6, 3, 3, box)

	h.SetBlockState(wld, doubleStoneSlab, 0, 8, 1, 1, box)

	h.SetBlockState(wld, glassPane, 0, 0, 2, 2, box)
	h.SetBlockState(wld, glassPane, 0, 0, 2, 4, box)
	h.SetBlockState(wld, glassPane, 0, 2, 2, 6, box)
	h.SetBlockState(wld, glassPane, 0, 4, 2, 6, box)

	h.SetBlockState(wld, fence, 0, 2, 1, 4, box)
	h.SetBlockState(wld, pressurePlate, 0, 2, 2, 4, box)
	h.SetBlockState(wld, planks, 0, 1, 1, 5, box)

	h.SetBlockState(wld, oakStairsNorth, stairsNorthMeta, 2, 1, 5, box)
	h.SetBlockState(wld, oakStairsWest, stairsWestMeta, 1, 1, 4, box)

	h.SetBlockState(wld, 54, 3, 5, 1, 5, box)

	for i := 6; i <= 8; i++ {
		h.SetBlockState(wld, stoneStairsNorth, stairsNorthMeta, i, 0, -1, box)
	}

	h.SetBlockState(wld, 0, 0, 1, 1, 0, box)
	h.SetBlockState(wld, 0, 0, 1, 2, 0, box)
	h.SetBlockState(wld, 64, 1, 1, 1, 0, box)
	h.SetBlockState(wld, 64, 8, 1, 2, 0, box)

	for l := 0; l < 6; l++ {
		for k := 0; k < 9; k++ {
			h.ClearCurrentPositionBlocksUpwards(wld, k, 9, l, box)
			h.ReplaceAirAndLiquidDownwards(wld, cobble, 0, k, -1, l, box)
		}
	}

	return true
}

func (h *VillageHouse) addComponentPartsHouse3(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	cobble := byte(4)
	planks := byte(5)
	log := byte(17)
	oakStairs := byte(53)
	glassPane := byte(102)
	torch := byte(50)

	stairsNorth := byte(2)
	stairsSouth := byte(3)
	stairsEast := byte(1)
	stairsWest := byte(0)

	h.FillWithBlocks(wld, box, 1, 1, 1, 7, 4, 4, 0, 0, 0, 0, false)
	h.FillWithBlocks(wld, box, 2, 1, 6, 8, 4, 10, 0, 0, 0, 0, false)

	h.FillWithBlocks(wld, box, 2, 0, 5, 8, 0, 10, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 1, 0, 1, 7, 0, 4, planks, 0, planks, 0, false)

	h.FillWithBlocks(wld, box, 0, 0, 0, 0, 3, 5, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 8, 0, 0, 8, 3, 10, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 1, 0, 0, 7, 2, 0, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 1, 0, 5, 2, 1, 5, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 2, 0, 6, 2, 3, 10, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 3, 0, 10, 7, 3, 10, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 1, 2, 0, 7, 3, 0, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 1, 2, 5, 2, 3, 5, planks, 0, planks, 0, false)

	h.FillWithBlocks(wld, box, 0, 4, 1, 8, 4, 1, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 0, 4, 4, 3, 4, 4, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 0, 5, 2, 8, 5, 3, planks, 0, planks, 0, false)
	h.SetBlockState(wld, planks, 0, 0, 4, 2, box)
	h.SetBlockState(wld, planks, 0, 0, 4, 3, box)
	h.SetBlockState(wld, planks, 0, 8, 4, 2, box)
	h.SetBlockState(wld, planks, 0, 8, 4, 3, box)
	h.SetBlockState(wld, planks, 0, 8, 4, 4, box)

	for i := -1; i <= 2; i++ {
		for j := 0; j <= 8; j++ {
			h.SetBlockState(wld, oakStairs, stairsNorth, j, 4+i, i, box)
			if (i > -1 || j <= 1) && (i > 0 || j <= 3) && (i > 1 || j <= 4 || j >= 6) {
				h.SetBlockState(wld, oakStairs, stairsSouth, j, 4+i, 5-i, box)
			}
		}
	}

	h.FillWithBlocks(wld, box, 3, 4, 5, 3, 4, 10, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 7, 4, 2, 7, 4, 10, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 4, 5, 4, 4, 5, 10, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 6, 5, 4, 6, 5, 10, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 5, 6, 3, 5, 6, 10, planks, 0, planks, 0, false)

	for k := 4; k >= 1; k-- {
		h.SetBlockState(wld, planks, 0, k, 2+k, 7-k, box)
		for k1 := 8 - k; k1 <= 10; k1++ {
			h.SetBlockState(wld, oakStairs, stairsEast, k, 2+k, k1, box)
		}
	}
	h.SetBlockState(wld, planks, 0, 6, 6, 3, box)
	h.SetBlockState(wld, planks, 0, 7, 5, 4, box)
	h.SetBlockState(wld, oakStairs, stairsWest, 6, 6, 4, box)

	for l := 6; l <= 8; l++ {
		for l1 := 5; l1 <= 10; l1++ {
			h.SetBlockState(wld, oakStairs, stairsWest, l, 12-l, l1, box)
		}
	}

	h.SetBlockState(wld, log, 0, 0, 2, 1, box)
	h.SetBlockState(wld, log, 0, 0, 2, 4, box)
	h.SetBlockState(wld, glassPane, 0, 0, 2, 2, box)
	h.SetBlockState(wld, glassPane, 0, 0, 2, 3, box)
	h.SetBlockState(wld, log, 0, 4, 2, 0, box)
	h.SetBlockState(wld, glassPane, 0, 5, 2, 0, box)
	h.SetBlockState(wld, log, 0, 6, 2, 0, box)
	h.SetBlockState(wld, log, 0, 8, 2, 1, box)
	h.SetBlockState(wld, glassPane, 0, 8, 2, 2, box)
	h.SetBlockState(wld, glassPane, 0, 8, 2, 3, box)
	h.SetBlockState(wld, log, 0, 8, 2, 4, box)
	h.SetBlockState(wld, planks, 0, 8, 2, 5, box)
	h.SetBlockState(wld, log, 0, 8, 2, 6, box)
	h.SetBlockState(wld, glassPane, 0, 8, 2, 7, box)
	h.SetBlockState(wld, glassPane, 0, 8, 2, 8, box)
	h.SetBlockState(wld, log, 0, 8, 2, 9, box)
	h.SetBlockState(wld, log, 0, 2, 2, 6, box)
	h.SetBlockState(wld, glassPane, 0, 2, 2, 7, box)
	h.SetBlockState(wld, glassPane, 0, 2, 2, 8, box)
	h.SetBlockState(wld, log, 0, 2, 2, 9, box)
	h.SetBlockState(wld, log, 0, 4, 4, 10, box)
	h.SetBlockState(wld, glassPane, 0, 5, 4, 10, box)
	h.SetBlockState(wld, log, 0, 6, 4, 10, box)
	h.SetBlockState(wld, planks, 0, 5, 5, 10, box)

	h.SetBlockState(wld, 0, 0, 2, 1, 0, box)
	h.SetBlockState(wld, 0, 0, 2, 2, 0, box)
	h.SetBlockState(wld, torch, 4, 2, 3, 1, box)
	h.SetBlockState(wld, 64, 1, 2, 1, 0, box)
	h.SetBlockState(wld, 64, 8, 2, 2, 0, box)

	for i1 := 0; i1 < 5; i1++ {
		for i2 := 0; i2 < 9; i2++ {
			h.ClearCurrentPositionBlocksUpwards(wld, i2, 7, i1, box)
			h.ReplaceAirAndLiquidDownwards(wld, cobble, 0, i2, -1, i1, box)
		}
	}
	for j1 := 5; j1 < 11; j1++ {
		for j2 := 2; j2 < 9; j2++ {
			h.ClearCurrentPositionBlocksUpwards(wld, j2, 7, j1, box)
			h.ReplaceAirAndLiquidDownwards(wld, cobble, 0, j2, -1, j1, box)
		}
	}

	return true
}

func (h *VillageHouse) addComponentPartsHouse4Garden(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	cobble := byte(4)
	planks := byte(5)
	log := byte(17)
	stoneStairs := byte(67)
	glassPane := byte(102)
	fence := byte(85)
	ladder := byte(65)
	torch := byte(50)

	stairsNorth := byte(2)

	isRoofAccessible := rnd.NextBoundedInt(2) == 0

	h.FillWithBlocks(wld, box, 0, 0, 0, 4, 0, 4, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 0, 4, 0, 4, 4, 4, log, 0, log, 0, false)
	h.FillWithBlocks(wld, box, 1, 4, 1, 3, 4, 3, planks, 0, planks, 0, false)

	h.SetBlockState(wld, cobble, 0, 0, 1, 0, box)
	h.SetBlockState(wld, cobble, 0, 0, 2, 0, box)
	h.SetBlockState(wld, cobble, 0, 0, 3, 0, box)
	h.SetBlockState(wld, cobble, 0, 4, 1, 0, box)
	h.SetBlockState(wld, cobble, 0, 4, 2, 0, box)
	h.SetBlockState(wld, cobble, 0, 4, 3, 0, box)
	h.SetBlockState(wld, cobble, 0, 0, 1, 4, box)
	h.SetBlockState(wld, cobble, 0, 0, 2, 4, box)
	h.SetBlockState(wld, cobble, 0, 0, 3, 4, box)
	h.SetBlockState(wld, cobble, 0, 4, 1, 4, box)
	h.SetBlockState(wld, cobble, 0, 4, 2, 4, box)
	h.SetBlockState(wld, cobble, 0, 4, 3, 4, box)

	h.FillWithBlocks(wld, box, 0, 1, 1, 0, 3, 3, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 4, 1, 1, 4, 3, 3, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 1, 1, 4, 3, 3, 4, planks, 0, planks, 0, false)

	h.SetBlockState(wld, glassPane, 0, 0, 2, 2, box)
	h.SetBlockState(wld, glassPane, 0, 2, 2, 4, box)
	h.SetBlockState(wld, glassPane, 0, 4, 2, 2, box)

	h.SetBlockState(wld, planks, 0, 1, 1, 0, box)
	h.SetBlockState(wld, planks, 0, 1, 2, 0, box)
	h.SetBlockState(wld, planks, 0, 1, 3, 0, box)
	h.SetBlockState(wld, planks, 0, 2, 3, 0, box)
	h.SetBlockState(wld, planks, 0, 3, 3, 0, box)
	h.SetBlockState(wld, planks, 0, 3, 2, 0, box)
	h.SetBlockState(wld, planks, 0, 3, 1, 0, box)

	h.FillWithBlocks(wld, box, 1, 1, 1, 3, 3, 3, 0, 0, 0, 0, false)

	if isRoofAccessible {
		h.SetBlockState(wld, fence, 0, 0, 5, 0, box)
		h.SetBlockState(wld, fence, 0, 1, 5, 0, box)
		h.SetBlockState(wld, fence, 0, 2, 5, 0, box)
		h.SetBlockState(wld, fence, 0, 3, 5, 0, box)
		h.SetBlockState(wld, fence, 0, 4, 5, 0, box)
		h.SetBlockState(wld, fence, 0, 0, 5, 4, box)
		h.SetBlockState(wld, fence, 0, 1, 5, 4, box)
		h.SetBlockState(wld, fence, 0, 2, 5, 4, box)
		h.SetBlockState(wld, fence, 0, 3, 5, 4, box)
		h.SetBlockState(wld, fence, 0, 4, 5, 4, box)
		h.SetBlockState(wld, fence, 0, 4, 5, 1, box)
		h.SetBlockState(wld, fence, 0, 4, 5, 2, box)
		h.SetBlockState(wld, fence, 0, 4, 5, 3, box)
		h.SetBlockState(wld, fence, 0, 0, 5, 1, box)
		h.SetBlockState(wld, fence, 0, 0, 5, 2, box)
		h.SetBlockState(wld, fence, 0, 0, 5, 3, box)
	}

	if isRoofAccessible {

		h.SetBlockState(wld, ladder, 3, 3, 1, 3, box)
		h.SetBlockState(wld, ladder, 3, 3, 2, 3, box)
		h.SetBlockState(wld, ladder, 3, 3, 3, 3, box)
		h.SetBlockState(wld, ladder, 3, 3, 4, 3, box)
	}

	h.SetBlockState(wld, torch, 4, 2, 3, 1, box)

	h.SetBlockState(wld, stoneStairs, stairsNorth, 2, 0, -1, box)

	h.SetBlockState(wld, 0, 0, 2, 1, 0, box)
	h.SetBlockState(wld, 0, 0, 2, 2, 0, box)
	h.SetBlockState(wld, 64, 1, 2, 1, 0, box)
	h.SetBlockState(wld, 64, 8, 2, 2, 0, box)

	for j := 0; j < 5; j++ {
		for i := 0; i < 5; i++ {
			h.ClearCurrentPositionBlocksUpwards(wld, i, 6, j, box)
			h.ReplaceAirAndLiquidDownwards(wld, cobble, 0, i, -1, j, box)
		}
	}

	return true
}

func (h *VillageHouse) addComponentPartsField2(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	log := byte(17)
	farmland := byte(60)
	water := byte(9)
	dirt := byte(3)
	wheat := byte(59)
	carrots := byte(141)
	potatoes := byte(142)
	beetroot := byte(207)

	h.FillWithBlocks(wld, box, 0, 1, 0, 6, 4, 8, 0, 0, 0, 0, false)

	h.FillWithBlocks(wld, box, 1, 0, 1, 2, 0, 7, farmland, 0, farmland, 0, false)
	h.FillWithBlocks(wld, box, 4, 0, 1, 5, 0, 7, farmland, 0, farmland, 0, false)

	h.FillWithBlocks(wld, box, 0, 0, 0, 0, 0, 8, log, 0, log, 0, false)
	h.FillWithBlocks(wld, box, 6, 0, 0, 6, 0, 8, log, 0, log, 0, false)
	h.FillWithBlocks(wld, box, 1, 0, 0, 5, 0, 0, log, 0, log, 0, false)
	h.FillWithBlocks(wld, box, 1, 0, 8, 5, 0, 8, log, 0, log, 0, false)

	h.FillWithBlocks(wld, box, 3, 0, 1, 3, 0, 7, water, 0, water, 0, false)

	cropTypes := []byte{wheat, carrots, potatoes, beetroot}
	cropA := cropTypes[rnd.NextBoundedInt(4)]
	cropB := cropTypes[rnd.NextBoundedInt(4)]

	for i := 1; i <= 7; i++ {

		metaA := byte(rnd.NextBoundedInt(8))
		h.SetBlockState(wld, cropA, metaA, 1, 1, i, box)
		h.SetBlockState(wld, cropA, metaA, 2, 1, i, box)

		metaB := byte(rnd.NextBoundedInt(8))
		h.SetBlockState(wld, cropB, metaB, 4, 1, i, box)
		h.SetBlockState(wld, cropB, metaB, 5, 1, i, box)
	}

	for j1 := 0; j1 < 9; j1++ {
		for k1 := 0; k1 < 7; k1++ {
			h.ClearCurrentPositionBlocksUpwards(wld, k1, 4, j1, box)
			h.ReplaceAirAndLiquidDownwards(wld, dirt, 0, k1, -1, j1, box)
		}
	}

	return true
}

func (h *VillageHouse) addComponentPartsHall(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	cobble := byte(4)
	planks := byte(5)
	log := byte(17)
	dirt := byte(3)
	fence := byte(85)
	glassPane := byte(102)
	pressurePlate := byte(72)
	oakStairs := byte(53)
	doubleStoneSlab := byte(43)
	torch := byte(50)

	stairsNorth := byte(2)
	stairsSouth := byte(3)
	stairsWest := byte(0)

	h.FillWithBlocks(wld, box, 1, 1, 1, 7, 4, 4, 0, 0, 0, 0, false)
	h.FillWithBlocks(wld, box, 2, 1, 6, 8, 4, 10, 0, 0, 0, 0, false)

	h.FillWithBlocks(wld, box, 2, 0, 6, 8, 0, 10, dirt, 0, dirt, 0, false)
	h.SetBlockState(wld, cobble, 0, 6, 0, 6, box)

	h.FillWithBlocks(wld, box, 2, 1, 6, 2, 1, 10, fence, 0, fence, 0, false)
	h.FillWithBlocks(wld, box, 8, 1, 6, 8, 1, 10, fence, 0, fence, 0, false)
	h.FillWithBlocks(wld, box, 3, 1, 10, 7, 1, 10, fence, 0, fence, 0, false)

	h.FillWithBlocks(wld, box, 1, 0, 1, 7, 0, 4, planks, 0, planks, 0, false)

	h.FillWithBlocks(wld, box, 0, 0, 0, 0, 3, 5, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 8, 0, 0, 8, 3, 5, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 1, 0, 0, 7, 1, 0, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 1, 0, 5, 7, 1, 5, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 1, 2, 0, 7, 3, 0, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 1, 2, 5, 7, 3, 5, planks, 0, planks, 0, false)

	h.FillWithBlocks(wld, box, 0, 4, 1, 8, 4, 1, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 0, 4, 4, 8, 4, 4, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 0, 5, 2, 8, 5, 3, planks, 0, planks, 0, false)
	h.SetBlockState(wld, planks, 0, 0, 4, 2, box)
	h.SetBlockState(wld, planks, 0, 0, 4, 3, box)
	h.SetBlockState(wld, planks, 0, 8, 4, 2, box)
	h.SetBlockState(wld, planks, 0, 8, 4, 3, box)

	for i := -1; i <= 2; i++ {
		for j := 0; j <= 8; j++ {
			h.SetBlockState(wld, oakStairs, stairsNorth, j, 4+i, i, box)
			h.SetBlockState(wld, oakStairs, stairsSouth, j, 4+i, 5-i, box)
		}
	}

	h.SetBlockState(wld, log, 0, 0, 2, 1, box)
	h.SetBlockState(wld, log, 0, 0, 2, 4, box)
	h.SetBlockState(wld, log, 0, 8, 2, 1, box)
	h.SetBlockState(wld, log, 0, 8, 2, 4, box)

	h.SetBlockState(wld, glassPane, 0, 0, 2, 2, box)
	h.SetBlockState(wld, glassPane, 0, 0, 2, 3, box)
	h.SetBlockState(wld, glassPane, 0, 8, 2, 2, box)
	h.SetBlockState(wld, glassPane, 0, 8, 2, 3, box)
	h.SetBlockState(wld, glassPane, 0, 2, 2, 5, box)
	h.SetBlockState(wld, glassPane, 0, 3, 2, 5, box)
	h.SetBlockState(wld, glassPane, 0, 5, 2, 0, box)
	h.SetBlockState(wld, glassPane, 0, 6, 2, 5, box)

	h.SetBlockState(wld, fence, 0, 2, 1, 3, box)
	h.SetBlockState(wld, pressurePlate, 0, 2, 2, 3, box)
	h.SetBlockState(wld, planks, 0, 1, 1, 4, box)
	h.SetBlockState(wld, oakStairs, stairsNorth, 2, 1, 4, box)
	h.SetBlockState(wld, oakStairs, stairsWest, 1, 1, 3, box)
	h.FillWithBlocks(wld, box, 5, 0, 1, 7, 0, 3, doubleStoneSlab, 0, doubleStoneSlab, 0, false)
	h.SetBlockState(wld, doubleStoneSlab, 0, 6, 1, 1, box)
	h.SetBlockState(wld, doubleStoneSlab, 0, 6, 1, 2, box)

	h.SetBlockState(wld, 0, 0, 2, 1, 0, box)
	h.SetBlockState(wld, 0, 0, 2, 2, 0, box)
	h.SetBlockState(wld, 64, 1, 2, 1, 0, box)
	h.SetBlockState(wld, 64, 8, 2, 2, 0, box)

	h.SetBlockState(wld, torch, 4, 2, 3, 1, box)

	h.SetBlockState(wld, 0, 0, 6, 1, 5, box)
	h.SetBlockState(wld, 0, 0, 6, 2, 5, box)
	h.SetBlockState(wld, 64, 3, 6, 1, 5, box)
	h.SetBlockState(wld, 64, 8, 6, 2, 5, box)

	h.SetBlockState(wld, torch, 3, 6, 3, 4, box)

	for k := 0; k < 5; k++ {
		for l := 0; l < 9; l++ {
			h.ClearCurrentPositionBlocksUpwards(wld, l, 7, k, box)
			h.ReplaceAirAndLiquidDownwards(wld, cobble, 0, l, -1, k, box)
		}
	}

	return true
}

func (h *VillageHouse) addComponentPartsWoodHut(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	cobble := byte(4)
	planks := byte(5)
	log := byte(17)
	glassPane := byte(102)
	fence := byte(85)
	pressurePlate := byte(72)
	stoneStairs := byte(67)
	dirt := byte(3)

	h.FillWithBlocks(wld, box, 1, 1, 1, 3, 5, 4, 0, 0, 0, 0, false)

	h.FillWithBlocks(wld, box, 0, 0, 0, 3, 0, 4, cobble, 0, cobble, 0, false)

	h.FillWithBlocks(wld, box, 1, 0, 1, 2, 0, 3, dirt, 0, dirt, 0, false)

	h.FillWithBlocks(wld, box, 1, 5, 1, 2, 5, 3, log, 0, log, 0, false)

	h.SetBlockState(wld, log, 0, 1, 4, 0, box)
	h.SetBlockState(wld, log, 0, 2, 4, 0, box)
	h.SetBlockState(wld, log, 0, 1, 4, 4, box)
	h.SetBlockState(wld, log, 0, 2, 4, 4, box)
	h.SetBlockState(wld, log, 0, 0, 4, 1, box)
	h.SetBlockState(wld, log, 0, 0, 4, 2, box)
	h.SetBlockState(wld, log, 0, 0, 4, 3, box)
	h.SetBlockState(wld, log, 0, 3, 4, 1, box)
	h.SetBlockState(wld, log, 0, 3, 4, 2, box)
	h.SetBlockState(wld, log, 0, 3, 4, 3, box)

	h.FillWithBlocks(wld, box, 0, 1, 0, 0, 3, 0, log, 0, log, 0, false)
	h.FillWithBlocks(wld, box, 3, 1, 0, 3, 3, 0, log, 0, log, 0, false)
	h.FillWithBlocks(wld, box, 0, 1, 4, 0, 3, 4, log, 0, log, 0, false)
	h.FillWithBlocks(wld, box, 3, 1, 4, 3, 3, 4, log, 0, log, 0, false)

	h.FillWithBlocks(wld, box, 0, 1, 1, 0, 3, 3, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 3, 1, 1, 3, 3, 3, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 1, 1, 0, 2, 3, 0, planks, 0, planks, 0, false)
	h.FillWithBlocks(wld, box, 1, 1, 4, 2, 3, 4, planks, 0, planks, 0, false)

	h.SetBlockState(wld, glassPane, 0, 0, 2, 2, box)
	h.SetBlockState(wld, glassPane, 0, 3, 2, 2, box)

	tablePos := rnd.NextBoundedInt(3)
	if tablePos > 0 {
		h.SetBlockState(wld, fence, 0, tablePos, 1, 3, box)
		h.SetBlockState(wld, pressurePlate, 0, tablePos, 2, 3, box)
	}

	h.SetBlockState(wld, 0, 0, 1, 1, 0, box)
	h.SetBlockState(wld, 0, 0, 1, 2, 0, box)

	h.SetBlockState(wld, 64, 1, 1, 1, 0, box)
	h.SetBlockState(wld, 64, 8, 1, 2, 0, box)

	h.SetBlockState(wld, stoneStairs, 2, 1, 0, -1, box)

	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			h.ReplaceAirAndLiquidDownwards(wld, cobble, 0, j, -1, i, box)
		}
	}

	return true
}

func GenerateAndAddRoadPiece(start *VillageStartPiece, components *[]StructureComponent, rnd *rand.Random, x, y, z int, facing int, typeInt int) *VillagePath {

	l := 7 + rnd.NextBoundedInt(6)

	bb := GetComponentToAddBoundingBox(x, y, z, 0, 0, 0, 3, 3, l, facing)

	for _, comp := range *components {
		if comp.GetBoundingBox().IntersectsWith(bb) {
			return nil
		}
	}

	path := NewVillagePath(start, typeInt, rnd, bb, facing)
	*components = append(*components, path)

	path.BuildComponent(start, components, rnd)
	return path
}

func FindAndCreateComponentFactory(start *VillageStartPiece, weights []*PieceWeight, components *[]StructureComponent, rnd *rand.Random, x, y, z int, facing int, depth int) StructureComponent {

	if depth > 50 {
		return nil
	}

	if start.VillageWell != nil && start.VillageWell.BoundingBox != nil {
		startBB := start.VillageWell.BoundingBox
		if abs(x-startBB.MinX) > 112 || abs(z-startBB.MinZ) > 112 {
			return nil
		}
	}

	totalWeight := GetTotalWeight(weights)
	if totalWeight <= 0 {
		return nil
	}

	t := rnd.NextBoundedInt(totalWeight)

	var selected *PieceWeight
	for _, pw := range weights {
		if pw.CheckLimit(0) {
			t -= pw.Weight
			if t < 0 {
				selected = pw
				break
			}
		}
	}

	if selected == nil {
		return nil
	}

	w, h, l := GetVillagePieceSize(selected.PieceClass)
	bb := GetComponentToAddBoundingBox(x, y, z, 0, 0, 0, w, h, l, facing)

	for _, comp := range *components {
		if comp.GetBoundingBox().IntersectsWith(bb) {
			return nil
		}
	}

	selected.InstancesSpawned++
	piece := NewVillageHouse(start, selected.PieceClass, rnd, bb, facing)
	*components = append(*components, piece)

	start.PendingHouses = append(start.PendingHouses, piece)

	return piece
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (h *VillageHouse) addComponentPartsChurch(wld WorldAccess, rnd *rand.Random, box *BoundingBox) bool {

	cobble := byte(4)
	stoneStairs := byte(67)
	glassPane := byte(102)
	ladder := byte(65)
	torch := byte(50)

	stairsNorth := byte(2)
	stairsWest := byte(0)
	stairsEast := byte(1)

	h.FillWithBlocks(wld, box, 1, 1, 1, 3, 3, 7, 0, 0, 0, 0, false)
	h.FillWithBlocks(wld, box, 1, 5, 1, 3, 9, 3, 0, 0, 0, 0, false)

	h.FillWithBlocks(wld, box, 1, 0, 0, 3, 0, 8, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 1, 1, 0, 3, 10, 0, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 0, 1, 1, 0, 10, 3, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 4, 1, 1, 4, 10, 3, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 0, 0, 4, 0, 4, 7, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 4, 0, 4, 4, 4, 7, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 1, 1, 8, 3, 4, 8, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 1, 5, 4, 3, 10, 4, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 1, 5, 5, 3, 5, 7, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 0, 9, 0, 4, 9, 4, cobble, 0, cobble, 0, false)
	h.FillWithBlocks(wld, box, 0, 4, 0, 4, 4, 4, cobble, 0, cobble, 0, false)

	h.SetBlockState(wld, cobble, 0, 0, 11, 2, box)
	h.SetBlockState(wld, cobble, 0, 4, 11, 2, box)
	h.SetBlockState(wld, cobble, 0, 2, 11, 0, box)
	h.SetBlockState(wld, cobble, 0, 2, 11, 4, box)

	h.SetBlockState(wld, cobble, 0, 1, 1, 6, box)
	h.SetBlockState(wld, cobble, 0, 1, 1, 7, box)
	h.SetBlockState(wld, cobble, 0, 2, 1, 7, box)
	h.SetBlockState(wld, cobble, 0, 3, 1, 6, box)
	h.SetBlockState(wld, cobble, 0, 3, 1, 7, box)
	h.SetBlockState(wld, stoneStairs, stairsNorth, 1, 1, 5, box)
	h.SetBlockState(wld, stoneStairs, stairsNorth, 2, 1, 6, box)
	h.SetBlockState(wld, stoneStairs, stairsNorth, 3, 1, 5, box)
	h.SetBlockState(wld, stoneStairs, stairsWest, 1, 2, 7, box)
	h.SetBlockState(wld, stoneStairs, stairsEast, 3, 2, 7, box)

	h.SetBlockState(wld, glassPane, 0, 0, 2, 2, box)
	h.SetBlockState(wld, glassPane, 0, 0, 3, 2, box)
	h.SetBlockState(wld, glassPane, 0, 4, 2, 2, box)
	h.SetBlockState(wld, glassPane, 0, 4, 3, 2, box)
	h.SetBlockState(wld, glassPane, 0, 0, 6, 2, box)
	h.SetBlockState(wld, glassPane, 0, 0, 7, 2, box)
	h.SetBlockState(wld, glassPane, 0, 4, 6, 2, box)
	h.SetBlockState(wld, glassPane, 0, 4, 7, 2, box)
	h.SetBlockState(wld, glassPane, 0, 2, 6, 0, box)
	h.SetBlockState(wld, glassPane, 0, 2, 7, 0, box)
	h.SetBlockState(wld, glassPane, 0, 2, 6, 4, box)
	h.SetBlockState(wld, glassPane, 0, 2, 7, 4, box)
	h.SetBlockState(wld, glassPane, 0, 0, 3, 6, box)
	h.SetBlockState(wld, glassPane, 0, 4, 3, 6, box)
	h.SetBlockState(wld, glassPane, 0, 2, 3, 8, box)

	h.SetBlockState(wld, torch, 3, 2, 4, 7, box)
	h.SetBlockState(wld, torch, 1, 1, 4, 6, box)
	h.SetBlockState(wld, torch, 2, 3, 4, 6, box)
	h.SetBlockState(wld, torch, 4, 2, 4, 5, box)

	for i := 1; i <= 9; i++ {
		h.SetBlockState(wld, ladder, 4, 3, i, 3, box)
	}

	h.SetBlockState(wld, 0, 0, 2, 1, 0, box)
	h.SetBlockState(wld, 0, 0, 2, 2, 0, box)

	h.SetBlockState(wld, 64, 1, 2, 1, 0, box)
	h.SetBlockState(wld, 64, 8, 2, 2, 0, box)

	h.SetBlockState(wld, stoneStairs, stairsNorth, 2, 0, -1, box)

	for k := 0; k < 9; k++ {
		for j := 0; j < 5; j++ {
			h.ClearCurrentPositionBlocksUpwards(wld, j, 12, k, box)
			h.ReplaceAirAndLiquidDownwards(wld, cobble, 0, j, -1, k, box)
		}
	}

	return true
}

func GetVillagePieceSize(classID int) (int, int, int) {
	switch classID {
	case VillagePieceHouse4Garden:
		return 9, 7, 7
	case VillagePieceChurch:
		return 12, 12, 9
	case VillagePieceLibrary:
		return 9, 6, 7
	case VillagePieceWoodHut:
		return 4, 6, 5
	case VillagePieceHall:
		return 9, 7, 5
	case VillagePieceField1:
		return 13, 4, 9
	case VillagePieceField2:
		return 7, 4, 9
	case VillagePieceHouse1:
		return 9, 7, 9
	case VillagePieceHouse2:
		return 10, 6, 8
	case VillagePieceHouse3:
		return 5, 6, 5
	}
	return 3, 3, 3
}
