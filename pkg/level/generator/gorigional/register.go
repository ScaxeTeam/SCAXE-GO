package gorigional

import (
	"github.com/scaxe/scaxe-go/pkg/level/generator"
)

func init() {
	factory := func(settings map[string]interface{}) generator.Generator {
		return NewChunkGeneratorOverworld(0)
	}
	generator.RegisterGenerator("gorigional", factory)
	generator.RegisterGenerator("normal", factory)
}
