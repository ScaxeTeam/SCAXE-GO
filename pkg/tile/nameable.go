package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
)
type Nameable interface {
	GetCustomName() string
	SetCustomName(name string)
	HasCustomName() bool
}
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
func (n *NameableBase) LoadNameFromNBT(nbtData *nbt.CompoundTag) {
	if name := nbtData.GetString("CustomName"); name != "" {
		n.customName = name
	}
}
func (n *NameableBase) SaveNameToNBT(nbtData *nbt.CompoundTag) {
	if n.customName != "" {
		nbtData.Set(nbt.NewStringTag("CustomName", n.customName))
	}
}
