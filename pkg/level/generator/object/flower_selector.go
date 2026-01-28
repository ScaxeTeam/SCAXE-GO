package object

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/world"
)

func GetFlowerTypeForNoise(noiseVal float64) PlantType {

	n := math.Max(-1.0, math.Min(1.0, noiseVal))

	if n < -0.8 {
		return PlantType{BlockID: block.RED_FLOWER, Meta: FlowerWhiteTulip}
	} else if n < -0.6 {
		return PlantType{BlockID: block.RED_FLOWER, Meta: FlowerOrangeTulip}
	} else if n < -0.4 {
		return PlantType{BlockID: block.RED_FLOWER, Meta: FlowerRedTulip}
	} else if n < -0.2 {
		return PlantType{BlockID: block.RED_FLOWER, Meta: FlowerPinkTulip}
	} else if n < 0.0 {
		return PlantType{BlockID: block.RED_FLOWER, Meta: FlowerOxeyeDaisy}
	} else if n < 0.2 {
		return PlantType{BlockID: block.RED_FLOWER, Meta: FlowerPoppy}
	} else if n < 0.4 {
		return PlantType{BlockID: block.RED_FLOWER, Meta: FlowerAllium}
	} else if n < 0.6 {
		return PlantType{BlockID: block.RED_FLOWER, Meta: FlowerAzureBluet}
	} else if n < 0.8 {

		return PlantType{BlockID: block.RED_FLOWER, Meta: FlowerRedTulip}
	} else {

		return PlantType{BlockID: block.DANDELION, Meta: 0}
	}
}

type FlowerSelector interface {
	Select(pos world.BlockPos, r float64) PlantType
}
