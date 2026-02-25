package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
)

// Nameable 接口表示可以拥有自定义名称的 Tile（箱子、熔炉等）
// 对应 PHP interface Nameable
type Nameable interface {
	// GetCustomName 获取自定义名称
	GetCustomName() string

	// SetCustomName 设置自定义名称
	SetCustomName(name string)

	// HasCustomName 是否拥有自定义名称
	HasCustomName() bool
}

// ---------- NameableBase 提供基于 NBT "CustomName" 的默认实现 ----------

// NameableBase 嵌入到需要 Nameable 功能的 Tile 中
type NameableBase struct {
	customName string
}

func (n *NameableBase) GetCustomName() string {
	return n.customName
}

func (n *NameableBase) SetCustomName(name string) {
	n.customName = name
}

func (n *NameableBase) HasCustomName() bool {
	return n.customName != ""
}

// LoadNameFromNBT 从 NBT 加载 CustomName
func (n *NameableBase) LoadNameFromNBT(nbtData *nbt.CompoundTag) {
	if name := nbtData.GetString("CustomName"); name != "" {
		n.customName = name
	}
}

// SaveNameToNBT 将 CustomName 保存到 NBT
func (n *NameableBase) SaveNameToNBT(nbtData *nbt.CompoundTag) {
	if n.customName != "" {
		nbtData.Set(nbt.NewStringTag("CustomName", n.customName))
	}
}
