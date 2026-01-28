package block

const (
	ToolTypeNone    = 0
	ToolTypeSword   = 1
	ToolTypeShovel  = 2
	ToolTypePickaxe = 3
	ToolTypeAxe     = 4
	ToolTypeShears  = 5
)

const (
	TierWooden  = 1
	TierGolden  = 2
	TierStone   = 3
	TierIron    = 4
	TierDiamond = 5
)

type Drop struct {
	ID    int
	Meta  int
	Count int
}

func NewDrop(id, meta, count int) Drop {
	return Drop{ID: id, Meta: meta, Count: count}
}

type BlockDrops struct {
	Drops []Drop
	WrongToolDrops []Drop
	MinTier int
	RequiredToolType int
}

func GetDefaultDrops(blockID, meta int) BlockDrops {
	return BlockDrops{
		Drops:            []Drop{{ID: blockID, Meta: meta, Count: 1}},
		WrongToolDrops:   nil,
		MinTier:          0,
		RequiredToolType: ToolTypeNone,
	}
}

type BlockBreakInfo struct {
	Hardness         float64
	BlastResistance  float64
	RequiredToolType int
	RequiredToolTier int
}

func (b BlockBreakInfo) CanHarvestWith(toolType, toolTier int) bool {
	if b.RequiredToolType == ToolTypeNone {
		return true
	}
	if toolType != b.RequiredToolType {
		return false
	}
	return toolTier >= b.RequiredToolTier
}

func (b BlockBreakInfo) GetBreakTime(toolType, toolTier int, efficiency float64) float64 {
	base := b.Hardness
	if base < 0 {
		return -1
	}

	if b.CanHarvestWith(toolType, toolTier) {
		base *= 1.5
	} else {
		base *= 5.0
	}

	if efficiency > 0 {
		base /= efficiency
	}

	return base
}
