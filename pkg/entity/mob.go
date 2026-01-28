package entity

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/entity/ai"
)

type Mob struct {
	*Living

	BehaviorManager *ai.BehaviorManager
	Jumping         bool
}

func NewMob() *Mob {
	m := &Mob{
		Living:          NewLiving(),
		BehaviorManager: ai.NewBehaviorManager(),
		Jumping:         false,
	}
	return m
}

func (m *Mob) AddBehavior(b ai.Behavior) {
	m.BehaviorManager.AddBehavior(b)
}

func (m *Mob) SetAIEnabled(enabled bool) {
	m.BehaviorManager.SetEnabled(enabled)
}

func (m *Mob) IsAIEnabled() bool {
	return m.BehaviorManager.IsEnabled()
}

func (m *Mob) GetCurrentBehavior() ai.Behavior {
	return m.BehaviorManager.GetCurrentBehavior()
}

func (m *Mob) Tick(currentTick int64) bool {
	if !m.Entity.Tick(currentTick) {
		return false
	}

	if m.AttackTime > 0 {
		m.AttackTime--
	}

	m.BehaviorManager.Tick()

	if m.Jumping && m.OnGround {
		m.Jump()
	} else if m.Jumping && !m.OnGround {
		m.Jumping = false
	}

	return true
}

func (m *Mob) GetPosition() (x, y, z float64) {
	return m.Position.X, m.Position.Y, m.Position.Z
}

func (m *Mob) SetPosition(x, y, z float64) {
	m.Position.X = x
	m.Position.Y = y
	m.Position.Z = z
	m.recalculateBoundingBox()
}

func (m *Mob) GetMotion() (x, y, z float64) {
	return m.Motion.X, m.Motion.Y, m.Motion.Z
}

func (m *Mob) SetMotion(x, y, z float64) {
	m.Motion.X = x
	m.Motion.Y = y
	m.Motion.Z = z
}

func (m *Mob) GetYaw() float64 {
	return m.Yaw
}

func (m *Mob) SetYaw(yaw float64) {
	m.Yaw = yaw
}

func (m *Mob) GetPitch() float64 {
	return m.Pitch
}

func (m *Mob) SetPitch(pitch float64) {
	m.Pitch = pitch
}

func (m *Mob) GetHeight() float64 {
	return m.Height
}

func (m *Mob) IsInsideOfWater() bool {

	return false
}

func (m *Mob) GetLevel() ai.LevelAccess {

	return nil
}

func (m *Mob) Move(dx, dy, dz float64) {
	m.Position.X += dx
	m.Position.Y += dy
	m.Position.Z += dz
	m.recalculateBoundingBox()
}

func (m *Mob) GetDirectionVector() (x, y, z float64) {
	pitchRad := m.Pitch * (math.Pi / 180)
	yawRad := m.Yaw * (math.Pi / 180)

	y = -math.Sin(pitchRad)
	xz := math.Cos(pitchRad)
	x = -xz * math.Sin(yawRad)
	z = xz * math.Cos(yawRad)

	return x, y, z
}

func (m *Mob) IsOnGround() bool {
	return m.OnGround
}
