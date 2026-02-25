package block

// rail.go — 4 种铁轨方块
// Rail (66): 普通铁轨，可弯曲 (meta 0-9)
// PoweredRail (27): 充能铁轨，bit 0x08 = powered
// DetectorRail (28): 探测铁轨，bit 0x08 = active
// ActivatorRail (126): 激活铁轨，bit 0x08 = powered

// ── Rail (ID 66) ────────────────────────────────────────────────
// 10 directions (meta 0-9): straight + curved.

type railBlock struct{ DefaultBlockInteraction }

func (b *railBlock) GetID() uint8                { return RAIL }
func (b *railBlock) GetName() string             { return "Rail" }
func (b *railBlock) GetHardness() float64        { return 0.7 }
func (b *railBlock) GetBlastResistance() float64 { return 3.5 }
func (b *railBlock) GetLightLevel() uint8        { return 0 }
func (b *railBlock) GetLightFilter() uint8       { return 0 }
func (b *railBlock) IsSolid() bool               { return false }
func (b *railBlock) IsTransparent() bool         { return true }
func (b *railBlock) CanBePlaced() bool           { return true }
func (b *railBlock) CanBeReplaced() bool         { return false }
func (b *railBlock) GetToolType() int            { return ToolTypePickaxe }
func (b *railBlock) GetToolTier() int            { return 0 }
func (b *railBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(RAIL), Meta: 0, Count: 1}}
}

// ── Powered Rail (ID 27) ────────────────────────────────────────
// Meta bit 0x08 = powered. Lower 3 bits = direction (0-5).

type poweredRailBlock struct{ DefaultBlockInteraction }

func (b *poweredRailBlock) GetID() uint8                { return POWERED_RAIL }
func (b *poweredRailBlock) GetName() string             { return "Powered Rail" }
func (b *poweredRailBlock) GetHardness() float64        { return 0.7 }
func (b *poweredRailBlock) GetBlastResistance() float64 { return 3.5 }
func (b *poweredRailBlock) GetLightLevel() uint8        { return 0 }
func (b *poweredRailBlock) GetLightFilter() uint8       { return 0 }
func (b *poweredRailBlock) IsSolid() bool               { return false }
func (b *poweredRailBlock) IsTransparent() bool         { return true }
func (b *poweredRailBlock) CanBePlaced() bool           { return true }
func (b *poweredRailBlock) CanBeReplaced() bool         { return false }
func (b *poweredRailBlock) GetToolType() int            { return ToolTypePickaxe }
func (b *poweredRailBlock) GetToolTier() int            { return 0 }
func (b *poweredRailBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(POWERED_RAIL), Meta: 0, Count: 1}}
}

// ── Detector Rail (ID 28) ───────────────────────────────────────
// Meta bit 0x08 = entity on rail (active). Lower 3 bits = direction.

type detectorRailBlock struct{ DefaultBlockInteraction }

func (b *detectorRailBlock) GetID() uint8                { return DETECTOR_RAIL }
func (b *detectorRailBlock) GetName() string             { return "Detector Rail" }
func (b *detectorRailBlock) GetHardness() float64        { return 0.7 }
func (b *detectorRailBlock) GetBlastResistance() float64 { return 3.5 }
func (b *detectorRailBlock) GetLightLevel() uint8        { return 0 }
func (b *detectorRailBlock) GetLightFilter() uint8       { return 0 }
func (b *detectorRailBlock) IsSolid() bool               { return false }
func (b *detectorRailBlock) IsTransparent() bool         { return true }
func (b *detectorRailBlock) CanBePlaced() bool           { return true }
func (b *detectorRailBlock) CanBeReplaced() bool         { return false }
func (b *detectorRailBlock) GetToolType() int            { return ToolTypePickaxe }
func (b *detectorRailBlock) GetToolTier() int            { return 0 }
func (b *detectorRailBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(DETECTOR_RAIL), Meta: 0, Count: 1}}
}

// ── Activator Rail (ID 126) ─────────────────────────────────────
// Meta bit 0x08 = powered. Lower 3 bits = direction.

type activatorRailBlock struct{ DefaultBlockInteraction }

func (b *activatorRailBlock) GetID() uint8                { return ACTIVATOR_RAIL }
func (b *activatorRailBlock) GetName() string             { return "Activator Rail" }
func (b *activatorRailBlock) GetHardness() float64        { return 0.7 }
func (b *activatorRailBlock) GetBlastResistance() float64 { return 3.5 }
func (b *activatorRailBlock) GetLightLevel() uint8        { return 0 }
func (b *activatorRailBlock) GetLightFilter() uint8       { return 0 }
func (b *activatorRailBlock) IsSolid() bool               { return false }
func (b *activatorRailBlock) IsTransparent() bool         { return true }
func (b *activatorRailBlock) CanBePlaced() bool           { return true }
func (b *activatorRailBlock) CanBeReplaced() bool         { return false }
func (b *activatorRailBlock) GetToolType() int            { return ToolTypePickaxe }
func (b *activatorRailBlock) GetToolTier() int            { return 0 }
func (b *activatorRailBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(ACTIVATOR_RAIL), Meta: 0, Count: 1}}
}

// ── Rail meta direction constants ───────────────────────────────

const (
	RailStraightNorthSouth = 0
	RailStraightEastWest   = 1
	RailAscendEast         = 2
	RailAscendWest         = 3
	RailAscendNorth        = 4
	RailAscendSouth        = 5
	RailCurvedSouthEast    = 6 // only normal rail
	RailCurvedSouthWest    = 7
	RailCurvedNorthWest    = 8
	RailCurvedNorthEast    = 9
)

// RailIsPowered checks if a powered/activator rail meta has the powered bit set.
func RailIsPowered(meta uint8) bool {
	return meta&0x08 != 0
}

// ── Registration ────────────────────────────────────────────────

func init() {
	Registry.Register(&railBlock{})
	Registry.Register(&poweredRailBlock{})
	Registry.Register(&detectorRailBlock{})
	Registry.Register(&activatorRailBlock{})
}
