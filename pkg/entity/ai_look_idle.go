package entity

import (
	"math"
	"math/rand"
)

type AILookIdle struct {
	BaseAIGoal
	entity     *Entity
	lookHelper *LookHelper
	lookX      float64
	lookZ      float64
	idleTime   int
}

func NewAILookIdle(e *Entity, lh *LookHelper) *AILookIdle {
	g := &AILookIdle{
		entity:     e,
		lookHelper: lh,
	}
	g.SetMutexBits(3)
	return g
}
func (g *AILookIdle) ShouldExecute() bool {
	return rand.Float32() < 0.02
}

func (g *AILookIdle) ShouldContinueExecuting() bool {
	return g.idleTime >= 0
}
func (g *AILookIdle) StartExecuting() {
	angle := math.Pi * 2.0 * rand.Float64()
	g.lookX = math.Cos(angle)
	g.lookZ = math.Sin(angle)
	g.idleTime = 20 + rand.Intn(20)
}
func (g *AILookIdle) UpdateTask() {
	g.idleTime--
	g.lookHelper.SetLookPosition(
		g.entity.Position.X+g.lookX,
		g.entity.Position.Y+g.entity.EyeHeight,
		g.entity.Position.Z+g.lookZ,
		10.0,
		40.0,
	)
}
