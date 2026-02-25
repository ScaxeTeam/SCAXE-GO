package block

type saplingBlock struct{ DefaultBlockInteraction }

func (b *saplingBlock) GetID() uint8                { return SAPLING }
func (b *saplingBlock) GetName() string             { return "Sapling" }
func (b *saplingBlock) GetHardness() float64        { return 0 }
func (b *saplingBlock) GetBlastResistance() float64 { return 0 }
func (b *saplingBlock) GetLightLevel() uint8        { return 0 }
func (b *saplingBlock) GetLightFilter() uint8       { return 0 }
func (b *saplingBlock) IsSolid() bool               { return false }
func (b *saplingBlock) IsTransparent() bool         { return true }
func (b *saplingBlock) CanBePlaced() bool           { return true }
func (b *saplingBlock) CanBeReplaced() bool         { return false }
func (b *saplingBlock) GetToolType() int            { return ToolTypeNone }
func (b *saplingBlock) GetToolTier() int            { return 0 }
func (b *saplingBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(SAPLING), Meta: 0, Count: 1}}
}

type tallGrassBlock struct{ DefaultBlockInteraction }

func (b *tallGrassBlock) GetID() uint8                { return TALL_GRASS }
func (b *tallGrassBlock) GetName() string             { return "Tall Grass" }
func (b *tallGrassBlock) GetHardness() float64        { return 0 }
func (b *tallGrassBlock) GetBlastResistance() float64 { return 0 }
func (b *tallGrassBlock) GetLightLevel() uint8        { return 0 }
func (b *tallGrassBlock) GetLightFilter() uint8       { return 0 }
func (b *tallGrassBlock) IsSolid() bool               { return false }
func (b *tallGrassBlock) IsTransparent() bool         { return true }
func (b *tallGrassBlock) CanBePlaced() bool           { return true }
func (b *tallGrassBlock) CanBeReplaced() bool         { return true }
func (b *tallGrassBlock) GetToolType() int            { return ToolTypeNone }
func (b *tallGrassBlock) GetToolTier() int            { return 0 }
func (b *tallGrassBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType == ToolTypeShears {
		return []Drop{{ID: int(TALL_GRASS), Meta: 0, Count: 1}}
	}
	return nil
}

type bedBlock struct{ DefaultBlockInteraction }

func (b *bedBlock) GetID() uint8                { return BED_BLOCK }
func (b *bedBlock) GetName() string             { return "Bed Block" }
func (b *bedBlock) GetHardness() float64        { return 0.2 }
func (b *bedBlock) GetBlastResistance() float64 { return 1.0 }
func (b *bedBlock) GetLightLevel() uint8        { return 0 }
func (b *bedBlock) GetLightFilter() uint8       { return 0 }
func (b *bedBlock) IsSolid() bool               { return true }
func (b *bedBlock) IsTransparent() bool         { return true }
func (b *bedBlock) CanBePlaced() bool           { return true }
func (b *bedBlock) CanBeReplaced() bool         { return false }
func (b *bedBlock) GetToolType() int            { return ToolTypeNone }
func (b *bedBlock) GetToolTier() int            { return 0 }
func (b *bedBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 355, Meta: 0, Count: 1}}
}

type doublePlantBlock struct{ DefaultBlockInteraction }

func (b *doublePlantBlock) GetID() uint8                { return DOUBLE_PLANT }
func (b *doublePlantBlock) GetName() string             { return "Double Plant" }
func (b *doublePlantBlock) GetHardness() float64        { return 0 }
func (b *doublePlantBlock) GetBlastResistance() float64 { return 0 }
func (b *doublePlantBlock) GetLightLevel() uint8        { return 0 }
func (b *doublePlantBlock) GetLightFilter() uint8       { return 0 }
func (b *doublePlantBlock) IsSolid() bool               { return false }
func (b *doublePlantBlock) IsTransparent() bool         { return true }
func (b *doublePlantBlock) CanBePlaced() bool           { return true }
func (b *doublePlantBlock) CanBeReplaced() bool         { return true }
func (b *doublePlantBlock) GetToolType() int            { return ToolTypeNone }
func (b *doublePlantBlock) GetToolTier() int            { return 0 }
func (b *doublePlantBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType == ToolTypeShears {
		return []Drop{{ID: int(DOUBLE_PLANT), Meta: 0, Count: 1}}
	}
	return nil
}

