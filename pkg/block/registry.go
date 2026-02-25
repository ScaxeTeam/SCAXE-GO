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
	Place(ctx *BlockContext) bool
	OnBreak(ctx *BlockContext, toolType, toolTier int) bool
	OnUpdate(ctx *BlockContext, updateType int) bool
	OnActivate(ctx *BlockContext, playerID int64) bool
	CanBeActivated() bool
	IsBreakable(toolType, toolTier int) bool
	GetBreakTime(toolType, toolTier int) float64
	TickRate() int
	GetFrictionFactor() float64
	HasEntityCollision() bool
	OnEntityCollide(ctx *BlockContext, entityID int64)
	GetBurnChance() int
	GetBurnAbility() int
	CanPassThrough() bool
	IsPowerSource() bool
	GetStrongPower(face int, meta uint8) int
	GetWeakPower(face int, meta uint8) int
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
	r.behaviors[COBBLESTONE] = &simpleBlock{id: COBBLESTONE, name: "Cobblestone", hardness: 2.0}
	r.behaviors[PLANKS] = &simpleBlock{id: PLANKS, name: "Planks", hardness: 2.0, blastResistance: 15}
	r.behaviors[BEDROCK] = &simpleBlock{id: BEDROCK, name: "Bedrock", hardness: -1, blastResistance: 18000000}
	r.behaviors[WOOD] = &simpleBlock{id: WOOD, name: "Wood", hardness: 2.0}
	r.behaviors[LEAVES] = &leavesBlock{id: LEAVES, name: "Leaves"}
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
