package entity

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

const EndermanNetworkID = 38
type Enderman struct {
	*Monster
	CarriedBlockID int
	CarriedBlockMeta int
	Trembling bool
}
func NewEnderman() *Enderman {
	m := NewMonster(EndermanNetworkID, "Enderman", 40, 0.6, 2.9, 7)
	m.DropExpMin = 5
	m.DropExpMax = 5

	return &Enderman{
		Monster: m,
	}
}
func (e *Enderman) SetBlockInHand(blockID, blockMeta int) {
	e.CarriedBlockID = blockID
	e.CarriedBlockMeta = blockMeta
}
func (e *Enderman) GetBlockInHand() (blockID, blockMeta int) {
	return e.CarriedBlockID, e.CarriedBlockMeta
}
func (e *Enderman) IsCarryingBlock() bool {
	return e.CarriedBlockID > 0
}
func (e *Enderman) SetTremble(trembling bool) {
	e.Trembling = trembling
}
func (e *Enderman) SaveEndermanNBT() {
	e.Monster.SaveMonsterNBT()
	tag := e.Monster.Mob.Living.Entity.NamedTag
	tag.Set(nbt.NewShortTag("carried", int16(e.CarriedBlockID)))
	tag.Set(nbt.NewShortTag("carriedData", int16(e.CarriedBlockMeta)))
}
func (e *Enderman) LoadEndermanFromNBT() {
	e.Monster.LoadMonsterFromNBT()
	tag := e.Monster.Mob.Living.Entity.NamedTag
	if tag != nil {
		e.CarriedBlockID = int(tag.GetShort("carried"))
		e.CarriedBlockMeta = int(tag.GetShort("carriedData"))
	}
}
func EndermanDrops() []ZombieDropItem {
	const EnderPearl = 368
	count := rand.Intn(2)
	if count == 0 {
		return nil
	}
	return []ZombieDropItem{{ItemID: EnderPearl, Count: count}}
}

const BlazeNetworkID = 43
type Blaze struct {
	*Monster
	Charging bool
}
func NewBlaze() *Blaze {
	m := NewMonster(BlazeNetworkID, "Blaze", 20, 0.6, 1.8, 6)
	m.DropExpMin = 10
	m.DropExpMax = 10

	return &Blaze{
		Monster: m,
	}
}
func (b *Blaze) IsFireImmune() bool {
	return true
}
func (b *Blaze) SetCharging(charging bool) {
	b.Charging = charging
}
func (b *Blaze) IsCharging() bool {
	return b.Charging
}
func BlazeDrops() []ZombieDropItem {
	const BlazeRod = 369
	count := rand.Intn(2)
	if count == 0 {
		return nil
	}
	return []ZombieDropItem{{ItemID: BlazeRod, Count: count}}
}

const GhastNetworkID = 41
type Ghast struct {
	*Monster
	Charging bool
}
func NewGhast() *Ghast {
	m := NewMonster(GhastNetworkID, "Ghast", 10, 4.0, 4.0, 0)
	m.DropExpMin = 5
	m.DropExpMax = 5

	return &Ghast{
		Monster: m,
	}
}
func (g *Ghast) SetCharging(charging bool) {
	g.Charging = charging
}
func (g *Ghast) IsCharging() bool {
	return g.Charging
}
func GhastDrops() []ZombieDropItem {
	const (
		GhastTear = 370
		Gunpowder = 289
	)
	drops := []ZombieDropItem{
		{ItemID: Gunpowder, Count: rand.Intn(2) + 1},
	}
	if rand.Intn(2) == 0 {
		drops = append(drops, ZombieDropItem{ItemID: GhastTear, Count: 1})
	}
	return drops
}

const SlimeNetworkID = 37
type Slime struct {
	*Monster
	Size int
}
func NewSlime() *Slime {
	size := 1 + rand.Intn(4)
	m := NewMonster(SlimeNetworkID, "Slime", slimeHealthForSize(size), 0.6, 0.6, slimeDamageForSize(size))
	m.DropExpMin = 1
	m.DropExpMax = 4

	return &Slime{
		Monster: m,
		Size:    size,
	}
}
func NewSlimeWithSize(size int) *Slime {
	if size < 1 {
		size = 1
	}
	if size > 4 {
		size = 4
	}
	m := NewMonster(SlimeNetworkID, "Slime", slimeHealthForSize(size), 0.6, 0.6, slimeDamageForSize(size))
	m.DropExpMin = 1
	m.DropExpMax = 4

	return &Slime{
		Monster: m,
		Size:    size,
	}
}
func slimeHealthForSize(size int) int {
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
func slimeDamageForSize(size int) int {
	switch size {
	case 1:
		return 0
	case 2:
		return 2
	case 3:
		return 3
	case 4:
		return 4
	default:
		return 0
	}
}
func (s *Slime) GetSize() int {
	return s.Size
}
func (s *Slime) SetSize(size int) {
	s.Size = size
}
func (s *Slime) ShouldSplit() bool {
	return s.Size > 1
}
func (s *Slime) GetSplitSize() int {
	return s.Size - 1
}
func (s *Slime) GetSplitCount() int {
	return 4
}
func SlimeDrops(size int) []ZombieDropItem {
	const Slimeball = 341
	if size == 1 {
		return []ZombieDropItem{{ItemID: Slimeball, Count: 1}}
	}
	return nil
}
