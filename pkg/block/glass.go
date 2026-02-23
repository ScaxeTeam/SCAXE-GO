package block

type GlassFullBlock struct {
	TransparentBase
}

func newGlass(blockID uint8, name string) *GlassFullBlock {
	return &GlassFullBlock{
		TransparentBase: TransparentBase{
			BlockID:       blockID,
			BlockName:     name,
			BlockHardness: 0.3,
			BlockCanPlace: true,
		},
	}
}

func (b *GlassFullBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

type GlassPaneBlock struct {
	TransparentBase
}

func newGlassPane(blockID uint8, name string) *GlassPaneBlock {
	return &GlassPaneBlock{
		TransparentBase: TransparentBase{
			BlockID:       blockID,
			BlockName:     name,
			BlockHardness: 0.3,
			BlockCanPlace: true,
		},
	}
}

func (b *GlassPaneBlock) IsSolid() bool { return false }

func (b *GlassPaneBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

type ThinConnections struct {
	North, South, West, East bool
}

type ThinBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}

const ThinThickness = 0.125

func CanThinConnect(targetID uint8, targetIsSolid, targetIsTransparent bool) bool {
	if IsThinBlock(targetID) {
		return true
	}
	return targetIsSolid && !targetIsTransparent
}

func IsThinBlock(blockID uint8) bool {
	return blockID == GLASS_PANE || blockID == STAINED_GLASS_PANE || blockID == IRON_BARS
}

func GetThinBoundingBox(x, y, z int, conn ThinConnections) ThinBoundingBox {
	fx, fy, fz := float64(x), float64(y), float64(z)
	inset := 0.5 - ThinThickness/2

	n := inset
	if conn.North {
		n = 0
	}
	s := 1 - inset
	if conn.South {
		s = 1
	}
	w := inset
	if conn.West {
		w = 0
	}
	e := 1 - inset
	if conn.East {
		e = 1
	}

	if !conn.North && !conn.South && !conn.West && !conn.East {
		return ThinBoundingBox{fx + inset, fy, fz + inset, fx + 1 - inset, fy + 1, fz + 1 - inset}
	}

	return ThinBoundingBox{fx + w, fy, fz + n, fx + e, fy + 1, fz + s}
}

var DyeColorNames = [16]string{
	"White", "Orange", "Magenta", "Light Blue",
	"Yellow", "Lime", "Pink", "Gray",
	"Light Gray", "Cyan", "Purple", "Blue",
	"Brown", "Green", "Red", "Black",
}

func init() {

	Registry.Register(newGlass(GLASS, "Glass"))
	Registry.Register(newGlass(STAINED_GLASS, "Stained Glass"))

	Registry.Register(newGlassPane(GLASS_PANE, "Glass Pane"))
	Registry.Register(newGlassPane(STAINED_GLASS_PANE, "Stained Glass Pane"))

	Registry.Register(&IronBarsBlock{
		TransparentBase: TransparentBase{
			BlockID:       IRON_BARS,
			BlockName:     "Iron Bars",
			BlockHardness: 5,
			BlockToolType: ToolTypePickaxe,
			BlockCanPlace: true,
		},
	})
}

type IronBarsBlock struct {
	TransparentBase
}

func (b *IronBarsBlock) IsSolid() bool { return false }

func (b *IronBarsBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe {
		return nil
	}
	return []Drop{{ID: int(IRON_BARS), Meta: 0, Count: 1}}
}
