package entity

// projectile.go — 投射物基类
// 对应 PHP: entity/Projectile.php (~200行)
//
// 投射物继承 Entity，核心功能:
//   - shootingEntity 追踪发射者
//   - entityBaseTick 碰撞检测（方块碰撞停止/实体射线检测最近目标）
//   - onHit 命中实体时计算伤害（速度×基础伤害+暴击随机）、燃烧传递
//   - 自动朝向更新（yaw/pitch 随运动方向）
//   - 不受除虚空外的伤害
//   - 先应用阻力再应用重力（applyDragBeforeGravity）

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

// ============ 投射物基类 ============

// Projectile 投射物实体基类
// 嵌入 Entity，提供发射者追踪、碰撞检测、命中处理
type Projectile struct {
	*Entity

	// ShootingEntityID 发射者实体ID（0=无发射者）
	ShootingEntityID int64

	// BaseDamage 基础伤害（命中伤害 = 速度 × BaseDamage）
	BaseDamage float64

	// HadCollision 是否已经与方块碰撞（插在方块里）
	HadCollision bool

	// IsCritical 是否为暴击（满蓄力箭）
	IsCritical bool

	// DragBeforeGravity 是否先应用阻力再应用重力（投射物默认true）
	DragBeforeGravity bool

	// Age 投射物存活时间（tick），超时后消失
	ProjectileAge int

	// MaxAge 最大存活时间（tick），默认1200（60秒）
	MaxAge int
}

// NewProjectile 创建投射物实体
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

	// 投射物默认属性
	p.Entity.Health = 1
	p.Entity.MaxHealth = 1
	p.Entity.CanCollide = false
	p.Entity.Width = 0.25
	p.Entity.Height = 0.25
	p.Entity.Gravity = 0.03
	p.Entity.Drag = 0.01

	return p
}

// ============ Tick 逻辑 ============

// ProjectileTickResult 投射物 tick 结果
type ProjectileTickResult struct {
	HasUpdate   bool    // 是否有更新（需要同步客户端）
	HitEntityID int64   // 命中的实体ID（0=未命中）
	HitBlock    bool    // 是否新碰撞到方块
	ShouldClose bool    // 是否应该关闭（超时/命中）
	Damage      float64 // 对命中实体的伤害
}

