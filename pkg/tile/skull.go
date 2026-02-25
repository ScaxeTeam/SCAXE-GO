package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

// Skull 头颅方块 TileEntity
// 对应 PHP class Skull extends Spawnable
type Skull struct {
	SpawnableBase
}

// NewSkull 创建 Skull 实例
// 对应 PHP Skull::__construct(FullChunk $chunk, CompoundTag $nbt)
func NewSkull(chunk *world.Chunk, nbtData *nbt.CompoundTag) *Skull {
	s := &Skull{}

	if nbtData.Get("SkullType") == nil {
		nbtData.Set(nbt.NewByteTag("SkullType", 0))
	}
	if nbtData.Get("Rot") == nil {
		nbtData.Set(nbt.NewByteTag("Rot", 0))
	}

	InitSpawnableBase(&s.SpawnableBase, TypeSkull, chunk, nbtData)
	return s
}

// GetName 返回名称
func (s *Skull) GetName() string {
	return "Skull"
}

// GetSkullType 获取头颅类型
// 0=骷髅, 1=凋灵骷髅, 2=僵尸, 3=玩家, 4=苦力怕
func (s *Skull) GetSkullType() int8 {
	return s.NBT.GetByte("SkullType")
}

// SetSkullType 设置头颅类型
func (s *Skull) SetSkullType(skullType int8) {
	s.NBT.Set(nbt.NewByteTag("SkullType", skullType))
}

// GetRot 获取旋转角度
func (s *Skull) GetRot() int8 {
	return s.NBT.GetByte("Rot")
}

// SaveNBT 保存 NBT（移除 Creator 标签）
// 对应 PHP Skull::saveNBT()
func (s *Skull) SaveNBT() {
	s.BaseTile.SaveNBT()
	s.NBT.Remove("Creator")
}

// GetSpawnCompound 返回发送给客户端的 NBT 数据
// 对应 PHP Skull::getSpawnCompound()
func (s *Skull) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeSkull))
	compound.Set(nbt.NewIntTag("x", s.X))
	compound.Set(nbt.NewIntTag("y", s.Y))
	compound.Set(nbt.NewIntTag("z", s.Z))
	compound.Set(nbt.NewByteTag("SkullType", s.GetSkullType()))
	compound.Set(nbt.NewByteTag("Rot", s.GetRot()))
	return compound
}

// UpdateCompoundTag 处理客户端发来的更新
func (s *Skull) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return false
}

// SpawnTo 向指定玩家发送数据包
func (s *Skull) SpawnTo(sender PacketSender) bool {
	return SpawnTo(s, sender)
}

// SpawnToAll 向区块内所有玩家广播
func (s *Skull) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(s, broadcaster)
}

func init() {
	RegisterTile(TypeSkull, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewSkull(chunk, nbtData)
	})
}
