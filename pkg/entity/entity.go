package entity

import (
	_ "math"
	"sync/atomic"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/nbt"
)

const MotionThreshold = 0.0001

type IEntity interface {
	GetID() int64
	Tick(currentTick int64) bool
	Close()
	GetNetworkID() int
	GetPosition() *Vector3
	GetYaw() float64
	GetPitch() float64
	HasMovementUpdate() bool
	GetBoundingBox() *AxisAlignedBB
}

type ILevel interface {
	GetBlock(x, y, z int32) block.BlockState
	GetCollisionCubes(e IEntity, bb *AxisAlignedBB, includeEntities bool) []*AxisAlignedBB
	GetNearbyEntities(bb *AxisAlignedBB, except IEntity) []IEntity
	GetEntities() []IEntity
	AddEntity(e IEntity)
	RemoveEntity(e IEntity)
}

var entityIDCounter int64 = 1

func NextEntityID() int64 {
	return atomic.AddInt64(&entityIDCounter, 1)
}

type Entity struct {
	ID        int64
	NetworkID int

	Position   *Vector3
	LastPos    *Vector3
	Motion     *Vector3
	LastMotion *Vector3
	Yaw        float64
	Pitch      float64
	LastYaw    float64
	LastPitch  float64
	OnGround   bool

	BoundingBox *AxisAlignedBB
	Width       float64
	Height      float64
	EyeHeight   float64
	StepHeight  float64

	Health        int
	MaxHealth     int
	FireTicks     int
	MaxFireTicks  int
	Age           int
	FallDistance  float64
	TicksLived    int
	NoDamageTicks int
	CanCollide    bool
	Closed        bool
	Invulnerable  bool

	Metadata   *MetadataStore
	Attributes *AttributeMap
	NamedTag   *nbt.CompoundTag

	Gravity float64
	Drag    float64

	Level ILevel
}

func NewEntity() *Entity {
	e := &Entity{
		ID:            NextEntityID(),
		NetworkID:     -1,
		Position:      NewVector3(0, 0, 0),
		LastPos:       NewVector3(0, 0, 0),
		Motion:        NewVector3(0, 0, 0),
		LastMotion:    NewVector3(0, 0, 0),
		Yaw:           0,
		Pitch:         0,
		LastYaw:       0,
		LastPitch:     0,
		OnGround:      false,
		BoundingBox:   NewAxisAlignedBB(0, 0, 0, 0, 0, 0),
		Width:         0.6,
		Height:        1.8,
		EyeHeight:     1.62,
		StepHeight:    0.6,
		Health:        20,
		MaxHealth:     20,
		FireTicks:     0,
		MaxFireTicks:  200,
		FallDistance:  0,
		TicksLived:    0,
		NoDamageTicks: 0,
		CanCollide:    true,
		Invulnerable:  false,
		Metadata:      NewMetadataStore(),
		Attributes:    NewAttributeMap(),
		NamedTag:      nbt.NewCompoundTag(""),
		Gravity:       0.04,
		Drag:          0.02,
	}
	e.recalculateBoundingBox()
	return e
}

func (e *Entity) GetID() int64 {
	return e.ID
}

func (e *Entity) GetNetworkID() int {
	return e.NetworkID
}

func (e *Entity) GetPosition() *Vector3 {
	return e.Position
}

func (e *Entity) GetHealth() int {
	return e.Health
}

func (e *Entity) GetBoundingBox() *AxisAlignedBB {
	return e.BoundingBox
}

func (e *Entity) SetHealth(health int) {
	e.Health = health
	if e.Health < 0 {
		e.Health = 0
	}
	if e.Health > e.MaxHealth {
		e.Health = e.MaxHealth
	}
}

func (e *Entity) GetMaxHealth() int {
	return e.MaxHealth
}

func (e *Entity) SetMaxHealth(maxHealth int) {
	e.MaxHealth = maxHealth
}

func (e *Entity) SetPosition(pos *Vector3) {
	e.Position = pos
	e.recalculateBoundingBox()
}

func (e *Entity) SetRotation(yaw, pitch float64) {
	e.Yaw = yaw
	e.Pitch = pitch
}

func (e *Entity) GetYaw() float64 {
	return e.Yaw
}

func (e *Entity) GetPitch() float64 {
	return e.Pitch
}

func (e *Entity) SetMotion(motion *Vector3) {
	e.Motion = motion
}

func (e *Entity) Tick(currentTick int64) bool {
	if e.Closed {
		return false
	}
	e.TicksLived++
	return true
}

func (e *Entity) Close() {
	e.Closed = true
}

func (e *Entity) recalculateBoundingBox() {
	halfWidth := e.Width / 2
	e.BoundingBox.SetBounds(
		e.Position.X-halfWidth,
		e.Position.Y,
		e.Position.Z-halfWidth,
		e.Position.X+halfWidth,
		e.Position.Y+e.Height,
		e.Position.Z+halfWidth,
	)
}

func (e *Entity) HandleAction(action int32) {

}

func (e *Entity) HandleMove(x, y, z float64, yaw, bodyYaw, pitch float32, onGround bool) {
	e.LastPos.X = e.Position.X
	e.LastPos.Y = e.Position.Y
	e.LastPos.Z = e.Position.Z

	e.Position.X = x
	e.Position.Y = y
	e.Position.Z = z
	e.Yaw = float64(yaw)
	e.Pitch = float64(pitch)
	e.OnGround = onGround

	e.recalculateBoundingBox()
}

