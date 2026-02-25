package block
type HardenResult uint8

const (
	HardenNone     HardenResult = iota
	HardenObsidian
	HardenCobble
	HardenStone
)
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
func IsWaterBlock(blockID uint8) bool {
	return blockID == WATER || blockID == STILL_WATER
}
func IsLavaBlock(blockID uint8) bool {
	return blockID == LAVA || blockID == STILL_LAVA
}
func IsLiquidBlock(blockID uint8) bool {
	return IsWaterBlock(blockID) || IsLavaBlock(blockID)
}
func CheckLavaFlowIntoWater(targetBlockID uint8) HardenResult {
	if IsWaterBlock(targetBlockID) {
		return HardenStone
	}
	return HardenNone
}
var LavaHardenCheckSides = [5]int{1, 2, 3, 4, 5}
func CheckAdjacentWater(checker BlockChecker, x, y, z int) bool {
	id, _ := checker.GetBlockIDMeta(x, y+1, z)
	if IsWaterBlock(id) {
		return true
	}
	id, _ = checker.GetBlockIDMeta(x, y, z-1)
	if IsWaterBlock(id) {
		return true
	}
	id, _ = checker.GetBlockIDMeta(x, y, z+1)
	if IsWaterBlock(id) {
		return true
	}
	id, _ = checker.GetBlockIDMeta(x-1, y, z)
	if IsWaterBlock(id) {
		return true
	}
	id, _ = checker.GetBlockIDMeta(x+1, y, z)
	if IsWaterBlock(id) {
		return true
	}
	return false
}
type WaterEntityEffect struct {
	ExtinguishFire    bool
	ResetFallDistance bool
}
func GetWaterEntityEffect() WaterEntityEffect {
	return WaterEntityEffect{
		ExtinguishFire:    true,
		ResetFallDistance: true,
	}
}
type LavaEntityEffect struct {
	Damage            float64
	HalveFallDistance bool
	SetOnFireDuration int
	ResetFallDistance bool
}
func GetLavaEntityEffect() LavaEntityEffect {
	return LavaEntityEffect{
		Damage:            4,
		HalveFallDistance: true,
		SetOnFireDuration: 15,
		ResetFallDistance: true,
	}
}
func GetLavaTickRate(isNether bool) int {
	if isNether {
		return 5
	}
	return 30
}
