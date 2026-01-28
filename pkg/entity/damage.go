package entity

type DamageSource interface {
	Cause() int
	GetDamage() float64
}

type EntityDamageSource struct {
	DamageCause int
	Damage      float64
	Attacker    IEntity
}

func NewEntityDamageSource(cause int, damage float64, attacker IEntity) *EntityDamageSource {
	return &EntityDamageSource{
		DamageCause: cause,
		Damage:      damage,
		Attacker:    attacker,
	}
}

func (e *EntityDamageSource) Cause() int {
	return e.DamageCause
}

func (e *EntityDamageSource) GetDamage() float64 {
	return e.Damage
}

func (e *EntityDamageSource) GetAttacker() IEntity {
	return e.Attacker
}

type SimpleDamageSource struct {
	DamageCause int
	Damage      float64
}

func NewSimpleDamageSource(cause int, damage float64) *SimpleDamageSource {
	return &SimpleDamageSource{
		DamageCause: cause,
		Damage:      damage,
	}
}

func (s *SimpleDamageSource) Cause() int {
	return s.DamageCause
}

func (s *SimpleDamageSource) GetDamage() float64 {
	return s.Damage
}

type DamageCalculator struct{}

func NewDamageCalculator() *DamageCalculator {
	return &DamageCalculator{}
}

func (d *DamageCalculator) CalculateDamage(baseDamage float64, armorPoints int, toughness float64) float64 {

	armor := float64(armorPoints)

	reduction := armor * 0.04
	if reduction > 0.8 {
		reduction = 0.8
	}

	finalDamage := baseDamage * (1 - reduction)
	if finalDamage < 0 {
		finalDamage = 0
	}

	return finalDamage
}

func WeaponDamage(itemID int) float64 {
	damageTable := map[int]float64{

		268: 4,
		283: 4,
		272: 5,
		267: 6,
		276: 7,

		271: 3,
		286: 3,
		275: 3,
		258: 5,
		279: 6,

		270: 2,
		285: 2,
		274: 3,
		257: 4,
		278: 5,

		269: 1,
		284: 1,
		273: 2,
		256: 3,
		277: 4,
	}

	if damage, ok := damageTable[itemID]; ok {
		return damage
	}
	return 1
}

func ArmorProtection(itemID int) int {
	protectionTable := map[int]int{

		298: 1,
		299: 3,
		300: 2,
		301: 1,

		302: 2,
		303: 5,
		304: 4,
		305: 1,

		306: 2,
		307: 6,
		308: 5,
		309: 2,

		314: 2,
		315: 5,
		316: 3,
		317: 1,

		310: 3,
		311: 8,
		312: 6,
		313: 3,
	}

	if protection, ok := protectionTable[itemID]; ok {
		return protection
	}
	return 0
}

const (
	DamageCauseLightning    = 16
	DamageCauseFreezing     = 17
	DamageCauseCampfire     = 18
	DamageCauseSonicBoom    = 19
	DamageCauseFlyIntoWall  = 20
	DamageCauseWitherEffect = 21
	DamageCauseThorns       = 22
	DamageCauseAnvil        = 23
	DamageCauseStalactite   = 24
	DamageCauseStalagmite   = 25
)

func DamageCauseName(cause int) string {
	names := map[int]string{
		DamageCauseContact:         "contact",
		DamageCauseEntityAttack:    "entity_attack",
		DamageCauseProjectile:      "projectile",
		DamageCauseSuffocation:     "suffocation",
		DamageCauseFall:            "fall",
		DamageCauseFire:            "fire",
		DamageCauseFireTick:        "fire_tick",
		DamageCauseLava:            "lava",
		DamageCauseDrowning:        "drowning",
		DamageCauseBlockExplosion:  "block_explosion",
		DamageCauseEntityExplosion: "entity_explosion",
		DamageCauseVoid:            "void",
		DamageCauseSuicide:         "suicide",
		DamageCauseMagic:           "magic",
		DamageCauseStarvation:      "starvation",
		DamageCauseCustom:          "custom",
		DamageCauseLightning:       "lightning",
	}

	if name, ok := names[cause]; ok {
		return name
	}
	return "unknown"
}
