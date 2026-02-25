package entity

// arrow.go — 箭实体
// 对应 PHP: entity/Arrow.php (~110行)
//
// 继承 Projectile，具体属性:
//   - NetworkID = 80
//   - gravity = 0.05, drag = 0.01
//   - baseDamage = 2.0
//   - punchKnockback（冲击附魔击退加成）
//   - isCritical（满蓄力暴击，飞行时显示粒子）
//   - 1200 tick 后消失

import (
	"math"
)

// ArrowNetworkID 箭实体的网络ID
const ArrowNetworkID = 80

// ============ Arrow 实体 ============

// Arrow 箭实体
type Arrow struct {
	*Projectile

	// PunchKnockback 冲击附魔的击退加成
	PunchKnockback float64
}

// NewArrow 创建箭实体
// 参数:
//   - shooterID: 发射者实体ID
//   - critical: 是否为暴击箭（满蓄力）
func NewArrow(shooterID int64, critical bool) *Arrow {
	a := &Arrow{
		Projectile:     NewProjectile(shooterID),
		PunchKnockback: 0,
	}

	// 箭的特定属性
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

// ============ Tick ============

// ArrowTickResult 箭 tick 结果（扩展 ProjectileTickResult）
type ArrowTickResult struct {
	ProjectileTickResult
	ShowCriticalParticle bool    // 是否显示暴击粒子
	ParticleX            float64 // 粒子位置
	ParticleY            float64
	ParticleZ            float64
}

// TickArrow 箭的每 tick 逻辑
// 对应 PHP Arrow::entityBaseTick()
func (a *Arrow) TickArrow(nearbyEntities []IEntity, isCollided bool) ArrowTickResult {
	base := a.Projectile.TickProjectile(nearbyEntities, isCollided)

	result := ArrowTickResult{
		ProjectileTickResult: base,
	}

	// 暴击粒子效果
	if !a.Projectile.HadCollision && a.Projectile.IsCritical {
		result.ShowCriticalParticle = true
		result.ParticleX = a.Entity.Position.X
		result.ParticleY = a.Entity.Position.Y + a.Entity.Height/2
		result.ParticleZ = a.Entity.Position.Z
	} else if a.Entity.OnGround {
		// 落地后取消暴击状态
		a.Projectile.IsCritical = false
	}

	return result
}

// ============ 命中处理 ============

// ArrowHitResult 箭命中结果
type ArrowHitResult struct {
	Damage       float64 // 伤害
	KnockbackX   float64 // 额外击退 X
	KnockbackY   float64 // 额外击退 Y
	KnockbackZ   float64 // 额外击退 Z
	HasKnockback bool    // 是否有额外击退
	TransferFire bool    // 是否传递燃烧
	FireDuration int     // 燃烧持续时间(tick)
}

// CalcArrowHit 计算箭命中效果
// 对应 PHP Arrow::onHit()
func (a *Arrow) CalcArrowHit() ArrowHitResult {
	result := ArrowHitResult{
		Damage: a.Projectile.CalcHitDamage(),
	}

	// 冲击附魔击退
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

	// 燃烧传递
	if a.Projectile.ShouldTransferFire() {
		result.TransferFire = true
		result.FireDuration = a.Projectile.GetFireTransferDuration()
	}

	return result
}

// ============ 辅助 ============

// SetPunchKnockback 设置冲击附魔击退值
func (a *Arrow) SetPunchKnockback(value float64) {
	a.PunchKnockback = value
}

// GetPunchKnockback 获取冲击附魔击退值
func (a *Arrow) GetPunchKnockback() float64 {
	return a.PunchKnockback
}
