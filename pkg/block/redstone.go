package block

type redstoneTorchBlock struct{ DefaultBlockInteraction }

func (b *redstoneTorchBlock) GetID() uint8                { return REDSTONE_TORCH }
func (b *redstoneTorchBlock) GetName() string             { return "Redstone Torch" }
func (b *redstoneTorchBlock) GetHardness() float64        { return 0 }
func (b *redstoneTorchBlock) GetBlastResistance() float64 { return 0 }
func (b *redstoneTorchBlock) GetLightLevel() uint8        { return 7 }
func (b *redstoneTorchBlock) GetLightFilter() uint8       { return 0 }
func (b *redstoneTorchBlock) IsSolid() bool               { return false }
func (b *redstoneTorchBlock) IsTransparent() bool         { return true }
func (b *redstoneTorchBlock) CanBePlaced() bool           { return true }
func (b *redstoneTorchBlock) CanBeReplaced() bool         { return false }
func (b *redstoneTorchBlock) GetToolType() int            { return ToolTypeNone }
func (b *redstoneTorchBlock) GetToolTier() int            { return 0 }
func (b *redstoneTorchBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(REDSTONE_TORCH), Meta: 0, Count: 1}}
}
func (b *redstoneTorchBlock) IsPowerSource() bool { return true }
func (b *redstoneTorchBlock) GetWeakPower(face int, meta uint8) int {
	return 15
}
func (b *redstoneTorchBlock) GetStrongPower(face int, meta uint8) int {
	return 0
}
func (b *redstoneTorchBlock) OnUpdate(ctx *BlockContext, updateType int) bool {
	if updateType == BlockUpdateNormal && ctx.Powered {
		ctx.ScheduleDelay = 2
	}
	if updateType == BlockUpdateScheduled && ctx.Powered {
		ctx.ReplaceBlockID = UNLIT_REDSTONE_TORCH
		ctx.ReplaceBlockMeta = ctx.Meta
	}
	return false
}

type unlitRedstoneTorchBlock struct{ DefaultBlockInteraction }

func (b *unlitRedstoneTorchBlock) GetID() uint8                { return UNLIT_REDSTONE_TORCH }
func (b *unlitRedstoneTorchBlock) GetName() string             { return "Unlit Redstone Torch" }
func (b *unlitRedstoneTorchBlock) GetHardness() float64        { return 0 }
func (b *unlitRedstoneTorchBlock) GetBlastResistance() float64 { return 0 }
func (b *unlitRedstoneTorchBlock) GetLightLevel() uint8        { return 0 }
func (b *unlitRedstoneTorchBlock) GetLightFilter() uint8       { return 0 }
func (b *unlitRedstoneTorchBlock) IsSolid() bool               { return false }
func (b *unlitRedstoneTorchBlock) IsTransparent() bool         { return true }
func (b *unlitRedstoneTorchBlock) CanBePlaced() bool           { return true }
func (b *unlitRedstoneTorchBlock) CanBeReplaced() bool         { return false }
func (b *unlitRedstoneTorchBlock) GetToolType() int            { return ToolTypeNone }
func (b *unlitRedstoneTorchBlock) GetToolTier() int            { return 0 }
func (b *unlitRedstoneTorchBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(REDSTONE_TORCH), Meta: 0, Count: 1}}
}
func (b *unlitRedstoneTorchBlock) OnUpdate(ctx *BlockContext, updateType int) bool {
	if updateType == BlockUpdateNormal && !ctx.Powered {
		ctx.ScheduleDelay = 2
	}
	if updateType == BlockUpdateScheduled && !ctx.Powered {
		ctx.ReplaceBlockID = REDSTONE_TORCH
		ctx.ReplaceBlockMeta = ctx.Meta
	}
	return false
}

type inactiveRedstoneLampBlock struct{ DefaultBlockInteraction }

