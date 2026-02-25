package block

// glass.go — 玻璃/玻璃板/染色玻璃 方块
// 对应 PHP: Glass, GlassPane, StainedGlass, StainedGlassPane
//
// 共同特点: 透明、不掉落（精准采集才掉）、硬度 0.3
// GlassPane/StainedGlassPane 是薄板型（Thin），碰撞箱像围栏一样连接

// ============ 玻璃（完整方块） ============

// GlassFullBlock 完整玻璃方块（不掉落）
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

// GetDrops 玻璃不掉落（精准采集由 Level 层处理）
func (b *GlassFullBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

// ============ 玻璃板（薄板型） ============

// GlassPaneBlock 薄型玻璃板
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

// IsSolid 玻璃板不是完整实心
func (b *GlassPaneBlock) IsSolid() bool { return false }

// GetDrops 玻璃板不掉落
func (b *GlassPaneBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

// ============ 薄板碰撞箱 ============

// ThinConnections 薄板连接方向
type ThinConnections struct {
	North, South, West, East bool
}

// ThinBoundingBox 薄板碰撞箱
type ThinBoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}

// ThinThickness 薄板厚度
const ThinThickness = 0.125 // 2/16

// CanThinConnect 判断薄板是否可以连接目标方块
// 玻璃板连接: 同类薄板、完整实心方块
func CanThinConnect(targetID uint8, targetIsSolid, targetIsTransparent bool) bool {
	if IsThinBlock(targetID) {
		return true
	}
	return targetIsSolid && !targetIsTransparent
}

// IsThinBlock 判断是否为薄板型方块（玻璃板/铁栏杆）
func IsThinBlock(blockID uint8) bool {
	return blockID == GLASS_PANE || blockID == STAINED_GLASS_PANE || blockID == IRON_BARS
}

// GetThinBoundingBox 计算薄板碰撞箱
// 类似围栏但只有1格高
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

	// 无连接时只是中心柱
	if !conn.North && !conn.South && !conn.West && !conn.East {
		return ThinBoundingBox{fx + inset, fy, fz + inset, fx + 1 - inset, fy + 1, fz + 1 - inset}
	}

	return ThinBoundingBox{fx + w, fy, fz + n, fx + e, fy + 1, fz + s}
}

// ============ 染色玻璃颜色 ============

// DyeColorNames 16 种染色名称 (meta 0-15)
var DyeColorNames = [16]string{
	"White", "Orange", "Magenta", "Light Blue",
	"Yellow", "Lime", "Pink", "Gray",
	"Light Gray", "Cyan", "Purple", "Blue",
	"Brown", "Green", "Red", "Black",
}

// ============ 注册 ============

func init() {
	// 普通玻璃（替换 registry.go 中的旧注册）
	Registry.Register(newGlass(GLASS, "Glass"))
	Registry.Register(newGlass(STAINED_GLASS, "Stained Glass"))

	// 玻璃板
	Registry.Register(newGlassPane(GLASS_PANE, "Glass Pane"))
	Registry.Register(newGlassPane(STAINED_GLASS_PANE, "Stained Glass Pane"))

	// 铁栏杆（行为类似玻璃板但需镐）
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

// IronBarsBlock 铁栏杆（行为同薄板，需镐掉落）
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
