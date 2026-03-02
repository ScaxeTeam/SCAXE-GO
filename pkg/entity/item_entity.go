package entity

import (
	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/nbt"
)

type ItemEntity struct {
	*Entity
	Item        item.Item
	PickupDelay int
	Owner       string
	Thrower     string
	Age         int
	Settled     bool
}

func NewItemEntity(it item.Item) *ItemEntity {
	e := &ItemEntity{
		Entity:      NewEntity(),
		Item:        it,
		PickupDelay: 10,
		Age:         0,
	}

	e.Width = 0.25
	e.Height = 0.25
	e.EyeHeight = 0
	e.StepHeight = 0
	e.recalculateBoundingBox()

	e.Gravity = 0.04
	e.Drag = 0.02

	return e
}

func (e *ItemEntity) Tick(currentTick int64) bool {
	if e.Closed {
		return false
	}

	e.LastPos.X = e.Position.X
	e.LastPos.Y = e.Position.Y
	e.LastPos.Z = e.Position.Z

	e.Age++
	if e.Age > 6000 {
		e.Close()
		return false
	}

	if e.PickupDelay > 0 {
		e.PickupDelay--

		if e.PickupDelay > 0 && e.PickupDelay%8 == 0 {
			e.tryMerge()
		}
	}

	e.Motion.Y -= e.Gravity

	e.Move(e.Motion.X, e.Motion.Y, e.Motion.Z)

	friction := 1.0 - e.Drag

	if e.OnGround {
		friction *= 0.6
	}

	e.Motion.X *= friction
	e.Motion.Y *= (1.0 - e.Drag)
	e.Motion.Z *= friction

	if e.OnGround {
		e.Motion.Y *= -0.5
	}

	motionSq := e.Motion.X*e.Motion.X + e.Motion.Y*e.Motion.Y + e.Motion.Z*e.Motion.Z
	if e.OnGround && motionSq < 0.0001 {
		if !e.Settled {
			e.Settled = true

		}
	} else {
		e.Settled = false
	}

	return true
}

func (e *ItemEntity) HasMovementUpdate() bool {
	if e.Settled {

		dx := e.Position.X - e.LastPos.X
		dy := e.Position.Y - e.LastPos.Y
		dz := e.Position.Z - e.LastPos.Z
		return dx*dx+dy*dy+dz*dz > 0.0001
	}

	return false
}

func (e *ItemEntity) CanCollideWith(ent *Entity) bool {
	return false
}

func (e *ItemEntity) SaveNBT() {
	e.Entity.SaveNBT()
	e.NamedTag.Set(e.Item.NBTSerialize(-1))
	e.NamedTag.Set(nbt.NewShortTag("Health", int16(e.Health)))
	e.NamedTag.Set(nbt.NewShortTag("Age", int16(e.Age)))
	e.NamedTag.Set(nbt.NewShortTag("PickupDelay", int16(e.PickupDelay)))
	if e.Owner != "" {
		e.NamedTag.Set(nbt.NewStringTag("Owner", e.Owner))
	}
	if e.Thrower != "" {
		e.NamedTag.Set(nbt.NewStringTag("Thrower", e.Thrower))
	}
}

func (e *ItemEntity) tryMerge() {
	if e.Level == nil {
		return
	}

	bb := e.BoundingBox
	if bb == nil {
		return
	}
	searchBB := bb.Grow(0.5, 0.5, 0.5)

	entities := e.Level.GetNearbyEntities(searchBB, e)

	for _, ent := range entities {
		target, ok := ent.(*ItemEntity)
		if !ok || target.Closed {
			continue
		}

		if e.isMergeable(target) {

			newCount := target.Item.Count + e.Item.Count
			target.Item.Count = newCount

			if e.PickupDelay > target.PickupDelay {
				target.PickupDelay = e.PickupDelay
			}

			e.Close()
			return
		}
	}
}

func (e *ItemEntity) isMergeable(target *ItemEntity) bool {
	if target == e {
		return false
	}

	if target.PickupDelay == 32767 {
		return false
	}

	if !e.Item.Equals(target.Item, true, true) {
		return false
	}

	if target.Item.Count+e.Item.Count > target.Item.GetMaxStackSize() {
		return false
	}
	return true
}
