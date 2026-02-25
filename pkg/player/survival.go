package player

// survival.go — 饥饿系统 + 死亡处理 + 重生流程
// 对应 PHP Human.php entityBaseTick() 饥饿部分 + Player.php 死亡/重生
//
// 饥饿机制 (80-tick 循环):
//   - food >= 18: 自然回血 + 消耗疲劳
//   - food <= 6:  强制停止冲刺
//   - food <= 0:  按难度造成饥饿伤害
//   - 和平模式:   自动恢复食物和血量
//
// 死亡流程:
//   1. health <= 0 → 广播死亡动画 → 掉落物品 → 标记死亡
//   2. 客户端发 ActionRespawn → 重置状态 → 传送到出生点

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/entity"
	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

// ── 常量 ──────────────────────────────────────

const (
	// EntityEventHurtAnimation 受伤动画事件 ID
	EntityEventHurtAnimation byte = 2
	// EntityEventDeathAnimation 死亡动画事件 ID
	EntityEventDeathAnimation byte = 3
	// EntityEventRespawn 重生事件 ID
	EntityEventRespawn byte = 18

	// FoodTickCycle 饥饿 tick 周期 (对应 PHP foodTickTimer >= 80)
	FoodTickCycle = 80

	// ExhaustionPerSprint 冲刺每方块疲劳
	ExhaustionPerSprint = 0.1
	// ExhaustionPerJump 普通跳跃疲劳
	ExhaustionPerJump = 0.05
	// ExhaustionPerSprintJump 冲刺跳疲劳
	ExhaustionPerSprintJump = 0.2
	// ExhaustionPerHealthRegen 自然回血疲劳
	ExhaustionPerHealthRegen = 3.0
)

// ── SurvivalState ──────────────────────────────

// SurvivalState 存储玩家生存模式相关状态
type SurvivalState struct {
	dead      bool
	deathTime int
}

func newSurvivalState() *SurvivalState {
	return &SurvivalState{}
}

// ── 饥饿系统 tick ──────────────────────────────

// tickSurvival 每 tick 处理饥饿/死亡逻辑
// 对应 PHP Human::entityBaseTick() L503-543
func (p *Player) tickSurvival() {
	// 死亡状态不处理
	if p.survival.dead {
		p.survival.deathTime++
		return
	}

	// 检测死亡
	if p.GetHealth() <= 0 {
		p.onDeath()
		return
	}

	// 创造/观察者模式不消耗饥饿
	if p.Gamemode == 1 || p.Gamemode == 3 {
		return
	}

	if !p.Human.FoodEnabled {
		return
	}

	food := p.GetFood()
	health := p.GetHealth()
	maxHealth := p.GetMaxHealth()
	difficulty := p.Difficulty

	// 递增 foodTickTimer
	p.Human.FoodTickTimer++
	if p.Human.FoodTickTimer >= FoodTickCycle {
		p.Human.FoodTickTimer = 0
	}

	// 和平难度特殊处理
	if difficulty == 0 {
		if p.Human.FoodTickTimer%10 == 0 {
			if food < 20 {
				p.AddFood(1.0)
			}
			if p.Human.FoodTickTimer%20 == 0 && health < maxHealth {
				p.Heal(1)
			}
		}
	}

	// 每 80 tick 周期触发
	if p.Human.FoodTickTimer == 0 {
		if food >= 18 {
			// 饱食回血
			if health < maxHealth {
				p.Heal(1)
				p.Exhaust(ExhaustionPerHealthRegen)
			}
		} else if food <= 0 {
			// 饥饿伤害 (按难度判定)
			shouldDamage := false
			switch difficulty {
			case 1: // Easy: health > 10 时扣血
				shouldDamage = health > 10
			case 2: // Normal: health > 1 时扣血
				shouldDamage = health > 1
			case 3: // Hard: 无限扣血
				shouldDamage = true
			}
			if shouldDamage {
				p.Attack(1, entity.DamageCauseStarvation)
				p.sendHurtAnimation()
			}
		}
	}

	// 食物 <= 6 停止冲刺
	if food <= 6 {
		if p.IsSprinting() {
			p.SetSprinting(false)
		}
	}
}

// ── 死亡处理 ──────────────────────────────────

// onDeath 处理玩家死亡
func (p *Player) onDeath() {
	p.survival.dead = true
	p.survival.deathTime = 0

	logger.Info("Player died", "player", p.Username)

	// 广播死亡动画
	p.broadcastEntityEvent(EntityEventDeathAnimation)

	// 生存模式掉落物品
	if p.Gamemode == 0 {
		p.dropAllItems()
	}

	// 发送重生包给自己（让客户端显示死亡界面）
	respawnPk := protocol.NewRespawnPacket()
	if lvl, ok := p.Human.Level.(*level.Level); ok {
		spawn := lvl.GetSafeSpawn()
		respawnPk.X = float32(spawn.X)
		respawnPk.Y = float32(spawn.Y)
		respawnPk.Z = float32(spawn.Z)
	}
	p.SendPacket(respawnPk)
}

