package player

// block_interaction.go — 玩家方块交互处理
// 对应 PHP Player.php 的 USE_ITEM_PACKET / REMOVE_BLOCK_PACKET / PLAYER_ACTION_PACKET 处理器
//
// 三大交互类型：
//   1. 放置方块 (UseItem face=0-5) — 验证距离+角度、调用 BlockBehavior.Place、扣减物品
//   2. 破坏方块 (RemoveBlock) — 验证距离、调用 BlockBehavior.OnBreak、掉落物品
//   3. 玩家动作 (PlayerAction) — 开始/停止/中止挖掘、冲刺、潜行、跳跃等

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/logger"
)

// 动作常量定义在 player.go 中（ActionStartBreak 等）

// ── 方块放置 ──────────────────────────────────────

// HandleUseItem 处理 UseItemPacket (face=0-5: 放置方块, face=0xFF: 右键空气)
// 对应 PHP Player::handleDataPacket() case USE_ITEM_PACKET (L2568-2780)
func (p *Player) HandleUseItem(x, y, z int32, face int, fx, fy, fz float32) {
	if !p.Spawned || !p.Connected {
		return
	}

	lvl, ok := p.Human.Level.(*level.Level)
	if !ok || lvl == nil {
		return
	}

	if face >= 0 && face <= 5 {
		// 放置方块
		p.handleBlockPlace(lvl, x, y, z, face, fx, fy, fz)
	} else if face == 0xFF {
		// 右键空气（使用物品 — 食物/弓/雪球等）
		p.handleItemUse()
	}
}

// handleBlockPlace 放置方块逻辑
// 对应 PHP USE_ITEM_PACKET face=0-5 (L2578-2611)
func (p *Player) handleBlockPlace(lvl *level.Level, x, y, z int32, face int, fx, fy, fz float32) {
	// 距离验证（13 格以内）
	dx := float64(x) + 0.5 - p.Position.X
	dy := float64(y) + 0.5 - p.Position.Y
	dz := float64(z) + 0.5 - p.Position.Z
	distSq := dx*dx + dy*dy + dz*dz

	if distSq > 169 { // 13^2
		logger.Debug("Block place too far",
			"player", p.Username,
			"distSq", distSq)
		return
	}

	// 朝向验证 (canInteract)
	if !p.canInteract(float64(x)+0.5, float64(y)+0.5, float64(z)+0.5, 13) {
		return
	}

	// 计算放置目标坐标
	targetX, targetY, targetZ := getBlockSide(x, y, z, face)

	// 验证放置位置范围
	if targetY < level.YMin || targetY >= level.YMax {
		return
	}

	// 获取目标位置和点击位置的方块
	clickedBlock := lvl.GetBlock(x, y, z)
	targetBlock := lvl.GetBlock(targetX, targetY, targetZ)

	// 尝试激活点击的方块（箱子/门/按钮等）
	clickedBehavior := block.Registry.GetBehavior(clickedBlock.ID)
	if clickedBehavior != nil && clickedBehavior.CanBeActivated() {
		ctx := &block.BlockContext{
			X: int(x), Y: int(y), Z: int(z),
			Meta:   clickedBlock.Meta,
			Face:   face,
			ClickX: float64(fx), ClickY: float64(fy), ClickZ: float64(fz),
		}
		if clickedBehavior.OnActivate(ctx, p.GetID()) {
			return // 方块消耗了交互（如打开箱子）
		}
	}

	// 手持物品的方块ID
	heldItem := p.Inventory.GetItemInHand()
	if heldItem.ID == 0 {
		return // 空手不能放置
	}

	// 目标位置必须可替换（空气/水/草等）
	targetBehavior := block.Registry.GetBehavior(targetBlock.ID)
	if targetBehavior != nil && !targetBehavior.CanBeReplaced() {
		// 回滚: 发送原始方块给客户端
		return
	}

	// 不能放置在自身碰撞箱内（防止把自己卡住）
	// TODO: 完善碰撞箱检查

	// 调用 BlockBehavior.Place
	placeBehavior := block.Registry.GetBehavior(uint8(heldItem.ID))
	if placeBehavior != nil && placeBehavior.CanBePlaced() {
		ctx := &block.BlockContext{
			X: int(targetX), Y: int(targetY), Z: int(targetZ),
			Meta:   uint8(heldItem.Meta),
			Face:   face,
			ClickX: float64(fx), ClickY: float64(fy), ClickZ: float64(fz),
		}

		if placeBehavior.Place(ctx) {
			// 在世界中设置方块
			lvl.SetBlock(targetX, targetY, targetZ, byte(heldItem.ID), byte(heldItem.Meta), true)

			// 生存模式扣减物品
			if p.IsSurvival() {
				heldItem.Count--
				if heldItem.Count <= 0 {
					heldItem.ID = 0
					heldItem.Meta = 0
					heldItem.Count = 0
				}
				p.Inventory.SetItemInHand(heldItem)
			}

			logger.Debug("Block placed",
				"player", p.Username,
				"x", targetX, "y", targetY, "z", targetZ,
				"id", heldItem.ID, "meta", heldItem.Meta)
		}
	}
}

