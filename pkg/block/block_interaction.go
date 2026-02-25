package block

// block_interaction.go — 方块交互框架
// 对应 PHP Block.php 的 place/onBreak/onUpdate/onActivate 等核心交互方法。
// DefaultBlockInteraction 提供与 PHP 一致的默认行为，所有方块 struct 通过嵌入它
// 自动获得默认实现，只需覆盖需要自定义行为的方块。

// BlockUpdateType 定义方块更新类型，对应 PHP Level::BLOCK_UPDATE_* 常量。
const (
	BlockUpdateNormal    = 1
	BlockUpdateRandom    = 2
	BlockUpdateScheduled = 3
	BlockUpdateWeak      = 4
)

// BlockContext 提供方块交互时的上下文信息。
// 在 Go 中我们不像 PHP 那样把 Level/Position 嵌入 Block 对象本身，
// 而是通过 context 参数传递，这更符合 Go 的设计风格。
type BlockContext struct {
	// X, Y, Z 是方块在世界中的坐标
	X, Y, Z int
	// Meta 是方块的元数据 (0-15)
	Meta uint8
	// Face 是操作的面 (0-5)
	Face int
	// 点击偏移 (0.0-1.0)
	ClickX, ClickY, ClickZ float64
}

// DefaultBlockInteraction 提供所有交互方法的默认（无操作）实现。
// 对应 PHP Block.php 中这些方法的默认 return 值。
// 所有方块 struct 嵌入此类型即可自动满足 BlockBehavior 接口中的新增方法。
type DefaultBlockInteraction struct{}

// Place 默认放置行为：允许放置（返回 true）。
// PHP: return $this->getLevel()->setBlock($this, $this, true);
func (d *DefaultBlockInteraction) Place(ctx *BlockContext) bool {
	return true
}

// OnBreak 默认破坏行为：允许破坏。
// PHP: return $this->getLevel()->setBlock($this, new Air(), true);
func (d *DefaultBlockInteraction) OnBreak(ctx *BlockContext, toolType, toolTier int) bool {
	return true
}

// OnUpdate 方块更新回调。返回值含义：
// - false: 无需进一步更新
// - true: 需要调度下一次更新
// PHP: public function onUpdate($type) { }  (默认空)
func (d *DefaultBlockInteraction) OnUpdate(ctx *BlockContext, updateType int) bool {
	return false
}

// OnActivate 右键交互回调。返回 true 表示已处理（消耗交互）。
// PHP: return false; (默认不做任何事)
func (d *DefaultBlockInteraction) OnActivate(ctx *BlockContext, playerID int64) bool {
	return false
}

// CanBeActivated 是否可以被右键交互。
// PHP: return false;
func (d *DefaultBlockInteraction) CanBeActivated() bool {
	return false
}

// IsBreakable 是否可以被破坏（bedrock 等不可破坏方块覆盖此方法）。
// PHP: return true;
func (d *DefaultBlockInteraction) IsBreakable(toolType, toolTier int) bool {
	return true
}

// GetBreakTime 破坏方块所需时间（秒）。
// PHP: base = hardness; canBreak ? base*1.5 : base*5; base /= efficiency
// 这里返回-1表示使用默认计算（调用方根据 hardness 自行计算）。
func (d *DefaultBlockInteraction) GetBreakTime(toolType, toolTier int) float64 {
	return -1 // 表示使用默认计算
}

// TickRate 方块的更新速率（tick 数），0 表示不参与计划更新。
// PHP 中不同方块覆盖此方法，如 Fire=30，Crops=random tick
func (d *DefaultBlockInteraction) TickRate() int {
	return 0
}

// GetFrictionFactor 返回方块的摩擦系数。
// PHP: return 0.6; (冰块覆盖为 0.98)
func (d *DefaultBlockInteraction) GetFrictionFactor() float64 {
	return 0.6
}

// HasEntityCollision 是否有实体碰撞特殊处理（如蜘蛛网减速、仙人掌伤害）。
// PHP: return false;
func (d *DefaultBlockInteraction) HasEntityCollision() bool {
	return false
}

// OnEntityCollide 实体碰撞回调（如蜘蛛网/仙人掌/岩浆）。
// PHP: public function onEntityCollide(Entity $entity) { }
func (d *DefaultBlockInteraction) OnEntityCollide(ctx *BlockContext, entityID int64) {
	// 默认无操作
}

// GetBurnChance 火焰蔓延概率。0=不可燃。
// PHP: return 0;
func (d *DefaultBlockInteraction) GetBurnChance() int {
	return 0
}

// GetBurnAbility 燃烧能力。0=不可被火烧毁。
// PHP: return 0;
func (d *DefaultBlockInteraction) GetBurnAbility() int {
	return 0
}

// CanPassThrough 是否可穿过（如空气、花）。
// PHP: return false;
func (d *DefaultBlockInteraction) CanPassThrough() bool {
	return false
}

// IsPowerSource 是否为红石信号源。
// PHP: return false;
func (d *DefaultBlockInteraction) IsPowerSource() bool {
	return false
}

// GetStrongPower 返回强红石信号强度 (0-15)。
// PHP: return 0;
func (d *DefaultBlockInteraction) GetStrongPower(face int) int {
	return 0
}

// GetWeakPower 返回弱红石信号强度 (0-15)。
// PHP: return 0;
func (d *DefaultBlockInteraction) GetWeakPower(face int) int {
	return 0
}
