package entity

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

const (
	PigNetworkID   = 12
	SheepNetworkID = 13
)

const (
	ItemRawPorkchop    = 319
	ItemCookedPorkchop = 320
	ItemRawChicken     = 365
	ItemCookedChicken  = 366
	ItemFeather        = 288
	ItemCarrot         = 391

	BlockWool = 35
)

func NewPig() *Animal {
	pig := NewAnimal(PigNetworkID, "Pig", 10, 0.9, 0.9)
	pig.FeedFoodID = ItemCarrot
	pig.DropExpMin = 1
	pig.DropExpMax = 3
	return pig
}

func PigDrops(isOnFire bool, count int) (int, int, int) {
	if isOnFire {
		return ItemCookedPorkchop, 0, count
	}
	return ItemRawPorkchop, 0, count
}

type Chicken struct {
	*Animal

	EggTimer int
}

func NewChicken() *Chicken {
	c := &Chicken{
		Animal: NewAnimal(ChickenNetworkID, "Chicken", 4, 0.6, 0.7),
	}
	c.Animal.FeedFoodID = ItemWheat
	c.Animal.DropExpMin = 1
	c.Animal.DropExpMax = 3
	c.Animal.Entity.Gravity = 0.04
	c.ResetEggTimer()
	return c
}

func (c *Chicken) ResetEggTimer() {
	c.EggTimer = 3000 + rand.Intn(3001)
}

type ChickenTickResult struct {
	AnimalTickResult
	ShouldLayEgg bool
}

func (c *Chicken) TickChicken() ChickenTickResult {
	base := c.Animal.TickAnimal()
	result := ChickenTickResult{AnimalTickResult: base}

	c.EggTimer--
	if c.EggTimer <= 0 {
		result.ShouldLayEgg = true
		c.ResetEggTimer()
	}

	return result
}

func (c *Chicken) IsImmuneToFallDamage() bool {
	return true
}

func ChickenDrops(isOnFire bool, rand01 int, count int) (int, int, int) {
	if rand01 == 0 {
		if isOnFire {
			return ItemCookedChicken, 0, count
		}
		return ItemRawChicken, 0, count
	}
	return ItemFeather, 0, count
}

type Sheep struct {
	*Animal

	Color int

	Sheared bool
}

var SheepColorWeights = []struct {
	Color  int
	Weight int
}{
	{0, 20},
	{1, 5},
	{2, 5},
	{5, 5},
	{3, 5},
	{4, 5},
	{6, 5},
	{7, 10},
	{8, 5},
	{9, 5},
	{10, 5},
	{11, 5},
	{12, 5},
	{13, 5},
	{14, 5},
	{15, 5},
}

func GetRandomSheepColor() int {
	totalWeight := 0
	for _, cw := range SheepColorWeights {
		totalWeight += cw.Weight
	}

	r := rand.Intn(totalWeight)
	cumulative := 0
	for _, cw := range SheepColorWeights {
		cumulative += cw.Weight
		if r < cumulative {
			return cw.Color
		}
	}
	return 0
}

func NewSheep() *Sheep {
	s := &Sheep{
		Animal:  NewAnimal(SheepNetworkID, "Sheep", 8, 0.9, 1.3),
		Color:   GetRandomSheepColor(),
		Sheared: false,
	}
	s.Animal.FeedFoodID = ItemWheat
	s.Animal.DropExpMin = 1
	s.Animal.DropExpMax = 3
	return s
}

func NewSheepWithColor(color int) *Sheep {
	s := &Sheep{
		Animal:  NewAnimal(SheepNetworkID, "Sheep", 8, 0.9, 1.3),
		Color:   color,
		Sheared: false,
	}
	s.Animal.FeedFoodID = ItemWheat
	return s
}

func (s *Sheep) SetColor(color int) {
	s.Color = color
}

func (s *Sheep) GetColor() int {
	return s.Color
}

func (s *Sheep) SetSheared(sheared bool) {
	s.Sheared = sheared
}

func (s *Sheep) IsSheared() bool {
	return s.Sheared
}

func (s *Sheep) RegrowWool() {
	s.Sheared = false
}

func (s *Sheep) SheepDrops() (int, int, int) {
	if s.Sheared {
		return 0, 0, 0
	}
	return BlockWool, s.Color, 1
}

func (s *Sheep) ShearDrops(count int) (int, int, int) {
	s.Sheared = true
	return BlockWool, s.Color, count
}

func (s *Sheep) SaveSheepNBT() {
	s.Animal.SaveAnimalNBT()
	s.Animal.Entity.NamedTag.Set(nbt.NewByteTag("Color", int8(s.Color)))
	sheared := int8(0)
	if s.Sheared {
		sheared = 1
	}
	s.Animal.Entity.NamedTag.Set(nbt.NewByteTag("Sheared", sheared))
}

func (s *Sheep) LoadSheepFromNBT() {
	s.Animal.LoadAnimalFromNBT()
	if s.Animal.Entity.NamedTag != nil {
		s.Color = int(s.Animal.Entity.NamedTag.GetByte("Color"))
		s.Sheared = s.Animal.Entity.NamedTag.GetByte("Sheared") == 1
	}
}

var MobFactory = map[int]func() *Animal{
	CowNetworkID: NewCow,
	PigNetworkID: NewPig,
}

func CreatePassiveMob(networkID int) *Animal {
	if factory, ok := MobFactory[networkID]; ok {
		return factory()
	}
	return nil
}
