package player

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/entity"
	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

const (
	EntityEventHurtAnimation byte = 2

	EntityEventDeathAnimation byte = 3

	EntityEventRespawn byte = 18

	FoodTickCycle = 80

	ExhaustionPerSprint = 0.1

	ExhaustionPerJump = 0.05

	ExhaustionPerSprintJump = 0.2

	ExhaustionPerHealthRegen = 3.0
)

type SurvivalState struct {
	dead      bool
	deathTime int
}

func newSurvivalState() *SurvivalState {
	return &SurvivalState{}
}

func (p *Player) tickSurvival() {

	if p.survival.dead {
		p.survival.deathTime++
		return
	}

	if p.GetHealth() <= 0 {
		p.onDeath()
		return
	}

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

	p.Human.FoodTickTimer++
	if p.Human.FoodTickTimer >= FoodTickCycle {
		p.Human.FoodTickTimer = 0
	}

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

	if p.Human.FoodTickTimer == 0 {
		if food >= 18 {

			if health < maxHealth {
				p.Heal(1)
				p.Exhaust(ExhaustionPerHealthRegen)
			}
		} else if food <= 0 {

			shouldDamage := false
			switch difficulty {
			case 1:
				shouldDamage = health > 10
			case 2:
				shouldDamage = health > 1
			case 3:
				shouldDamage = true
			}
			if shouldDamage {
				p.Attack(1, entity.DamageCauseStarvation)
				p.sendHurtAnimation()
			}
		}
	}

	if food <= 6 {
		if p.IsSprinting() {
			p.SetSprinting(false)
		}
	}
}

func (p *Player) onDeath() {
	p.survival.dead = true
	p.survival.deathTime = 0

	logger.Info("Player died", "player", p.Username)

	p.broadcastEntityEvent(EntityEventDeathAnimation)

	if p.Gamemode == 0 {
		p.dropAllItems()
	}

	respawnPk := protocol.NewRespawnPacket()
	if lvl, ok := p.Human.Level.(*level.Level); ok {
		spawn := lvl.GetSafeSpawn()
		respawnPk.X = float32(spawn.X)
		respawnPk.Y = float32(spawn.Y)
		respawnPk.Z = float32(spawn.Z)
	}
	p.SendPacket(respawnPk)
}

func (p *Player) dropAllItems() {
	if p.Inventory == nil {
		return
	}

	contents := p.Inventory.GetContents()
	for slot, it := range contents {
		if it.ID == 0 || it.Count == 0 {
			continue
		}

		p.Inventory.ClearSlot(slot, false)
	}
}

func (p *Player) handleRespawn() {
	if !p.survival.dead {
		return
	}

	p.survival.dead = false
	p.survival.deathTime = 0

	p.SetHealth(p.GetMaxHealth())
	p.SetFood(20)
	p.SetSaturation(20)
	p.SetExhaustion(0)
	p.Human.FoodTickTimer = 0

	lvl, ok := p.Human.Level.(*level.Level)
	if !ok || lvl == nil {
		return
	}

	spawn := lvl.GetSafeSpawn()
	p.Teleport(spawn.X, spawn.Y, spawn.Z)

	healthPk := protocol.NewSetHealthPacket()
	healthPk.Health = int32(p.GetMaxHealth())
	p.SendPacket(healthPk)

	p.broadcastEntityEvent(EntityEventRespawn)

	p.sendInventoryContents()

	logger.Info("Player respawned", "player", p.Username)
}

func (p *Player) IsDead() bool {
	return p.survival.dead
}

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

func (p *Player) sendHurtAnimation() {
	p.broadcastEntityEvent(EntityEventHurtAnimation)
}

func (p *Player) broadcastEntityEvent(event byte) {
	pk := protocol.NewEntityEventPacket()
	pk.EntityID = p.GetID()
	pk.Event = event

	selfPk := protocol.NewEntityEventPacket()
	selfPk.EntityID = 0
	selfPk.Event = event
	p.SendPacket(selfPk)

	viewers := p.getViewers()
	for _, viewer := range viewers {
		if viewer != p {
			viewer.SendPacket(pk)
		}
	}
}

func (p *Player) ExhaustFromSprint(horizontalDist float64) {
	if p.Gamemode != 0 || !p.Human.FoodEnabled {
		return
	}
	p.Exhaust(ExhaustionPerSprint * horizontalDist)
}

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

func (p *Player) Exhaust(amount float64) {
	if p.Gamemode != 0 || !p.Human.FoodEnabled {
		return
	}

	oldFood := p.GetFood()
	p.Human.Exhaust(amount)
	newFood := p.GetFood()

	if math.Abs(oldFood-newFood) > 0.001 {
		p.syncFoodAttributes()
	}
}

func (p *Player) syncFoodAttributes() {

}
