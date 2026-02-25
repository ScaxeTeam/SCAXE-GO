package entity

// special_entities.go — 特殊实体: FishingHook, ExperienceOrb, Fireballs, LavaSlime (岩浆怪)
// 对应 PHP: entity/FishingHook.php, entity/ExperienceOrb.php,
//          entity/BigFireball.php, entity/SmallFireball.php, entity/LavaSlime.php

import "math/rand"

// ============================================================
//                  FishingHook (ID 77)
// ============================================================

const FishingHookNetworkID = 77

// FishingHook 钓鱼钩实体
type FishingHook struct {
	Entity

	// OwnerID 抛出钓钩的玩家 ID
	OwnerID int64

	// WaitTimer 等待咬钩倒计时 (tick)
	WaitTimer int

	// Hooked 是否已咬钩
	Hooked bool
}

// NewFishingHook 创建钓鱼钩
// - 抛出后受重力下落
// - 落入水中后开始等待咬钩
// - 咬钩后玩家收竿获得物品
func NewFishingHook(ownerID int64) *FishingHook {
	fh := &FishingHook{
		OwnerID:   ownerID,
		WaitTimer: 100 + rand.Intn(400), // 5-25 秒随机
	}
	fh.Entity.NetworkID = FishingHookNetworkID
	fh.Entity.Width = 0.25
	fh.Entity.Height = 0.25
	fh.Entity.Gravity = 0.04
	fh.Entity.MaxHealth = 1
	fh.Entity.Health = 1
	return fh
}

// TickFishingHook 钓鱼钩 tick
func (fh *FishingHook) TickFishingHook() bool {
	if fh.WaitTimer > 0 {
		fh.WaitTimer--
		if fh.WaitTimer <= 0 {
			fh.Hooked = true
			return true // 通知玩家咬钩了
		}
	}
	return false
}

// IsHooked 是否已咬钩
func (fh *FishingHook) IsHooked() bool {
	return fh.Hooked
}

// ============================================================
//                ExperienceOrb (ID 69)
// ============================================================

const ExperienceOrbNetworkID = 69

// ExperienceOrb 经验球实体
type ExperienceOrb struct {
	Entity

	// Experience 包含的经验值
	Experience int

	// PickupDelay 拾取延迟 (tick)
	PickupDelay int

	// Age 存在时间 (tick), 超过 6000 消失
	Age int
}

const (
	ExperienceOrbMaxAge     = 6000 // 5 分钟
	ExperienceOrbPickupDist = 2.0  // 拾取距离
)

// NewExperienceOrb 创建经验球
func NewExperienceOrb(experience int) *ExperienceOrb {
	orb := &ExperienceOrb{
		Experience:  experience,
		PickupDelay: 10, // 0.5 秒后可拾取
	}
	orb.Entity.NetworkID = ExperienceOrbNetworkID
	orb.Entity.Width = 0.25
	orb.Entity.Height = 0.25
	orb.Entity.Gravity = 0.04
	orb.Entity.MaxHealth = 1
	orb.Entity.Health = 1
	return orb
}

// TickExperienceOrb 经验球 tick
// 返回是否应消除
func (o *ExperienceOrb) TickExperienceOrb() bool {
	o.Age++
	if o.Age >= ExperienceOrbMaxAge {
		return true // 超时消失
	}
	if o.PickupDelay > 0 {
		o.PickupDelay--
	}
	return false
}

// CanPickup 是否可被拾取
func (o *ExperienceOrb) CanPickup() bool {
	return o.PickupDelay <= 0
}

// GetExperience 获取经验值
func (o *ExperienceOrb) GetExperience() int {
	return o.Experience
}

// ============================================================
//                 BigFireball (ID 85) — 恶魂火球
// ============================================================

const BigFireballNetworkID = 85

// BigFireball 恶魂火球（大型，落地产生爆炸）
type BigFireball struct {
	Entity

	// ShooterID 发射者实体 ID
	ShooterID int64

	// ExplosionPower 爆炸威力
	ExplosionPower float64
}

