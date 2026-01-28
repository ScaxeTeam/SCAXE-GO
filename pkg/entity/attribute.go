package entity

import (
	"math"
)

const (
	AttributeAbsorption      = 0
	AttributeSaturation      = 1
	AttributeExhaustion      = 2
	AttributeKnockbackResist = 3
	AttributeHealth          = 4
	AttributeMovementSpeed   = 5
	AttributeFollowRange     = 6
	AttributeHunger          = 7
	AttributeFood            = 7
	AttributeAttackDamage    = 8
	AttributeExperienceLevel = 9
	AttributeExperience      = 10
)

type Attribute struct {
	ID           int
	Name         string
	MinValue     float64
	MaxValue     float64
	DefaultValue float64
	CurrentValue float64
	ShouldSend   bool
	Dirty        bool
}

func NewAttribute(id int, name string, minValue, maxValue, defaultValue float64, shouldSend bool) *Attribute {
	return &Attribute{
		ID:           id,
		Name:         name,
		MinValue:     minValue,
		MaxValue:     maxValue,
		DefaultValue: defaultValue,
		CurrentValue: defaultValue,
		ShouldSend:   shouldSend,
		Dirty:        true,
	}
}

func (a *Attribute) SetValue(value float64) *Attribute {
	value = math.Max(a.MinValue, math.Min(a.MaxValue, value))
	if a.CurrentValue != value {
		a.CurrentValue = value
		a.Dirty = true
	}
	return a
}

func (a *Attribute) GetValue() float64 {
	return a.CurrentValue
}

func (a *Attribute) SetMinValue(value float64) *Attribute {
	if a.MinValue != value {
		a.MinValue = value
		a.Dirty = true
	}
	return a
}

func (a *Attribute) SetMaxValue(value float64) *Attribute {
	if a.MaxValue != value {
		a.MaxValue = value
		a.Dirty = true
	}
	return a
}

func (a *Attribute) SetDefaultValue(value float64) *Attribute {
	if a.DefaultValue != value {
		a.DefaultValue = value
		a.Dirty = true
	}
	return a
}

func (a *Attribute) ResetToDefault() *Attribute {
	return a.SetValue(a.DefaultValue)
}

func (a *Attribute) IsDirty() bool {
	return a.ShouldSend && a.Dirty
}

func (a *Attribute) MarkClean() {
	a.Dirty = false
}

func (a *Attribute) Clone() *Attribute {
	return &Attribute{
		ID:           a.ID,
		Name:         a.Name,
		MinValue:     a.MinValue,
		MaxValue:     a.MaxValue,
		DefaultValue: a.DefaultValue,
		CurrentValue: a.CurrentValue,
		ShouldSend:   a.ShouldSend,
		Dirty:        a.Dirty,
	}
}

type AttributeMap struct {
	Attributes map[int]*Attribute
}

func NewAttributeMap() *AttributeMap {
	return &AttributeMap{
		Attributes: make(map[int]*Attribute),
	}
}

func (am *AttributeMap) AddAttribute(attr *Attribute) {
	am.Attributes[attr.ID] = attr
}

func (am *AttributeMap) GetAttribute(id int) *Attribute {
	return am.Attributes[id]
}

func (am *AttributeMap) GetDirtyAttributes() []*Attribute {
	var dirty []*Attribute
	for _, attr := range am.Attributes {
		if attr.IsDirty() {
			dirty = append(dirty, attr)
		}
	}
	return dirty
}

func (am *AttributeMap) MarkAllClean() {
	for _, attr := range am.Attributes {
		attr.MarkClean()
	}
}

const maxFloat = 3.4028234663852886e+38

var DefaultAttributes = map[int]*Attribute{
	AttributeAbsorption:      NewAttribute(AttributeAbsorption, "generic.absorption", 0, maxFloat, 0, true),
	AttributeSaturation:      NewAttribute(AttributeSaturation, "player.saturation", 0, 20, 20, true),
	AttributeExhaustion:      NewAttribute(AttributeExhaustion, "player.exhaustion", 0, 5, 0, true),
	AttributeKnockbackResist: NewAttribute(AttributeKnockbackResist, "generic.knockbackResistance", 0, 1, 0, true),
	AttributeHealth:          NewAttribute(AttributeHealth, "generic.health", 0, 20, 20, true),
	AttributeMovementSpeed:   NewAttribute(AttributeMovementSpeed, "generic.movementSpeed", 0, maxFloat, 0.1, true),
	AttributeFollowRange:     NewAttribute(AttributeFollowRange, "generic.followRange", 0, 2048, 16, false),
	AttributeHunger:          NewAttribute(AttributeHunger, "player.hunger", 0, 20, 20, true),
	AttributeAttackDamage:    NewAttribute(AttributeAttackDamage, "generic.attackDamage", 0, maxFloat, 1, false),
	AttributeExperienceLevel: NewAttribute(AttributeExperienceLevel, "player.level", 0, 2147483647, 0, true),
	AttributeExperience:      NewAttribute(AttributeExperience, "player.experience", 0, 1, 0, true),
}

func GetDefaultAttribute(id int) *Attribute {
	if attr, ok := DefaultAttributes[id]; ok {
		return attr.Clone()
	}
	return nil
}

func InitDefaultAttributes() *AttributeMap {
	am := NewAttributeMap()
	for id := range DefaultAttributes {
		am.AddAttribute(GetDefaultAttribute(id))
	}
	return am
}