func (b *inactiveRedstoneLampBlock) GetID() uint8                { return INACTIVE_REDSTONE_LAMP }
func (b *inactiveRedstoneLampBlock) GetName() string             { return "Inactive Redstone Lamp" }
func (b *inactiveRedstoneLampBlock) GetHardness() float64        { return 0.3 }
func (b *inactiveRedstoneLampBlock) GetBlastResistance() float64 { return 1.5 }
func (b *inactiveRedstoneLampBlock) GetLightLevel() uint8        { return 0 }
func (b *inactiveRedstoneLampBlock) GetLightFilter() uint8       { return 15 }
func (b *inactiveRedstoneLampBlock) IsSolid() bool               { return true }
func (b *inactiveRedstoneLampBlock) IsTransparent() bool         { return false }
func (b *inactiveRedstoneLampBlock) CanBePlaced() bool           { return true }
func (b *inactiveRedstoneLampBlock) CanBeReplaced() bool         { return false }
func (b *inactiveRedstoneLampBlock) GetToolType() int            { return ToolTypeNone }
func (b *inactiveRedstoneLampBlock) GetToolTier() int            { return 0 }
func (b *inactiveRedstoneLampBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(INACTIVE_REDSTONE_LAMP), Meta: 0, Count: 1}}
}
func (b *inactiveRedstoneLampBlock) OnUpdate(ctx *BlockContext, updateType int) bool {
	if updateType == BlockUpdateNormal && ctx.Powered {
		ctx.ReplaceBlockID = ACTIVE_REDSTONE_LAMP
		ctx.ReplaceBlockMeta = 0
	}
	return false
}

type activeRedstoneLampBlock struct{ DefaultBlockInteraction }

func (b *activeRedstoneLampBlock) GetID() uint8                { return ACTIVE_REDSTONE_LAMP }
func (b *activeRedstoneLampBlock) GetName() string             { return "Active Redstone Lamp" }
func (b *activeRedstoneLampBlock) GetHardness() float64        { return 0.3 }
func (b *activeRedstoneLampBlock) GetBlastResistance() float64 { return 1.5 }
func (b *activeRedstoneLampBlock) GetLightLevel() uint8        { return 15 }
func (b *activeRedstoneLampBlock) GetLightFilter() uint8       { return 15 }
func (b *activeRedstoneLampBlock) IsSolid() bool               { return true }
func (b *activeRedstoneLampBlock) IsTransparent() bool         { return false }
func (b *activeRedstoneLampBlock) CanBePlaced() bool           { return true }
func (b *activeRedstoneLampBlock) CanBeReplaced() bool         { return false }
func (b *activeRedstoneLampBlock) GetToolType() int            { return ToolTypeNone }
func (b *activeRedstoneLampBlock) GetToolTier() int            { return 0 }
func (b *activeRedstoneLampBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(INACTIVE_REDSTONE_LAMP), Meta: 0, Count: 1}}
}
func (b *activeRedstoneLampBlock) OnUpdate(ctx *BlockContext, updateType int) bool {
	if updateType == BlockUpdateNormal && !ctx.Powered {
		ctx.ScheduleDelay = 4
	}
	if updateType == BlockUpdateScheduled && !ctx.Powered {
		ctx.ReplaceBlockID = INACTIVE_REDSTONE_LAMP
		ctx.ReplaceBlockMeta = 0
	}
	return false
}

type leverBlock struct{ DefaultBlockInteraction }

