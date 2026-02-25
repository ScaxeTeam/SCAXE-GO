package tile

// beacon.go — 信标 TileEntity
// MCPE 0.14 中 Beacon 作为 TileEntity 存在
// Spawnable + Nameable: 存储信标等级和选择的效果
// 完整的信标效果逻辑在 Phase 3 实现

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	TypeBeacon = "Beacon"
)

// Beacon 信标 TileEntity
type Beacon struct {
	SpawnableBase
	NameableBase

	// Primary 和 Secondary 效果 ID (药水效果ID)
	Primary   int32
	Secondary int32
}

// NewBeacon 创建信标实例
func NewBeacon(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
	b := &Beacon{}
	InitSpawnableBase(&b.SpawnableBase, TypeBeacon, chunk, nbtData)
	b.NameableBase.LoadNameFromNBT(nbtData)

	b.Primary = nbtData.GetInt("primary")
	b.Secondary = nbtData.GetInt("secondary")

	return b
}

func (b *Beacon) GetName() string {
	if b.HasCustomName() {
		return b.GetCustomName()
	}
	return "Beacon"
}

// OnUpdate 信标 tick — 检查金字塔等级并应用效果
// TODO: 完整实现需要检测下方金字塔方块 + 应用药水效果
func (b *Beacon) OnUpdate() bool {
	return false
}

// GetSpawnCompound 客户端渲染数据
func (b *Beacon) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeBeacon))
	compound.Set(nbt.NewIntTag("x", b.X))
	compound.Set(nbt.NewIntTag("y", b.Y))
	compound.Set(nbt.NewIntTag("z", b.Z))
	compound.Set(nbt.NewIntTag("primary", b.Primary))
	compound.Set(nbt.NewIntTag("secondary", b.Secondary))

	if b.HasCustomName() {
		compound.Set(nbt.NewStringTag("CustomName", b.GetCustomName()))
	}
	return compound
}

func (b *Beacon) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	// 客户端可以发送信标效果选择
	b.Primary = nbtData.GetInt("primary")
	b.Secondary = nbtData.GetInt("secondary")
	return true
}

func (b *Beacon) SpawnTo(sender PacketSender) bool {
	return SpawnTo(b, sender)
}

func (b *Beacon) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(b, broadcaster)
}

func (b *Beacon) SaveNBT() {
	b.SpawnableBase.SaveNBT()
	b.NameableBase.SaveNameToNBT(b.NBT)
	b.NBT.Set(nbt.NewIntTag("primary", b.Primary))
	b.NBT.Set(nbt.NewIntTag("secondary", b.Secondary))
}

func init() {
	RegisterTile(TypeBeacon, NewBeacon)
}
