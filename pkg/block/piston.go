package block

// ---------- Piston (活塞) ----------
// MCPE 方块 ID 33 — 普通活塞
// meta 低 3 位 = 朝向 (0=下 1=上 2=北 3=南 4=西 5=东)
// meta bit 3 (0x08) = 是否已伸出

type PistonBlock struct {
	SolidBase
}

func NewPistonBlock() *PistonBlock {
	return &PistonBlock{
		SolidBase: SolidBase{
			BlockID:       PISTON,
			BlockName:     "Piston",
			BlockHardness: 0.5,
			BlockToolType: ToolTypeNone,
		},
	}
}

// GetPlacementMeta 活塞根据玩家朝向放置（朝向为面向玩家的方向）
func (b *PistonBlock) GetPlacementMeta(playerDirection int) uint8 {
	// 水平方向映射: 0南→3, 1西→4, 2北→2, 3东→5
	switch playerDirection {
	case 0:
		return 3
	case 1:
		return 4
	case 2:
		return 2
	case 3:
		return 5
	default:
		return 1 // 默认朝上
	}
}

// GetFacing 从 meta 获取活塞朝向 (0-5)
func PistonGetFacing(meta uint8) int {
	return int(meta & 0x07)
}

// IsExtended 活塞是否已伸出
func PistonIsExtended(meta uint8) bool {
	return meta&0x08 != 0
}

func (b *PistonBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(PISTON), Meta: 0, Count: 1}}
}

// ---------- StickyPiston (粘性活塞) ----------
// MCPE 方块 ID 29 — meta 编码与普通活塞完全相同

type StickyPistonBlock struct {
	SolidBase
}

func NewStickyPistonBlock() *StickyPistonBlock {
	return &StickyPistonBlock{
		SolidBase: SolidBase{
			BlockID:       STICKY_PISTON,
			BlockName:     "Sticky Piston",
			BlockHardness: 0.5,
			BlockToolType: ToolTypeNone,
		},
	}
}

func (b *StickyPistonBlock) GetPlacementMeta(playerDirection int) uint8 {
	switch playerDirection {
	case 0:
		return 3
	case 1:
		return 4
	case 2:
		return 2
	case 3:
		return 5
	default:
		return 1
	}
}

func (b *StickyPistonBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(STICKY_PISTON), Meta: 0, Count: 1}}
}

// ---------- PistonHead (活塞臂) ----------
// MCPE 方块 ID 34 — 活塞伸出时生成的方块
// meta 低 3 位 = 朝向 (与活塞相同)
// meta bit 3 (0x08) = 是否为粘性 (1=粘性活塞臂)
//
// 特性:
//   - 不可由玩家放置
//   - 不可被推动
//   - 破坏时不掉落自身（但回缩对应活塞）

type PistonHeadBlock struct {
	TransparentBase
}

func NewPistonHeadBlock() *PistonHeadBlock {
	return &PistonHeadBlock{
		TransparentBase: TransparentBase{
			BlockID:       PISTON_HEAD,
			BlockName:     "Piston Head",
			BlockHardness: 0.5,
			BlockToolType: ToolTypeNone,
			BlockCanPlace: false, // 不可由玩家手动放置
		},
	}
}

// PistonHeadIsSticky 判断活塞臂是否属于粘性活塞
func PistonHeadIsSticky(meta uint8) bool {
	return meta&0x08 != 0
}

// GetDrops 活塞臂不掉落任何物品
func (b *PistonHeadBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

// ---------- 活塞推送逻辑常量 ----------

const (
	PistonFacingDown  = 0
	PistonFacingUp    = 1
	PistonFacingNorth = 2
	PistonFacingSouth = 3
	PistonFacingWest  = 4
	PistonFacingEast  = 5

	PistonMaxPushDistance = 12 // 活塞最大推送距离（方块数）
)

// PistonFacingOffset 返回活塞朝向对应的XYZ偏移量
func PistonFacingOffset(facing int) (dx, dy, dz int) {
	switch facing {
	case PistonFacingDown:
		return 0, -1, 0
	case PistonFacingUp:
		return 0, 1, 0
	case PistonFacingNorth:
		return 0, 0, -1
	case PistonFacingSouth:
		return 0, 0, 1
	case PistonFacingWest:
		return -1, 0, 0
	case PistonFacingEast:
		return 1, 0, 0
	default:
		return 0, 0, 0
	}
}

// ---------- 注册 ----------

func init() {
	Registry.Register(NewPistonBlock())
	Registry.Register(NewStickyPistonBlock())
	Registry.Register(NewPistonHeadBlock())
}
