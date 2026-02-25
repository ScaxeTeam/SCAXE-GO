package block

type carpetBlock struct{ DefaultBlockInteraction }

func (b *carpetBlock) GetID() uint8                { return CARPET }
func (b *carpetBlock) GetName() string             { return "Carpet" }
func (b *carpetBlock) GetHardness() float64        { return 0.1 }
func (b *carpetBlock) GetBlastResistance() float64 { return 0.5 }
func (b *carpetBlock) GetLightLevel() uint8        { return 0 }
func (b *carpetBlock) GetLightFilter() uint8       { return 0 }
func (b *carpetBlock) IsSolid() bool               { return false }
func (b *carpetBlock) IsTransparent() bool         { return true }
func (b *carpetBlock) CanBePlaced() bool           { return true }
func (b *carpetBlock) CanBeReplaced() bool         { return false }
func (b *carpetBlock) GetToolType() int            { return ToolTypeNone }
func (b *carpetBlock) GetToolTier() int            { return 0 }
func (b *carpetBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(CARPET), Meta: 0, Count: 1}}
}

type snowLayerBlock struct{ DefaultBlockInteraction }

func (b *snowLayerBlock) GetID() uint8                { return SNOW_LAYER }
func (b *snowLayerBlock) GetName() string             { return "Snow Layer" }
func (b *snowLayerBlock) GetHardness() float64        { return 0.1 }
func (b *snowLayerBlock) GetBlastResistance() float64 { return 0.5 }
func (b *snowLayerBlock) GetLightLevel() uint8        { return 0 }
func (b *snowLayerBlock) GetLightFilter() uint8       { return 0 }
func (b *snowLayerBlock) IsSolid() bool               { return false }
func (b *snowLayerBlock) IsTransparent() bool         { return true }
func (b *snowLayerBlock) CanBePlaced() bool           { return true }
func (b *snowLayerBlock) CanBeReplaced() bool         { return true }
func (b *snowLayerBlock) GetToolType() int            { return ToolTypeShovel }
func (b *snowLayerBlock) GetToolTier() int            { return 0 }
func (b *snowLayerBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypeShovel {
		return nil
	}
	return []Drop{{ID: 332, Meta: 0, Count: 1}}
}

type snowBlockBlock struct{ DefaultBlockInteraction }

func (b *snowBlockBlock) GetID() uint8                { return SNOW_BLOCK }
func (b *snowBlockBlock) GetName() string             { return "Snow Block" }
func (b *snowBlockBlock) GetHardness() float64        { return 0.2 }
func (b *snowBlockBlock) GetBlastResistance() float64 { return 1.0 }
func (b *snowBlockBlock) GetLightLevel() uint8        { return 0 }
func (b *snowBlockBlock) GetLightFilter() uint8       { return 15 }
func (b *snowBlockBlock) IsSolid() bool               { return true }
func (b *snowBlockBlock) IsTransparent() bool         { return false }
func (b *snowBlockBlock) CanBePlaced() bool           { return true }
func (b *snowBlockBlock) CanBeReplaced() bool         { return false }
func (b *snowBlockBlock) GetToolType() int            { return ToolTypeShovel }
func (b *snowBlockBlock) GetToolTier() int            { return 0 }
func (b *snowBlockBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypeShovel {
		return nil
	}
	return []Drop{{ID: 332, Meta: 0, Count: 4}}
}

type vineBlock struct{ DefaultBlockInteraction }

func (b *vineBlock) GetID() uint8                { return VINE }
func (b *vineBlock) GetName() string             { return "Vine" }
func (b *vineBlock) GetHardness() float64        { return 0.2 }
func (b *vineBlock) GetBlastResistance() float64 { return 1.0 }
func (b *vineBlock) GetLightLevel() uint8        { return 0 }
func (b *vineBlock) GetLightFilter() uint8       { return 0 }
func (b *vineBlock) IsSolid() bool               { return false }
func (b *vineBlock) IsTransparent() bool         { return true }
func (b *vineBlock) CanBePlaced() bool           { return true }
func (b *vineBlock) CanBeReplaced() bool         { return true }
func (b *vineBlock) GetToolType() int            { return ToolTypeShears }
func (b *vineBlock) GetToolTier() int            { return 0 }
func (b *vineBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType == ToolTypeShears {
		return []Drop{{ID: int(VINE), Meta: 0, Count: 1}}
	}
	return nil
}

type ladderBlock struct{ DefaultBlockInteraction }

