package block

type SlabBlock struct {
	TransparentBase
	DoubleSlabID uint8
	VariantNames [8]string
}

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

func (b *SlabBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(b.BlockID), Meta: 0, Count: 1}}
}

func SlabIsUpper(meta uint8) bool     { return meta&SlabMaskUpper != 0 }
func SlabGetVariant(meta uint8) uint8 { return meta & SlabMaskVariant }

type SlabPlacementResult struct {
	PlaceBlock    bool
	ResultBlockID uint8
	ResultMeta    uint8
	MergeToDouble bool
}

func GetSlabPlacementResult(
	slabID, doubleSlabID uint8,
	placingVariant uint8,
	face int,
	clickY float64,
	targetID, targetMeta uint8,
	blockID, blockMeta uint8,
) SlabPlacementResult {

	if face == 0 {
		if targetID == slabID && SlabIsUpper(targetMeta) && SlabGetVariant(targetMeta) == placingVariant {
			return SlabPlacementResult{PlaceBlock: true, ResultBlockID: doubleSlabID, ResultMeta: placingVariant, MergeToDouble: true}
		}
		if blockID == slabID && SlabGetVariant(blockMeta) == placingVariant {
			return SlabPlacementResult{PlaceBlock: true, ResultBlockID: doubleSlabID, ResultMeta: placingVariant, MergeToDouble: true}
		}

		return SlabPlacementResult{PlaceBlock: true, ResultBlockID: slabID, ResultMeta: placingVariant | SlabMaskUpper}
	}

	if face == 1 {
		if targetID == slabID && !SlabIsUpper(targetMeta) && SlabGetVariant(targetMeta) == placingVariant {
			return SlabPlacementResult{PlaceBlock: true, ResultBlockID: doubleSlabID, ResultMeta: placingVariant, MergeToDouble: true}
		}
		if blockID == slabID && SlabGetVariant(blockMeta) == placingVariant {
			return SlabPlacementResult{PlaceBlock: true, ResultBlockID: doubleSlabID, ResultMeta: placingVariant, MergeToDouble: true}
		}

		return SlabPlacementResult{PlaceBlock: true, ResultBlockID: slabID, ResultMeta: placingVariant}
	}

	if blockID == slabID && SlabGetVariant(blockMeta) == placingVariant {
		return SlabPlacementResult{PlaceBlock: true, ResultBlockID: doubleSlabID, ResultMeta: placingVariant, MergeToDouble: true}
	}
	meta := placingVariant
	if clickY > 0.5 {
		meta |= SlabMaskUpper
	}
	return SlabPlacementResult{PlaceBlock: true, ResultBlockID: slabID, ResultMeta: meta}
}

type SlabBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}

func GetSlabBoundingBox(x, y, z int, meta uint8) SlabBoundingBox {
	fx, fy, fz := float64(x), float64(y), float64(z)
	if SlabIsUpper(meta) {
		return SlabBoundingBox{fx, fy + 0.5, fz, fx + 1, fy + 1, fz + 1}
	}
	return SlabBoundingBox{fx, fy, fz, fx + 1, fy + 0.5, fz + 1}
}

type DoubleSlabBlock struct {
	SolidBase
	SingleSlabID uint8
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

func (b *DoubleSlabBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(b.SingleSlabID), Meta: 0, Count: 2}}
}

var stoneSlabVariants = [8]string{
	"Stone", "Sandstone", "Wooden", "Cobblestone",
	"Brick", "Stone Brick", "Quartz", "Nether Brick",
}

func init() {

	Registry.Register(newSlab(SLAB, "Slab", DOUBLE_SLAB, stoneSlabVariants))
	Registry.Register(newDoubleSlab(DOUBLE_SLAB, "Double Slab", SLAB))
}