func (e *Entity) HasMovementUpdate() bool {
	if e.Motion.X != 0 || e.Motion.Y != 0 || e.Motion.Z != 0 {
		return true
	}
	dx := e.Position.X - e.LastPos.X
	dy := e.Position.Y - e.LastPos.Y
	dz := e.Position.Z - e.LastPos.Z
	distSq := dx*dx + dy*dy + dz*dz
	return distSq > MotionThreshold*MotionThreshold
}

func (e *Entity) HasRotationUpdate() bool {
	return e.Yaw != e.LastYaw || e.Pitch != e.LastPitch
}

func (e *Entity) GetEyePosition() *Vector3 {
	return NewVector3(e.Position.X, e.Position.Y+e.EyeHeight, e.Position.Z)
}

func (e *Entity) GetLocation() *Location {
	return NewLocation(e.Position.X, e.Position.Y, e.Position.Z, float32(e.Yaw), float32(e.Pitch))
}

func (e *Entity) getBlocksAround() []block.BlockState {
	if e.Level == nil {
		return nil
	}

	inset := 0.001
	var blocks []block.BlockState

	minX := int32(e.BoundingBox.MinX + inset)
	minY := int32(e.BoundingBox.MinY + inset)
	minZ := int32(e.BoundingBox.MinZ + inset)
	maxX := int32(e.BoundingBox.MaxX - inset)
	maxY := int32(e.BoundingBox.MaxY - inset)
	maxZ := int32(e.BoundingBox.MaxZ - inset)

	for z := minZ; z <= maxZ; z++ {
		for x := minX; x <= maxX; x++ {
			for y := minY; y <= maxY; y++ {
				blocks = append(blocks, e.Level.GetBlock(x, y, z))
			}
		}
	}
	return blocks
}

func (e *Entity) checkBlockCollision() {

}

func (e *Entity) Move(dx, dy, dz float64) bool {
	if dx == 0 && dy == 0 && dz == 0 {
		return true
	}

	if e.Level == nil {
		return false
	}

	wantedX := dx
	wantedY := dy
	wantedZ := dz

	moveBB := e.BoundingBox.Clone()

	targetBB := moveBB.AddCoord(dx, dy, dz)
	list := e.Level.GetCollisionCubes(e, targetBB, false)

	for _, bb := range list {
		dy = bb.CalculateYOffset(moveBB, dy)
	}
	moveBB = moveBB.Offset(0, dy, 0)

	for _, bb := range list {
		dx = bb.CalculateXOffset(moveBB, dx)
	}
	moveBB = moveBB.Offset(dx, 0, 0)

	for _, bb := range list {
		dz = bb.CalculateZOffset(moveBB, dz)
	}
	moveBB = moveBB.Offset(0, 0, dz)

	e.BoundingBox = moveBB
	e.Position.X = (e.BoundingBox.MinX + e.BoundingBox.MaxX) / 2
	e.Position.Y = e.BoundingBox.MinY
	e.Position.Z = (e.BoundingBox.MinZ + e.BoundingBox.MaxZ) / 2

	e.checkBlockCollision()

	e.OnGround = wantedY != dy && wantedY < 0

	if wantedX != dx {
		e.Motion.X = 0
	}
	if wantedY != dy {
		e.Motion.Y = 0
	}
	if wantedZ != dz {
		e.Motion.Z = 0
	}

	return true
}

func (e *Entity) SaveNBT() {
	if e.NamedTag == nil {
		e.NamedTag = nbt.NewCompoundTag("")
	}

	posList := nbt.NewListTag("Pos", nbt.TagDouble)
	posList.Add(nbt.NewDoubleTag("", e.Position.X))
	posList.Add(nbt.NewDoubleTag("", e.Position.Y))
	posList.Add(nbt.NewDoubleTag("", e.Position.Z))
	e.NamedTag.Set(posList)

	motList := nbt.NewListTag("Motion", nbt.TagDouble)
	motList.Add(nbt.NewDoubleTag("", e.Motion.X))
	motList.Add(nbt.NewDoubleTag("", e.Motion.Y))
	motList.Add(nbt.NewDoubleTag("", e.Motion.Z))
	e.NamedTag.Set(motList)

	rotList := nbt.NewListTag("Rotation", nbt.TagFloat)
	rotList.Add(nbt.NewFloatTag("", float32(e.Yaw)))
	rotList.Add(nbt.NewFloatTag("", float32(e.Pitch)))
	e.NamedTag.Set(rotList)

	e.NamedTag.Set(nbt.NewShortTag("Health", int16(e.Health)))
	e.NamedTag.Set(nbt.NewShortTag("Fire", int16(e.FireTicks)))
	e.NamedTag.Set(nbt.NewShortTag("Air", 300))

	val := int8(0)
	if e.OnGround {
		val = 1
	}
	e.NamedTag.Set(nbt.NewByteTag("OnGround", val))
}

func (e *Entity) SetSprinting(value bool) {
	e.Metadata.SetFlag(DataFlags, DataFlagSprinting, value)

}

func (e *Entity) SetSneaking(value bool) {
	e.Metadata.SetFlag(DataFlags, DataFlagSneaking, value)
}
