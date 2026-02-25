package item
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
func IsBow(id int) bool {
	return id == BOW
}
func GetBowInfo() *BowInfo {
	return &bowInfo
}
type BowChargeState int

const (
	BowChargeNone   BowChargeState = iota
	BowChargeWeak
	BowChargeMedium
	BowChargeFull
)
const BowForceMultiplier = 3.0
func CalcBowForce(useDurationTicks int) (force float64, state BowChargeState) {
	if useDurationTicks < 0 {
		return 0, BowChargeNone
	}
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
func CalcBowDamage(force float64, isCritical bool) int {
	damage := int(force*2 + 0.5)
	if damage < 1 {
		damage = 1
	}
	if isCritical {
		damage += damage / 4
	}
	return damage
}
func IsBowCritical(force float64) bool {
	return force >= 2.8
}
type ArrowInfo struct {
	ID   int
	Name string
}

var arrowInfo = ArrowInfo{
	ID:   ARROW,
	Name: "Arrow",
}
func IsArrow(id int) bool {
	return id == ARROW
}
func GetArrowInfo() *ArrowInfo {
	return &arrowInfo
}
const ArrowMaxStackSize = 64
const ArrowEntityNetworkID = 80
