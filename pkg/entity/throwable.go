package entity

const (
	SnowballNetworkID = 81
	EggNetworkID      = 82
	ChickenNetworkID  = 10

	EggChickenSpawnChance = 8
)
type Snowball struct {
	*Projectile
}
func NewSnowball(shooterID int64) *Snowball {
	s := &Snowball{
		Projectile: NewProjectile(shooterID),
	}

	s.Entity.NetworkID = SnowballNetworkID
	s.Entity.Width = 0.25
	s.Entity.Height = 0.25
	s.Entity.Gravity = 0.03
	s.Entity.Drag = 0.01
	s.Projectile.BaseDamage = 0
	s.Projectile.MaxAge = 1200

	return s
}
type ThrowableTickResult struct {
	ProjectileTickResult
	ShouldSpawnChicken bool
	SpawnX             float64
	SpawnY             float64
	SpawnZ             float64
}
func (s *Snowball) TickSnowball(nearbyEntities []IEntity, isCollided bool) ThrowableTickResult {
	base := s.Projectile.TickProjectile(nearbyEntities, isCollided)

	result := ThrowableTickResult{
		ProjectileTickResult: base,
	}
	if isCollided || s.Entity.OnGround {
		result.ShouldClose = true
	}

	return result
}
type Egg struct {
	*Projectile
}
func NewEgg(shooterID int64) *Egg {
	e := &Egg{
		Projectile: NewProjectile(shooterID),
	}

	e.Entity.NetworkID = EggNetworkID
	e.Entity.Width = 0.25
	e.Entity.Height = 0.25
	e.Entity.Gravity = 0.03
	e.Entity.Drag = 0.01
	e.Projectile.BaseDamage = 0
	e.Projectile.MaxAge = 1200

	return e
}
func (e *Egg) TickEgg(nearbyEntities []IEntity, isCollided bool, shouldSpawnChicken bool) ThrowableTickResult {
	base := e.Projectile.TickProjectile(nearbyEntities, isCollided)

	result := ThrowableTickResult{
		ProjectileTickResult: base,
	}
	if isCollided {
		result.ShouldClose = true
		if shouldSpawnChicken {
			result.ShouldSpawnChicken = true
			result.SpawnX = e.Entity.Position.X
			result.SpawnY = e.Entity.Position.Y
			result.SpawnZ = e.Entity.Position.Z
		}
	}

	return result
}
