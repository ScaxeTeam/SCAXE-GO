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

type EntityDamageByBlockEvent struct {
	*EntityDamageEvent
	BlockID   int
	BlockMeta int
}

var entityDamageByBlockHandlers = NewHandlerList()

func NewEntityDamageByBlockEvent(entityID int64, cause int, damage float64, blockID, blockMeta int) *EntityDamageByBlockEvent {
	return &EntityDamageByBlockEvent{
		EntityDamageEvent: NewEntityDamageEvent(entityID, cause, damage),
		BlockID:           blockID,
		BlockMeta:         blockMeta,
	}
}

func (e *EntityDamageByBlockEvent) GetHandlers() *HandlerList { return entityDamageByBlockHandlers }

type EntityDamageByChildEntityEvent struct {
	*EntityDamageByEntityEvent
	ChildEntityID int64
}

var entityDamageByChildEntityHandlers = NewHandlerList()

func NewEntityDamageByChildEntityEvent(entityID, damagerID, childEntityID int64, cause int, damage float64) *EntityDamageByChildEntityEvent {
	return &EntityDamageByChildEntityEvent{
		EntityDamageByEntityEvent: NewEntityDamageByEntityEvent(entityID, damagerID, cause, damage),
		ChildEntityID:             childEntityID,
	}
}

func (e *EntityDamageByChildEntityEvent) GetHandlers() *HandlerList {
	return entityDamageByChildEntityHandlers
}

type EntityCombustEvent struct {
	*EntityEvent
	Duration int
}

var entityCombustHandlers = NewHandlerList()

func NewEntityCombustEvent(entityID int64, duration int) *EntityCombustEvent {
	return &EntityCombustEvent{
		EntityEvent: NewEntityEvent("EntityCombustEvent", entityID),
		Duration:    duration,
	}
}

func (e *EntityCombustEvent) GetHandlers() *HandlerList { return entityCombustHandlers }
func (e *EntityCombustEvent) SetDuration(d int)         { e.Duration = d }

type EntityCombustByBlockEvent struct {
	*EntityCombustEvent
	BlockID int
}

var entityCombustByBlockHandlers = NewHandlerList()

func NewEntityCombustByBlockEvent(entityID int64, duration int, blockID int) *EntityCombustByBlockEvent {
	return &EntityCombustByBlockEvent{
		EntityCombustEvent: NewEntityCombustEvent(entityID, duration),
		BlockID:            blockID,
	}
}

func (e *EntityCombustByBlockEvent) GetHandlers() *HandlerList { return entityCombustByBlockHandlers }

type EntityCombustByEntityEvent struct {
	*EntityCombustEvent
	CombusterID int64
}

var entityCombustByEntityHandlers = NewHandlerList()

func NewEntityCombustByEntityEvent(entityID int64, duration int, combusterID int64) *EntityCombustByEntityEvent {
	return &EntityCombustByEntityEvent{
		EntityCombustEvent: NewEntityCombustEvent(entityID, duration),
		CombusterID:        combusterID,
	}
}

func (e *EntityCombustByEntityEvent) GetHandlers() *HandlerList { return entityCombustByEntityHandlers }

type EntityExplodeEvent struct {
	*EntityEvent
	X, Y, Z   float64
	Force     float64
	BlockList [][3]int
	Yield     float64
}

var entityExplodeHandlers = NewHandlerList()

func NewEntityExplodeEvent(entityID int64, x, y, z, force, yield float64) *EntityExplodeEvent {
	return &EntityExplodeEvent{
		EntityEvent: NewEntityEvent("EntityExplodeEvent", entityID),
		X:           x, Y: y, Z: z,
		Force: force,
		Yield: yield,
	}
}

func (e *EntityExplodeEvent) GetHandlers() *HandlerList { return entityExplodeHandlers }
func (e *EntityExplodeEvent) SetYield(y float64)        { e.Yield = y }

type ExplosionPrimeEvent struct {
	*EntityEvent
	Force      float64
	BlockBreak bool
}

var explosionPrimeHandlers = NewHandlerList()

