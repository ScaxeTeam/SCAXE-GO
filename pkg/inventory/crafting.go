package inventory

import (
	"github.com/scaxe/scaxe-go/pkg/item"
)

type ItemDropper interface {
	DropItem(it item.Item)
}

type TemporaryInventory struct {
	*ContainerInventory
	resultSlotIndex int
}

func NewTemporaryInventory(holder InventoryHolder, invType *InventoryType, resultSlotIndex int) *TemporaryInventory {
	return &TemporaryInventory{
		ContainerInventory: NewContainerInventory(holder, invType, 0, ""),
		resultSlotIndex:    resultSlotIndex,
	}
}

func (t *TemporaryInventory) GetResultSlotIndex() int {
	return t.resultSlotIndex
}

func (t *TemporaryInventory) OnClose(who Viewer) {
	if dropper, ok := who.(ItemDropper); ok {
		for slot, it := range t.GetContents() {
			if slot == t.resultSlotIndex {

				continue
			}
			dropper.DropItem(it)
		}
	}
	t.Clear(false)
	t.ContainerInventory.OnClose(who)
}

type CraftingInventory struct {
	*TemporaryInventory
}

func NewCraftingInventory(holder InventoryHolder) *CraftingInventory {
	invType := GetInventoryType(TypeCrafting)
	return &CraftingInventory{
		TemporaryInventory: NewTemporaryInventory(holder, invType, 0),
	}
}

func NewWorkbenchInventory(holder InventoryHolder) *CraftingInventory {
	invType := GetInventoryType(TypeWorkbench)
	return &CraftingInventory{
		TemporaryInventory: NewTemporaryInventory(holder, invType, 0),
	}
}

type FakeBlockMenu struct {
	inv     Inventory
	x, y, z int
}

func NewFakeBlockMenu(inv Inventory, x, y, z int) *FakeBlockMenu {
	return &FakeBlockMenu{inv: inv, x: x, y: y, z: z}
}

func (f *FakeBlockMenu) GetInventory() Inventory { return f.inv }
func (f *FakeBlockMenu) GetX() int               { return f.x }
func (f *FakeBlockMenu) GetY() int               { return f.y }
func (f *FakeBlockMenu) GetZ() int               { return f.z }
