package entity

// throwable.go — 投掷物实体（雪球/鸡蛋）
// 对应 PHP: entity/Snowball.php (~66行), entity/Egg.php (~93行)
//
// 两者继承 Projectile，碰撞/落地即消失。
// Egg 额外有 12.5% 概率在碰撞位置生成小鸡（Baby Chicken）。

// ============ 常量 ============

const (
	SnowballNetworkID = 81
	EggNetworkID      = 82
	ChickenNetworkID  = 10 // 小鸡的网络ID

	EggChickenSpawnChance = 8 // 1/8 = 12.5% 概率生成小鸡
)

// ============ Snowball 雪球 ============

// Snowball 雪球实体
type Snowball struct {
	*Projectile
}

// NewSnowball 创建雪球实体
func NewSnowball(shooterID int64) *Snowball {
	s := &Snowball{
		Projectile: NewProjectile(shooterID),
	}

	s.Entity.NetworkID = SnowballNetworkID
	s.Entity.Width = 0.25
	s.Entity.Height = 0.25
	s.Entity.Gravity = 0.03
	s.Entity.Drag = 0.01
	s.Projectile.BaseDamage = 0 // 雪球不造成伤害
	s.Projectile.MaxAge = 1200

	return s
}

// ThrowableTickResult 投掷物 tick 结果
type ThrowableTickResult struct {
	ProjectileTickResult
	ShouldSpawnChicken bool    // 是否应生成小鸡（仅 Egg）
	SpawnX             float64 // 小鸡生成位置
	SpawnY             float64
	SpawnZ             float64
}

// TickSnowball 雪球 tick 逻辑
// 对应 PHP Snowball::entityBaseTick()
// 碰撞或落地即消失
func (s *Snowball) TickSnowball(nearbyEntities []IEntity, isCollided bool) ThrowableTickResult {
	base := s.Projectile.TickProjectile(nearbyEntities, isCollided)

	result := ThrowableTickResult{
		ProjectileTickResult: base,
	}

	// 碰撞或落地即消失
	if isCollided || s.Entity.OnGround {
		result.ShouldClose = true
	}

	return result
}

// ============ Egg 鸡蛋 ============

// Egg 鸡蛋实体
type Egg struct {
	*Projectile
}

// NewEgg 创建鸡蛋实体
func NewEgg(shooterID int64) *Egg {
	e := &Egg{
		Projectile: NewProjectile(shooterID),
	}

	e.Entity.NetworkID = EggNetworkID
	e.Entity.Width = 0.25
	e.Entity.Height = 0.25
	e.Entity.Gravity = 0.03
	e.Entity.Drag = 0.01
	e.Projectile.BaseDamage = 0 // 鸡蛋不造成伤害
	e.Projectile.MaxAge = 1200

	return e
}

// TickEgg 鸡蛋 tick 逻辑
// 对应 PHP Egg::entityBaseTick()
// 碰撞即消失，12.5% 概率生成小鸡
//
// 参数:
//   - shouldSpawnChicken: 由调用方传入随机结果（rand.Intn(8)==0），
//     避免在纯逻辑层引入随机依赖
func (e *Egg) TickEgg(nearbyEntities []IEntity, isCollided bool, shouldSpawnChicken bool) ThrowableTickResult {
	base := e.Projectile.TickProjectile(nearbyEntities, isCollided)

	result := ThrowableTickResult{
		ProjectileTickResult: base,
	}

	// 碰撞即消失
	if isCollided {
		result.ShouldClose = true

		// 12.5% 概率生成小鸡
		if shouldSpawnChicken {
			result.ShouldSpawnChicken = true
			result.SpawnX = e.Entity.Position.X
			result.SpawnY = e.Entity.Position.Y
			result.SpawnZ = e.Entity.Position.Z
		}
	}

	return result
}