func NewExplosionPrimeEvent(entityID int64, force float64) *ExplosionPrimeEvent {
	return &ExplosionPrimeEvent{
		EntityEvent: NewEntityEvent("ExplosionPrimeEvent", entityID),
		Force:       force,
		BlockBreak:  true,
	}
}

func (e *ExplosionPrimeEvent) GetHandlers() *HandlerList { return explosionPrimeHandlers }
func (e *ExplosionPrimeEvent) SetForce(f float64)        { e.Force = f }

const (
	RegainEating       = 0
	RegainEffect       = 1
	RegainRegeneration = 2
	RegainCustom       = 3
)

type EntityRegainHealthEvent struct {
	*EntityEvent
	Amount float64
	Reason int
}

var entityRegainHealthHandlers = NewHandlerList()

func NewEntityRegainHealthEvent(entityID int64, amount float64, reason int) *EntityRegainHealthEvent {
	return &EntityRegainHealthEvent{
		EntityEvent: NewEntityEvent("EntityRegainHealthEvent", entityID),
		Amount:      amount,
		Reason:      reason,
	}
}

func (e *EntityRegainHealthEvent) GetHandlers() *HandlerList { return entityRegainHealthHandlers }
func (e *EntityRegainHealthEvent) SetAmount(a float64)       { e.Amount = a }

type EntityMotionEvent struct {
	*EntityEvent
	MotionX, MotionY, MotionZ float64
}

var entityMotionHandlers = NewHandlerList()

func NewEntityMotionEvent(entityID int64, mx, my, mz float64) *EntityMotionEvent {
	return &EntityMotionEvent{
		EntityEvent: NewEntityEvent("EntityMotionEvent", entityID),
		MotionX:     mx, MotionY: my, MotionZ: mz,
	}
}

func (e *EntityMotionEvent) GetHandlers() *HandlerList { return entityMotionHandlers }

type EntityTeleportEvent struct {
	*EntityEvent
	FromX, FromY, FromZ float64
	ToX, ToY, ToZ       float64
}

var entityTeleportHandlers = NewHandlerList()

func NewEntityTeleportEvent(entityID int64, fromX, fromY, fromZ, toX, toY, toZ float64) *EntityTeleportEvent {
	return &EntityTeleportEvent{
		EntityEvent: NewEntityEvent("EntityTeleportEvent", entityID),
		FromX:       fromX, FromY: fromY, FromZ: fromZ,
		ToX: toX, ToY: toY, ToZ: toZ,
	}
}

func (e *EntityTeleportEvent) GetHandlers() *HandlerList { return entityTeleportHandlers }

type EntityLevelChangeEvent struct {
	*EntityEvent
	OriginLevelName string
	TargetLevelName string
}

var entityLevelChangeHandlers = NewHandlerList()

func NewEntityLevelChangeEvent(entityID int64, originLevel, targetLevel string) *EntityLevelChangeEvent {
	return &EntityLevelChangeEvent{
		EntityEvent:     NewEntityEvent("EntityLevelChangeEvent", entityID),
		OriginLevelName: originLevel,
		TargetLevelName: targetLevel,
	}
}

func (e *EntityLevelChangeEvent) GetHandlers() *HandlerList { return entityLevelChangeHandlers }

type EntityShootBowEvent struct {
	*EntityEvent
	Force        float64
	ProjectileID int64
}

var entityShootBowHandlers = NewHandlerList()

func NewEntityShootBowEvent(entityID int64, force float64, projectileID int64) *EntityShootBowEvent {
	return &EntityShootBowEvent{
		EntityEvent:  NewEntityEvent("EntityShootBowEvent", entityID),
		Force:        force,
		ProjectileID: projectileID,
	}
}

func (e *EntityShootBowEvent) GetHandlers() *HandlerList { return entityShootBowHandlers }

type EntityArmorChangeEvent struct {
	*EntityEvent
	Slot      int
	OldItemID int
	NewItemID int
}

var entityArmorChangeHandlers = NewHandlerList()

