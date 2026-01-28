package trees

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	BlockAir     = 0
	BlockLog     = 17
	BlockLeaves  = 18
	BlockLog2    = 162
	BlockLeaves2 = 161
	BlockDirt    = 3
	BlockGrass   = 2
	BlockSapling = 6
)

const (
	WoodOak     = 0
	WoodSpruce  = 1
	WoodBirch   = 2
	WoodJungle  = 3
	WoodAcacia  = 0
	WoodDarkOak = 1
)

type Tree interface {
	Generate(chunk *world.Chunk, random *rand.Rand, x, y, z int) bool
}

type GenericTree struct {
	LogType    int
	LogMeta    byte
	LeavesType int
	LeavesMeta byte
	MinHeight  int
	MaxHeight  int
}

func NewOakTree() *GenericTree {
	return &GenericTree{
		LogType:    BlockLog,
		LogMeta:    WoodOak,
		LeavesType: BlockLeaves,
		LeavesMeta: WoodOak,
		MinHeight:  4,
		MaxHeight:  7,
	}
}

func NewBirchTree() *GenericTree {
	return &GenericTree{
		LogType:    BlockLog,
		LogMeta:    WoodBirch,
		LeavesType: BlockLeaves,
		LeavesMeta: WoodBirch,
		MinHeight:  5,
		MaxHeight:  8,
	}
}

func NewSpruceTree() *GenericTree {
	return &GenericTree{
		LogType:    BlockLog,
		LogMeta:    WoodSpruce,
		LeavesType: BlockLeaves,
		LeavesMeta: WoodSpruce,
		MinHeight:  6,
		MaxHeight:  10,
	}
}

func NewJungleTree() *GenericTree {
	return &GenericTree{
		LogType:    BlockLog,
		LogMeta:    WoodJungle,
		LeavesType: BlockLeaves,
		LeavesMeta: WoodJungle,
		MinHeight:  4,
		MaxHeight:  12,
	}
}

func NewAcaciaTree() *GenericTree {
	return &GenericTree{
		LogType:    BlockLog2,
		LogMeta:    WoodAcacia,
		LeavesType: BlockLeaves2,
		LeavesMeta: WoodAcacia,
		MinHeight:  5,
		MaxHeight:  8,
	}
}

func (t *GenericTree) Generate(chunk *world.Chunk, random *rand.Rand, x, y, z int) bool {
	height := t.MinHeight + random.Intn(t.MaxHeight-t.MinHeight+1)

	if y+height+1 > 127 {
		return false
	}

	groundID, _ := chunk.GetBlock(x, y-1, z)
	if groundID != BlockDirt && groundID != BlockGrass {
		return false
	}

	for dy := 0; dy < height; dy++ {
		blockID, _ := chunk.GetBlock(x, y+dy, z)
		if blockID != BlockAir && blockID != BlockLeaves && blockID != BlockLeaves2 {
			return false
		}
	}

	for dy := 0; dy < height; dy++ {
		chunk.SetBlock(x, y+dy, z, byte(t.LogType), t.LogMeta)
	}

	t.generateLeaves(chunk, x, y+height-3, z, height)

	return true
}

func (t *GenericTree) generateLeaves(chunk *world.Chunk, x, y, z, height int) {

	for dy := 0; dy < 4; dy++ {
		radius := 2
		if dy == 0 || dy == 3 {
			radius = 1
		}

		for dx := -radius; dx <= radius; dx++ {
			for dz := -radius; dz <= radius; dz++ {

				if (dx == -radius || dx == radius) && (dz == -radius || dz == radius) {
					continue
				}

				lx := x + dx
				ly := y + dy
				lz := z + dz

				if lx < 0 || lx >= 16 || lz < 0 || lz >= 16 || ly >= 128 {
					continue
				}

				blockID, _ := chunk.GetBlock(lx, ly, lz)
				if blockID == BlockAir {
					chunk.SetBlock(lx, ly, lz, byte(t.LeavesType), t.LeavesMeta)
				}
			}
		}
	}
}
