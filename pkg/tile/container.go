package tile

import (
	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/nbt"
)
type Container interface {
	GetItem(index int) item.Item
	SetItem(index int, it item.Item)
	GetSize() int
}
type ContainerBase struct {
	items []item.Item
}
func InitContainerBase(c *ContainerBase, size int) {
	c.items = make([]item.Item, size)
}

func (c *ContainerBase) GetItem(index int) item.Item {
	if index < 0 || index >= len(c.items) {
		return item.Air()
	}
	return c.items[index]
}

func (c *ContainerBase) SetItem(index int, it item.Item) {
	if index < 0 || index >= len(c.items) {
		return
	}
	c.items[index] = it
}

func (c *ContainerBase) GetSize() int {
	return len(c.items)
}
func (c *ContainerBase) GetContents() []item.Item {
	result := make([]item.Item, len(c.items))
	copy(result, c.items)
	return result
}
func (c *ContainerBase) SetContents(items []item.Item) {
	for i := 0; i < len(c.items); i++ {
		if i < len(items) {
			c.items[i] = items[i]
		} else {
			c.items[i] = item.Air()
		}
	}
}
func (c *ContainerBase) ClearAll() {
	for i := range c.items {
		c.items[i] = item.Air()
	}
}
func (c *ContainerBase) SaveItemsToNBT(nbtData *nbt.CompoundTag) {
	itemsList := nbt.NewListTag("Items", nbt.TagCompound)
	for i, it := range c.items {
		if !it.IsAir() {
			itemsList.Add(it.NBTSerialize(i))
		}
	}
	nbtData.Set(itemsList)
}
func (c *ContainerBase) LoadItemsFromNBT(nbtData *nbt.CompoundTag) {
	itemsList := nbtData.GetList("Items")
	if itemsList == nil {
		return
	}
	for i := 0; i < itemsList.Len(); i++ {
		tag, ok := itemsList.Get(i).(*nbt.CompoundTag)
		if !ok {
			continue
		}
		slot := int(tag.GetByte("Slot"))
		if slot >= 0 && slot < len(c.items) {
			c.items[slot] = item.NBTDeserialize(tag)
		}
	}
}
