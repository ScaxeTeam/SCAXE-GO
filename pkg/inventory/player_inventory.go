package inventory

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

const (
	hotbarSize      = 9
	playerInvSize   = 36
	armorSlots      = 4
	playerTotalSize = playerInvSize + armorSlots
	SpecialArmor byte = 0x78
)
type PlayerInventory struct {
	*BaseInventory
	itemInHandIndex int
	hotbar          []int
}
func NewPlayerInventory() *PlayerInventory {
	inv := &PlayerInventory{
		BaseInventory:   NewSimpleInventory("Player Inventory", playerTotalSize),
		itemInHandIndex: 0,
		hotbar:          make([]int, hotbarSize),
	}
	for i := 0; i < hotbarSize; i++ {
		inv.hotbar[i] = i
	}
	return inv
}
func (inv *PlayerInventory) GetSize() int {
	return playerInvSize
}
func (inv *PlayerInventory) GetHotbarSize() int {
	return hotbarSize
}
func (inv *PlayerInventory) GetHotbarSlotIndex(index int) int {
	if index >= 0 && index < hotbarSize {
		return inv.hotbar[index]
	}
	return -1
}
func (inv *PlayerInventory) SetHotbarSlotIndex(index, slot int) {
	if index >= 0 && index < hotbarSize && slot >= -1 && slot < playerInvSize {
		inv.hotbar[index] = slot
	}
}
func (inv *PlayerInventory) GetHotbar() []int32 {
	hotbar := make([]int32, hotbarSize)
	for i := 0; i < hotbarSize; i++ {
		idx := inv.hotbar[i]
		if idx <= -1 {
			hotbar[i] = -1
		} else {
			hotbar[i] = int32(idx + hotbarSize)
		}
	}
	return hotbar
}
func (inv *PlayerInventory) GetHeldItemIndex() int {
	return inv.itemInHandIndex
}
func (inv *PlayerInventory) GetHeldItemSlot() int {
	return inv.GetHotbarSlotIndex(inv.itemInHandIndex)
}
func (inv *PlayerInventory) SetHeldItemIndex(hotbarSlotIndex int) error {
	if hotbarSlotIndex < 0 || hotbarSlotIndex >= hotbarSize {
		return fmt.Errorf("invalid hotbar index: %d", hotbarSlotIndex)
	}
	inv.itemInHandIndex = hotbarSlotIndex
	return nil
}
func (inv *PlayerInventory) SetHeldItemIndexWithMapping(hotbarSlotIndex, slotMapping int) {
	if hotbarSlotIndex < 0 || hotbarSlotIndex >= hotbarSize {
		return
	}
	inv.itemInHandIndex = hotbarSlotIndex
	slotMapping -= hotbarSize
	if slotMapping < 0 || slotMapping >= playerInvSize {
		slotMapping = -1
	}
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
func (inv *PlayerInventory) GetItemInHand() item.Item {
	slot := inv.GetHeldItemSlot()
	if slot < 0 {
		return item.Item{}
	}
	return inv.GetItem(slot)
}
func (inv *PlayerInventory) SetItemInHand(it item.Item) error {
	slot := inv.GetHeldItemSlot()
	if slot < 0 {
		return fmt.Errorf("no held item slot")
	}
	return inv.SetItem(slot, it)
}
func (inv *PlayerInventory) GetArmorItem(index int) item.Item {
	return inv.GetItem(playerInvSize + index)
}
func (inv *PlayerInventory) SetArmorItem(slot int, it item.Item) error {
	if slot < 0 || slot > 3 {
		return fmt.Errorf("invalid armor slot: %d", slot)
	}
	return inv.BaseInventory.SetItem(playerInvSize+slot, it)
}
func (inv *PlayerInventory) GetArmorContents() []item.Item {
	armor := make([]item.Item, armorSlots)
	for i := 0; i < armorSlots; i++ {
		armor[i] = inv.GetItem(playerInvSize + i)
	}
	return armor
}
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
func (inv *PlayerInventory) SendContents(targets ...Viewer) {
	totalSend := playerInvSize + hotbarSize
	items := make([]item.Item, totalSend)
	for i := 0; i < playerInvSize; i++ {
		items[i] = inv.GetItem(i)
	}
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
func (inv *PlayerInventory) SendSlot(index int, targets ...Viewer) {
	it := inv.GetItem(index)
	for _, viewer := range targets {
		windowID := viewer.GetWindowID(inv)
		pk := protocol.NewContainerSetSlotPacket(windowID, uint16(index), it)
		viewer.SendDataPacket(pk)
	}
}
func (inv *PlayerInventory) SendArmorContents(targets ...Viewer) {
	armor := inv.GetArmorContents()
	for _, viewer := range targets {
		pk := protocol.NewContainerSetContentPacket(SpecialArmor, armor)
		viewer.SendDataPacket(pk)
	}
}
func (inv *PlayerInventory) SendArmorSlot(index int, targets ...Viewer) {
	it := inv.GetItem(index)
	armorIdx := index - playerInvSize
	for _, viewer := range targets {
		pk := protocol.NewContainerSetSlotPacket(SpecialArmor, uint16(armorIdx), it)
		viewer.SendDataPacket(pk)
	}
}
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
func (inv *PlayerInventory) ClearAll() {
	for i := 0; i < playerTotalSize; i++ {
		inv.ClearSlot(i, false)
	}
	for i := 0; i < hotbarSize; i++ {
		inv.hotbar[i] = i
	}
}
