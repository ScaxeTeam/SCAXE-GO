package objects

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	BlockGrass = 2

	BlockTallGrass    = 31
	BlockYellowFlower = 37
	BlockRedFlower    = 38
	BlockCactus       = 81
	BlockSugarCane    = 83
	BlockPumpkin      = 86
	BlockMelon        = 103
	BlockDeadBush     = 32
)

type TallGrass struct{}

func (g *TallGrass) Generate(chunk *world.Chunk, random *rand.Rand, x, y, z int) bool {
	for i := 0; i < 128; i++ {
		nx := x + random.Intn(8) - random.Intn(8)
		ny := y + random.Intn(4) - random.Intn(4)
		nz := z + random.Intn(8) - random.Intn(8)

		if nx < 0 || nx >= 16 || nz < 0 || nz >= 16 || ny < 1 || ny >= 127 {
			continue
		}

		blockBelow, _ := chunk.GetBlock(nx, ny-1, nz)
		blockAt, _ := chunk.GetBlock(nx, ny, nz)

		if blockBelow == BlockGrass && blockAt == 0 {
			meta := byte(1)
			if random.Intn(3) == 0 {
				meta = 2
			}
			chunk.SetBlock(nx, ny, nz, BlockTallGrass, meta)
		}
	}
	return true
}

type Flower struct {
	FlowerType int
}

func NewYellowFlower() *Flower {
	return &Flower{FlowerType: BlockYellowFlower}
}

func NewRedFlower() *Flower {
	return &Flower{FlowerType: BlockRedFlower}
}

func (f *Flower) Generate(chunk *world.Chunk, random *rand.Rand, x, y, z int) bool {
	for i := 0; i < 64; i++ {
		nx := x + random.Intn(8) - random.Intn(8)
		ny := y + random.Intn(4) - random.Intn(4)
		nz := z + random.Intn(8) - random.Intn(8)

		if nx < 0 || nx >= 16 || nz < 0 || nz >= 16 || ny < 1 || ny >= 127 {
			continue
		}

		blockBelow, _ := chunk.GetBlock(nx, ny-1, nz)
		blockAt, _ := chunk.GetBlock(nx, ny, nz)

		if blockBelow == BlockGrass && blockAt == 0 {
			chunk.SetBlock(nx, ny, nz, byte(f.FlowerType), 0)
		}
	}
	return true
}

type Cactus struct{}

func (c *Cactus) Generate(chunk *world.Chunk, random *rand.Rand, x, y, z int) bool {
	if x < 1 || x >= 15 || z < 1 || z >= 15 || y >= 125 {
		return false
	}

	blockBelow, _ := chunk.GetBlock(x, y-1, z)
	if blockBelow != BlockSand {
		return false
	}

	for _, dx := range []int{-1, 1} {
		adj, _ := chunk.GetBlock(x+dx, y, z)
		if adj != 0 {
			return false
		}
	}
	for _, dz := range []int{-1, 1} {
		adj, _ := chunk.GetBlock(x, y, z+dz)
		if adj != 0 {
			return false
		}
	}

	height := 1 + random.Intn(3)
	for dy := 0; dy < height && y+dy < 128; dy++ {
		chunk.SetBlock(x, y+dy, z, BlockCactus, 0)
	}

	return true
}

type SugarCane struct{}

func (s *SugarCane) Generate(chunk *world.Chunk, random *rand.Rand, x, y, z int) bool {
	if y >= 126 {
		return false
	}

	blockBelow, _ := chunk.GetBlock(x, y-1, z)
	if blockBelow != BlockSand && blockBelow != BlockDirt && blockBelow != BlockGrass {
		return false
	}

	hasWater := false
	for _, dx := range []int{-1, 0, 1} {
		for _, dz := range []int{-1, 0, 1} {
			if dx == 0 && dz == 0 {
				continue
			}
			bx := x + dx
			bz := z + dz
			if bx >= 0 && bx < 16 && bz >= 0 && bz < 16 {
				block, _ := chunk.GetBlock(bx, y-1, bz)
				if block == 9 || block == 8 {
					hasWater = true
					break
				}
			}
		}
	}

	if !hasWater {
		return false
	}

	blockAt, _ := chunk.GetBlock(x, y, z)
	if blockAt != 0 {
		return false
	}

	height := 1 + random.Intn(3)
	for dy := 0; dy < height && y+dy < 128; dy++ {
		chunk.SetBlock(x, y+dy, z, BlockSugarCane, 0)
	}

	return true
}
