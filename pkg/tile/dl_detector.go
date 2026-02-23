package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type DLDetector struct {
	SpawnableBase
	lastType int32
}

const (
	TimeDay     int32 = 0
	TimeSunset  int32 = 12000
	TimeNight   int32 = 14000
	TimeSunrise int32 = 23000
	TimeFull    int32 = 24000
)

func NewDLDetector(chunk *world.Chunk, nbtData *nbt.CompoundTag) *DLDetector {
	d := &DLDetector{}
	InitSpawnableBase(&d.SpawnableBase, TypeDLDetector, chunk, nbtData)
	return d
}

func (d *DLDetector) GetName() string {
	return "Daylight Detector"
}

func (d *DLDetector) OnUpdate() bool {
	if d.IsClosed() {
		return false
	}

	return true
}

func GetLightByTime(time int32) int {

	if (time >= TimeDay && time <= TimeSunset) ||
		(time >= TimeSunrise && time <= TimeFull) {
		return 15
	}
	return 0
}

func (d *DLDetector) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeDLDetector))
	compound.Set(nbt.NewIntTag("x", d.X))
	compound.Set(nbt.NewIntTag("y", d.Y))
	compound.Set(nbt.NewIntTag("z", d.Z))
	return compound
}

func (d *DLDetector) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return false
}

func (d *DLDetector) SpawnTo(sender PacketSender) bool {
	return SpawnTo(d, sender)
}

func (d *DLDetector) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(d, broadcaster)
}

func init() {
	RegisterTile(TypeDLDetector, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewDLDetector(chunk, nbtData)
	})
}
