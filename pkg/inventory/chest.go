package inventory

import (
	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)
type ChestInventory struct {
	*ContainerInventory
}
func NewChestInventory(holder InventoryHolder) *ChestInventory {
	return &ChestInventory{
		ContainerInventory: NewContainerInventory(
			holder,
			GetInventoryType(TypeChest),
			0, "",
		),
	}
}
func (c *ChestInventory) OnOpen(who Viewer) {
	c.ContainerInventory.OnOpen(who)
	if len(c.GetViewers()) == 1 {
		c.sendBlockEvent(1, 2)
	}
}
func (c *ChestInventory) OnClose(who Viewer) {
	if len(c.GetViewers()) == 1 {
		c.sendBlockEvent(1, 0)
	}
	c.ContainerInventory.OnClose(who)
}
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
	for _, viewer := range c.GetViewers() {
		viewer.SendDataPacket(pk)
	}
}
type DoubleChestInventory struct {
	*ContainerInventory
	left  *ChestInventory
	right *ChestInventory
}
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
func (d *DoubleChestInventory) GetItem(slot int) item.Item {
	if slot < singleChestSize {
		return d.left.GetItem(slot)
	}
	return d.right.GetItem(slot - singleChestSize)
}
func (d *DoubleChestInventory) SetItem(slot int, it item.Item) error {
	if slot < singleChestSize {
		return d.left.BaseInventory.SetItem(slot, it)
	}
	return d.right.BaseInventory.SetItem(slot-singleChestSize, it)
}
func (d *DoubleChestInventory) ClearSlot(index int, send bool) error {
	if index < singleChestSize {
		return d.left.BaseInventory.ClearSlot(index, send)
	}
	return d.right.BaseInventory.ClearSlot(index-singleChestSize, send)
}
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
func (d *DoubleChestInventory) OnOpen(who Viewer) {
	d.ContainerInventory.OnOpen(who)
	if len(d.GetViewers()) == 1 {
		d.sendBlockEventForRight(1, 2)
	}
}
func (d *DoubleChestInventory) OnClose(who Viewer) {
	if len(d.GetViewers()) == 1 {
		d.sendBlockEventForRight(1, 0)
	}
	d.ContainerInventory.OnClose(who)
}
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