// NewBigFireball 创建恶魂火球
func NewBigFireball(shooterID int64) *BigFireball {
	fb := &BigFireball{
		ShooterID:      shooterID,
		ExplosionPower: 1.0,
	}
	fb.Entity.NetworkID = BigFireballNetworkID
	fb.Entity.Width = 1.0
	fb.Entity.Height = 1.0
	fb.Entity.Gravity = 0.0 // 火球无重力
	fb.Entity.MaxHealth = 1
	fb.Entity.Health = 1
	return fb
}

// ============================================================
//                SmallFireball (ID 94) — 烈焰人火球
// ============================================================

const SmallFireballNetworkID = 94

// SmallFireball 烈焰人火球（小型，点燃目标）
type SmallFireball struct {
	Entity

	// ShooterID 发射者实体 ID
	ShooterID int64
}

// NewSmallFireball 创建烈焰人火球
func NewSmallFireball(shooterID int64) *SmallFireball {
	fb := &SmallFireball{
		ShooterID: shooterID,
	}
	fb.Entity.NetworkID = SmallFireballNetworkID
	fb.Entity.Width = 0.3125
	fb.Entity.Height = 0.3125
	fb.Entity.Gravity = 0.0
	fb.Entity.MaxHealth = 1
	fb.Entity.Health = 1
	return fb
}

// ============================================================
//               LavaSlime / MagmaCube (ID 42)
// ============================================================

const LavaSlimeNetworkID = 42

// LavaSlime 岩浆怪（地狱版史莱姆，免疫火焰）
type LavaSlime struct {
	*Monster

	// Size 大小 (1-4)
	Size int
}

// NewLavaSlime 创建岩浆怪
// - 与史莱姆类似: 可分裂
// - 免疫火焰/岩浆
// - 掉落: 岩浆膏 (最小才掉)
func NewLavaSlime() *LavaSlime {
	size := 1 + rand.Intn(4) // 1-4
	m := NewMonster(LavaSlimeNetworkID, "Magma Cube", lavaSlimeHealthForSize(size),
		0.6, 0.6, lavaSlimeDamageForSize(size))
	m.DropExpMin = 1
	m.DropExpMax = 4

	return &LavaSlime{
		Monster: m,
		Size:    size,
	}
}

// NewLavaSlimeWithSize 创建指定大小的岩浆怪
func NewLavaSlimeWithSize(size int) *LavaSlime {
	if size < 1 {
		size = 1
	}
	if size > 4 {
		size = 4
	}
	m := NewMonster(LavaSlimeNetworkID, "Magma Cube", lavaSlimeHealthForSize(size),
		0.6, 0.6, lavaSlimeDamageForSize(size))
	m.DropExpMin = 1
	m.DropExpMax = 4

	return &LavaSlime{
		Monster: m,
		Size:    size,
	}
}

func lavaSlimeHealthForSize(size int) int {
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

func lavaSlimeDamageForSize(size int) int {
	switch size {
	case 1:
		return 3
	case 2:
		return 4
	case 3:
		return 5
	case 4:
		return 6
	default:
		return 3
	}
}

// GetSize 获取大小
func (l *LavaSlime) GetSize() int {
	return l.Size
}

// SetSize 设置大小
func (l *LavaSlime) SetSize(size int) {
	l.Size = size
}

// ShouldSplit 死亡时是否应分裂
func (l *LavaSlime) ShouldSplit() bool {
	return l.Size > 1
}

// GetSplitSize 分裂后的子代大小
func (l *LavaSlime) GetSplitSize() int {
	return l.Size - 1
}

// IsFireImmune 岩浆怪免疫火焰
func (l *LavaSlime) IsFireImmune() bool {
	return true
}

// LavaSlimeDrops 岩浆怪掉落物
// 只有最小的 (size=1) 才掉岩浆膏
func LavaSlimeDrops(size int) []ZombieDropItem {
	const MagmaCream = 378
	if size == 1 {
		return []ZombieDropItem{{ItemID: MagmaCream, Count: 1}}
	}
	return nil
}
