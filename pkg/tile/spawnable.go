package tile

import (
	"bytes"

	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/protocol"
	"github.com/scaxe/scaxe-go/pkg/world"
)

// Spawnable 接口表示可以将自身 NBT 数据发送给客户端的 Tile
// 对应 PHP abstract class Spawnable extends Tile
type Spawnable interface {
	Tile

	// GetSpawnCompound 返回发送给客户端的 NBT 数据
	// 每个具体 Tile 必须实现，用于客户端渲染
	GetSpawnCompound() *nbt.CompoundTag

	// UpdateCompoundTag 处理玩家发来的 NBT 更新（如编辑告示牌）
	// 返回 true 表示成功，false 会导致重新发送 Tile 数据给玩家
	UpdateCompoundTag(nbtData *nbt.CompoundTag) bool

	// SpawnTo 向指定玩家发送 BlockEntityDataPacket
	SpawnTo(sender PacketSender) bool

	// SpawnToAll 向区块内所有玩家发送
	SpawnToAll(broadcaster ChunkBroadcaster)
}

// PacketSender 是能接收数据包的对象（通常是 Player）
type PacketSender interface {
	SendPacket(pk protocol.DataPacket)
}

// ChunkBroadcaster 广播数据包给区块内的玩家
type ChunkBroadcaster interface {
	GetChunkPlayers(chunkX, chunkZ int32) []PacketSender
}

// ---------- SpawnableBase 实现 ----------

// SpawnableBase 为 Spawnable Tile 提供公共功能
// 具体 Tile（如 Sign, Chest）嵌入此结构体
type SpawnableBase struct {
	BaseTile
}

// InitSpawnableBase 初始化 SpawnableBase
func InitSpawnableBase(s *SpawnableBase, saveID string, chunk *world.Chunk, nbtData *nbt.CompoundTag) {
	InitBaseTile(&s.BaseTile, saveID, chunk, nbtData)
}

// CreateSpawnPacket 根据 SpawnCompound 创建 BlockEntityDataPacket
// 对应 PHP Spawnable::spawnTo() 中的序列化逻辑
func CreateSpawnPacket(s Spawnable) *protocol.BlockEntityDataPacket {
	compound := s.GetSpawnCompound()

	// 序列化 NBT 为 LittleEndian 字节
	buf := new(bytes.Buffer)
	writer := nbt.NewWriter(buf, nbt.LittleEndian)
	writer.WriteTag(compound)

	pk := protocol.NewBlockEntityDataPacket()
	x, y, z := s.GetPosition()
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.NBTData = buf.Bytes()
	return pk
}

// SpawnTo 向指定玩家发送此 Tile 的数据包
// 对应 PHP Spawnable::spawnTo(Player $player)
func SpawnTo(s Spawnable, sender PacketSender) bool {
	if s.IsClosed() {
		return false
	}
	pk := CreateSpawnPacket(s)
	sender.SendPacket(pk)
	return true
}

// SpawnToAll 向区块内所有玩家广播
// 对应 PHP Spawnable::spawnToAll()
func SpawnToAll(s Spawnable, broadcaster ChunkBroadcaster) {
	if s.IsClosed() {
		return
	}
	chunk := s.GetChunk()
	if chunk == nil {
		return
	}
	players := broadcaster.GetChunkPlayers(chunk.X, chunk.Z)
	for _, p := range players {
		SpawnTo(s, p)
	}
}

// DefaultUpdateCompoundTag 默认的 UpdateCompoundTag 实现（返回 false）
// 对应 PHP Spawnable::updateCompoundTag() { return false; }
func DefaultUpdateCompoundTag(_ *nbt.CompoundTag) bool {
	return false
}
