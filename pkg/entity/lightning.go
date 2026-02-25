package entity

// lightning.go — 闪电实体
// 对应 PHP: entity/Lightning.php
//
// 闪电是短暂的天气实体:
//   - NetworkID: 93
//   - 存活 20 tick 后自动关闭
//   - 生成时发送 ExplodePacket 视觉效果
//   - 在落点放置火方块（若有 lightningFire 配置）
//   - 对周围 4×3×4 范围内玩家造成 8-20 闪电伤害 + 着火 3-8 秒
//   - 对 Creeper 充能 (setPowered)

import "math/rand"

const LightningNetworkID = 93

// ============ 闪电伤害常量 ============

const (
	LightningDamageMin = 8  // 最小伤害
	LightningDamageMax = 20 // 最大伤害
	LightningFireMin   = 3  // 着火最短时间（秒）
	LightningFireMax   = 8  // 着火最长时间（秒）
	LightningLifetime  = 20 // 存活 tick 数

	// 影响范围（以闪电位置为中心）
	LightningRangeX = 4.0
	LightningRangeY = 3.0
	LightningRangeZ = 4.0

	// ExplodePacket 视觉效果半径
	LightningExplodeRadius = 10.0
)

// ============ Lightning 实体 ============

// Lightning 闪电实体
type Lightning struct {
	*Entity

	// Age 存活 tick 计数
	Age int
}

// NewLightning 创建闪电实体
func NewLightning() *Lightning {
	l := &Lightning{
		Entity: NewEntity(),
		Age:    0,
	}

	l.Entity.NetworkID = LightningNetworkID
	l.Entity.Width = 0.3
	l.Entity.Height = 1.8
	l.Entity.MaxHealth = 2
	l.Entity.Health = 2
	l.Entity.Gravity = 0
	l.Entity.Drag = 0

	return l
}

// GetName 获取名称
func (l *Lightning) GetName() string {
	return "Lightning"
}

// ============ Tick ============

// LightningTickResult 闪电 tick 结果
type LightningTickResult struct {
	ShouldClose bool // 是否应该关闭（超过存活时间）
}

// TickLightning 闪电 tick 逻辑
// 对应 PHP Lightning::entityBaseTick()
func (l *Lightning) TickLightning() LightningTickResult {
	l.Age++
	if l.Age > LightningLifetime {
		return LightningTickResult{ShouldClose: true}
	}
	return LightningTickResult{}
}

// ============ 伤害计算 ============

// CalcLightningDamage 计算闪电伤害（随机 8-20）
func CalcLightningDamage() int {
	return LightningDamageMin + rand.Intn(LightningDamageMax-LightningDamageMin+1)
}

// CalcLightningFireDuration 计算着火时间（随机 3-8 秒）
func CalcLightningFireDuration() int {
	return LightningFireMin + rand.Intn(LightningFireMax-LightningFireMin+1)
}

// ============ 影响范围 ============

// LightningImpactInfo 闪电落地时的影响信息
type LightningImpactInfo struct {
	// ShouldPlaceFire 是否应在落点放火
	ShouldPlaceFire bool

	// FireX/Y/Z 火方块放置坐标
	FireX, FireY, FireZ int

	// DamageRange 伤害范围 (以落点为中心)
	RangeX, RangeY, RangeZ float64
}

// CalcImpact 计算闪电落地影响
// 参数:
//   - lightningFire: 服务器是否启用闪电着火配置
//   - hitBlockSolid: 落点方块是否为固体
//   - hitBlockLiquid: 落点方块是否为液体
//   - x, y, z: 闪电位置
func CalcImpact(lightningFire bool, hitBlockSolid, hitBlockLiquid bool, x, y, z float64) LightningImpactInfo {
	info := LightningImpactInfo{
		RangeX: LightningRangeX,
		RangeY: LightningRangeY,
		RangeZ: LightningRangeZ,
	}

	if !lightningFire {
		return info
	}

	// 液体上不放火
	if hitBlockLiquid {
		return info
	}

	if hitBlockSolid {
		// 固体方块上面放火
		info.ShouldPlaceFire = true
		info.FireX = int(x)
		info.FireY = int(y) + 1
		info.FireZ = int(z)
	} else {
		// 非固体则在当前位置放火
		info.ShouldPlaceFire = true
		info.FireX = int(x)
		info.FireY = int(y)
		info.FireZ = int(z)
	}

	return info
}
