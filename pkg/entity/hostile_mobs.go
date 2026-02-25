package entity

// hostile_mobs.go — 常见敌对怪物（上）: Zombie, Skeleton, Creeper, Spider
// 对应 PHP: entity/Zombie.php, entity/Skeleton.php, entity/Creeper.php, entity/Spider.php

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

// ============================================================
//                       Zombie (ID 32)
// ============================================================

const ZombieNetworkID = 32

// NewZombie 创建僵尸
// 对应 PHP Zombie.php
// - MaxHealth: 20, 攻击伤害: 4 (普通难度)
// - 白天燃烧
// - 可穿戴随机盔甲（难度相关）
// - 掉落: 腐肉 1-2 + 10% 概率掉铁锭/胡萝卜/土豆
func NewZombie() *Monster {
	m := NewMonster(ZombieNetworkID, "Zombie", 20, 0.6, 1.8, 4)
	m.DropExpMin = 5
	m.DropExpMax = 5
	return m
}

// ZombieDropItem 僵尸掉落物结构
type ZombieDropItem struct {
	ItemID int
	Meta   int
	Count  int
}

// ZombieDrops 计算僵尸掉落物
// 对应 PHP Zombie::getDrops()
//
// 返回主掉落（腐肉）+ 可选稀有掉落
func ZombieDrops() []ZombieDropItem {
	const (
		RottenFlesh = 367
		IronIngot   = 265
		Carrot      = 391
		Potato      = 392
	)

	drops := []ZombieDropItem{
		{ItemID: RottenFlesh, Count: 1 + rand.Intn(2)}, // 1-2
	}

	// 10% 稀有掉落
	if rand.Intn(100) < 10 {
		switch rand.Intn(3) {
		case 0:
			drops = append(drops, ZombieDropItem{ItemID: IronIngot, Count: 1})
		case 1:
			drops = append(drops, ZombieDropItem{ItemID: Carrot, Count: 1})
		case 2:
			drops = append(drops, ZombieDropItem{ItemID: Potato, Count: 1})
		}
	}

	return drops
}

// ============================================================
//                      Skeleton (ID 34)
// ============================================================

const SkeletonNetworkID = 34

// NewSkeleton 创建骷髅
// 对应 PHP Skeleton.php
// - MaxHealth: 20, 远程攻击（弓箭）
// - 白天燃烧
// - 手持弓
// - 可穿戴随机盔甲（难度相关）
// - 掉落: 骨头 1-2 + 盔甲掉落 25%
func NewSkeleton() *Monster {
	m := NewMonster(SkeletonNetworkID, "Skeleton", 20, 0.6, 1.8, 4)
	m.DropExpMin = 5
	m.DropExpMax = 5
	return m
}

// SkeletonDrops 计算骷髅掉落物
// 对应 PHP Skeleton::getDrops()
func SkeletonDrops() []ZombieDropItem {
	const (
		Bone  = 352
		Arrow = 262
	)

	drops := []ZombieDropItem{
		{ItemID: Bone, Count: 1 + rand.Intn(2)}, // 1-2 骨头
	}

	// 50% 概率掉箭
	if rand.Intn(2) == 0 {
		drops = append(drops, ZombieDropItem{ItemID: Arrow, Count: 1 + rand.Intn(2)})
	}

	return drops
}

// ============================================================
//                      Creeper (ID 33)
// ============================================================

const CreeperNetworkID = 33

// Creeper 苦力怕（扩展 Monster，增加 Powered/Swell 状态）
type Creeper struct {
	*Monster

	// Powered 是否被闪电充能
	Powered bool

	// SwellDirection 膨胀方向（0=未膨胀，1=膨胀中）
	SwellDirection int

	// SwellCounter 膨胀计数器（0-30，到 30 爆炸）
	SwellCounter int
}

