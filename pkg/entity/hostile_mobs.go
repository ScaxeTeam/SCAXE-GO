package entity

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

const ZombieNetworkID = 32
func NewZombie() *Monster {
	m := NewMonster(ZombieNetworkID, "Zombie", 20, 0.6, 1.8, 4)
	m.DropExpMin = 5
	m.DropExpMax = 5
	return m
}
type ZombieDropItem struct {
	ItemID int
	Meta   int
	Count  int
}
func ZombieDrops() []ZombieDropItem {
	const (
		RottenFlesh = 367
		IronIngot   = 265
		Carrot      = 391
		Potato      = 392
	)

	drops := []ZombieDropItem{
		{ItemID: RottenFlesh, Count: 1 + rand.Intn(2)},
	}
	if rand.Intn(100) < 10 {
		switch rand.Intn(3) {
		case 0:
			drops = append(drops, ZombieDropItem{ItemID: IronIngot, Count: 1})
		case 1:
			drops = append(drops, ZombieDropItem{ItemID: Carrot, Count: 1})
		case 2:
			drops = append(drops, ZombieDropItem{ItemID: Potato, Count: 1})
		}
	}

	return drops
}

const SkeletonNetworkID = 34
func NewSkeleton() *Monster {
	m := NewMonster(SkeletonNetworkID, "Skeleton", 20, 0.6, 1.8, 4)
	m.DropExpMin = 5
	m.DropExpMax = 5
	return m
}
func SkeletonDrops() []ZombieDropItem {
	const (
		Bone  = 352
		Arrow = 262
	)

	drops := []ZombieDropItem{
		{ItemID: Bone, Count: 1 + rand.Intn(2)},
	}
	if rand.Intn(2) == 0 {
		drops = append(drops, ZombieDropItem{ItemID: Arrow, Count: 1 + rand.Intn(2)})
	}

	return drops
}

const CreeperNetworkID = 33
type Creeper struct {
	*Monster
	Powered bool
	SwellDirection int
	SwellCounter int
}
func NewCreeper() *Creeper {
	m := NewMonster(CreeperNetworkID, "Creeper", 20, 0.6, 1.8, 0)
	m.DropExpMin = 5
	m.DropExpMax = 5

	return &Creeper{
		Monster: m,
	}
}
func (c *Creeper) SetPowered(powered bool) {
	c.Powered = powered
}
func (c *Creeper) IsPowered() bool {
	return c.Powered
}
func (c *Creeper) SetSwelled(swelled bool) {
	if swelled {
		c.SwellDirection = 1
	} else {
		c.SwellDirection = 0
	}
}
func (c *Creeper) IsSwelled() bool {
	return c.SwellDirection == 1
}
type CreeperTickResult struct {
	HasUpdate     bool
	ShouldExplode bool
}
func (c *Creeper) TickCreeper() CreeperTickResult {
	result := CreeperTickResult{}

	if c.IsSwelled() {
		if c.SwellCounter < 30 {
			c.SwellCounter++
			result.HasUpdate = true
		} else {
			c.SetSwelled(false)
			result.ShouldExplode = true
			result.HasUpdate = true
		}
	} else if c.SwellCounter > 0 {
		c.SwellCounter--
		result.HasUpdate = true
	}

	return result
}
func (c *Creeper) GetExplosionPower() float64 {
	if c.Powered {
		return 6.0
	}
	return 3.0
}
func (c *Creeper) SaveCreeperNBT() {
	c.Monster.SaveMonsterNBT()
	powered := int8(0)
	if c.Powered {
		powered = 1
	}
	c.Monster.Mob.Living.Entity.NamedTag.Set(nbt.NewByteTag("powered", powered))
}
func (c *Creeper) LoadCreeperFromNBT() {
	c.Monster.LoadMonsterFromNBT()
	tag := c.Monster.Mob.Living.Entity.NamedTag
	if tag != nil && tag.GetByte("powered") == 1 {
		c.Powered = true
	}
}
func CreeperDrops() []ZombieDropItem {
	const Gunpowder = 289
	return []ZombieDropItem{
		{ItemID: Gunpowder, Count: 1},
	}
}

const SpiderNetworkID = 35
func NewSpider() *Monster {
	m := NewMonster(SpiderNetworkID, "Spider", 16, 1.4, 0.9, 3)
	m.DropExpMin = 5
	m.DropExpMax = 5
	return m
}
func SpiderDrops() []ZombieDropItem {
	const (
		SpiderEye = 375
		String    = 287
	)
	if rand.Intn(3) < 1 {
		return []ZombieDropItem{{ItemID: SpiderEye, Count: 1}}
	}
	return []ZombieDropItem{{ItemID: String, Count: 1 + rand.Intn(2)}}
}
