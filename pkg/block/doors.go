package block

// doors.go — 所有门变种（6种木门 + 铁门）
// 对应 PHP: WoodDoor, SpruceDoor, BirchDoor, JungleDoor, AcaciaDoor, DarkOakDoor, IronDoor
//
// 木门: 嵌入 DoorBase, 斧子, 硬度3, 下半掉落对应门物品
// 铁门: 嵌入 DoorBase, 镐子, 硬度5, 下半+需镐才掉落, 不可右键开（仅红石）

// ---------- 木门通用 ----------

// WoodDoorBlock 木门方块（通用）
type WoodDoorBlock struct {
	DoorBase
	dropItemID int // 掉落的物品ID（门物品和门方块ID不同）
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

// GetDrops 只有下半掉落门物品
// 对应 PHP WoodDoor::getDrops(): 只有 (meta & 0x08) === 0 时掉落
// 注意：具体 meta 检查在 Level 层，这里返回标准掉落
func (b *WoodDoorBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: b.dropItemID, Meta: 0, Count: 1}}
}

// GetFuelTime 木门可以作燃料
func (b *WoodDoorBlock) GetFuelTime() int {
	return 200
}

// ---------- 铁门 ----------

// IronDoorBlock 铁门方块
type IronDoorBlock struct {
	DoorBase
}

// NewIronDoorBlock 创建铁门
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

// CanBeActivated 铁门不可右键开（仅红石）
// 对应 PHP IronDoor::onActivate() 中对玩家直接返回 true（不执行开关）
func (b *IronDoorBlock) CanBeActivated() bool {
	return false
}

// GetDrops 铁门需要镐才掉落
// 对应 PHP IronDoor::getDrops()
func (b *IronDoorBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: ItemIronDoor, Meta: 0, Count: 1}}
}

// ---------- 物品ID常量（门物品与门方块ID不同） ----------

const (
	ItemWoodenDoor  = 324
	ItemIronDoor    = 330 // 已在 block_ids 中定义为 IRON_DOOR
	ItemSpruceDoor  = 427
	ItemBirchDoor   = 428
	ItemJungleDoor  = 429
	ItemAcaciaDoor  = 430
	ItemDarkOakDoor = 431
)

// ---------- 注册 ----------

func init() {
	// 6 种木门
	Registry.Register(newWoodDoor(WOOD_DOOR_BLOCK, "Oak Door Block", ItemWoodenDoor))
	Registry.Register(newWoodDoor(SPRUCE_DOOR_BLOCK, "Spruce Door Block", ItemSpruceDoor))
	Registry.Register(newWoodDoor(BIRCH_DOOR_BLOCK, "Birch Door Block", ItemBirchDoor))
	Registry.Register(newWoodDoor(JUNGLE_DOOR_BLOCK, "Jungle Door Block", ItemJungleDoor))
	Registry.Register(newWoodDoor(ACACIA_DOOR_BLOCK, "Acacia Door Block", ItemAcaciaDoor))
	Registry.Register(newWoodDoor(DARK_OAK_DOOR_BLOCK, "Dark Oak Door Block", ItemDarkOakDoor))

	// 铁门
	Registry.Register(NewIronDoorBlock())
}
