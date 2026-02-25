package entity

import "math"

type MoveHelper struct {
	Entity   *Entity
	PosX     float64
	PosY     float64
	PosZ     float64
	Speed    float64
	IsMoving bool
}

func NewMoveHelper(e *Entity) *MoveHelper {
	return &MoveHelper{Entity: e}
}
func (m *MoveHelper) SetMoveTo(x, y, z, speed float64) {
	m.PosX = x
	m.PosY = y
	m.PosZ = z
	m.Speed = speed
	m.IsMoving = true
}
func (m *MoveHelper) OnUpdateMoveHelper() {
	if !m.IsMoving {
		return
	}

	e := m.Entity
	dx := m.PosX - e.Position.X
	dz := m.PosZ - e.Position.Z
	distSqXZ := dx*dx + dz*dz

	if distSqXZ < 0.25 {
		m.IsMoving = false
		e.Motion.X = 0
		e.Motion.Z = 0
		return
	}

	dist := math.Sqrt(distSqXZ)
	yaw := math.Atan2(dz, dx)*180.0/math.Pi - 90.0
	e.Yaw = yaw
	speed := m.Speed * e.MovementSpeed * 0.45
	e.Motion.X = dx / dist * speed
	e.Motion.Z = dz / dist * speed

	dy := m.PosY - e.Position.Y
	if dy > e.StepHeight && distSqXZ < 1.0 && e.OnGround {
		e.Motion.Y = 0.42
	}
}

type LookHelper struct {
	Entity       *Entity
	DeltaX       float64
	DeltaZ       float64
	IsLooking    bool
	LookPosX     float64
	LookPosY     float64
	LookPosZ     float64
	MaxYawTurn   float64
	MaxPitchTurn float64
}

func NewLookHelper(e *Entity) *LookHelper {
	return &LookHelper{
		Entity:       e,
		MaxYawTurn:   10.0,
		MaxPitchTurn: 40.0,
	}
}
func (l *LookHelper) SetLookPosition(x, y, z, maxYaw, maxPitch float64) {
	l.LookPosX = x
	l.LookPosY = y
	l.LookPosZ = z
	l.MaxYawTurn = maxYaw
	l.MaxPitchTurn = maxPitch
	l.IsLooking = true
}
func (l *LookHelper) OnUpdateLook() {
	if !l.IsLooking {
		return
	}

	l.IsLooking = false

	e := l.Entity
	dx := l.LookPosX - e.Position.X
	dy := l.LookPosY - (e.Position.Y + e.EyeHeight)
	dz := l.LookPosZ - e.Position.Z
	distHoriz := math.Sqrt(dx*dx + dz*dz)

	targetYaw := math.Atan2(dz, dx)*180.0/math.Pi - 90.0
	targetPitch := -math.Atan2(dy, distHoriz) * 180.0 / math.Pi

	e.Pitch = l.updateRotation(e.Pitch, targetPitch, l.MaxPitchTurn)
	e.Yaw = l.updateRotation(e.Yaw, targetYaw, l.MaxYawTurn)
}
func (l *LookHelper) updateRotation(current, target, maxDelta float64) float64 {
	diff := wrapDegrees(target - current)
	if diff > maxDelta {
		diff = maxDelta
	}
	if diff < -maxDelta {
		diff = -maxDelta
	}
	return current + diff
}
func wrapDegrees(deg float64) float64 {
	deg = math.Mod(deg, 360.0)
	if deg >= 180.0 {
		deg -= 360.0
	}
	if deg < -180.0 {
		deg += 360.0
	}
	return deg
}

type JumpHelper struct {
	Entity    *Entity
	IsJumping bool
}

func NewJumpHelper(e *Entity) *JumpHelper {
	return &JumpHelper{Entity: e}
}

func (j *JumpHelper) SetJumping() {
	j.IsJumping = true
}
func (j *JumpHelper) DoJump() {
	if j.IsJumping {
		if j.Entity.OnGround {
			j.Entity.Motion.Y = 0.42
		}
		j.IsJumping = false
	}
}