func (b *ladderBlock) GetID() uint8                { return LADDER }
func (b *ladderBlock) GetName() string             { return "Ladder" }
func (b *ladderBlock) GetHardness() float64        { return 0.4 }
func (b *ladderBlock) GetBlastResistance() float64 { return 2.0 }
func (b *ladderBlock) GetLightLevel() uint8        { return 0 }
func (b *ladderBlock) GetLightFilter() uint8       { return 0 }
func (b *ladderBlock) IsSolid() bool               { return false }
func (b *ladderBlock) IsTransparent() bool         { return true }
func (b *ladderBlock) CanBePlaced() bool           { return true }
func (b *ladderBlock) CanBeReplaced() bool         { return false }
func (b *ladderBlock) GetToolType() int            { return ToolTypeAxe }
func (b *ladderBlock) GetToolTier() int            { return 0 }
func (b *ladderBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(LADDER), Meta: 0, Count: 1}}
}

type cobwebBlock struct{ DefaultBlockInteraction }

func (b *cobwebBlock) GetID() uint8                { return COBWEB }
func (b *cobwebBlock) GetName() string             { return "Cobweb" }
func (b *cobwebBlock) GetHardness() float64        { return 4.0 }
func (b *cobwebBlock) GetBlastResistance() float64 { return 20.0 }
func (b *cobwebBlock) GetLightLevel() uint8        { return 0 }
func (b *cobwebBlock) GetLightFilter() uint8       { return 0 }
func (b *cobwebBlock) IsSolid() bool               { return false }
func (b *cobwebBlock) IsTransparent() bool         { return true }
func (b *cobwebBlock) CanBePlaced() bool           { return true }
func (b *cobwebBlock) CanBeReplaced() bool         { return false }
func (b *cobwebBlock) GetToolType() int            { return ToolTypeSword }
func (b *cobwebBlock) GetToolTier() int            { return 0 }
func (b *cobwebBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType == ToolTypeShears {
		return []Drop{{ID: int(COBWEB), Meta: 0, Count: 1}}
	}
	return []Drop{{ID: 287, Meta: 0, Count: 1}}
}

type dandelionBlock struct{ DefaultBlockInteraction }

func (b *dandelionBlock) GetID() uint8                { return DANDELION }
func (b *dandelionBlock) GetName() string             { return "Dandelion" }
func (b *dandelionBlock) GetHardness() float64        { return 0 }
func (b *dandelionBlock) GetBlastResistance() float64 { return 0 }
func (b *dandelionBlock) GetLightLevel() uint8        { return 0 }
func (b *dandelionBlock) GetLightFilter() uint8       { return 0 }
func (b *dandelionBlock) IsSolid() bool               { return false }
func (b *dandelionBlock) IsTransparent() bool         { return true }
func (b *dandelionBlock) CanBePlaced() bool           { return true }
func (b *dandelionBlock) CanBeReplaced() bool         { return false }
func (b *dandelionBlock) GetToolType() int            { return ToolTypeNone }
func (b *dandelionBlock) GetToolTier() int            { return 0 }
func (b *dandelionBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(DANDELION), Meta: 0, Count: 1}}
}

type redFlowerBlock struct{ DefaultBlockInteraction }

func (b *redFlowerBlock) GetID() uint8                { return RED_FLOWER }
func (b *redFlowerBlock) GetName() string             { return "Red Flower" }
func (b *redFlowerBlock) GetHardness() float64        { return 0 }
func (b *redFlowerBlock) GetBlastResistance() float64 { return 0 }
func (b *redFlowerBlock) GetLightLevel() uint8        { return 0 }
func (b *redFlowerBlock) GetLightFilter() uint8       { return 0 }
func (b *redFlowerBlock) IsSolid() bool               { return false }
func (b *redFlowerBlock) IsTransparent() bool         { return true }
func (b *redFlowerBlock) CanBePlaced() bool           { return true }
func (b *redFlowerBlock) CanBeReplaced() bool         { return false }
func (b *redFlowerBlock) GetToolType() int            { return ToolTypeNone }
func (b *redFlowerBlock) GetToolTier() int            { return 0 }
func (b *redFlowerBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(RED_FLOWER), Meta: 0, Count: 1}}
}

type brownMushroomBlock struct{ DefaultBlockInteraction }

