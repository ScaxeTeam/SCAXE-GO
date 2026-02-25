package block

import (
	"fmt"
	"sync"
)

type BlockState struct {
	ID   uint8
	Meta uint8
}

func NewBlockState(id, meta uint8) BlockState {
	return BlockState{ID: id, Meta: meta & 0x0F}
}

func (b BlockState) FullID() int {
	return (int(b.ID) << 4) | int(b.Meta)
}

func (b BlockState) String() string {
	return fmt.Sprintf("Block{ID: %d, Meta: %d}", b.ID, b.Meta)
}

type BlockBehavior interface {
	GetID() uint8

	GetName() string

	GetHardness() float64

	GetBlastResistance() float64

	GetLightLevel() uint8

	GetLightFilter() uint8

	IsSolid() bool

	IsTransparent() bool

	CanBePlaced() bool

	CanBeReplaced() bool

	GetToolType() int

	GetToolTier() int

	GetDrops(toolType, toolTier int) []Drop

	// ── 交互方法 (PHP Block.php 核心回调) ─────────────

	// Place 放置方块。返回 true 表示放置成功。
	Place(ctx *BlockContext) bool

	// OnBreak 方块被破坏时的回调。
	OnBreak(ctx *BlockContext, toolType, toolTier int) bool

	// OnUpdate 方块更新回调（normal/random/scheduled/weak）。
	OnUpdate(ctx *BlockContext, updateType int) bool

	// OnActivate 右键交互。返回 true 表示已处理。
	OnActivate(ctx *BlockContext, playerID int64) bool

	// CanBeActivated 此方块是否支持右键交互。
	CanBeActivated() bool

	// IsBreakable 是否可以被破坏。
	IsBreakable(toolType, toolTier int) bool

	// GetBreakTime 破坏时间（秒），-1 表示使用默认计算。
	GetBreakTime(toolType, toolTier int) float64

	// TickRate 计划更新的间隔 tick 数，0=不参与计划更新。
	TickRate() int

	// GetFrictionFactor 摩擦系数（冰=0.98, 默认=0.6）。
	GetFrictionFactor() float64

	// HasEntityCollision 是否有实体碰撞特殊处理。
	HasEntityCollision() bool

	// OnEntityCollide 实体碰撞回调。
	OnEntityCollide(ctx *BlockContext, entityID int64)

	// GetBurnChance 火焰蔓延概率。
	GetBurnChance() int

	// GetBurnAbility 燃烧能力。
	GetBurnAbility() int

	// CanPassThrough 是否可穿过。
	CanPassThrough() bool

	// IsPowerSource 是否为红石信号源。
	IsPowerSource() bool

	// GetStrongPower 强红石信号强度。
	GetStrongPower(face int) int

	// GetWeakPower 弱红石信号强度。
	GetWeakPower(face int) int
}

var Registry = &blockRegistry{}

type blockRegistry struct {
	mu sync.RWMutex

	behaviors [256]BlockBehavior

	fullList [4096]BlockState

	solid           [256]bool
	transparent     [256]bool
	hardness        [256]float64
	lightLevel      [256]uint8
	lightFilter     [256]uint8
	blastResistance [256]float64

	initialized bool
}

func (r *blockRegistry) Init() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.initialized {
		return
	}

	r.registerVanillaBlocks()

	for id := 0; id < 256; id++ {
		behavior := r.behaviors[id]
		if behavior != nil {

			r.solid[id] = behavior.IsSolid()
			r.transparent[id] = behavior.IsTransparent()
			r.hardness[id] = behavior.GetHardness()
			r.lightLevel[id] = behavior.GetLightLevel()
			r.lightFilter[id] = min(behavior.GetLightFilter()+1, 15)
			r.blastResistance[id] = behavior.GetBlastResistance()
		} else {

			r.solid[id] = true
			r.transparent[id] = false
			r.hardness[id] = 10
			r.lightLevel[id] = 0
			r.lightFilter[id] = 1
			r.blastResistance[id] = 50
		}

		for meta := 0; meta < 16; meta++ {
			r.fullList[(id<<4)|meta] = BlockState{ID: uint8(id), Meta: uint8(meta)}
		}
	}

	r.initialized = true
}

