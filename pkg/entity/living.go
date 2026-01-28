package entity

import (
	"math"
)

type Living struct {
	*Entity

	AttackTime   int
	LastDamage   float64
	DeathTime    int
	JumpVelocity float64

	HeadYaw float32
}

func NewLiving() *Living {
	l := &Living{
		Entity:       NewEntity(),
		AttackTime:   0,
		LastDamage:   0,
		DeathTime:    0,
		JumpVelocity: 0.42,
		HeadYaw:      0,
	}
	l.initLivingAttributes()
	return l
}

func (l *Living) initLivingAttributes() {
	l.Attributes.AddAttribute(GetDefaultAttribute(AttributeHealth))
	l.Attributes.AddAttribute(GetDefaultAttribute(AttributeMovementSpeed))
	l.Attributes.AddAttribute(GetDefaultAttribute(AttributeKnockbackResist))
	l.Attributes.AddAttribute(GetDefaultAttribute(AttributeFollowRange))
	l.Attributes.AddAttribute(GetDefaultAttribute(AttributeAttackDamage))
}

func (l *Living) GetName() string {
	return ""
}

func (l *Living) Attack(damage float64, cause int) bool {
	if l.Invulnerable || l.NoDamageTicks > 0 {
		return false
	}

	if attr := l.Attributes.GetAttribute(AttributeKnockbackResist); attr != nil {
		resistance := attr.GetValue()
		if resistance >= 1 {

		}
	}

	finalDamage := int(math.Round(damage))
	if finalDamage < 1 {
		finalDamage = 1
	}

	l.SetHealth(l.GetHealth() - finalDamage)
	l.LastDamage = damage
	l.NoDamageTicks = 10

	return true
}

func (l *Living) Heal(amount float64) {
	newHealth := l.GetHealth() + int(math.Round(amount))
	if newHealth > l.GetMaxHealth() {
		newHealth = l.GetMaxHealth()
	}
	l.SetHealth(newHealth)
}

func (l *Living) KnockBack(attackerX, attackerZ float64, base, verticalLimit float64) {

	dx := l.Position.X - attackerX
	dz := l.Position.Z - attackerZ
	dist := math.Sqrt(dx*dx + dz*dz)

	if dist <= 0 {
		return
	}

	f := 1 / dist

	l.Motion.X /= 2
	l.Motion.Y /= 2
	l.Motion.Z /= 2

	l.Motion.X += dx * f * base
	l.Motion.Y += base
	l.Motion.Z += dz * f * base

	if attr := l.Attributes.GetAttribute(AttributeKnockbackResist); attr != nil {
		resistance := attr.GetValue()
		if resistance > 0 {
			l.Motion.X *= (1 - resistance)
			l.Motion.Z *= (1 - resistance)

		}
	}

	if l.Motion.Y > verticalLimit {
		l.Motion.Y = verticalLimit
	}
}

func (l *Living) Kill() {
	l.SetHealth(0)
	l.DeathTime = 0
}

func (l *Living) GetJumpVelocity() float64 {
	return l.JumpVelocity
}

func (l *Living) Jump() {
	l.Motion.Y = l.GetJumpVelocity()
}

func (l *Living) Fall(fallDistance float64) {
	if fallDistance > 3 {
		damage := fallDistance - 3
		l.Attack(damage, DamageCauseFall)
	}
}

func (l *Living) UpdateMovementSpeed(speed float64) {
	if attr := l.Attributes.GetAttribute(AttributeMovementSpeed); attr != nil {
		attr.SetValue(speed)
	}
}

func (l *Living) GetMovementSpeed() float64 {
	if attr := l.Attributes.GetAttribute(AttributeMovementSpeed); attr != nil {
		return attr.GetValue()
	}
	return 0.1
}

const (
	DamageCauseContact         = 0
	DamageCauseEntityAttack    = 1
	DamageCauseProjectile      = 2
	DamageCauseSuffocation     = 3
	DamageCauseFall            = 4
	DamageCauseFire            = 5
	DamageCauseFireTick        = 6
	DamageCauseLava            = 7
	DamageCauseDrowning        = 8
	DamageCauseBlockExplosion  = 9
	DamageCauseEntityExplosion = 10
	DamageCauseVoid            = 11
	DamageCauseSuicide         = 12
	DamageCauseMagic           = 13
	DamageCauseStarvation      = 14
	DamageCauseCustom          = 15
)

func (l *Living) IsAlive() bool {
	return l.Health > 0
}

func (l *Living) EntityBaseTick(tickDiff int) bool {
	hasUpdate := false

	if l.IsAlive() {

		if l.AttackTime > 0 {
			l.AttackTime -= tickDiff
			if l.AttackTime < 0 {
				l.AttackTime = 0
			}
		}

		if l.NoDamageTicks > 0 {
			l.NoDamageTicks -= tickDiff
			if l.NoDamageTicks < 0 {
				l.NoDamageTicks = 0
			}
		}

		hasUpdate = true
	}

	return hasUpdate
}

func (l *Living) OnDeathUpdate(tickDiff int) bool {
	l.DeathTime += tickDiff

	return l.DeathTime >= 10
}

func (l *Living) GetDrops() []interface{} {
	return []interface{}{}
}

func (l *Living) HasLineOfSight(target IEntity) bool {

	return true
}

func (l *Living) GetDirectionVector() *Vector3 {
	pitchRad := l.Pitch * (math.Pi / 180)
	yawRad := l.Yaw * (math.Pi / 180)

	y := -math.Sin(pitchRad)
	xz := math.Cos(pitchRad)
	x := -xz * math.Sin(yawRad)
	z := xz * math.Cos(yawRad)

	return NewVector3(x, y, z)
}

func (l *Living) Distance(x, y, z float64) float64 {
	dx := l.Position.X - x
	dy := l.Position.Y - y
	dz := l.Position.Z - z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (l *Living) DistanceSquared(x, y, z float64) float64 {
	dx := l.Position.X - x
	dy := l.Position.Y - y
	dz := l.Position.Z - z
	return dx*dx + dy*dy + dz*dz
}

func (l *Living) ApplyEntityCollision(other IEntity) bool {
	otherPos := other.GetPosition()
	dx := otherPos.X - l.Position.X
	dz := otherPos.Z - l.Position.Z
	absMax := math.Max(math.Abs(dx), math.Abs(dz))

	if absMax >= 0.001 {
		dist := math.Sqrt(absMax)
		dx /= dist
		dz /= dist

		d3 := 1.0 / dist
		if d3 > 1.0 {
			d3 = 1.0
		}

		dx *= d3
		dz *= d3
		dx *= 0.05
		dz *= 0.05

		l.Motion.X -= dx
		l.Motion.Z -= dz
		return true
	}
	return false
}

func (l *Living) SetOnFire(seconds int) {
	ticks := seconds * 20
	if ticks > l.FireTicks {
		l.FireTicks = ticks
	}
}

func (l *Living) ExtinguishFire() {
	l.FireTicks = 0
}

func (l *Living) IsOnFire() bool {
	return l.FireTicks > 0
}
