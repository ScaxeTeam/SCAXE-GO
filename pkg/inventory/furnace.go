package inventory

import (
	"github.com/scaxe/scaxe-go/pkg/item"
)

// Furnace slot indices
const (
	FurnaceSlotSmelting = 0 // Input item being smelted
	FurnaceSlotFuel     = 1 // Fuel item
	FurnaceSlotResult   = 2 // Output/result item
)

// FurnaceUpdateNotifier is implemented by the Furnace tile to receive
// notification when inventory slots change (triggering cook checks).
type FurnaceUpdateNotifier interface {
	SetNeedUpdate()
}

// FurnaceInventory is a 3-slot ContainerInventory for furnaces.
type FurnaceInventory struct {
	*ContainerInventory
}

// NewFurnaceInventory creates a furnace inventory (3 slots).
func NewFurnaceInventory(holder InventoryHolder) *FurnaceInventory {
	f := &FurnaceInventory{
		ContainerInventory: NewContainerInventory(
			holder,
			GetInventoryType(TypeFurnace),
			0, "",
		),
	}
	return f
}

// GetSmelting returns the item being smelted (slot 0).
func (f *FurnaceInventory) GetSmelting() item.Item {
	return f.GetItem(FurnaceSlotSmelting)
}

// SetSmelting sets the item being smelted (slot 0).
func (f *FurnaceInventory) SetSmelting(it item.Item) error {
	return f.SetItem(FurnaceSlotSmelting, it)
}

// GetFuel returns the fuel item (slot 1).
func (f *FurnaceInventory) GetFuel() item.Item {
	return f.GetItem(FurnaceSlotFuel)
}

// SetFuel sets the fuel item (slot 1).
func (f *FurnaceInventory) SetFuel(it item.Item) error {
	return f.SetItem(FurnaceSlotFuel, it)
}

// GetResult returns the smelting result (slot 2).
func (f *FurnaceInventory) GetResult() item.Item {
	return f.GetItem(FurnaceSlotResult)
}

// SetResult sets the smelting result (slot 2).
func (f *FurnaceInventory) SetResult(it item.Item) error {
	return f.SetItem(FurnaceSlotResult, it)
}

// OnSlotChange notifies the furnace tile that a slot changed,
// triggering a re-check of the smelting process.
func (f *FurnaceInventory) OnSlotChange(index int, before item.Item, send bool) {
	f.ContainerInventory.OnSlotChange(index, before, send)

	if notifier, ok := f.GetHolder().(FurnaceUpdateNotifier); ok {
		notifier.SetNeedUpdate()
	}
}
