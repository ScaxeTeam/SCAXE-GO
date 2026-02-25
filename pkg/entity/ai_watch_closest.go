package entity

import (
	"math"
	"math/rand"
)

type AIWatchClosest struct {
	BaseAIGoal
	entity        *Entity
	lookHelper    *LookHelper
	closestEntity *Entity
	maxDistance   float64
	lookTime      int
	chance        float32
}

func NewAIWatchClosest(e *Entity, lh *LookHelper, maxDist float64) *AIWatchClosest {
	return NewAIWatchClosestWithChance(e, lh, maxDist, 0.02)
}

func NewAIWatchClosestWithChance(e *Entity, lh *LookHelper, maxDist float64, chance float32) *AIWatchClosest {
	g := &AIWatchClosest{
		entity:      e,
		lookHelper:  lh,
		maxDistance: maxDist,
		chance:      chance,
	}
	g.SetMutexBits(2)
	return g
}
func (g *AIWatchClosest) ShouldExecute() bool {
	if rand.Float32() >= g.chance {
		return false
	}

	if g.entity.Level == nil {
		return false
	}

	bb := g.entity.BoundingBox.Grow(g.maxDistance, 3.0, g.maxDistance)
	nearby := g.entity.Level.GetNearbyEntities(bb, g.entity)

	var closest *Entity
	closestDistSq := g.maxDistance * g.maxDistance

	for _, ne := range nearby {
		other, ok := ne.(*Entity)
		if !ok {
			continue
		}
		if other.NetworkID != 0 {
			continue
		}
		dx := other.Position.X - g.entity.Position.X
		dy := other.Position.Y - g.entity.Position.Y
		dz := other.Position.Z - g.entity.Position.Z
		distSq := dx*dx + dy*dy + dz*dz
		if distSq < closestDistSq {
			closestDistSq = distSq
			closest = other
		}
	}

	if closest == nil {
		return false
	}

	g.closestEntity = closest
	return true
}
func (g *AIWatchClosest) ShouldContinueExecuting() bool {
	if g.closestEntity == nil || g.closestEntity.Closed {
		return false
	}

	dx := g.closestEntity.Position.X - g.entity.Position.X
	dy := g.closestEntity.Position.Y - g.entity.Position.Y
	dz := g.closestEntity.Position.Z - g.entity.Position.Z
	distSq := dx*dx + dy*dy + dz*dz

	if distSq > g.maxDistance*g.maxDistance {
		return false
	}

	return g.lookTime > 0
}

func (g *AIWatchClosest) StartExecuting() {
	g.lookTime = 40 + rand.Intn(40)
}

func (g *AIWatchClosest) ResetTask() {
	g.closestEntity = nil
}
func (g *AIWatchClosest) UpdateTask() {
	if g.closestEntity == nil {
		return
	}
	g.lookHelper.SetLookPosition(
		g.closestEntity.Position.X,
		g.closestEntity.Position.Y+g.closestEntity.EyeHeight,
		g.closestEntity.Position.Z,
		10.0,
		math.Abs(float64(g.entity.GetPitch())),
	)
	g.lookTime--
}