// ── 方块破坏 ──────────────────────────────────────

// HandleRemoveBlock 处理 RemoveBlockPacket (破坏方块)
// 对应 PHP Player::handleDataPacket() case REMOVE_BLOCK_PACKET (L3055-3090)
func (p *Player) HandleRemoveBlock(x, y, z int32) {
	if !p.Spawned || !p.Connected {
		return
	}

	lvl, ok := p.Human.Level.(*level.Level)
	if !ok || lvl == nil {
		return
	}

	// 距离验证
	maxDist := 6.0
	if p.IsCreative() {
		maxDist = 13.0
	}
	if !p.canInteract(float64(x)+0.5, float64(y)+0.5, float64(z)+0.5, maxDist) {
		return
	}

	// 获取目标方块
	bs := lvl.GetBlock(x, y, z)
	if bs.ID == block.AIR {
		return
	}

	behavior := block.Registry.GetBehavior(bs.ID)
	if behavior == nil {
		return
	}

	// 获取手持工具
	_ = p.Inventory.GetItemInHand()
	toolType := 0 // TODO: 从 Item 获取工具类型
	toolTier := 0 // TODO: 从 Item 获取工具等级

	// 检查是否可破坏
	if !behavior.IsBreakable(toolType, toolTier) && !p.IsCreative() {
		// 回滚方块
		return
	}

	// 调用 OnBreak
	ctx := &block.BlockContext{
		X: int(x), Y: int(y), Z: int(z),
		Meta: bs.Meta,
	}

	if behavior.OnBreak(ctx, toolType, toolTier) {
		// 设置为空气
		lvl.SetBlock(x, y, z, block.AIR, 0, true)

		// 生存模式处理掉落物
		if p.IsSurvival() {
			drops := behavior.GetDrops(toolType, toolTier)
			for _, drop := range drops {
				if drop.ID != 0 && drop.Count > 0 {
					// TODO: 在世界中生成掉落物实体
					_ = drop
				}
			}
		}

		logger.Debug("Block broken",
			"player", p.Username,
			"x", x, "y", y, "z", z,
			"id", bs.ID, "meta", bs.Meta)
	}
}

// ── 玩家动作 ──────────────────────────────────────

