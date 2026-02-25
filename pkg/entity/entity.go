package entity

import (
	"fmt"
	_ "math"
	"sync/atomic"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/logger"
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
	GetEyeHeight() float64
	HasMovementUpdate() bool
	HasRotationUpdate() bool
	GetBoundingBox() *AxisAlignedBB
}

type ILevel interface {
	GetBlock(x, y, z int32) block.BlockState
	GetCollisionCubes(e IEntity, bb *AxisAlignedBB, includeEntities bool) []*AxisAlignedBB
	GetNearbyEntities(bb *AxisAlignedBB, except IEntity) []IEntity
	GetEntities() []IEntity
	AddEntity(e IEntity)
	RemoveEntity(e IEntity)
	FindGroundY(x, z, startY int32) int32
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

	Gravity       float64
	Drag          float64
	MovementSpeed float64
	SlowFall      bool
	YSize         float64

	Level ILevel

	Tasks      *AITasks
	MoveHelper *MoveHelper
	LookHelper *LookHelper
	JumpHelper *JumpHelper
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
		MovementSpeed: 0.25,
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

func (e *Entity) GetEyeHeight() float64 {
	return e.EyeHeight
}

func (e *Entity) SetMotion(motion *Vector3) {
	e.Motion = motion
}

func (e *Entity) Tick(currentTick int64) bool {
	if e.Closed {
		return false
	}
	e.TicksLived++

	e.LastPos.X = e.Position.X
	e.LastPos.Y = e.Position.Y
	e.LastPos.Z = e.Position.Z
	e.LastYaw = e.Yaw
	e.LastPitch = e.Pitch

	if !e.OnGround && !e.IsInWater() {
		e.Motion.Y -= e.Gravity
		if e.SlowFall && e.Motion.Y < -0.08 {
			e.Motion.Y = -0.08
		}
	}

	e.UpdateAI()

	moved := false
	if e.Motion.X != 0 || e.Motion.Y != 0 || e.Motion.Z != 0 {
		e.Move(e.Motion.X, e.Motion.Y, e.Motion.Z)
		moved = true
	}

	if e.OnGround {
		e.Motion.X *= 0.546
		e.Motion.Z *= 0.546
	}

	if e.NoDamageTicks > 0 {
		e.NoDamageTicks--
	}

	if e.NetworkID > 0 && e.TicksLived%20 == 1 {
		activeGoals := 0
		var goalNames string
		if e.Tasks != nil {
			for _, entry := range e.Tasks.executingEntries {
				if entry.Using {
					activeGoals++
				}
			}
			goalNames = fmt.Sprintf("total=%d/exec=%d", len(e.Tasks.taskEntries), len(e.Tasks.executingEntries))
		}
		mhMoving := false
		mhTarget := ""
		if e.MoveHelper != nil {
			mhMoving = e.MoveHelper.IsMoving
			if mhMoving {
				mhTarget = fmt.Sprintf("%.1f,%.1f,%.1f", e.MoveHelper.PosX, e.MoveHelper.PosY, e.MoveHelper.PosZ)
			}
		}
		dx := e.Position.X - e.LastPos.X
		dy := e.Position.Y - e.LastPos.Y
		dz := e.Position.Z - e.LastPos.Z
		logger.Info("AI-DEBUG",
			"eid", e.ID, "nid", e.NetworkID, "tick", e.TicksLived,
			"pos", fmt.Sprintf("%.2f,%.2f,%.2f", e.Position.X, e.Position.Y, e.Position.Z),
			"delta", fmt.Sprintf("%.4f,%.4f,%.4f", dx, dy, dz),
			"motion", fmt.Sprintf("%.4f,%.4f,%.4f", e.Motion.X, e.Motion.Y, e.Motion.Z),
			"onGround", e.OnGround, "moved", moved,
			"goals", fmt.Sprintf("%d(%s)", activeGoals, goalNames),
			"mhMoving", mhMoving, "mhTarget", mhTarget, "level", e.Level != nil)
	}

	return true
}

func (e *Entity) Close() {
	e.Closed = true
}

func (e *Entity) IsInWater() bool {
	if e.Level == nil {
		return false
	}
	feetBlock := e.Level.GetBlock(int32(e.Position.X), int32(e.Position.Y), int32(e.Position.Z))
	return feetBlock.ID == 8 || feetBlock.ID == 9
}

func (e *Entity) InitAI() {
	e.Tasks = NewAITasks()
	e.MoveHelper = NewMoveHelper(e)
	e.LookHelper = NewLookHelper(e)
	e.JumpHelper = NewJumpHelper(e)
}

func (e *Entity) UpdateAI() {
	if e.Tasks != nil {
		e.Tasks.OnUpdateTasks()
	}
	if e.MoveHelper != nil {
		e.MoveHelper.OnUpdateMoveHelper()
	}
	if e.LookHelper != nil {
		e.LookHelper.OnUpdateLook()
	}
	if e.JumpHelper != nil {
		e.JumpHelper.DoJump()
	}
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
	e.YSize *= 0.4

	moveBB := e.BoundingBox.Clone()

	targetBB := moveBB.AddCoord(dx, dy, dz)
	list := e.Level.GetCollisionCubes(e, targetBB, false)

	for _, bb := range list {
		dy = bb.CalculateYOffset(moveBB, dy)
	}
	moveBB = moveBB.Offset(0, dy, 0)

	fallingFlag := e.OnGround || (dy != wantedY && wantedY < 0)

	for _, bb := range list {
		dx = bb.CalculateXOffset(moveBB, dx)
	}
	moveBB = moveBB.Offset(dx, 0, 0)

	for _, bb := range list {
		dz = bb.CalculateZOffset(moveBB, dz)
	}
	moveBB = moveBB.Offset(0, 0, dz)

	if e.StepHeight > 0 && fallingFlag && (wantedX != dx || wantedZ != dz) {
		cx := dx
		cy := dy
		cz := dz
		dx = wantedX
		dy = e.StepHeight
		dz = wantedZ

		stepBB := e.BoundingBox.Clone()
		stepTargetBB := stepBB.AddCoord(dx, dy, dz)
		stepList := e.Level.GetCollisionCubes(e, stepTargetBB, false)

		for _, bb := range stepList {
			dy = bb.CalculateYOffset(stepBB, dy)
		}
		stepBB = stepBB.Offset(0, dy, 0)

		for _, bb := range stepList {
			dx = bb.CalculateXOffset(stepBB, dx)
		}
		stepBB = stepBB.Offset(dx, 0, 0)

		for _, bb := range stepList {
			dz = bb.CalculateZOffset(stepBB, dz)
		}
		stepBB = stepBB.Offset(0, 0, dz)

		reverseDY := -dy
		for _, bb := range stepList {
			reverseDY = bb.CalculateYOffset(stepBB, reverseDY)
		}
		dy += reverseDY
		stepBB = stepBB.Offset(0, reverseDY, 0)

		if (cx*cx + cz*cz) >= (dx*dx + dz*dz) {
			dx = cx
			dy = cy
			dz = cz
		} else {
			moveBB = stepBB
			e.YSize += 0.5
		}
	}

	e.BoundingBox = moveBB
	e.Position.X = (e.BoundingBox.MinX + e.BoundingBox.MaxX) / 2
	e.Position.Y = e.BoundingBox.MinY - e.YSize
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
