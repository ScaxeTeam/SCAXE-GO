package entity

import "math/rand"

type AIPanic struct {
	BaseAIGoal
	entity     *Entity
	moveHelper *MoveHelper
	speed      float64
	randPosX   float64
	randPosY   float64
	randPosZ   float64
}

func NewAIPanic(e *Entity, mh *MoveHelper, speed float64) *AIPanic {
	g := &AIPanic{
		entity:     e,
		moveHelper: mh,
		speed:      speed,
	}
	g.SetMutexBits(1)
	return g
}
func (g *AIPanic) ShouldExecute() bool {
	if g.entity.NoDamageTicks <= 0 && g.entity.FireTicks <= 0 {
		return false
	}

	return g.findRandomPosition()
}
func (g *AIPanic) findRandomPosition() bool {
	e := g.entity
	for i := 0; i < 5; i++ {
		x := e.Position.X + float64(rand.Intn(11)-5)
		y := e.Position.Y + float64(rand.Intn(9)-4)
		z := e.Position.Z + float64(rand.Intn(11)-5)

		if y < 1 {
			y = 1
		}
		if y > 127 {
			y = 127
		}

		g.randPosX = x
		g.randPosY = y
		g.randPosZ = z
		return true
	}
	return false
}

func (g *AIPanic) StartExecuting() {
	g.moveHelper.SetMoveTo(g.randPosX, g.randPosY, g.randPosZ, g.speed)
}
func (g *AIPanic) ShouldContinueExecuting() bool {
	return g.moveHelper.IsMoving
}
