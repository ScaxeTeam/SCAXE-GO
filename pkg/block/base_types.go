package block
type SolidBase struct {
	DefaultBlockInteraction
	BlockID         uint8
	BlockName       string
	BlockHardness   float64
	BlockResistance float64
	BlockLightLevel uint8
	BlockToolType   int
	BlockToolTier   int
}

func (b *SolidBase) GetID() uint8    { return b.BlockID }
func (b *SolidBase) GetName() string { return b.BlockName }
func (b *SolidBase) GetHardness() float64 {
	return b.BlockHardness
}
func (b *SolidBase) GetBlastResistance() float64 {
	if b.BlockResistance > 0 {
		return b.BlockResistance
	}
	return b.BlockHardness * 5
}
func (b *SolidBase) GetLightLevel() uint8  { return b.BlockLightLevel }
func (b *SolidBase) GetLightFilter() uint8 { return 15 }
func (b *SolidBase) IsSolid() bool         { return true }
func (b *SolidBase) IsTransparent() bool   { return false }
func (b *SolidBase) CanBePlaced() bool     { return true }
func (b *SolidBase) CanBeReplaced() bool   { return false }
func (b *SolidBase) GetToolType() int      { return b.BlockToolType }
func (b *SolidBase) GetToolTier() int      { return b.BlockToolTier }
func (b *SolidBase) GetDrops(toolType, toolTier int) []Drop {
	if b.BlockToolType != ToolTypeNone && (toolType != b.BlockToolType || toolTier < b.BlockToolTier) {
		return nil
	}
	return []Drop{{ID: int(b.BlockID), Meta: 0, Count: 1}}
}
type TransparentBase struct {
	DefaultBlockInteraction
	BlockID          uint8
	BlockName        string
	BlockHardness    float64
	BlockResistance  float64
	BlockLightLevel  uint8
	BlockLightFilter uint8
	BlockToolType    int
	BlockToolTier    int
	BlockCanPlace    bool
}

func (b *TransparentBase) GetID() uint8    { return b.BlockID }
func (b *TransparentBase) GetName() string { return b.BlockName }
func (b *TransparentBase) GetHardness() float64 {
	return b.BlockHardness
}
func (b *TransparentBase) GetBlastResistance() float64 {
	if b.BlockResistance > 0 {
		return b.BlockResistance
	}
	return b.BlockHardness * 5
}
func (b *TransparentBase) GetLightLevel() uint8  { return b.BlockLightLevel }
func (b *TransparentBase) GetLightFilter() uint8 { return b.BlockLightFilter }
func (b *TransparentBase) IsSolid() bool         { return true }
func (b *TransparentBase) IsTransparent() bool   { return true }
func (b *TransparentBase) CanBePlaced() bool {
	if b.BlockCanPlace {
		return true
	}
	return b.BlockID != AIR
}
func (b *TransparentBase) CanBeReplaced() bool { return false }
func (b *TransparentBase) GetToolType() int    { return b.BlockToolType }
func (b *TransparentBase) GetToolTier() int    { return b.BlockToolTier }
func (b *TransparentBase) GetDrops(toolType, toolTier int) []Drop {
	if b.BlockToolType != ToolTypeNone && (toolType != b.BlockToolType || toolTier < b.BlockToolTier) {
		return nil
	}
	return []Drop{{ID: int(b.BlockID), Meta: 0, Count: 1}}
}
type FlowableBase struct {
	DefaultBlockInteraction
	BlockID         uint8
	BlockName       string
	BlockHardness   float64
	BlockLightLevel uint8
	BlockToolType   int
}

func (b *FlowableBase) GetID() uint8                { return b.BlockID }
func (b *FlowableBase) GetName() string             { return b.BlockName }
func (b *FlowableBase) GetHardness() float64        { return b.BlockHardness }
func (b *FlowableBase) GetBlastResistance() float64 { return 0 }
func (b *FlowableBase) GetLightLevel() uint8        { return b.BlockLightLevel }
func (b *FlowableBase) GetLightFilter() uint8       { return 0 }
func (b *FlowableBase) IsSolid() bool               { return false }
func (b *FlowableBase) IsTransparent() bool         { return true }
func (b *FlowableBase) CanBePlaced() bool           { return true }
func (b *FlowableBase) CanBeReplaced() bool         { return false }
func (b *FlowableBase) GetToolType() int            { return b.BlockToolType }
func (b *FlowableBase) GetToolTier() int            { return 0 }
func (b *FlowableBase) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(b.BlockID), Meta: 0, Count: 1}}
}
type FallableBase struct {
	SolidBase
}
func (b *FallableBase) IsFallable() bool { return true }
