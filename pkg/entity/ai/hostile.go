package ai

import (
	"math"
)

type DamageableEntity interface {
	Attack(damage float64, source interface{}) bool
	IsAlive() bool
}

type AttackEnemyBehavior struct {
	*BehaviorBase
	Speed            float64
	SpeedMultiplier  float64
	LookDistance     float64
	AttackCooldown   int
	TimeLeft         int
	AttackPlayer     bool
	TargetNetworkIDs []int
	Enemy            interface{}
}

func NewAttackEnemyBehavior(mob MobEntity, targetIDs []int, attackPlayer bool, speed, multiplier float64) *AttackEnemyBehavior {
	return &AttackEnemyBehavior{
		BehaviorBase:     NewBehaviorBase(mob),
		Speed:            speed,
		SpeedMultiplier:  multiplier,
		LookDistance:     16.0,
		AttackCooldown:   35,
		AttackPlayer:     attackPlayer,
		TargetNetworkIDs: targetIDs,
	}
}

func NewPlayerAttackBehavior(mob MobEntity, speed, multiplier float64) *AttackEnemyBehavior {
	return NewAttackEnemyBehavior(mob, nil, true, speed, multiplier)
}

func NewDefaultAttackEnemyBehavior(mob MobEntity) *AttackEnemyBehavior {
	return NewPlayerAttackBehavior(mob, 0.25, 0.75)
}

func (a *AttackEnemyBehavior) Name() string {
	return "AttackEnemy"
}

func (a *AttackEnemyBehavior) ShouldStart() bool {
	a.Enemy = nil
	minDistance := a.LookDistance
	mx, my, mz := a.Mob.GetPosition()

	if a.AttackPlayer {

		player := a.Mob.GetLevel().GetNearestPlayer(mx, my, mz, a.LookDistance)
		if player != nil && player.IsSurvival() {
			px, py, pz := player.GetPosition()
			dist := distance3D(mx, my, mz, px, py, pz)
			if dist < minDistance {
				a.Enemy = player
				return true
			}
		}
	} else {

		entities := a.Mob.GetLevel().GetEntities()
		for _, entity := range entities {
			if entity == a.Mob {
				continue
			}
			ex, ey, ez := entity.GetPosition()
			dist := distance3D(mx, my, mz, ex, ey, ez)
			if dist < minDistance {
				a.Enemy = entity
				minDistance = dist
			}
		}
	}

	return a.Enemy != nil
}

func (a *AttackEnemyBehavior) CanContinue() bool {
	if a.Enemy == nil {
		return false
	}

	if player, ok := a.Enemy.(PlayerEntity); ok {
		if !player.IsAlive() || !player.IsConnected() {
			return false
		}
		px, py, pz := player.GetPosition()
		mx, my, mz := a.Mob.GetPosition()
		return distance3D(mx, my, mz, px, py, pz) < a.LookDistance
	}

	if damageable, ok := a.Enemy.(DamageableEntity); ok {
		if !damageable.IsAlive() {
			return false
		}
	}

	return true
}

func (a *AttackEnemyBehavior) OnTick() {
	if a.Enemy == nil {
		return
	}

	var ex, ey, ez float64
	if player, ok := a.Enemy.(PlayerEntity); ok {
		ex, ey, ez = player.GetPosition()
	} else if mob, ok := a.Enemy.(MobEntity); ok {
		ex, ey, ez = mob.GetPosition()
	} else {
		return
	}

	mx, my, mz := a.Mob.GetPosition()
	distance := distance3D(mx, my, mz, ex, ey, ez)

	a.AimAt(ex, ey, ez)

	if distance >= 1.5 {

		a.chaseTarget()
	} else if a.TimeLeft <= 0 {

		a.attackTarget()
		a.TimeLeft = a.AttackCooldown
	}

	if a.TimeLeft > 0 {
		a.TimeLeft--
	}

	a.CheckSwimming()
}

