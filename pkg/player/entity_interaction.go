package player

// entity_interaction.go — 玩家实体交互处理
// 对应 PHP Player.php 的 INTERACT_PACKET / ANIMATE_PACKET 处理器
//
// 三大交互类型：
//   1. 攻击实体 — 伤害计算 + 附魔加成 + 暴击 + 击退 + 工具耐久
//   2. 骑乘交互 — 上船/矿车、离开载具
//   3. 右键实体 — 动物喂食/剪毛/挤奶、NPC 交互

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/entity"
	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

// ── 交互动作常量 ──────────────────────────────────

const (
	InteractActionLeftClick    byte = 1 // 左键攻击
	InteractActionRightClick   byte = 2 // 右键交互
	InteractActionLeaveVehicle byte = 3 // 离开载具

	// 攻击冷却 (ticks)
	AttackCooldownTicks = 10

	// 默认击退力度
	DefaultKnockback = 0.4
)

// ── 攻击状态 ──────────────────────────────────────

// CombatState 存储玩家战斗相关状态
type CombatState struct {
	CPS            int   // 每秒点击数计数器
	AttackCooldown int   // 攻击冷却剩余 tick
	MaxCPS         int   // 最大允许 CPS
	LastAttackTick int64 // 上次攻击时间
}

// newCombatState 创建战斗状态
func newCombatState() *CombatState {
	return &CombatState{
		MaxCPS: 20,
	}
}

// ── InteractPacket 处理 ──────────────────────────

// HandleInteract 处理实体交互包
// 对应 PHP Player::handleDataPacket() case INTERACT_PACKET (L3091-3283)
func (p *Player) HandleInteract(targetEID int64, action byte) {
	if !p.Spawned || !p.Connected {
		return
	}

	lvl, ok := p.Human.Level.(*level.Level)
	if !ok || lvl == nil {
		return
	}

	// 查找目标实体
	target := lvl.GetEntityByID(targetEID)
	if target == nil {
		return
	}

	switch action {
	case InteractActionRightClick:
		p.handleEntityRightClick(target)

	case InteractActionLeftClick:
		p.handleEntityAttack(target)

	case InteractActionLeaveVehicle:
		p.handleLeaveVehicle(target)
	}
}

// ── 实体攻击 ──────────────────────────────────────

// handleEntityAttack 攻击实体
// 对应 PHP INTERACT_PACKET ACTION_LEFT_CLICK (L3192-3281)
func (p *Player) handleEntityAttack(target entity.IEntity) {
	if p.IsSpectator() {
		return
	}

	// CPS 限制
	p.combat.CPS++
	if p.combat.CPS > p.combat.MaxCPS {
		logger.Debug("CPS exceeded", "player", p.Username, "cps", p.combat.CPS)
		return
	}

	// 攻击冷却
	if p.combat.AttackCooldown > 0 {
		return
	}

	// 距离验证
	maxDist := 4.1
	if p.IsCreative() {
		maxDist = 8.0
	}

	targetPos := target.GetPosition()
	if !p.canInteract(targetPos.X, targetPos.Y, targetPos.Z, maxDist) {
		return
	}

	// 基础伤害
	heldItem := p.Inventory.GetItemInHand()
	damageBase := 1.0 // 默认空手伤害
	_ = heldItem      // TODO: 从 Item 获取攻击点数

	// 击退力度
	knockback := DefaultKnockback

	// TODO: 附魔加成
	// - 锋利 (Sharpness): damageBase += floor(level * 1.25)
	// - 火焰附加 (Fire Aspect): 点燃目标
	// - 击退 (Knockback): knockback += 0.1 * level

	// 暴击检测: 下落中 + 非地面 + 非水中
	isCritical := false
	if !p.IsOnGround() && p.movement.SpeedY < 0 && !p.IsSwimming() {
		damageBase *= 1.5
		isCritical = true
	}

	// 应用伤害（直接操作 Entity.Health）
	if ent, ok := target.(*entity.Entity); ok {
		ent.Health -= int(damageBase)
		if ent.Health < 0 {
			ent.Health = 0
		}
	}

	// 应用击退
	dx := targetPos.X - p.Position.X
	dz := targetPos.Z - p.Position.Z
	dist := math.Sqrt(dx*dx + dz*dz)
	if dist > 0 {
		kbX := dx / dist * knockback
		kbZ := dz / dist * knockback
		if ent, ok := target.(*entity.Entity); ok {
			ent.SetMotion(entity.NewVector3(kbX, knockback*0.4, kbZ))
		}
	}

	// 暴击特效
	if isCritical {
		p.broadcastCriticalHit(target)
	}

	// 设置攻击冷却
	p.combat.AttackCooldown = AttackCooldownTicks

	// 生存模式: 工具耐久消耗 + 饥饿消耗
	if p.IsSurvival() {
		// TODO: item.useOn(target) — 工具耐久
		// TODO: p.exhaust(0.1, CauseAttack)
	}

	logger.Debug("Entity attacked",
		"player", p.Username,
		"target", target.GetID(),
		"damage", damageBase,
		"critical", isCritical)
}

// ── 骑乘交互 ──────────────────────────────────────

// handleEntityRightClick 右键实体交互
// 对应 PHP INTERACT_PACKET ACTION_RIGHT_CLICK (L3109-3189)
func (p *Player) handleEntityRightClick(target entity.IEntity) {
	if p.IsSpectator() {
		return
	}

	// TODO: 检查是否为可骑乘实体（船/矿车）并处理骑乘
	// TODO: 动物交互（喂食/剪毛/挤奶等）

	logger.Debug("Entity right-click",
		"player", p.Username,
		"target", target.GetID())
}

// handleLeaveVehicle 离开载具
// 对应 PHP L3144-3147
func (p *Player) handleLeaveVehicle(target entity.IEntity) {
	p.Human.Metadata.SetFlag(entity.DataFlags, entity.DataFlagRiding, false)

	logger.Debug("Left vehicle",
		"player", p.Username,
		"target", target.GetID())
}

// ── 动画广播 ──────────────────────────────────────

// HandleAnimate 处理 AnimatePacket（挥手等动画）
// 对应 PHP ANIMATE_PACKET (L3284-3299)
func (p *Player) HandleAnimate(animAction byte) {
	if !p.Spawned || !p.Connected {
		return
	}

	viewers := p.getViewers()
	if len(viewers) == 0 {
		return
	}

	pk := protocol.NewAnimatePacket()
	pk.EntityID = p.GetID()
	pk.Action = animAction

	for _, viewer := range viewers {
		if viewer != p {
			viewer.SendPacket(pk)
		}
	}
}

// broadcastCriticalHit 广播暴击特效
func (p *Player) broadcastCriticalHit(target entity.IEntity) {
	viewers := p.getViewers()

	pk := protocol.NewAnimatePacket()
	pk.EntityID = target.GetID()
	pk.Action = protocol.AnimateActionCriticalHit

	for _, viewer := range viewers {
		viewer.SendPacket(pk)
	}
}

// ── 战斗 Tick ──────────────────────────────────────

// tickCombat 每 tick 更新战斗状态（冷却递减）
// 应在 Player.Tick() 中调用
func (p *Player) tickCombat() {
	if p.combat.AttackCooldown > 0 {
		p.combat.AttackCooldown--
	}
}

// ResetCPS 重置 CPS 计数器（每秒调用一次）
func (p *Player) ResetCPS() {
	p.combat.CPS = 0
}

// ── 辅助方法 ──────────────────────────────────────

// IsAlive 返回玩家是否存活
func (p *Player) IsAlive() bool {
	return p.Human.Health > 0
}
