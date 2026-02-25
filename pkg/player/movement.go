package player

// movement.go — 玩家移动处理 + 碰撞检测
// 对应 PHP Player.php 的 handleMovement/processMovement/checkGroundState/checkBlockCollision
//
// 移动流程：
//   1. 客户端发 MovePlayerPacket → HandleMovePacket (pitch 验证)
//   2. handleMovement (速率限制 + 距离检查 + 调用 Level.GetCollisionCubes)
//   3. processMovement (每 tick, PlayerMoveEvent, 实体移动广播, 饥饿消耗)
//   4. checkGroundState (地面碰撞检测, 反飞行)
//   5. checkBlockCollision (游泳/攀爬状态, 方块实体碰撞)

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/entity"
	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

// ── 移动常量 ──────────────────────────────────────

const (
	// MoveBacklogSize 最大积压移动包数量
	MoveBacklogSize = 100

	// MovesPerTick 每 tick 允许处理的移动包数
	MovesPerTick = 2

	// MaxMoveDistanceSq 单次移动最大距离平方（反瞬移）
	// 对应 PHP: distanceSquared > 115
	MaxMoveDistanceSq = 115.0

	// EyeHeight 玩家眼睛高度偏移
	EyeHeight = 1.62
)

// ── 移动状态 ──────────────────────────────────────

// MovementState 存储玩家移动相关的状态
type MovementState struct {
	// 上次确认的位置/角度
	LastX, LastY, LastZ float64
	LastYaw, LastPitch  float64
	lastPosInitialized  bool

	// 移动速率限制
	MoveRateLimit float64

	// 速度向量
	SpeedX, SpeedY, SpeedZ float64

	// 状态标记
	Moving   bool
	OnGround bool
	Swimming bool
	Climbing bool

	// 反飞行
	FlightPossibility float64
	AllowFlight       bool

	// 碰撞
	IsCollided bool
}

// newMovementState 创建移动状态
func newMovementState() *MovementState {
	return &MovementState{
		MoveRateLimit: MoveBacklogSize,
	}
}

// ── MovePlayerPacket 处理 ──────────────────────────

// HandleMovePacket 处理客户端发来的 MovePlayerPacket
// 对应 PHP Player::handleDataPacket() case MOVE_PLAYER_PACKET (L2508-2538)
func (p *Player) HandleMovePacket(x, y, z float64, yaw, bodyYaw, pitch float32) {
	if !p.Spawned || !p.Connected {
		return
	}

	// 修正 Y 坐标: 客户端发送的是眼睛高度，需要减去 EyeHeight
	newY := float64(y) - EyeHeight

	// Pitch 范围验证（-90 ~ 90，超出为非法移动）
	if pitch > 90 || pitch < -90 {
		logger.Warn("Invalid pitch, kicking player",
			"player", p.Username, "pitch", pitch)
		p.Kick("非法移动", false)
		return
	}

	// 标准化 Yaw (0-360)
	normalizedYaw := float32(math.Mod(float64(yaw), 360))
	if normalizedYaw < 0 {
		normalizedYaw += 360
	}

	// 设置旋转
	p.Yaw = float64(normalizedYaw)
	p.Pitch = float64(pitch)

	// 处理移动
	p.handleMovement(float64(x), newY, float64(z))
}

// ── 移动处理核心 ──────────────────────────────────

