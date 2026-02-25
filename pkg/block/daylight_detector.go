package block

type DaylightDetectorBlock struct {
	TransparentBase
	inverted bool
}

func NewDaylightDetectorBlock() *DaylightDetectorBlock {
	return &DaylightDetectorBlock{
		TransparentBase: TransparentBase{
			BlockID:       DAYLIGHT_SENSOR,
			BlockName:     "Daylight Sensor",
			BlockHardness: 0.2,
			BlockToolType: ToolTypeAxe,
		},
		inverted: false,
	}
}

func NewInvertedDaylightDetectorBlock() *DaylightDetectorBlock {
	return &DaylightDetectorBlock{
		TransparentBase: TransparentBase{
			BlockID:       DAYLIGHT_SENSOR_INVERTED,
			BlockName:     "Daylight Sensor Inverted",
			BlockHardness: 0.2,
			BlockToolType: ToolTypeAxe,
		},
		inverted: true,
	}
}

func (b *DaylightDetectorBlock) CanBeActivated() bool {
	return true
}
func (b *DaylightDetectorBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	return true
}

func (b *DaylightDetectorBlock) IsInverted() bool {
	return b.inverted
}
func (b *DaylightDetectorBlock) GetFuelTime() int {
	return 300
}

func (b *DaylightDetectorBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(DAYLIGHT_SENSOR), Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewDaylightDetectorBlock())
	Registry.Register(NewInvertedDaylightDetectorBlock())
}
