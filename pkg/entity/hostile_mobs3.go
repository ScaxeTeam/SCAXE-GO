package entity

// hostile_mobs3.go — 剩余敌对/中立怪物: CaveSpider, PigZombie, Witch, Silverfish, Bat
// 对应 PHP: entity/CaveSpider.php, entity/PigZombie.php, entity/Witch.php,
//          entity/Silverfish.php, entity/Bat.php

import "math/rand"

// ============================================================
//                   CaveSpider (ID 40)
// ============================================================

const CaveSpiderNetworkID = 40

// NewCaveSpider 创建洞穴蜘蛛
// 对应 PHP CaveSpider.php
// - MaxHealth: 12, 攻击伤害: 2
// - 比普通蜘蛛小 (0.7×0.5)
// - 攻击附带中毒效果 (普通/困难难度)
func NewCaveSpider() *Monster {
	m := NewMonster(CaveSpiderNetworkID, "Cave Spider", 12, 0.7, 0.5, 2)
	m.DropExpMin = 5
	m.DropExpMax = 5
	return m
}

// CaveSpiderDrops 洞穴蜘蛛掉落物
// 与普通蜘蛛相同: 蛛丝/蜘蛛眼
func CaveSpiderDrops() []ZombieDropItem {
	const (
		SpiderEye = 375
		String    = 287
	)
	if rand.Intn(3) < 1 {
		return []ZombieDropItem{{ItemID: SpiderEye, Count: 1}}
	}
	return []ZombieDropItem{{ItemID: String, Count: 1 + rand.Intn(2)}}
}

// ============================================================
//                   PigZombie / ZombiePigman (ID 36)
// ============================================================

const PigZombieNetworkID = 36

// PigZombie 僵尸猪人（中立怪物，被攻击后变敌对）
type PigZombie struct {
	*Monster

	// Angry 是否被激怒
	Angry bool

	// AngerTimer 愤怒持续时间 (tick)
	AngerTimer int
}

// NewPigZombie 创建僵尸猪人
// 对应 PHP PigZombie.php
// - MaxHealth: 20, 攻击伤害: 5
// - 中立怪物，被攻击后敌对
// - 火焰免疫
// - 掉落: 金粒 1-3 + 10% 金剑
func NewPigZombie() *PigZombie {
	m := NewMonster(PigZombieNetworkID, "Zombie Pigman", 20, 0.6, 1.8, 5)
	m.DropExpMin = 5
	m.DropExpMax = 5

	return &PigZombie{
		Monster: m,
	}
}

// SetAngry 设置激怒状态
func (p *PigZombie) SetAngry(angry bool) {
	p.Angry = angry
	if angry {
		p.AngerTimer = 400 // 20秒
	} else {
		p.AngerTimer = 0
	}
}

// IsAngry 是否被激怒
func (p *PigZombie) IsAngry() bool {
	return p.Angry
}

// IsFireImmune 僵尸猪人免疫火焰
func (p *PigZombie) IsFireImmune() bool {
	return true
}

// TickAnger 愤怒 tick
func (p *PigZombie) TickAnger() {
	if p.Angry && p.AngerTimer > 0 {
		p.AngerTimer--
		if p.AngerTimer <= 0 {
			p.Angry = false
		}
	}
}

// PigZombieDrops 僵尸猪人掉落物
func PigZombieDrops() []ZombieDropItem {
	const (
		GoldNugget = 371
		GoldSword  = 283
	)

	drops := []ZombieDropItem{
		{ItemID: GoldNugget, Count: 1 + rand.Intn(3)}, // 1-3 金粒
	}

	// 10% 掉金剑
	if rand.Intn(100) < 10 {
		drops = append(drops, ZombieDropItem{ItemID: GoldSword, Count: 1})
	}

	return drops
}

// ============================================================
//                      Witch (ID 45)
// ============================================================

const WitchNetworkID = 45

// NewWitch 创建女巫
// 对应 PHP Witch.php
// - MaxHealth: 26, 远程攻击（投掷药水）
// - 掉落: 药水瓶/荧石粉/火药/红石/蜘蛛眼/糖/木棍
func NewWitch() *Monster {
	m := NewMonster(WitchNetworkID, "Witch", 26, 0.6, 1.8, 0)
	m.DropExpMin = 5
	m.DropExpMax = 5
	return m
}

// WitchDrops 女巫掉落物 (随机 1-3 种)
func WitchDrops() []ZombieDropItem {
	const (
		GlassBottle   = 374
		GlowstoneDust = 348
		Gunpowder     = 289
		Redstone      = 331
		SpiderEye     = 375
		Sugar         = 353
		Stick         = 280
	)

	possibleDrops := []int{GlassBottle, GlowstoneDust, Gunpowder, Redstone, SpiderEye, Sugar, Stick}
	drops := make([]ZombieDropItem, 0, 3)

	// 随机掉 1-3 种
	count := 1 + rand.Intn(3)
	for i := 0; i < count && i < len(possibleDrops); i++ {
		idx := rand.Intn(len(possibleDrops))
		drops = append(drops, ZombieDropItem{
			ItemID: possibleDrops[idx],
			Count:  1 + rand.Intn(2),
		})
	}

	return drops
}

// ============================================================
//                    Silverfish (ID 39)
// ============================================================

const SilverfishNetworkID = 39

// NewSilverfish 创建蠹虫
// 对应 PHP Silverfish.php
// - MaxHealth: 8, 攻击伤害: 1
// - 尺寸小 (0.4×0.3)
// - 被攻击时会召唤附近蠹虫方块中的同伴
func NewSilverfish() *Monster {
	m := NewMonster(SilverfishNetworkID, "Silverfish", 8, 0.4, 0.3, 1)
	m.DropExpMin = 5
	m.DropExpMax = 5
	return m
}

// SilverfishDrops 蠹虫无掉落物
func SilverfishDrops() []ZombieDropItem {
	return nil
}

// ============================================================
//                       Bat (ID 19)
// ============================================================

const BatNetworkID = 19

// Bat 蝙蝠（被动飞行生物）
type Bat struct {
	*Monster

	// Hanging 是否倒挂在方块上
	Hanging bool
}

// NewBat 创建蝙蝠
// 对应 PHP Bat.php
// - MaxHealth: 6, 无攻击
// - 飞行生物
// - 在光照等级 ≤4 的区域生成
func NewBat() *Bat {
	m := NewMonster(BatNetworkID, "Bat", 6, 0.5, 0.9, 0)
	m.DropExpMin = 0
	m.DropExpMax = 0

	return &Bat{
		Monster: m,
	}
}

// SetHanging 设置倒挂状态
func (b *Bat) SetHanging(hanging bool) {
	b.Hanging = hanging
}

// IsHanging 是否倒挂
func (b *Bat) IsHanging() bool {
	return b.Hanging
}

// BatDrops 蝙蝠无掉落物
func BatDrops() []ZombieDropItem {
	return nil
}
