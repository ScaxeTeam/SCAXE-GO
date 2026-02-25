package block

// water_lava.go — 水和岩浆的特化逻辑
// 对应 PHP: Water, StillWater, Lava, StillLava
//
// 主要特化:
//   Water: 碰到实体 → 灭火+重置摔落
//   Lava:  checkForHarden → 岩浆碰水产生黑曜石/圆石
//          flowIntoBlock → 流入水中产生石头
//          下界维度 tickRate 降为 5（同水）

// ============ 水岩浆碰撞结果 ============

// HardenResult 岩浆硬化检查结果
type HardenResult uint8

const (
	HardenNone     HardenResult = iota // 无碰撞
	HardenObsidian                     // 岩浆源碰水 → 黑曜石 (ID 49)
	HardenCobble                       // 岩浆流动碰水 → 圆石 (ID 4)
	HardenStone                        // 岩浆流入水中 → 石头 (ID 1)
)

// HardenResultBlockID 硬化结果对应的方块ID
func HardenResultBlockID(r HardenResult) uint8 {
	switch r {
	case HardenObsidian:
		return OBSIDIAN
	case HardenCobble:
		return COBBLESTONE
	case HardenStone:
		return STONE
	default:
		return AIR
	}
}

// ============ Lava checkForHarden ============

// CheckLavaHarden 检查岩浆方块是否应该硬化
// 对应 PHP Lava::checkForHarden()
//
// 规则:
//   - 检查上方(1)和四个侧面(2-5)是否有水，不检查下方(0)
//   - 如果岩浆源头(meta==0)碰到水 → 黑曜石
//   - 如果流动岩浆(meta<=4)碰到水 → 圆石
//
// 参数:
//   lavaMeta: 岩浆方块的 meta
//   hasAdjacentWater: 上方或四个侧面是否有水方块
func CheckLavaHarden(lavaMeta uint8, hasAdjacentWater bool) HardenResult {
	if !hasAdjacentWater {
		return HardenNone
	}

	if lavaMeta == 0 {
		return HardenObsidian
	}
	if lavaMeta <= 4 {
		return HardenCobble
	}
	return HardenNone
}

// IsWaterBlock 判断方块 ID 是否为水
func IsWaterBlock(blockID uint8) bool {
	return blockID == WATER || blockID == STILL_WATER
}

// IsLavaBlock 判断方块 ID 是否为岩浆
func IsLavaBlock(blockID uint8) bool {
	return blockID == LAVA || blockID == STILL_LAVA
}

// IsLiquidBlock 判断方块 ID 是否为液体
func IsLiquidBlock(blockID uint8) bool {
	return IsWaterBlock(blockID) || IsLavaBlock(blockID)
}

// ============ Lava flowIntoBlock 特化 ============

// CheckLavaFlowIntoWater 检查岩浆流入目标方块时是否应产生石头
// 对应 PHP Lava::flowIntoBlock(): 如果目标是水 → 产生石头（而非正常流入）
func CheckLavaFlowIntoWater(targetBlockID uint8) HardenResult {
	if IsWaterBlock(targetBlockID) {
		return HardenStone
	}
	return HardenNone
}

// ============ 检查相邻水的方向列表 ============

// LavaHardenCheckSides 岩浆硬化检查的方向（上方 + 四侧面，不含下方）
// Side 索引: 0=下, 1=上, 2=北(z-), 3=南(z+), 4=西(x-), 5=东(x+)
var LavaHardenCheckSides = [5]int{1, 2, 3, 4, 5}

// CheckAdjacentWater 检查指定坐标周围（上方+四侧面）是否有水
// 由 Level 层调用，传入方块坐标
func CheckAdjacentWater(checker BlockChecker, x, y, z int) bool {
	// 上方
	id, _ := checker.GetBlockIDMeta(x, y+1, z)
	if IsWaterBlock(id) {
		return true
	}
	// 北(z-)
	id, _ = checker.GetBlockIDMeta(x, y, z-1)
	if IsWaterBlock(id) {
		return true
	}
	// 南(z+)
	id, _ = checker.GetBlockIDMeta(x, y, z+1)
	if IsWaterBlock(id) {
		return true
	}
	// 西(x-)
	id, _ = checker.GetBlockIDMeta(x-1, y, z)
	if IsWaterBlock(id) {
		return true
	}
	// 东(x+)
	id, _ = checker.GetBlockIDMeta(x+1, y, z)
	if IsWaterBlock(id) {
		return true
	}
	return false
}

// ============ 水特化：实体碰撞效果 ============

// WaterEntityEffect 水方块对实体的效果
type WaterEntityEffect struct {
	ExtinguishFire    bool // 灭火
	ResetFallDistance bool // 重置摔落距离
}

// GetWaterEntityEffect 获取水对实体的效果
// 对应 PHP Water::onEntityCollide()
func GetWaterEntityEffect() WaterEntityEffect {
	return WaterEntityEffect{
		ExtinguishFire:    true,
		ResetFallDistance: true,
	}
}

// ============ 岩浆特化：实体碰撞效果 ============

// LavaEntityEffect 岩浆方块对实体的效果
type LavaEntityEffect struct {
	Damage            float64 // 伤害值
	HalveFallDistance bool    // 减半摔落距离
	SetOnFireDuration int     // 着火时间（秒）
	ResetFallDistance bool    // 重置摔落距离
}

// GetLavaEntityEffect 获取岩浆对实体的效果
// 对应 PHP Lava::onEntityCollide()
func GetLavaEntityEffect() LavaEntityEffect {
	return LavaEntityEffect{
		Damage:            4,
		HalveFallDistance: true,
		SetOnFireDuration: 15,
		ResetFallDistance: true,
	}
}

// ============ 下界维度 Lava tickRate 特化 ============

// GetLavaTickRate 获取岩浆更新间隔
// 对应 PHP Lava::tickRate(): 下界=5, 其他=30
func GetLavaTickRate(isNether bool) int {
	if isNether {
		return 5
	}
	return 30
}
