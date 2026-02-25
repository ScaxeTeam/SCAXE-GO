package block
type ActivateResult struct {
	Handled bool
	NewMeta    uint8
	MetaChange bool
	SyncPositions [][3]int
	OpenInventory bool
	InventoryType int
	PlaySound   string
	SoundVolume float64
	SoundPitch  float64
}
const (
	InventoryTypeChest        = 0
	InventoryTypeCrafting     = 1
	InventoryTypeEnchant      = 3
	InventoryTypeFurnace      = 2
	InventoryTypeAnvil        = 5
	InventoryTypeBrewingStand = 6
)
func DoorOnActivate(meta uint8, x, y, z int) ActivateResult {
	isTop := DoorIsTopHalf(meta)

	if isTop {
		bottomY := y - 1
		return ActivateResult{
			Handled:       true,
			MetaChange:    false,
			SyncPositions: [][3]int{{x, bottomY, z}},
			PlaySound:     "random.door_open",
			SoundVolume:   1.0,
			SoundPitch:    1.0,
		}
	}
	newMeta := DoorToggleOpen(meta)
	return ActivateResult{
		Handled:       true,
		NewMeta:       newMeta,
		MetaChange:    true,
		SyncPositions: [][3]int{{x, y + 1, z}},
		PlaySound:     "random.door_open",
		SoundVolume:   1.0,
		SoundPitch:    1.0,
	}
}

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

func FenceGateOnActivate(meta uint8, playerDirection int) ActivateResult {
	isOpen := FenceGateIsOpen(meta)
	var newMeta uint8

	if isOpen {
		newMeta = meta &^ FenceGateMaskOpen
	} else {
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

func ChestOnActivate() ActivateResult {
	return ActivateResult{
		Handled:       true,
		OpenInventory: true,
		InventoryType: InventoryTypeChest,
	}
}

func FurnaceOnActivate() ActivateResult {
	return ActivateResult{
		Handled:       true,
		OpenInventory: true,
		InventoryType: InventoryTypeFurnace,
	}
}

func CraftingTableOnActivate() ActivateResult {
	return ActivateResult{
		Handled:       true,
		OpenInventory: true,
		InventoryType: InventoryTypeCrafting,
	}
}
