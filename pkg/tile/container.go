package tile

import (
	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/nbt"
)

// Container 接口表示可以存储物品的 Tile（箱子、熔炉、漏斗等）
// 对应 PHP interface Container
type Container interface {
	// GetItem 获取指定槽位的物品
	GetItem(index int) item.Item

	// SetItem 设置指定槽位的物品
	SetItem(index int, it item.Item)

	// GetSize 获取容器大小（槽位数）
	GetSize() int
}

// ---------- ContainerBase 提供公共的存储实现 ----------

// ContainerBase 为实现 Container 接口提供基于切片的存储
type ContainerBase struct {
	items []item.Item
}

// InitContainerBase 初始化容器，分配指定大小的物品槽
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

// GetContents 获取容器内所有物品的副本
func (c *ContainerBase) GetContents() []item.Item {
	result := make([]item.Item, len(c.items))
	copy(result, c.items)
	return result
}

// SetContents 批量设置容器内容
func (c *ContainerBase) SetContents(items []item.Item) {
	for i := 0; i < len(c.items); i++ {
		if i < len(items) {
			c.items[i] = items[i]
		} else {
			c.items[i] = item.Air()
		}
	}
}

// ClearAll 清空容器
func (c *ContainerBase) ClearAll() {
	for i := range c.items {
		c.items[i] = item.Air()
	}
}

// SaveItemsToNBT 将容器物品保存到 NBT CompoundTag（Items 列表）
func (c *ContainerBase) SaveItemsToNBT(nbtData *nbt.CompoundTag) {
	itemsList := nbt.NewListTag("Items", nbt.TagCompound)
	for i, it := range c.items {
		if !it.IsAir() {
			itemsList.Add(it.NBTSerialize(i))
		}
	}
	nbtData.Set(itemsList)
}

// LoadItemsFromNBT 从 NBT CompoundTag 的 Items 列表加载物品
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
