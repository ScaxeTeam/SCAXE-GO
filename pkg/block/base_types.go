package block

// base_types.go — 方块基类（对应 PHP Block/Solid/Transparent 继承层次）
//
// Go 用组合（嵌入）替代 PHP 的继承。具体方块类型只需嵌入 SolidBase 或 TransparentBase，
// 然后覆盖需要自定义的方法，大幅减少样板代码。
//
// 对应关系:
//   PHP Solid extends Block       → Go SolidBase
//   PHP Transparent extends Block → Go TransparentBase

// ---------- SolidBase（对应 PHP abstract class Solid extends Block） ----------

// SolidBase 实心方块基类
// IsSolid() = true, IsTransparent() = false, LightFilter = 15
// 嵌入到具体方块中，只需覆盖需要自定义的方法
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
	// 需要正确工具才能掉落
	if b.BlockToolType != ToolTypeNone && (toolType != b.BlockToolType || toolTier < b.BlockToolTier) {
		return nil
	}
	return []Drop{{ID: int(b.BlockID), Meta: 0, Count: 1}}
}

// ---------- TransparentBase（对应 PHP abstract class Transparent extends Block） ----------

// TransparentBase 透明方块基类
// IsSolid() = true（物理碰撞仍然存在，如玻璃/栅栏）, IsTransparent() = true, LightFilter = 0
// 注意：PHP 的 Transparent 类 isSolid() 仍然默认为 true（继承 Block 默认值）
type TransparentBase struct {
	DefaultBlockInteraction
	BlockID          uint8
	BlockName        string
	BlockHardness    float64
	BlockResistance  float64
	BlockLightLevel  uint8
	BlockLightFilter uint8 // 默认0（完全透光），可覆盖
	BlockToolType    int
	BlockToolTier    int
	BlockCanPlace    bool // 默认 true
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
	// 默认为 true（零值时也返回 true）
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

// ---------- FlowableBase（对应 PHP abstract class Flowable extends Transparent） ----------

// FlowableBase 可被液体冲走的非实体方块（花/草/作物等）
// IsSolid() = false, IsTransparent() = true, LightFilter = 0, CanBeReplaced = false
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

// ---------- FallableBase（对应 PHP abstract class Fallable extends Solid） ----------

// FallableBase 受重力影响的实心方块（沙子/砂砾等）
// 继承 SolidBase 的所有属性，额外标记为可下落
// 实际下落逻辑由 Level/Physics 系统驱动
type FallableBase struct {
	SolidBase
}

// IsFallable 标记方块受重力影响
func (b *FallableBase) IsFallable() bool { return true }
