package entity

// hostile.go — 敌对生物基类 + 敌对怪物属性表
// 对应 PHP: entity/Monster.php, entity/Creature.php
//
// PHP继承链: Living → Creature → Mob → Monster → (Zombie/Creeper/...)
// Go继承链: Entity → (Living) → Mob → Monster
//
// Monster 在 PHP 中非常简单，只多了 getHurt() 返回攻击伤害。
// Go 端 mob.go 已实现 Mob (AI System)，此文件在 Mob 之上添加:
//   - 攻击伤害值
//   - 攻击动画 tick
//   - 目标玩家跟踪
//   - 敌对怪物属性查询表

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
)

// ============ Monster 基类 ============

// Monster 敌对生物基类
type Monster struct {
	*Mob

	// AttackDamage 攻击伤害
	AttackDamage int

	// AttackingTick 攻击动画计时（受击后 20 tick 红色闪烁）
	AttackingTick int

	// TargetEntityID 当前追击目标实体ID（0=无目标）
	TargetEntityID int64

	// MobName 怪物名称
	MobName string

	// DropExpMin/Max 掉落经验范围
	DropExpMin int
	DropExpMax int

	// IsBabyFlag 是否为幼崽（僵尸等可以有幼崽变种）
	IsBabyFlag bool
}

// NewMonster 创建敌对生物
func NewMonster(networkID int, name string, maxHealth int, width, height float64, attackDamage int) *Monster {
	m := &Monster{
		Mob:          NewMob(),
		MobName:      name,
		AttackDamage: attackDamage,
	}

	m.Mob.Living.Entity.NetworkID = networkID
	m.Mob.Living.Entity.Width = width
	m.Mob.Living.Entity.Height = height
	m.Mob.Living.Entity.MaxHealth = maxHealth
	m.Mob.Living.Entity.Health = maxHealth
	m.Mob.Living.Entity.Gravity = 0.04
	m.Mob.Living.Entity.Drag = 0.02
	m.Mob.Living.Entity.StepHeight = 0.6

	return m
}

// ============ 攻击系统 ============

// GetHurt 获取攻击伤害值
// 对应 PHP Monster::getHurt()
func (m *Monster) GetHurt() int {
	return m.AttackDamage
}

// SetHurt 设置攻击伤害值
func (m *Monster) SetHurt(damage int) {
	m.AttackDamage = damage
}

// ============ 目标系统 ============

// SetTarget 设置追击目标
func (m *Monster) SetTarget(entityID int64) {
	m.TargetEntityID = entityID
}

// GetTarget 获取追击目标ID
func (m *Monster) GetTarget() int64 {
	return m.TargetEntityID
}

// HasTarget 是否有追击目标
func (m *Monster) HasTarget() bool {
	return m.TargetEntityID != 0
}

// ClearTarget 清除追击目标
func (m *Monster) ClearTarget() {
	m.TargetEntityID = 0
}

// ============ 幼崽系统（僵尸等） ============

// SetBaby 设置幼崽状态
func (m *Monster) SetBaby(baby bool) {
	m.IsBabyFlag = baby
}

// IsBaby 是否为幼崽
func (m *Monster) IsBaby() bool {
	return m.IsBabyFlag
}

// ============ Tick ============

// MonsterTickResult 敌对生物 tick 结果
type MonsterTickResult struct {
	HasUpdate bool
}

// TickMonster 敌对生物 tick 逻辑
func (m *Monster) TickMonster() MonsterTickResult {
	result := MonsterTickResult{}

	// 攻击动画计时递减
	if m.AttackingTick > 0 {
		m.AttackingTick--
		result.HasUpdate = true
	}

	return result
}

// OnAttacked 受到攻击时（设置攻击动画计时）
// 对应 PHP Creature::attack() 中的 attackingTick = 20
func (m *Monster) OnAttacked() {
	m.AttackingTick = 20
}

// ============ NBT ============

// SaveMonsterNBT 保存敌对生物 NBT
func (m *Monster) SaveMonsterNBT() {
	m.Mob.Living.Entity.SaveNBT()
	if m.IsBabyFlag {
		m.Mob.Living.Entity.NamedTag.Set(nbt.NewByteTag("IsBaby", 1))
	}
}

// LoadMonsterFromNBT 从 NBT 加载敌对生物数据
func (m *Monster) LoadMonsterFromNBT() {
	if m.Mob.Living.Entity.NamedTag == nil {
		return
	}
	if m.Mob.Living.Entity.NamedTag.GetByte("IsBaby") == 1 {
		m.SetBaby(true)
	}
}

