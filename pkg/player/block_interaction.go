package player

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/logger"
)
func (p *Player) HandleUseItem(x, y, z int32, face int, fx, fy, fz float32) {
	if !p.Spawned || !p.Connected {
		return
	}

	lvl, ok := p.Human.Level.(*level.Level)
	if !ok || lvl == nil {
		return
	}

	if face >= 0 && face <= 5 {
		p.handleBlockPlace(lvl, x, y, z, face, fx, fy, fz)
	} else if face == 0xFF {
		p.handleItemUse()
	}
}
func (p *Player) handleBlockPlace(lvl *level.Level, x, y, z int32, face int, fx, fy, fz float32) {
	dx := float64(x) + 0.5 - p.Position.X
	dy := float64(y) + 0.5 - p.Position.Y
	dz := float64(z) + 0.5 - p.Position.Z
	distSq := dx*dx + dy*dy + dz*dz

	if distSq > 169 {
		logger.Debug("Block place too far",
			"player", p.Username,
			"distSq", distSq)
		return
	}
	if !p.canInteract(float64(x)+0.5, float64(y)+0.5, float64(z)+0.5, 13) {
		return
	}
	targetX, targetY, targetZ := getBlockSide(x, y, z, face)
	if targetY < level.YMin || targetY >= level.YMax {
		return
	}
	clickedBlock := lvl.GetBlock(x, y, z)
	targetBlock := lvl.GetBlock(targetX, targetY, targetZ)
	clickedBehavior := block.Registry.GetBehavior(clickedBlock.ID)
	if clickedBehavior != nil && clickedBehavior.CanBeActivated() {
		ctx := &block.BlockContext{
			X: int(x), Y: int(y), Z: int(z),
			Meta:   clickedBlock.Meta,
			Face:   face,
			ClickX: float64(fx), ClickY: float64(fy), ClickZ: float64(fz),
		}
		if clickedBehavior.OnActivate(ctx, p.GetID()) {
			return
		}
	}
	heldItem := p.Inventory.GetItemInHand()
	if heldItem.ID == 0 {
		return
	}
	targetBehavior := block.Registry.GetBehavior(targetBlock.ID)
	if targetBehavior != nil && !targetBehavior.CanBeReplaced() {
		return
	}
	placeBehavior := block.Registry.GetBehavior(uint8(heldItem.ID))
	if placeBehavior != nil && placeBehavior.CanBePlaced() {
		ctx := &block.BlockContext{
			X: int(targetX), Y: int(targetY), Z: int(targetZ),
			Meta:   uint8(heldItem.Meta),
			Face:   face,
			ClickX: float64(fx), ClickY: float64(fy), ClickZ: float64(fz),
		}

		if placeBehavior.Place(ctx) {
			lvl.SetBlock(targetX, targetY, targetZ, byte(heldItem.ID), byte(heldItem.Meta), true)
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
func (p *Player) HandleRemoveBlock(x, y, z int32) {
	if !p.Spawned || !p.Connected {
		return
	}

	lvl, ok := p.Human.Level.(*level.Level)
	if !ok || lvl == nil {
		return
	}
	maxDist := 6.0
	if p.IsCreative() {
		maxDist = 13.0
	}
	if !p.canInteract(float64(x)+0.5, float64(y)+0.5, float64(z)+0.5, maxDist) {
		return
	}
	bs := lvl.GetBlock(x, y, z)
	if bs.ID == block.AIR {
		return
	}

	behavior := block.Registry.GetBehavior(bs.ID)
	if behavior == nil {
		return
	}
	_ = p.Inventory.GetItemInHand()
	toolType := 0
	toolTier := 0
	if !behavior.IsBreakable(toolType, toolTier) && !p.IsCreative() {
		return
	}
	ctx := &block.BlockContext{
		X: int(x), Y: int(y), Z: int(z),
		Meta: bs.Meta,
	}

	if behavior.OnBreak(ctx, toolType, toolTier) {
		lvl.SetBlock(x, y, z, block.AIR, 0, true)
		if p.IsSurvival() {
			drops := behavior.GetDrops(toolType, toolTier)
			for _, drop := range drops {
				if drop.ID != 0 && drop.Count > 0 {
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
func (p *Player) HandlePlayerAction(action int32, x, y, z int32, face int) {
	if !p.Spawned || !p.Connected {
		return
	}

	switch action {
	case ActionStartBreak:
		logger.Debug("Start break",
			"player", p.Username,
			"x", x, "y", y, "z", z)

	case ActionAbortBreak:
		logger.Debug("Abort break", "player", p.Username)

	case ActionStopBreak:

	case ActionReleaseItem:
		logger.Debug("Release item", "player", p.Username)

	case ActionJump:

	case ActionStartSprint:
		p.Human.SetSprinting(true)
		logger.Debug("Start sprint", "player", p.Username)

	case ActionStopSprint:
		p.Human.SetSprinting(false)
		logger.Debug("Stop sprint", "player", p.Username)

	case ActionStartSneak:
		p.Human.Metadata.SetFlag(1, 1, true)
		logger.Debug("Start sneak", "player", p.Username)

	case ActionStopSneak:
		p.Human.Metadata.SetFlag(1, 1, false)
		logger.Debug("Stop sneak", "player", p.Username)

	case ActionRespawn:
		logger.Debug("Respawn", "player", p.Username)
	}
}
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
	dirX := -math.Sin(p.Yaw/180*math.Pi) * math.Cos(p.Pitch/180*math.Pi)
	dirY := -math.Sin(p.Pitch / 180 * math.Pi)
	dirZ := math.Cos(p.Yaw/180*math.Pi) * math.Cos(p.Pitch/180*math.Pi)

	eyeDot := dirX*eyeX + dirY*eyeY + dirZ*eyeZ
	targetDot := dirX*x + dirY*y + dirZ*z

	return (targetDot - eyeDot) >= -math.Sqrt(3)/2
}
func getBlockSide(x, y, z int32, face int) (int32, int32, int32) {
	switch face {
	case 0:
		return x, y - 1, z
	case 1:
		return x, y + 1, z
	case 2:
		return x, y, z - 1
	case 3:
		return x, y, z + 1
	case 4:
		return x - 1, y, z
	case 5:
		return x + 1, y, z
	default:
		return x, y, z
	}
}
func (p *Player) IsCreative() bool {
	return (p.Gamemode & 0x01) > 0
}
func (p *Player) IsSurvival() bool {
	return (p.Gamemode & 0x01) == 0
}
func (p *Player) IsAdventure() bool {
	return (p.Gamemode & 0x02) > 0
}
func (p *Player) IsSpectator() bool {
	return p.Gamemode == 3
}
func (p *Player) handleItemUse() {
	if p.IsSpectator() {
		return
	}

	heldItem := p.Inventory.GetItemInHand()
	if heldItem.ID == 0 {
		return
	}
	logger.Debug("Item use (air)",
		"player", p.Username,
		"itemID", heldItem.ID)
}
