package inventory

import (
	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

type ContainerInventory struct {
	*BaseInventory
}

func NewContainerInventory(holder InventoryHolder, invType *InventoryType, overrideSize int, overrideTitle string) *ContainerInventory {
	return &ContainerInventory{
		BaseInventory: NewBaseInventory(holder, invType, overrideSize, overrideTitle),
	}
}

func (c *ContainerInventory) OnOpen(who Viewer) {
	c.BaseInventory.OnOpen(who)

	pk := protocol.NewContainerOpenPacket()
	pk.WindowID = who.GetWindowID(c)
	pk.Type = c.GetType().GetNetworkType()
	pk.Slots = int16(c.GetSize())

	holder := c.GetHolder()
	if ph, ok := holder.(PositionHolder); ok {
		pk.X = int32(ph.GetX())
		pk.Y = int32(ph.GetY())
		pk.Z = int32(ph.GetZ())
	}
	if eh, ok := holder.(EntityHolder); ok {
		pk.EntityID = eh.GetEntityID()
	}

	who.SendDataPacket(pk)

	c.SendContents(who)
}

func (c *ContainerInventory) OnClose(who Viewer) {
	pk := protocol.NewContainerClosePacket()
	pk.WindowID = who.GetWindowID(c)
	who.SendDataPacket(pk)

	c.BaseInventory.OnClose(who)
}

func (c *ContainerInventory) SendContents(targets ...Viewer) {
	items := make([]item.Item, c.GetSize())
	for i := 0; i < c.GetSize(); i++ {
		items[i] = c.GetItem(i)
	}

	for _, viewer := range targets {
		windowID := viewer.GetWindowID(c)
		if windowID == 0xFF && !viewer.IsSpawned() {
			c.Close(viewer)
			continue
		}
		pk := protocol.NewContainerSetContentPacket(windowID, items)
		viewer.SendDataPacket(pk)
	}
}

func (c *ContainerInventory) SendSlot(index int, targets ...Viewer) {
	it := c.GetItem(index)

	for _, viewer := range targets {
		windowID := viewer.GetWindowID(c)
		if windowID == 0xFF {
			c.Close(viewer)
			continue
		}
		pk := protocol.NewContainerSetSlotPacket(windowID, uint16(index), it)
		viewer.SendDataPacket(pk)
	}
}
