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

func init() {
	Registry.Register(&redstoneTorchBlock{})
	Registry.Register(&unlitRedstoneTorchBlock{})
	Registry.Register(&inactiveRedstoneLampBlock{})
	Registry.Register(&activeRedstoneLampBlock{})
	Registry.Register(&leverBlock{})
	Registry.Register(&stoneButtonBlock{})
	Registry.Register(&woodenButtonBlock{})
	Registry.Register(&redstoneWireBlock{})
}
