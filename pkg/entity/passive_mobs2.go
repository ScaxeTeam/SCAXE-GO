package entity

import "math/rand"

const IronGolemNetworkID = 20
type IronGolem struct {
	*Monster
	PlayerCreated bool
}
func NewIronGolem() *IronGolem {
	m := NewMonster(IronGolemNetworkID, "Iron Golem", 100, 1.4, 2.7, 21)
	m.DropExpMin = 0
	m.DropExpMax = 0

	return &IronGolem{
		Monster: m,
	}
}
func (g *IronGolem) IsPlayerCreated() bool {
	return g.PlayerCreated
}
func IronGolemDrops() []ZombieDropItem {
	const (
		IronIngot = 265
		Poppy     = 38
	)
	drops := []ZombieDropItem{
		{ItemID: IronIngot, Count: 3 + rand.Intn(3)},
	}
	poppyCount := rand.Intn(3)
	if poppyCount > 0 {
		drops = append(drops, ZombieDropItem{ItemID: Poppy, Count: poppyCount})
	}
	return drops
}

const SnowGolemNetworkID = 21
type SnowGolem struct {
	*Monster
	Sheared bool
}
func NewSnowGolem() *SnowGolem {
	m := NewMonster(SnowGolemNetworkID, "Snow Golem", 4, 0.7, 1.9, 0)
	m.DropExpMin = 0
	m.DropExpMax = 0

	return &SnowGolem{
		Monster: m,
	}
}
func (s *SnowGolem) SetSheared(sheared bool) {
	s.Sheared = sheared
}
func (s *SnowGolem) IsSheared() bool {
	return s.Sheared
}
func SnowGolemDrops() []ZombieDropItem {
	const Snowball = 332
	count := rand.Intn(16)
	if count == 0 {
		return nil
	}
	return []ZombieDropItem{{ItemID: Snowball, Count: count}}
}

const SquidNetworkID = 17
func NewSquid() *Monster {
	m := NewMonster(SquidNetworkID, "Squid", 10, 0.8, 0.8, 0)
	m.DropExpMin = 1
	m.DropExpMax = 3
	return m
}
func SquidDrops() []ZombieDropItem {
	const InkSac = 351
	return []ZombieDropItem{{ItemID: InkSac, Meta: 0, Count: 1 + rand.Intn(3)}}
}

const NPCNetworkID = 15
type NPC struct {
	*Monster
	Profession int
}

const (
	ProfessionFarmer     = 0
	ProfessionLibrarian  = 1
	ProfessionPriest     = 2
	ProfessionBlacksmith = 3
	ProfessionButcher    = 4
)
func NewNPC() *NPC {
	m := NewMonster(NPCNetworkID, "Villager", 20, 0.6, 1.8, 0)
	m.DropExpMin = 0
	m.DropExpMax = 0

	return &NPC{
		Monster:    m,
		Profession: rand.Intn(5),
	}
}
func NewNPCWithProfession(profession int) *NPC {
	m := NewMonster(NPCNetworkID, "Villager", 20, 0.6, 1.8, 0)
	m.DropExpMin = 0
	m.DropExpMax = 0

	return &NPC{
		Monster:    m,
		Profession: profession,
	}
}
func (n *NPC) GetProfession() int {
	return n.Profession
}
func (n *NPC) SetProfession(profession int) {
	n.Profession = profession
}
func NPCDrops() []ZombieDropItem {
	return nil
}
type FlyingAnimal struct {
	*Animal
	Flying bool
}
func NewFlyingAnimal(networkID int, name string, maxHealth int, width, height float64) *FlyingAnimal {
	a := NewAnimal(networkID, name, maxHealth, width, height, 0.25, 1.25)
	a.Entity.Gravity = 0.02

	return &FlyingAnimal{
		Animal: a,
	}
}
func (f *FlyingAnimal) SetFlying(flying bool) {
	f.Flying = flying
}
func (f *FlyingAnimal) IsFlying() bool {
	return f.Flying
}
