package ai

import (
	"math"
	"math/rand"
)

type FleeEnemyBehavior struct {
	BehaviorBase
	fleeSpeed    float64
	fleeDistance float64
	lastAttacker int64
}

func NewFleeEnemyBehavior(mob MobEntity) *FleeEnemyBehavior {
	return &FleeEnemyBehavior{
		BehaviorBase: BehaviorBase{Mob: mob},
		fleeSpeed:    1.3,
		fleeDistance: 16.0,
	}
}

func (b *FleeEnemyBehavior) Name() string { return "FleeEnemy" }

func (b *FleeEnemyBehavior) ShouldStart() bool { return b.lastAttacker != 0 }

func (b *FleeEnemyBehavior) CanContinue() bool { return b.lastAttacker != 0 }

func (b *FleeEnemyBehavior) OnTick() {

}

func (b *FleeEnemyBehavior) OnEnd() { b.lastAttacker = 0 }

func (b *FleeEnemyBehavior) SetLastAttacker(id int64) { b.lastAttacker = id }

type FollowOwnerBehavior struct {
	BehaviorBase
	ownerID        int64
	followDistance float64
	maxDistance    float64
}

func NewFollowOwnerBehavior(mob MobEntity, ownerID int64) *FollowOwnerBehavior {
	return &FollowOwnerBehavior{
		BehaviorBase:   BehaviorBase{Mob: mob},
		ownerID:        ownerID,
		followDistance: 3.0,
		maxDistance:    10.0,
	}
}

func (b *FollowOwnerBehavior) Name() string { return "FollowOwner" }

func (b *FollowOwnerBehavior) ShouldStart() bool { return true }

func (b *FollowOwnerBehavior) CanContinue() bool { return true }

func (b *FollowOwnerBehavior) OnTick() {}

func (b *FollowOwnerBehavior) OnEnd() {}

type AvoidPlayerBehavior struct {
	BehaviorBase
	avoidDistance float64
	avoidSpeed    float64
}

func NewAvoidPlayerBehavior(mob MobEntity) *AvoidPlayerBehavior {
	return &AvoidPlayerBehavior{
		BehaviorBase:  BehaviorBase{Mob: mob},
		avoidDistance: 8.0,
		avoidSpeed:    1.2,
	}
}

func (b *AvoidPlayerBehavior) Name() string { return "AvoidPlayer" }

func (b *AvoidPlayerBehavior) ShouldStart() bool { return false }

func (b *AvoidPlayerBehavior) CanContinue() bool { return false }

func (b *AvoidPlayerBehavior) OnTick() {}

func (b *AvoidPlayerBehavior) OnEnd() {}

type BreedBehavior struct {
	BehaviorBase
	inLove   bool
	loveTick int
}

func NewBreedBehavior(mob MobEntity) *BreedBehavior {
	return &BreedBehavior{
		BehaviorBase: BehaviorBase{Mob: mob},
		inLove:       false,
		loveTick:     0,
	}
}

func (b *BreedBehavior) Name() string { return "Breed" }

func (b *BreedBehavior) SetInLove(val bool) {
	b.inLove = val
	if val {
		b.loveTick = 600
	}
}

func (b *BreedBehavior) IsInLove() bool {
	return b.inLove && b.loveTick > 0
}

func (b *BreedBehavior) ShouldStart() bool {
	if b.loveTick > 0 {
		b.loveTick--
	}
	return b.IsInLove()
}

func (b *BreedBehavior) CanContinue() bool { return b.IsInLove() }

func (b *BreedBehavior) OnTick() {}

func (b *BreedBehavior) OnEnd() { b.inLove = false }

type TemptBehavior struct {
	BehaviorBase
	temptItems []int
	speed      float64
	targetID   int64
}

func NewTemptBehavior(mob MobEntity, temptItems []int) *TemptBehavior {
	return &TemptBehavior{
		BehaviorBase: BehaviorBase{Mob: mob},
		temptItems:   temptItems,
		speed:        0.8,
	}
}

func (b *TemptBehavior) Name() string { return "Tempt" }

func (b *TemptBehavior) ShouldStart() bool {

	return false
}

func (b *TemptBehavior) CanContinue() bool { return b.targetID != 0 }

func (b *TemptBehavior) OnTick() {
	if b.targetID == 0 || b.Mob == nil {
		return
	}

}

func (b *TemptBehavior) OnEnd() { b.targetID = 0 }

var _ = math.Pi
var _ = rand.Int