// TickProjectile 投射物的每 tick 逻辑
// 对应 PHP Projectile::entityBaseTick()
//
// 参数:
//   - nearbyEntities: 附近实体列表（已由 Level 层过滤范围）
//   - isCollided: 当前 tick 是否与方块碰撞
func (p *Projectile) TickProjectile(nearbyEntities []IEntity, isCollided bool) ProjectileTickResult {
	result := ProjectileTickResult{}

	p.Entity.TicksLived++
	p.ProjectileAge++

	// 超时消失
	if p.ProjectileAge >= p.MaxAge {
		result.ShouldClose = true
		return result
	}

	// ====== 方块碰撞处理 ======
	if isCollided && !p.HadCollision {
		// 首次碰撞 → 停止运动
		p.HadCollision = true
		p.Entity.Motion.X = 0
		p.Entity.Motion.Y = 0
		p.Entity.Motion.Z = 0
		result.HitBlock = true
		result.HasUpdate = true
	} else if !isCollided && p.HadCollision {
		// 碰撞的方块被移除 → 恢复物理
		p.HadCollision = false
	}

	// ====== 实体碰撞检测（射线检测最近目标）======
	if !p.HadCollision {
		nearDist := math.MaxFloat64
		var nearEntityID int64

		moveX := p.Entity.Position.X + p.Entity.Motion.X
		moveY := p.Entity.Position.Y + p.Entity.Motion.Y
		moveZ := p.Entity.Position.Z + p.Entity.Motion.Z

		for _, ent := range nearbyEntities {
			entID := ent.GetID()

			// 跳过发射者（前5 tick）
			if entID == p.ShootingEntityID && p.Entity.TicksLived < 5 {
				continue
			}

			// 简化的碰撞：距离检查
			entPos := ent.GetPosition()
			dx := entPos.X - moveX
			dy := entPos.Y - moveY
			dz := entPos.Z - moveZ
			distSq := dx*dx + dy*dy + dz*dz

			// 碰撞半径 = 实体碰撞箱膨胀 0.3
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
			// 命中实体
			result.HitEntityID = nearEntityID
			result.Damage = p.CalcHitDamage()
			result.ShouldClose = true
			result.HasUpdate = true
			p.HadCollision = true
			return result
		}
	}

	// ====== 朝向更新 ======
	mx := p.Entity.Motion.X
	my := p.Entity.Motion.Y
	mz := p.Entity.Motion.Z
	if !p.Entity.OnGround || math.Abs(mx) > 0.00001 || math.Abs(my) > 0.00001 || math.Abs(mz) > 0.00001 {
		f := math.Sqrt(mx*mx + mz*mz)
		p.Entity.Yaw = math.Atan2(mx, mz) * 180 / math.Pi
		p.Entity.Pitch = math.Atan2(my, f) * 180 / math.Pi
		result.HasUpdate = true
	}

	// ====== 物理：阻力 + 重力 ======
	if p.DragBeforeGravity {
		// 先阻力
		p.Entity.Motion.X *= 1 - p.Entity.Drag
		p.Entity.Motion.Y *= 1 - p.Entity.Drag
		p.Entity.Motion.Z *= 1 - p.Entity.Drag
		// 后重力
		p.Entity.Motion.Y -= p.Entity.Gravity
	} else {
		// 先重力
		p.Entity.Motion.Y -= p.Entity.Gravity
		// 后阻力
		p.Entity.Motion.X *= 1 - p.Entity.Drag
		p.Entity.Motion.Y *= 1 - p.Entity.Drag
		p.Entity.Motion.Z *= 1 - p.Entity.Drag
	}

	return result
}

// ============ 伤害计算 ============

// CalcHitDamage 计算投射物命中伤害
// 对应 PHP Projectile::onHit()
// 伤害 = ceil(速度 × BaseDamage)，暴击时 +random(0, damage/2+1)
func (p *Projectile) CalcHitDamage() float64 {
	mx := p.Entity.Motion.X
	my := p.Entity.Motion.Y
	mz := p.Entity.Motion.Z
	speed := math.Sqrt(mx*mx + my*my + mz*mz)
	damage := math.Ceil(speed * p.BaseDamage)

	if p.IsCritical && damage > 0 {
		// 暴击：+random(0, damage/2+1)
		bonus := int(damage/2) + 1
		if bonus > 0 {
			damage += float64(bonus / 2) // 简化: 取中间值
		}
	}

	return damage
}

// ============ 辅助 ============

// SetShooter 设置发射者
func (p *Projectile) SetShooter(entityID int64) {
	p.ShootingEntityID = entityID
}

// IsOnGround 投射物是否在地面（插在方块里）
func (p *Projectile) IsOnGround() bool {
	return p.HadCollision
}

// GetSpeed 获取当前速度
func (p *Projectile) GetSpeed() float64 {
	mx := p.Entity.Motion.X
	my := p.Entity.Motion.Y
	mz := p.Entity.Motion.Z
	return math.Sqrt(mx*mx + my*my + mz*mz)
}

// ShouldTransferFire 是否应该传递燃烧（投射物着火→命中目标着火）
// 对应 PHP Projectile::onHit() 中的 fireTicks 检查
func (p *Projectile) ShouldTransferFire() bool {
	return p.Entity.FireTicks > 0
}

// GetFireTransferDuration 获取传递燃烧的持续时间（tick）
func (p *Projectile) GetFireTransferDuration() int {
	return 100 // 5秒
}

// SaveNBT 保存投射物 NBT
func (p *Projectile) SaveProjectileNBT() {
	p.Entity.SaveNBT()
	p.Entity.NamedTag.Set(nbt.NewShortTag("Age", int16(p.ProjectileAge)))
}
