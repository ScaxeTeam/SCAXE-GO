package block
type DirtBlock struct {
	SolidBase
}

func NewDirtBlock() *DirtBlock {
	return &DirtBlock{
		SolidBase: SolidBase{
			BlockID:       DIRT,
			BlockName:     "Dirt",
			BlockHardness: 0.5,
			BlockToolType: ToolTypeShovel,
		},
	}
}
func (b *DirtBlock) CanBeActivated() bool {
	return true
}
type OnActivateResult struct {
	Handled      bool
	ReplaceBlock uint8
	UseTool      bool
}
func OnDirtActivate(isHoe bool) OnActivateResult {
	if isHoe {
		return OnActivateResult{Handled: true, ReplaceBlock: FARMLAND, UseTool: true}
	}
	return OnActivateResult{Handled: false}
}
type GrassBlock struct {
	SolidBase
}

func NewGrassBlock() *GrassBlock {
	return &GrassBlock{
		SolidBase: SolidBase{
			BlockID:       GRASS,
			BlockName:     "Grass",
			BlockHardness: 0.6,
			BlockToolType: ToolTypeShovel,
		},
	}
}
func (b *GrassBlock) CanBeActivated() bool {
	return true
}
func (b *GrassBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(DIRT), Meta: 0, Count: 1}}
}
type GrassActivateType uint8

const (
	GrassActivateNone     GrassActivateType = iota
	GrassActivateBoneMeal
	GrassActivateHoe
	GrassActivateShovel
)
func OnGrassActivate(isBoneMeal bool, isHoe bool, isShovel bool, topBlockIsAir bool) OnActivateResult {
	if isBoneMeal {
		return OnActivateResult{Handled: true, UseTool: true}
	}
	if isHoe {
		return OnActivateResult{Handled: true, ReplaceBlock: FARMLAND, UseTool: true}
	}
	if isShovel && topBlockIsAir {
		return OnActivateResult{Handled: true, ReplaceBlock: GRASS_PATH, UseTool: true}
	}
	return OnActivateResult{Handled: false}
}
type GrassRandomTickResult uint8

const (
	GrassTickNoChange GrassRandomTickResult = iota
	GrassTickDie
	GrassTickSpread
)
func CheckGrassRandomTick(lightAbove int, lightFilterAbove int) GrassRandomTickResult {
	if lightAbove < 4 && lightFilterAbove >= 3 {
		return GrassTickDie
	}
	if lightAbove >= 9 {
		return GrassTickSpread
	}
	return GrassTickNoChange
}
func CanGrassSpreadTo(targetID, targetMeta uint8, lightAboveTarget int, filterAboveTarget int, aboveTargetIsAir bool) bool {
	return targetID == DIRT &&
		targetMeta != 1 &&
		lightAboveTarget >= 4 &&
		filterAboveTarget < 3 &&
		aboveTargetIsAir
}
type SandBlock struct {
	FallableBase
}

func NewSandBlock() *SandBlock {
	return &SandBlock{
		FallableBase: FallableBase{
			SolidBase: SolidBase{
				BlockID:       SAND,
				BlockName:     "Sand",
				BlockHardness: 0.5,
				BlockToolType: ToolTypeShovel,
			},
		},
	}
}
type GravelBlock struct {
	FallableBase
}

func NewGravelBlock() *GravelBlock {
	return &GravelBlock{
		FallableBase: FallableBase{
			SolidBase: SolidBase{
				BlockID:       GRAVEL,
				BlockName:     "Gravel",
				BlockHardness: 0.6,
				BlockToolType: ToolTypeShovel,
			},
		},
	}
}
const ItemFlint = 318
func (b *GravelBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(GRAVEL), Meta: 0, Count: 1}}
}
func GetGravelFlintDrop() Drop {
	return Drop{ID: ItemFlint, Meta: 0, Count: 1}
}
const GravelFlintChance = 10

func init() {
	Registry.Register(NewDirtBlock())
	Registry.Register(NewGrassBlock())
	Registry.Register(NewSandBlock())
	Registry.Register(NewGravelBlock())
}