func (b *leverBlock) GetID() uint8                { return LEVER }
func (b *leverBlock) GetName() string             { return "Lever" }
func (b *leverBlock) GetHardness() float64        { return 0.5 }
func (b *leverBlock) GetBlastResistance() float64 { return 2.5 }
func (b *leverBlock) GetLightLevel() uint8        { return 0 }
func (b *leverBlock) GetLightFilter() uint8       { return 0 }
func (b *leverBlock) IsSolid() bool               { return false }
func (b *leverBlock) IsTransparent() bool         { return true }
func (b *leverBlock) CanBePlaced() bool           { return true }
func (b *leverBlock) CanBeReplaced() bool         { return false }
func (b *leverBlock) GetToolType() int            { return ToolTypeNone }
func (b *leverBlock) GetToolTier() int            { return 0 }
func (b *leverBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(LEVER), Meta: 0, Count: 1}}
}
func (b *leverBlock) CanBeActivated() bool                              { return true }
func (b *leverBlock) OnActivate(ctx *BlockContext, playerID int64) bool { return true }
func (b *leverBlock) IsPowerSource() bool                               { return true }
func (b *leverBlock) GetStrongPower(face int, meta uint8) int {
	if meta&0x08 != 0 {
		return 15
	}
	return 0
}
func (b *leverBlock) GetWeakPower(face int, meta uint8) int {
	if meta&0x08 != 0 {
		return 15
	}
	return 0
}

type stoneButtonBlock struct{ DefaultBlockInteraction }

func (b *stoneButtonBlock) GetID() uint8                { return STONE_BUTTON }
func (b *stoneButtonBlock) GetName() string             { return "Stone Button" }
func (b *stoneButtonBlock) GetHardness() float64        { return 0.5 }
func (b *stoneButtonBlock) GetBlastResistance() float64 { return 2.5 }
func (b *stoneButtonBlock) GetLightLevel() uint8        { return 0 }
func (b *stoneButtonBlock) GetLightFilter() uint8       { return 0 }
func (b *stoneButtonBlock) IsSolid() bool               { return false }
func (b *stoneButtonBlock) IsTransparent() bool         { return true }
func (b *stoneButtonBlock) CanBePlaced() bool           { return true }
func (b *stoneButtonBlock) CanBeReplaced() bool         { return false }
func (b *stoneButtonBlock) GetToolType() int            { return ToolTypeNone }
func (b *stoneButtonBlock) GetToolTier() int            { return 0 }
func (b *stoneButtonBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(STONE_BUTTON), Meta: 0, Count: 1}}
}
func (b *stoneButtonBlock) CanBeActivated() bool                              { return true }
func (b *stoneButtonBlock) OnActivate(ctx *BlockContext, playerID int64) bool { return true }
func (b *stoneButtonBlock) IsPowerSource() bool                               { return true }
func (b *stoneButtonBlock) GetStrongPower(face int, meta uint8) int {
	if meta&0x08 != 0 {
		return 15
	}
	return 0
}
func (b *stoneButtonBlock) GetWeakPower(face int, meta uint8) int {
	if meta&0x08 != 0 {
		return 15
	}
	return 0
}
func (b *stoneButtonBlock) OnUpdate(ctx *BlockContext, updateType int) bool {
	if updateType == BlockUpdateScheduled && ctx.Meta&0x08 != 0 {
		ctx.ReplaceBlockID = STONE_BUTTON
		ctx.ReplaceBlockMeta = ctx.Meta &^ 0x08
	}
	return false
}

type woodenButtonBlock struct{ DefaultBlockInteraction }

