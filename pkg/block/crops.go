package block

// ── cropBlock (base for wheat, carrot, potato, beetroot) ────────
// All crops: meta 0-7 = growth stage. Placed on FARMLAND only.
// Not solid, transparent, no tool requirement.

type cropBlock struct {
	DefaultBlockInteraction
	id       uint8
	name     string
	seedItem int // item ID dropped as seeds
	cropItem int // item ID dropped when mature (meta >= 7)
}

func (b *cropBlock) GetID() uint8                { return b.id }
func (b *cropBlock) GetName() string             { return b.name }
func (b *cropBlock) GetHardness() float64        { return 0 }
func (b *cropBlock) GetBlastResistance() float64 { return 0 }
func (b *cropBlock) GetLightLevel() uint8        { return 0 }
func (b *cropBlock) GetLightFilter() uint8       { return 0 }
func (b *cropBlock) IsSolid() bool               { return false }
func (b *cropBlock) IsTransparent() bool         { return true }
func (b *cropBlock) CanBePlaced() bool           { return true }
func (b *cropBlock) CanBeReplaced() bool         { return false }
func (b *cropBlock) GetToolType() int            { return ToolTypeNone }
func (b *cropBlock) GetToolTier() int            { return 0 }
func (b *cropBlock) GetDrops(toolType, toolTier int) []Drop {
	// Mature (meta 7) drops crop item + seeds; immature only seeds
	// This is the default behavior, actual meta check done at runtime
	return []Drop{{ID: b.seedItem, Meta: 0, Count: 1}}
}

// ── Wheat (ID 59) ───────────────────────────────────────────────
// Seed item: 295 (WHEAT_SEEDS). Crop item: 296 (WHEAT).

func newWheatBlock() *cropBlock {
	return &cropBlock{
		id:       WHEAT_BLOCK,
		name:     "Wheat Block",
		seedItem: 295, // WHEAT_SEEDS
		cropItem: 296, // WHEAT
	}
}

// ── Carrot (ID 141) ─────────────────────────────────────────────
// Drops: carrot (item 391)

func newCarrotBlock() *cropBlock {
	return &cropBlock{
		id:       CARROT_BLOCK,
		name:     "Carrot Block",
		seedItem: 391, // CARROT (also seed)
		cropItem: 391,
	}
}

// ── Potato (ID 142) ─────────────────────────────────────────────
// Drops: potato (item 392)

func newPotatoBlock() *cropBlock {
	return &cropBlock{
		id:       POTATO_BLOCK,
		name:     "Potato Block",
		seedItem: 392, // POTATO (also seed)
		cropItem: 392,
	}
}

// ── Beetroot (ID 244) ───────────────────────────────────────────
// Seed item: 458 (BEETROOT_SEEDS). Crop item: 457 (BEETROOT).

func newBeetrootBlock() *cropBlock {
	return &cropBlock{
		id:       BEETROOT_BLOCK,
		name:     "Beetroot Block",
		seedItem: 458, // BEETROOT_SEEDS
		cropItem: 457, // BEETROOT
	}
}

// ── Sugarcane (ID 83) ───────────────────────────────────────────
// Grows up to 3 blocks tall. Placed on grass/dirt/sand near water.
// Drops: sugarcane item (338).

type sugarcaneBlock struct{ DefaultBlockInteraction }

func (b *sugarcaneBlock) GetID() uint8                { return SUGARCANE_BLOCK }
func (b *sugarcaneBlock) GetName() string             { return "Sugarcane" }
func (b *sugarcaneBlock) GetHardness() float64        { return 0 }
func (b *sugarcaneBlock) GetBlastResistance() float64 { return 0 }
func (b *sugarcaneBlock) GetLightLevel() uint8        { return 0 }
func (b *sugarcaneBlock) GetLightFilter() uint8       { return 0 }
func (b *sugarcaneBlock) IsSolid() bool               { return false }
func (b *sugarcaneBlock) IsTransparent() bool         { return true }
func (b *sugarcaneBlock) CanBePlaced() bool           { return true }
func (b *sugarcaneBlock) CanBeReplaced() bool         { return false }
func (b *sugarcaneBlock) GetToolType() int            { return ToolTypeNone }
func (b *sugarcaneBlock) GetToolTier() int            { return 0 }
func (b *sugarcaneBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 338, Meta: 0, Count: 1}} // SUGARCANE item
}

