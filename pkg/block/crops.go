package block

type cropBlock struct {
	DefaultBlockInteraction
	id       uint8
	name     string
	seedItem int
	cropItem int
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
	return []Drop{{ID: b.seedItem, Meta: 0, Count: 1}}
}

func newWheatBlock() *cropBlock {
	return &cropBlock{
		id:       WHEAT_BLOCK,
		name:     "Wheat Block",
		seedItem: 295,
		cropItem: 296,
	}
}

func newCarrotBlock() *cropBlock {
	return &cropBlock{
		id:       CARROT_BLOCK,
		name:     "Carrot Block",
		seedItem: 391,
		cropItem: 391,
	}
}

func newPotatoBlock() *cropBlock {
	return &cropBlock{
		id:       POTATO_BLOCK,
		name:     "Potato Block",
		seedItem: 392,
		cropItem: 392,
	}
}

func newBeetrootBlock() *cropBlock {
	return &cropBlock{
		id:       BEETROOT_BLOCK,
		name:     "Beetroot Block",
		seedItem: 458,
		cropItem: 457,
	}
}

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
	return []Drop{{ID: 338, Meta: 0, Count: 1}}
}

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
	return []Drop{{ID: 361, Meta: 0, Count: 1}}
}

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
	return []Drop{{ID: 362, Meta: 0, Count: 1}}
}

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
