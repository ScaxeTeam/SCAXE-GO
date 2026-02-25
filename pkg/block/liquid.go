package block

// liquid.go — 液体方块基类
// 对应 PHP abstract class Liquid extends Transparent
//
// 液体物理核心:
//   Meta 编码: 0=源头, 1-7=流动衰减(越大越浅), >=8=从上方流下
//   流动逻辑: 优先向下，水平方向用 BFS 找最近可下落点
//   无穷水源: 水方块如果 2+ 相邻源且底部是实心/水源 → 变成源头
//
// 此文件提供数据结构和纯函数，不直接依赖 Level 实例。
// 实际的 tick 调度、方块设置、事件触发由 Level 层驱动。

// LiquidType 液体类型
type LiquidType uint8

const (
	LiquidTypeWater LiquidType = iota
	LiquidTypeLava
)

// LiquidConfig 液体配置（不同类型的参数差异）
type LiquidConfig struct {
	Type              LiquidType
	FlowingID         uint8 // 流动方块ID (WATER / LAVA)
	StillID           uint8 // 静止方块ID (STILL_WATER / STILL_LAVA)
	TickRate          int   // 更新间隔（tick），水=5, 岩浆=30
	FlowDecayPerBlock int   // 每格衰减量，水=1, 岩浆=2
	LightLevel        uint8 // 发光等级，水=0, 岩浆=15
	LightFilter       uint8 // 光过滤，水=2, 岩浆=0
	InfiniteSource    bool  // 是否支持无穷水源（仅水）
}

// WaterConfig 水
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

// LavaConfig 岩浆
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

// ---------- Meta 工具函数 ----------

// LiquidIsSource 是否为液体源头 (meta == 0)
func LiquidIsSource(meta uint8) bool {
	return meta == 0
}

// LiquidIsFalling 是否从上方流下 (meta >= 8)
func LiquidIsFalling(meta uint8) bool {
	return meta >= 8
}

// LiquidGetDecay 获取有效衰减值 (falling 视为 0)
func LiquidGetDecay(meta uint8) int {
	if meta >= 8 {
		return 0
	}
	return int(meta)
}

// LiquidGetFluidHeight 获取液面高度百分比 (0.0~1.0)
// 对应 PHP Liquid::getFluidHeightPercent()
func LiquidGetFluidHeight(meta uint8) float64 {
	d := int(meta)
	if d >= 8 {
		d = 0
	}
	return float64(d+1) / 9.0
}

// ---------- 流动衰减判断 ----------

// FlowCheckResult 流动方向的检查结果
type FlowCheckResult int

const (
	FlowBlocked FlowCheckResult = -1 // 被阻挡
	FlowCanFlow FlowCheckResult = 0  // 可以水平流动
	FlowCanDown FlowCheckResult = 1  // 可以向下流动
)

// LiquidFlowDecay 获取方块的流动衰减
// 对应 PHP Liquid::getFlowDecay()
// 返回 -1 表示不是同类液体（被阻挡）
func LiquidFlowDecay(blockID, liquidFlowingID, liquidStillID, blockMeta uint8) int {
	if blockID != liquidFlowingID && blockID != liquidStillID {
		return -1
	}
	return int(blockMeta)
}

// LiquidEffectiveFlowDecay 获取有效流动衰减（falling 视为 0）
// 对应 PHP Liquid::getEffectiveFlowDecay()
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

// ---------- 流向量计算 ----------

// LiquidFlowVector 液体流向量
type LiquidFlowVector struct {
	X, Y, Z float64
}

// ---------- 最优流动方向（BFS） ----------

// FlowDirection 四个水平方向: 0=-X, 1=+X, 2=-Z, 3=+Z
type FlowDirection int

const (
	FlowNegX FlowDirection = 0
	FlowPosX FlowDirection = 1
	FlowNegZ FlowDirection = 2
	FlowPosZ FlowDirection = 3
)

// FlowDirectionOffset 方向偏移
var FlowDirectionOffset = [4][2]int{
	{-1, 0}, // -X
	{1, 0},  // +X
	{0, -1}, // -Z
	{0, 1},  // +Z
}

// OppositeDirection 对方向
func OppositeDirection(dir FlowDirection) FlowDirection {
	return dir ^ 1 // 0↔1, 2↔3
}

// ---------- Scheduled 更新逻辑（纯数据） ----------

// SmallestFlowDecayResult 计算最小流动衰减结果
type SmallestFlowDecayResult struct {
	Decay           int // 最小衰减值
	AdjacentSources int // 相邻源头数量
}

// GetSmallestFlowDecay 从一个方向获取最小流动衰减
// 对应 PHP Liquid::getSmallestFlowDecay()
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

// ---------- 流动成本搜索接口（BFS） ----------

// BlockChecker 方块检查器接口
// Level 层实现此接口，供液体流动算法查询方块状态
type BlockChecker interface {
	// GetBlockIDMeta 获取指定坐标的方块 ID 和 Meta
	GetBlockIDMeta(x, y, z int) (id uint8, meta uint8)
	// CanFlowInto 判断液体能否流入目标坐标
	CanFlowInto(x, y, z int, liquidFlowingID, liquidStillID uint8) bool
	// CanBeFlowedInto 判断目标坐标的方块能否被液体冲走
	CanBeFlowedInto(x, y, z int) bool
}

// CalculateFlowCost BFS 计算流动成本
// 对应 PHP Liquid::calculateFlowCost()
// 返回从 (bx,by,bz) 到最近可下落点的最短路径长度
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

// GetOptimalFlowDirections 获取最优流动方向
// 对应 PHP Liquid::getOptimalFlowDirections()
// 返回 [4]bool，指示 -X, +X, -Z, +Z 四个方向是否应该流动
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

	// 找最小值
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

// blockHash 坐标哈希（与 Level.blockHash 兼容）
func blockHash(x, y, z int) int64 {
	return (int64(x) << 32) | (int64(z) & 0xFFFFFFFF) ^ (int64(y) << 48)
}

// ---------- LiquidBlock 行为 ----------

// LiquidBlock 液体方块行为
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
	return nil // 液体不掉落
}

// ---------- 注册 ----------

func init() {
	// 流动水/静水
	Registry.Register(newLiquidBlock(WATER, "Water", WaterConfig))
	Registry.Register(newLiquidBlock(STILL_WATER, "Still Water", WaterConfig))
	// 流动岩浆/静止岩浆
	Registry.Register(newLiquidBlock(LAVA, "Lava", LavaConfig))
	Registry.Register(newLiquidBlock(STILL_LAVA, "Still Lava", LavaConfig))
}
