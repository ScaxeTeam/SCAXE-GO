package item

// placeable.go — 可放置物品（门物品→方块 映射 + 船物品→实体 映射）
// 对应 PHP: item/WoodenDoor.php, item/IronDoor.php, item/SpruceDoor.php, ...
//           item/Boat.php
//
// 注意: 不能 import block 包（会造成循环依赖），所以直接使用方块ID常量。

// ============ 方块 ID 常量（避免循环依赖） ============

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

// ============ 门物品 → 方块 映射 ============

// doorItemToBlock 门物品ID → 对应方块ID
var doorItemToBlock = map[int]uint8{
	WOODEN_DOOR:   blockWoodenDoor,
	IRON_DOOR:     blockIronDoor,
	SPRUCE_DOOR:   blockSpruceDoor,
	BIRCH_DOOR:    blockBirchDoor,
	JUNGLE_DOOR:   blockJungleDoor,
	ACACIA_DOOR:   blockAcaciaDoor,
	DARK_OAK_DOOR: blockDarkOakDoor,
}

// IsDoorItem 判断物品是否为门物品
func IsDoorItem(id int) bool {
	_, ok := doorItemToBlock[id]
	return ok
}

// GetDoorBlockID 获取门物品对应的方块ID，非门返回 0
func GetDoorBlockID(id int) uint8 {
	if bid, ok := doorItemToBlock[id]; ok {
		return bid
	}
	return 0
}

// ============ 船物品 → 实体 映射 ============

// 船木材类型（对应 Boat 实体的 WoodID NBT 字段）
const (
	BoatWoodOak     = 0
	BoatWoodSpruce  = 1
	BoatWoodBirch   = 2
	BoatWoodJungle  = 3
	BoatWoodAcacia  = 4
	BoatWoodDarkOak = 5
)

// boatItemToWood 船物品ID → 木材类型
var boatItemToWood = map[int]int{
	BOAT:          BoatWoodOak,
	SPRUCE_BOAT:   BoatWoodSpruce,
	BIRCH_BOAT:    BoatWoodBirch,
	JUNGLE_BOAT:   BoatWoodJungle,
	ACACIA_BOAT:   BoatWoodAcacia,
	DARK_OAK_BOAT: BoatWoodDarkOak,
}

// IsBoatItem 判断物品是否为船物品
func IsBoatItem(id int) bool {
	_, ok := boatItemToWood[id]
	return ok
}

// GetBoatWoodType 获取船物品的木材类型, 非船返回 -1
func GetBoatWoodType(id int) int {
	if w, ok := boatItemToWood[id]; ok {
		return w
	}
	return -1
}

// ============ 其他可放置物品 ============

// placeableItemToBlock 其他可放置物品 → 对应方块ID
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

// IsPlaceableItem 判断物品是否为可放置物品（包括门、船、和其他可放置物品）
func IsPlaceableItem(id int) bool {
	if IsDoorItem(id) || IsBoatItem(id) {
		return true
	}
	_, ok := placeableItemToBlock[id]
	return ok
}

// GetPlaceableBlockID 获取物品对应的方块ID，门物品走 doorItemToBlock，其他走 placeableItemToBlock
// 非可放置物品返回 0
func GetPlaceableBlockID(id int) uint8 {
	if bid := GetDoorBlockID(id); bid != 0 {
		return bid
	}
	if bid, ok := placeableItemToBlock[id]; ok {
		return bid
	}
	return 0
}
