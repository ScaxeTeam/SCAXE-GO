package inventory

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

const (
	hotbarSize      = 9
	playerInvSize   = 36                         // container slots (0-35)
	armorSlots      = 4                          // armor slots (36-39)
	playerTotalSize = playerInvSize + armorSlots // 40

	// Special window IDs for armor
	SpecialArmor byte = 0x78
)

// PlayerInventory represents a player's personal inventory (36 slots + 4 armor).
// Hotbar slots [0-8] map to inventory slot indices via the hotbar array.
type PlayerInventory struct {
	*BaseInventory
	itemInHandIndex int
	hotbar          []int // hotbar[i] = inventory slot index linked to hotbar slot i
}

// NewPlayerInventory creates a new player inventory with default hotbar mapping.
func NewPlayerInventory() *PlayerInventory {
	inv := &PlayerInventory{
		BaseInventory:   NewSimpleInventory("Player Inventory", playerTotalSize),
		itemInHandIndex: 0,
		hotbar:          make([]int, hotbarSize),
	}
	// Default hotbar mapping: hotbar[i] = i
	for i := 0; i < hotbarSize; i++ {
		inv.hotbar[i] = i
	}
	return inv
}

// -- Size --

// GetSize returns the usable container size (36), excluding armor.
func (inv *PlayerInventory) GetSize() int {
	return playerInvSize
}

// GetHotbarSize returns the hotbar size (9).
func (inv *PlayerInventory) GetHotbarSize() int {
	return hotbarSize
}

// -- Hotbar mapping --

// GetHotbarSlotIndex returns the inventory slot index linked to a hotbar slot.
func (inv *PlayerInventory) GetHotbarSlotIndex(index int) int {
	if index >= 0 && index < hotbarSize {
		return inv.hotbar[index]
	}
	return -1
}

// SetHotbarSlotIndex changes the inventory slot linked to a hotbar slot.
func (inv *PlayerInventory) SetHotbarSlotIndex(index, slot int) {
	if index >= 0 && index < hotbarSize && slot >= -1 && slot < playerInvSize {
		inv.hotbar[index] = slot
	}
}

// GetHotbar returns the full hotbar mapping array.
func (inv *PlayerInventory) GetHotbar() []int32 {
	hotbar := make([]int32, hotbarSize)
	for i := 0; i < hotbarSize; i++ {
		idx := inv.hotbar[i]
		if idx <= -1 {
			hotbar[i] = -1
		} else {
			hotbar[i] = int32(idx + hotbarSize) // offset for PE protocol
		}
	}
	return hotbar
}

// -- Held item --

// GetHeldItemIndex returns which hotbar slot is currently selected.
func (inv *PlayerInventory) GetHeldItemIndex() int {
	return inv.itemInHandIndex
}

// GetHeldItemSlot returns the actual inventory slot of the held item.
func (inv *PlayerInventory) GetHeldItemSlot() int {
	return inv.GetHotbarSlotIndex(inv.itemInHandIndex)
}

// SetHeldItemIndex sets the active hotbar slot with optional slot remapping.
// slotMapping is the raw inventory slot sent by the client (9-44 range),
// or -1 to skip remapping.
func (inv *PlayerInventory) SetHeldItemIndex(hotbarSlotIndex int) error {
	if hotbarSlotIndex < 0 || hotbarSlotIndex >= hotbarSize {
		return fmt.Errorf("invalid hotbar index: %d", hotbarSlotIndex)
	}
	inv.itemInHandIndex = hotbarSlotIndex
	return nil
}

// SetHeldItemIndexWithMapping sets the held item with PE slot remapping.
func (inv *PlayerInventory) SetHeldItemIndexWithMapping(hotbarSlotIndex, slotMapping int) {
	if hotbarSlotIndex < 0 || hotbarSlotIndex >= hotbarSize {
		return
	}
	inv.itemInHandIndex = hotbarSlotIndex

	// Convert raw PE slot to inventory index
	slotMapping -= hotbarSize
	if slotMapping < 0 || slotMapping >= playerInvSize {
		slotMapping = -1
	}

	// Swap hotbar links if the target slot is already linked elsewhere
	if slotMapping != -1 {
		for i, linked := range inv.hotbar {
			if linked == slotMapping {
				inv.hotbar[i] = inv.hotbar[inv.itemInHandIndex]
				break
			}
		}
	}
	inv.hotbar[inv.itemInHandIndex] = slotMapping
}

// GetItemInHand returns the item currently being held.
func (inv *PlayerInventory) GetItemInHand() item.Item {
	slot := inv.GetHeldItemSlot()
	if slot < 0 {
		return item.Item{}
	}
	return inv.GetItem(slot)
}

// SetItemInHand sets the item in the held slot.
func (inv *PlayerInventory) SetItemInHand(it item.Item) error {
	slot := inv.GetHeldItemSlot()
	if slot < 0 {
		return fmt.Errorf("no held item slot")
	}
	return inv.SetItem(slot, it)
}

// -- Armor --

// GetArmorItem returns the armor item in the given armor slot (0-3).
func (inv *PlayerInventory) GetArmorItem(index int) item.Item {
	return inv.GetItem(playerInvSize + index)
}

