package tile

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)
type MobSpawner struct {
	SpawnableBase
}
func NewMobSpawner(chunk *world.Chunk, nbtData *nbt.CompoundTag) *MobSpawner {
	s := &MobSpawner{}
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
func (s *MobSpawner) GetName() string {
	return "Monster Spawner"
}
func (s *MobSpawner) GetEntityId() int32 {
	return s.NBT.GetInt("EntityId")
}
func (s *MobSpawner) SetEntityId(id int32) {
	s.NBT.Set(nbt.NewIntTag("EntityId", id))
}
func (s *MobSpawner) GetSpawnCount() int32 {
	return s.NBT.GetInt("SpawnCount")
}
func (s *MobSpawner) SetSpawnCount(value int32) {
	s.NBT.Set(nbt.NewIntTag("SpawnCount", value))
}
func (s *MobSpawner) GetSpawnRange() int32 {
	return s.NBT.GetInt("SpawnRange")
}
func (s *MobSpawner) SetSpawnRange(value int32) {
	s.NBT.Set(nbt.NewIntTag("SpawnRange", value))
}
func (s *MobSpawner) GetDelay() int32 {
	return s.NBT.GetInt("Delay")
}
func (s *MobSpawner) SetDelay(value int32) {
	s.NBT.Set(nbt.NewIntTag("Delay", value))
}
func (s *MobSpawner) GetMinSpawnDelay() int32 {
	return s.NBT.GetInt("MinSpawnDelay")
}
func (s *MobSpawner) SetMinSpawnDelay(value int32) {
	s.NBT.Set(nbt.NewIntTag("MinSpawnDelay", value))
}
func (s *MobSpawner) GetMaxSpawnDelay() int32 {
	return s.NBT.GetInt("MaxSpawnDelay")
}
func (s *MobSpawner) SetMaxSpawnDelay(value int32) {
	s.NBT.Set(nbt.NewIntTag("MaxSpawnDelay", value))
}
func (s *MobSpawner) OnUpdate() bool {
	if s.IsClosed() {
		return false
	}

	if s.GetEntityId() == 0 {
		return true
	}

	delay := s.GetDelay()
	if delay > 0 {
		s.SetDelay(delay - 1)
	}

	return true
}
func (s *MobSpawner) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeMobSpawner))
	compound.Set(nbt.NewIntTag("x", s.X))
	compound.Set(nbt.NewIntTag("y", s.Y))
	compound.Set(nbt.NewIntTag("z", s.Z))
	compound.Set(nbt.NewIntTag("EntityId", s.GetEntityId()))
	return compound
}
func (s *MobSpawner) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return false
}
func (s *MobSpawner) SpawnTo(sender PacketSender) bool {
	return SpawnTo(s, sender)
}
func (s *MobSpawner) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(s, broadcaster)
}

func init() {
	RegisterTile(TypeMobSpawner, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewMobSpawner(chunk, nbtData)
	})
}