type waterLilyBlock struct{ DefaultBlockInteraction }

func (b *waterLilyBlock) GetID() uint8                { return WATER_LILY }
func (b *waterLilyBlock) GetName() string             { return "Lily Pad" }
func (b *waterLilyBlock) GetHardness() float64        { return 0 }
func (b *waterLilyBlock) GetBlastResistance() float64 { return 0 }
func (b *waterLilyBlock) GetLightLevel() uint8        { return 0 }
func (b *waterLilyBlock) GetLightFilter() uint8       { return 0 }
func (b *waterLilyBlock) IsSolid() bool               { return false }
func (b *waterLilyBlock) IsTransparent() bool         { return true }
func (b *waterLilyBlock) CanBePlaced() bool           { return true }
func (b *waterLilyBlock) CanBeReplaced() bool         { return false }
func (b *waterLilyBlock) GetToolType() int            { return ToolTypeNone }
func (b *waterLilyBlock) GetToolTier() int            { return 0 }
func (b *waterLilyBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: int(WATER_LILY), Meta: 0, Count: 1}}
}

type signPostBlock struct{ DefaultBlockInteraction }

func (b *signPostBlock) GetID() uint8                { return SIGN_POST }
func (b *signPostBlock) GetName() string             { return "Sign Post" }
func (b *signPostBlock) GetHardness() float64        { return 1.0 }
func (b *signPostBlock) GetBlastResistance() float64 { return 5.0 }
func (b *signPostBlock) GetLightLevel() uint8        { return 0 }
func (b *signPostBlock) GetLightFilter() uint8       { return 0 }
func (b *signPostBlock) IsSolid() bool               { return false }
func (b *signPostBlock) IsTransparent() bool         { return true }
func (b *signPostBlock) CanBePlaced() bool           { return true }
func (b *signPostBlock) CanBeReplaced() bool         { return false }
func (b *signPostBlock) GetToolType() int            { return ToolTypeAxe }
func (b *signPostBlock) GetToolTier() int            { return 0 }
func (b *signPostBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 323, Meta: 0, Count: 1}}
}

type wallSignBlock struct{ DefaultBlockInteraction }

func (b *wallSignBlock) GetID() uint8                { return WALL_SIGN }
func (b *wallSignBlock) GetName() string             { return "Wall Sign" }
func (b *wallSignBlock) GetHardness() float64        { return 1.0 }
func (b *wallSignBlock) GetBlastResistance() float64 { return 5.0 }
func (b *wallSignBlock) GetLightLevel() uint8        { return 0 }
func (b *wallSignBlock) GetLightFilter() uint8       { return 0 }
func (b *wallSignBlock) IsSolid() bool               { return false }
func (b *wallSignBlock) IsTransparent() bool         { return true }
func (b *wallSignBlock) CanBePlaced() bool           { return true }
func (b *wallSignBlock) CanBeReplaced() bool         { return false }
func (b *wallSignBlock) GetToolType() int            { return ToolTypeAxe }
func (b *wallSignBlock) GetToolTier() int            { return 0 }
func (b *wallSignBlock) GetDrops(toolType, toolTier int) []Drop {
	return []Drop{{ID: 323, Meta: 0, Count: 1}}
}

