package object

import "github.com/scaxe/scaxe-go/pkg/block"

const (
	FlowerPoppy       = 0
	FlowerBlueOrchid  = 1
	FlowerAllium      = 2
	FlowerAzureBluet  = 3
	FlowerRedTulip    = 4
	FlowerOrangeTulip = 5
	FlowerWhiteTulip  = 6
	FlowerPinkTulip   = 7
	FlowerOxeyeDaisy  = 8
)

const (
	DoublePlantSunflower = 0
	DoublePlantLilac     = 1
	DoublePlantGrass     = 2
	DoublePlantFern      = 3
	DoublePlantRoseBush  = 4
	DoublePlantPeony     = 5
)

type PlantType struct {
	BlockID byte
	Meta    byte
}

func GetFlowerType(meta byte) PlantType {
	return PlantType{BlockID: block.RED_FLOWER, Meta: meta}
}

func Dandelion() PlantType {
	return PlantType{BlockID: block.DANDELION, Meta: 0}
}
