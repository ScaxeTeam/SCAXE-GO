package gorigional

import (
	"github.com/scaxe/scaxe-go/pkg/level/generator"
)

func init() {
	generator.RegisterGenerator("gorigional", func(settings map[string]interface{}) generator.Generator {

		return NewChunkGeneratorOverworld(0)
	})
}
