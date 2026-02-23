package tile

import (
	"bytes"

	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/protocol"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Spawnable interface {
	Tile

	GetSpawnCompound() *nbt.CompoundTag

	UpdateCompoundTag(nbtData *nbt.CompoundTag) bool

	SpawnTo(sender PacketSender) bool

	SpawnToAll(broadcaster ChunkBroadcaster)
}

type PacketSender interface {
	SendPacket(pk protocol.DataPacket)
}

type ChunkBroadcaster interface {
	GetChunkPlayers(chunkX, chunkZ int32) []PacketSender
}

type SpawnableBase struct {
	BaseTile
}

func InitSpawnableBase(s *SpawnableBase, saveID string, chunk *world.Chunk, nbtData *nbt.CompoundTag) {
	InitBaseTile(&s.BaseTile, saveID, chunk, nbtData)
}

func CreateSpawnPacket(s Spawnable) *protocol.BlockEntityDataPacket {
	compound := s.GetSpawnCompound()

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

func SpawnTo(s Spawnable, sender PacketSender) bool {
	if s.IsClosed() {
		return false
	}
	pk := CreateSpawnPacket(s)
	sender.SendPacket(pk)
	return true
}

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

func DefaultUpdateCompoundTag(_ *nbt.CompoundTag) bool {
	return false
}
