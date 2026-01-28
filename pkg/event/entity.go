package event

type EntityEvent struct {
	*BaseEvent
	EntityID int64
}

func NewEntityEvent(name string, entityID int64) *EntityEvent {
	return &EntityEvent{
		BaseEvent: NewBaseEvent(name),
		EntityID:  entityID,
	}
}

func (e *EntityEvent) GetEntityID() int64 {
	return e.EntityID
}

const (
	ModifierBase       = 0
	ModifierResistance = 1
	ModifierArmor      = 2
	ModifierProtection = 3
	ModifierStrength   = 4
	ModifierWeakness   = 5
	ModifierCritical   = 6
)

const (
	CauseContact         = 0
	CauseEntityAttack    = 1
	CauseProjectile      = 2
	CauseSuffocation     = 3
	CauseFall            = 4
	CauseFire            = 5
	CauseFireTick        = 6
	CauseLava            = 7
	CauseDrowning        = 8
	CauseBlockExplosion  = 9
	CauseEntityExplosion = 10
	CauseVoid            = 11
	CauseSuicide         = 12
	CauseMagic           = 13
	CauseCustom          = 14
	CauseStarvation      = 15
	CauseLightning       = 16
)

type EntityDamageEvent struct {
	*EntityEvent
	cause     int
	modifiers map[int]float64
	originals map[int]float64
	knockBack float64
}

var entityDamageHandlers = NewHandlerList()

func NewEntityDamageEvent(entityID int64, cause int, damage float64) *EntityDamageEvent {
	e := &EntityDamageEvent{
		EntityEvent: NewEntityEvent("EntityDamageEvent", entityID),
		cause:       cause,
		modifiers:   map[int]float64{ModifierBase: damage},
		originals:   map[int]float64{ModifierBase: damage},
		knockBack:   0.4,
	}
	return e
}

func (e *EntityDamageEvent) GetHandlers() *HandlerList {
	return entityDamageHandlers
}

func (e *EntityDamageEvent) GetCause() int {
	return e.cause
}

func (e *EntityDamageEvent) GetDamage(modifierType int) float64 {
	if d, ok := e.modifiers[modifierType]; ok {
		return d
	}
	return 0
}

func (e *EntityDamageEvent) SetDamage(damage float64, modifierType int) {
	e.modifiers[modifierType] = damage
}

func (e *EntityDamageEvent) GetOriginalDamage(modifierType int) float64 {
	if d, ok := e.originals[modifierType]; ok {
		return d
	}
	return 0
}

func (e *EntityDamageEvent) GetFinalDamage() float64 {
	damage := 1.0
	for _, d := range e.modifiers {
		damage *= d
	}
	return damage
}

func (e *EntityDamageEvent) GetKnockBack() float64 {
	return e.knockBack
}

func (e *EntityDamageEvent) SetKnockBack(kb float64) {
	e.knockBack = kb
}

type EntityDamageByEntityEvent struct {
	*EntityDamageEvent
	DamagerID int64
}

var entityDamageByEntityHandlers = NewHandlerList()

func NewEntityDamageByEntityEvent(entityID, damagerID int64, cause int, damage float64) *EntityDamageByEntityEvent {
	return &EntityDamageByEntityEvent{
		EntityDamageEvent: NewEntityDamageEvent(entityID, cause, damage),
		DamagerID:         damagerID,
	}
}

func (e *EntityDamageByEntityEvent) GetHandlers() *HandlerList {
	return entityDamageByEntityHandlers
}

func (e *EntityDamageByEntityEvent) GetDamagerID() int64 {
	return e.DamagerID
}

type EntityDeathEvent struct {
	*EntityEvent
	Drops []interface{}
}

var entityDeathHandlers = NewHandlerList()

func NewEntityDeathEvent(entityID int64, drops []interface{}) *EntityDeathEvent {
	return &EntityDeathEvent{
		EntityEvent: NewEntityEvent("EntityDeathEvent", entityID),
		Drops:       drops,
	}
}

func (e *EntityDeathEvent) GetHandlers() *HandlerList {
	return entityDeathHandlers
}

type EntitySpawnEvent struct {
	*EntityEvent
}

var entitySpawnHandlers = NewHandlerList()

func NewEntitySpawnEvent(entityID int64) *EntitySpawnEvent {
	return &EntitySpawnEvent{
		EntityEvent: NewEntityEvent("EntitySpawnEvent", entityID),
	}
}

func (e *EntitySpawnEvent) GetHandlers() *HandlerList {
	return entitySpawnHandlers
}

type EntityDespawnEvent struct {
	*EntityEvent
}

var entityDespawnHandlers = NewHandlerList()

func NewEntityDespawnEvent(entityID int64) *EntityDespawnEvent {
	return &EntityDespawnEvent{
		EntityEvent: NewEntityEvent("EntityDespawnEvent", entityID),
	}
}

func (e *EntityDespawnEvent) GetHandlers() *HandlerList {
	return entityDespawnHandlers
}
