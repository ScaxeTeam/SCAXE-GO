package block

type BlockProperty struct {
	ID              uint8
	Name            string
	Hardness        float64
	BlastResistance float64
	LightLevel      uint8
	LightFilter     uint8
	Solid           bool
	Transparent     bool
	Replaceable     bool
	ToolType        int
	ToolTier        int
	FuelTime        int
	FlammableChance int
	BurnChance      int
}

var defaultProperty = BlockProperty{
	Name:            "Unknown",
	Hardness:        10,
	BlastResistance: 50,
	LightLevel:      0,
	LightFilter:     15,
	Solid:           true,
	Transparent:     false,
	Replaceable:     false,
	ToolType:        ToolTypeNone,
	ToolTier:        0,
}

var vanillaBlocks = [256]BlockProperty{

	AIR: {ID: AIR, Name: "Air", Hardness: 0, BlastResistance: 0, LightFilter: 0, Solid: false, Transparent: true, Replaceable: true},

	STONE: {ID: STONE, Name: "Stone", Hardness: 1.5, BlastResistance: 30, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	GRASS: {ID: GRASS, Name: "Grass", Hardness: 0.6, BlastResistance: 3, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},

	DIRT: {ID: DIRT, Name: "Dirt", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},

	COBBLESTONE: {ID: COBBLESTONE, Name: "Cobblestone", Hardness: 2.0, BlastResistance: 30, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	PLANKS: {ID: PLANKS, Name: "Planks", Hardness: 2.0, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe, FuelTime: 300, FlammableChance: 5, BurnChance: 20},

	SAPLING: {ID: SAPLING, Name: "Sapling", Hardness: 0, BlastResistance: 0, LightFilter: 0, Solid: false, Transparent: true, Replaceable: true},

	BEDROCK: {ID: BEDROCK, Name: "Bedrock", Hardness: -1, BlastResistance: 18000000, LightFilter: 15, Solid: true},

	WATER: {ID: WATER, Name: "Water", Hardness: 100, BlastResistance: 500, LightFilter: 2, Solid: false, Transparent: true, Replaceable: true},

	STILL_WATER: {ID: STILL_WATER, Name: "Still Water", Hardness: 100, BlastResistance: 500, LightFilter: 2, Solid: false, Transparent: true, Replaceable: true},

	LAVA: {ID: LAVA, Name: "Lava", Hardness: 100, BlastResistance: 500, LightLevel: 15, LightFilter: 0, Solid: false, Transparent: true, Replaceable: true},

	STILL_LAVA: {ID: STILL_LAVA, Name: "Still Lava", Hardness: 100, BlastResistance: 500, LightLevel: 15, LightFilter: 0, Solid: false, Transparent: true, Replaceable: true},

	SAND: {ID: SAND, Name: "Sand", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},

	GRAVEL: {ID: GRAVEL, Name: "Gravel", Hardness: 0.6, BlastResistance: 3, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},

	GOLD_ORE: {ID: GOLD_ORE, Name: "Gold Ore", Hardness: 3.0, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierIron},

	IRON_ORE: {ID: IRON_ORE, Name: "Iron Ore", Hardness: 3.0, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierStone},

	COAL_ORE: {ID: COAL_ORE, Name: "Coal Ore", Hardness: 3.0, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	WOOD: {ID: WOOD, Name: "Wood", Hardness: 2.0, BlastResistance: 10, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe, FuelTime: 300, FlammableChance: 5, BurnChance: 5},

	LEAVES: {ID: LEAVES, Name: "Leaves", Hardness: 0.2, BlastResistance: 1, LightFilter: 1, Solid: false, Transparent: true, ToolType: ToolTypeShears, FlammableChance: 30, BurnChance: 60},

	SPONGE: {ID: SPONGE, Name: "Sponge", Hardness: 0.6, BlastResistance: 3, LightFilter: 15, Solid: true},

	GLASS: {ID: GLASS, Name: "Glass", Hardness: 0.3, BlastResistance: 1.5, LightFilter: 0, Solid: true, Transparent: true},

	LAPIS_ORE: {ID: LAPIS_ORE, Name: "Lapis Lazuli Ore", Hardness: 3.0, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierStone},

	LAPIS_BLOCK: {ID: LAPIS_BLOCK, Name: "Lapis Lazuli Block", Hardness: 3.0, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierStone},

	SANDSTONE: {ID: SANDSTONE, Name: "Sandstone", Hardness: 0.8, BlastResistance: 4, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	COBWEB: {ID: COBWEB, Name: "Cobweb", Hardness: 4.0, BlastResistance: 20, LightFilter: 0, Solid: false, Transparent: true, ToolType: ToolTypeSword},

	TALL_GRASS: {ID: TALL_GRASS, Name: "Tall Grass", Hardness: 0, BlastResistance: 0, LightFilter: 0, Solid: false, Transparent: true, Replaceable: true, ToolType: ToolTypeShears, FlammableChance: 60, BurnChance: 100},

	DEAD_BUSH: {ID: DEAD_BUSH, Name: "Dead Bush", Hardness: 0, BlastResistance: 0, LightFilter: 0, Solid: false, Transparent: true, Replaceable: true, FlammableChance: 60, BurnChance: 100},

	WOOL: {ID: WOOL, Name: "Wool", Hardness: 0.8, BlastResistance: 4, LightFilter: 15, Solid: true, ToolType: ToolTypeShears, FlammableChance: 30, BurnChance: 60},

	DANDELION: {ID: DANDELION, Name: "Dandelion", Hardness: 0, BlastResistance: 0, LightFilter: 0, Solid: false, Transparent: true, Replaceable: true},

	RED_FLOWER: {ID: RED_FLOWER, Name: "Flower", Hardness: 0, BlastResistance: 0, LightFilter: 0, Solid: false, Transparent: true, Replaceable: true},

	GOLD_BLOCK: {ID: GOLD_BLOCK, Name: "Gold Block", Hardness: 3.0, BlastResistance: 30, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierIron},

	IRON_BLOCK: {ID: IRON_BLOCK, Name: "Iron Block", Hardness: 5.0, BlastResistance: 30, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierStone},

	BRICKS_BLOCK: {ID: BRICKS_BLOCK, Name: "Bricks", Hardness: 2.0, BlastResistance: 30, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	TNT: {ID: TNT, Name: "TNT", Hardness: 0, BlastResistance: 0, LightFilter: 15, Solid: true, FlammableChance: 15, BurnChance: 100},

	BOOKSHELF: {ID: BOOKSHELF, Name: "Bookshelf", Hardness: 1.5, BlastResistance: 7.5, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe, FlammableChance: 30, BurnChance: 20},

	MOSS_STONE: {ID: MOSS_STONE, Name: "Moss Stone", Hardness: 2.0, BlastResistance: 30, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	OBSIDIAN: {ID: OBSIDIAN, Name: "Obsidian", Hardness: 50.0, BlastResistance: 6000, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierDiamond},

	TORCH: {ID: TORCH, Name: "Torch", Hardness: 0, BlastResistance: 0, LightLevel: 14, LightFilter: 0, Solid: false, Transparent: true},

	MONSTER_SPAWNER: {ID: MONSTER_SPAWNER, Name: "Monster Spawner", Hardness: 5.0, BlastResistance: 25, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	CHEST: {ID: CHEST, Name: "Chest", Hardness: 2.5, BlastResistance: 12.5, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe},

	DIAMOND_ORE: {ID: DIAMOND_ORE, Name: "Diamond Ore", Hardness: 3.0, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierIron},

	DIAMOND_BLOCK: {ID: DIAMOND_BLOCK, Name: "Diamond Block", Hardness: 5.0, BlastResistance: 30, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierIron},

	WORKBENCH: {ID: WORKBENCH, Name: "Crafting Table", Hardness: 2.5, BlastResistance: 12.5, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe, FlammableChance: 5, BurnChance: 20},

	FARMLAND: {ID: FARMLAND, Name: "Farmland", Hardness: 0.6, BlastResistance: 3, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},

	FURNACE: {ID: FURNACE, Name: "Furnace", Hardness: 3.5, BlastResistance: 17.5, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	BURNING_FURNACE: {ID: BURNING_FURNACE, Name: "Burning Furnace", Hardness: 3.5, BlastResistance: 17.5, LightLevel: 13, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	REDSTONE_ORE: {ID: REDSTONE_ORE, Name: "Redstone Ore", Hardness: 3.0, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierIron},

	GLOWING_REDSTONE_ORE: {ID: GLOWING_REDSTONE_ORE, Name: "Glowing Redstone Ore", Hardness: 3.0, BlastResistance: 15, LightLevel: 9, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierIron},

	ICE: {ID: ICE, Name: "Ice", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 2, Solid: true, Transparent: true, ToolType: ToolTypePickaxe},

	SNOW_BLOCK: {ID: SNOW_BLOCK, Name: "Snow Block", Hardness: 0.2, BlastResistance: 1, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel, ToolTier: TierWooden},

	CACTUS: {ID: CACTUS, Name: "Cactus", Hardness: 0.4, BlastResistance: 2, LightFilter: 0, Solid: true, Transparent: true},

	CLAY_BLOCK: {ID: CLAY_BLOCK, Name: "Clay", Hardness: 0.6, BlastResistance: 3, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},

	NETHERRACK: {ID: NETHERRACK, Name: "Netherrack", Hardness: 0.4, BlastResistance: 2, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	SOUL_SAND: {ID: SOUL_SAND, Name: "Soul Sand", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},

	GLOWSTONE_BLOCK: {ID: GLOWSTONE_BLOCK, Name: "Glowstone", Hardness: 0.3, BlastResistance: 1.5, LightLevel: 15, LightFilter: 15, Solid: true},

	STONE_BRICKS: {ID: STONE_BRICKS, Name: "Stone Bricks", Hardness: 1.5, BlastResistance: 30, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	BROWN_MUSHROOM_BLOCK: {ID: BROWN_MUSHROOM_BLOCK, Name: "Brown Mushroom Block", Hardness: 0.2, BlastResistance: 1, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe},

	RED_MUSHROOM_BLOCK: {ID: RED_MUSHROOM_BLOCK, Name: "Red Mushroom Block", Hardness: 0.2, BlastResistance: 1, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe},

	IRON_BARS: {ID: IRON_BARS, Name: "Iron Bars", Hardness: 5.0, BlastResistance: 30, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	GLASS_PANE: {ID: GLASS_PANE, Name: "Glass Pane", Hardness: 0.3, BlastResistance: 1.5, LightFilter: 0, Solid: true, Transparent: true},

	NETHER_BRICKS: {ID: NETHER_BRICKS, Name: "Nether Bricks", Hardness: 2.0, BlastResistance: 30, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	END_STONE: {ID: END_STONE, Name: "End Stone", Hardness: 3.0, BlastResistance: 45, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	EMERALD_ORE: {ID: EMERALD_ORE, Name: "Emerald Ore", Hardness: 3.0, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierIron},

	EMERALD_BLOCK: {ID: EMERALD_BLOCK, Name: "Emerald Block", Hardness: 5.0, BlastResistance: 30, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierIron},

	REDSTONE_BLOCK: {ID: REDSTONE_BLOCK, Name: "Block of Redstone", Hardness: 5.0, BlastResistance: 30, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	NETHER_QUARTZ_ORE: {ID: NETHER_QUARTZ_ORE, Name: "Nether Quartz Ore", Hardness: 3.0, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	QUARTZ_BLOCK: {ID: QUARTZ_BLOCK, Name: "Block of Quartz", Hardness: 0.8, BlastResistance: 4, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden},

	SLIME_BLOCK: {ID: SLIME_BLOCK, Name: "Slime Block", Hardness: 0, BlastResistance: 0, LightFilter: 0, Solid: true, Transparent: true},

	COAL_BLOCK: {ID: COAL_BLOCK, Name: "Block of Coal", Hardness: 5.0, BlastResistance: 30, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, ToolTier: TierWooden, FuelTime: 16000, FlammableChance: 5, BurnChance: 5},

	PACKED_ICE: {ID: PACKED_ICE, Name: "Packed Ice", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
}

func GetProperty(id uint8) BlockProperty {
	prop := vanillaBlocks[id]
	if prop.Name == "" {

		result := defaultProperty
		result.ID = id
		return result
	}
	return prop
}

func GetName(id uint8) string {
	return GetProperty(id).Name
}

func GetToolType(id uint8) int {
	return GetProperty(id).ToolType
}

func GetToolTier(id uint8) int {
	return GetProperty(id).ToolTier
}

func IsFlammable(id uint8) bool {
	return GetProperty(id).FlammableChance > 0
}

func CanBeFuel(id uint8) bool {
	return GetProperty(id).FuelTime > 0
}

func GetFuelTime(id uint8) int {
	return GetProperty(id).FuelTime
}