// handleMovement 核心移动处理
// 对应 PHP Player::handleMovement() L1721-1771
func (p *Player) handleMovement(newX, newY, newZ float64) {
	ms := p.movement

	// 速率限制
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

	// 防瞬移: 单次移动距离过大
	if distSq > MaxMoveDistanceSq {
		logger.Warn("Player moved too fast, reverting",
			"player", p.Username,
			"distSq", distSq)
		revert = true
	}

	if !revert && distSq > 0.0001 {
		// 更新位置
		p.mu.Lock()
		p.Position = entity.NewVector3(newX, newY, newZ)
		p.mu.Unlock()

		// 记录速度
		ms.SpeedX = dx
		ms.SpeedY = dy
		ms.SpeedZ = dz
		ms.Moving = true

		// 地面状态检测
		p.checkGroundState(dx, dy, dz)

		// 方块碰撞检测
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

// processMovement 每 tick 移动处理（广播 + 事件）
// 对应 PHP Player::processMovement() L1776-1849
// 应在 Player.Tick() 中调用
func (p *Player) processMovement() {
	ms := p.movement

	// 恢复速率限制额度
	if ms.MoveRateLimit < MoveBacklogSize {
		ms.MoveRateLimit += MovesPerTick
		if ms.MoveRateLimit > MoveBacklogSize {
			ms.MoveRateLimit = MoveBacklogSize
		}
	}

	curX := p.Position.X
	curY := p.Position.Y
	curZ := p.Position.Z

	// 计算与上次确认位置的距离
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

		// 广播实体移动给其他玩家
		p.broadcastMovement()

		// 计算饥饿消耗
		horizontalDist := math.Sqrt(dx*dx + dz*dz)
		if horizontalDist > 0.01 {
			if p.IsSprinting() {
				p.ExhaustFromSprint(horizontalDist)
			}
		}
	}

	// 地面 tick
	if ms.OnGround {
		// inAirTicks = 0
	}
}

// ── 碰撞检测 ──────────────────────────────────────

// checkGroundState 检查玩家是否在地面上
// 对应 PHP Player::checkGroundState() L1584-1630
func (p *Player) checkGroundState(dx, dy, dz float64) {
	ms := p.movement

	lvl, ok := p.Human.Level.(*level.Level)
	if !ok || lvl == nil {
		return
	}

	if !ms.OnGround || dx != 0 || dy != 0 || dz != 0 {
		// 构造脚下检测区域
		bb := p.BoundingBox
		if bb == nil {
			return
		}

		// 在脚下 0.2 格范围内检测碰撞
		checkBB := entity.NewAxisAlignedBB(
			bb.MinX, p.Position.Y-0.2, bb.MinZ,
			bb.MaxX, p.Position.Y+0.2, bb.MaxZ,
		)

		// 扩展到移动路径（跑下楼梯时需要）
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

	// 反飞行检测
	if !ms.AllowFlight {
		if ms.OnGround || ms.Swimming {
			ms.FlightPossibility = 0
		} else if ms.Climbing {
			ms.FlightPossibility = 0
		} else if dy >= -0.4 {
			// 上升或平飞，权重 +1.5
			ms.FlightPossibility += 1.5
		} else if ms.FlightPossibility > 0 {
			// 下降，权重 -1
			ms.FlightPossibility--
		}

		// 1.5 秒连续飞行 → 拉回
		if ms.FlightPossibility >= 30 { // 20 tps * 1.5s
			if int(ms.FlightPossibility)%10 == 0 {
				p.Teleport(p.Position.X, p.Position.Y-1, p.Position.Z)
			}
			// 2.5 秒 → 踢出
			if ms.FlightPossibility >= 50 { // 20 tps * 2.5s
				p.Kick("飞行在此服务器中不被允许", false)
			}
		}
	}

	ms.IsCollided = ms.OnGround
}

// checkBlockCollision 检查玩家是否在水/藤蔓中
// 对应 PHP Player::checkBlockCollision() L1647-1658
func (p *Player) checkBlockCollision() {
	ms := p.movement
	ms.Swimming = false
	ms.Climbing = false

	lvl, ok := p.Human.Level.(*level.Level)
	if !ok || lvl == nil {
		return
	}

	// 检查玩家脚部和头部位置的方块
	feetX := int32(math.Floor(p.Position.X))
	feetY := int32(math.Floor(p.Position.Y))
	feetZ := int32(math.Floor(p.Position.Z))

	// 检查脚部方块
	checkPositions := [][3]int32{
		{feetX, feetY, feetZ},     // 脚部
		{feetX, feetY + 1, feetZ}, // 身体
	}

	for _, pos := range checkPositions {
		bs := lvl.GetBlock(pos[0], pos[1], pos[2])
		switch bs.ID {
		case 8, 9: // FLOWING_WATER, STILL_WATER
			ms.Swimming = true
		case 106, 65: // VINE, LADDER
			ms.Climbing = true
		}
	}
}

// ── 辅助方法 ──────────────────────────────────────

// revertMovement 回滚移动到指定位置
// 对应 PHP Player::revertMovement() L1852-1862
func (p *Player) revertMovement(x, y, z float64) {
	ms := p.movement
	ms.LastX = x
	ms.LastY = y
	ms.LastZ = z
	ms.FlightPossibility = 0

	p.mu.Lock()
	p.Position = entity.NewVector3(x, y, z)
	p.mu.Unlock()

	// 发送位置重置包
	pk := protocol.NewMovePlayerPacket()
	pk.EntityID = p.GetID()
	pk.X = float32(x)
	pk.Y = float32(y) + EyeHeight
	pk.Z = float32(z)
	pk.Yaw = float32(p.Yaw)
	pk.BodyYaw = float32(p.Yaw)
	pk.Pitch = float32(p.Pitch)
	pk.Mode = 1 // MODE_RESET
	pk.OnGround = p.movement.OnGround
	p.SendPacket(pk)
}

// broadcastMovement 广播移动给其他玩家
// 对应 PHP Level::addEntityMovement()
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
	pk.Mode = 0 // MODE_NORMAL
	pk.OnGround = p.movement.OnGround

	for _, viewer := range viewers {
		if viewer != p {
			viewer.SendPacket(pk)
		}
	}
}

// IsSprinting 返回是否在冲刺（通过 metadata flag）
func (p *Player) IsSprinting() bool {
	return p.Human.Metadata.GetFlag(entity.DataFlags, entity.DataFlagSprinting)
}

// IsSwimming 返回是否在水中
func (p *Player) IsSwimming() bool {
	return p.movement.Swimming
}

// IsClimbing 返回是否在攀爬
func (p *Player) IsClimbing() bool {
	return p.movement.Climbing
}

// IsOnGround 返回是否在地面
func (p *Player) IsOnGround() bool {
	return p.movement.OnGround
}

// IsMoving 返回是否在移动
func (p *Player) IsMoving() bool {
	return p.movement.Moving
}

// SetAllowFlight 设置是否允许飞行
func (p *Player) SetAllowFlight(allow bool) {
	p.movement.AllowFlight = allow
	p.movement.FlightPossibility = 0
}
