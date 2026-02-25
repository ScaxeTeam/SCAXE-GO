package block
type LiquidType uint8

const (
	LiquidTypeWater LiquidType = iota
	LiquidTypeLava
)
type LiquidConfig struct {
	Type              LiquidType
	FlowingID         uint8
	StillID           uint8
	TickRate          int
	FlowDecayPerBlock int
	LightLevel        uint8
	LightFilter       uint8
	InfiniteSource    bool
}
var WaterConfig = LiquidConfig{
	Type:              LiquidTypeWater,
	FlowingID:         WATER,
	StillID:           STILL_WATER,
	TickRate:          5,
	FlowDecayPerBlock: 1,
	LightLevel:        0,
	LightFilter:       2,
	InfiniteSource:    true,
}
var LavaConfig = LiquidConfig{
	Type:              LiquidTypeLava,
	FlowingID:         LAVA,
	StillID:           STILL_LAVA,
	TickRate:          30,
	FlowDecayPerBlock: 2,
	LightLevel:        15,
	LightFilter:       0,
	InfiniteSource:    false,
}
func LiquidIsSource(meta uint8) bool {
	return meta == 0
}
func LiquidIsFalling(meta uint8) bool {
	return meta >= 8
}
func LiquidGetDecay(meta uint8) int {
	if meta >= 8 {
		return 0
	}
	return int(meta)
}
func LiquidGetFluidHeight(meta uint8) float64 {
	d := int(meta)
	if d >= 8 {
		d = 0
	}
	return float64(d+1) / 9.0
}
type FlowCheckResult int

const (
	FlowBlocked FlowCheckResult = -1
	FlowCanFlow FlowCheckResult = 0
	FlowCanDown FlowCheckResult = 1
)
func LiquidFlowDecay(blockID, liquidFlowingID, liquidStillID, blockMeta uint8) int {
	if blockID != liquidFlowingID && blockID != liquidStillID {
		return -1
	}
	return int(blockMeta)
}
func LiquidEffectiveFlowDecay(blockID, liquidFlowingID, liquidStillID, blockMeta uint8) int {
	if blockID != liquidFlowingID && blockID != liquidStillID {
		return -1
	}
	decay := int(blockMeta)
	if decay >= 8 {
		decay = 0
	}
	return decay
}
type LiquidFlowVector struct {
	X, Y, Z float64
}
type FlowDirection int

