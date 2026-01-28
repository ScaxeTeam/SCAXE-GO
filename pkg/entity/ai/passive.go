package ai

import (
	"math"
	"math/rand"
)

type StrollBehavior struct {
	*BehaviorBase
	Duration        int
	TimeLeft        int
	Speed           float64
	SpeedMultiplier float64
}

func NewStrollBehavior(mob MobEntity, duration int, speed, speedMultiplier float64) *StrollBehavior {
	return &StrollBehavior{
		BehaviorBase:    NewBehaviorBase(mob),
		Duration:        duration,
		TimeLeft:        duration,
		Speed:           speed,
		SpeedMultiplier: speedMultiplier,
	}
}

func NewDefaultStrollBehavior(mob MobEntity) *StrollBehavior {
	return NewStrollBehavior(mob, 80, 0.25, 0.75)
}

func (s *StrollBehavior) Name() string {
	return "Stroll"
}

func (s *StrollBehavior) ShouldStart() bool {
	return RandomChance(10)
}

func (s *StrollBehavior) CanContinue() bool {
	s.TimeLeft--
	return s.TimeLeft > 0
}

func (s *StrollBehavior) OnTick() {

	speedFactor := s.Speed * s.SpeedMultiplier * 0.7
	if s.Mob.IsInsideOfWater() {
		speedFactor *= 0.3
	} else {
		speedFactor *= 0.8
	}

	level := s.Mob.GetLevel()
	x, y, z := s.Mob.GetPosition()
	dirX, _, dirZ := s.Mob.GetDirectionVector()

	blockDown := level.GetBlock(int(x), int(y)-1, int(z))
	_, motionY, _ := s.Mob.GetMotion()
	if motionY < 0 && blockDown.IsAir() {
		return
	}

	coordX := x + dirX*speedFactor + dirX*0.5
	coordZ := z + dirZ*speedFactor + dirZ*0.5

	block := level.GetBlock(int(coordX), int(y), int(coordZ))
	blockUp := level.GetBlock(int(coordX), int(y)+1, int(coordZ))
	blockUpUp := level.GetBlock(int(coordX), int(y)+2, int(coordZ))

	height := s.Mob.GetHeight()
	colliding := block.IsSolid() || (height >= 1 && blockUp.IsSolid())

	if !colliding {

		s.Mob.Move(dirX*speedFactor, 0, dirZ*speedFactor)
	} else {

		if !blockUp.IsSolid() && !(height > 1 && blockUpUp.IsSolid()) && !RandomChance(6) {

			s.Mob.SetMotion(0, 0.42, 0)
		} else {

			yaw := s.Mob.GetYaw()
			s.Mob.SetYaw(yaw + 180)
		}
	}

	s.CheckSwimming()
}

func (s *StrollBehavior) OnEnd() {
	s.TimeLeft = s.Duration
	s.Mob.SetMotion(0, 0, 0)
}

type RandomLookaroundBehavior struct {
	*BehaviorBase
	Duration  int
	TimeLeft  int
	TargetYaw float64
}

func NewRandomLookaroundBehavior(mob MobEntity) *RandomLookaroundBehavior {
	return &RandomLookaroundBehavior{
		BehaviorBase: NewBehaviorBase(mob),
		Duration:     40,
		TimeLeft:     0,
	}
}

func (r *RandomLookaroundBehavior) Name() string {
	return "RandomLookaround"
}

func (r *RandomLookaroundBehavior) ShouldStart() bool {
	if RandomChance(20) {
		r.TimeLeft = r.Duration
		r.TargetYaw = float64(randRange(-180, 180))
		return true
	}
	return false
}

func (r *RandomLookaroundBehavior) CanContinue() bool {
	r.TimeLeft--
	return r.TimeLeft > 0
}

