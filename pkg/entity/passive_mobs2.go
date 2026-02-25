package entity

// passive_mobs2.go — 额外被动/中立生物: IronGolem, SnowGolem, Squid, NPC
// 对应 PHP: entity/IronGolem.php, entity/SnowGolem.php, entity/Squid.php

import "math/rand"

// ============================================================
//                   IronGolem (ID 20)
// ============================================================

const IronGolemNetworkID = 20

// IronGolem 铁傀儡（中立生物，玩家建造或村庄生成）
type IronGolem struct {
	*Monster

	// PlayerCreated 是否由玩家建造
	PlayerCreated bool
}

// NewIronGolem 创建铁傀儡
// - MaxHealth: 100, 攻击伤害: 21
// - 尺寸: 1.4×2.7
// - 不受击退影响
// - 掉落: 铁锭 3-5 + 罂粟 0-2
func NewIronGolem() *IronGolem {
	m := NewMonster(IronGolemNetworkID, "Iron Golem", 100, 1.4, 2.7, 21)
	m.DropExpMin = 0
	m.DropExpMax = 0 // 铁傀儡不掉经验

	return &IronGolem{
		Monster: m,
	}
}

// IsPlayerCreated 是否由玩家建造
func (g *IronGolem) IsPlayerCreated() bool {
	return g.PlayerCreated
}

// IronGolemDrops 铁傀儡掉落物
func IronGolemDrops() []ZombieDropItem {
	const (
		IronIngot = 265
		Poppy     = 38 // RED_FLOWER
	)
	drops := []ZombieDropItem{
		{ItemID: IronIngot, Count: 3 + rand.Intn(3)}, // 3-5
	}
	// 0-2 罂粟
	poppyCount := rand.Intn(3)
	if poppyCount > 0 {
		drops = append(drops, ZombieDropItem{ItemID: Poppy, Count: poppyCount})
	}
	return drops
}

// ============================================================
//                   SnowGolem (ID 21)
// ============================================================

const SnowGolemNetworkID = 21

// SnowGolem 雪傀儡（玩家建造，投雪球攻击）
type SnowGolem struct {
	*Monster

	// Sheared 是否被剪毛（去掉南瓜头）
	Sheared bool
}

// NewSnowGolem 创建雪傀儡
// - MaxHealth: 4, 无近战伤害
// - 经过的路径会放雪
// - 在热带/下界受伤
// - 投掷雪球攻击敌对怪物
func NewSnowGolem() *SnowGolem {
	m := NewMonster(SnowGolemNetworkID, "Snow Golem", 4, 0.7, 1.9, 0)
	m.DropExpMin = 0
	m.DropExpMax = 0

	return &SnowGolem{
		Monster: m,
	}
}

// SetSheared 去掉南瓜头
func (s *SnowGolem) SetSheared(sheared bool) {
	s.Sheared = sheared
}

// IsSheared 是否已去掉南瓜头
func (s *SnowGolem) IsSheared() bool {
	return s.Sheared
}

// SnowGolemDrops 雪傀儡掉落物 (0-15 个雪球)
func SnowGolemDrops() []ZombieDropItem {
	const Snowball = 332
	count := rand.Intn(16) // 0-15
	if count == 0 {
		return nil
	}
	return []ZombieDropItem{{ItemID: Snowball, Count: count}}
}

// ============================================================
//                     Squid (ID 17)
// ============================================================

const SquidNetworkID = 17

// NewSquid 创建鱿鱼
// - MaxHealth: 10, 无攻击
// - 水生生物，出水窒息
// - 掉落: 墨囊 1-3
func NewSquid() *Monster {
	m := NewMonster(SquidNetworkID, "Squid", 10, 0.8, 0.8, 0)
	m.DropExpMin = 1
	m.DropExpMax = 3
	return m
}

// SquidDrops 鱿鱼掉落物
func SquidDrops() []ZombieDropItem {
	const InkSac = 351                                                          // DYE meta 0
	return []ZombieDropItem{{ItemID: InkSac, Meta: 0, Count: 1 + rand.Intn(3)}} // 1-3
}

// ============================================================
//                      NPC (ID 15)
// ============================================================

const NPCNetworkID = 15

// NPC 村民 (可交易的中立实体)
type NPC struct {
	*Monster

	// Profession 职业 (0=农民 1=图书管理员 2=牧师 3=铁匠 4=屠户)
	Profession int
}

const (
	ProfessionFarmer     = 0
	ProfessionLibrarian  = 1
	ProfessionPriest     = 2
	ProfessionBlacksmith = 3
	ProfessionButcher    = 4
)

// NewNPC 创建村民
// - MaxHealth: 20, 无攻击
// - 随机职业
func NewNPC() *NPC {
	m := NewMonster(NPCNetworkID, "Villager", 20, 0.6, 1.8, 0)
	m.DropExpMin = 0
	m.DropExpMax = 0

	return &NPC{
		Monster:    m,
		Profession: rand.Intn(5),
	}
}

// NewNPCWithProfession 创建指定职业的村民
func NewNPCWithProfession(profession int) *NPC {
	m := NewMonster(NPCNetworkID, "Villager", 20, 0.6, 1.8, 0)
	m.DropExpMin = 0
	m.DropExpMax = 0

	return &NPC{
		Monster:    m,
		Profession: profession,
	}
}

// GetProfession 获取职业
func (n *NPC) GetProfession() int {
	return n.Profession
}

// SetProfession 设置职业
func (n *NPC) SetProfession(profession int) {
	n.Profession = profession
}

// NPCDrops 村民无掉落物
func NPCDrops() []ZombieDropItem {
	return nil
}

// ============================================================
//                  FlyingAnimal 飞行动物基类
// ============================================================

// FlyingAnimal 飞行动物（如鹦鹉等）
// 基础结构，可被具体飞行生物扩展
type FlyingAnimal struct {
	*Animal

	// Flying 是否在飞行
	Flying bool
}

// NewFlyingAnimal 创建飞行动物基类
func NewFlyingAnimal(networkID int, name string, maxHealth int, width, height float64) *FlyingAnimal {
	a := NewAnimal(networkID, name, maxHealth, width, height)
	a.Entity.Gravity = 0.02 // 飞行生物重力更小

	return &FlyingAnimal{
		Animal: a,
	}
}

// SetFlying 设置飞行状态
func (f *FlyingAnimal) SetFlying(flying bool) {
	f.Flying = flying
}

// IsFlying 是否在飞行
func (f *FlyingAnimal) IsFlying() bool {
	return f.Flying
}