// ── Cactus (ID 81) ──────────────────────────────────────────────
// Grows up to 3 blocks tall. Placed on sand.

type cactusBlock struct{ DefaultBlockInteraction }

func (b *cactusBlock) GetID() uint8                { return CACTUS }
func (b *cactusBlock) GetName() string             { return "Cactus" }
func (b *cactusBlock) GetHardness() float64        { return 0.4 }
func (b *cactusBlock) GetBlastResistance() float64 { return 2.0 }
func (b *cactusBlock) GetLightLevel() uint8        { return 0 }
func (b *cactusBlock) GetLightFilter() uint8       { return 0 }
func (b *cactusBlock) IsSolid() bool               { return true }
func (b *cactusBlock) IsTransparent() bool         { return true }
func (b *cactusBlock) CanBePlaced() bool           { return true }
func (b *cactusBlock) CanBeReplaced() bool         { return false }
func (b *cactusBlock) GetToolType() int            { return ToolTypeNone }
func (b *cactusBlock) GetToolTier() int            { return 0 }
func (b *cactusBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(CACTUS), Meta: 0, Count: 1}}
}

// ── Pumpkin Stem (ID 104) ───────────────────────────────────────

type pumpkinStemBlock struct{ DefaultBlockInteraction }

func (b *pumpkinStemBlock) GetID() uint8                { return PUMPKIN_STEM }
func (b *pumpkinStemBlock) GetName() string             { return "Pumpkin Stem" }
func (b *pumpkinStemBlock) GetHardness() float64        { return 0 }
func (b *pumpkinStemBlock) GetBlastResistance() float64 { return 0 }
func (b *pumpkinStemBlock) GetLightLevel() uint8        { return 0 }
func (b *pumpkinStemBlock) GetLightFilter() uint8       { return 0 }
func (b *pumpkinStemBlock) IsSolid() bool               { return false }
func (b *pumpkinStemBlock) IsTransparent() bool         { return true }
func (b *pumpkinStemBlock) CanBePlaced() bool           { return true }
func (b *pumpkinStemBlock) CanBeReplaced() bool         { return false }
func (b *pumpkinStemBlock) GetToolType() int            { return ToolTypeNone }
func (b *pumpkinStemBlock) GetToolTier() int            { return 0 }
func (b *pumpkinStemBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 361, Meta: 0, Count: 1}} // PUMPKIN_SEEDS
}

// ── Melon Stem (ID 105) ────────────────────────────────────────

type melonStemBlock struct{ DefaultBlockInteraction }

func (b *melonStemBlock) GetID() uint8                { return MELON_STEM }
func (b *melonStemBlock) GetName() string             { return "Melon Stem" }
func (b *melonStemBlock) GetHardness() float64        { return 0 }
func (b *melonStemBlock) GetBlastResistance() float64 { return 0 }
func (b *melonStemBlock) GetLightLevel() uint8        { return 0 }
func (b *melonStemBlock) GetLightFilter() uint8       { return 0 }
func (b *melonStemBlock) IsSolid() bool               { return false }
func (b *melonStemBlock) IsTransparent() bool         { return true }
func (b *melonStemBlock) CanBePlaced() bool           { return true }
func (b *melonStemBlock) CanBeReplaced() bool         { return false }
func (b *melonStemBlock) GetToolType() int            { return ToolTypeNone }
func (b *melonStemBlock) GetToolTier() int            { return 0 }
func (b *melonStemBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 362, Meta: 0, Count: 1}} // MELON_SEEDS
}

// ── Registration ────────────────────────────────────────────────

func init() {
	Registry.Register(newWheatBlock())
	Registry.Register(newCarrotBlock())
	Registry.Register(newPotatoBlock())
	Registry.Register(newBeetrootBlock())
	Registry.Register(&sugarcaneBlock{})
	Registry.Register(&cactusBlock{})
	Registry.Register(&pumpkinStemBlock{})
	Registry.Register(&melonStemBlock{})
}
