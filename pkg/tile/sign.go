package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

// Sign 告示牌 TileEntity
// 对应 PHP class Sign extends Spawnable
// 存储 4 行文本（Text1 ~ Text4），支持客户端渲染和玩家编辑
type Sign struct {
	SpawnableBase

	text [4]string // Text1, Text2, Text3, Text4
}

// NewSign 创建 Sign 实例
// 对应 PHP Sign::__construct(FullChunk $chunk, CompoundTag $nbt)
func NewSign(chunk *world.Chunk, nbtData *nbt.CompoundTag) *Sign {
	s := &Sign{}

	// 确保 NBT 中存在 Text1-Text4（对应 PHP 构造函数中的默认值设置）
	for i, key := range signTextKeys {
		if val := nbtData.GetString(key); val != "" {
			s.text[i] = val
		} else {
			nbtData.Set(nbt.NewStringTag(key, ""))
		}
	}

	InitSpawnableBase(&s.SpawnableBase, TypeSign, chunk, nbtData)
	return s
}

// signTextKeys NBT 中告示牌文本的键名
var signTextKeys = [4]string{"Text1", "Text2", "Text3", "Text4"}

// ---------- 文本读写 ----------

// GetText 获取 4 行文本
// 对应 PHP Sign::getText()
func (s *Sign) GetText() [4]string {
	return s.text
}

// GetLine 获取指定行文本（0-3）
func (s *Sign) GetLine(index int) string {
	if index < 0 || index > 3 {
		return ""
	}
	return s.text[index]
}

// SetText 设置 4 行文本，并同步到 NBT
// 对应 PHP Sign::setText($line1, $line2, $line3, $line4)
// 注意：调用者应在之后调用 SpawnToAll 广播给附近玩家
func (s *Sign) SetText(line1, line2, line3, line4 string) {
	s.text[0] = line1
	s.text[1] = line2
	s.text[2] = line3
	s.text[3] = line4

	// 同步到 NBT
	for i, key := range signTextKeys {
		s.NBT.Set(nbt.NewStringTag(key, s.text[i]))
	}
}

// SetLine 设置指定行文本（0-3）
func (s *Sign) SetLine(index int, line string) {
	if index < 0 || index > 3 {
		return
	}
	s.text[index] = line
	s.NBT.Set(nbt.NewStringTag(signTextKeys[index], line))
}

// ---------- Spawnable 接口实现 ----------

// GetSpawnCompound 返回发送给客户端的 NBT 数据
// 对应 PHP Sign::getSpawnCompound()
func (s *Sign) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeSign))
	compound.Set(nbt.NewIntTag("x", s.X))
	compound.Set(nbt.NewIntTag("y", s.Y))
	compound.Set(nbt.NewIntTag("z", s.Z))
	for i, key := range signTextKeys {
		compound.Set(nbt.NewStringTag(key, s.text[i]))
	}
	return compound
}

// UpdateCompoundTag 处理玩家发来的告示牌编辑数据
// 对应 PHP Sign::updateCompoundTag(CompoundTag $nbt, Player $player)
// 简化版：不包含事件系统，仅做基本验证和文本更新
func (s *Sign) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	if nbtData.GetString("id") != TypeSign {
		return false
	}

	var lines [4]string
	for i, key := range signTextKeys {
		lines[i] = nbtData.GetString(key)
	}

	s.SetText(lines[0], lines[1], lines[2], lines[3])
	return true
}

// SpawnTo 向指定玩家发送此告示牌的数据包
func (s *Sign) SpawnTo(sender PacketSender) bool {
	return SpawnTo(s, sender)
}

// SpawnToAll 向区块内所有玩家广播
func (s *Sign) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(s, broadcaster)
}

// ---------- Tile 接口补充 ----------

// SaveNBT 将当前状态写入 NBT
// 对应 PHP Sign::saveNBT()
func (s *Sign) SaveNBT() {
	s.BaseTile.SaveNBT()

	// 同步文本到 NBT
	for i, key := range signTextKeys {
		s.NBT.Set(nbt.NewStringTag(key, s.text[i]))
	}

	// 对应 PHP: unset($this->namedtag->Creator)
	s.NBT.Remove("Creator")
}

// ---------- 注册 ----------

func init() {
	RegisterTile(TypeSign, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewSign(chunk, nbtData)
	})
}
