package tile

import (
	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

// ChestSize 单箱子槽位数
const ChestSize = 27

// Chest 箱子 TileEntity
// 对应 PHP class Chest extends Spawnable implements InventoryHolder, Container, Nameable
// 实现 Spawnable + Container + Nameable，支持大箱子配对
type Chest struct {
	SpawnableBase
	ContainerBase
	NameableBase

	// 大箱子配对坐标（-1 表示未配对）
	pairX int32
	pairZ int32
}

// NewChest 创建 Chest 实例
// 对应 PHP Chest::__construct(FullChunk $chunk, CompoundTag $nbt)
func NewChest(chunk *world.Chunk, nbtData *nbt.CompoundTag) *Chest {
	c := &Chest{
		pairX: -1,
		pairZ: -1,
	}

	InitSpawnableBase(&c.SpawnableBase, TypeChest, chunk, nbtData)
	InitContainerBase(&c.ContainerBase, ChestSize)

	// 从 NBT 加载物品
	c.ContainerBase.LoadItemsFromNBT(nbtData)

	// 从 NBT 加载自定义名称
	c.NameableBase.LoadNameFromNBT(nbtData)

	// 从 NBT 加载配对信息
	// 对应 PHP: isset($this->namedtag->pairx) && isset($this->namedtag->pairz)
	if nbtData.Has("pairx") && nbtData.Has("pairz") {
		c.pairX = nbtData.GetInt("pairx")
		c.pairZ = nbtData.GetInt("pairz")
	}

	return c
}

// ---------- Container 接口（覆盖 ContainerBase 的 GetSize） ----------

// GetSize 返回箱子容量 (27)
// 对应 PHP Chest::getSize() { return 27; }
func (c *Chest) GetSize() int {
	return ChestSize
}

// ---------- Nameable 接口实现 ----------

// GetName 返回箱子名称
// 对应 PHP Chest::getName() 优先返回 CustomName，否则返回 "Chest"
func (c *Chest) GetName() string {
	if c.HasCustomName() {
		return c.GetCustomName()
	}
	return "Chest"
}

// ---------- 大箱子配对 ----------

// IsPaired 是否已与另一个箱子配对
// 对应 PHP Chest::isPaired()
func (c *Chest) IsPaired() bool {
	return c.pairX != -1 && c.pairZ != -1
}

// GetPairPosition 获取配对箱子的坐标
func (c *Chest) GetPairPosition() (x, z int32) {
	return c.pairX, c.pairZ
}

// PairWith 与另一个箱子配对，形成大箱子
// 对应 PHP Chest::pairWith(Chest $tile)
func (c *Chest) PairWith(other *Chest) bool {
	if c.IsPaired() || other.IsPaired() {
		return false
	}

	c.createPair(other)
	return true
}

// Unpair 解除配对
// 对应 PHP Chest::unpair()
func (c *Chest) Unpair(getPairFunc func(x, z int32) *Chest) bool {
	if !c.IsPaired() {
		return false
	}

	// 获取配对的箱子
	pair := getPairFunc(c.pairX, c.pairZ)

	// 清除自身配对
	c.pairX = -1
	c.pairZ = -1
	c.NBT.Remove("pairx")
	c.NBT.Remove("pairz")

	// 清除对方配对
	if pair != nil {
		pair.pairX = -1
		pair.pairZ = -1
		pair.NBT.Remove("pairx")
		pair.NBT.Remove("pairz")
	}

	return true
}

// createPair 创建配对关系（内部方法）
// 对应 PHP Chest::createPair(Chest $tile)
func (c *Chest) createPair(other *Chest) {
	c.pairX = other.X
	c.pairZ = other.Z
	c.NBT.Set(nbt.NewIntTag("pairx", other.X))
	c.NBT.Set(nbt.NewIntTag("pairz", other.Z))

	other.pairX = c.X
	other.pairZ = c.Z
	other.NBT.Set(nbt.NewIntTag("pairx", c.X))
	other.NBT.Set(nbt.NewIntTag("pairz", c.Z))
}

// ---------- Spawnable 接口实现 ----------

// GetSpawnCompound 返回发送给客户端的 NBT 数据
// 对应 PHP Chest::getSpawnCompound()
func (c *Chest) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeChest))
	compound.Set(nbt.NewIntTag("x", c.X))
	compound.Set(nbt.NewIntTag("y", c.Y))
	compound.Set(nbt.NewIntTag("z", c.Z))

	// 如果已配对，附加配对坐标（客户端需要这些信息来渲染大箱子）
	if c.IsPaired() {
		compound.Set(nbt.NewIntTag("pairx", c.pairX))
		compound.Set(nbt.NewIntTag("pairz", c.pairZ))
	}

	// 自定义名称
	if c.HasCustomName() {
		compound.Set(nbt.NewStringTag("CustomName", c.GetCustomName()))
	}

	return compound
}

// UpdateCompoundTag 处理玩家发来的 NBT 更新
func (c *Chest) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	if nbtData.GetString("id") != TypeChest {
		return false
	}
	return true
}

// SpawnTo 向指定玩家发送此箱子的数据包
func (c *Chest) SpawnTo(sender PacketSender) bool {
	return SpawnTo(c, sender)
}

// SpawnToAll 向区块内所有玩家广播
func (c *Chest) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(c, broadcaster)
}

// ---------- Tile 接口补充 ----------

// SaveNBT 将当前状态写入 NBT
// 对应 PHP Chest::saveNBT()
func (c *Chest) SaveNBT() {
	c.BaseTile.SaveNBT()

	// 保存物品到 NBT
	c.ContainerBase.SaveItemsToNBT(c.NBT)

	// 保存自定义名称
	c.NameableBase.SaveNameToNBT(c.NBT)

	// 配对信息已在 createPair/Unpair 中同步到 NBT，无需重复
}

// Close 关闭箱子 Tile
// 对应 PHP Chest::close()
func (c *Chest) Close() {
	if c.IsClosed() {
		return
	}
	// 清空容器
	c.ClearAll()
	c.BaseTile.Close()
}

// ---------- 便捷方法 ----------

// GetContentsForDoubleChest 获取大箱子合并后的所有物品（54 槽位）
// 需要传入配对的 Chest 实例
func GetContentsForDoubleChest(left, right *Chest) []item.Item {
	result := make([]item.Item, ChestSize*2)
	for i := 0; i < ChestSize; i++ {
		result[i] = left.GetItem(i)
		result[ChestSize+i] = right.GetItem(i)
	}
	return result
}

// ---------- 注册 ----------

func init() {
	RegisterTile(TypeChest, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewChest(chunk, nbtData)
	})
}