func (r *RandomLookaroundBehavior) OnTick() {
	currentYaw := r.Mob.GetYaw()
	diff := r.TargetYaw - currentYaw

	for diff > 180 {
		diff -= 360
	}
	for diff < -180 {
		diff += 360
	}

	if math.Abs(diff) > 5 {
		if diff > 0 {
			r.Mob.SetYaw(currentYaw + 5)
		} else {
			r.Mob.SetYaw(currentYaw - 5)
		}
	}
}

func (r *RandomLookaroundBehavior) OnEnd() {}

type LookAtPlayerBehavior struct {
	*BehaviorBase
	LookDistance float64
	Duration     int
	TimeLeft     int
	Target       PlayerEntity
}

func NewLookAtPlayerBehavior(mob MobEntity, lookDistance float64) *LookAtPlayerBehavior {
	return &LookAtPlayerBehavior{
		BehaviorBase: NewBehaviorBase(mob),
		LookDistance: lookDistance,
		Duration:     60,
		TimeLeft:     0,
	}
}

func NewDefaultLookAtPlayerBehavior(mob MobEntity) *LookAtPlayerBehavior {
	return NewLookAtPlayerBehavior(mob, 6.0)
}

func (l *LookAtPlayerBehavior) Name() string {
	return "LookAtPlayer"
}

func (l *LookAtPlayerBehavior) ShouldStart() bool {
	x, y, z := l.Mob.GetPosition()
	l.Target = l.Mob.GetLevel().GetNearestPlayer(x, y, z, l.LookDistance)
	if l.Target != nil {
		l.TimeLeft = l.Duration
		return true
	}
	return false
}

func (l *LookAtPlayerBehavior) CanContinue() bool {
	l.TimeLeft--
	if l.TimeLeft <= 0 {
		return false
	}
	if l.Target == nil || !l.Target.IsAlive() || !l.Target.IsConnected() {
		return false
	}
	return true
}

func (l *LookAtPlayerBehavior) OnTick() {
	if l.Target != nil {
		tx, ty, tz := l.Target.GetPosition()
		l.AimAt(tx, ty, tz)
	}
	l.CheckSwimming()
}

func (l *LookAtPlayerBehavior) OnEnd() {
	l.Target = nil
	l.Mob.SetPitch(0)
}

type PanicBehavior struct {
	*BehaviorBase
	Speed           float64
	SpeedMultiplier float64
	Duration        int
	TimeLeft        int
	Active          bool
}

func NewPanicBehavior(mob MobEntity, speed, speedMultiplier float64) *PanicBehavior {
	return &PanicBehavior{
		BehaviorBase:    NewBehaviorBase(mob),
		Speed:           speed,
		SpeedMultiplier: speedMultiplier,
		Duration:        60,
		TimeLeft:        0,
		Active:          false,
	}
}

func NewDefaultPanicBehavior(mob MobEntity) *PanicBehavior {
	return NewPanicBehavior(mob, 0.4, 1.0)
}

func (p *PanicBehavior) Name() string {
	return "Panic"
}

func (p *PanicBehavior) Trigger() {
	p.Active = true
	p.TimeLeft = p.Duration

	p.Mob.SetYaw(float64(randRange(-180, 180)))
}

func (p *PanicBehavior) ShouldStart() bool {
	if p.Active {
		p.Active = false
		return true
	}
	return false
}

func (p *PanicBehavior) CanContinue() bool {
	p.TimeLeft--
	return p.TimeLeft > 0
}

func (p *PanicBehavior) OnTick() {
	speedFactor := p.Speed * p.SpeedMultiplier
	if p.Mob.IsInsideOfWater() {
		speedFactor *= 0.5
	}

	dirX, _, dirZ := p.Mob.GetDirectionVector()
	p.Mob.Move(dirX*speedFactor, 0, dirZ*speedFactor)
	p.CheckSwimming()
}

func (p *PanicBehavior) OnEnd() {
	p.Mob.SetMotion(0, 0, 0)
}

func randRange(min, max int) int {
	if max <= min {
		return min
	}
	return min + rand.Intn(max-min)
}