func (b *brownMushroomBlock) GetID() uint8                { return BROWN_MUSHROOM }
func (b *brownMushroomBlock) GetName() string             { return "Brown Mushroom" }
func (b *brownMushroomBlock) GetHardness() float64        { return 0 }
func (b *brownMushroomBlock) GetBlastResistance() float64 { return 0 }
func (b *brownMushroomBlock) GetLightLevel() uint8        { return 1 }
func (b *brownMushroomBlock) GetLightFilter() uint8       { return 0 }
func (b *brownMushroomBlock) IsSolid() bool               { return false }
func (b *brownMushroomBlock) IsTransparent() bool         { return true }
func (b *brownMushroomBlock) CanBePlaced() bool           { return true }
func (b *brownMushroomBlock) CanBeReplaced() bool         { return false }
func (b *brownMushroomBlock) GetToolType() int            { return ToolTypeNone }
func (b *brownMushroomBlock) GetToolTier() int            { return 0 }
func (b *brownMushroomBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(BROWN_MUSHROOM), Meta: 0, Count: 1}}
}

type redMushroomBlock struct{ DefaultBlockInteraction }

func (b *redMushroomBlock) GetID() uint8                { return RED_MUSHROOM }
func (b *redMushroomBlock) GetName() string             { return "Red Mushroom" }
func (b *redMushroomBlock) GetHardness() float64        { return 0 }
func (b *redMushroomBlock) GetBlastResistance() float64 { return 0 }
func (b *redMushroomBlock) GetLightLevel() uint8        { return 0 }
func (b *redMushroomBlock) GetLightFilter() uint8       { return 0 }
func (b *redMushroomBlock) IsSolid() bool               { return false }
func (b *redMushroomBlock) IsTransparent() bool         { return true }
func (b *redMushroomBlock) CanBePlaced() bool           { return true }
func (b *redMushroomBlock) CanBeReplaced() bool         { return false }
func (b *redMushroomBlock) GetToolType() int            { return ToolTypeNone }
func (b *redMushroomBlock) GetToolTier() int            { return 0 }
func (b *redMushroomBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(RED_MUSHROOM), Meta: 0, Count: 1}}
}

type deadBushBlock struct{ DefaultBlockInteraction }

func (b *deadBushBlock) GetID() uint8                { return DEAD_BUSH }
func (b *deadBushBlock) GetName() string             { return "Dead Bush" }
func (b *deadBushBlock) GetHardness() float64        { return 0 }
func (b *deadBushBlock) GetBlastResistance() float64 { return 0 }
func (b *deadBushBlock) GetLightLevel() uint8        { return 0 }
func (b *deadBushBlock) GetLightFilter() uint8       { return 0 }
func (b *deadBushBlock) IsSolid() bool               { return false }
func (b *deadBushBlock) IsTransparent() bool         { return true }
func (b *deadBushBlock) CanBePlaced() bool           { return true }
func (b *deadBushBlock) CanBeReplaced() bool         { return true }
func (b *deadBushBlock) GetToolType() int            { return ToolTypeNone }
func (b *deadBushBlock) GetToolTier() int            { return 0 }
func (b *deadBushBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType == ToolTypeShears {
		return []Drop{{ID: int(DEAD_BUSH), Meta: 0, Count: 1}}
	}
	return []Drop{{ID: 280, Meta: 0, Count: 1}}
}

type pumpkinBlock struct{ DefaultBlockInteraction }

func (b *pumpkinBlock) GetID() uint8                { return PUMPKIN }
func (b *pumpkinBlock) GetName() string             { return "Pumpkin" }
func (b *pumpkinBlock) GetHardness() float64        { return 1.0 }
func (b *pumpkinBlock) GetBlastResistance() float64 { return 5.0 }
func (b *pumpkinBlock) GetLightLevel() uint8        { return 0 }
func (b *pumpkinBlock) GetLightFilter() uint8       { return 15 }
func (b *pumpkinBlock) IsSolid() bool               { return true }
func (b *pumpkinBlock) IsTransparent() bool         { return false }
func (b *pumpkinBlock) CanBePlaced() bool           { return true }
func (b *pumpkinBlock) CanBeReplaced() bool         { return false }
func (b *pumpkinBlock) GetToolType() int            { return ToolTypeAxe }
func (b *pumpkinBlock) GetToolTier() int            { return 0 }
func (b *pumpkinBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(PUMPKIN), Meta: 0, Count: 1}}
}

type litPumpkinBlock struct{ DefaultBlockInteraction }

