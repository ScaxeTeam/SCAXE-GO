package block

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

const (
	RailStraightNorthSouth = 0
	RailStraightEastWest   = 1
	RailAscendEast         = 2
	RailAscendWest         = 3
	RailAscendNorth        = 4
	RailAscendSouth        = 5
	RailCurvedSouthEast    = 6
	RailCurvedSouthWest    = 7
	RailCurvedNorthWest    = 8
	RailCurvedNorthEast    = 9
)
func RailIsPowered(meta uint8) bool {
	return meta&0x08 != 0
}

func init() {
	Registry.Register(&railBlock{})
	Registry.Register(&poweredRailBlock{})
	Registry.Register(&detectorRailBlock{})
	Registry.Register(&activatorRailBlock{})
}
