package ai

import (
	"math"
	"math/rand"
)

type Behavior interface {
	Name() string

	ShouldStart() bool

	CanContinue() bool

	OnTick()

	OnEnd()
}

type BehaviorBase struct {
	Mob MobEntity
}

type MobEntity interface {
	GetPosition() (x, y, z float64)
	SetPosition(x, y, z float64)
	GetMotion() (x, y, z float64)
	SetMotion(x, y, z float64)
	GetYaw() float64
	SetYaw(yaw float64)
	GetPitch() float64
	SetPitch(pitch float64)
	GetHeight() float64
	IsInsideOfWater() bool
	GetLevel() LevelAccess
	Move(dx, dy, dz float64)
	GetDirectionVector() (x, y, z float64)
	IsOnGround() bool
}

type LevelAccess interface {
	GetBlock(x, y, z int) BlockInfo
	GetNearestPlayer(x, y, z float64, maxDistance float64) PlayerEntity
	GetEntities() []MobEntity
}

type BlockInfo interface {
	IsSolid() bool
	IsAir() bool
}

type PlayerEntity interface {
	GetPosition() (x, y, z float64)
	IsAlive() bool
	IsConnected() bool
	IsSurvival() bool
}

func NewBehaviorBase(mob MobEntity) *BehaviorBase {
	return &BehaviorBase{Mob: mob}
}

func (b *BehaviorBase) CheckSwimming() {
	if b.Mob.IsInsideOfWater() {
		_, my, _ := b.Mob.GetMotion()
		if my <= 0 {
			b.Mob.SetMotion(0, 0.2, 0)
		}
	}
}

func (b *BehaviorBase) AimAt(targetX, targetY, targetZ float64) {
	ex, ey, ez := b.Mob.GetPosition()
	dx := targetX - ex
	dy := targetY - ey
	dz := targetZ - ez

	dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
	if dist == 0 {
		return
	}

	pitch := -math.Asin(dy/dist) * 180 / math.Pi

	yaw := -math.Atan2(dx, dz) * 180 / math.Pi

	b.Mob.SetPitch(pitch)
	b.Mob.SetYaw(yaw)
}

func (b *BehaviorBase) Distance(x, y, z float64) float64 {
	ex, ey, ez := b.Mob.GetPosition()
	dx := x - ex
	dy := y - ey
	dz := z - ez
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

type BehaviorManager struct {
	behaviors       []Behavior
	currentBehavior Behavior
	enabled         bool
}

func NewBehaviorManager() *BehaviorManager {
	return &BehaviorManager{
		behaviors: make([]Behavior, 0),
		enabled:   true,
	}
}

func (m *BehaviorManager) AddBehavior(b Behavior) {
	m.behaviors = append(m.behaviors, b)
}

func (m *BehaviorManager) SetEnabled(enabled bool) {
	m.enabled = enabled
}

func (m *BehaviorManager) IsEnabled() bool {
	return m.enabled
}

func (m *BehaviorManager) GetCurrentBehavior() Behavior {
	return m.currentBehavior
}

func (m *BehaviorManager) Tick() {
	if !m.enabled {
		return
	}

	m.currentBehavior = m.checkBehavior()
	if m.currentBehavior != nil {
		m.currentBehavior.OnTick()
	}
}

func (m *BehaviorManager) checkBehavior() Behavior {
	for i, behavior := range m.behaviors {

		if behavior == m.currentBehavior {
			if behavior.CanContinue() {
				return behavior
			}
			behavior.OnEnd()
			m.currentBehavior = nil
		}

		if behavior.ShouldStart() {

			if m.currentBehavior == nil || m.getBehaviorIndex(m.currentBehavior) > i {
				if m.currentBehavior != nil {
					m.currentBehavior.OnEnd()
				}
				return behavior
			}
		}
	}
	return nil
}

func (m *BehaviorManager) getBehaviorIndex(b Behavior) int {
	for i, behavior := range m.behaviors {
		if behavior == b {
			return i
		}
	}
	return -1
}

func RandomChance(n int) bool {
	return rand.Intn(n) == 0
}
