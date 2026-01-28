package populators

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/object"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type TreePopulator struct {
	BaseAmount   int
	RandomAmount int
	Type         int
}

func NewTreePopulator(treeType int) *TreePopulator {
	return &TreePopulator{
		BaseAmount:   0,
		RandomAmount: 1,
		Type:         treeType,
	}
}

func (tp *TreePopulator) SetBaseAmount(amount int) {
	tp.BaseAmount = amount
}

func (tp *TreePopulator) SetRandomAmount(amount int) {
	tp.RandomAmount = amount
}

func (tp *TreePopulator) Populate(level populator.ChunkManager, chunk *world.Chunk, cx, cz int32, random *rand.Random) {

	amount := tp.BaseAmount
	if tp.RandomAmount > 0 {
		amount += random.NextRange(0, tp.RandomAmount+1)
	}

	chunkX := int32(chunk.X) * 16
	chunkZ := int32(chunk.Z) * 16

	for i := 0; i < amount; i++ {
		rx := random.NextRange(0, 16)
		rz := random.NextRange(0, 16)
		realX := chunkX + int32(rx)
		realZ := chunkZ + int32(rz)

		y := tp.getHighestWorkableBlock(level, realX, realZ)
		if y != -1 && y < 254 {
			var treeObj interface {
				Generate(populator.ChunkManager, *rand.Random, world.BlockPos) bool
			}

			switch tp.Type {
			case 1:
				treeObj = object.NewSpruceTree()
			case 2:
				treeObj = object.NewBirchTree(false)
			default:
				treeObj = object.NewOakTree()
			}

			treeObj.Generate(level, random, world.NewBlockPos(realX, int32(y), realZ))
		}
	}
}

func (p *TreePopulator) getHighestWorkableBlock(level populator.ChunkManager, x, z int32) int {

	for y := 127; y >= 0; y-- {
		id := level.GetBlockId(x, int32(y), z)

		if id == block.DIRT || id == block.GRASS || id == block.PODZOL {

			if y < 63 {

				aboveId := level.GetBlockId(x, int32(y+1), z)
				if aboveId == block.WATER || aboveId == block.STILL_WATER {
					return -1
				}
			}
			return y + 1
		}

		if id == block.AIR || id == block.SNOW_LAYER || id == block.LEAVES || id == block.LEAVES2 {
			continue
		}

		if id == block.WATER || id == block.STILL_WATER {
			return -1
		}

		return -1
	}
	return -1
}
