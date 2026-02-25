package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

// Dropper 投掷器 TileEntity
// 对应 PHP class Dropper extends Spawnable implements InventoryHolder, Container, Nameable
type Dropper struct {
	SpawnableBase
}

// DropperSize 投掷器容器大小（9格）
const DropperSize = 9

// NewDropper 创建 Dropper 实例
// 对应 PHP Dropper::__construct(FullChunk $chunk, CompoundTag $nbt)
func NewDropper(chunk *world.Chunk, nbtData *nbt.CompoundTag) *Dropper {
	d := &Dropper{}

	// 确保 Items 列表存在
	if nbtData.Get("Items") == nil {
		nbtData.Set(nbt.NewListTag("Items", nbt.TagCompound))
	}

	InitSpawnableBase(&d.SpawnableBase, TypeDropper, chunk, nbtData)
	return d
}

// GetName 返回名称（支持自定义名称）
// 对应 PHP Dropper::getName()
func (d *Dropper) GetName() string {
	if customName := d.NBT.GetString("CustomName"); customName != "" {
		return customName
	}
	return "Dropper"
}

// HasName 是否有自定义名称
// 对应 PHP Dropper::hasName()
func (d *Dropper) HasName() bool {
	return d.NBT.Get("CustomName") != nil && d.NBT.GetString("CustomName") != ""
}

// SetName 设置自定义名称
// 对应 PHP Dropper::setName($str)
func (d *Dropper) SetName(name string) {
	if name == "" {
		d.NBT.Remove("CustomName")
		return
	}
	d.NBT.Set(nbt.NewStringTag("CustomName", name))
}

// GetSize 返回容器大小
func (d *Dropper) GetSize() int {
	return DropperSize
}

// GetMotion 根据方块朝向返回投射方向向量
// 对应 PHP Dropper::getMotion()
// meta 为方块的 damage/meta 值
func GetDropperMotion(meta int) (x, y, z int) {
	switch meta {
	case 0: // SIDE_DOWN
		return 0, -1, 0
	case 1: // SIDE_UP
		return 0, 1, 0
	case 2: // SIDE_NORTH
		return 0, 0, -1
	case 3: // SIDE_SOUTH
		return 0, 0, 1
	case 4: // SIDE_WEST
		return -1, 0, 0
	case 5: // SIDE_EAST
		return 1, 0, 0
	default:
		return 0, 0, 0
	}
}

// GetSpawnCompound 返回发送给客户端的 NBT 数据
// 对应 PHP Dropper::getSpawnCompound()
func (d *Dropper) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeDropper))
	compound.Set(nbt.NewIntTag("x", d.X))
	compound.Set(nbt.NewIntTag("y", d.Y))
	compound.Set(nbt.NewIntTag("z", d.Z))

	if d.HasName() {
		compound.Set(nbt.NewStringTag("CustomName", d.GetName()))
	}

	return compound
}

// UpdateCompoundTag 处理客户端发来的更新
func (d *Dropper) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return false
}

// SpawnTo 向指定玩家发送数据包
func (d *Dropper) SpawnTo(sender PacketSender) bool {
	return SpawnTo(d, sender)
}

// SpawnToAll 向区块内所有玩家广播
func (d *Dropper) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(d, broadcaster)
}

func init() {
	RegisterTile(TypeDropper, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewDropper(chunk, nbtData)
	})
}
