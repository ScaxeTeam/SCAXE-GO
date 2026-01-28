package structure

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type VillageStart struct {
	*StructureStart
	Valid bool
}

func NewVillageStart(worldSeed int64, rnd *rand.Random, chunkX, chunkZ, size int) *VillageStart {
	start := &VillageStart{
		StructureStart: NewStructureStart(chunkX, chunkZ),
	}

	list := GetStructureVillageWeightedPieceList(rnd, size)

	x := (chunkX * 16) + 2
	z := (chunkZ * 16) + 2

	startPiece := NewVillageStartPiece(nil, 0, rnd, x, z, list, size)

	start.Components = append(start.Components, startPiece)

	startPiece.BuildComponent(startPiece, &start.Components, rnd)

	for len(startPiece.PendingRoads) > 0 || len(startPiece.PendingHouses) > 0 {
		if len(startPiece.PendingRoads) == 0 {

			i := rnd.NextBoundedInt(len(startPiece.PendingHouses))
			comp := startPiece.PendingHouses[i]

			startPiece.PendingHouses = append(
				startPiece.PendingHouses[:i],
				startPiece.PendingHouses[i+1:]...,
			)
			comp.BuildComponent(startPiece, &start.Components, rnd)
		} else {

			j := rnd.NextBoundedInt(len(startPiece.PendingRoads))
			comp := startPiece.PendingRoads[j]

			startPiece.PendingRoads = append(
				startPiece.PendingRoads[:j],
				startPiece.PendingRoads[j+1:]...,
			)
			comp.BuildComponent(startPiece, &start.Components, rnd)
		}
	}

	start.UpdateBoundingBox()

	nonRoadCount := 0
	for _, comp := range start.Components {
		if comp.GetComponentType() != VillagePiecePath {
			nonRoadCount++
		}
	}
	start.Valid = nonRoadCount > 2

	return start
}

func (v *VillageStart) UpdateBoundingBox() {
	if len(v.Components) > 0 {

		first := v.Components[0].GetBoundingBox()
		if first != nil && v.BoundingBox == nil {
			v.BoundingBox = &BoundingBox{
				MinX: first.MinX,
				MinY: first.MinY,
				MinZ: first.MinZ,
				MaxX: first.MaxX,
				MaxY: first.MaxY,
				MaxZ: first.MaxZ,
			}
		}

		for _, comp := range v.Components {
			bb := comp.GetBoundingBox()
			if bb != nil && v.BoundingBox != nil {
				if bb.MinX < v.BoundingBox.MinX {
					v.BoundingBox.MinX = bb.MinX
				}
				if bb.MinY < v.BoundingBox.MinY {
					v.BoundingBox.MinY = bb.MinY
				}
				if bb.MinZ < v.BoundingBox.MinZ {
					v.BoundingBox.MinZ = bb.MinZ
				}
				if bb.MaxX > v.BoundingBox.MaxX {
					v.BoundingBox.MaxX = bb.MaxX
				}
				if bb.MaxY > v.BoundingBox.MaxY {
					v.BoundingBox.MaxY = bb.MaxY
				}
				if bb.MaxZ > v.BoundingBox.MaxZ {
					v.BoundingBox.MaxZ = bb.MaxZ
				}
			}
		}
	}
}
