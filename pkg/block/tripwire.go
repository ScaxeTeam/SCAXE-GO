package block

type TripwireHookBlock struct {
	TransparentBase
}

func NewTripwireHookBlock() *TripwireHookBlock {
	return &TripwireHookBlock{
		TransparentBase: TransparentBase{
			BlockID:       TRIPWIRE_HOOK,
			BlockName:     "Tripwire Hook",
			BlockHardness: 0,
			BlockToolType: ToolTypeNone,
		},
	}
}

func TripwireHookGetDirection(meta uint8) int {
	return int(meta & 0x03)
}

func TripwireHookIsConnected(meta uint8) bool {
	return meta&0x04 != 0
}

func TripwireHookIsTriggered(meta uint8) bool {
	return meta&0x08 != 0
}

func (b *TripwireHookBlock) GetPlacementMeta(playerDirection int) uint8 {
	if playerDirection < 0 || playerDirection > 3 {
		return 0
	}
	return uint8(playerDirection)
}

func (b *TripwireHookBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(TRIPWIRE_HOOK), Meta: 0, Count: 1}}
}

type TripwireBlock struct {
	TransparentBase
}

func NewTripwireBlock() *TripwireBlock {
	return &TripwireBlock{
		TransparentBase: TransparentBase{
			BlockID:       TRIPWIRE,
			BlockName:     "Tripwire",
			BlockHardness: 0,
			BlockToolType: ToolTypeNone,
		},
	}
}

func TripwireIsTriggered(meta uint8) bool {
	return meta&0x01 != 0
}

func TripwireIsConnected(meta uint8) bool {
	return meta&0x04 != 0
}

func (b *TripwireBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 287, Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewTripwireHookBlock())
	Registry.Register(NewTripwireBlock())
}
