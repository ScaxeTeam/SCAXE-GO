package block

// ---------- TripwireHookBlock 绊线钩 ----------
// MCPE 方块 ID 131
// meta 低 2 位 = 朝向 (0=南 1=西 2=北 3=东)
// meta bit 2 (0x04) = 是否已连接
// meta bit 3 (0x08) = 是否被触发

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

// TripwireHookGetDirection 从 meta 获取绊线钩朝向
func TripwireHookGetDirection(meta uint8) int {
	return int(meta & 0x03)
}

// TripwireHookIsConnected 是否已连接到另一个钩
func TripwireHookIsConnected(meta uint8) bool {
	return meta&0x04 != 0
}

// TripwireHookIsTriggered 是否被实体触发
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

// ---------- TripwireBlock 绊线 ----------
// MCPE 方块 ID 132
// meta bit 0 (0x01) = 是否被触发
// meta bit 2 (0x04) = 是否已连接到钩
// meta bit 3 (0x08) = 是否被解除 (拆除时 disarmed)

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

// TripwireIsTriggered 绊线是否被触发
func TripwireIsTriggered(meta uint8) bool {
	return meta&0x01 != 0
}

// TripwireIsConnected 绊线是否连接到钩
func TripwireIsConnected(meta uint8) bool {
	return meta&0x04 != 0
}

func (b *TripwireBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 287, Meta: 0, Count: 1}} // STRING item
}

func init() {
	Registry.Register(NewTripwireHookBlock())
	Registry.Register(NewTripwireBlock())
}
