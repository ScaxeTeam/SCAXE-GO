package generator

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/biome"
)

func InitBiomes() {

	biome.InitBiomes()
}

func init() {
	InitBiomes()
}

func c(id uint8, meta uint8) block.BlockState {
	return block.NewBlockState(id, meta)
}
