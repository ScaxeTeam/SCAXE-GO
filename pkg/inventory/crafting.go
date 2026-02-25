package inventory

import (
	"github.com/scaxe/scaxe-go/pkg/item"
)

// ItemDropper is implemented by entities (players) that can drop items
// into the world when a temporary inventory is closed.
type ItemDropper interface {
	DropItem(it item.Item)
}

// TemporaryInventory is a ContainerInventory that drops its contents when
// closed, except for the result slot (virtual item). Used for crafting tables,
// enchanting tables, anvils, etc.
type TemporaryInventory struct {
	*ContainerInventory
	resultSlotIndex int
}

// NewTemporaryInventory creates a temporary inventory.
func NewTemporaryInventory(holder InventoryHolder, invType *InventoryType, resultSlotIndex int) *TemporaryInventory {
	return &TemporaryInventory{
		ContainerInventory: NewContainerInventory(holder, invType, 0, ""),
		resultSlotIndex:    resultSlotIndex,
	}
}

// GetResultSlotIndex returns the index of the virtual result slot.
func (t *TemporaryInventory) GetResultSlotIndex() int {
	return t.resultSlotIndex
}

// OnClose drops all items except the result slot, then clears the inventory.
func (t *TemporaryInventory) OnClose(who Viewer) {
	if dropper, ok := who.(ItemDropper); ok {
		for slot, it := range t.GetContents() {
			if slot == t.resultSlotIndex {
				// Result slot is virtual — do not drop
				continue
			}
			dropper.DropItem(it)
		}
	}
	t.Clear(false)
	t.ContainerInventory.OnClose(who)
}

// CraftingInventory manages crafting grid operations.
// The 2x2 player crafting grid (5 slots = 4 input + 1 result) or
// the 3x3 workbench grid (10 slots = 9 input + 1 result).
type CraftingInventory struct {
	*TemporaryInventory
}

// NewCraftingInventory creates a 2x2 player crafting inventory (4 input + 1 result = 5 slots).
func NewCraftingInventory(holder InventoryHolder) *CraftingInventory {
	invType := GetInventoryType(TypeCrafting)
	return &CraftingInventory{
		TemporaryInventory: NewTemporaryInventory(holder, invType, 0),
	}
}

// NewWorkbenchInventory creates a 3x3 workbench crafting inventory (9 input + 1 result = 10 slots).
func NewWorkbenchInventory(holder InventoryHolder) *CraftingInventory {
	invType := GetInventoryType(TypeWorkbench)
	return &CraftingInventory{
		TemporaryInventory: NewTemporaryInventory(holder, invType, 0),
	}
}

// FakeBlockMenu is a position-based InventoryHolder for virtual block menus
// like crafting tables and enchanting tables. It does not correspond to a real
// tile entity in the world.
type FakeBlockMenu struct {
	inv     Inventory
	x, y, z int
}

// NewFakeBlockMenu creates a fake block menu at the given position.
func NewFakeBlockMenu(inv Inventory, x, y, z int) *FakeBlockMenu {
	return &FakeBlockMenu{inv: inv, x: x, y: y, z: z}
}

func (f *FakeBlockMenu) GetInventory() Inventory { return f.inv }
func (f *FakeBlockMenu) GetX() int               { return f.x }
func (f *FakeBlockMenu) GetY() int               { return f.y }
func (f *FakeBlockMenu) GetZ() int               { return f.z }
