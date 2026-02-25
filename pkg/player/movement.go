package player

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/entity"
	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

const (
	MoveBacklogSize = 100
	MovesPerTick = 2
	MaxMoveDistanceSq = 115.0
	EyeHeight = 1.62
)
type MovementState struct {
	LastX, LastY, LastZ float64
	LastYaw, LastPitch  float64
	lastPosInitialized  bool
	MoveRateLimit float64
	SpeedX, SpeedY, SpeedZ float64
	Moving   bool
	OnGround bool
	Swimming bool
	Climbing bool
	FlightPossibility float64
	AllowFlight       bool
	IsCollided bool
}
func newMovementState() *MovementState {
	return &MovementState{
		MoveRateLimit: MoveBacklogSize,
	}
}
func (p *Player) HandleMovePacket(x, y, z float64, yaw, bodyYaw, pitch float32) {
	if !p.Spawned || !p.Connected {
		return
	}
	newY := float64(y) - EyeHeight
	if pitch > 90 || pitch < -90 {
		logger.Warn("Invalid pitch, kicking player",
			"player", p.Username, "pitch", pitch)
		p.Kick("非法移动", false)
		return
	}
	normalizedYaw := float32(math.Mod(float64(yaw), 360))
	if normalizedYaw < 0 {
		normalizedYaw += 360
	}
	p.Yaw = float64(normalizedYaw)
	p.Pitch = float64(pitch)
	p.handleMovement(float64(x), newY, float64(z))
}
func (p *Player) handleMovement(newX, newY, newZ float64) {
	ms := p.movement
	ms.MoveRateLimit--
	if ms.MoveRateLimit < 0 {
		return
	}

	oldX := p.Position.X
	oldY := p.Position.Y
	oldZ := p.Position.Z

	dx := newX - oldX
	dy := newY - oldY
	dz := newZ - oldZ
	distSq := dx*dx + dy*dy + dz*dz

	revert := false
	if distSq > MaxMoveDistanceSq {
		logger.Warn("Player moved too fast, reverting",
			"player", p.Username,
			"distSq", distSq)
		revert = true
	}

	if !revert && distSq > 0.0001 {
		p.mu.Lock()
		p.Position = entity.NewVector3(newX, newY, newZ)
		p.mu.Unlock()
		ms.SpeedX = dx
		ms.SpeedY = dy
		ms.SpeedZ = dz
		ms.Moving = true
		p.checkGroundState(dx, dy, dz)
		p.checkBlockCollision()

	} else if distSq <= 0.0001 {
		ms.SpeedX = 0
		ms.SpeedY = 0
		ms.SpeedZ = 0
		ms.Moving = false
	}

	if revert {
		p.revertMovement(oldX, oldY, oldZ)
	}
}
func (p *Player) processMovement() {
	ms := p.movement
	if ms.MoveRateLimit < MoveBacklogSize {
		ms.MoveRateLimit += MovesPerTick
		if ms.MoveRateLimit > MoveBacklogSize {
			ms.MoveRateLimit = MoveBacklogSize
		}
	}

	curX := p.Position.X
	curY := p.Position.Y
	curZ := p.Position.Z
	if !ms.lastPosInitialized {
		ms.LastX = curX
		ms.LastY = curY
		ms.LastZ = curZ
		ms.LastYaw = p.Yaw
		ms.LastPitch = p.Pitch
		ms.lastPosInitialized = true
		return
	}

	dx := curX - ms.LastX
	dy := curY - ms.LastY
	dz := curZ - ms.LastZ
	distSq := dx*dx + dy*dy + dz*dz
	deltaAngle := math.Abs(ms.LastYaw-p.Yaw) + math.Abs(ms.LastPitch-p.Pitch)

	if distSq > 0.0001 || deltaAngle > 1.0 {
		ms.LastX = curX
		ms.LastY = curY
		ms.LastZ = curZ
		ms.LastYaw = p.Yaw
		ms.LastPitch = p.Pitch
		p.broadcastMovement()
		horizontalDist := math.Sqrt(dx*dx + dz*dz)
		if horizontalDist > 0.01 {
			if p.IsSprinting() {
				p.ExhaustFromSprint(horizontalDist)
			}
		}
	}
	if ms.OnGround {
	}
}
func (p *Player) checkGroundState(dx, dy, dz float64) {
	ms := p.movement

	lvl, ok := p.Human.Level.(*level.Level)
	if !ok || lvl == nil {
		return
	}

	if !ms.OnGround || dx != 0 || dy != 0 || dz != 0 {
		bb := p.BoundingBox
		if bb == nil {
			return
		}
		checkBB := entity.NewAxisAlignedBB(
			bb.MinX, p.Position.Y-0.2, bb.MinZ,
			bb.MaxX, p.Position.Y+0.2, bb.MaxZ,
		)
		if dx < 0 {
			checkBB.MinX += dx
		} else {
			checkBB.MaxX += dx
		}
		if dy < 0 {
			checkBB.MinY += dy
		} else {
			checkBB.MaxY += dy
		}
		if dz < 0 {
			checkBB.MinZ += dz
		} else {
			checkBB.MaxZ += dz
		}

		collisions := lvl.GetCollisionCubes(p, checkBB, false)
		ms.OnGround = len(collisions) > 0
		ms.IsCollided = ms.OnGround
	}
	if !ms.AllowFlight {
		if ms.OnGround || ms.Swimming {
			ms.FlightPossibility = 0
		} else if ms.Climbing {
			ms.FlightPossibility = 0
		} else if dy >= -0.4 {
			ms.FlightPossibility += 1.5
		} else if ms.FlightPossibility > 0 {
			ms.FlightPossibility--
		}
		if ms.FlightPossibility >= 30 {
			if int(ms.FlightPossibility)%10 == 0 {
				p.Teleport(p.Position.X, p.Position.Y-1, p.Position.Z)
			}
			if ms.FlightPossibility >= 50 {
				p.Kick("飞行在此服务器中不被允许", false)
			}
		}
	}

	ms.IsCollided = ms.OnGround
}
func (p *Player) checkBlockCollision() {
	ms := p.movement
	ms.Swimming = false
	ms.Climbing = false

	lvl, ok := p.Human.Level.(*level.Level)
	if !ok || lvl == nil {
		return
	}
	feetX := int32(math.Floor(p.Position.X))
	feetY := int32(math.Floor(p.Position.Y))
	feetZ := int32(math.Floor(p.Position.Z))
	checkPositions := [][3]int32{
		{feetX, feetY, feetZ},
		{feetX, feetY + 1, feetZ},
	}

	for _, pos := range checkPositions {
		bs := lvl.GetBlock(pos[0], pos[1], pos[2])
		switch bs.ID {
		case 8, 9:
			ms.Swimming = true
		case 106, 65:
			ms.Climbing = true
		}
	}
}
func (p *Player) revertMovement(x, y, z float64) {
	ms := p.movement
	ms.LastX = x
	ms.LastY = y
	ms.LastZ = z
	ms.FlightPossibility = 0

	p.mu.Lock()
	p.Position = entity.NewVector3(x, y, z)
	p.mu.Unlock()
	pk := protocol.NewMovePlayerPacket()
	pk.EntityID = p.GetID()
	pk.X = float32(x)
	pk.Y = float32(y) + EyeHeight
	pk.Z = float32(z)
	pk.Yaw = float32(p.Yaw)
	pk.BodyYaw = float32(p.Yaw)
	pk.Pitch = float32(p.Pitch)
	pk.Mode = 1
	pk.OnGround = p.movement.OnGround
	p.SendPacket(pk)
}
func (p *Player) broadcastMovement() {
	viewers := p.getViewers()

	if len(viewers) == 0 {
		return
	}

	pk := protocol.NewMovePlayerPacket()
	pk.EntityID = p.GetID()
	pk.X = float32(p.Position.X)
	pk.Y = float32(p.Position.Y) + EyeHeight
	pk.Z = float32(p.Position.Z)
	pk.Yaw = float32(p.Yaw)
	pk.BodyYaw = float32(p.Yaw)
	pk.Pitch = float32(p.Pitch)
	pk.Mode = 0
	pk.OnGround = p.movement.OnGround

	for _, viewer := range viewers {
		if viewer != p {
			viewer.SendPacket(pk)
		}
	}
}
func (p *Player) IsSprinting() bool {
	return p.Human.Metadata.GetFlag(entity.DataFlags, entity.DataFlagSprinting)
}
func (p *Player) IsSwimming() bool {
	return p.movement.Swimming
}
func (p *Player) IsClimbing() bool {
	return p.movement.Climbing
}
func (p *Player) IsOnGround() bool {
	return p.movement.OnGround
}
func (p *Player) IsMoving() bool {
	return p.movement.Moving
}
func (p *Player) SetAllowFlight(allow bool) {
	p.movement.AllowFlight = allow
	p.movement.FlightPossibility = 0
}
