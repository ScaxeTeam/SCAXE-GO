package entity

// hostile_mobs2.go — 常见敌对怪物（下）: Enderman, Blaze, Ghast, Slime
// 对应 PHP: entity/Enderman.php, entity/Blaze.php, entity/Ghast.php, entity/Slime.php

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

// ============================================================
//                     Enderman (ID 38)
// ============================================================

const EndermanNetworkID = 38

// Enderman 末影人（扩展 Monster，增加搬运方块和颤抖状态）
type Enderman struct {
	*Monster

	// CarriedBlockID 搬运的方块ID
	CarriedBlockID int

	// CarriedBlockMeta 搬运的方块 meta
	CarriedBlockMeta int

	// Trembling 是否在颤抖（被看/被攻击时）
	Trembling bool
}

// NewEnderman 创建末影人
// 对应 PHP Enderman.php
// - MaxHealth: 40, 攻击伤害: 7
// - 可搬运方块
// - 被水伤害
func NewEnderman() *Enderman {
	m := NewMonster(EndermanNetworkID, "Enderman", 40, 0.6, 2.9, 7)
	m.DropExpMin = 5
	m.DropExpMax = 5

	return &Enderman{
		Monster: m,
	}
}

// SetBlockInHand 设置搬运的方块
// 对应 PHP Enderman::setBlockInHand()
func (e *Enderman) SetBlockInHand(blockID, blockMeta int) {
	e.CarriedBlockID = blockID
	e.CarriedBlockMeta = blockMeta
}

// GetBlockInHand 获取搬运的方块
func (e *Enderman) GetBlockInHand() (blockID, blockMeta int) {
	return e.CarriedBlockID, e.CarriedBlockMeta
}

// IsCarryingBlock 是否正在搬运方块
func (e *Enderman) IsCarryingBlock() bool {
	return e.CarriedBlockID > 0
}

// SetTremble 设置颤抖状态
func (e *Enderman) SetTremble(trembling bool) {
	e.Trembling = trembling
}

// SaveEndermanNBT 保存末影人 NBT
func (e *Enderman) SaveEndermanNBT() {
	e.Monster.SaveMonsterNBT()
	tag := e.Monster.Mob.Living.Entity.NamedTag
	tag.Set(nbt.NewShortTag("carried", int16(e.CarriedBlockID)))
	tag.Set(nbt.NewShortTag("carriedData", int16(e.CarriedBlockMeta)))
}

// LoadEndermanFromNBT 从 NBT 加载末影人数据
func (e *Enderman) LoadEndermanFromNBT() {
	e.Monster.LoadMonsterFromNBT()
	tag := e.Monster.Mob.Living.Entity.NamedTag
	if tag != nil {
		e.CarriedBlockID = int(tag.GetShort("carried"))
		e.CarriedBlockMeta = int(tag.GetShort("carriedData"))
	}
}

// EndermanDrops 末影人掉落物
func EndermanDrops() []ZombieDropItem {
	const EnderPearl = 368
	// 0-1 末影珍珠
	count := rand.Intn(2)
	if count == 0 {
		return nil
	}
	return []ZombieDropItem{{ItemID: EnderPearl, Count: count}}
}

// ============================================================
//                       Blaze (ID 43)
// ============================================================

const BlazeNetworkID = 43

// Blaze 烈焰人（飞行怪物，火焰免疫，水中受伤）
type Blaze struct {
	*Monster

	// Charging 是否正在蓄力发射
	Charging bool
}

// NewBlaze 创建烈焰人
// 对应 PHP Blaze.php
// - MaxHealth: 20, 远程攻击（火球）
// - 免疫火焰/岩浆
// - 水中每 10 tick 受伤
func NewBlaze() *Blaze {
	m := NewMonster(BlazeNetworkID, "Blaze", 20, 0.6, 1.8, 6)
	m.DropExpMin = 10
	m.DropExpMax = 10

	return &Blaze{
		Monster: m,
	}
}

// IsFireImmune 烈焰人免疫火焰伤害
func (b *Blaze) IsFireImmune() bool {
	return true
}

// SetCharging 设置蓄力状态
func (b *Blaze) SetCharging(charging bool) {
	b.Charging = charging
}

// IsCharging 是否正在蓄力
func (b *Blaze) IsCharging() bool {
	return b.Charging
}

