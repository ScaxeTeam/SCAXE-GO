package block

// ── Fire (ID 51) ────────────────────────────────────────────────
// Light level 15. Has entity collision (burning). Replaceable.
// Meta 0-15 = fire age. Spreads to neighbors based on burnChance/burnAbility.
// Tick rate: 30 ticks.

type fireBlock struct{ DefaultBlockInteraction }

func (b *fireBlock) GetID() uint8                { return FIRE }
func (b *fireBlock) GetName() string             { return "Fire Block" }
func (b *fireBlock) GetHardness() float64        { return 0 }
func (b *fireBlock) GetBlastResistance() float64 { return 0 }
func (b *fireBlock) GetLightLevel() uint8        { return 15 }
func (b *fireBlock) GetLightFilter() uint8       { return 0 }
func (b *fireBlock) IsSolid() bool               { return false }
func (b *fireBlock) IsTransparent() bool         { return true }
func (b *fireBlock) CanBePlaced() bool           { return false }
func (b *fireBlock) CanBeReplaced() bool         { return true }
func (b *fireBlock) GetToolType() int            { return ToolTypeNone }
func (b *fireBlock) GetToolTier() int            { return 0 }
func (b *fireBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil // fire drops nothing
}

// FireTickRate is the base tick rate for fire updates (30 ticks).
const FireTickRate = 30

func init() {
	Registry.Register(&fireBlock{})
}
