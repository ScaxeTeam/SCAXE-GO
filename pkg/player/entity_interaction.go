package player

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/entity"
	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

const (
	InteractActionLeftClick    byte = 1
	InteractActionRightClick   byte = 2
	InteractActionLeaveVehicle byte = 3

	AttackCooldownTicks = 10

	DefaultKnockback = 0.4
)

type CombatState struct {
	CPS            int
	AttackCooldown int
	MaxCPS         int
	LastAttackTick int64
}

func newCombatState() *CombatState {
	return &CombatState{
		MaxCPS: 20,
	}
}

func (p *Player) HandleInteract(targetEID int64, action byte) {
	if !p.Spawned || !p.Connected {
		return
	}

	lvl, ok := p.Human.Level.(*level.Level)
	if !ok || lvl == nil {
		return
	}

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

func (p *Player) handleEntityAttack(target entity.IEntity) {
	if p.IsSpectator() {
		return
	}

	p.combat.CPS++
	if p.combat.CPS > p.combat.MaxCPS {
		logger.Debug("CPS exceeded", "player", p.Username, "cps", p.combat.CPS)
		return
	}

	if p.combat.AttackCooldown > 0 {
		return
	}

	maxDist := 4.1
	if p.IsCreative() {
		maxDist = 8.0
	}

	targetPos := target.GetPosition()
	if !p.canInteract(targetPos.X, targetPos.Y, targetPos.Z, maxDist) {
		return
	}

	heldItem := p.Inventory.GetItemInHand()
	damageBase := 1.0
	_ = heldItem

	knockback := DefaultKnockback

	isCritical := false
	if !p.IsOnGround() && p.movement.SpeedY < 0 && !p.IsSwimming() {
		damageBase *= 1.5
		isCritical = true
	}

	if ent, ok := target.(*entity.Entity); ok {
		ent.Health -= int(damageBase)
		if ent.Health < 0 {
			ent.Health = 0
		}
	}

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

	if isCritical {
		p.broadcastCriticalHit(target)
	}

	p.combat.AttackCooldown = AttackCooldownTicks

	if p.IsSurvival() {

	}

	logger.Debug("Entity attacked",
		"player", p.Username,
		"target", target.GetID(),
		"damage", damageBase,
		"critical", isCritical)
}

func (p *Player) handleEntityRightClick(target entity.IEntity) {
	if p.IsSpectator() {
		return
	}

	logger.Debug("Entity right-click",
		"player", p.Username,
		"target", target.GetID())
}

func (p *Player) handleLeaveVehicle(target entity.IEntity) {
	p.Human.Metadata.SetFlag(entity.DataFlags, entity.DataFlagRiding, false)

	logger.Debug("Left vehicle",
		"player", p.Username,
		"target", target.GetID())
}

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

func (p *Player) broadcastCriticalHit(target entity.IEntity) {
	viewers := p.getViewers()

	pk := protocol.NewAnimatePacket()
	pk.EntityID = target.GetID()
	pk.Action = protocol.AnimateActionCriticalHit

	for _, viewer := range viewers {
		viewer.SendPacket(pk)
	}
}

func (p *Player) tickCombat() {
	if p.combat.AttackCooldown > 0 {
		p.combat.AttackCooldown--
	}
}

func (p *Player) ResetCPS() {
	p.combat.CPS = 0
}

func (p *Player) IsAlive() bool {
	return p.Human.Health > 0
}
