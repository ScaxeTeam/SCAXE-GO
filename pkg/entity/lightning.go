package entity

import "math/rand"

const LightningNetworkID = 93

const (
	LightningDamageMin = 8
	LightningDamageMax = 20
	LightningFireMin   = 3
	LightningFireMax   = 8
	LightningLifetime  = 20

	LightningRangeX = 4.0
	LightningRangeY = 3.0
	LightningRangeZ = 4.0

	LightningExplodeRadius = 10.0
)

type Lightning struct {
	*Entity

	Age int
}

func NewLightning() *Lightning {
	l := &Lightning{
		Entity: NewEntity(),
		Age:    0,
	}

	l.Entity.NetworkID = LightningNetworkID
	l.Entity.Width = 0.3
	l.Entity.Height = 1.8
	l.Entity.MaxHealth = 2
	l.Entity.Health = 2
	l.Entity.Gravity = 0
	l.Entity.Drag = 0

	return l
}

func (l *Lightning) GetName() string {
	return "Lightning"
}

type LightningTickResult struct {
	ShouldClose bool
}

func (l *Lightning) TickLightning() LightningTickResult {
	l.Age++
	if l.Age > LightningLifetime {
		return LightningTickResult{ShouldClose: true}
	}
	return LightningTickResult{}
}

func CalcLightningDamage() int {
	return LightningDamageMin + rand.Intn(LightningDamageMax-LightningDamageMin+1)
}

func CalcLightningFireDuration() int {
	return LightningFireMin + rand.Intn(LightningFireMax-LightningFireMin+1)
}

type LightningImpactInfo struct {
	ShouldPlaceFire bool

	FireX, FireY, FireZ int

	RangeX, RangeY, RangeZ float64
}

func CalcImpact(lightningFire bool, hitBlockSolid, hitBlockLiquid bool, x, y, z float64) LightningImpactInfo {
	info := LightningImpactInfo{
		RangeX: LightningRangeX,
		RangeY: LightningRangeY,
		RangeZ: LightningRangeZ,
	}

	if !lightningFire {
		return info
	}

	if hitBlockLiquid {
		return info
	}

	if hitBlockSolid {

		info.ShouldPlaceFire = true
		info.FireX = int(x)
		info.FireY = int(y) + 1
		info.FireZ = int(z)
	} else {

		info.ShouldPlaceFire = true
		info.FireX = int(x)
		info.FireY = int(y)
		info.FireZ = int(z)
	}

	return info
}
