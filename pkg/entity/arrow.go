package entity

import (
	"math"
)

const ArrowNetworkID = 80

type Arrow struct {
	*Projectile

	PunchKnockback float64
}

func NewArrow(shooterID int64, critical bool) *Arrow {
	a := &Arrow{
		Projectile:     NewProjectile(shooterID),
		PunchKnockback: 0,
	}

	a.Entity.NetworkID = ArrowNetworkID
	a.Entity.Width = 0.25
	a.Entity.Height = 0.25
	a.Entity.Gravity = 0.05
	a.Entity.Drag = 0.01
	a.Projectile.BaseDamage = 2.0
	a.Projectile.IsCritical = critical
	a.Projectile.MaxAge = 1200

	return a
}

type ArrowTickResult struct {
	ProjectileTickResult
	ShowCriticalParticle bool
	ParticleX            float64
	ParticleY            float64
	ParticleZ            float64
}

func (a *Arrow) TickArrow(nearbyEntities []IEntity, isCollided bool) ArrowTickResult {
	base := a.Projectile.TickProjectile(nearbyEntities, isCollided)

	result := ArrowTickResult{
		ProjectileTickResult: base,
	}

	if !a.Projectile.HadCollision && a.Projectile.IsCritical {
		result.ShowCriticalParticle = true
		result.ParticleX = a.Entity.Position.X
		result.ParticleY = a.Entity.Position.Y + a.Entity.Height/2
		result.ParticleZ = a.Entity.Position.Z
	} else if a.Entity.OnGround {

		a.Projectile.IsCritical = false
	}

	return result
}

type ArrowHitResult struct {
	Damage       float64
	KnockbackX   float64
	KnockbackY   float64
	KnockbackZ   float64
	HasKnockback bool
	TransferFire bool
	FireDuration int
}

func (a *Arrow) CalcArrowHit() ArrowHitResult {
	result := ArrowHitResult{
		Damage: a.Projectile.CalcHitDamage(),
	}

	if a.PunchKnockback > 0 {
		mx := a.Entity.Motion.X
		mz := a.Entity.Motion.Z
		horizontalSpeed := math.Sqrt(mx*mx + mz*mz)
		if horizontalSpeed > 0 {
			multiplier := a.PunchKnockback * 0.2 / horizontalSpeed
			result.HasKnockback = true
			result.KnockbackX = mx * multiplier
			result.KnockbackY = 0.1
			result.KnockbackZ = mz * multiplier
		}
	}

	if a.Projectile.ShouldTransferFire() {
		result.TransferFire = true
		result.FireDuration = a.Projectile.GetFireTransferDuration()
	}

	return result
}

func (a *Arrow) SetPunchKnockback(value float64) {
	a.PunchKnockback = value
}

func (a *Arrow) GetPunchKnockback() float64 {
	return a.PunchKnockback
}
