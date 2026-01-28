package inventory

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/item"
)

type Inventory interface {
	GetSize() int

	GetItem(slot int) item.Item

	SetItem(slot int, item item.Item) error

	GetContents() map[int]item.Item

	Clear(send bool)

	GetName() string

	Contains(item item.Item) bool

	First(item item.Item) int

	RemoveItem(items ...item.Item) []item.Item
}

type BaseInventory struct {
	slots        []item.Item
	name         string
	size         int
	OnSlotChange func(slot int, item item.Item)
}

func NewBaseInventory(name string, size int) *BaseInventory {
	return &BaseInventory{
		slots: make([]item.Item, size),
		name:  name,
		size:  size,
	}
}

func (inv *BaseInventory) GetSize() int {
	return inv.size
}

func (inv *BaseInventory) GetItem(slot int) item.Item {
	if slot < 0 || slot >= inv.size {
		return item.Item{}
	}
	return inv.slots[slot]
}

func (inv *BaseInventory) SetItem(slot int, it item.Item) error {
	if slot < 0 || slot >= inv.size {
		return fmt.Errorf("slot index out of bounds: %d", slot)
	}
	inv.slots[slot] = it
	if inv.OnSlotChange != nil {
		inv.OnSlotChange(slot, it)
	}
	return nil
}

func (inv *BaseInventory) GetContents() map[int]item.Item {
	contents := make(map[int]item.Item)
	for i, it := range inv.slots {
		if it.ID != 0 {
			contents[i] = it
		}
	}
	return contents
}

func (inv *BaseInventory) Clear(send bool) {
	for i := range inv.slots {
		if inv.slots[i].ID != 0 {
			inv.SetItem(i, item.Item{})
		}
	}
}

func (inv *BaseInventory) GetName() string {
	return inv.name
}

func (inv *BaseInventory) Contains(it item.Item) bool {
	count := it.Count
	if count <= 0 {
		count = 1
	}

	for _, i := range inv.slots {

		if i.Equals(it, true, true) {
			count -= i.Count
			if count <= 0 {
				return true
			}
		}
	}
	return false
}

func (inv *BaseInventory) CanAddItem(it item.Item) bool {
	if it.Count <= 0 {
		return true
	}

	clone := it
	for _, i := range inv.slots {
		if i.ID == 0 {
			return true
		}

		if i.Equals(clone, true, true) {
			maxSize := i.GetMaxStackSize()
			if i.Count < maxSize {
				room := maxSize - i.Count
				if room >= clone.Count {
					return true
				}
				clone.Count -= room
			}
		}
	}
	return false
}

func (inv *BaseInventory) AddItem(items ...item.Item) []item.Item {
	var leftovers []item.Item

	for _, it := range items {
		if it.Count <= 0 {
			continue
		}

		for i := 0; i < inv.size; i++ {
			slotItem := inv.slots[i]

			if slotItem.Equals(it, true, true) {
				maxSize := slotItem.GetMaxStackSize()
				if slotItem.Count < maxSize {
					room := maxSize - slotItem.Count
					toAdd := room
					if it.Count < room {
						toAdd = it.Count
					}

					updatedItem := inv.slots[i]
					updatedItem.Count += toAdd
					inv.SetItem(i, updatedItem)

					it.Count -= toAdd

					if it.Count <= 0 {
						break
					}
				}
			}
		}

		if it.Count > 0 {
			for i := 0; i < inv.size; i++ {
				if inv.slots[i].ID == 0 {
					maxSize := it.GetMaxStackSize()
					toAdd := maxSize
					if it.Count < maxSize {
						toAdd = it.Count
					}

					inv.SetItem(i, it)

					finalItem := inv.GetItem(i)
					finalItem.Count = toAdd
					inv.SetItem(i, finalItem)
					it.Count -= toAdd

					if it.Count <= 0 {
						break
					}
				}
			}
		}

		if it.Count > 0 {
			leftovers = append(leftovers, it)
		}
	}

	return leftovers
}

func (inv *BaseInventory) First(it item.Item) int {
	for i, slotItem := range inv.slots {

		if slotItem.ID == it.ID && slotItem.Meta == it.Meta {
			return i
		}
	}
	return -1
}

func (inv *BaseInventory) RemoveItem(items ...item.Item) []item.Item {
	var leftovers []item.Item

	for _, it := range items {
		if it.Count <= 0 {
			continue
		}

		toRemove := it.Count

		for i := 0; i < inv.size; i++ {
			slotItem := inv.slots[i]

			if slotItem.Equals(it, true, true) {
				if slotItem.Count > toRemove {

					updatedItem := inv.slots[i]
					updatedItem.Count -= toRemove
					inv.SetItem(i, updatedItem)
					toRemove = 0
					break
				} else {

					toRemove -= slotItem.Count

					inv.SetItem(i, item.NewItem(0, 0, 0))
				}
			}
		}

		if toRemove > 0 {
			it.Count = toRemove
			leftovers = append(leftovers, it)
		}
	}

	return leftovers
}
