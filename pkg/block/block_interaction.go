package block

const (
	BlockUpdateNormal    = 1
	BlockUpdateRandom    = 2
	BlockUpdateScheduled = 3
	BlockUpdateWeak      = 4
)

type BlockContext struct {
	X, Y, Z int

	Meta uint8

	Face int

	ClickX, ClickY, ClickZ float64
}

type DefaultBlockInteraction struct{}

func (d *DefaultBlockInteraction) Place(ctx *BlockContext) bool {
	return true
}

func (d *DefaultBlockInteraction) OnBreak(ctx *BlockContext, toolType, toolTier int) bool {
	return true
}

func (d *DefaultBlockInteraction) OnUpdate(ctx *BlockContext, updateType int) bool {
	return false
}

func (d *DefaultBlockInteraction) OnActivate(ctx *BlockContext, playerID int64) bool {
	return false
}

func (d *DefaultBlockInteraction) CanBeActivated() bool {
	return false
}

func (d *DefaultBlockInteraction) IsBreakable(toolType, toolTier int) bool {
	return true
}

func (d *DefaultBlockInteraction) GetBreakTime(toolType, toolTier int) float64 {
	return -1
}

func (d *DefaultBlockInteraction) TickRate() int {
	return 0
}

func (d *DefaultBlockInteraction) GetFrictionFactor() float64 {
	return 0.6
}

func (d *DefaultBlockInteraction) HasEntityCollision() bool {
	return false
}

func (d *DefaultBlockInteraction) OnEntityCollide(ctx *BlockContext, entityID int64) {

}

func (d *DefaultBlockInteraction) GetBurnChance() int {
	return 0
}

func (d *DefaultBlockInteraction) GetBurnAbility() int {
	return 0
}

func (d *DefaultBlockInteraction) CanPassThrough() bool {
	return false
}

func (d *DefaultBlockInteraction) IsPowerSource() bool {
	return false
}

func (d *DefaultBlockInteraction) GetStrongPower(face int) int {
	return 0
}

func (d *DefaultBlockInteraction) GetWeakPower(face int) int {
	return 0
}
