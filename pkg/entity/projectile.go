package entity

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)
type Projectile struct {
	*Entity
	ShootingEntityID int64
	BaseDamage float64
	HadCollision bool
	IsCritical bool
	DragBeforeGravity bool
	ProjectileAge int
	MaxAge int
}
func NewProjectile(shooterID int64) *Projectile {
	p := &Projectile{
		Entity:            NewEntity(),
		ShootingEntityID:  shooterID,
		BaseDamage:        0,
		HadCollision:      false,
		IsCritical:        false,
		DragBeforeGravity: true,
		ProjectileAge:     0,
		MaxAge:            1200,
	}
	p.Entity.Health = 1
	p.Entity.MaxHealth = 1
	p.Entity.CanCollide = false
	p.Entity.Width = 0.25
	p.Entity.Height = 0.25
	p.Entity.Gravity = 0.03
	p.Entity.Drag = 0.01

	return p
}
type ProjectileTickResult struct {
	HasUpdate   bool
	HitEntityID int64
	HitBlock    bool
	ShouldClose bool
	Damage      float64
}
func (p *Projectile) TickProjectile(nearbyEntities []IEntity, isCollided bool) ProjectileTickResult {
	result := ProjectileTickResult{}

	p.Entity.TicksLived++
	p.ProjectileAge++
	if p.ProjectileAge >= p.MaxAge {
		result.ShouldClose = true
		return result
	}
	if isCollided && !p.HadCollision {
		p.HadCollision = true
		p.Entity.Motion.X = 0
		p.Entity.Motion.Y = 0
		p.Entity.Motion.Z = 0
		result.HitBlock = true
		result.HasUpdate = true
	} else if !isCollided && p.HadCollision {
		p.HadCollision = false
	}
	if !p.HadCollision {
		nearDist := math.MaxFloat64
		var nearEntityID int64

		moveX := p.Entity.Position.X + p.Entity.Motion.X
		moveY := p.Entity.Position.Y + p.Entity.Motion.Y
		moveZ := p.Entity.Position.Z + p.Entity.Motion.Z

		for _, ent := range nearbyEntities {
			entID := ent.GetID()
			if entID == p.ShootingEntityID && p.Entity.TicksLived < 5 {
				continue
			}
			entPos := ent.GetPosition()
			dx := entPos.X - moveX
			dy := entPos.Y - moveY
			dz := entPos.Z - moveZ
			distSq := dx*dx + dy*dy + dz*dz
			bb := ent.GetBoundingBox()
			if bb == nil {
				continue
			}
			hitRadius := ((bb.MaxX - bb.MinX) / 2) + 0.3

			if distSq < hitRadius*hitRadius && distSq < nearDist {
				nearDist = distSq
				nearEntityID = entID
			}
		}

		if nearEntityID != 0 {
			result.HitEntityID = nearEntityID
			result.Damage = p.CalcHitDamage()
			result.ShouldClose = true
			result.HasUpdate = true
			p.HadCollision = true
			return result
		}
	}
	mx := p.Entity.Motion.X
	my := p.Entity.Motion.Y
	mz := p.Entity.Motion.Z
	if !p.Entity.OnGround || math.Abs(mx) > 0.00001 || math.Abs(my) > 0.00001 || math.Abs(mz) > 0.00001 {
		f := math.Sqrt(mx*mx + mz*mz)
		p.Entity.Yaw = math.Atan2(mx, mz) * 180 / math.Pi
		p.Entity.Pitch = math.Atan2(my, f) * 180 / math.Pi
		result.HasUpdate = true
	}
	if p.DragBeforeGravity {
		p.Entity.Motion.X *= 1 - p.Entity.Drag
		p.Entity.Motion.Y *= 1 - p.Entity.Drag
		p.Entity.Motion.Z *= 1 - p.Entity.Drag
		p.Entity.Motion.Y -= p.Entity.Gravity
	} else {
		p.Entity.Motion.Y -= p.Entity.Gravity
		p.Entity.Motion.X *= 1 - p.Entity.Drag
		p.Entity.Motion.Y *= 1 - p.Entity.Drag
		p.Entity.Motion.Z *= 1 - p.Entity.Drag
	}

	return result
}
func (p *Projectile) CalcHitDamage() float64 {
	mx := p.Entity.Motion.X
	my := p.Entity.Motion.Y
	mz := p.Entity.Motion.Z
	speed := math.Sqrt(mx*mx + my*my + mz*mz)
	damage := math.Ceil(speed * p.BaseDamage)

	if p.IsCritical && damage > 0 {
		bonus := int(damage/2) + 1
		if bonus > 0 {
			damage += float64(bonus / 2)
		}
	}

	return damage
}
func (p *Projectile) SetShooter(entityID int64) {
	p.ShootingEntityID = entityID
}
func (p *Projectile) IsOnGround() bool {
	return p.HadCollision
}
func (p *Projectile) GetSpeed() float64 {
	mx := p.Entity.Motion.X
	my := p.Entity.Motion.Y
	mz := p.Entity.Motion.Z
	return math.Sqrt(mx*mx + my*my + mz*mz)
}
func (p *Projectile) ShouldTransferFire() bool {
	return p.Entity.FireTicks > 0
}
func (p *Projectile) GetFireTransferDuration() int {
	return 100
}
func (p *Projectile) SaveProjectileNBT() {
	p.Entity.SaveNBT()
	p.Entity.NamedTag.Set(nbt.NewShortTag("Age", int16(p.ProjectileAge)))
}
