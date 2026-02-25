package block

// DaylightDetectorBlock 日光传感器
// MCPE 方块 ID 151 (普通) / 178 (反转)
//
// meta 0-15 = 输出红石信号强度 (由光照等级决定)
// 右键切换正常/反转模式

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

// OnActivate 右键切换正常/反转模式
func (b *DaylightDetectorBlock) OnActivate(ctx *BlockContext, playerID int64) bool {
	// 切换 DAYLIGHT_SENSOR <-> DAYLIGHT_SENSOR_INVERTED
	return true
}

func (b *DaylightDetectorBlock) IsInverted() bool {
	return b.inverted
}

// GetFuelTime 日光传感器可作燃料 (300 tick)
func (b *DaylightDetectorBlock) GetFuelTime() int {
	return 300
}

func (b *DaylightDetectorBlock) GetDrops(toolType, toolTier int) []Drop {
	// 不论正反模式，掉落普通日光传感器
	return []Drop{{ID: int(DAYLIGHT_SENSOR), Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(NewDaylightDetectorBlock())
	Registry.Register(NewInvertedDaylightDetectorBlock())
}
