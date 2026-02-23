package entity

import "math/rand"

const CaveSpiderNetworkID = 40

func NewCaveSpider() *Monster {
	m := NewMonster(CaveSpiderNetworkID, "Cave Spider", 12, 0.7, 0.5, 2)
	m.DropExpMin = 5
	m.DropExpMax = 5
	return m
}

func CaveSpiderDrops() []ZombieDropItem {
	const (
		SpiderEye = 375
		String    = 287
	)
	if rand.Intn(3) < 1 {
		return []ZombieDropItem{{ItemID: SpiderEye, Count: 1}}
	}
	return []ZombieDropItem{{ItemID: String, Count: 1 + rand.Intn(2)}}
}

const PigZombieNetworkID = 36

type PigZombie struct {
	*Monster

	Angry bool

	AngerTimer int
}

func NewPigZombie() *PigZombie {
	m := NewMonster(PigZombieNetworkID, "Zombie Pigman", 20, 0.6, 1.8, 5)
	m.DropExpMin = 5
	m.DropExpMax = 5

	return &PigZombie{
		Monster: m,
	}
}

func (p *PigZombie) SetAngry(angry bool) {
	p.Angry = angry
	if angry {
		p.AngerTimer = 400
	} else {
		p.AngerTimer = 0
	}
}

func (p *PigZombie) IsAngry() bool {
	return p.Angry
}

func (p *PigZombie) IsFireImmune() bool {
	return true
}

func (p *PigZombie) TickAnger() {
	if p.Angry && p.AngerTimer > 0 {
		p.AngerTimer--
		if p.AngerTimer <= 0 {
			p.Angry = false
		}
	}
}

func PigZombieDrops() []ZombieDropItem {
	const (
		GoldNugget = 371
		GoldSword  = 283
	)

	drops := []ZombieDropItem{
		{ItemID: GoldNugget, Count: 1 + rand.Intn(3)},
	}

	if rand.Intn(100) < 10 {
		drops = append(drops, ZombieDropItem{ItemID: GoldSword, Count: 1})
	}

	return drops
}

const WitchNetworkID = 45

func NewWitch() *Monster {
	m := NewMonster(WitchNetworkID, "Witch", 26, 0.6, 1.8, 0)
	m.DropExpMin = 5
	m.DropExpMax = 5
	return m
}

func WitchDrops() []ZombieDropItem {
	const (
		GlassBottle   = 374
		GlowstoneDust = 348
		Gunpowder     = 289
		Redstone      = 331
		SpiderEye     = 375
		Sugar         = 353
		Stick         = 280
	)

	possibleDrops := []int{GlassBottle, GlowstoneDust, Gunpowder, Redstone, SpiderEye, Sugar, Stick}
	drops := make([]ZombieDropItem, 0, 3)

	count := 1 + rand.Intn(3)
	for i := 0; i < count && i < len(possibleDrops); i++ {
		idx := rand.Intn(len(possibleDrops))
		drops = append(drops, ZombieDropItem{
			ItemID: possibleDrops[idx],
			Count:  1 + rand.Intn(2),
		})
	}

	return drops
}

const SilverfishNetworkID = 39

func NewSilverfish() *Monster {
	m := NewMonster(SilverfishNetworkID, "Silverfish", 8, 0.4, 0.3, 1)
	m.DropExpMin = 5
	m.DropExpMax = 5
	return m
}

func SilverfishDrops() []ZombieDropItem {
	return nil
}

const BatNetworkID = 19

type Bat struct {
	*Monster

	Hanging bool
}

func NewBat() *Bat {
	m := NewMonster(BatNetworkID, "Bat", 6, 0.5, 0.9, 0)
	m.DropExpMin = 0
	m.DropExpMax = 0

	return &Bat{
		Monster: m,
	}
}

func (b *Bat) SetHanging(hanging bool) {
	b.Hanging = hanging
}

func (b *Bat) IsHanging() bool {
	return b.Hanging
}

func BatDrops() []ZombieDropItem {
	return nil
}