func (b *woodenButtonBlock) GetID() uint8                { return WOODEN_BUTTON }
func (b *woodenButtonBlock) GetName() string             { return "Wooden Button" }
func (b *woodenButtonBlock) GetHardness() float64        { return 0.5 }
func (b *woodenButtonBlock) GetBlastResistance() float64 { return 2.5 }
func (b *woodenButtonBlock) GetLightLevel() uint8        { return 0 }
func (b *woodenButtonBlock) GetLightFilter() uint8       { return 0 }
func (b *woodenButtonBlock) IsSolid() bool               { return false }
func (b *woodenButtonBlock) IsTransparent() bool         { return true }
func (b *woodenButtonBlock) CanBePlaced() bool           { return true }
func (b *woodenButtonBlock) CanBeReplaced() bool         { return false }
func (b *woodenButtonBlock) GetToolType() int            { return ToolTypeNone }
func (b *woodenButtonBlock) GetToolTier() int            { return 0 }
func (b *woodenButtonBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(WOODEN_BUTTON), Meta: 0, Count: 1}}
}
func (b *woodenButtonBlock) CanBeActivated() bool                              { return true }
func (b *woodenButtonBlock) OnActivate(ctx *BlockContext, playerID int64) bool { return true }
func (b *woodenButtonBlock) IsPowerSource() bool                               { return true }
func (b *woodenButtonBlock) GetStrongPower(face int, meta uint8) int {
	if meta&0x08 != 0 {
		return 15
	}
	return 0
}
func (b *woodenButtonBlock) GetWeakPower(face int, meta uint8) int {
	if meta&0x08 != 0 {
		return 15
	}
	return 0
}
func (b *woodenButtonBlock) OnUpdate(ctx *BlockContext, updateType int) bool {
	if updateType == BlockUpdateScheduled && ctx.Meta&0x08 != 0 {
		ctx.ReplaceBlockID = WOODEN_BUTTON
		ctx.ReplaceBlockMeta = ctx.Meta &^ 0x08
	}
	return false
}

type redstoneWireBlock struct{ DefaultBlockInteraction }

func (b *redstoneWireBlock) GetID() uint8                { return REDSTONE_WIRE }
func (b *redstoneWireBlock) GetName() string             { return "Redstone Wire" }
func (b *redstoneWireBlock) GetHardness() float64        { return 0 }
func (b *redstoneWireBlock) GetBlastResistance() float64 { return 0 }
func (b *redstoneWireBlock) GetLightLevel() uint8        { return 0 }
func (b *redstoneWireBlock) GetLightFilter() uint8       { return 0 }
func (b *redstoneWireBlock) IsSolid() bool               { return false }
func (b *redstoneWireBlock) IsTransparent() bool         { return true }
func (b *redstoneWireBlock) CanBePlaced() bool           { return true }
func (b *redstoneWireBlock) CanBeReplaced() bool         { return false }
func (b *redstoneWireBlock) GetToolType() int            { return ToolTypeNone }
func (b *redstoneWireBlock) GetToolTier() int            { return 0 }
func (b *redstoneWireBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 331, Meta: 0, Count: 1}}
}
func (b *redstoneWireBlock) IsPowerSource() bool { return true }
func (b *redstoneWireBlock) GetWeakPower(face int, meta uint8) int {
	return int(meta)
}
func (b *redstoneWireBlock) GetStrongPower(face int, meta uint8) int {
	return 0
}

type stonePressurePlateBlock struct{ DefaultBlockInteraction }

func (b *stonePressurePlateBlock) GetID() uint8                { return STONE_PRESSURE_PLATE }
func (b *stonePressurePlateBlock) GetName() string             { return "Stone Pressure Plate" }
func (b *stonePressurePlateBlock) GetHardness() float64        { return 0.5 }
func (b *stonePressurePlateBlock) GetBlastResistance() float64 { return 2.5 }
func (b *stonePressurePlateBlock) GetLightLevel() uint8        { return 0 }
func (b *stonePressurePlateBlock) GetLightFilter() uint8       { return 0 }
func (b *stonePressurePlateBlock) IsSolid() bool               { return false }
func (b *stonePressurePlateBlock) IsTransparent() bool         { return true }
func (b *stonePressurePlateBlock) CanBePlaced() bool           { return true }
func (b *stonePressurePlateBlock) CanBeReplaced() bool         { return false }
func (b *stonePressurePlateBlock) GetToolType() int            { return ToolTypePickaxe }
func (b *stonePressurePlateBlock) GetToolTier() int            { return 0 }
func (b *stonePressurePlateBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(STONE_PRESSURE_PLATE), Meta: 0, Count: 1}}
}
func (b *stonePressurePlateBlock) IsPowerSource() bool { return true }
func (b *stonePressurePlateBlock) GetWeakPower(face int, meta uint8) int {
	if meta&0x01 != 0 {
		return 15
	}
	return 0
}
func (b *stonePressurePlateBlock) GetStrongPower(face int, meta uint8) int {
	if meta&0x01 != 0 && face == 0 {
		return 15
	}
	return 0
}
func (b *stonePressurePlateBlock) OnUpdate(ctx *BlockContext, updateType int) bool {
	if updateType == BlockUpdateScheduled && ctx.Meta&0x01 != 0 {
		ctx.ReplaceBlockID = STONE_PRESSURE_PLATE
		ctx.ReplaceBlockMeta = 0
	}
	return false
}

