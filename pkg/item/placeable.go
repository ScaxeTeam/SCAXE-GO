package item

const (
	blockWoodenDoor          uint8 = 64
	blockIronDoor            uint8 = 71
	blockSpruceDoor          uint8 = 193
	blockBirchDoor           uint8 = 194
	blockJungleDoor          uint8 = 195
	blockAcaciaDoor          uint8 = 196
	blockDarkOakDoor         uint8 = 197
	blockSignPost            uint8 = 63
	blockBedBlock            uint8 = 26
	blockCakeBlock           uint8 = 92
	blockRedstoneWire        uint8 = 55
	blockSugarcane           uint8 = 83
	blockFlowerPot           uint8 = 140
	blockCauldron            uint8 = 118
	blockBrewingStand        uint8 = 117
	blockHopper              uint8 = 154
	blockUnpoweredRepeater   uint8 = 93
	blockUnpoweredComparator uint8 = 149
)

var doorItemToBlock = map[int]uint8{
	WOODEN_DOOR:   blockWoodenDoor,
	IRON_DOOR:     blockIronDoor,
	SPRUCE_DOOR:   blockSpruceDoor,
	BIRCH_DOOR:    blockBirchDoor,
	JUNGLE_DOOR:   blockJungleDoor,
	ACACIA_DOOR:   blockAcaciaDoor,
	DARK_OAK_DOOR: blockDarkOakDoor,
}

func IsDoorItem(id int) bool {
	_, ok := doorItemToBlock[id]
	return ok
}

func GetDoorBlockID(id int) uint8 {
	if bid, ok := doorItemToBlock[id]; ok {
		return bid
	}
	return 0
}

const (
	BoatWoodOak     = 0
	BoatWoodSpruce  = 1
	BoatWoodBirch   = 2
	BoatWoodJungle  = 3
	BoatWoodAcacia  = 4
	BoatWoodDarkOak = 5
)

var boatItemToWood = map[int]int{
	BOAT:          BoatWoodOak,
	SPRUCE_BOAT:   BoatWoodSpruce,
	BIRCH_BOAT:    BoatWoodBirch,
	JUNGLE_BOAT:   BoatWoodJungle,
	ACACIA_BOAT:   BoatWoodAcacia,
	DARK_OAK_BOAT: BoatWoodDarkOak,
}

func IsBoatItem(id int) bool {
	_, ok := boatItemToWood[id]
	return ok
}

func GetBoatWoodType(id int) int {
	if w, ok := boatItemToWood[id]; ok {
		return w
	}
	return -1
}

var placeableItemToBlock = map[int]uint8{
	SIGN:          blockSignPost,
	BED:           blockBedBlock,
	CAKE:          blockCakeBlock,
	REDSTONE:      blockRedstoneWire,
	SUGARCANE:     blockSugarcane,
	FLOWER_POT:    blockFlowerPot,
	CAULDRON:      blockCauldron,
	BREWING_STAND: blockBrewingStand,
	HOPPER:        blockHopper,
	REPEATER:      blockUnpoweredRepeater,
	COMPARATOR:    blockUnpoweredComparator,
}

func IsPlaceableItem(id int) bool {
	if IsDoorItem(id) || IsBoatItem(id) {
		return true
	}
	_, ok := placeableItemToBlock[id]
	return ok
}

func GetPlaceableBlockID(id int) uint8 {
	if bid := GetDoorBlockID(id); bid != 0 {
		return bid
	}
	if bid, ok := placeableItemToBlock[id]; ok {
		return bid
	}
	return 0
}
