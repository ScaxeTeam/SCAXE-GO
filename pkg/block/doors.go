package block

type WoodDoorBlock struct {
	DoorBase
	dropItemID int
}

func newWoodDoor(blockID uint8, name string, dropItemID int) *WoodDoorBlock {
	return &WoodDoorBlock{
		DoorBase: DoorBase{
			TransparentBase: TransparentBase{
				BlockID:       blockID,
				BlockName:     name,
				BlockHardness: 3,
				BlockToolType: ToolTypeAxe,
				BlockCanPlace: true,
			},
		},
		dropItemID: dropItemID,
	}
}

func (b *WoodDoorBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: b.dropItemID, Meta: 0, Count: 1}}
}

func (b *WoodDoorBlock) GetFuelTime() int {
	return 200
}

type IronDoorBlock struct {
	DoorBase
}

func NewIronDoorBlock() *IronDoorBlock {
	return &IronDoorBlock{
		DoorBase: DoorBase{
			TransparentBase: TransparentBase{
				BlockID:       IRON_DOOR_BLOCK,
				BlockName:     "Iron Door Block",
				BlockHardness: 5,
				BlockToolType: ToolTypePickaxe,
				BlockCanPlace: true,
			},
		},
	}
}

func (b *IronDoorBlock) CanBeActivated() bool {
	return false
}

func (b *IronDoorBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: ItemIronDoor, Meta: 0, Count: 1}}
}

const (
	ItemWoodenDoor  = 324
	ItemIronDoor    = 330
	ItemSpruceDoor  = 427
	ItemBirchDoor   = 428
	ItemJungleDoor  = 429
	ItemAcaciaDoor  = 430
	ItemDarkOakDoor = 431
)

func init() {

	Registry.Register(newWoodDoor(WOOD_DOOR_BLOCK, "Oak Door Block", ItemWoodenDoor))
	Registry.Register(newWoodDoor(SPRUCE_DOOR_BLOCK, "Spruce Door Block", ItemSpruceDoor))
	Registry.Register(newWoodDoor(BIRCH_DOOR_BLOCK, "Birch Door Block", ItemBirchDoor))
	Registry.Register(newWoodDoor(JUNGLE_DOOR_BLOCK, "Jungle Door Block", ItemJungleDoor))
	Registry.Register(newWoodDoor(ACACIA_DOOR_BLOCK, "Acacia Door Block", ItemAcaciaDoor))
	Registry.Register(newWoodDoor(DARK_OAK_DOOR_BLOCK, "Dark Oak Door Block", ItemDarkOakDoor))

	Registry.Register(NewIronDoorBlock())
}