// BlazeDrops 烈焰人掉落物
// 对应 PHP Blaze::getDrops()
func BlazeDrops() []ZombieDropItem {
	const BlazeRod = 369
	count := rand.Intn(2) // 0-1
	if count == 0 {
		return nil
	}
	return []ZombieDropItem{{ItemID: BlazeRod, Count: count}}
}

// ============================================================
//                       Ghast (ID 41)
// ============================================================

const GhastNetworkID = 41

// Ghast 恶魂（大型飞行怪物，发射火球）
type Ghast struct {
	*Monster

	// Charging 是否正在蓄力发射火球
	Charging bool
}

// NewGhast 创建恶魂
// 对应 PHP Ghast.php
// - MaxHealth: 10, 远程攻击（火球，伤害由爆炸决定）
// - 尺寸: 4×4×4 (大型)
func NewGhast() *Ghast {
	m := NewMonster(GhastNetworkID, "Ghast", 10, 4.0, 4.0, 0)
	m.DropExpMin = 5
	m.DropExpMax = 5

	return &Ghast{
		Monster: m,
	}
}

// SetCharging 设置蓄力状态
func (g *Ghast) SetCharging(charging bool) {
	g.Charging = charging
}

// IsCharging 是否正在蓄力
func (g *Ghast) IsCharging() bool {
	return g.Charging
}

// GhastDrops 恶魂掉落物
func GhastDrops() []ZombieDropItem {
	const (
		GhastTear = 370
		Gunpowder = 289
	)
	drops := []ZombieDropItem{
		{ItemID: Gunpowder, Count: rand.Intn(2) + 1}, // 1-2 火药
	}
	// 50% 恶魂之泪
	if rand.Intn(2) == 0 {
		drops = append(drops, ZombieDropItem{ItemID: GhastTear, Count: 1})
	}
	return drops
}

// ============================================================
//                       Slime (ID 37)
// ============================================================

const SlimeNetworkID = 37

// Slime 史莱姆（可分裂的怪物）
type Slime struct {
	*Monster

	// Size 大小(1-4), 1=最小, 4=最大
	Size int
}

// NewSlime 创建史莱姆
// 对应 PHP Slime.php
// - 大小随机 1-4
// - 死亡时如果 size > 1，分裂为 4 个 size-1 的小史莱姆
func NewSlime() *Slime {
	size := 1 + rand.Intn(4) // 1-4
	m := NewMonster(SlimeNetworkID, "Slime", slimeHealthForSize(size), 0.6, 0.6, slimeDamageForSize(size))
	m.DropExpMin = 1
	m.DropExpMax = 4

	return &Slime{
		Monster: m,
		Size:    size,
	}
}

// NewSlimeWithSize 创建指定大小的史莱姆
func NewSlimeWithSize(size int) *Slime {
	if size < 1 {
		size = 1
	}
	if size > 4 {
		size = 4
	}
	m := NewMonster(SlimeNetworkID, "Slime", slimeHealthForSize(size), 0.6, 0.6, slimeDamageForSize(size))
	m.DropExpMin = 1
	m.DropExpMax = 4

	return &Slime{
		Monster: m,
		Size:    size,
	}
}

// slimeHealthForSize 不同大小的生命值
func slimeHealthForSize(size int) int {
	switch size {
	case 1:
		return 1
	case 2:
		return 4
	case 3:
		return 8
	case 4:
		return 16
	default:
		return 1
	}
}

// slimeDamageForSize 不同大小的攻击伤害
func slimeDamageForSize(size int) int {
	switch size {
	case 1:
		return 0 // 最小的不攻击
	case 2:
		return 2
	case 3:
		return 3
	case 4:
		return 4
	default:
		return 0
	}
}

// GetSize 获取大小
func (s *Slime) GetSize() int {
	return s.Size
}

// SetSize 设置大小
func (s *Slime) SetSize(size int) {
	s.Size = size
}

// ShouldSplit 死亡时是否应分裂（size > 1）
func (s *Slime) ShouldSplit() bool {
	return s.Size > 1
}

// GetSplitSize 分裂后的子代大小
func (s *Slime) GetSplitSize() int {
	return s.Size - 1
}

// GetSplitCount 分裂数量（始终为 4）
func (s *Slime) GetSplitCount() int {
	return 4
}

// SlimeDrops 史莱姆掉落物
// 对应 PHP Slime::getDrops()
// 只有最小的 (size=1) 才掉粘液球
func SlimeDrops(size int) []ZombieDropItem {
	const Slimeball = 341
	if size == 1 {
		return []ZombieDropItem{{ItemID: Slimeball, Count: 1}}
	}
	return nil
}