// HandlePlayerAction 处理 PlayerActionPacket
// 对应 PHP Player::handleDataPacket() case PLAYER_ACTION_PACKET (L2842-2990)
func (p *Player) HandlePlayerAction(action int32, x, y, z int32, face int) {
	if !p.Spawned || !p.Connected {
		return
	}

	switch action {
	case ActionStartBreak:
		// 开始挖掘 — 记录开始时间，触发 PlayerInteractEvent
		// TODO: 验证距离，触发事件
		logger.Debug("Start break",
			"player", p.Username,
			"x", x, "y", y, "z", z)

	case ActionAbortBreak:
		// 中止挖掘 — 重置挖掘状态
		logger.Debug("Abort break", "player", p.Username)

	case ActionStopBreak:
		// 停止挖掘 — 无特殊处理（移植自 PHP 空实现）

	case ActionReleaseItem:
		// 释放物品（弓箭射出等）
		// TODO: 弓箭逻辑
		logger.Debug("Release item", "player", p.Username)

	case ActionJump:
		// 跳跃 — 饥饿消耗
		// TODO: p.exhaust(0.2, ExhaustCauseJump)

	case ActionStartSprint:
		p.Human.SetSprinting(true)
		logger.Debug("Start sprint", "player", p.Username)

	case ActionStopSprint:
		p.Human.SetSprinting(false)
		logger.Debug("Stop sprint", "player", p.Username)

	case ActionStartSneak:
		p.Human.Metadata.SetFlag(1, 1, true) // DATA_FLAGS, DATA_FLAG_SNEAKING
		logger.Debug("Start sneak", "player", p.Username)

	case ActionStopSneak:
		p.Human.Metadata.SetFlag(1, 1, false)
		logger.Debug("Stop sneak", "player", p.Username)

	case ActionRespawn:
		// 重生 — 由死亡/重生模块处理
		logger.Debug("Respawn", "player", p.Username)
	}
}

// ── 辅助方法 ──────────────────────────────────────

// canInteract 验证玩家是否可以与指定位置交互 (距离 + 朝向)
// 对应 PHP Player::canInteract() L2130-2140
func (p *Player) canInteract(x, y, z float64, maxDistance float64) bool {
	eyeX := p.Position.X
	eyeY := p.Position.Y + EyeHeight
	eyeZ := p.Position.Z

	dx := x - eyeX
	dy := y - eyeY
	dz := z - eyeZ
	distSq := dx*dx + dy*dy + dz*dz

	if distSq > maxDistance*maxDistance {
		return false
	}

	// 朝向验证: 目标点必须在视线前方
	dirX := -math.Sin(p.Yaw/180*math.Pi) * math.Cos(p.Pitch/180*math.Pi)
	dirY := -math.Sin(p.Pitch / 180 * math.Pi)
	dirZ := math.Cos(p.Yaw/180*math.Pi) * math.Cos(p.Pitch/180*math.Pi)

	eyeDot := dirX*eyeX + dirY*eyeY + dirZ*eyeZ
	targetDot := dirX*x + dirY*y + dirZ*z

	return (targetDot - eyeDot) >= -math.Sqrt(3)/2
}

// getBlockSide 获取方块某个面的相邻坐标
// 对应 PHP Block::getSide()
func getBlockSide(x, y, z int32, face int) (int32, int32, int32) {
	switch face {
	case 0: // Down
		return x, y - 1, z
	case 1: // Up
		return x, y + 1, z
	case 2: // North (Z-)
		return x, y, z - 1
	case 3: // South (Z+)
		return x, y, z + 1
	case 4: // West (X-)
		return x - 1, y, z
	case 5: // East (X+)
		return x + 1, y, z
	default:
		return x, y, z
	}
}

// IsCreative 返回是否为创造模式
func (p *Player) IsCreative() bool {
	return (p.Gamemode & 0x01) > 0
}

// IsSurvival 返回是否为生存模式
func (p *Player) IsSurvival() bool {
	return (p.Gamemode & 0x01) == 0
}

// IsAdventure 返回是否为冒险模式
func (p *Player) IsAdventure() bool {
	return (p.Gamemode & 0x02) > 0
}

// IsSpectator 返回是否为观察者模式
func (p *Player) IsSpectator() bool {
	return p.Gamemode == 3
}

// handleItemUse 右键空气使用物品 (食物/弓/雪球等)
// 对应 PHP USE_ITEM_PACKET face=0xFF (L2612-2780)
func (p *Player) handleItemUse() {
	if p.IsSpectator() {
		return
	}

	heldItem := p.Inventory.GetItemInHand()
	if heldItem.ID == 0 {
		return
	}

	// TODO: 根据物品类型分发处理
	// - 食物: 开始进食
	// - 弓: 开始蓄力 (ACTION_RELEASE_ITEM 时射出)
	// - 雪球/鸡蛋: 抛射
	// - 药水: 饮用/投掷
	logger.Debug("Item use (air)",
		"player", p.Username,
		"itemID", heldItem.ID)
}
