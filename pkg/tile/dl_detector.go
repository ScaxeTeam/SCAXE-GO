package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

// DLDetector 日光传感器 TileEntity
// 对应 PHP class DLDetector extends Spawnable
type DLDetector struct {
	SpawnableBase
	lastType int32
}

// TIME 常量，对应 PHP Level::TIME_*
const (
	TimeDay     int32 = 0
	TimeSunset  int32 = 12000
	TimeNight   int32 = 14000
	TimeSunrise int32 = 23000
	TimeFull    int32 = 24000
)

// NewDLDetector 创建 DLDetector 实例
// 对应 PHP DLDetector::__construct(FullChunk $chunk, CompoundTag $nbt)
func NewDLDetector(chunk *world.Chunk, nbtData *nbt.CompoundTag) *DLDetector {
	d := &DLDetector{}
	InitSpawnableBase(&d.SpawnableBase, TypeDLDetector, chunk, nbtData)
	return d
}

// GetName 返回名称
func (d *DLDetector) GetName() string {
	return "Daylight Detector"
}

// OnUpdate 每 tick 更新（实际每 3 tick 检测一次时间变化）
// 对应 PHP DLDetector::onUpdate()
func (d *DLDetector) OnUpdate() bool {
	if d.IsClosed() {
		return false
	}
	// 实际的周期性检测和红石信号更新由 Level tick 层驱动
	return true
}

// GetLightByTime 根据游戏时间返回光照等级
// 对应 PHP DLDetector::getLightByTime()
func GetLightByTime(time int32) int {
	// 白天时段返回15，夜晚返回0
	if (time >= TimeDay && time <= TimeSunset) ||
		(time >= TimeSunrise && time <= TimeFull) {
		return 15
	}
	return 0
}

// GetSpawnCompound 返回发送给客户端的 NBT 数据
// 对应 PHP DLDetector::getSpawnCompound()
func (d *DLDetector) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeDLDetector))
	compound.Set(nbt.NewIntTag("x", d.X))
	compound.Set(nbt.NewIntTag("y", d.Y))
	compound.Set(nbt.NewIntTag("z", d.Z))
	return compound
}

// UpdateCompoundTag 处理客户端发来的更新
func (d *DLDetector) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return false
}

// SpawnTo 向指定玩家发送数据包
func (d *DLDetector) SpawnTo(sender PacketSender) bool {
	return SpawnTo(d, sender)
}

// SpawnToAll 向区块内所有玩家广播
func (d *DLDetector) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(d, broadcaster)
}

func init() {
	RegisterTile(TypeDLDetector, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewDLDetector(chunk, nbtData)
	})
}