type woodenPressurePlateBlock struct{ DefaultBlockInteraction }

func (b *woodenPressurePlateBlock) GetID() uint8                { return WOODEN_PRESSURE_PLATE }
func (b *woodenPressurePlateBlock) GetName() string             { return "Wooden Pressure Plate" }
func (b *woodenPressurePlateBlock) GetHardness() float64        { return 0.5 }
func (b *woodenPressurePlateBlock) GetBlastResistance() float64 { return 2.5 }
func (b *woodenPressurePlateBlock) GetLightLevel() uint8        { return 0 }
func (b *woodenPressurePlateBlock) GetLightFilter() uint8       { return 0 }
func (b *woodenPressurePlateBlock) IsSolid() bool               { return false }
func (b *woodenPressurePlateBlock) IsTransparent() bool         { return true }
func (b *woodenPressurePlateBlock) CanBePlaced() bool           { return true }
func (b *woodenPressurePlateBlock) CanBeReplaced() bool         { return false }
func (b *woodenPressurePlateBlock) GetToolType() int            { return ToolTypeAxe }
func (b *woodenPressurePlateBlock) GetToolTier() int            { return 0 }
func (b *woodenPressurePlateBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(WOODEN_PRESSURE_PLATE), Meta: 0, Count: 1}}
}
func (b *woodenPressurePlateBlock) IsPowerSource() bool { return true }
func (b *woodenPressurePlateBlock) GetWeakPower(face int, meta uint8) int {
	if meta&0x01 != 0 {
		return 15
	}
	return 0
}
func (b *woodenPressurePlateBlock) GetStrongPower(face int, meta uint8) int {
	if meta&0x01 != 0 && face == 0 {
		return 15
	}
	return 0
}
func (b *woodenPressurePlateBlock) OnUpdate(ctx *BlockContext, updateType int) bool {
	if updateType == BlockUpdateScheduled && ctx.Meta&0x01 != 0 {
		ctx.ReplaceBlockID = WOODEN_PRESSURE_PLATE
		ctx.ReplaceBlockMeta = 0
	}
	return false
}

type lightWeightedPressurePlateBlock struct{ DefaultBlockInteraction }

func (b *lightWeightedPressurePlateBlock) GetID() uint8                { return LIGHT_WEIGHTED_PRESSURE_PLATE }
func (b *lightWeightedPressurePlateBlock) GetName() string             { return "Light Weighted Pressure Plate" }
func (b *lightWeightedPressurePlateBlock) GetHardness() float64        { return 0.5 }
func (b *lightWeightedPressurePlateBlock) GetBlastResistance() float64 { return 2.5 }
func (b *lightWeightedPressurePlateBlock) GetLightLevel() uint8        { return 0 }
func (b *lightWeightedPressurePlateBlock) GetLightFilter() uint8       { return 0 }
func (b *lightWeightedPressurePlateBlock) IsSolid() bool               { return false }
func (b *lightWeightedPressurePlateBlock) IsTransparent() bool         { return true }
func (b *lightWeightedPressurePlateBlock) CanBePlaced() bool           { return true }
func (b *lightWeightedPressurePlateBlock) CanBeReplaced() bool         { return false }
func (b *lightWeightedPressurePlateBlock) GetToolType() int            { return ToolTypePickaxe }
func (b *lightWeightedPressurePlateBlock) GetToolTier() int            { return 0 }
func (b *lightWeightedPressurePlateBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(LIGHT_WEIGHTED_PRESSURE_PLATE), Meta: 0, Count: 1}}
}
func (b *lightWeightedPressurePlateBlock) IsPowerSource() bool                     { return true }
func (b *lightWeightedPressurePlateBlock) GetWeakPower(face int, meta uint8) int   { return int(meta) }
func (b *lightWeightedPressurePlateBlock) GetStrongPower(face int, meta uint8) int { return 0 }
func (b *lightWeightedPressurePlateBlock) OnUpdate(ctx *BlockContext, updateType int) bool {
	if updateType == BlockUpdateScheduled && ctx.Meta > 0 {
		ctx.ReplaceBlockID = LIGHT_WEIGHTED_PRESSURE_PLATE
		ctx.ReplaceBlockMeta = 0
	}
	return false
}

