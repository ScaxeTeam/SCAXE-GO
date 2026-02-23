package inventory

import (
	"github.com/scaxe/scaxe-go/pkg/item"
)

const (
	FurnaceSlotSmelting = 0
	FurnaceSlotFuel     = 1
	FurnaceSlotResult   = 2
)

type FurnaceUpdateNotifier interface {
	SetNeedUpdate()
}

type FurnaceInventory struct {
	*ContainerInventory
}

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

func (f *FurnaceInventory) GetSmelting() item.Item {
	return f.GetItem(FurnaceSlotSmelting)
}

func (f *FurnaceInventory) SetSmelting(it item.Item) error {
	return f.SetItem(FurnaceSlotSmelting, it)
}

func (f *FurnaceInventory) GetFuel() item.Item {
	return f.GetItem(FurnaceSlotFuel)
}

func (f *FurnaceInventory) SetFuel(it item.Item) error {
	return f.SetItem(FurnaceSlotFuel, it)
}

func (f *FurnaceInventory) GetResult() item.Item {
	return f.GetItem(FurnaceSlotResult)
}

func (f *FurnaceInventory) SetResult(it item.Item) error {
	return f.SetItem(FurnaceSlotResult, it)
}

func (f *FurnaceInventory) OnSlotChange(index int, before item.Item, send bool) {
	f.ContainerInventory.OnSlotChange(index, before, send)

	if notifier, ok := f.GetHolder().(FurnaceUpdateNotifier); ok {
		notifier.SetNeedUpdate()
	}
}