// ============ 辅助 ============

// GetName 获取怪物名称
func (m *Monster) GetName() string {
	return m.MobName
}

// ============ 敌对怪物属性表 ============

// HostileMobInfo 敌对怪物的注册信息
type HostileMobInfo struct {
	NetworkID    int
	Name         string
	MaxHealth    int
	Width        float64
	Height       float64
	AttackDamage int // 近战攻击伤害
	DropExpMin   int
	DropExpMax   int
	BurnInDay    bool // 是否白天燃烧（僵尸/骷髅）
}

// hostileMobs 所有敌对怪物属性表
var hostileMobs = map[int]HostileMobInfo{
	// 常见近战怪物  (E.2)
	32: {NetworkID: 32, Name: "Zombie", MaxHealth: 20, Width: 0.6, Height: 1.8, AttackDamage: 4, DropExpMin: 5, DropExpMax: 5, BurnInDay: true},
	34: {NetworkID: 34, Name: "Skeleton", MaxHealth: 20, Width: 0.6, Height: 1.8, AttackDamage: 4, DropExpMin: 5, DropExpMax: 5, BurnInDay: true},
	33: {NetworkID: 33, Name: "Creeper", MaxHealth: 20, Width: 0.6, Height: 1.8, AttackDamage: 0, DropExpMin: 5, DropExpMax: 5},
	35: {NetworkID: 35, Name: "Spider", MaxHealth: 16, Width: 1.4, Height: 0.9, AttackDamage: 3, DropExpMin: 5, DropExpMax: 5},

	// 特殊怪物  (E.3)
	38: {NetworkID: 38, Name: "Enderman", MaxHealth: 40, Width: 0.6, Height: 2.9, AttackDamage: 7, DropExpMin: 5, DropExpMax: 5},
	43: {NetworkID: 43, Name: "Blaze", MaxHealth: 20, Width: 0.6, Height: 1.8, AttackDamage: 6, DropExpMin: 10, DropExpMax: 10},
	41: {NetworkID: 41, Name: "Ghast", MaxHealth: 10, Width: 4.0, Height: 4.0, AttackDamage: 0, DropExpMin: 5, DropExpMax: 5},
	37: {NetworkID: 37, Name: "Slime", MaxHealth: 16, Width: 2.0, Height: 2.0, AttackDamage: 4, DropExpMin: 1, DropExpMax: 4},

	// 变种
	36: {NetworkID: 36, Name: "Cave Spider", MaxHealth: 12, Width: 0.7, Height: 0.5, AttackDamage: 2, DropExpMin: 5, DropExpMax: 5},
	40: {NetworkID: 40, Name: "Silverfish", MaxHealth: 8, Width: 0.4, Height: 0.3, AttackDamage: 1, DropExpMin: 5, DropExpMax: 5},
	42: {NetworkID: 42, Name: "Magma Cube", MaxHealth: 16, Width: 2.0, Height: 2.0, AttackDamage: 3, DropExpMin: 1, DropExpMax: 4},
	39: {NetworkID: 39, Name: "Zombie Pigman", MaxHealth: 20, Width: 0.6, Height: 1.8, AttackDamage: 5, DropExpMin: 5, DropExpMax: 5},
}

// GetHostileMobInfo 获取敌对怪物信息
func GetHostileMobInfo(networkID int) *HostileMobInfo {
	if info, ok := hostileMobs[networkID]; ok {
		return &info
	}
	return nil
}

// IsHostileMob 判断网络ID是否为敌对怪物
func IsHostileMob(networkID int) bool {
	_, ok := hostileMobs[networkID]
	return ok
}

// AllHostileMobIDs 返回所有敌对怪物的网络ID列表
func AllHostileMobIDs() []int {
	ids := make([]int, 0, len(hostileMobs))
	for id := range hostileMobs {
		ids = append(ids, id)
	}
	return ids
}

// NewMonsterFromInfo 根据网络ID创建预配置的 Monster 实例
func NewMonsterFromInfo(networkID int) *Monster {
	info := GetHostileMobInfo(networkID)
	if info == nil {
		return nil
	}
	m := NewMonster(info.NetworkID, info.Name, info.MaxHealth, info.Width, info.Height, info.AttackDamage)
	m.DropExpMin = info.DropExpMin
	m.DropExpMax = info.DropExpMax
	return m
}
