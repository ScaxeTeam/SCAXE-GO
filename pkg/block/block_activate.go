package block

// block_activate.go — 方块 OnActivate 交互实现
// 对应 PHP 各方块类中的 onActivate() 方法。
// 这些方法定义了右键方块时的行为逻辑。

// ActivateResult 交互结果，描述 OnActivate 应该产生的效果。
// 服务器层根据此结果执行实际的世界修改（setBlock/openInventory 等）。
type ActivateResult struct {
	Handled bool // 是否消耗了交互

	// ── 方块状态变更 ──
	// 非零时表示需要将方块的 meta 设置为此值
	// 例如门的开关、活板门的翻转
	NewMeta    uint8
	MetaChange bool // 显式标记 meta 是否有变化（因为 NewMeta=0 也可能是有效值）

	// ── 额外方块操作 ──
	// 需要更新的相邻方块坐标（如门需要同步上下半）
	// 格式: [x, y, z] 的列表
	SyncPositions [][3]int

	// ── 背包/UI 操作 ──
	// 非零时表示需要打开容器背包界面
	OpenInventory bool
	InventoryType int // 背包类型（对应 PHP InventoryType 常量）

	// ── 音效 ──
	PlaySound   string // 播放的音效ID（如 "random.door_open"）
	SoundVolume float64
	SoundPitch  float64
}

// 背包类型常量（对应 PHP WindowTypes）
const (
	InventoryTypeChest        = 0
	InventoryTypeCrafting     = 1
	InventoryTypeEnchant      = 3
	InventoryTypeFurnace      = 2
	InventoryTypeAnvil        = 5
	InventoryTypeBrewingStand = 6
)

// ---- Door OnActivate ----
// 对应 PHP Door::onActivate()
// 切换门的开/关状态，需要同步对半门

// DoorOnActivate 门的右键交互逻辑
// meta: 当前方块 meta
// otherHalfMeta: 对半门方块的 meta
// x, y, z: 当前方块坐标
func DoorOnActivate(meta uint8, x, y, z int) ActivateResult {
	isTop := DoorIsTopHalf(meta)

	if isTop {
		// 上半门：需要切换下半门的 meta
		bottomY := y - 1
		return ActivateResult{
			Handled:       true,
			MetaChange:    false,                     // 上半门的 meta 不变
			SyncPositions: [][3]int{{x, bottomY, z}}, // 告诉服务器同步下半的 meta
			PlaySound:     "random.door_open",
			SoundVolume:   1.0,
			SoundPitch:    1.0,
		}
	}

	// 下半门：切换自身的 meta
	newMeta := DoorToggleOpen(meta)
	return ActivateResult{
		Handled:       true,
		NewMeta:       newMeta,
		MetaChange:    true,
		SyncPositions: [][3]int{{x, y + 1, z}}, // 通知上半门
		PlaySound:     "random.door_open",
		SoundVolume:   1.0,
		SoundPitch:    1.0,
	}
}

// ---- Trapdoor OnActivate ----
// 对应 PHP Trapdoor::onActivate()

func TrapdoorOnActivate(meta uint8) ActivateResult {
	return ActivateResult{
		Handled:     true,
		NewMeta:     TrapdoorToggleOpen(meta),
		MetaChange:  true,
		PlaySound:   "random.door_open",
		SoundVolume: 1.0,
		SoundPitch:  1.0,
	}
}

// ---- FenceGate OnActivate ----
// 对应 PHP FenceGate::onActivate()
// 切换开/关态，同时如果打开，朝向会面向玩家

func FenceGateOnActivate(meta uint8, playerDirection int) ActivateResult {
	isOpen := FenceGateIsOpen(meta)
	var newMeta uint8

	if isOpen {
		// 关闭：保持朝向，去掉 open 位
		newMeta = meta &^ FenceGateMaskOpen
	} else {
		// 打开：方向面向玩家，加 open 位
		newMeta = (uint8(playerDirection) & FenceGateMaskDirection) | FenceGateMaskOpen
	}

	return ActivateResult{
		Handled:     true,
		NewMeta:     newMeta,
		MetaChange:  true,
		PlaySound:   "random.door_open",
		SoundVolume: 1.0,
		SoundPitch:  1.0,
	}
}

// ---- Chest OnActivate ----
// 对应 PHP Chest::onActivate()
// 打开背包界面（实际的背包创建在 server 层）

func ChestOnActivate() ActivateResult {
	return ActivateResult{
		Handled:       true,
		OpenInventory: true,
		InventoryType: InventoryTypeChest,
	}
}

// ---- Furnace OnActivate ----
// 对应 PHP Furnace::onActivate()

func FurnaceOnActivate() ActivateResult {
	return ActivateResult{
		Handled:       true,
		OpenInventory: true,
		InventoryType: InventoryTypeFurnace,
	}
}

// ---- CraftingTable OnActivate ----
// 对应 PHP Workbench::onActivate()
// 打开 3x3 合成界面

func CraftingTableOnActivate() ActivateResult {
	return ActivateResult{
		Handled:       true,
		OpenInventory: true,
		InventoryType: InventoryTypeCrafting,
	}
}
