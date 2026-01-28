package block

type simpleBlock struct {
	id              uint8
	name            string
	hardness        float64
	blastResistance float64
	lightLevel      uint8
	lightFilter     uint8
	toolType        int
	toolTier        int
}

func (b *simpleBlock) GetID() uint8          { return b.id }
func (b *simpleBlock) GetName() string       { return b.name }
func (b *simpleBlock) GetHardness() float64  { return b.hardness }
func (b *simpleBlock) GetLightLevel() uint8  { return b.lightLevel }
func (b *simpleBlock) GetLightFilter() uint8 { return 15 }
func (b *simpleBlock) IsSolid() bool         { return true }
func (b *simpleBlock) IsTransparent() bool   { return false }
func (b *simpleBlock) CanBePlaced() bool     { return true }
func (b *simpleBlock) CanBeReplaced() bool   { return false }
func (b *simpleBlock) GetToolType() int      { return b.toolType }
func (b *simpleBlock) GetToolTier() int      { return b.toolTier }
func (b *simpleBlock) GetBlastResistance() float64 {
	if b.blastResistance > 0 {
		return b.blastResistance
	}
	return b.hardness * 5
}
func (b *simpleBlock) GetDrops(toolType, toolTier int) []Drop {

	if b.toolType != ToolTypeNone && (toolType != b.toolType || toolTier < b.toolTier) {
		return nil
	}
	return []Drop{{ID: int(b.id), Meta: 0, Count: 1}}
}

type transparentBlock struct {
	id          uint8
	name        string
	hardness    float64
	lightLevel  uint8
	lightFilter uint8
}

func (b *transparentBlock) GetID() uint8                { return b.id }
func (b *transparentBlock) GetName() string             { return b.name }
func (b *transparentBlock) GetHardness() float64        { return b.hardness }
func (b *transparentBlock) GetLightLevel() uint8        { return b.lightLevel }
func (b *transparentBlock) GetLightFilter() uint8       { return b.lightFilter }
func (b *transparentBlock) IsSolid() bool               { return true }
func (b *transparentBlock) IsTransparent() bool         { return true }
func (b *transparentBlock) CanBePlaced() bool           { return true }
func (b *transparentBlock) CanBeReplaced() bool         { return false }
func (b *transparentBlock) GetBlastResistance() float64 { return b.hardness * 5 }
func (b *transparentBlock) GetToolType() int            { return ToolTypeNone }
func (b *transparentBlock) GetToolTier() int            { return 0 }
func (b *transparentBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

type airBlock struct{}

func (b *airBlock) GetID() uint8                           { return AIR }
func (b *airBlock) GetName() string                        { return "Air" }
func (b *airBlock) GetHardness() float64                   { return 0 }
func (b *airBlock) GetBlastResistance() float64            { return 0 }
func (b *airBlock) GetLightLevel() uint8                   { return 0 }
func (b *airBlock) GetLightFilter() uint8                  { return 0 }
func (b *airBlock) IsSolid() bool                          { return false }
func (b *airBlock) IsTransparent() bool                    { return true }
func (b *airBlock) CanBePlaced() bool                      { return false }
func (b *airBlock) CanBeReplaced() bool                    { return true }
func (b *airBlock) GetToolType() int                       { return ToolTypeNone }
func (b *airBlock) GetToolTier() int                       { return 0 }
func (b *airBlock) GetDrops(toolType, toolTier int) []Drop { return nil }

const (
	StoneNormal           = 0
	StoneGranite          = 1
	StonePolishedGranite  = 2
	StoneDiorite          = 3
	StonePolishedDiorite  = 4
	StoneAndesite         = 5
	StonePolishedAndesite = 6
)

type stoneBlock struct{}

func (b *stoneBlock) GetID() uint8                { return STONE }
func (b *stoneBlock) GetName() string             { return "Stone" }
func (b *stoneBlock) GetHardness() float64        { return 1.5 }
func (b *stoneBlock) GetBlastResistance() float64 { return 30 }
func (b *stoneBlock) GetLightLevel() uint8        { return 0 }
func (b *stoneBlock) GetLightFilter() uint8       { return 15 }
func (b *stoneBlock) IsSolid() bool               { return true }
func (b *stoneBlock) IsTransparent() bool         { return false }
func (b *stoneBlock) CanBePlaced() bool           { return true }
func (b *stoneBlock) CanBeReplaced() bool         { return false }
func (b *stoneBlock) GetToolType() int            { return ToolTypePickaxe }
func (b *stoneBlock) GetToolTier() int            { return TierWooden }
func (b *stoneBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe || toolTier < TierWooden {
		return nil
	}

	return []Drop{{ID: COBBLESTONE, Meta: 0, Count: 1}}
}

type grassBlock struct{}

func (b *grassBlock) GetID() uint8                { return GRASS }
func (b *grassBlock) GetName() string             { return "Grass" }
func (b *grassBlock) GetHardness() float64        { return 0.6 }
func (b *grassBlock) GetBlastResistance() float64 { return 3.0 }
func (b *grassBlock) GetLightLevel() uint8        { return 0 }
func (b *grassBlock) GetLightFilter() uint8       { return 15 }
func (b *grassBlock) IsSolid() bool               { return true }
func (b *grassBlock) IsTransparent() bool         { return false }
func (b *grassBlock) CanBePlaced() bool           { return true }
func (b *grassBlock) CanBeReplaced() bool         { return false }
func (b *grassBlock) GetToolType() int            { return ToolTypeShovel }
func (b *grassBlock) GetToolTier() int            { return 0 }
func (b *grassBlock) GetDrops(toolType, toolTier int) []Drop {

	return []Drop{{ID: DIRT, Meta: 0, Count: 1}}
}

type dirtBlock struct{}

func (b *dirtBlock) GetID() uint8                { return DIRT }
func (b *dirtBlock) GetName() string             { return "Dirt" }
func (b *dirtBlock) GetHardness() float64        { return 0.5 }
func (b *dirtBlock) GetBlastResistance() float64 { return 2.5 }
func (b *dirtBlock) GetLightLevel() uint8        { return 0 }
func (b *dirtBlock) GetLightFilter() uint8       { return 15 }
func (b *dirtBlock) IsSolid() bool               { return true }
func (b *dirtBlock) IsTransparent() bool         { return false }
func (b *dirtBlock) CanBePlaced() bool           { return true }
func (b *dirtBlock) CanBeReplaced() bool         { return false }
func (b *dirtBlock) GetToolType() int            { return ToolTypeShovel }
func (b *dirtBlock) GetToolTier() int            { return 0 }
func (b *dirtBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: DIRT, Meta: 0, Count: 1}}
}