func NewEntityArmorChangeEvent(entityID int64, slot, oldItemID, newItemID int) *EntityArmorChangeEvent {
	return &EntityArmorChangeEvent{
		EntityEvent: NewEntityEvent("EntityArmorChangeEvent", entityID),
		Slot:        slot,
		OldItemID:   oldItemID,
		NewItemID:   newItemID,
	}
}

func (e *EntityArmorChangeEvent) GetHandlers() *HandlerList { return entityArmorChangeHandlers }

type EntityEatEvent struct {
	*EntityEvent
	HealAmount     float64
	SaturationGain float64
}

var entityEatHandlers = NewHandlerList()

func NewEntityEatEvent(entityID int64, healAmount, saturationGain float64) *EntityEatEvent {
	return &EntityEatEvent{
		EntityEvent:    NewEntityEvent("EntityEatEvent", entityID),
		HealAmount:     healAmount,
		SaturationGain: saturationGain,
	}
}

func (e *EntityEatEvent) GetHandlers() *HandlerList { return entityEatHandlers }

type EntityEatBlockEvent struct {
	*EntityEatEvent
	BlockID   int
	BlockMeta int
}

var entityEatBlockHandlers = NewHandlerList()

func NewEntityEatBlockEvent(entityID int64, healAmount, saturationGain float64, blockID, blockMeta int) *EntityEatBlockEvent {
	return &EntityEatBlockEvent{
		EntityEatEvent: NewEntityEatEvent(entityID, healAmount, saturationGain),
		BlockID:        blockID,
		BlockMeta:      blockMeta,
	}
}

func (e *EntityEatBlockEvent) GetHandlers() *HandlerList { return entityEatBlockHandlers }

type EntityEatItemEvent struct {
	*EntityEatEvent
	ItemID   int
	ItemMeta int
}

var entityEatItemHandlers = NewHandlerList()

func NewEntityEatItemEvent(entityID int64, healAmount, saturationGain float64, itemID, itemMeta int) *EntityEatItemEvent {
	return &EntityEatItemEvent{
		EntityEatEvent: NewEntityEatEvent(entityID, healAmount, saturationGain),
		ItemID:         itemID,
		ItemMeta:       itemMeta,
	}
}

func (e *EntityEatItemEvent) GetHandlers() *HandlerList { return entityEatItemHandlers }

type EntityDrinkPotionEvent struct {
	*EntityEvent
	PotionID int
}

var entityDrinkPotionHandlers = NewHandlerList()

func NewEntityDrinkPotionEvent(entityID int64, potionID int) *EntityDrinkPotionEvent {
	return &EntityDrinkPotionEvent{
		EntityEvent: NewEntityEvent("EntityDrinkPotionEvent", entityID),
		PotionID:    potionID,
	}
}

func (e *EntityDrinkPotionEvent) GetHandlers() *HandlerList { return entityDrinkPotionHandlers }

type EntityInventoryChangeEvent struct {
	*EntityEvent
	Slot      int
	OldItemID int
	NewItemID int
}

var entityInventoryChangeHandlers = NewHandlerList()

func NewEntityInventoryChangeEvent(entityID int64, slot, oldItemID, newItemID int) *EntityInventoryChangeEvent {
	return &EntityInventoryChangeEvent{
		EntityEvent: NewEntityEvent("EntityInventoryChangeEvent", entityID),
		Slot:        slot,
		OldItemID:   oldItemID,
		NewItemID:   newItemID,
	}
}

func (e *EntityInventoryChangeEvent) GetHandlers() *HandlerList { return entityInventoryChangeHandlers }

type EntityBlockChangeEvent struct {
	*EntityEvent
	BlockX, BlockY, BlockZ int
	OldBlockID             int
	NewBlockID             int
}

var entityBlockChangeHandlers = NewHandlerList()

func NewEntityBlockChangeEvent(entityID int64, bx, by, bz, oldBlockID, newBlockID int) *EntityBlockChangeEvent {
	return &EntityBlockChangeEvent{
		EntityEvent: NewEntityEvent("EntityBlockChangeEvent", entityID),
		BlockX:      bx, BlockY: by, BlockZ: bz,
		OldBlockID: oldBlockID,
		NewBlockID: newBlockID,
	}
}

