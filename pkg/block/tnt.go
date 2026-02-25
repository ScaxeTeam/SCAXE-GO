package block

// ── TNT (ID 46) ─────────────────────────────────────────────────
// Activated by flint & steel or redstone. Spawns PrimedTNT entity.
// Burns easily: burnChance=15, burnAbility=100.

type tntBlock struct{ DefaultBlockInteraction }

func (b *tntBlock) GetID() uint8                { return TNT }
func (b *tntBlock) GetName() string             { return "TNT" }
func (b *tntBlock) GetHardness() float64        { return 0 }
func (b *tntBlock) GetBlastResistance() float64 { return 0 }
func (b *tntBlock) GetLightLevel() uint8        { return 0 }
func (b *tntBlock) GetLightFilter() uint8       { return 15 }
func (b *tntBlock) IsSolid() bool               { return true }
func (b *tntBlock) IsTransparent() bool         { return false }
func (b *tntBlock) CanBePlaced() bool           { return true }
func (b *tntBlock) CanBeReplaced() bool         { return false }
func (b *tntBlock) GetToolType() int            { return ToolTypeNone }
func (b *tntBlock) GetToolTier() int            { return 0 }
func (b *tntBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(TNT), Meta: 0, Count: 1}}
}

// TNTBurnChance and TNTBurnAbility for fire spread calculations
const (
	TNTBurnChance  = 15
	TNTBurnAbility = 100
)

func init() {
	Registry.Register(&tntBlock{})
}