func (r *blockRegistry) Register(behavior BlockBehavior) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := behavior.GetID()
	r.behaviors[id] = behavior

	if r.initialized {
		r.solid[id] = behavior.IsSolid()
		r.transparent[id] = behavior.IsTransparent()
		r.hardness[id] = behavior.GetHardness()
		r.lightLevel[id] = behavior.GetLightLevel()
		r.lightFilter[id] = min(behavior.GetLightFilter()+1, 15)
		r.blastResistance[id] = behavior.GetBlastResistance()
	}
}

func (r *blockRegistry) Get(id, meta uint8) BlockState {
	return r.fullList[(int(id)<<4)|int(meta&0x0F)]
}

func (r *blockRegistry) GetBehavior(id uint8) BlockBehavior {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.behaviors[id]
}

func (r *blockRegistry) IsSolid(id uint8) bool {
	return r.solid[id]
}

func (r *blockRegistry) IsTransparent(id uint8) bool {
	return r.transparent[id]
}

func (r *blockRegistry) GetHardness(id uint8) float64 {
	return r.hardness[id]
}

func (r *blockRegistry) GetLightLevel(id uint8) uint8 {
	return r.lightLevel[id]
}

func (r *blockRegistry) GetLightFilter(id uint8) uint8 {
	return r.lightFilter[id]
}

func (r *blockRegistry) GetBlastResistance(id uint8) float64 {
	return r.blastResistance[id]
}

func (r *blockRegistry) registerVanillaBlocks() {

	r.behaviors[AIR] = &airBlock{}

	r.behaviors[STONE] = &stoneBlock{}
	// GRASS/DIRT 由 natural.go init() 注册
	r.behaviors[COBBLESTONE] = &simpleBlock{id: COBBLESTONE, name: "Cobblestone", hardness: 2.0}
	r.behaviors[PLANKS] = &simpleBlock{id: PLANKS, name: "Planks", hardness: 2.0, blastResistance: 15}
	r.behaviors[BEDROCK] = &simpleBlock{id: BEDROCK, name: "Bedrock", hardness: -1, blastResistance: 18000000}
	// WATER/STILL_WATER/LAVA/STILL_LAVA 由 liquid.go init() 注册
	// SAND/GRAVEL 由 natural.go init() 注册
	// GOLD_ORE/IRON_ORE/COAL_ORE/DIAMOND_ORE 由 ores.go init() 注册
	r.behaviors[WOOD] = &simpleBlock{id: WOOD, name: "Wood", hardness: 2.0}
	r.behaviors[LEAVES] = &leavesBlock{id: LEAVES, name: "Leaves"}
	// GLASS/GLASS_PANE/STAINED_GLASS/IRON_BARS 由 glass.go init() 注册
	r.behaviors[OBSIDIAN] = &simpleBlock{id: OBSIDIAN, name: "Obsidian", hardness: 50.0, blastResistance: 6000}
	r.behaviors[TORCH] = &torchBlock{}
	r.behaviors[GLOWSTONE_BLOCK] = &simpleBlock{id: GLOWSTONE_BLOCK, name: "Glowstone", hardness: 0.3, lightLevel: 15}
	r.behaviors[DIAMOND_BLOCK] = &simpleBlock{id: DIAMOND_BLOCK, name: "Diamond Block", hardness: 5.0}
	r.behaviors[GOLD_BLOCK] = &simpleBlock{id: GOLD_BLOCK, name: "Gold Block", hardness: 3.0}
	r.behaviors[IRON_BLOCK] = &simpleBlock{id: IRON_BLOCK, name: "Iron Block", hardness: 5.0}

}

func min(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}
