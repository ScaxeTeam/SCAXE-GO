package entity

import "math/rand"

type AIWander struct {
	BaseAIGoal
	entity          *Entity
	moveHelper      *MoveHelper
	x, y, z         float64
	speed           float64
	executionChance int
	mustUpdate      bool
}

func NewAIWander(e *Entity, mh *MoveHelper, speed float64) *AIWander {
	return NewAIWanderWithChance(e, mh, speed, 20)
}

func NewAIWanderWithChance(e *Entity, mh *MoveHelper, speed float64, chance int) *AIWander {
	g := &AIWander{
		entity:          e,
		moveHelper:      mh,
		speed:           speed,
		executionChance: chance,
	}
	g.SetMutexBits(1)
	return g
}
func (g *AIWander) ShouldExecute() bool {
	if !g.mustUpdate {
		if rand.Intn(g.executionChance) != 0 {
			return false
		}
	}

	pos := g.getPosition()
	if pos == nil {
		return false
	}

	g.x = pos.X
	g.y = pos.Y
	g.z = pos.Z
	g.mustUpdate = false
	return true
}
func (g *AIWander) getPosition() *Vector3 {
	e := g.entity
	x := e.Position.X + float64(rand.Intn(21)-10)
	z := e.Position.Z + float64(rand.Intn(21)-10)
	y := e.Position.Y

	if e.Level != nil {
		groundY := e.Level.FindGroundY(int32(x), int32(z), int32(y)+4)
		if groundY > 0 {
			y = float64(groundY)
		}
	}

	return NewVector3(x, y, z)
}
func (g *AIWander) ShouldContinueExecuting() bool {
	return g.moveHelper.IsMoving
}
func (g *AIWander) StartExecuting() {
	g.moveHelper.SetMoveTo(g.x, g.y, g.z, g.speed)
}

func (g *AIWander) MakeUpdate() {
	g.mustUpdate = true
}

func (g *AIWander) SetExecutionChance(chance int) {
	g.executionChance = chance
}
