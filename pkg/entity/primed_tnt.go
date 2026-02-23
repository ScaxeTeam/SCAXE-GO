package entity

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
)

const (
	PrimedTNTNetworkID = 65

	DefaultFuse = 80

	ShortFuseMin = 10
	ShortFuseMax = 30

	DefaultExplosionForce = 4.0
)

type PrimedTNT struct {
	*Entity

	Fuse int

	ExplosionForce float64

	BlockBreaking bool
}

func NewPrimedTNT(fuse int) *PrimedTNT {
	t := &PrimedTNT{
		Entity:         NewEntity(),
		Fuse:           fuse,
		ExplosionForce: DefaultExplosionForce,
		BlockBreaking:  true,
	}

	t.Entity.NetworkID = PrimedTNTNetworkID
	t.Entity.Width = 0.98
	t.Entity.Height = 0.98
	t.Entity.Gravity = 0.04
	t.Entity.Drag = 0.02
	t.Entity.Health = 1
	t.Entity.MaxHealth = 1
	t.Entity.CanCollide = false

	return t
}

type PrimedTNTTickResult struct {
	HasUpdate     bool
	ShouldExplode bool
	Force         float64
	BlockBreaking bool
	ExplodeX      float64
	ExplodeY      float64
	ExplodeZ      float64
}

func (t *PrimedTNT) TickTNT() PrimedTNTTickResult {
	result := PrimedTNTTickResult{}

	t.Entity.TicksLived++
	t.Fuse--

	result.HasUpdate = true

	if t.Fuse <= 0 {
		result.ShouldExplode = true
		result.Force = t.ExplosionForce
		result.BlockBreaking = t.BlockBreaking
		result.ExplodeX = t.Entity.Position.X
		result.ExplodeY = t.Entity.Position.Y + t.Entity.Height/2
		result.ExplodeZ = t.Entity.Position.Z
	}

	return result
}

func (t *PrimedTNT) SavePrimedTNTNBT() {
	t.Entity.SaveNBT()
	t.Entity.NamedTag.Set(nbt.NewByteTag("Fuse", int8(t.Fuse)))
}

func (t *PrimedTNT) LoadFuseFromNBT() {
	if t.Entity.NamedTag != nil {
		fuse := t.Entity.NamedTag.GetByte("Fuse")
		if fuse > 0 {
			t.Fuse = int(fuse)
		}
	}
}

func (t *PrimedTNT) GetFuse() int {
	return t.Fuse
}

func (t *PrimedTNT) SetFuse(fuse int) {
	t.Fuse = fuse
}

func (t *PrimedTNT) SetExplosionForce(force float64) {
	t.ExplosionForce = force
}

func (t *PrimedTNT) SetBlockBreaking(breaking bool) {
	t.BlockBreaking = breaking
}
