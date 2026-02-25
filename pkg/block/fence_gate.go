package block

// fence_gate.go — 栅栏门方块（6种木材变种）
// 对应 PHP class FenceGate extends Transparent implements IRedstone
// 及 FenceGateSpruce, FenceGateBirch, FenceGateJungle, FenceGateDarkOak, FenceGateAcacia
//
// Meta 位编码:
//   bit0-1: 朝向 (0-3)
//   bit2:   打开 (0x04)

// FenceGateBlock 栅栏门方块
type FenceGateBlock struct {
	TransparentBase
}

// FenceGate meta 掩码
const (
	FenceGateMaskDirection = 0x03
	FenceGateMaskOpen      = 0x04
)

func newFenceGate(blockID uint8, name string) *FenceGateBlock {
	return &FenceGateBlock{
		TransparentBase: TransparentBase{
			BlockID:         blockID,
			BlockName:       name,
			BlockHardness:   2,
			BlockResistance: 10,
			BlockToolType:   ToolTypeAxe,
			BlockCanPlace:   true,
		},
	}
}

// CanBeActivated 栅栏门可以被右键激活
func (b *FenceGateBlock) CanBeActivated() bool {
	return true
}

// OnActivate 栅栏门右键交互 — 开/关切换
// 实际 meta 变更由服务器层通过 FenceGateOnActivate() 执行
func (b *FenceGateBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

// GetFuelTime 木质栅栏门可作燃料
func (b *FenceGateBlock) GetFuelTime() int {
	return 300
}

// GetDrops 栅栏门掉落自身
func (b *FenceGateBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(b.BlockID), Meta: 0, Count: 1}}
}

// ---------- Meta 工具函数 ----------

// FenceGateIsOpen 判断栅栏门是否打开
func FenceGateIsOpen(meta uint8) bool {
	return meta&FenceGateMaskOpen != 0
}

// FenceGateGetDirection 获取栅栏门朝向 (0-3)
func FenceGateGetDirection(meta uint8) uint8 {
	return meta & FenceGateMaskDirection
}

// FenceGateToggleOpen 切换开关，返回新 meta
func FenceGateToggleOpen(meta uint8) uint8 {
	return meta ^ FenceGateMaskOpen
}

// ---------- 放置 ----------

// FenceGateDirectionToMeta 玩家朝向 → meta 映射
// 对应 PHP FenceGate::place() 中的 $faces
var FenceGateDirectionToMeta = [4]uint8{3, 0, 1, 2}

// GetFenceGatePlacementMeta 获取放置栅栏门的 meta
func GetFenceGatePlacementMeta(playerDirection int) uint8 {
	return FenceGateDirectionToMeta[playerDirection&0x03]
}

// ---------- 碰撞箱 ----------

// FenceGateBoundingBox 栅栏门碰撞箱
type FenceGateBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
	HasCollision     bool // 打开时无碰撞
}

// GetFenceGateBoundingBox 计算栅栏门碰撞箱
// 对应 PHP FenceGate::recalculateBoundingBox()
// 打开时返回 HasCollision=false（无碰撞）
// 关闭时根据朝向返回竖直薄墙（高1.5格）
func GetFenceGateBoundingBox(x, y, z int, meta uint8) FenceGateBoundingBox {
	fx, fy, fz := float64(x), float64(y), float64(z)

	if FenceGateIsOpen(meta) {
		return FenceGateBoundingBox{HasCollision: false}
	}

	dir := FenceGateGetDirection(meta)
	// 方向 0 和 2: 南北朝向，薄墙沿 Z 轴中间
	// 方向 1 和 3: 东西朝向，薄墙沿 X 轴中间
	if dir == 0 || dir == 2 {
		return FenceGateBoundingBox{
			MinX: fx, MinY: fy, MinZ: fz + 0.375,
			MaxX: fx + 1, MaxY: fy + 1.5, MaxZ: fz + 0.625,
			HasCollision: true,
		}
	}
	return FenceGateBoundingBox{
		MinX: fx + 0.375, MinY: fy, MinZ: fz,
		MaxX: fx + 0.625, MaxY: fy + 1.5, MaxZ: fz + 1,
		HasCollision: true,
	}
}

// ---------- 注册 ----------

func init() {
	Registry.Register(newFenceGate(FENCE_GATE, "Oak Fence Gate"))
	Registry.Register(newFenceGate(FENCE_GATE_SPRUCE, "Spruce Fence Gate"))
	Registry.Register(newFenceGate(FENCE_GATE_BIRCH, "Birch Fence Gate"))
	Registry.Register(newFenceGate(FENCE_GATE_JUNGLE, "Jungle Fence Gate"))
	Registry.Register(newFenceGate(FENCE_GATE_DARK_OAK, "Dark Oak Fence Gate"))
	Registry.Register(newFenceGate(FENCE_GATE_ACACIA, "Acacia Fence Gate"))
}
