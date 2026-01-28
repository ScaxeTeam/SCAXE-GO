package populators

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/biome"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type GroundCoverPopulator struct{}

func NewGroundCoverPopulator() *GroundCoverPopulator {
	return &GroundCoverPopulator{}
}

func (gc *GroundCoverPopulator) Populate(level populator.ChunkManager, chunk *world.Chunk, cx, cz int32, random *rand.Random) {

	waterHeight := 62

	for colX := 0; colX < 16; colX++ {
		for colZ := 0; colZ < 16; colZ++ {
			biomeID := chunk.GetBiomeID(colX, colZ)
			b := biome.GetBiome(biomeID)
			cover := b.GetGroundCover()

			if len(cover) > 0 {
				diffY := 0

				if !isSolid(cover[0].ID) {
					diffY = 1
				}

				y := 127
				for ; y > 0; y-- {
					id := chunk.GetBlockId(colX, y, colZ)
					if id != 0 && !isTransparent(id) {
						break
					}
				}

				startY := y + diffY
				if startY > 127 {
					startY = 127
				}
				endY := startY - len(cover)

				for y := startY; y > endY && y >= 0; y-- {
					idx := startY - y
					if idx >= len(cover) {
						break
					}
					bState := cover[idx]

					currentID := chunk.GetBlockId(colX, y, colZ)
					if currentID == 0 && isSolid(bState.ID) {
						break
					}

					if y <= waterHeight && (bState.ID == block.GRASS || bState.ID == block.SNOW_LAYER) {

						aboveID := chunk.GetBlockId(colX, y+1, colZ)
						if aboveID == block.WATER {
							bState = block.NewBlockState(block.DIRT, 0)
						}
					}

					if y == waterHeight && bState.ID == block.SNOW_LAYER {
						bState = block.NewBlockState(block.ICE, 0)
					}

					chunk.SetBlock(colX, y, colZ, byte(bState.ID), byte(bState.Meta))
				}
			}
		}
	}
}

func isSolid(id uint8) bool {

	return id != block.AIR && id != block.SNOW_LAYER && id != block.TALL_GRASS && id != block.DANDELION && id != block.RED_FLOWER
}

func isTransparent(id uint8) bool {

	return id == block.AIR || id == block.SNOW_LAYER || id == block.TALL_GRASS
}
