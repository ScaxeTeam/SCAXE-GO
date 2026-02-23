package entity

import (
	"math"
)

const (
	BoatNetworkID = 90

	BoatWoodOak     = 0
	BoatWoodSpruce  = 1
	BoatWoodBirch   = 2
	BoatWoodJungle  = 3
	BoatWoodAcacia  = 4
	BoatWoodDarkOak = 5

	BoatDespawnAge = 1500

	BoatYawThreshold = 5.0
)

type Boat struct {
	*Entity

	WoodID int

	LinkedEntityID int64

	BoatAge int
}

func NewBoat(woodID int) *Boat {
	b := &Boat{
		Entity:  NewEntity(),
		WoodID:  woodID,
		BoatAge: 0,
	}

	b.Entity.NetworkID = BoatNetworkID
	b.Entity.Width = 1.6
	b.Entity.Height = 0.7
	b.Entity.Gravity = 0.5
	b.Entity.Drag = 0.1
	b.Entity.MaxHealth = 10
	b.Entity.Health = 10

	return b
}

type BoatTickResult struct {
	HasUpdate   bool
	ShouldClose bool
	YawUpdate   bool
	NewYaw      float64
}

func (b *Boat) TickBoat(riderYaw float64, hasRider bool) BoatTickResult {
	result := BoatTickResult{}

	b.Entity.TicksLived++

	if !hasRider {
		b.BoatAge++

		if b.BoatAge > BoatDespawnAge {
			result.ShouldClose = true
			result.HasUpdate = true
			return result
		}
	} else {
		b.BoatAge = 0

		if math.Abs(riderYaw-b.Entity.Yaw) > BoatYawThreshold {
			b.Entity.Yaw = riderYaw
			result.YawUpdate = true
			result.NewYaw = riderYaw
			result.HasUpdate = true
		}
	}

	return result
}

type BoatGravityResult struct {
	MotionY float64
}

func ApplyBoatGravity(isInWater bool, hasBlockBelow bool) BoatGravityResult {
	if hasBlockBelow || isInWater {
		return BoatGravityResult{MotionY: 0.1}
	}
	return BoatGravityResult{MotionY: -0.08}
}

const BoatItemID = 333

func (b *Boat) GetBoatDropItemID() (int, int) {
	return BoatItemID, b.WoodID
}

func (b *Boat) SetRider(entityID int64) {
	b.LinkedEntityID = entityID
	b.BoatAge = 0
}

func (b *Boat) RemoveRider() {
	b.LinkedEntityID = 0
}

func (b *Boat) HasRider() bool {
	return b.LinkedEntityID != 0
}

func (b *Boat) GetWoodID() int {
	return b.WoodID
}