func init() {
	Registry.Register(&saplingBlock{})
	Registry.Register(&tallGrassBlock{})
	Registry.Register(&bedBlock{})
	Registry.Register(&doublePlantBlock{})
	Registry.Register(&waterLilyBlock{})
	Registry.Register(&signPostBlock{})
	Registry.Register(&wallSignBlock{})
	Registry.Register(&simpleBlock{id: SPONGE, name: "Sponge", hardness: 0.6})
	Registry.Register(&simpleBlock{id: SANDSTONE, name: "Sandstone", hardness: 0.8, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: NOTEBLOCK, name: "Noteblock", hardness: 0.8, toolType: ToolTypeAxe})
	Registry.Register(&simpleBlock{id: WOOL, name: "Wool", hardness: 0.8})
	Registry.Register(&simpleBlock{id: BRICKS_BLOCK, name: "Bricks", hardness: 2.0, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: BOOKSHELF, name: "Bookshelf", hardness: 1.5, toolType: ToolTypeAxe})
	Registry.Register(&simpleBlock{id: MOSS_STONE, name: "Moss Stone", hardness: 2.0, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: MONSTER_SPAWNER, name: "Monster Spawner", hardness: 5.0, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: WORKBENCH, name: "Crafting Table", hardness: 2.5, toolType: ToolTypeAxe})
	Registry.Register(&simpleBlock{id: FARMLAND, name: "Farmland", hardness: 0.6})
	Registry.Register(&simpleBlock{id: STONE_PRESSURE_PLATE, name: "Stone Pressure Plate", hardness: 0.5, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: WOODEN_PRESSURE_PLATE, name: "Wooden Pressure Plate", hardness: 0.5, toolType: ToolTypeAxe})
	Registry.Register(&simpleBlock{id: ICE, name: "Ice", hardness: 0.5, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: CLAY_BLOCK, name: "Clay", hardness: 0.6})
	Registry.Register(&simpleBlock{id: NETHERRACK, name: "Netherrack", hardness: 0.4, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: SOUL_SAND, name: "Soul Sand", hardness: 0.5})
	Registry.Register(&simpleBlock{id: CAKE_BLOCK, name: "Cake Block", hardness: 0.5})
	Registry.Register(&simpleBlock{id: INVISIBLE_BEDROCK, name: "Invisible Bedrock", hardness: -1, blastResistance: 18000000})
	Registry.Register(&simpleBlock{id: MONSTER_EGG, name: "Monster Egg", hardness: 0.75, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: STONE_BRICKS, name: "Stone Bricks", hardness: 1.5, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: BROWN_MUSHROOM_BLOCK, name: "Brown Mushroom Block", hardness: 0.2, toolType: ToolTypeAxe})
	Registry.Register(&simpleBlock{id: RED_MUSHROOM_BLOCK, name: "Red Mushroom Block", hardness: 0.2, toolType: ToolTypeAxe})
	Registry.Register(&simpleBlock{id: MYCELIUM, name: "Mycelium", hardness: 0.6})
	Registry.Register(&simpleBlock{id: NETHER_BRICKS, name: "Nether Bricks", hardness: 2.0, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: NETHER_WART_BLOCK, name: "Nether Wart Block", hardness: 0})
	Registry.Register(&simpleBlock{id: ENCHANTING_TABLE, name: "Enchanting Table", hardness: 5.0, lightLevel: 12, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: BREWING_STAND_BLOCK, name: "Brewing Stand", hardness: 0.5, lightLevel: 1, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: CAULDRON_BLOCK, name: "Cauldron", hardness: 2.0, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: END_PORTAL_FRAME, name: "End Portal Frame", hardness: -1, lightLevel: 1, blastResistance: 18000000})
	Registry.Register(&simpleBlock{id: END_STONE, name: "End Stone", hardness: 3.0, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: DROPPER, name: "Dropper", hardness: 3.5, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: COCOA_BLOCK, name: "Cocoa Block", hardness: 0.2})
	Registry.Register(&simpleBlock{id: TRIPWIRE_HOOK, name: "Tripwire Hook", hardness: 0})
	Registry.Register(&simpleBlock{id: TRIPWIRE, name: "Tripwire", hardness: 0})
	Registry.Register(&simpleBlock{id: EMERALD_BLOCK, name: "Emerald Block", hardness: 5.0, toolType: ToolTypePickaxe, toolTier: TierIron})
	Registry.Register(&simpleBlock{id: SKULL_BLOCK, name: "Skull", hardness: 1.0})
	Registry.Register(&simpleBlock{id: ANVIL, name: "Anvil", hardness: 5.0, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: TRAPPED_CHEST, name: "Trapped Chest", hardness: 2.5, toolType: ToolTypeAxe})
	Registry.Register(&simpleBlock{id: LIGHT_WEIGHTED_PRESSURE_PLATE, name: "Light Weighted Pressure Plate", hardness: 0.5, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: HEAVY_WEIGHTED_PRESSURE_PLATE, name: "Heavy Weighted Pressure Plate", hardness: 0.5, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: UNPOWERED_COMPARATOR_BLOCK, name: "Comparator", hardness: 0})
	Registry.Register(&simpleBlock{id: POWERED_COMPARATOR_BLOCK, name: "Powered Comparator", hardness: 0})
	Registry.Register(&simpleBlock{id: DAYLIGHT_SENSOR, name: "Daylight Sensor", hardness: 0.2, toolType: ToolTypeAxe})
	Registry.Register(&simpleBlock{id: REDSTONE_BLOCK, name: "Redstone Block", hardness: 5.0, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: HOPPER_BLOCK, name: "Hopper", hardness: 3.0, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: QUARTZ_BLOCK, name: "Quartz Block", hardness: 0.8, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: STAINED_CLAY, name: "Stained Clay", hardness: 1.25, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: LEAVES2, name: "Leaves", hardness: 0.2})
	Registry.Register(&simpleBlock{id: WOOD2, name: "Wood", hardness: 2.0, toolType: ToolTypeAxe})
	Registry.Register(&simpleBlock{id: SLIME_BLOCK, name: "Slime Block", hardness: 0})
	Registry.Register(&simpleBlock{id: HARDENED_CLAY, name: "Hardened Clay", hardness: 1.25, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: COAL_BLOCK, name: "Coal Block", hardness: 5.0, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: PACKED_ICE, name: "Packed Ice", hardness: 0.5, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: DAYLIGHT_SENSOR_INVERTED, name: "Inverted Daylight Sensor", hardness: 0.2, toolType: ToolTypeAxe})
	Registry.Register(&simpleBlock{id: RED_SANDSTONE, name: "Red Sandstone", hardness: 0.8, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: UNPOWERED_REPEATER, name: "Repeater", hardness: 0})
	Registry.Register(&simpleBlock{id: POWERED_REPEATER, name: "Powered Repeater", hardness: 0})
	Registry.Register(&simpleBlock{id: GRASS_PATH, name: "Grass Path", hardness: 0.65})
	Registry.Register(&simpleBlock{id: ITEM_FRAME_BLOCK, name: "Item Frame", hardness: 0})
	Registry.Register(&simpleBlock{id: PODZOL, name: "Podzol", hardness: 0.5})
	Registry.Register(&simpleBlock{id: STONECUTTER, name: "Stonecutter", hardness: 3.5, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: GLOWING_OBSIDIAN, name: "Glowing Obsidian", hardness: 50.0, lightLevel: 12, blastResistance: 6000, toolType: ToolTypePickaxe, toolTier: TierDiamond})
	Registry.Register(&simpleBlock{id: NETHER_REACTOR, name: "Nether Reactor Core", hardness: 3.0, toolType: ToolTypePickaxe})
	Registry.Register(&simpleBlock{id: NETHER_BRICK_FENCE, name: "Nether Brick Fence", hardness: 2.0, toolType: ToolTypePickaxe})
}