type liquidBlock struct {
	id          uint8
	name        string
	lightFilter uint8
}

func (b *liquidBlock) GetID() uint8                           { return b.id }
func (b *liquidBlock) GetName() string                        { return b.name }
func (b *liquidBlock) GetHardness() float64                   { return 100 }
func (b *liquidBlock) GetBlastResistance() float64            { return 500 }
func (b *liquidBlock) GetLightLevel() uint8                   { return 0 }
func (b *liquidBlock) GetLightFilter() uint8                  { return b.lightFilter }
func (b *liquidBlock) IsSolid() bool                          { return false }
func (b *liquidBlock) IsTransparent() bool                    { return true }
func (b *liquidBlock) CanBePlaced() bool                      { return false }
func (b *liquidBlock) CanBeReplaced() bool                    { return true }
func (b *liquidBlock) GetToolType() int                       { return ToolTypeNone }
func (b *liquidBlock) GetToolTier() int                       { return 0 }
func (b *liquidBlock) GetDrops(toolType, toolTier int) []Drop { return nil }

type lavaBlock struct {
	id   uint8
	name string
}

func (b *lavaBlock) GetID() uint8                           { return b.id }
func (b *lavaBlock) GetName() string                        { return b.name }
func (b *lavaBlock) GetHardness() float64                   { return 100 }
func (b *lavaBlock) GetBlastResistance() float64            { return 500 }
func (b *lavaBlock) GetLightLevel() uint8                   { return 15 }
func (b *lavaBlock) GetLightFilter() uint8                  { return 0 }
func (b *lavaBlock) IsSolid() bool                          { return false }
func (b *lavaBlock) IsTransparent() bool                    { return true }
func (b *lavaBlock) CanBePlaced() bool                      { return false }
func (b *lavaBlock) CanBeReplaced() bool                    { return true }
func (b *lavaBlock) GetToolType() int                       { return ToolTypeNone }
func (b *lavaBlock) GetToolTier() int                       { return 0 }
func (b *lavaBlock) GetDrops(toolType, toolTier int) []Drop { return nil }

type leavesBlock struct {
	id   uint8
	name string
}

func (b *leavesBlock) GetID() uint8                { return b.id }
func (b *leavesBlock) GetName() string             { return b.name }
func (b *leavesBlock) GetHardness() float64        { return 0.2 }
func (b *leavesBlock) GetBlastResistance() float64 { return 1.0 }
func (b *leavesBlock) GetLightLevel() uint8        { return 0 }
func (b *leavesBlock) GetLightFilter() uint8       { return 1 }
func (b *leavesBlock) IsSolid() bool               { return false }
func (b *leavesBlock) IsTransparent() bool         { return true }
func (b *leavesBlock) CanBePlaced() bool           { return true }
func (b *leavesBlock) CanBeReplaced() bool         { return false }
func (b *leavesBlock) GetToolType() int            { return ToolTypeShears }
func (b *leavesBlock) GetToolTier() int            { return 0 }
func (b *leavesBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType == ToolTypeShears {
		return []Drop{{ID: int(b.id), Meta: 0, Count: 1}}
	}

	return nil
}

type torchBlock struct{}

func (b *torchBlock) GetID() uint8                { return TORCH }
func (b *torchBlock) GetName() string             { return "Torch" }
func (b *torchBlock) GetHardness() float64        { return 0 }
func (b *torchBlock) GetBlastResistance() float64 { return 0 }
func (b *torchBlock) GetLightLevel() uint8        { return 14 }
func (b *torchBlock) GetLightFilter() uint8       { return 0 }
func (b *torchBlock) IsSolid() bool               { return false }
func (b *torchBlock) IsTransparent() bool         { return true }
func (b *torchBlock) CanBePlaced() bool           { return true }
func (b *torchBlock) CanBeReplaced() bool         { return false }
func (b *torchBlock) GetToolType() int            { return ToolTypeNone }
func (b *torchBlock) GetToolTier() int            { return 0 }
func (b *torchBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: TORCH, Meta: 0, Count: 1}}
}