func (e *EntityBlockChangeEvent) GetHandlers() *HandlerList { return entityBlockChangeHandlers }

type EntityGenerateEvent struct {
	*EntityEvent
	X, Y, Z float64
	Cause   int
}

const (
	GenerateCauseNatural  = 0
	GenerateCauseSpawner  = 1
	GenerateCauseSpawnEgg = 2
	GenerateCauseCommand  = 3
)

var entityGenerateHandlers = NewHandlerList()

func NewEntityGenerateEvent(entityID int64, x, y, z float64, cause int) *EntityGenerateEvent {
	return &EntityGenerateEvent{
		EntityEvent: NewEntityEvent("EntityGenerateEvent", entityID),
		X:           x, Y: y, Z: z,
		Cause: cause,
	}
}

func (e *EntityGenerateEvent) GetHandlers() *HandlerList { return entityGenerateHandlers }

type ProjectileHitEvent struct {
	*EntityEvent
	HitEntityID int64
}

var projectileHitHandlers = NewHandlerList()

func NewProjectileHitEvent(projectileID int64, hitEntityID int64) *ProjectileHitEvent {
	return &ProjectileHitEvent{
		EntityEvent: NewEntityEvent("ProjectileHitEvent", projectileID),
		HitEntityID: hitEntityID,
	}
}

func (e *ProjectileHitEvent) GetHandlers() *HandlerList { return projectileHitHandlers }

type ProjectileLaunchEvent struct {
	*EntityEvent
	ShooterID int64
}

var projectileLaunchHandlers = NewHandlerList()

func NewProjectileLaunchEvent(projectileID, shooterID int64) *ProjectileLaunchEvent {
	return &ProjectileLaunchEvent{
		EntityEvent: NewEntityEvent("ProjectileLaunchEvent", projectileID),
		ShooterID:   shooterID,
	}
}

func (e *ProjectileLaunchEvent) GetHandlers() *HandlerList { return projectileLaunchHandlers }

type ItemSpawnEvent struct {
	*EntityEvent
}

var itemSpawnHandlers = NewHandlerList()

func NewItemSpawnEvent(entityID int64) *ItemSpawnEvent {
	return &ItemSpawnEvent{EntityEvent: NewEntityEvent("ItemSpawnEvent", entityID)}
}

func (e *ItemSpawnEvent) GetHandlers() *HandlerList { return itemSpawnHandlers }

type ItemDespawnEvent struct {
	*EntityEvent
}

var itemDespawnHandlers = NewHandlerList()

func NewItemDespawnEvent(entityID int64) *ItemDespawnEvent {
	return &ItemDespawnEvent{EntityEvent: NewEntityEvent("ItemDespawnEvent", entityID)}
}

func (e *ItemDespawnEvent) GetHandlers() *HandlerList { return itemDespawnHandlers }

type ItemMergeEvent struct {
	*EntityEvent
	TargetEntityID int64
}

var itemMergeHandlers = NewHandlerList()

func NewItemMergeEvent(entityID, targetEntityID int64) *ItemMergeEvent {
	return &ItemMergeEvent{
		EntityEvent:    NewEntityEvent("ItemMergeEvent", entityID),
		TargetEntityID: targetEntityID,
	}
}

func (e *ItemMergeEvent) GetHandlers() *HandlerList { return itemMergeHandlers }

const (
	CreeperPowerCauseLightning = 0
	CreeperPowerCauseSetOn     = 1
	CreeperPowerCauseSetOff    = 2
)

type CreeperPowerEvent struct {
	*EntityEvent
	Cause       int
	LightningID int64
}

var creeperPowerHandlers = NewHandlerList()

func NewCreeperPowerEvent(entityID int64, cause int, lightningID int64) *CreeperPowerEvent {
	return &CreeperPowerEvent{
		EntityEvent: NewEntityEvent("CreeperPowerEvent", entityID),
		Cause:       cause,
		LightningID: lightningID,
	}
}

func (e *CreeperPowerEvent) GetHandlers() *HandlerList { return creeperPowerHandlers }
