package entity

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
)

type Monster struct {
	*Mob

	AttackDamage int

	AttackingTick int

	TargetEntityID int64

	MobName string

	DropExpMin int
	DropExpMax int

	IsBabyFlag bool
}

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

func (m *Monster) GetHurt() int {
	return m.AttackDamage
}

func (m *Monster) SetHurt(damage int) {
	m.AttackDamage = damage
}

func (m *Monster) SetTarget(entityID int64) {
	m.TargetEntityID = entityID
}

func (m *Monster) GetTarget() int64 {
	return m.TargetEntityID
}

func (m *Monster) HasTarget() bool {
	return m.TargetEntityID != 0
}

func (m *Monster) ClearTarget() {
	m.TargetEntityID = 0
}

func (m *Monster) SetBaby(baby bool) {
	m.IsBabyFlag = baby
}

func (m *Monster) IsBaby() bool {
	return m.IsBabyFlag
}

type MonsterTickResult struct {
	HasUpdate bool
}

func (m *Monster) TickMonster() MonsterTickResult {
	result := MonsterTickResult{}

	if m.AttackingTick > 0 {
		m.AttackingTick--
		result.HasUpdate = true
	}

	return result
}

func (m *Monster) OnAttacked() {
	m.AttackingTick = 20
}

func (m *Monster) SaveMonsterNBT() {
	m.Mob.Living.Entity.SaveNBT()
	if m.IsBabyFlag {
		m.Mob.Living.Entity.NamedTag.Set(nbt.NewByteTag("IsBaby", 1))
	}
}

func (m *Monster) LoadMonsterFromNBT() {
	if m.Mob.Living.Entity.NamedTag == nil {
		return
	}
	if m.Mob.Living.Entity.NamedTag.GetByte("IsBaby") == 1 {
		m.SetBaby(true)
	}
}

func (m *Monster) GetName() string {
	return m.MobName
}

type HostileMobInfo struct {
	NetworkID    int
	Name         string
	MaxHealth    int
	Width        float64
	Height       float64
	AttackDamage int
	DropExpMin   int
	DropExpMax   int
	BurnInDay    bool
}

var hostileMobs = map[int]HostileMobInfo{

	32: {NetworkID: 32, Name: "Zombie", MaxHealth: 20, Width: 0.6, Height: 1.8, AttackDamage: 4, DropExpMin: 5, DropExpMax: 5, BurnInDay: true},
	34: {NetworkID: 34, Name: "Skeleton", MaxHealth: 20, Width: 0.6, Height: 1.8, AttackDamage: 4, DropExpMin: 5, DropExpMax: 5, BurnInDay: true},
	33: {NetworkID: 33, Name: "Creeper", MaxHealth: 20, Width: 0.6, Height: 1.8, AttackDamage: 0, DropExpMin: 5, DropExpMax: 5},
	35: {NetworkID: 35, Name: "Spider", MaxHealth: 16, Width: 1.4, Height: 0.9, AttackDamage: 3, DropExpMin: 5, DropExpMax: 5},

	38: {NetworkID: 38, Name: "Enderman", MaxHealth: 40, Width: 0.6, Height: 2.9, AttackDamage: 7, DropExpMin: 5, DropExpMax: 5},
	43: {NetworkID: 43, Name: "Blaze", MaxHealth: 20, Width: 0.6, Height: 1.8, AttackDamage: 6, DropExpMin: 10, DropExpMax: 10},
	41: {NetworkID: 41, Name: "Ghast", MaxHealth: 10, Width: 4.0, Height: 4.0, AttackDamage: 0, DropExpMin: 5, DropExpMax: 5},
	37: {NetworkID: 37, Name: "Slime", MaxHealth: 16, Width: 2.0, Height: 2.0, AttackDamage: 4, DropExpMin: 1, DropExpMax: 4},

	36: {NetworkID: 36, Name: "Cave Spider", MaxHealth: 12, Width: 0.7, Height: 0.5, AttackDamage: 2, DropExpMin: 5, DropExpMax: 5},
	40: {NetworkID: 40, Name: "Silverfish", MaxHealth: 8, Width: 0.4, Height: 0.3, AttackDamage: 1, DropExpMin: 5, DropExpMax: 5},
	42: {NetworkID: 42, Name: "Magma Cube", MaxHealth: 16, Width: 2.0, Height: 2.0, AttackDamage: 3, DropExpMin: 1, DropExpMax: 4},
	39: {NetworkID: 39, Name: "Zombie Pigman", MaxHealth: 20, Width: 0.6, Height: 1.8, AttackDamage: 5, DropExpMin: 5, DropExpMax: 5},
}

func GetHostileMobInfo(networkID int) *HostileMobInfo {
	if info, ok := hostileMobs[networkID]; ok {
		return &info
	}
	return nil
}

func IsHostileMob(networkID int) bool {
	_, ok := hostileMobs[networkID]
	return ok
}

func AllHostileMobIDs() []int {
	ids := make([]int, 0, len(hostileMobs))
	for id := range hostileMobs {
		ids = append(ids, id)
	}
	return ids
}

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