const (
	FlowNegX FlowDirection = 0
	FlowPosX FlowDirection = 1
	FlowNegZ FlowDirection = 2
	FlowPosZ FlowDirection = 3
)
var FlowDirectionOffset = [4][2]int{
	{-1, 0},
	{1, 0},
	{0, -1},
	{0, 1},
}
func OppositeDirection(dir FlowDirection) FlowDirection {
	return dir ^ 1
}
type SmallestFlowDecayResult struct {
	Decay           int
	AdjacentSources int
}
func GetSmallestFlowDecay(blockDecay int, currentDecay int, adjacentSources int) (int, int) {
	if blockDecay < 0 {
		return currentDecay, adjacentSources
	}
	if blockDecay == 0 {
		adjacentSources++
	} else if blockDecay >= 8 {
		blockDecay = 0
	}

	if currentDecay >= 0 && blockDecay >= currentDecay {
		return currentDecay, adjacentSources
	}
	return blockDecay, adjacentSources
}
type BlockChecker interface {
	GetBlockIDMeta(x, y, z int) (id uint8, meta uint8)
	CanFlowInto(x, y, z int, liquidFlowingID, liquidStillID uint8) bool
	CanBeFlowedInto(x, y, z int) bool
}
func CalculateFlowCost(
	checker BlockChecker,
	bx, by, bz int,
	accumulatedCost int,
	maxCost int,
	originOpposite FlowDirection,
	lastOpposite FlowDirection,
	liquidFlowingID, liquidStillID uint8,
	visited map[int64]FlowCheckResult,
) int {
	cost := 1000

	for j := FlowDirection(0); j < 4; j++ {
		if j == originOpposite || j == lastOpposite {
			continue
		}

		nx := bx + FlowDirectionOffset[j][0]
		nz := bz + FlowDirectionOffset[j][1]

		hash := blockHash(nx, by, nz)
		if _, exists := visited[hash]; !exists {
			if !checker.CanFlowInto(nx, by, nz, liquidFlowingID, liquidStillID) {
				visited[hash] = FlowBlocked
			} else if checker.CanBeFlowedInto(nx, by-1, nz) {
				visited[hash] = FlowCanDown
			} else {
				visited[hash] = FlowCanFlow
			}
		}

		status := visited[hash]
		if status == FlowBlocked {
			continue
		} else if status == FlowCanDown {
			return accumulatedCost
		}

		if accumulatedCost >= maxCost {
			continue
		}

		realCost := CalculateFlowCost(checker, nx, by, nz, accumulatedCost+1, maxCost, originOpposite, j^1, liquidFlowingID, liquidStillID, visited)
		if realCost < cost {
			cost = realCost
		}
	}

	return cost
}
func GetOptimalFlowDirections(
	checker BlockChecker,
	x, y, z int,
	flowDecayPerBlock int,
	liquidFlowingID, liquidStillID uint8,
) [4]bool {
	flowCost := [4]int{1000, 1000, 1000, 1000}
	maxCost := 4 / flowDecayPerBlock
	visited := make(map[int64]FlowCheckResult)

	for j := FlowDirection(0); j < 4; j++ {
		nx := x + FlowDirectionOffset[j][0]
		nz := z + FlowDirectionOffset[j][1]

		hash := blockHash(nx, y, nz)
		if !checker.CanFlowInto(nx, y, nz, liquidFlowingID, liquidStillID) {
			visited[hash] = FlowBlocked
		} else if checker.CanBeFlowedInto(nx, y-1, nz) {
			visited[hash] = FlowCanDown
			flowCost[j] = 0
			maxCost = 0
		} else if maxCost > 0 {
			visited[hash] = FlowCanFlow
			flowCost[j] = CalculateFlowCost(checker, nx, y, nz, 1, maxCost, j^1, j^1, liquidFlowingID, liquidStillID, visited)
			if flowCost[j] < maxCost {
				maxCost = flowCost[j]
			}
		}
	}
	minCost := flowCost[0]
	for _, c := range flowCost[1:] {
		if c < minCost {
			minCost = c
		}
	}

	var result [4]bool
	for i := 0; i < 4; i++ {
		result[i] = flowCost[i] == minCost
	}
	return result
}
func blockHash(x, y, z int) int64 {
	return (int64(x) << 32) | (int64(z) & 0xFFFFFFFF) ^ (int64(y) << 48)
}
type LiquidBlock struct {
	TransparentBase
	Config LiquidConfig
}

func newLiquidBlock(id uint8, name string, config LiquidConfig) *LiquidBlock {
	return &LiquidBlock{
		TransparentBase: TransparentBase{
			BlockID:          id,
			BlockName:        name,
			BlockHardness:    100,
			BlockResistance:  500,
			BlockLightLevel:  config.LightLevel,
			BlockLightFilter: config.LightFilter,
		},
		Config: config,
	}
}

func (b *LiquidBlock) IsSolid() bool       { return false }
func (b *LiquidBlock) IsTransparent() bool { return true }
func (b *LiquidBlock) CanBePlaced() bool   { return false }
func (b *LiquidBlock) CanBeReplaced() bool { return true }

func (b *LiquidBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

func init() {
	Registry.Register(newLiquidBlock(WATER, "Water", WaterConfig))
	Registry.Register(newLiquidBlock(STILL_WATER, "Still Water", WaterConfig))
	Registry.Register(newLiquidBlock(LAVA, "Lava", LavaConfig))
	Registry.Register(newLiquidBlock(STILL_LAVA, "Still Lava", LavaConfig))
}