type heavyWeightedPressurePlateBlock struct{ DefaultBlockInteraction }

func (b *heavyWeightedPressurePlateBlock) GetID() uint8                { return HEAVY_WEIGHTED_PRESSURE_PLATE }
func (b *heavyWeightedPressurePlateBlock) GetName() string             { return "Heavy Weighted Pressure Plate" }
func (b *heavyWeightedPressurePlateBlock) GetHardness() float64        { return 0.5 }
func (b *heavyWeightedPressurePlateBlock) GetBlastResistance() float64 { return 2.5 }
func (b *heavyWeightedPressurePlateBlock) GetLightLevel() uint8        { return 0 }
func (b *heavyWeightedPressurePlateBlock) GetLightFilter() uint8       { return 0 }
func (b *heavyWeightedPressurePlateBlock) IsSolid() bool               { return false }
func (b *heavyWeightedPressurePlateBlock) IsTransparent() bool         { return true }
func (b *heavyWeightedPressurePlateBlock) CanBePlaced() bool           { return true }
func (b *heavyWeightedPressurePlateBlock) CanBeReplaced() bool         { return false }
func (b *heavyWeightedPressurePlateBlock) GetToolType() int            { return ToolTypePickaxe }
func (b *heavyWeightedPressurePlateBlock) GetToolTier() int            { return 0 }
func (b *heavyWeightedPressurePlateBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(HEAVY_WEIGHTED_PRESSURE_PLATE), Meta: 0, Count: 1}}
}
func (b *heavyWeightedPressurePlateBlock) IsPowerSource() bool                     { return true }
func (b *heavyWeightedPressurePlateBlock) GetWeakPower(face int, meta uint8) int   { return int(meta) }
func (b *heavyWeightedPressurePlateBlock) GetStrongPower(face int, meta uint8) int { return 0 }
func (b *heavyWeightedPressurePlateBlock) OnUpdate(ctx *BlockContext, updateType int) bool {
	if updateType == BlockUpdateScheduled && ctx.Meta > 0 {
		ctx.ReplaceBlockID = HEAVY_WEIGHTED_PRESSURE_PLATE
		ctx.ReplaceBlockMeta = 0
	}
	return false
}

func init() {
	Registry.Register(&redstoneTorchBlock{})
	Registry.Register(&unlitRedstoneTorchBlock{})
	Registry.Register(&inactiveRedstoneLampBlock{})
	Registry.Register(&activeRedstoneLampBlock{})
	Registry.Register(&leverBlock{})
	Registry.Register(&stoneButtonBlock{})
	Registry.Register(&woodenButtonBlock{})
	Registry.Register(&redstoneWireBlock{})
	Registry.Register(&stonePressurePlateBlock{})
	Registry.Register(&woodenPressurePlateBlock{})
	Registry.Register(&lightWeightedPressurePlateBlock{})
	Registry.Register(&heavyWeightedPressurePlateBlock{})
}
