package item

// bow.go — 弓 + 箭 物品
// 对应 PHP: item/Bow.php, item/Arrow.php
//
// PHP 中 Bow 继承 Tool（耐久385），Arrow 继承 Item（可堆叠）。
// 弓的射击逻辑在 Player::releaseUseItem() 中（蓄力→生成 Arrow 实体），
// 此处只定义物品层的属性和便捷查询。
// Arrow 实体的投射物逻辑属于阶段4（entity/arrow.go）。

// ============ 弓 ============

// BowInfo 弓的属性
type BowInfo struct {
	ID         int
	Name       string
	Durability int
}

var bowInfo = BowInfo{
	ID:         BOW,
	Name:       "Bow",
	Durability: 385,
}

// IsBow 判断物品是否为弓
func IsBow(id int) bool {
	return id == BOW
}

// GetBowInfo 获取弓的属性
func GetBowInfo() *BowInfo {
	return &bowInfo
}

// ============ 弓蓄力 ============

// BowChargeState 弓蓄力阶段
type BowChargeState int

const (
	BowChargeNone   BowChargeState = iota // 未蓄力
	BowChargeWeak                         // 弱蓄力 (<0.2s)
	BowChargeMedium                       // 中蓄力 (0.2-1.0s)
	BowChargeFull                         // 满蓄力 (>=1.0s)
)

// BowForceMultiplier 最大蓄力倍率（满蓄力时力度=3.0）
const BowForceMultiplier = 3.0

// CalcBowForce 根据蓄力时长（tick）计算弓的力度
// 对应 PHP Player::releaseUseItem() 中的蓄力计算
//
// 参数:
//   - useDurationTicks: 蓄力持续的 tick 数（20 tick = 1秒）
//
// 返回:
//   - force: 射击力度 (0.0 ~ 3.0)
//   - state: 蓄力阶段
func CalcBowForce(useDurationTicks int) (force float64, state BowChargeState) {
	if useDurationTicks < 0 {
		return 0, BowChargeNone
	}

	// 蓄力公式: force = ticks / 20 (最大3.0)
	f := float64(useDurationTicks) / 20.0
	if f > BowForceMultiplier {
		f = BowForceMultiplier
	}

	switch {
	case f < 0.1:
		return f, BowChargeNone
	case f < 0.6:
		return f, BowChargeWeak
	case f < 2.0:
		return f, BowChargeMedium
	default:
		return f, BowChargeFull
	}
}

// CalcBowDamage 根据力度计算箭的伤害
// 基础伤害 = ceil(force * 2)，满蓄力有暴击（+25%）
func CalcBowDamage(force float64, isCritical bool) int {
	damage := int(force*2 + 0.5) // 近似 ceil
	if damage < 1 {
		damage = 1
	}
	if isCritical {
		damage += damage / 4 // +25%
	}
	return damage
}

// IsBowCritical 判断弓射击是否触发暴击（满蓄力时）
func IsBowCritical(force float64) bool {
	return force >= 2.8 // 接近满蓄力时暴击
}

// ============ 箭 ============

// ArrowInfo 箭的属性
type ArrowInfo struct {
	ID   int
	Name string
}

var arrowInfo = ArrowInfo{
	ID:   ARROW,
	Name: "Arrow",
}

// IsArrow 判断物品是否为箭
func IsArrow(id int) bool {
	return id == ARROW
}

// GetArrowInfo 获取箭的属性
func GetArrowInfo() *ArrowInfo {
	return &arrowInfo
}

// ArrowMaxStackSize 箭的最大堆叠数
const ArrowMaxStackSize = 64

// ArrowEntityNetworkID 箭实体的网络ID（用于 AddEntityPacket）
const ArrowEntityNetworkID = 80
