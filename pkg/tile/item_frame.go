package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

// ItemFrame 物品展示框 TileEntity
// 对应 PHP class ItemFrame extends Spawnable
type ItemFrame struct {
	SpawnableBase
}

// NewItemFrame 创建 ItemFrame 实例
// 对应 PHP ItemFrame::__construct(FullChunk $chunk, CompoundTag $nbt)
func NewItemFrame(chunk *world.Chunk, nbtData *nbt.CompoundTag) *ItemFrame {
	f := &ItemFrame{}

	// 确保 NBT 中存在必要字段
	if nbtData.Get("Item") == nil {
		// 默认空物品：Air 的 NBT 表示
		itemTag := nbt.NewCompoundTag("Item")
		itemTag.Set(nbt.NewShortTag("id", 0))
		itemTag.Set(nbt.NewShortTag("Damage", 0))
		itemTag.Set(nbt.NewByteTag("Count", 0))
		nbtData.Set(itemTag)
	}

	if nbtData.Get("ItemRotation") == nil {
		nbtData.Set(nbt.NewByteTag("ItemRotation", 0))
	}

	if nbtData.Get("ItemDropChance") == nil {
		nbtData.Set(nbt.NewFloatTag("ItemDropChance", 1.0))
	}

	InitSpawnableBase(&f.SpawnableBase, TypeItemFrame, chunk, nbtData)
	return f
}

// GetName 返回名称
func (f *ItemFrame) GetName() string {
	return "Item Frame"
}

// GetItemRotation 获取物品旋转值
func (f *ItemFrame) GetItemRotation() int8 {
	return f.NBT.GetByte("ItemRotation")
}

// SetItemRotation 设置物品旋转值
func (f *ItemFrame) SetItemRotation(rotation int8) {
	f.NBT.Set(nbt.NewByteTag("ItemRotation", rotation))
}

// GetItemDropChance 获取物品掉落几率
func (f *ItemFrame) GetItemDropChance() float32 {
	return f.NBT.GetFloat("ItemDropChance")
}

// SetItemDropChance 设置物品掉落几率
func (f *ItemFrame) SetItemDropChance(chance float32) {
	f.NBT.Set(nbt.NewFloatTag("ItemDropChance", chance))
}

// GetSpawnCompound 返回发送给客户端的 NBT 数据
// 对应 PHP ItemFrame::getSpawnCompound()
func (f *ItemFrame) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeItemFrame))
	compound.Set(nbt.NewIntTag("x", f.X))
	compound.Set(nbt.NewIntTag("y", f.Y))
	compound.Set(nbt.NewIntTag("z", f.Z))
	compound.Set(nbt.NewByteTag("ItemRotation", f.GetItemRotation()))
	compound.Set(nbt.NewFloatTag("ItemDropChance", f.GetItemDropChance()))

	// 如果有物品，则包含 Item 标签
	if itemTag := f.NBT.Get("Item"); itemTag != nil {
		if ct, ok := itemTag.(*nbt.CompoundTag); ok {
			// 检查物品 id 是否为 0 (Air)
			if ct.GetShort("id") != 0 {
				clone := ct.Clone().(*nbt.CompoundTag)
				clone.SetName("Item")
				compound.Set(clone)
			}
		}
	}

	return compound
}

// UpdateCompoundTag 处理客户端发来的更新
func (f *ItemFrame) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return false
}

// SpawnTo 向指定玩家发送数据包
func (f *ItemFrame) SpawnTo(sender PacketSender) bool {
	return SpawnTo(f, sender)
}

// SpawnToAll 向区块内所有玩家广播
func (f *ItemFrame) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(f, broadcaster)
}

func init() {
	RegisterTile(TypeItemFrame, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewItemFrame(chunk, nbtData)
	})
}
