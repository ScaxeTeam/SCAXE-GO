package inventory

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/item"
)

type PlayerInventory struct {
	*BaseInventory
	itemInHandIndex int
}

func NewPlayerInventory() *PlayerInventory {
	return &PlayerInventory{
		BaseInventory:   NewBaseInventory("Player Inventory", 40),
		itemInHandIndex: 0,
	}
}

func (inv *PlayerInventory) GetItemInHand() item.Item {
	return inv.GetItem(inv.itemInHandIndex)
}

func (inv *PlayerInventory) SetItemInHand(it item.Item) error {
	return inv.SetItem(inv.itemInHandIndex, it)
}

func (inv *PlayerInventory) GetSize() int {
	return 36
}

func (inv *PlayerInventory) GetHeldItemIndex() int {
	return inv.itemInHandIndex
}

func (inv *PlayerInventory) SetHeldItemIndex(index int) error {
	if index < 0 || index > 8 {
		return fmt.Errorf("invalid hotbar index: %d", index)
	}
	inv.itemInHandIndex = index
	return nil
}

func (inv *PlayerInventory) GetArmorContents() []item.Item {
	armor := make([]item.Item, 4)
	for i := 0; i < 4; i++ {
		armor[i] = inv.GetItem(36 + i)
	}
	return armor
}

func (inv *PlayerInventory) GetHotbarSlotIndex(index int) int {
	if index < 0 || index > 8 {
		return -1
	}

	return index
}

func (inv *PlayerInventory) GetHotbar() []int32 {
	hotbar := make([]int32, 9)
	for i := 0; i < 9; i++ {
		hotbar[i] = int32(i)
	}
	return hotbar
}