// SetArmorItem sets the armor item in the given armor slot (0-3).
func (inv *PlayerInventory) SetArmorItem(slot int, it item.Item) error {
	if slot < 0 || slot > 3 {
		return fmt.Errorf("invalid armor slot: %d", slot)
	}
	return inv.BaseInventory.SetItem(playerInvSize+slot, it)
}

// GetArmorContents returns all 4 armor pieces [helmet, chestplate, leggings, boots].
func (inv *PlayerInventory) GetArmorContents() []item.Item {
	armor := make([]item.Item, armorSlots)
	for i := 0; i < armorSlots; i++ {
		armor[i] = inv.GetItem(playerInvSize + i)
	}
	return armor
}

// SetArmorContents sets all 4 armor pieces.
func (inv *PlayerInventory) SetArmorContents(items []item.Item) {
	for i := 0; i < armorSlots; i++ {
		if i < len(items) && items[i].ID != 0 {
			inv.SetArmorItem(i, items[i])
		} else {
			inv.SetArmorItem(i, item.Item{})
		}
	}
}

func (inv *PlayerInventory) GetHelmet() item.Item     { return inv.GetArmorItem(0) }
func (inv *PlayerInventory) GetChestplate() item.Item { return inv.GetArmorItem(1) }
func (inv *PlayerInventory) GetLeggings() item.Item   { return inv.GetArmorItem(2) }
func (inv *PlayerInventory) GetBoots() item.Item      { return inv.GetArmorItem(3) }

func (inv *PlayerInventory) SetHelmet(it item.Item) error     { return inv.SetArmorItem(0, it) }
func (inv *PlayerInventory) SetChestplate(it item.Item) error { return inv.SetArmorItem(1, it) }
func (inv *PlayerInventory) SetLeggings(it item.Item) error   { return inv.SetArmorItem(2, it) }
func (inv *PlayerInventory) SetBoots(it item.Item) error      { return inv.SetArmorItem(3, it) }

// -- Network sync --

// SendContents sends the full inventory contents (36 container + 9 dummy padding)
// to the specified viewers. PE requires 9 extra empty slots be sent.
func (inv *PlayerInventory) SendContents(targets ...Viewer) {
	// 36 real slots + 9 dummy = 45 sent to PE
	totalSend := playerInvSize + hotbarSize
	items := make([]item.Item, totalSend)
	for i := 0; i < playerInvSize; i++ {
		items[i] = inv.GetItem(i)
	}
	// Slots 36-44 are dummy empty slots for PE display
	for i := playerInvSize; i < totalSend; i++ {
		items[i] = item.Item{}
	}

	for _, viewer := range targets {
		windowID := viewer.GetWindowID(inv)
		pk := protocol.NewContainerSetContentPacket(windowID, items)
		pk.HotbarTypes = inv.GetHotbar()
		viewer.SendDataPacket(pk)
	}
}

// SendSlot sends a single slot update.
func (inv *PlayerInventory) SendSlot(index int, targets ...Viewer) {
	it := inv.GetItem(index)
	for _, viewer := range targets {
		windowID := viewer.GetWindowID(inv)
		pk := protocol.NewContainerSetSlotPacket(windowID, uint16(index), it)
		viewer.SendDataPacket(pk)
	}
}

// SendArmorContents sends all armor pieces to the specified viewers.
// Self receives ContainerSetContentPacket with SPECIAL_ARMOR window ID.
// Others receive MobArmorEquipmentPacket.
func (inv *PlayerInventory) SendArmorContents(targets ...Viewer) {
	armor := inv.GetArmorContents()
	for _, viewer := range targets {
		pk := protocol.NewContainerSetContentPacket(SpecialArmor, armor)
		viewer.SendDataPacket(pk)
	}
}

// SendArmorSlot sends a single armor slot update.
func (inv *PlayerInventory) SendArmorSlot(index int, targets ...Viewer) {
	it := inv.GetItem(index)
	armorIdx := index - playerInvSize
	for _, viewer := range targets {
		pk := protocol.NewContainerSetSlotPacket(SpecialArmor, uint16(armorIdx), it)
		viewer.SendDataPacket(pk)
	}
}

// SendHeldItem broadcasts the held item to the specified viewer via MobEquipmentPacket.
func (inv *PlayerInventory) SendHeldItem(viewer Viewer, entityID int64) {
	it := inv.GetItemInHand()
	pk := protocol.NewMobEquipmentPacket()
	pk.EntityID = entityID
	pk.ItemID = int16(it.ID)
	pk.ItemCount = int8(it.Count)
	pk.ItemMeta = uint16(it.Meta)
	slot := inv.GetHeldItemSlot()
	if slot < 0 {
		slot = 0
	}
	pk.Slot = uint8(slot)
	pk.SelectedSlot = uint8(inv.GetHeldItemIndex())
	viewer.SendDataPacket(pk)
}

// ClearAll clears all 40 slots and resets hotbar to default.
func (inv *PlayerInventory) ClearAll() {
	for i := 0; i < playerTotalSize; i++ {
		inv.ClearSlot(i, false)
	}
	for i := 0; i < hotbarSize; i++ {
		inv.hotbar[i] = i
	}
}
