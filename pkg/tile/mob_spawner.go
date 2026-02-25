package tile

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

// MobSpawner 刷怪笼 TileEntity
// 对应 PHP class MobSpawner extends Spawnable
type MobSpawner struct {
	SpawnableBase
}

// NewMobSpawner 创建 MobSpawner 实例
// 对应 PHP MobSpawner::__construct(FullChunk $chunk, CompoundTag $nbt)
func NewMobSpawner(chunk *world.Chunk, nbtData *nbt.CompoundTag) *MobSpawner {
	s := &MobSpawner{}

	// 设置默认 NBT 值
	if nbtData.Get("EntityId") == nil {
		nbtData.Set(nbt.NewIntTag("EntityId", 0))
	}
	if nbtData.Get("SpawnCount") == nil {
		nbtData.Set(nbt.NewIntTag("SpawnCount", 4))
	}
	if nbtData.Get("SpawnRange") == nil {
		nbtData.Set(nbt.NewIntTag("SpawnRange", 4))
	}
	if nbtData.Get("MinSpawnDelay") == nil {
		nbtData.Set(nbt.NewIntTag("MinSpawnDelay", 200))
	}
	if nbtData.Get("MaxSpawnDelay") == nil {
		nbtData.Set(nbt.NewIntTag("MaxSpawnDelay", 799))
	}
	if nbtData.Get("Delay") == nil {
		minDelay := nbtData.GetInt("MinSpawnDelay")
		maxDelay := nbtData.GetInt("MaxSpawnDelay")
		delay := minDelay
		if maxDelay > minDelay {
			delay = minDelay + int32(rand.Intn(int(maxDelay-minDelay+1)))
		}
		nbtData.Set(nbt.NewIntTag("Delay", delay))
	}

	InitSpawnableBase(&s.SpawnableBase, TypeMobSpawner, chunk, nbtData)
	return s
}

// GetName 返回名称
func (s *MobSpawner) GetName() string {
	return "Monster Spawner"
}

// GetEntityId 获取要生成的实体类型 ID
func (s *MobSpawner) GetEntityId() int32 {
	return s.NBT.GetInt("EntityId")
}

// SetEntityId 设置要生成的实体类型 ID
func (s *MobSpawner) SetEntityId(id int32) {
	s.NBT.Set(nbt.NewIntTag("EntityId", id))
}

// GetSpawnCount 获取每次生成数量
func (s *MobSpawner) GetSpawnCount() int32 {
	return s.NBT.GetInt("SpawnCount")
}

// SetSpawnCount 设置每次生成数量
func (s *MobSpawner) SetSpawnCount(value int32) {
	s.NBT.Set(nbt.NewIntTag("SpawnCount", value))
}

// GetSpawnRange 获取生成范围
func (s *MobSpawner) GetSpawnRange() int32 {
	return s.NBT.GetInt("SpawnRange")
}

// SetSpawnRange 设置生成范围
func (s *MobSpawner) SetSpawnRange(value int32) {
	s.NBT.Set(nbt.NewIntTag("SpawnRange", value))
}

// GetDelay 获取当前延迟
func (s *MobSpawner) GetDelay() int32 {
	return s.NBT.GetInt("Delay")
}

// SetDelay 设置当前延迟
func (s *MobSpawner) SetDelay(value int32) {
	s.NBT.Set(nbt.NewIntTag("Delay", value))
}

// GetMinSpawnDelay 获取最小生成延迟
func (s *MobSpawner) GetMinSpawnDelay() int32 {
	return s.NBT.GetInt("MinSpawnDelay")
}

// SetMinSpawnDelay 设置最小生成延迟
func (s *MobSpawner) SetMinSpawnDelay(value int32) {
	s.NBT.Set(nbt.NewIntTag("MinSpawnDelay", value))
}

// GetMaxSpawnDelay 获取最大生成延迟
func (s *MobSpawner) GetMaxSpawnDelay() int32 {
	return s.NBT.GetInt("MaxSpawnDelay")
}

// SetMaxSpawnDelay 设置最大生成延迟
func (s *MobSpawner) SetMaxSpawnDelay(value int32) {
	s.NBT.Set(nbt.NewIntTag("MaxSpawnDelay", value))
}

// OnUpdate 每 tick 更新
// 对应 PHP MobSpawner::onUpdate()
// 注意：实际的生物生成逻辑需要 Level 层配合，这里只处理延迟计数器
func (s *MobSpawner) OnUpdate() bool {
	if s.IsClosed() {
		return false
	}

	if s.GetEntityId() == 0 {
		return true // 没有设置实体类型，继续 tick 但不生成
	}

	delay := s.GetDelay()
	if delay > 0 {
		s.SetDelay(delay - 1)
	}
	// 实际生成逻辑由 Level tick 层调用 SpawnMobs() 驱动

	return true
}

// GetSpawnCompound 返回发送给客户端的 NBT 数据
// 对应 PHP MobSpawner::getSpawnCompound()
func (s *MobSpawner) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeMobSpawner))
	compound.Set(nbt.NewIntTag("x", s.X))
	compound.Set(nbt.NewIntTag("y", s.Y))
	compound.Set(nbt.NewIntTag("z", s.Z))
	compound.Set(nbt.NewIntTag("EntityId", s.GetEntityId()))
	return compound
}

// UpdateCompoundTag 处理客户端发来的更新
func (s *MobSpawner) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return false
}

// SpawnTo 向指定玩家发送数据包
func (s *MobSpawner) SpawnTo(sender PacketSender) bool {
	return SpawnTo(s, sender)
}

// SpawnToAll 向区块内所有玩家广播
func (s *MobSpawner) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(s, broadcaster)
}

func init() {
	RegisterTile(TypeMobSpawner, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewMobSpawner(chunk, nbtData)
	})
}
