package inventory

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/item"
)
type InventoryHolder interface {
	GetInventory() Inventory
}
type PositionHolder interface {
	GetX() int
	GetY() int
	GetZ() int
}
type EntityHolder interface {
	GetEntityID() int64
}
type Viewer interface {
	GetWindowID(inv Inventory) byte
	SendDataPacket(pk interface{})
	IsSpawned() bool
	GetViewerID() string
}
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
	AddItem(items ...item.Item) []item.Item
	CanAddItem(item item.Item) bool
	GetViewers() map[string]Viewer
	Open(who Viewer) bool
	Close(who Viewer)
	OnOpen(who Viewer)
	OnClose(who Viewer)
	SendContents(targets ...Viewer)
	SendSlot(index int, targets ...Viewer)
	GetType() *InventoryType
	GetHolder() InventoryHolder
	OnSlotChange(index int, before item.Item, send bool)
}
type BaseInventory struct {
	slots    []item.Item
	name     string
	title    string
	size     int
	invType  *InventoryType
	holder   InventoryHolder
	viewers  map[string]Viewer
	maxStack int

	OnSlotChangeFunc func(slot int, item item.Item)
}
func NewBaseInventory(holder InventoryHolder, invType *InventoryType, overrideSize int, overrideTitle string) *BaseInventory {
	size := invType.GetDefaultSize()
	if overrideSize > 0 {
		size = overrideSize
	}
	title := invType.GetDefaultTitle()
	if overrideTitle != "" {
		title = overrideTitle
	}

	return &BaseInventory{
		slots:    make([]item.Item, size),
		name:     invType.GetDefaultTitle(),
		title:    title,
		size:     size,
		invType:  invType,
		holder:   holder,
		viewers:  make(map[string]Viewer),
		maxStack: 64,
	}
}
func NewSimpleInventory(name string, size int) *BaseInventory {
	return &BaseInventory{
		slots:    make([]item.Item, size),
		name:     name,
		title:    name,
		size:     size,
		viewers:  make(map[string]Viewer),
		maxStack: 64,
	}
}

func (inv *BaseInventory) GetSize() int                  { return inv.size }
func (inv *BaseInventory) GetName() string               { return inv.name }
func (inv *BaseInventory) GetTitle() string              { return inv.title }
func (inv *BaseInventory) GetMaxStackSize() int          { return inv.maxStack }
func (inv *BaseInventory) SetMaxStackSize(s int)         { inv.maxStack = s }
func (inv *BaseInventory) GetType() *InventoryType       { return inv.invType }
func (inv *BaseInventory) GetHolder() InventoryHolder    { return inv.holder }
func (inv *BaseInventory) GetViewers() map[string]Viewer { return inv.viewers }

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

	if it.ID == 0 || it.Count <= 0 {
		return inv.ClearSlot(slot, true)
	}

	old := inv.GetItem(slot)
	inv.slots[slot] = it
	inv.OnSlotChange(slot, old, true)
	return nil
}

func (inv *BaseInventory) ClearSlot(index int, send bool) error {
	if index < 0 || index >= inv.size {
		return fmt.Errorf("slot index out of bounds: %d", index)
	}
	old := inv.slots[index]
	inv.slots[index] = item.Item{}
	inv.OnSlotChange(index, old, send)
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

func (inv *BaseInventory) SetContents(items []item.Item, send bool) {
	for i := 0; i < inv.size; i++ {
		if i < len(items) && items[i].ID != 0 {
			inv.slots[i] = items[i]
		} else {
			inv.slots[i] = item.Item{}
		}
	}
	if send {
		viewers := inv.getViewerSlice()
		if len(viewers) > 0 {
			inv.SendContents(viewers...)
		}
	}
}

func (inv *BaseInventory) Clear(send bool) {
	for i := range inv.slots {
		if inv.slots[i].ID != 0 {
			inv.ClearSlot(i, send)
		}
	}
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

func (inv *BaseInventory) First(it item.Item) int {
	for i, slotItem := range inv.slots {
		if slotItem.ID == it.ID && slotItem.Meta == it.Meta {
			return i
		}
	}
	return -1
}

func (inv *BaseInventory) FirstEmpty() int {
	for i := 0; i < inv.size; i++ {
		if inv.slots[i].ID == 0 {
			return i
		}
	}
	return -1
}

func (inv *BaseInventory) CanAddItem(it item.Item) bool {
	if it.Count <= 0 {
		return true
	}

	count := it.Count
	for _, slot := range inv.slots {
		if slot.ID == 0 {
			count -= it.GetMaxStackSize()
		} else if slot.Equals(it, true, true) {
			maxSize := slot.GetMaxStackSize()
			if slot.Count < maxSize {
				count -= (maxSize - slot.Count)
			}
		}
		if count <= 0 {
			return true
		}
	}
	return false
}

func (inv *BaseInventory) AddItem(items ...item.Item) []item.Item {
	var leftovers []item.Item

	for _, it := range items {
		if it.Count <= 0 || it.ID == 0 {
			continue
		}
		for i := 0; i < inv.size; i++ {
			slotItem := inv.slots[i]
			if slotItem.Equals(it, true, true) && slotItem.Count < slotItem.GetMaxStackSize() {
				maxSize := slotItem.GetMaxStackSize()
				amount := maxSize - slotItem.Count
				if amount > it.Count {
					amount = it.Count
				}
				if amount > 0 {
					updated := inv.slots[i]
					updated.Count += amount
					inv.SetItem(i, updated)
					it.Count -= amount
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
					newItem := it
					newItem.Count = toAdd
					inv.SetItem(i, newItem)
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

func (inv *BaseInventory) RemoveItem(items ...item.Item) []item.Item {
	var leftovers []item.Item

	for _, it := range items {
		if it.Count <= 0 || it.ID == 0 {
			continue
		}

		toRemove := it.Count

		for i := 0; i < inv.size; i++ {
			slotItem := inv.slots[i]
			if slotItem.Equals(it, true, true) {
				if slotItem.Count > toRemove {
					updated := inv.slots[i]
					updated.Count -= toRemove
					inv.SetItem(i, updated)
					toRemove = 0
					break
				} else {
					toRemove -= slotItem.Count
					inv.ClearSlot(i, true)
				}
			}
		}

		if toRemove > 0 {
			leftover := it
			leftover.Count = toRemove
			leftovers = append(leftovers, leftover)
		}
	}

	return leftovers
}

func (inv *BaseInventory) Open(who Viewer) bool {
	inv.OnOpen(who)
	return true
}

func (inv *BaseInventory) Close(who Viewer) {
	inv.OnClose(who)
}

func (inv *BaseInventory) OnOpen(who Viewer) {
	inv.viewers[who.GetViewerID()] = who
}

func (inv *BaseInventory) OnClose(who Viewer) {
	delete(inv.viewers, who.GetViewerID())
}

func (inv *BaseInventory) OnSlotChange(index int, before item.Item, send bool) {
	if send {
		viewers := inv.getViewerSlice()
		if len(viewers) > 0 {
			inv.SendSlot(index, viewers...)
		}
	}
	if inv.OnSlotChangeFunc != nil {
		inv.OnSlotChangeFunc(index, inv.GetItem(index))
	}
}

func (inv *BaseInventory) SendContents(targets ...Viewer) {
}

func (inv *BaseInventory) SendSlot(index int, targets ...Viewer) {
}

func (inv *BaseInventory) getViewerSlice() []Viewer {
	viewers := make([]Viewer, 0, len(inv.viewers))
	for _, v := range inv.viewers {
		viewers = append(viewers, v)
	}
	return viewers
}
