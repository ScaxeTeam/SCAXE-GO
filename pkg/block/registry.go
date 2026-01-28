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
	r.behaviors[GRASS] = &grassBlock{}
	r.behaviors[DIRT] = &dirtBlock{}
	r.behaviors[COBBLESTONE] = &simpleBlock{id: COBBLESTONE, name: "Cobblestone", hardness: 2.0}
	r.behaviors[PLANKS] = &simpleBlock{id: PLANKS, name: "Planks", hardness: 2.0, blastResistance: 15}
	r.behaviors[BEDROCK] = &simpleBlock{id: BEDROCK, name: "Bedrock", hardness: -1, blastResistance: 18000000}
	r.behaviors[WATER] = &liquidBlock{id: WATER, name: "Water", lightFilter: 2}
	r.behaviors[STILL_WATER] = &liquidBlock{id: STILL_WATER, name: "Still Water", lightFilter: 2}
	r.behaviors[LAVA] = &lavaBlock{id: LAVA, name: "Lava"}
	r.behaviors[STILL_LAVA] = &lavaBlock{id: STILL_LAVA, name: "Still Lava"}
	r.behaviors[SAND] = &simpleBlock{id: SAND, name: "Sand", hardness: 0.5}
	r.behaviors[GRAVEL] = &simpleBlock{id: GRAVEL, name: "Gravel", hardness: 0.6}
	r.behaviors[GOLD_ORE] = &simpleBlock{id: GOLD_ORE, name: "Gold Ore", hardness: 3.0}
	r.behaviors[IRON_ORE] = &simpleBlock{id: IRON_ORE, name: "Iron Ore", hardness: 3.0}
	r.behaviors[COAL_ORE] = &simpleBlock{id: COAL_ORE, name: "Coal Ore", hardness: 3.0}
	r.behaviors[WOOD] = &simpleBlock{id: WOOD, name: "Wood", hardness: 2.0}
	r.behaviors[LEAVES] = &leavesBlock{id: LEAVES, name: "Leaves"}
	r.behaviors[GLASS] = &transparentBlock{id: GLASS, name: "Glass", hardness: 0.3}
	r.behaviors[OBSIDIAN] = &simpleBlock{id: OBSIDIAN, name: "Obsidian", hardness: 50.0, blastResistance: 6000}
	r.behaviors[TORCH] = &torchBlock{}
	r.behaviors[GLOWSTONE_BLOCK] = &simpleBlock{id: GLOWSTONE_BLOCK, name: "Glowstone", hardness: 0.3, lightLevel: 15}
	r.behaviors[DIAMOND_ORE] = &simpleBlock{id: DIAMOND_ORE, name: "Diamond Ore", hardness: 3.0}
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
