package inventory

import (
	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

// ContainerInventory is a BaseInventory that sends container open/close/content
// packets to viewing players. Used by Chest, Furnace, Workbench, etc.
type ContainerInventory struct {
	*BaseInventory
}

// NewContainerInventory creates a new container inventory.
func NewContainerInventory(holder InventoryHolder, invType *InventoryType, overrideSize int, overrideTitle string) *ContainerInventory {
	return &ContainerInventory{
		BaseInventory: NewBaseInventory(holder, invType, overrideSize, overrideTitle),
	}
}

// OnOpen sends a ContainerOpenPacket and the full contents to the player.
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

// OnClose sends a ContainerClosePacket to the player.
func (c *ContainerInventory) OnClose(who Viewer) {
	pk := protocol.NewContainerClosePacket()
	pk.WindowID = who.GetWindowID(c)
	who.SendDataPacket(pk)

	c.BaseInventory.OnClose(who)
}

// SendContents sends the full inventory contents to the specified viewers.
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

// SendSlot sends a single slot update to the specified viewers.
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
