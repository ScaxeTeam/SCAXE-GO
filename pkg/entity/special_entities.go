package entity

import "math/rand"

const FishingHookNetworkID = 77

type FishingHook struct {
	Entity

	OwnerID int64

	WaitTimer int

	Hooked bool
}

func NewFishingHook(ownerID int64) *FishingHook {
	fh := &FishingHook{
		OwnerID:   ownerID,
		WaitTimer: 100 + rand.Intn(400),
	}
	fh.Entity.NetworkID = FishingHookNetworkID
	fh.Entity.Width = 0.25
	fh.Entity.Height = 0.25
	fh.Entity.Gravity = 0.04
	fh.Entity.MaxHealth = 1
	fh.Entity.Health = 1
	return fh
}

func (fh *FishingHook) TickFishingHook() bool {
	if fh.WaitTimer > 0 {
		fh.WaitTimer--
		if fh.WaitTimer <= 0 {
			fh.Hooked = true
			return true
		}
	}
	return false
}

func (fh *FishingHook) IsHooked() bool {
	return fh.Hooked
}

const ExperienceOrbNetworkID = 69

type ExperienceOrb struct {
	Entity

	Experience int

	PickupDelay int

	Age int
}

const (
	ExperienceOrbMaxAge     = 6000
	ExperienceOrbPickupDist = 2.0
)

func NewExperienceOrb(experience int) *ExperienceOrb {
	orb := &ExperienceOrb{
		Experience:  experience,
		PickupDelay: 10,
	}
	orb.Entity.NetworkID = ExperienceOrbNetworkID
	orb.Entity.Width = 0.25
	orb.Entity.Height = 0.25
	orb.Entity.Gravity = 0.04
	orb.Entity.MaxHealth = 1
	orb.Entity.Health = 1
	return orb
}

func (o *ExperienceOrb) TickExperienceOrb() bool {
	o.Age++
	if o.Age >= ExperienceOrbMaxAge {
		return true
	}
	if o.PickupDelay > 0 {
		o.PickupDelay--
	}
	return false
}

func (o *ExperienceOrb) CanPickup() bool {
	return o.PickupDelay <= 0
}

func (o *ExperienceOrb) GetExperience() int {
	return o.Experience
}

const BigFireballNetworkID = 85

type BigFireball struct {
	Entity

	ShooterID int64

	ExplosionPower float64
}

func NewBigFireball(shooterID int64) *BigFireball {
	fb := &BigFireball{
		ShooterID:      shooterID,
		ExplosionPower: 1.0,
	}
	fb.Entity.NetworkID = BigFireballNetworkID
	fb.Entity.Width = 1.0
	fb.Entity.Height = 1.0
	fb.Entity.Gravity = 0.0
	fb.Entity.MaxHealth = 1
	fb.Entity.Health = 1
	return fb
}

const SmallFireballNetworkID = 94

type SmallFireball struct {
	Entity

	ShooterID int64
}

func NewSmallFireball(shooterID int64) *SmallFireball {
	fb := &SmallFireball{
		ShooterID: shooterID,
	}
	fb.Entity.NetworkID = SmallFireballNetworkID
	fb.Entity.Width = 0.3125
	fb.Entity.Height = 0.3125
	fb.Entity.Gravity = 0.0
	fb.Entity.MaxHealth = 1
	fb.Entity.Health = 1
	return fb
}

const LavaSlimeNetworkID = 42

type LavaSlime struct {
	*Monster

	Size int
}

func NewLavaSlime() *LavaSlime {
	size := 1 + rand.Intn(4)
	m := NewMonster(LavaSlimeNetworkID, "Magma Cube", lavaSlimeHealthForSize(size),
		0.6, 0.6, lavaSlimeDamageForSize(size))
	m.DropExpMin = 1
	m.DropExpMax = 4

	return &LavaSlime{
		Monster: m,
		Size:    size,
	}
}

func NewLavaSlimeWithSize(size int) *LavaSlime {
	if size < 1 {
		size = 1
	}
	if size > 4 {
		size = 4
	}
	m := NewMonster(LavaSlimeNetworkID, "Magma Cube", lavaSlimeHealthForSize(size),
		0.6, 0.6, lavaSlimeDamageForSize(size))
	m.DropExpMin = 1
	m.DropExpMax = 4

	return &LavaSlime{
		Monster: m,
		Size:    size,
	}
}

func lavaSlimeHealthForSize(size int) int {
	switch size {
	case 1:
		return 1
	case 2:
		return 4
	case 3:
		return 8
	case 4:
		return 16
	default:
		return 1
	}
}

func lavaSlimeDamageForSize(size int) int {
	switch size {
	case 1:
		return 3
	case 2:
		return 4
	case 3:
		return 5
	case 4:
		return 6
	default:
		return 3
	}
}

func (l *LavaSlime) GetSize() int {
	return l.Size
}

func (l *LavaSlime) SetSize(size int) {
	l.Size = size
}

func (l *LavaSlime) ShouldSplit() bool {
	return l.Size > 1
}

func (l *LavaSlime) GetSplitSize() int {
	return l.Size - 1
}

func (l *LavaSlime) IsFireImmune() bool {
	return true
}

func LavaSlimeDrops(size int) []ZombieDropItem {
	const MagmaCream = 378
	if size == 1 {
		return []ZombieDropItem{{ItemID: MagmaCream, Count: 1}}
	}
	return nil
}
