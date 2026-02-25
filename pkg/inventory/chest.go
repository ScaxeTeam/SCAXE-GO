package inventory

import (
	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

// ChestInventory is a ContainerInventory for single chests.
// Sends BlockEventPacket to animate chest lid open/close.
type ChestInventory struct {
	*ContainerInventory
}

// NewChestInventory creates a single chest inventory (27 slots).
func NewChestInventory(holder InventoryHolder) *ChestInventory {
	return &ChestInventory{
		ContainerInventory: NewContainerInventory(
			holder,
			GetInventoryType(TypeChest),
			0, "",
		),
	}
}

// OnOpen sends chest open animation via BlockEventPacket, then delegates.
func (c *ChestInventory) OnOpen(who Viewer) {
	c.ContainerInventory.OnOpen(who)

	// Send chest-open animation when first viewer opens
	if len(c.GetViewers()) == 1 {
		c.sendBlockEvent(1, 2)
	}
}

// OnClose sends chest close animation via BlockEventPacket, then delegates.
func (c *ChestInventory) OnClose(who Viewer) {
	// Send chest-close animation when last viewer closes
	if len(c.GetViewers()) == 1 {
		c.sendBlockEvent(1, 0)
	}
	c.ContainerInventory.OnClose(who)
}

// sendBlockEvent broadcasts a BlockEventPacket at the holder's position.
func (c *ChestInventory) sendBlockEvent(case1, case2 int32) {
	holder := c.GetHolder()
	ph, ok := holder.(PositionHolder)
	if !ok {
		return
	}

	pk := protocol.NewBlockEventPacket()
	pk.X = int32(ph.GetX())
	pk.Y = int32(ph.GetY())
	pk.Z = int32(ph.GetZ())
	pk.Case1 = case1
	pk.Case2 = case2

	// Broadcast to all viewers
	for _, viewer := range c.GetViewers() {
		viewer.SendDataPacket(pk)
	}
}

// DoubleChestInventory combines two ChestInventory halves into 54 slots.
// Get/Set operations delegate to left (0-26) or right (27-53).
type DoubleChestInventory struct {
	*ContainerInventory
	left  *ChestInventory
	right *ChestInventory
}

// NewDoubleChestInventory creates a double chest inventory from two chest
// inventories. The left chest's holder is used as the position for the
// ContainerOpenPacket.
func NewDoubleChestInventory(leftHolder, rightHolder InventoryHolder, left, right *ChestInventory) *DoubleChestInventory {
	d := &DoubleChestInventory{
		ContainerInventory: NewContainerInventory(
			leftHolder,
			GetInventoryType(TypeDoubleChest),
			0, "",
		),
		left:  left,
		right: right,
	}
	return d
}

func (d *DoubleChestInventory) GetLeftSide() *ChestInventory  { return d.left }
func (d *DoubleChestInventory) GetRightSide() *ChestInventory { return d.right }

const singleChestSize = 27

// GetItem delegates to left or right half based on index.
func (d *DoubleChestInventory) GetItem(slot int) item.Item {
	if slot < singleChestSize {
		return d.left.GetItem(slot)
	}
	return d.right.GetItem(slot - singleChestSize)
}

// SetItem delegates to left or right half based on index.
func (d *DoubleChestInventory) SetItem(slot int, it item.Item) error {
	if slot < singleChestSize {
		return d.left.BaseInventory.SetItem(slot, it)
	}
	return d.right.BaseInventory.SetItem(slot-singleChestSize, it)
}

// ClearSlot delegates to left or right half.
func (d *DoubleChestInventory) ClearSlot(index int, send bool) error {
	if index < singleChestSize {
		return d.left.BaseInventory.ClearSlot(index, send)
	}
	return d.right.BaseInventory.ClearSlot(index-singleChestSize, send)
}

// GetContents returns all 54 slots merged from both halves.
func (d *DoubleChestInventory) GetContents() map[int]item.Item {
	contents := make(map[int]item.Item)
	for i := 0; i < d.GetSize(); i++ {
		it := d.GetItem(i)
		if it.ID != 0 {
			contents[i] = it
		}
	}
	return contents
}

// OnOpen sends chest open animation for both halves.
func (d *DoubleChestInventory) OnOpen(who Viewer) {
	d.ContainerInventory.OnOpen(who)

	// Also send BlockEvent for the right half
	if len(d.GetViewers()) == 1 {
		d.sendBlockEventForRight(1, 2)
	}
}

// OnClose sends chest close animation for both halves.
func (d *DoubleChestInventory) OnClose(who Viewer) {
	if len(d.GetViewers()) == 1 {
		d.sendBlockEventForRight(1, 0)
	}
	d.ContainerInventory.OnClose(who)
}

// sendBlockEventForRight broadcasts a BlockEventPacket at the right half's position.
func (d *DoubleChestInventory) sendBlockEventForRight(case1, case2 int32) {
	rightHolder := d.right.GetHolder()
	ph, ok := rightHolder.(PositionHolder)
	if !ok {
		return
	}

	pk := protocol.NewBlockEventPacket()
	pk.X = int32(ph.GetX())
	pk.Y = int32(ph.GetY())
	pk.Z = int32(ph.GetZ())
	pk.Case1 = case1
	pk.Case2 = case2

	for _, viewer := range d.GetViewers() {
		viewer.SendDataPacket(pk)
	}
}
