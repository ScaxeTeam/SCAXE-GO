package structure

import (
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type ScatteredFeaturePiece struct {
	*StructureComponentBase
	Width  int
	Height int
	Depth  int
	HPos   int
}

func NewScatteredFeaturePiece(componentType int, rnd *rand.Random, x, y, z, width, height, depth int, facing int) *ScatteredFeaturePiece {
	box := GetComponentToAddBoundingBox(x, y, z, 0, 0, 0, width, height, depth, facing)

	return &ScatteredFeaturePiece{
		StructureComponentBase: &StructureComponentBase{
			ComponentType: componentType,
			BoundingBox:   box,
			CoordBaseMode: facing,
		},
		Width:  width,
		Height: height,
		Depth:  depth,
		HPos:   -1,
	}
}

func (p *ScatteredFeaturePiece) OffsetToAverageGroundLevel(w WorldAccess, box *BoundingBox, yOffset int) bool {
	if p.HPos >= 0 {
		return true
	}

	totalHeight := 0
	count := 0

	for z := p.BoundingBox.MinZ; z <= p.BoundingBox.MaxZ; z++ {
		for x := p.BoundingBox.MinX; x <= p.BoundingBox.MaxX; x++ {
			if box.ResultIsInside(x, 64, z) {

				y := p.getTopSolidBlockY(w, x, z)

				if y < 63 {

					y = 63
				}
				totalHeight += y
				count++
			}
		}
	}

	if count == 0 {
		return false
	}

	avgY := totalHeight / count
	p.HPos = avgY

	p.BoundingBox.Offset(0, p.HPos-p.BoundingBox.MinY+yOffset, 0)
	return true
}

func (p *ScatteredFeaturePiece) getTopSolidBlockY(w WorldAccess, x, z int) int {

	for y := 255; y >= 0; y-- {
		id, _ := w.GetBlock(x, y, z)
		if id != 0 && id != 9 && id != 11 && id != 10 && id != 8 {

			return y
		}
	}
	return 0
}

type DesertPyramid struct {
	*ScatteredFeaturePiece
	hasPlacedChest [4]bool
}

func NewDesertPyramid(rnd *rand.Random, x, z int) *DesertPyramid {
	facing := rnd.NextBoundedInt(4)
	return &DesertPyramid{
		ScatteredFeaturePiece: NewScatteredFeaturePiece(0, rnd, x, 64, z, 21, 15, 21, facing),
	}
}

func (d *DesertPyramid) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {
}

func (d *DesertPyramid) AddComponentParts(w WorldAccess, rnd *rand.Random, box *BoundingBox) bool {
	if !d.OffsetToAverageGroundLevel(w, box, 0) {
		return false
	}

	SANDSTONE := byte(24)
	META_SMOOTH := byte(2)

	STAINED_CLAY := byte(159)
	META_ORANGE := byte(1)
	META_BLUE := byte(11)

	STONE_PRESSURE_PLATE := byte(70)
	TNT := byte(46)

	d.FillWithBlocks(w, box, 0, -4, 0, 20, 0, 20, SANDSTONE, 0, SANDSTONE, 0, false)

	for i := 1; i <= 9; i++ {
		d.FillWithBlocks(w, box, i, i, i, 20-i, i, 20-i, SANDSTONE, 0, SANDSTONE, 0, false)
		d.FillWithBlocks(w, box, i+1, i, i+1, 20-i-1, i, 20-i-1, 0, 0, 0, 0, false)
	}

	for i2 := 0; i2 < 21; i2++ {
		for j := 0; j < 21; j++ {
			d.ReplaceAirAndLiquidDownwards(w, SANDSTONE, 0, i2, -5, j, box)
		}
	}

	d.FillWithBlocks(w, box, 0, 0, 0, 4, 9, 4, SANDSTONE, 0, 0, 0, false)
	d.FillWithBlocks(w, box, 1, 10, 1, 3, 10, 3, SANDSTONE, 0, SANDSTONE, 0, false)

	d.FillWithBlocks(w, box, 16, 0, 0, 20, 9, 4, SANDSTONE, 0, 0, 0, false)
	d.FillWithBlocks(w, box, 17, 10, 1, 19, 10, 3, SANDSTONE, 0, SANDSTONE, 0, false)

	d.FillWithBlocks(w, box, 8, 0, 0, 12, 4, 4, SANDSTONE, 0, 0, 0, false)
	d.FillWithBlocks(w, box, 9, 1, 0, 11, 3, 4, 0, 0, 0, 0, false)

	d.SetBlockState(w, SANDSTONE, META_SMOOTH, 9, 1, 1, box)
	d.SetBlockState(w, SANDSTONE, META_SMOOTH, 9, 2, 1, box)
	d.SetBlockState(w, SANDSTONE, META_SMOOTH, 9, 3, 1, box)
	d.SetBlockState(w, SANDSTONE, META_SMOOTH, 10, 3, 1, box)
	d.SetBlockState(w, SANDSTONE, META_SMOOTH, 11, 3, 1, box)
	d.SetBlockState(w, SANDSTONE, META_SMOOTH, 11, 2, 1, box)
	d.SetBlockState(w, SANDSTONE, META_SMOOTH, 11, 1, 1, box)

	d.FillWithBlocks(w, box, 4, 1, 1, 8, 3, 3, SANDSTONE, 0, 0, 0, false)
	d.FillWithBlocks(w, box, 4, 1, 2, 8, 2, 2, 0, 0, 0, 0, false)
	d.FillWithBlocks(w, box, 12, 1, 1, 16, 3, 3, SANDSTONE, 0, 0, 0, false)
	d.FillWithBlocks(w, box, 12, 1, 2, 16, 2, 2, 0, 0, 0, 0, false)

	for j2 := 0; j2 < 21; j2 += 20 {

		d.SetBlockState(w, SANDSTONE, META_SMOOTH, j2, 2, 1, box)
		d.SetBlockState(w, STAINED_CLAY, META_ORANGE, j2, 2, 2, box)
		d.SetBlockState(w, SANDSTONE, META_SMOOTH, j2, 2, 3, box)
		d.SetBlockState(w, SANDSTONE, META_SMOOTH, j2, 3, 1, box)
		d.SetBlockState(w, STAINED_CLAY, META_ORANGE, j2, 3, 2, box)
		d.SetBlockState(w, SANDSTONE, META_SMOOTH, j2, 3, 3, box)

	}

	d.SetBlockState(w, STAINED_CLAY, META_BLUE, 10, 0, 10, box)

	d.SetBlockState(w, STAINED_CLAY, META_ORANGE, 10, 0, 9, box)
	d.SetBlockState(w, STAINED_CLAY, META_ORANGE, 10, 0, 11, box)
	d.SetBlockState(w, STAINED_CLAY, META_ORANGE, 9, 0, 10, box)
	d.SetBlockState(w, STAINED_CLAY, META_ORANGE, 11, 0, 10, box)

	d.SetBlockState(w, STONE_PRESSURE_PLATE, 0, 10, -11, 10, box)
	d.FillWithBlocks(w, box, 9, -13, 9, 11, -13, 11, TNT, 0, 0, 0, false)

	for i := 0; i < 4; i++ {
		if !d.hasPlacedChest[i] {

			cx := 10
			cz := 10
			switch i {
			case 0:
				cz -= 2
			case 1:
				cx += 2
			case 2:
				cz += 2
			case 3:
				cx -= 2
			}

			if box.ResultIsInside(d.GetXWithOffset(cx, cz), d.GetYWithOffset(-11), d.GetZWithOffset(cx, cz)) {
				d.SetBlockState(w, 54, 0, cx, -11, cz, box)

				d.hasPlacedChest[i] = true
			}
		}
	}

	return true
}

type JunglePyramid struct {
	*ScatteredFeaturePiece
	placedMainChest   bool
	placedHiddenChest bool
	placedTrap1       bool
	placedTrap2       bool
}

func NewJunglePyramid(rnd *rand.Random, x, z int) *JunglePyramid {
	facing := rnd.NextBoundedInt(4)
	return &JunglePyramid{
		ScatteredFeaturePiece: NewScatteredFeaturePiece(1, rnd, x, 64, z, 12, 10, 15, facing),
	}
}

func (j *JunglePyramid) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {
}

func (j *JunglePyramid) AddComponentParts(w WorldAccess, rnd *rand.Random, box *BoundingBox) bool {
	if !j.OffsetToAverageGroundLevel(w, box, 0) {
		return false
	}

	COBBLE := byte(4)
	MOSSY := byte(48)
	STONE_STAIRS := byte(67)
	VINE := byte(106)
	TRIPWIRE_HOOK := byte(131)
	TRIPWIRE := byte(132)
	REDSTONE := byte(55)
	LEVER := byte(69)
	STICKY_PISTON := byte(29)
	REPEATER := byte(93)
	DISPENSER := byte(23)
	STONEBRICK := byte(98)
	CHISELED := byte(3)

	cobblestoneSelector := func(rnd *rand.Random, x, y, z int, wall bool) (byte, byte) {
		if rnd.NextFloat() < 0.4 {
			return COBBLE, 0
		}
		return MOSSY, 0
	}

	j.FillWithRandomizedBlocks(w, box, 0, -4, 0, j.Width-1, 0, j.Depth-1, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 2, 1, 2, 9, 2, 2, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 2, 1, 12, 9, 2, 12, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 2, 1, 3, 2, 2, 11, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 9, 1, 3, 9, 2, 11, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 1, 3, 1, 10, 6, 1, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 1, 3, 13, 10, 6, 13, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 1, 3, 2, 1, 6, 12, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 10, 3, 2, 10, 6, 12, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 2, 3, 2, 9, 3, 12, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 2, 6, 2, 9, 6, 12, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 3, 7, 3, 8, 7, 11, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 4, 8, 4, 7, 8, 10, false, rnd, cobblestoneSelector)

	j.FillWithAir(w, box, 3, 1, 3, 8, 2, 11)
	j.FillWithAir(w, box, 4, 3, 6, 7, 3, 9)
	j.FillWithAir(w, box, 2, 4, 2, 9, 5, 12)
	j.FillWithAir(w, box, 4, 6, 5, 7, 6, 9)
	j.FillWithAir(w, box, 5, 7, 6, 6, 7, 8)
	j.FillWithAir(w, box, 5, 1, 2, 6, 2, 2)
	j.FillWithAir(w, box, 5, 2, 12, 6, 2, 12)
	j.FillWithAir(w, box, 5, 5, 1, 6, 5, 1)
	j.FillWithAir(w, box, 5, 5, 13, 6, 5, 13)

	j.SetBlockState(w, 0, 0, 1, 5, 5, box)
	j.SetBlockState(w, 0, 0, 10, 5, 5, box)
	j.SetBlockState(w, 0, 0, 1, 5, 9, box)
	j.SetBlockState(w, 0, 0, 10, 5, 9, box)

	for i := 0; i <= 14; i += 14 {
		j.FillWithRandomizedBlocks(w, box, 2, 4, i, 2, 5, i, false, rnd, cobblestoneSelector)
		j.FillWithRandomizedBlocks(w, box, 4, 4, i, 4, 5, i, false, rnd, cobblestoneSelector)
		j.FillWithRandomizedBlocks(w, box, 7, 4, i, 7, 5, i, false, rnd, cobblestoneSelector)
		j.FillWithRandomizedBlocks(w, box, 9, 4, i, 9, 5, i, false, rnd, cobblestoneSelector)
	}

	j.FillWithRandomizedBlocks(w, box, 5, 6, 0, 6, 6, 0, false, rnd, cobblestoneSelector)

	for l := 0; l <= 11; l += 11 {
		for k := 2; k <= 12; k += 2 {
			j.FillWithRandomizedBlocks(w, box, l, 4, k, l, 5, k, false, rnd, cobblestoneSelector)
		}
		j.FillWithRandomizedBlocks(w, box, l, 6, 5, l, 6, 5, false, rnd, cobblestoneSelector)
		j.FillWithRandomizedBlocks(w, box, l, 6, 9, l, 6, 9, false, rnd, cobblestoneSelector)
	}

	j.FillWithRandomizedBlocks(w, box, 2, 7, 2, 2, 9, 2, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 9, 7, 2, 9, 9, 2, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 2, 7, 12, 2, 9, 12, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 9, 7, 12, 9, 9, 12, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 4, 9, 4, 4, 9, 4, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 7, 9, 4, 7, 9, 4, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 4, 9, 10, 4, 9, 10, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 7, 9, 10, 7, 9, 10, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 5, 9, 7, 6, 9, 7, false, rnd, cobblestoneSelector)

	STATE_EAST := byte(0)
	STATE_WEST := byte(1)
	STATE_SOUTH := byte(2)
	STATE_NORTH := byte(3)

	j.SetBlockState(w, STONE_STAIRS, STATE_NORTH, 5, 9, 6, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_NORTH, 6, 9, 6, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_SOUTH, 5, 9, 8, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_SOUTH, 6, 9, 8, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_NORTH, 4, 0, 0, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_NORTH, 5, 0, 0, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_NORTH, 6, 0, 0, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_NORTH, 7, 0, 0, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_NORTH, 4, 1, 8, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_NORTH, 4, 2, 9, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_NORTH, 4, 3, 10, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_NORTH, 7, 1, 8, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_NORTH, 7, 2, 9, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_NORTH, 7, 3, 10, box)

	j.FillWithRandomizedBlocks(w, box, 4, 1, 9, 4, 1, 9, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 7, 1, 9, 7, 1, 9, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 4, 1, 10, 7, 2, 10, false, rnd, cobblestoneSelector)

	j.FillWithRandomizedBlocks(w, box, 5, 4, 5, 6, 4, 5, false, rnd, cobblestoneSelector)
	j.SetBlockState(w, STONE_STAIRS, STATE_EAST, 4, 4, 5, box)
	j.SetBlockState(w, STONE_STAIRS, STATE_WEST, 7, 4, 5, box)

	for k := 0; k < 4; k++ {
		j.SetBlockState(w, STONE_STAIRS, STATE_SOUTH, 5, 0-k, 6+k, box)
		j.SetBlockState(w, STONE_STAIRS, STATE_SOUTH, 6, 0-k, 6+k, box)
		j.FillWithAir(w, box, 5, 0-k, 7+k, 6, 0-k, 9+k)
	}

	j.FillWithAir(w, box, 1, -3, 12, 10, -1, 13)
	j.FillWithAir(w, box, 1, -3, 1, 3, -1, 13)
	j.FillWithAir(w, box, 1, -3, 1, 9, -1, 5)

	for i := 1; i <= 13; i += 2 {
		j.FillWithRandomizedBlocks(w, box, 1, -3, i, 1, -2, i, false, rnd, cobblestoneSelector)
	}
	for i := 2; i <= 12; i += 2 {
		j.FillWithRandomizedBlocks(w, box, 1, -1, i, 3, -1, i, false, rnd, cobblestoneSelector)
	}

	j.FillWithRandomizedBlocks(w, box, 2, -2, 1, 5, -2, 1, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 7, -2, 1, 9, -2, 1, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 6, -3, 1, 6, -3, 1, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 6, -1, 1, 6, -1, 1, false, rnd, cobblestoneSelector)

	j.SetBlockState(w, TRIPWIRE_HOOK, 3|4, 1, -3, 8, box)
	j.SetBlockState(w, TRIPWIRE_HOOK, 1|4, 4, -3, 8, box)
	j.SetBlockState(w, TRIPWIRE, 4, 2, -3, 8, box)
	j.SetBlockState(w, TRIPWIRE, 4, 3, -3, 8, box)

	j.SetBlockState(w, REDSTONE, 0, 5, -3, 7, box)
	j.SetBlockState(w, REDSTONE, 0, 5, -3, 6, box)
	j.SetBlockState(w, REDSTONE, 0, 5, -3, 5, box)
	j.SetBlockState(w, REDSTONE, 0, 5, -3, 4, box)
	j.SetBlockState(w, REDSTONE, 0, 5, -3, 3, box)
	j.SetBlockState(w, REDSTONE, 0, 5, -3, 2, box)
	j.SetBlockState(w, REDSTONE, 0, 5, -3, 1, box)
	j.SetBlockState(w, REDSTONE, 0, 4, -3, 1, box)
	j.SetBlockState(w, MOSSY, 0, 3, -3, 1, box)

	if !j.placedTrap1 {
		j.SetBlockState(w, DISPENSER, 2, 3, -2, 1, box)
		j.placedTrap1 = true
	}

	j.SetBlockState(w, VINE, 4, 3, -2, 2, box)

	j.SetBlockState(w, TRIPWIRE_HOOK, 2|4, 7, -3, 1, box)
	j.SetBlockState(w, TRIPWIRE_HOOK, 0|4, 7, -3, 5, box)
	j.SetBlockState(w, TRIPWIRE, 4, 7, -3, 2, box)
	j.SetBlockState(w, TRIPWIRE, 4, 7, -3, 3, box)
	j.SetBlockState(w, TRIPWIRE, 4, 7, -3, 4, box)

	j.SetBlockState(w, REDSTONE, 0, 8, -3, 6, box)
	j.SetBlockState(w, REDSTONE, 0, 9, -3, 6, box)
	j.SetBlockState(w, REDSTONE, 0, 9, -3, 5, box)
	j.SetBlockState(w, MOSSY, 0, 9, -3, 4, box)
	j.SetBlockState(w, REDSTONE, 0, 9, -2, 4, box)

	if !j.placedTrap2 {
		j.SetBlockState(w, DISPENSER, 4, 9, -2, 3, box)
		j.placedTrap2 = true
	}

	j.SetBlockState(w, VINE, 8, 8, -1, 3, box)
	j.SetBlockState(w, VINE, 8, 8, -2, 3, box)

	if !j.placedMainChest {
		j.SetBlockState(w, 54, 2, 8, -3, 3, box)
		j.placedMainChest = true
	}

	j.SetBlockState(w, MOSSY, 0, 9, -3, 2, box)
	j.SetBlockState(w, MOSSY, 0, 8, -3, 1, box)
	j.SetBlockState(w, MOSSY, 0, 4, -3, 5, box)
	j.SetBlockState(w, MOSSY, 0, 5, -2, 5, box)
	j.SetBlockState(w, MOSSY, 0, 5, -1, 5, box)
	j.SetBlockState(w, MOSSY, 0, 6, -3, 5, box)
	j.SetBlockState(w, MOSSY, 0, 7, -2, 5, box)
	j.SetBlockState(w, MOSSY, 0, 7, -1, 5, box)
	j.SetBlockState(w, MOSSY, 0, 8, -3, 5, box)
	j.FillWithRandomizedBlocks(w, box, 9, -1, 1, 9, -1, 5, false, rnd, cobblestoneSelector)

	j.FillWithAir(w, box, 8, -3, 8, 10, -1, 10)
	j.SetBlockState(w, STONEBRICK, CHISELED, 8, -2, 11, box)
	j.SetBlockState(w, STONEBRICK, CHISELED, 9, -2, 11, box)
	j.SetBlockState(w, STONEBRICK, CHISELED, 10, -2, 11, box)
	j.SetBlockState(w, LEVER, 12, 8, -2, 12, box)
	j.SetBlockState(w, LEVER, 4, 8, -2, 12, box)
	j.SetBlockState(w, LEVER, 4, 9, -2, 12, box)
	j.SetBlockState(w, LEVER, 4, 10, -2, 12, box)

	j.FillWithRandomizedBlocks(w, box, 8, -3, 8, 8, -3, 10, false, rnd, cobblestoneSelector)
	j.FillWithRandomizedBlocks(w, box, 10, -3, 8, 10, -3, 10, false, rnd, cobblestoneSelector)
	j.SetBlockState(w, MOSSY, 0, 10, -2, 9, box)
	j.SetBlockState(w, REDSTONE, 0, 8, -2, 9, box)
	j.SetBlockState(w, REDSTONE, 0, 8, -2, 10, box)
	j.SetBlockState(w, REDSTONE, 0, 10, -1, 9, box)
	j.SetBlockState(w, STICKY_PISTON, 1, 9, -2, 8, box)
	j.SetBlockState(w, STICKY_PISTON, 4, 10, -2, 8, box)
	j.SetBlockState(w, STICKY_PISTON, 4, 10, -1, 8, box)
	j.SetBlockState(w, REPEATER, 0, 10, -2, 10, box)

	if !j.placedHiddenChest {
		j.SetBlockState(w, 54, 4, 9, -3, 10, box)
		j.placedHiddenChest = true
	}

	return true
}

type SwampHut struct {
	*ScatteredFeaturePiece
	hasWitch bool
}

func NewSwampHut(rnd *rand.Random, x, z int) *SwampHut {
	facing := rnd.NextBoundedInt(4)
	return &SwampHut{
		ScatteredFeaturePiece: NewScatteredFeaturePiece(2, rnd, x, 64, z, 7, 7, 9, facing),
	}
}

func (s *SwampHut) BuildComponent(component StructureComponent, components *[]StructureComponent, rnd *rand.Random) {
}

func (s *SwampHut) AddComponentParts(w WorldAccess, rnd *rand.Random, box *BoundingBox) bool {
	if !s.OffsetToAverageGroundLevel(w, box, 0) {
		return false
	}

	SPRUCE_PLANKS := byte(5)
	SPRUCE_PLANKS_META := byte(1)
	SPRUCE_LOG := byte(17)
	SPRUCE_LOG_META := byte(1)

	s.FillWithBlocks(w, box, 1, 0, 2, 1, 3, 2, SPRUCE_LOG, SPRUCE_LOG_META, SPRUCE_LOG, SPRUCE_LOG_META, false)
	s.FillWithBlocks(w, box, 5, 0, 2, 5, 3, 2, SPRUCE_LOG, SPRUCE_LOG_META, SPRUCE_LOG, SPRUCE_LOG_META, false)
	s.FillWithBlocks(w, box, 1, 0, 7, 1, 3, 7, SPRUCE_LOG, SPRUCE_LOG_META, SPRUCE_LOG, SPRUCE_LOG_META, false)
	s.FillWithBlocks(w, box, 5, 0, 7, 5, 3, 7, SPRUCE_LOG, SPRUCE_LOG_META, SPRUCE_LOG, SPRUCE_LOG_META, false)

	for i := 2; i <= 7; i += 5 {
		for j := 1; j <= 5; j += 4 {
			s.ReplaceAirAndLiquidDownwards(w, SPRUCE_LOG, SPRUCE_LOG_META, j, -1, i, box)
		}
	}

	s.FillWithBlocks(w, box, 1, 1, 1, 5, 1, 7, SPRUCE_PLANKS, SPRUCE_PLANKS_META, SPRUCE_PLANKS, SPRUCE_PLANKS_META, false)
	s.FillWithBlocks(w, box, 2, 1, 0, 4, 1, 0, SPRUCE_PLANKS, SPRUCE_PLANKS_META, SPRUCE_PLANKS, SPRUCE_PLANKS_META, false)

	s.FillWithBlocks(w, box, 1, 2, 3, 1, 3, 6, SPRUCE_PLANKS, SPRUCE_PLANKS_META, SPRUCE_PLANKS, SPRUCE_PLANKS_META, false)
	s.FillWithBlocks(w, box, 5, 2, 3, 5, 3, 6, SPRUCE_PLANKS, SPRUCE_PLANKS_META, SPRUCE_PLANKS, SPRUCE_PLANKS_META, false)
	s.FillWithBlocks(w, box, 2, 2, 7, 4, 3, 7, SPRUCE_PLANKS, SPRUCE_PLANKS_META, SPRUCE_PLANKS, SPRUCE_PLANKS_META, false)
	s.FillWithBlocks(w, box, 1, 0, 2, 1, 3, 2, SPRUCE_LOG, SPRUCE_LOG_META, SPRUCE_LOG, SPRUCE_LOG_META, false)
	s.FillWithBlocks(w, box, 5, 0, 2, 5, 3, 2, SPRUCE_LOG, SPRUCE_LOG_META, SPRUCE_LOG, SPRUCE_LOG_META, false)
	s.FillWithBlocks(w, box, 1, 0, 7, 1, 3, 7, SPRUCE_LOG, SPRUCE_LOG_META, SPRUCE_LOG, SPRUCE_LOG_META, false)
	s.FillWithBlocks(w, box, 5, 0, 7, 5, 3, 7, SPRUCE_LOG, SPRUCE_LOG_META, SPRUCE_LOG, SPRUCE_LOG_META, false)

	SPRUCE_STAIRS := byte(134)
	S_NORTH := byte(3)
	S_EAST := byte(0)
	S_WEST := byte(1)
	S_SOUTH := byte(2)
	s.FillWithBlocks(w, box, 0, 4, 1, 6, 4, 1, SPRUCE_STAIRS, S_NORTH, SPRUCE_STAIRS, S_NORTH, false)
	s.FillWithBlocks(w, box, 0, 4, 2, 0, 4, 7, SPRUCE_STAIRS, S_EAST, SPRUCE_STAIRS, S_EAST, false)
	s.FillWithBlocks(w, box, 6, 4, 2, 6, 4, 7, SPRUCE_STAIRS, S_WEST, SPRUCE_STAIRS, S_WEST, false)
	s.FillWithBlocks(w, box, 0, 4, 8, 6, 4, 8, SPRUCE_STAIRS, S_SOUTH, SPRUCE_STAIRS, S_SOUTH, false)

	s.FillWithBlocks(w, box, 1, 4, 2, 5, 4, 7, SPRUCE_PLANKS, SPRUCE_PLANKS_META, SPRUCE_PLANKS, SPRUCE_PLANKS_META, false)

	OAK_FENCE := byte(85)
	FLOWER_POT := byte(140)
	CRAFTING_TABLE := byte(58)
	CAULDRON := byte(118)

	s.SetBlockState(w, OAK_FENCE, 0, 2, 3, 2, box)
	s.SetBlockState(w, OAK_FENCE, 0, 3, 3, 7, box)

	s.SetBlockState(w, FLOWER_POT, 0, 1, 3, 5, box)
	s.SetBlockState(w, CRAFTING_TABLE, 0, 3, 2, 6, box)
	s.SetBlockState(w, CAULDRON, 0, 4, 2, 6, box)

	if !s.hasWitch {

		s.hasWitch = true
	}

	return true
}
