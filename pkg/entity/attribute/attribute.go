package attribute

const (
	Health              = "minecraft:health"
	MovementSpeed       = "minecraft:movement_speed"
	Absorption          = "minecraft:absorption"
	KnockbackResistance = "minecraft:knockback_resistance"
	FollowRange         = "minecraft:follow_range"
	AttackDamage        = "minecraft:attack_damage"
)

type Attribute struct {
	ID           uint32
	Name         string
	MinValue     float32
	MaxValue     float32
	DefaultValue float32
	CurrentValue float32
}

func NewAttribute(name string, min, max, def float32) *Attribute {
	return &Attribute{
		Name:         name,
		MinValue:     min,
		MaxValue:     max,
		DefaultValue: def,
		CurrentValue: def,
	}
}

func (a *Attribute) SetValue(val float32) {
	if val < a.MinValue {
		val = a.MinValue
	}
	if val > a.MaxValue {
		val = a.MaxValue
	}
	a.CurrentValue = val
}

func (a *Attribute) Value() float32 {
	return a.CurrentValue
}

type AttributeMap struct {
	attributes map[string]*Attribute
}

func NewAttributeMap() *AttributeMap {
	return &AttributeMap{
		attributes: make(map[string]*Attribute),
	}
}

func (am *AttributeMap) Add(attr *Attribute) {
	am.attributes[attr.Name] = attr
}

func (am *AttributeMap) Get(name string) *Attribute {
	return am.attributes[name]
}

func (am *AttributeMap) GetAll() []*Attribute {
	list := make([]*Attribute, 0, len(am.attributes))
	for _, a := range am.attributes {
		list = append(list, a)
	}
	return list
}

func GetDefaultAttributes() *AttributeMap {
	am := NewAttributeMap()
	am.Add(NewAttribute(Health, 0, 20, 20))
	am.Add(NewAttribute(MovementSpeed, 0, 24791, 0.1))
	am.Add(NewAttribute(KnockbackResistance, 0, 1, 0))
	am.Add(NewAttribute(AttackDamage, 0, 2048, 1))
	am.Add(NewAttribute(Absorption, 0, 4, 0))
	am.Add(NewAttribute(FollowRange, 0, 2048, 16))
	return am
}