// NewCreeper 创建苦力怕
// 对应 PHP Creeper.php
// - MaxHealth: 20, 攻击方式: 爆炸（非近战）
// - 掉落: 火药
func NewCreeper() *Creeper {
	m := NewMonster(CreeperNetworkID, "Creeper", 20, 0.6, 1.8, 0)
	m.DropExpMin = 5
	m.DropExpMax = 5

	return &Creeper{
		Monster: m,
	}
}

// SetPowered 设置充能状态
// 对应 PHP Creeper::setPowered()
func (c *Creeper) SetPowered(powered bool) {
	c.Powered = powered
}

// IsPowered 是否被充能
func (c *Creeper) IsPowered() bool {
	return c.Powered
}

// SetSwelled 设置膨胀状态
// 对应 PHP Creeper::setSwelled()
func (c *Creeper) SetSwelled(swelled bool) {
	if swelled {
		c.SwellDirection = 1
	} else {
		c.SwellDirection = 0
	}
}

// IsSwelled 是否正在膨胀
func (c *Creeper) IsSwelled() bool {
	return c.SwellDirection == 1
}

// CreeperTickResult Creeper tick 结果
type CreeperTickResult struct {
	HasUpdate     bool
	ShouldExplode bool // 膨胀计数器到 30 时应爆炸
}

// TickCreeper 苦力怕 tick 逻辑
// 对应 PHP Creeper::entityBaseTick()
func (c *Creeper) TickCreeper() CreeperTickResult {
	result := CreeperTickResult{}

	if c.IsSwelled() {
		// 膨胀中
		if c.SwellCounter < 30 {
			c.SwellCounter++
			result.HasUpdate = true
		} else {
			// 到 30 tick，应爆炸
			c.SetSwelled(false)
			result.ShouldExplode = true
			result.HasUpdate = true
		}
	} else if c.SwellCounter > 0 {
		// 取消膨胀后慢慢回复
		c.SwellCounter--
		result.HasUpdate = true
	}

	return result
}

// GetExplosionPower 获取爆炸威力
// 普通: 3, 充能: 6
func (c *Creeper) GetExplosionPower() float64 {
	if c.Powered {
		return 6.0
	}
	return 3.0
}

// SaveCreeperNBT 保存苦力怕 NBT
func (c *Creeper) SaveCreeperNBT() {
	c.Monster.SaveMonsterNBT()
	powered := int8(0)
	if c.Powered {
		powered = 1
	}
	c.Monster.Mob.Living.Entity.NamedTag.Set(nbt.NewByteTag("powered", powered))
}

// LoadCreeperFromNBT 从 NBT 加载苦力怕
func (c *Creeper) LoadCreeperFromNBT() {
	c.Monster.LoadMonsterFromNBT()
	tag := c.Monster.Mob.Living.Entity.NamedTag
	if tag != nil && tag.GetByte("powered") == 1 {
		c.Powered = true
	}
}

// CreeperDrops 苦力怕掉落物
// 对应 PHP Creeper::getDrops()
func CreeperDrops() []ZombieDropItem {
	const Gunpowder = 289
	return []ZombieDropItem{
		{ItemID: Gunpowder, Count: 1},
	}
}

// ============================================================
//                       Spider (ID 35)
// ============================================================

const SpiderNetworkID = 35

// NewSpider 创建蜘蛛
// 对应 PHP Spider.php
// - MaxHealth: 16, 攻击伤害: 3
// - 掉落: 蛛丝 或 蜘蛛眼 (33%)
func NewSpider() *Monster {
	m := NewMonster(SpiderNetworkID, "Spider", 16, 1.4, 0.9, 3)
	m.DropExpMin = 5
	m.DropExpMax = 5
	return m
}

// SpiderDrops 蜘蛛掉落物
// 对应 PHP Spider::getDrops()
func SpiderDrops() []ZombieDropItem {
	const (
		SpiderEye = 375
		String    = 287
	)

	// 33% 蜘蛛眼, 67% 蛛丝
	if rand.Intn(3) < 1 {
		return []ZombieDropItem{{ItemID: SpiderEye, Count: 1}}
	}
	return []ZombieDropItem{{ItemID: String, Count: 1 + rand.Intn(2)}}
}