func (b *litPumpkinBlock) GetID() uint8                { return LIT_PUMPKIN }
func (b *litPumpkinBlock) GetName() string             { return "Jack o'Lantern" }
func (b *litPumpkinBlock) GetHardness() float64        { return 1.0 }
func (b *litPumpkinBlock) GetBlastResistance() float64 { return 5.0 }
func (b *litPumpkinBlock) GetLightLevel() uint8        { return 15 }
func (b *litPumpkinBlock) GetLightFilter() uint8       { return 15 }
func (b *litPumpkinBlock) IsSolid() bool               { return true }
func (b *litPumpkinBlock) IsTransparent() bool         { return false }
func (b *litPumpkinBlock) CanBePlaced() bool           { return true }
func (b *litPumpkinBlock) CanBeReplaced() bool         { return false }
func (b *litPumpkinBlock) GetToolType() int            { return ToolTypeAxe }
func (b *litPumpkinBlock) GetToolTier() int            { return 0 }
func (b *litPumpkinBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(LIT_PUMPKIN), Meta: 0, Count: 1}}
}

type melonBlockBlock struct{ DefaultBlockInteraction }

func (b *melonBlockBlock) GetID() uint8                { return MELON_BLOCK }
func (b *melonBlockBlock) GetName() string             { return "Melon Block" }
func (b *melonBlockBlock) GetHardness() float64        { return 1.0 }
func (b *melonBlockBlock) GetBlastResistance() float64 { return 5.0 }
func (b *melonBlockBlock) GetLightLevel() uint8        { return 0 }
func (b *melonBlockBlock) GetLightFilter() uint8       { return 15 }
func (b *melonBlockBlock) IsSolid() bool               { return true }
func (b *melonBlockBlock) IsTransparent() bool         { return false }
func (b *melonBlockBlock) CanBePlaced() bool           { return true }
func (b *melonBlockBlock) CanBeReplaced() bool         { return false }
func (b *melonBlockBlock) GetToolType() int            { return ToolTypeAxe }
func (b *melonBlockBlock) GetToolTier() int            { return 0 }
func (b *melonBlockBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 360, Meta: 0, Count: 4}}
}

type hayBaleBlock struct{ DefaultBlockInteraction }

func (b *hayBaleBlock) GetID() uint8                { return HAY_BALE }
func (b *hayBaleBlock) GetName() string             { return "Hay Bale" }
func (b *hayBaleBlock) GetHardness() float64        { return 0.5 }
func (b *hayBaleBlock) GetBlastResistance() float64 { return 2.5 }
func (b *hayBaleBlock) GetLightLevel() uint8        { return 0 }
func (b *hayBaleBlock) GetLightFilter() uint8       { return 15 }
func (b *hayBaleBlock) IsSolid() bool               { return true }
func (b *hayBaleBlock) IsTransparent() bool         { return false }
func (b *hayBaleBlock) CanBePlaced() bool           { return true }
func (b *hayBaleBlock) CanBeReplaced() bool         { return false }
func (b *hayBaleBlock) GetToolType() int            { return ToolTypeNone }
func (b *hayBaleBlock) GetToolTier() int            { return 0 }
func (b *hayBaleBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(HAY_BALE), Meta: 0, Count: 1}}
}

type flowerPotBlock struct{ DefaultBlockInteraction }

func (b *flowerPotBlock) GetID() uint8                { return FLOWER_POT_BLOCK }
func (b *flowerPotBlock) GetName() string             { return "Flower Pot" }
func (b *flowerPotBlock) GetHardness() float64        { return 0 }
func (b *flowerPotBlock) GetBlastResistance() float64 { return 0 }
func (b *flowerPotBlock) GetLightLevel() uint8        { return 0 }
func (b *flowerPotBlock) GetLightFilter() uint8       { return 0 }
func (b *flowerPotBlock) IsSolid() bool               { return false }
func (b *flowerPotBlock) IsTransparent() bool         { return true }
func (b *flowerPotBlock) CanBePlaced() bool           { return true }
func (b *flowerPotBlock) CanBeReplaced() bool         { return false }
func (b *flowerPotBlock) GetToolType() int            { return ToolTypeNone }
func (b *flowerPotBlock) GetToolTier() int            { return 0 }
func (b *flowerPotBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 390, Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(&carpetBlock{})
	Registry.Register(&snowLayerBlock{})
	Registry.Register(&snowBlockBlock{})
	Registry.Register(&vineBlock{})
	Registry.Register(&ladderBlock{})
	Registry.Register(&cobwebBlock{})
	Registry.Register(&dandelionBlock{})
	Registry.Register(&redFlowerBlock{})
	Registry.Register(&brownMushroomBlock{})
	Registry.Register(&redMushroomBlock{})
	Registry.Register(&deadBushBlock{})
	Registry.Register(&pumpkinBlock{})
	Registry.Register(&litPumpkinBlock{})
	Registry.Register(&melonBlockBlock{})
	Registry.Register(&hayBaleBlock{})
	Registry.Register(&flowerPotBlock{})
}
