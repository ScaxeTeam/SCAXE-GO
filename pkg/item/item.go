package item

import (
	"bytes"
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

type Item struct {
	ID      int
	Meta    int
	Count   int
	Name    string
	NBTData *nbt.CompoundTag
}

func NewItem(id, meta, count int) Item {

	if id > 0 && count <= 0 {
		count = 1
	}
	return Item{
		ID:    id & 0xFFFF,
		Meta:  meta & 0xFFFF,
		Count: count,
		Name:  GetItemName(id),
	}
}

func Air() Item {
	return NewItem(0, 0, 0)
}

func (i Item) IsAir() bool {
	return i.ID == 0 || i.Count <= 0
}

func (i Item) GetID() int {
	return i.ID
}

func (i Item) GetDamage() int {
	return i.Meta
}

func (i *Item) SetDamage(meta int) {
	i.Meta = meta & 0xFFFF
}

func (i Item) GetCount() int {
	return i.Count
}

func (i *Item) SetCount(count int) {
	i.Count = count
}

func (i Item) GetName() string {

	if i.HasCustomName() {
		return i.GetCustomName()
	}
	if i.Name != "" {
		return i.Name
	}
	return GetItemName(i.ID)
}

func (i *Item) SetName(name string) {
	i.Name = name
}

func (i Item) GetMaxStackSize() int {

	if IsTool(i.ID) {
		return 1
	}

	if IsArmor(i.ID) {
		return 1
	}

	switch i.ID {
	case EGG, SNOWBALL:
		return 16
	case BUCKET, BED, SIGN, CAKE:
		return 1
	case POTION, SPLASH_POTION:
		return 1
	}

	return 64
}

func (i Item) IsTool() bool {
	return IsTool(i.ID)
}

func (i Item) GetToolType() int {
	return GetToolType(i.ID)
}

func (i Item) GetToolTier() int {
	return GetToolTier(i.ID)
}

func (i Item) GetMaxDurability() int {
	return GetMaxDurability(i.ID)
}

func (i Item) GetMiningEfficiency() float64 {
	return GetMiningEfficiency(i.ID)
}

func (i Item) HasNBT() bool {
	return i.NBTData != nil && i.NBTData.Count() > 0
}

func (i Item) GetCompoundTagBytes() []byte {
	if i.NBTData == nil || i.NBTData.Count() == 0 {
		return nil
	}

	buf := new(bytes.Buffer)
	w := nbt.NewWriter(buf, nbt.LittleEndian)

	if err := i.NBTData.Write(w); err != nil {
		return nil
	}
	return buf.Bytes()
}

func (i *Item) GetNBT() *nbt.CompoundTag {
	if i.NBTData == nil {
		i.NBTData = nbt.NewCompoundTag("")
	}
	return i.NBTData
}

func (i *Item) SetNBT(tag *nbt.CompoundTag) {
	i.NBTData = tag
}

func (i *Item) ClearNBT() {
	i.NBTData = nil
}

func (i Item) HasCustomName() bool {
	if !i.HasNBT() {
		return false
	}
	display := i.NBTData.GetCompound("display")
	if display == nil {
		return false
	}
	return display.Has("Name")
}

func (i Item) GetCustomName() string {
	if !i.HasNBT() {
		return ""
	}
	display := i.NBTData.GetCompound("display")
	if display == nil {
		return ""
	}
	return display.GetString("Name")
}

func (i *Item) SetCustomName(name string) {
	if name == "" {
		i.ClearCustomName()
		return
	}
	nbtTag := i.GetNBT()
	display := nbtTag.GetCompound("display")
	if display == nil {
		display = nbt.NewCompoundTag("display")
		nbtTag.Set(display)
	}
	display.Set(nbt.NewStringTag("Name", name))
}

func (i *Item) ClearCustomName() {
	if !i.HasNBT() {
		return
	}
	display := i.NBTData.GetCompound("display")
	if display != nil {
		display.Remove("Name")
		if display.Count() == 0 {
			i.NBTData.Remove("display")
		}
	}
}

func (i Item) GetLore() []string {
	if !i.HasNBT() {
		return nil
	}
	display := i.NBTData.GetCompound("display")
	if display == nil {
		return nil
	}
	loreTag := display.GetList("Lore")
	if loreTag == nil {
		return nil
	}
	lore := make([]string, loreTag.Len())
	for j := 0; j < loreTag.Len(); j++ {
		if tag, ok := loreTag.Get(j).(*nbt.StringTag); ok {
			lore[j] = tag.Value().(string)
		}
	}
	return lore
}

func (i *Item) SetLore(lines []string) {
	nbtTag := i.GetNBT()
	display := nbtTag.GetCompound("display")
	if display == nil {
		display = nbt.NewCompoundTag("display")
		nbtTag.Set(display)
	}
	loreList := nbt.NewListTag("Lore", nbt.TagString)
	for _, line := range lines {
		loreList.Add(nbt.NewStringTag("", line))
	}
	display.Set(loreList)
}

func (i Item) Clone() Item {
	clone := Item{
		ID:    i.ID,
		Meta:  i.Meta,
		Count: i.Count,
		Name:  i.Name,
	}
	if i.NBTData != nil {
		clone.NBTData = i.NBTData.Clone().(*nbt.CompoundTag)
	}
	return clone
}

func (i Item) Equals(other Item, checkMeta, checkNBT bool) bool {
	if i.ID != other.ID {
		return false
	}
	if checkMeta && i.Meta != other.Meta {
		return false
	}
	if checkNBT {

		if i.HasNBT() != other.HasNBT() {
			return false
		}

		if i.HasNBT() {

			myBytes := i.GetCompoundTagBytes()
			otherBytes := other.GetCompoundTagBytes()
			return bytes.Equal(myBytes, otherBytes)
		}
	}
	return true
}

func (i Item) EqualsExact(other Item) bool {
	if !i.Equals(other, true, true) {
		return false
	}
	return i.Count == other.Count
}

func (i Item) String() string {
	return fmt.Sprintf("Item{ID: %d, Meta: %d, Count: %d, Name: %q}", i.ID, i.Meta, i.Count, i.GetName())
}

func (i Item) NBTSerialize(slot int) *nbt.CompoundTag {
	tag := nbt.NewCompoundTag("")
	tag.Set(nbt.NewShortTag("id", int16(i.ID)))
	tag.Set(nbt.NewByteTag("Count", int8(i.Count)))
	tag.Set(nbt.NewShortTag("Damage", int16(i.Meta)))

	if slot >= 0 {
		tag.Set(nbt.NewByteTag("Slot", int8(slot)))
	}

	if i.HasNBT() {
		itemTag := i.NBTData.Clone().(*nbt.CompoundTag)
		itemTag.SetName("tag")
		tag.Set(itemTag)
	}

	return tag
}

func NBTDeserialize(tag *nbt.CompoundTag) Item {
	id := int(tag.GetShort("id"))
	count := int(tag.GetByte("Count"))
	meta := int(tag.GetShort("Damage"))

	item := NewItem(id, meta, count)

	if itemTag := tag.GetCompound("tag"); itemTag != nil {
		item.NBTData = itemTag
	}

	return item
}
