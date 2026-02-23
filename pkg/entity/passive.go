package entity

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
)

const (
	BabyGrowAge = -24000
)

type Animal struct {
	*Entity

	IsBabyFlag bool

	InLove bool

	AnimalAge int

	DropExpMin int
	DropExpMax int

	FeedFoodID int

	MobName string
}

func NewAnimal(networkID int, name string, maxHealth int, width, height float64) *Animal {
	a := &Animal{
		Entity:    NewEntity(),
		MobName:   name,
		AnimalAge: 0,
	}

	a.Entity.NetworkID = networkID
	a.Entity.Width = width
	a.Entity.Height = height
	a.Entity.MaxHealth = maxHealth
	a.Entity.Health = maxHealth
	a.Entity.Gravity = 0.04
	a.Entity.Drag = 0.02
	a.Entity.StepHeight = 0.6

	return a
}

func (a *Animal) SetBaby(baby bool) {
	a.IsBabyFlag = baby
	if baby {
		a.AnimalAge = BabyGrowAge
	}
}

func (a *Animal) IsBaby() bool {
	return a.IsBabyFlag
}

func (a *Animal) SetInLove(inLove bool) {
	a.InLove = inLove
}

func (a *Animal) IsInLove() bool {
	return a.InLove
}

type AnimalTickResult struct {
	HasUpdate bool
	GrewUp    bool
}

func (a *Animal) TickAnimal() AnimalTickResult {
	result := AnimalTickResult{}

	a.Entity.TicksLived++
	a.AnimalAge++

	if a.IsBabyFlag && a.AnimalAge >= 0 {
		a.SetBaby(false)
		result.GrewUp = true
		result.HasUpdate = true
	}

	return result
}

func (a *Animal) SaveAnimalNBT() {
	a.Entity.SaveNBT()
	baby := int8(0)
	if a.IsBabyFlag {
		baby = 1
	}
	a.Entity.NamedTag.Set(nbt.NewByteTag("IsBaby", baby))
	a.Entity.NamedTag.Set(nbt.NewShortTag("Age", int16(a.AnimalAge)))
}

func (a *Animal) LoadAnimalFromNBT() {
	if a.Entity.NamedTag == nil {
		return
	}
	if a.Entity.NamedTag.GetByte("IsBaby") == 1 {
		a.SetBaby(true)
	}
	age := a.Entity.NamedTag.GetShort("Age")
	if age != 0 {
		a.AnimalAge = int(age)
	}
}

func (a *Animal) GetName() string {
	return a.MobName
}

func (a *Animal) GetFeedFoodID() int {
	return a.FeedFoodID
}

func (a *Animal) CanBeFedWith(itemID int) bool {
	return a.FeedFoodID > 0 && itemID == a.FeedFoodID
}

const CowNetworkID = 11

const ItemWheat = 296

func NewCow() *Animal {
	cow := NewAnimal(CowNetworkID, "Cow", 8, 0.9, 1.3)
	cow.FeedFoodID = ItemWheat
	cow.DropExpMin = 1
	cow.DropExpMax = 3
	return cow
}

func CowDrops(isOnFire bool, rand01 int, count int) (int, int, int) {
	const (
		RawBeef    = 363
		CookedBeef = 364
		Leather    = 334
	)

	if rand01 == 0 {
		if isOnFire {
			return CookedBeef, 0, count
		}
		return RawBeef, 0, count
	}
	return Leather, 0, count
}
