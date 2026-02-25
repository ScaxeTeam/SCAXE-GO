package block

// slab.go — 台阶(半砖)方块
// 对应 PHP class Slab extends Transparent
//
// Meta 编码:
//   bit0-2: 材质变种 (0-7)
//   bit3:   上半 (0x08)
//
// 放置逻辑: 同种半砖可合并为双台阶(Double Slab)
// 双台阶是独立的方块ID，破坏后掉落2个半砖

// SlabBlock 台阶方块
type SlabBlock struct {
	TransparentBase
	DoubleSlabID uint8     // 双台阶方块 ID
	VariantNames [8]string // meta 0-7 对应的材质名称
}

// Slab meta 掩码
const (
	SlabMaskVariant = 0x07
	SlabMaskUpper   = 0x08
)

func newSlab(blockID uint8, name string, doubleSlabID uint8, variantNames [8]string) *SlabBlock {
	return &SlabBlock{
		TransparentBase: TransparentBase{
			BlockID:         blockID,
			BlockName:       name,
			BlockHardness:   2,
			BlockResistance: 30,
			BlockToolType:   ToolTypePickaxe,
			BlockCanPlace:   true,
		},
		DoubleSlabID: doubleSlabID,
		VariantNames: variantNames,
	}
}

// GetDrops 台阶需镐才掉落（掉落时 meta 不含上半标记）
func (b *SlabBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(b.BlockID), Meta: 0, Count: 1}}
}

// ---------- Meta 工具函数 ----------

func SlabIsUpper(meta uint8) bool     { return meta&SlabMaskUpper != 0 }
func SlabGetVariant(meta uint8) uint8 { return meta & SlabMaskVariant }

// ---------- 放置逻辑 ----------

// SlabPlacementResult 台阶放置结果
type SlabPlacementResult struct {
	PlaceBlock    bool  // 是否放置
	ResultBlockID uint8 // 结果方块 ID
	ResultMeta    uint8 // 结果 meta
	MergeToDouble bool  // 是否合并为双台阶
}

// GetSlabPlacementResult 计算台阶放置结果
// 对应 PHP Slab::place() 的完整逻辑
//
// 参数:
//   slabID: 台阶方块 ID
//   doubleSlabID: 双台阶方块 ID
//   placingVariant: 放置的半砖变种 (0-7)
//   face: 放置面 (0=下, 1=上, 2-5=侧面)
//   clickY: 点击 Y (0.0~1.0)
//   targetID, targetMeta: 目标方块（被点击的方块）
//   blockID, blockMeta: 放置位置的方块（可能是空气或已有半砖）
func GetSlabPlacementResult(
	slabID, doubleSlabID uint8,
	placingVariant uint8,
	face int,
	clickY float64,
	targetID, targetMeta uint8,
	blockID, blockMeta uint8,
) SlabPlacementResult {
	// 向下放置
	if face == 0 { // SIDE_DOWN
		if targetID == slabID && SlabIsUpper(targetMeta) && SlabGetVariant(targetMeta) == placingVariant {
			return SlabPlacementResult{PlaceBlock: true, ResultBlockID: doubleSlabID, ResultMeta: placingVariant, MergeToDouble: true}
		}
		if blockID == slabID && SlabGetVariant(blockMeta) == placingVariant {
			return SlabPlacementResult{PlaceBlock: true, ResultBlockID: doubleSlabID, ResultMeta: placingVariant, MergeToDouble: true}
		}
		// 上半台阶
		return SlabPlacementResult{PlaceBlock: true, ResultBlockID: slabID, ResultMeta: placingVariant | SlabMaskUpper}
	}

	// 向上放置
	if face == 1 { // SIDE_UP
		if targetID == slabID && !SlabIsUpper(targetMeta) && SlabGetVariant(targetMeta) == placingVariant {
			return SlabPlacementResult{PlaceBlock: true, ResultBlockID: doubleSlabID, ResultMeta: placingVariant, MergeToDouble: true}
		}
		if blockID == slabID && SlabGetVariant(blockMeta) == placingVariant {
			return SlabPlacementResult{PlaceBlock: true, ResultBlockID: doubleSlabID, ResultMeta: placingVariant, MergeToDouble: true}
		}
		// 下半台阶
		return SlabPlacementResult{PlaceBlock: true, ResultBlockID: slabID, ResultMeta: placingVariant}
	}

	// 侧面放置
	if blockID == slabID && SlabGetVariant(blockMeta) == placingVariant {
		return SlabPlacementResult{PlaceBlock: true, ResultBlockID: doubleSlabID, ResultMeta: placingVariant, MergeToDouble: true}
	}
	meta := placingVariant
	if clickY > 0.5 {
		meta |= SlabMaskUpper
	}
	return SlabPlacementResult{PlaceBlock: true, ResultBlockID: slabID, ResultMeta: meta}
}

// ---------- 碰撞箱 ----------

// SlabBoundingBox 台阶碰撞箱
type SlabBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}

// GetSlabBoundingBox 计算台阶碰撞箱
func GetSlabBoundingBox(x, y, z int, meta uint8) SlabBoundingBox {
	fx, fy, fz := float64(x), float64(y), float64(z)
	if SlabIsUpper(meta) {
		return SlabBoundingBox{fx, fy + 0.5, fz, fx + 1, fy + 1, fz + 1}
	}
	return SlabBoundingBox{fx, fy, fz, fx + 1, fy + 0.5, fz + 1}
}

// ---------- 双台阶 ----------

// DoubleSlabBlock 双台阶方块（破坏掉落2个半砖）
type DoubleSlabBlock struct {
	SolidBase
	SingleSlabID uint8 // 对应的半砖 ID
}

func newDoubleSlab(blockID uint8, name string, singleSlabID uint8) *DoubleSlabBlock {
	return &DoubleSlabBlock{
		SolidBase: SolidBase{
			BlockID:       blockID,
			BlockName:     name,
			BlockHardness: 2,
			BlockToolType: ToolTypePickaxe,
		},
		SingleSlabID: singleSlabID,
	}
}

// GetDrops 双台阶破坏掉落2个半砖
func (b *DoubleSlabBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(b.SingleSlabID), Meta: 0, Count: 2}}
}

// ---------- 材质名称 ----------

var stoneSlabVariants = [8]string{
	"Stone", "Sandstone", "Wooden", "Cobblestone",
	"Brick", "Stone Brick", "Quartz", "Nether Brick",
}

// ---------- 注册 ----------

func init() {
	// 石质台阶 + 双台阶
	Registry.Register(newSlab(SLAB, "Slab", DOUBLE_SLAB, stoneSlabVariants))
	Registry.Register(newDoubleSlab(DOUBLE_SLAB, "Double Slab", SLAB))
}
