package entity

// boat.go — 船载具实体
// 对应 PHP: entity/Boat.php (~119行)
//
// 继承 Entity（PHP 中继承 Vehicle → Entity），核心功能:
//   - WoodID 木材类型（0-5，对应6种木材）
//   - 水中浮力（在水中 motionY=+0.1，否则 motionY=-0.08）
//   - 无乘客时 1500 tick 后消失
//   - 有乘客时跟随玩家朝向（yaw 变化>5° 时同步）
//   - 掉落: 对应木材类型的船物品

import (
	"math"
)

// ============ 常量 ============

const (
	BoatNetworkID = 90

	// 船的木材类型
	BoatWoodOak     = 0
	BoatWoodSpruce  = 1
	BoatWoodBirch   = 2
	BoatWoodJungle  = 3
	BoatWoodAcacia  = 4
	BoatWoodDarkOak = 5

	// 船无乘客超时（tick）
	BoatDespawnAge = 1500

	// 朝向同步阈值（度）
	BoatYawThreshold = 5.0
)

// ============ Boat 实体 ============

// Boat 船载具实体
type Boat struct {
	*Entity

	// WoodID 木材类型 (0-5)
	WoodID int

	// LinkedEntityID 乘客实体ID（0=无乘客）
	LinkedEntityID int64

	// BoatAge 无乘客累计 tick（有乘客时重置）
	BoatAge int
}

// NewBoat 创建船实体
func NewBoat(woodID int) *Boat {
	b := &Boat{
		Entity:  NewEntity(),
		WoodID:  woodID,
		BoatAge: 0,
	}

	b.Entity.NetworkID = BoatNetworkID
	b.Entity.Width = 1.6
	b.Entity.Height = 0.7
	b.Entity.Gravity = 0.5
	b.Entity.Drag = 0.1
	b.Entity.MaxHealth = 10
	b.Entity.Health = 10

	return b
}

// ============ Tick ============

// BoatTickResult 船 tick 结果
type BoatTickResult struct {
	HasUpdate   bool    // 需要同步
	ShouldClose bool    // 应该消失
	YawUpdate   bool    // 是否需要同步朝向
	NewYaw      float64 // 新朝向
}

// TickBoat 船的逻辑 tick
// 对应 PHP Boat::entityBaseTick()
//
// 参数:
//   - riderYaw: 乘客的朝向（无乘客时传 0）
//   - hasRider: 是否有乘客
func (b *Boat) TickBoat(riderYaw float64, hasRider bool) BoatTickResult {
	result := BoatTickResult{}

	b.Entity.TicksLived++

	if !hasRider {
		b.BoatAge++
		// 无乘客 1500 tick 后消失
		if b.BoatAge > BoatDespawnAge {
			result.ShouldClose = true
			result.HasUpdate = true
			return result
		}
	} else {
		b.BoatAge = 0

		// 跟随乘客朝向（变化 > 5° 时同步）
		if math.Abs(riderYaw-b.Entity.Yaw) > BoatYawThreshold {
			b.Entity.Yaw = riderYaw
			result.YawUpdate = true
			result.NewYaw = riderYaw
			result.HasUpdate = true
		}
	}

	return result
}

// ============ 重力/浮力 ============

// BoatGravityResult 船的重力/浮力结果
type BoatGravityResult struct {
	MotionY float64 // 应用的 Y 轴运动
}

// ApplyBoatGravity 应用船的重力/浮力
// 对应 PHP Boat::applyGravity()
//
// 参数:
//   - isInWater: 船是否在水中
//   - hasBlockBelow: 船下方是否有碰撞箱（站在方块上）
func ApplyBoatGravity(isInWater bool, hasBlockBelow bool) BoatGravityResult {
	if hasBlockBelow || isInWater {
		return BoatGravityResult{MotionY: 0.1} // 浮力
	}
	return BoatGravityResult{MotionY: -0.08} // 重力
}

// ============ 掉落 ============

// BoatItemID 船物品的基础ID
const BoatItemID = 333 // item.BOAT

// GetBoatDropItemID 获取船掉落的物品ID
// 返回 (itemID, itemMeta)
func (b *Boat) GetBoatDropItemID() (int, int) {
	return BoatItemID, b.WoodID
}

// ============ 辅助 ============

// SetRider 设置乘客
func (b *Boat) SetRider(entityID int64) {
	b.LinkedEntityID = entityID
	b.BoatAge = 0
}

// RemoveRider 移除乘客
func (b *Boat) RemoveRider() {
	b.LinkedEntityID = 0
}

// HasRider 是否有乘客
func (b *Boat) HasRider() bool {
	return b.LinkedEntityID != 0
}

// GetWoodID 获取木材类型
func (b *Boat) GetWoodID() int {
	return b.WoodID
}
