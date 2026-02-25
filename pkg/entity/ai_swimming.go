package entity

import "math/rand"

type AISwimming struct {
	BaseAIGoal
	entity *Entity
}

func NewAISwimming(e *Entity) *AISwimming {
	g := &AISwimming{entity: e}
	g.SetMutexBits(4)
	return g
}

func (g *AISwimming) ShouldExecute() bool {
	return g.entity.IsInWater()
}

func (g *AISwimming) ShouldContinueExecuting() bool {
	return g.entity.IsInWater()
}

func (g *AISwimming) UpdateTask() {
	if rand.Float32() < 0.8 {
		g.entity.Motion.Y = 0.04
	}
}