func (a *AttackEnemyBehavior) chaseTarget() {
	speedFactor := a.Speed * a.SpeedMultiplier * 0.7
	if a.Mob.IsInsideOfWater() {
		speedFactor *= 0.3
	} else {
		speedFactor *= 0.8
	}

	level := a.Mob.GetLevel()
	x, y, z := a.Mob.GetPosition()
	dirX, _, dirZ := a.Mob.GetDirectionVector()

	blockDown := level.GetBlock(int(x), int(y)-1, int(z))
	_, motionY, _ := a.Mob.GetMotion()
	if motionY < 0 && blockDown.IsAir() {
		return
	}

	coordX := x + dirX*speedFactor + dirX*0.5
	coordZ := z + dirZ*speedFactor + dirZ*0.5

	block := level.GetBlock(int(coordX), int(y), int(coordZ))
	blockUp := level.GetBlock(int(coordX), int(y)+1, int(coordZ))
	blockUpUp := level.GetBlock(int(coordX), int(y)+2, int(coordZ))

	height := a.Mob.GetHeight()
	colliding := block.IsSolid() || (height >= 1 && blockUp.IsSolid())

	if !colliding {
		a.Mob.Move(dirX*speedFactor, 0, dirZ*speedFactor)
	} else {

		if !blockUp.IsSolid() && !(height > 1 && blockUpUp.IsSolid()) {
			a.Mob.SetMotion(0, 0.42, 0)
		}
	}
}

func (a *AttackEnemyBehavior) attackTarget() {

	damage := 2.0

	if damageable, ok := a.Enemy.(DamageableEntity); ok {
		damageable.Attack(damage, a.Mob)
	}
}

func (a *AttackEnemyBehavior) OnEnd() {
	a.Enemy = nil
	a.Mob.SetMotion(0, 0, 0)
}

type Exploder interface {
	Explode()
}

type ExplodeBehavior struct {
	*BehaviorBase
	FuseTime     int
	FuseLeft     int
	ExplodeRange float64
	Target       PlayerEntity
	Fusing       bool
}

func NewExplodeBehavior(mob MobEntity, fuseTime int, explodeRange float64) *ExplodeBehavior {
	return &ExplodeBehavior{
		BehaviorBase: NewBehaviorBase(mob),
		FuseTime:     fuseTime,
		ExplodeRange: explodeRange,
	}
}

func NewDefaultExplodeBehavior(mob MobEntity) *ExplodeBehavior {
	return NewExplodeBehavior(mob, 30, 3.0)
}

func (e *ExplodeBehavior) Name() string {
	return "Explode"
}

func (e *ExplodeBehavior) ShouldStart() bool {
	x, y, z := e.Mob.GetPosition()
	e.Target = e.Mob.GetLevel().GetNearestPlayer(x, y, z, e.ExplodeRange)
	if e.Target != nil && e.Target.IsSurvival() {
		e.Fusing = true
		e.FuseLeft = e.FuseTime
		return true
	}
	return false
}

func (e *ExplodeBehavior) CanContinue() bool {
	if !e.Fusing {
		return false
	}
	if e.Target == nil || !e.Target.IsAlive() {
		return false
	}

	tx, ty, tz := e.Target.GetPosition()
	mx, my, mz := e.Mob.GetPosition()
	if distance3D(mx, my, mz, tx, ty, tz) > e.ExplodeRange*1.5 {
		return false
	}
	return true
}

func (e *ExplodeBehavior) OnTick() {
	if e.Target != nil {
		tx, ty, tz := e.Target.GetPosition()
		e.AimAt(tx, ty, tz)
	}

	e.FuseLeft--
	if e.FuseLeft <= 0 {

		if exploder, ok := e.Mob.(Exploder); ok {
			exploder.Explode()
		}
		e.Fusing = false
	}
}

func (e *ExplodeBehavior) OnEnd() {
	e.Fusing = false
	e.Target = nil
}

func distance3D(x1, y1, z1, x2, y2, z2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	dz := z2 - z1
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