// dropAllItems 掉落玩家所有物品
func (p *Player) dropAllItems() {
	if p.Inventory == nil {
		return
	}

	contents := p.Inventory.GetContents()
	for slot, it := range contents {
		if it.ID == 0 || it.Count == 0 {
			continue
		}
		// TODO: 在玩家位置生成掉落物实体 (需要 Level.DropItem)
		// 当前仅清空背包
		p.Inventory.ClearSlot(slot, false)
	}
}

// ── 重生流程 ──────────────────────────────────

// handleRespawn 处理玩家重生
func (p *Player) handleRespawn() {
	if !p.survival.dead {
		return
	}

	// 重置状态
	p.survival.dead = false
	p.survival.deathTime = 0

	p.SetHealth(p.GetMaxHealth())
	p.SetFood(20)
	p.SetSaturation(20)
	p.SetExhaustion(0)
	p.Human.FoodTickTimer = 0

	// 传送到出生点
	lvl, ok := p.Human.Level.(*level.Level)
	if !ok || lvl == nil {
		return
	}

	spawn := lvl.GetSafeSpawn()
	p.Teleport(spawn.X, spawn.Y, spawn.Z)

	// 发送血量更新
	healthPk := protocol.NewSetHealthPacket()
	healthPk.Health = int32(p.GetMaxHealth())
	p.SendPacket(healthPk)

	// 广播重生事件
	p.broadcastEntityEvent(EntityEventRespawn)

	// 重新同步背包
	p.sendInventoryContents()

	logger.Info("Player respawned", "player", p.Username)
}

// IsDead 返回玩家是否死亡
func (p *Player) IsDead() bool {
	return p.survival.dead
}

// ── 辅助方法 ──────────────────────────────────

// Heal 恢复生命值（封装 Living.Heal 并发送更新）
func (p *Player) Heal(amount float64) {
	oldHealth := p.GetHealth()
	p.Human.Living.Heal(amount)
	newHealth := p.GetHealth()

	if newHealth != oldHealth {
		healthPk := protocol.NewSetHealthPacket()
		healthPk.Health = int32(newHealth)
		p.SendPacket(healthPk)
	}
}

// Attack 对玩家造成伤害（封装 Living.Attack 并发送更新）
func (p *Player) Attack(damage float64, cause int) bool {
	if p.survival.dead {
		return false
	}

	result := p.Human.Living.Attack(damage, cause)
	if result {
		healthPk := protocol.NewSetHealthPacket()
		healthPk.Health = int32(p.GetHealth())
		p.SendPacket(healthPk)

		if p.GetHealth() <= 0 {
			p.onDeath()
		}
	}
	return result
}

// sendHurtAnimation 发送受伤动画
func (p *Player) sendHurtAnimation() {
	p.broadcastEntityEvent(EntityEventHurtAnimation)
}

// broadcastEntityEvent 广播实体事件给所有观察者（含自己）
func (p *Player) broadcastEntityEvent(event byte) {
	pk := protocol.NewEntityEventPacket()
	pk.EntityID = p.GetID()
	pk.Event = event

	// 发给自己 (EntityID=0 表示自己)
	selfPk := protocol.NewEntityEventPacket()
	selfPk.EntityID = 0
	selfPk.Event = event
	p.SendPacket(selfPk)

	// 广播给其他玩家
	viewers := p.getViewers()
	for _, viewer := range viewers {
		if viewer != p {
			viewer.SendPacket(pk)
		}
	}
}

// ExhaustFromSprint 冲刺时的疲劳消耗
func (p *Player) ExhaustFromSprint(horizontalDist float64) {
	if p.Gamemode != 0 || !p.Human.FoodEnabled {
		return
	}
	p.Exhaust(ExhaustionPerSprint * horizontalDist)
}

// ExhaustFromJump 跳跃时的疲劳消耗
func (p *Player) ExhaustFromJump() {
	if p.Gamemode != 0 || !p.Human.FoodEnabled {
		return
	}
	if p.IsSprinting() {
		p.Exhaust(ExhaustionPerSprintJump)
	} else {
		p.Exhaust(ExhaustionPerJump)
	}
}

// Exhaust 增加疲劳值（封装 Human.Exhaust 并同步属性）
func (p *Player) Exhaust(amount float64) {
	if p.Gamemode != 0 || !p.Human.FoodEnabled {
		return
	}

	oldFood := p.GetFood()
	p.Human.Exhaust(amount)
	newFood := p.GetFood()

	// 食物发生变化时通知属性更新
	if math.Abs(oldFood-newFood) > 0.001 {
		p.syncFoodAttributes()
	}
}

// syncFoodAttributes 同步食物相关属性到客户端
func (p *Player) syncFoodAttributes() {
	// 属性通过 UpdateAttributesPacket 同步
	// 当前框架中属性变更通过 Attribute.SetValue 自动标记 dirty
	// 实际发送在 tick 中处理，此处预留
}
