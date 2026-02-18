package block

type BlockProperty struct {
	ID               uint8
	Name             string
	Hardness         float64
	BlastResistance  float64
	LightLevel       uint8
	LightFilter      uint8
	Solid            bool
	Transparent      bool
	Replaceable      bool
	Flowable         bool
	DiffusesSkyLight bool
	ToolType         int
	ToolTier         int
	FuelTime         int
	FlammableChance  int
	BurnChance       int
}

var defaultProperty = BlockProperty{
	Name:        "Unknown",
	Hardness:    0,
	LightFilter: 0,
	Solid:       true,
	Transparent: true,
}

var vanillaBlocks = [256]BlockProperty{
	AIR:                           {ID: AIR, Name: "Air", LightFilter: 0, Transparent: true, Replaceable: true, Flowable: true},
	STONE:                         {ID: STONE, Name: "Stone", Hardness: 1.5, BlastResistance: 7.5, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	GRASS:                         {ID: GRASS, Name: "Grass", Hardness: 0.6, BlastResistance: 3, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},
	DIRT:                          {ID: DIRT, Name: "Dirt", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},
	COBBLESTONE:                   {ID: COBBLESTONE, Name: "Cobblestone", Hardness: 2, BlastResistance: 10, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	PLANKS:                        {ID: PLANKS, Name: "Oak Wood Planks", Hardness: 2, BlastResistance: 10, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe, FlammableChance: 5, BurnChance: 20, FuelTime: 300},
	SAPLING:                       {ID: SAPLING, Name: "Oak Sapling", LightFilter: 0, Transparent: true, Flowable: true, FuelTime: 100},
	BEDROCK:                       {ID: BEDROCK, Name: "Bedrock", Hardness: -1, BlastResistance: 18000000, LightFilter: 15, Solid: true},
	WATER:                         {ID: WATER, Name: "Water", Hardness: 100, BlastResistance: 500, LightFilter: 2, Transparent: true, Replaceable: true, Flowable: true},
	STILL_WATER:                   {ID: STILL_WATER, Name: "Still Water", Hardness: 100, BlastResistance: 500, LightFilter: 2, Transparent: true, Replaceable: true, Flowable: true},
	LAVA:                          {ID: LAVA, Name: "Lava", Hardness: 100, BlastResistance: 500, LightLevel: 15, LightFilter: 0, Transparent: true, Replaceable: true, Flowable: true},
	STILL_LAVA:                    {ID: STILL_LAVA, Name: "Still Lava", Hardness: 100, BlastResistance: 500, LightLevel: 15, LightFilter: 0, Transparent: true, Replaceable: true, Flowable: true},
	SAND:                          {ID: SAND, Name: "Sand", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},
	GRAVEL:                        {ID: GRAVEL, Name: "Gravel", Hardness: 0.6, BlastResistance: 3, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},
	GOLD_ORE:                      {ID: GOLD_ORE, Name: "Gold Ore", Hardness: 3, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	IRON_ORE:                      {ID: IRON_ORE, Name: "Iron Ore", Hardness: 3, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	COAL_ORE:                      {ID: COAL_ORE, Name: "Coal Ore", Hardness: 3, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	WOOD:                          {ID: WOOD, Name: "Oak Wood", Hardness: 2, BlastResistance: 10, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe, FlammableChance: 5, BurnChance: 10, FuelTime: 300},
	LEAVES:                        {ID: LEAVES, Name: "Oak Leaves", Hardness: 0.2, BlastResistance: 1, LightFilter: 0, Solid: true, Transparent: true, DiffusesSkyLight: true, ToolType: ToolTypeShears, FlammableChance: 30, BurnChance: 60},
	SPONGE:                        {ID: SPONGE, Name: "Sponge", Hardness: 0.6, BlastResistance: 3, LightFilter: 15, Solid: true},
	GLASS:                         {ID: GLASS, Name: "Glass", Hardness: 0.3, BlastResistance: 1.5, LightFilter: 0, Solid: true, Transparent: true},
	LAPIS_ORE:                     {ID: LAPIS_ORE, Name: "Lapis Lazuli Ore", Hardness: 3, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	LAPIS_BLOCK:                   {ID: LAPIS_BLOCK, Name: "Lapis Lazuli Block", Hardness: 3, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	DISPENSER:                     {ID: DISPENSER, Name: "Dispenser", Hardness: 3.5, BlastResistance: 17.5, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	SANDSTONE:                     {ID: SANDSTONE, Name: "Sandstone", Hardness: 0.8, BlastResistance: 4, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	NOTEBLOCK:                     {ID: NOTEBLOCK, Name: "Noteblock", Hardness: 0.8, BlastResistance: 4, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe, FuelTime: 300},
	BED_BLOCK:                     {ID: BED_BLOCK, Name: "Bed Block", Hardness: 0.2, BlastResistance: 1, LightFilter: 0, Solid: true, Transparent: true},
	POWERED_RAIL:                  {ID: POWERED_RAIL, Name: "Powered Rail", Hardness: 0.7, LightFilter: 0, Transparent: true, Flowable: true},
	DETECTOR_RAIL:                 {ID: DETECTOR_RAIL, Name: "Detector Rail", Hardness: 0.7, LightFilter: 0, Transparent: true, Flowable: true},
	COBWEB:                        {ID: COBWEB, Name: "Cobweb", Hardness: 4, LightFilter: 0, Transparent: true, Flowable: true, DiffusesSkyLight: true, ToolType: ToolTypeShears},
	TALL_GRASS:                    {ID: TALL_GRASS, Name: "Dead Shrub", LightFilter: 0, Transparent: true, Replaceable: true, Flowable: true, ToolType: ToolTypeShears, FlammableChance: 60, BurnChance: 100},
	DEAD_BUSH:                     {ID: DEAD_BUSH, Name: "Dead Bush", LightFilter: 0, Transparent: true, Flowable: true},
	WOOL:                          {ID: WOOL, Name: "White Wool", Hardness: 0.8, BlastResistance: 4, LightFilter: 15, Solid: true, ToolType: ToolTypeShears, FlammableChance: 30, BurnChance: 60},
	DANDELION:                     {ID: DANDELION, Name: "Dandelion", LightFilter: 0, Transparent: true, Flowable: true},
	RED_FLOWER:                    {ID: RED_FLOWER, Name: "Poppy", LightFilter: 0, Transparent: true, Flowable: true},
	BROWN_MUSHROOM:                {ID: BROWN_MUSHROOM, Name: "Brown Mushroom", LightLevel: 1, LightFilter: 0, Transparent: true, Flowable: true},
	RED_MUSHROOM:                  {ID: RED_MUSHROOM, Name: "Red Mushroom", LightFilter: 0, Transparent: true, Flowable: true},
	GOLD_BLOCK:                    {ID: GOLD_BLOCK, Name: "Gold Block", Hardness: 3, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	IRON_BLOCK:                    {ID: IRON_BLOCK, Name: "Iron Block", Hardness: 5, BlastResistance: 25, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	DOUBLE_SLAB:                   {ID: DOUBLE_SLAB, Name: "Double Stone Slab", Hardness: 2, BlastResistance: 10, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	SLAB:                          {ID: SLAB, Name: "Stone Slab", Hardness: 2, BlastResistance: 10, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe},
	BRICKS_BLOCK:                  {ID: BRICKS_BLOCK, Name: "Bricks", Hardness: 2, BlastResistance: 30, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	TNT:                           {ID: TNT, Name: "TNT", LightFilter: 15, Solid: true, FlammableChance: 15, BurnChance: 100},
	BOOKSHELF:                     {ID: BOOKSHELF, Name: "Bookshelf", Hardness: 1.5, BlastResistance: 7.5, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe, FlammableChance: 30, BurnChance: 20, FuelTime: 300},
	MOSS_STONE:                    {ID: MOSS_STONE, Name: "Moss Stone", Hardness: 2, BlastResistance: 10, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	OBSIDIAN:                      {ID: OBSIDIAN, Name: "Obsidian", Hardness: 35, BlastResistance: 175, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	TORCH:                         {ID: TORCH, Name: "Torch", LightLevel: 14, LightFilter: 0, Transparent: true, Flowable: true},
	FIRE:                          {ID: FIRE, Name: "Fire Block", LightLevel: 15, LightFilter: 0, Transparent: true, Replaceable: true, Flowable: true},
	MONSTER_SPAWNER:               {ID: MONSTER_SPAWNER, Name: "Monster Spawner", Hardness: 5, BlastResistance: 25, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	WOOD_STAIRS:                   {ID: WOOD_STAIRS, Name: "Wood Stairs", Hardness: 2, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FlammableChance: 5, BurnChance: 20, FuelTime: 300},
	CHEST:                         {ID: CHEST, Name: "Chest", Hardness: 2.5, BlastResistance: 12.5, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FuelTime: 300},
	REDSTONE_WIRE:                 {ID: REDSTONE_WIRE, Name: "Redstone Wire", LightFilter: 0, Transparent: true, Flowable: true},
	DIAMOND_ORE:                   {ID: DIAMOND_ORE, Name: "Diamond Ore", Hardness: 3, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	DIAMOND_BLOCK:                 {ID: DIAMOND_BLOCK, Name: "Diamond Block", Hardness: 5, BlastResistance: 25, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	WORKBENCH:                     {ID: WORKBENCH, Name: "Crafting Table", Hardness: 2.5, BlastResistance: 12.5, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe, FuelTime: 300},
	WHEAT_BLOCK:                   {ID: WHEAT_BLOCK, Name: "Wheat Block", LightFilter: 0, Transparent: true, Flowable: true},
	FARMLAND:                      {ID: FARMLAND, Name: "Farmland", Hardness: 0.6, BlastResistance: 3, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},
	FURNACE:                       {ID: FURNACE, Name: "Furnace", Hardness: 3.5, BlastResistance: 17.5, LightLevel: 13, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	BURNING_FURNACE:               {ID: BURNING_FURNACE, Name: "Burning Furnace", Hardness: 3.5, BlastResistance: 17.5, LightLevel: 13, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	SIGN_POST:                     {ID: SIGN_POST, Name: "Sign Post", Hardness: 1, BlastResistance: 5, LightFilter: 0, Transparent: true, ToolType: ToolTypeAxe},
	WOOD_DOOR_BLOCK:               {ID: WOOD_DOOR_BLOCK, Name: "Wood Door Block", Hardness: 3, BlastResistance: 15, LightFilter: 0, Transparent: true, ToolType: ToolTypeAxe},
	LADDER:                        {ID: LADDER, Name: "Ladder", Hardness: 0.4, BlastResistance: 2, LightFilter: 0, Transparent: true, ToolType: ToolTypeAxe},
	RAIL:                          {ID: RAIL, Name: "Rail", Hardness: 0.7, LightFilter: 0, Transparent: true, Flowable: true},
	COBBLESTONE_STAIRS:            {ID: COBBLESTONE_STAIRS, Name: "Cobblestone Stairs", Hardness: 2, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe, FlammableChance: 5, BurnChance: 20},
	WALL_SIGN:                     {ID: WALL_SIGN, Name: "Wall Sign", Hardness: 1, BlastResistance: 5, LightFilter: 0, Transparent: true, ToolType: ToolTypeAxe},
	LEVER:                         {ID: LEVER, Name: "Lever", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 0, Transparent: true, Flowable: true},
	STONE_PRESSURE_PLATE:          {ID: STONE_PRESSURE_PLATE, Name: "Stone Pressure Plate", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 0, Solid: true, Transparent: true},
	IRON_DOOR_BLOCK:               {ID: IRON_DOOR_BLOCK, Name: "Iron Door Block", Hardness: 5, BlastResistance: 25, LightFilter: 0, Transparent: true, ToolType: ToolTypePickaxe},
	WOODEN_PRESSURE_PLATE:         {ID: WOODEN_PRESSURE_PLATE, Name: "Wooden Pressure Plate", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 0, Solid: true, Transparent: true, FuelTime: 300},
	REDSTONE_ORE:                  {ID: REDSTONE_ORE, Name: "Redstone Ore", Hardness: 3, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	GLOWING_REDSTONE_ORE:          {ID: GLOWING_REDSTONE_ORE, Name: "Glowing Redstone Ore", Hardness: 3, BlastResistance: 15, LightLevel: 9, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	UNLIT_REDSTONE_TORCH:          {ID: UNLIT_REDSTONE_TORCH, Name: "Redstone Torch", LightFilter: 0, Transparent: true, Flowable: true},
	REDSTONE_TORCH:                {ID: REDSTONE_TORCH, Name: "Redstone Torch", LightLevel: 7, LightFilter: 0, Transparent: true, Flowable: true},
	STONE_BUTTON:                  {ID: STONE_BUTTON, Name: "Stone Button", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 0, Solid: true, Transparent: true},
	SNOW_LAYER:                    {ID: SNOW_LAYER, Name: "Snow Layer", Hardness: 0.1, LightFilter: 0, Transparent: true, Replaceable: true, Flowable: true, ToolType: ToolTypeShovel},
	ICE:                           {ID: ICE, Name: "Ice", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 2, Solid: true, Transparent: true, ToolType: ToolTypePickaxe},
	SNOW_BLOCK:                    {ID: SNOW_BLOCK, Name: "Snow Block", Hardness: 0.2, BlastResistance: 1, LightFilter: 15, Solid: true},
	CACTUS:                        {ID: CACTUS, Name: "Cactus", Hardness: 0.4, BlastResistance: 2, LightFilter: 0, Solid: true, Transparent: true},
	CLAY_BLOCK:                    {ID: CLAY_BLOCK, Name: "Clay Block", Hardness: 0.6, BlastResistance: 3, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},
	SUGARCANE_BLOCK:               {ID: SUGARCANE_BLOCK, Name: "Sugarcane", LightFilter: 0, Transparent: true, Flowable: true},
	FENCE:                         {ID: FENCE, Name: "Oak Fence", Hardness: 2, BlastResistance: 10, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FlammableChance: 5, BurnChance: 20, FuelTime: 300},
	PUMPKIN:                       {ID: PUMPKIN, Name: "Pumpkin", Hardness: 1, BlastResistance: 5, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe},
	NETHERRACK:                    {ID: NETHERRACK, Name: "Netherrack", Hardness: 0.4, BlastResistance: 2, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	SOUL_SAND:                     {ID: SOUL_SAND, Name: "Soul Sand", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},
	GLOWSTONE_BLOCK:               {ID: GLOWSTONE_BLOCK, Name: "Glowstone", Hardness: 0.3, BlastResistance: 1.5, LightLevel: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe},
	PORTAL:                        {ID: PORTAL, Name: "Portal", Hardness: -1, LightLevel: 11, LightFilter: 0, Transparent: true},
	LIT_PUMPKIN:                   {ID: LIT_PUMPKIN, Name: "Jack o'Lantern", Hardness: 1, BlastResistance: 5, LightLevel: 15, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe},
	CAKE_BLOCK:                    {ID: CAKE_BLOCK, Name: "Cake Block", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 0, Solid: true, Transparent: true},
	UNPOWERED_REPEATER:            {ID: UNPOWERED_REPEATER, Name: "Unpowered Repeater", LightFilter: 0, Transparent: true, Flowable: true},
	POWERED_REPEATER:              {ID: POWERED_REPEATER, Name: "Powered Repeater", LightFilter: 0, Transparent: true, Flowable: true},
	INVISIBLE_BEDROCK:             {ID: INVISIBLE_BEDROCK, Name: "Invisible Bedrock", Hardness: -1, BlastResistance: 18000000, LightFilter: 15, Solid: true},
	TRAPDOOR:                      {ID: TRAPDOOR, Name: "Wooden Trapdoor", Hardness: 3, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FuelTime: 300},
	STONE_BRICKS:                  {ID: STONE_BRICKS, Name: "Stone Bricks", Hardness: 1.5, BlastResistance: 7.5, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	BROWN_MUSHROOM_BLOCK:          {ID: BROWN_MUSHROOM_BLOCK, Name: "Brown Mushroom Block", Hardness: 0.2, BlastResistance: 1, LightFilter: 15, Solid: true},
	RED_MUSHROOM_BLOCK:            {ID: RED_MUSHROOM_BLOCK, Name: "Red Mushroom Block", Hardness: 0.2, BlastResistance: 1, LightFilter: 15, Solid: true},
	IRON_BARS:                     {ID: IRON_BARS, Name: "Iron Bars", Hardness: 5, BlastResistance: 25, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe},
	GLASS_PANE:                    {ID: GLASS_PANE, Name: "Glass Pane", Hardness: 0.3, BlastResistance: 1.5, LightFilter: 0, Solid: true, Transparent: true},
	MELON_BLOCK:                   {ID: MELON_BLOCK, Name: "Melon Block", Hardness: 1, BlastResistance: 5, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe},
	PUMPKIN_STEM:                  {ID: PUMPKIN_STEM, Name: "Pumpkin Stem", LightFilter: 0, Transparent: true, Flowable: true},
	MELON_STEM:                    {ID: MELON_STEM, Name: "Melon Stem", LightFilter: 0, Transparent: true, Flowable: true},
	VINE:                          {ID: VINE, Name: "Vines", Hardness: 0.2, BlastResistance: 1, LightFilter: 0, Transparent: true, ToolType: ToolTypeShears},
	FENCE_GATE:                    {ID: FENCE_GATE, Name: "Oak Fence Gate", Hardness: 2, BlastResistance: 10, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FuelTime: 300},
	BRICK_STAIRS:                  {ID: BRICK_STAIRS, Name: "Brick Stairs", Hardness: 2, BlastResistance: 30, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe, FlammableChance: 5, BurnChance: 20},
	STONE_BRICK_STAIRS:            {ID: STONE_BRICK_STAIRS, Name: "Stone Brick Stairs", Hardness: 1.5, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe, FlammableChance: 5, BurnChance: 20},
	MYCELIUM:                      {ID: MYCELIUM, Name: "Mycelium", Hardness: 0.6, BlastResistance: 3, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},
	WATER_LILY:                    {ID: WATER_LILY, Name: "Lily Pad", LightFilter: 0, Transparent: true, Flowable: true},
	NETHER_BRICKS:                 {ID: NETHER_BRICKS, Name: "Nether Bricks", Hardness: 2, BlastResistance: 10, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	NETHER_BRICK_FENCE:            {ID: NETHER_BRICK_FENCE, Name: "Nether Brick Fence", Hardness: 2, BlastResistance: 10, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe, FlammableChance: 5, BurnChance: 20, FuelTime: 300},
	NETHER_BRICKS_STAIRS:          {ID: NETHER_BRICKS_STAIRS, Name: "Nether Bricks Stairs", Hardness: 2, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe, FlammableChance: 5, BurnChance: 20},
	NETHER_WART_BLOCK:             {ID: NETHER_WART_BLOCK, Name: "Nether Wart Block", LightFilter: 0, Transparent: true, Flowable: true},
	ENCHANTING_TABLE:              {ID: ENCHANTING_TABLE, Name: "Enchanting Table", Hardness: 5, BlastResistance: 6000, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe},
	BREWING_STAND_BLOCK:           {ID: BREWING_STAND_BLOCK, Name: "Brewing Stand", Hardness: 0.5, BlastResistance: 2.5, LightLevel: 1, LightFilter: 0, Solid: true, Transparent: true},
	CAULDRON_BLOCK:                {ID: CAULDRON_BLOCK, Name: "Cauldron", Hardness: 2, BlastResistance: 10, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	END_PORTAL_FRAME:              {ID: END_PORTAL_FRAME, Name: "End Portal Frame", Hardness: -1, BlastResistance: 18000000, LightLevel: 1, LightFilter: 15, Solid: true},
	END_STONE:                     {ID: END_STONE, Name: "End Stone", Hardness: 3, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	INACTIVE_REDSTONE_LAMP:        {ID: INACTIVE_REDSTONE_LAMP, Name: "Redstone Lamp", Hardness: 0.3, BlastResistance: 1.5, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	ACTIVE_REDSTONE_LAMP:          {ID: ACTIVE_REDSTONE_LAMP, Name: "Active Redstone Lamp", Hardness: 0.3, BlastResistance: 1.5, LightLevel: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	DROPPER:                       {ID: DROPPER, Name: "Dropper", Hardness: 3.5, BlastResistance: 17.5, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	ACTIVATOR_RAIL:                {ID: ACTIVATOR_RAIL, Name: "Activator Rail", Hardness: 0.7, LightFilter: 0, Transparent: true, Flowable: true},
	COCOA_BLOCK:                   {ID: COCOA_BLOCK, Name: "Cocoa Block", Hardness: 0.2, BlastResistance: 15, LightFilter: 0, Transparent: true, Flowable: true},
	SANDSTONE_STAIRS:              {ID: SANDSTONE_STAIRS, Name: "Sandstone Stairs", Hardness: 0.8, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe, FlammableChance: 5, BurnChance: 20},
	EMERALD_ORE:                   {ID: EMERALD_ORE, Name: "Emerald Ore", Hardness: 3, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	TRIPWIRE_HOOK:                 {ID: TRIPWIRE_HOOK, Name: "Tripwire Hook", LightFilter: 0, Solid: true, Transparent: true},
	TRIPWIRE:                      {ID: TRIPWIRE, Name: "Tripwire", LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeShears},
	EMERALD_BLOCK:                 {ID: EMERALD_BLOCK, Name: "Emerald Block", Hardness: 5, BlastResistance: 25, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	SPRUCE_WOOD_STAIRS:            {ID: SPRUCE_WOOD_STAIRS, Name: "Spruce Wood Stairs", Hardness: 2, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FlammableChance: 5, BurnChance: 20, FuelTime: 300},
	BIRCH_WOOD_STAIRS:             {ID: BIRCH_WOOD_STAIRS, Name: "Birch Wood Stairs", Hardness: 2, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FlammableChance: 5, BurnChance: 20, FuelTime: 300},
	JUNGLE_WOOD_STAIRS:            {ID: JUNGLE_WOOD_STAIRS, Name: "Jungle Wood Stairs", Hardness: 2, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FlammableChance: 5, BurnChance: 20, FuelTime: 300},
	STONE_WALL:                    {ID: STONE_WALL, Name: "Cobblestone Wall", Hardness: 2, BlastResistance: 10, LightFilter: 0, Transparent: true, ToolType: ToolTypePickaxe},
	FLOWER_POT_BLOCK:              {ID: FLOWER_POT_BLOCK, Name: "Flower Pot Block", LightFilter: 0, Transparent: true, Flowable: true},
	CARROT_BLOCK:                  {ID: CARROT_BLOCK, Name: "Carrot Block", LightFilter: 0, Transparent: true, Flowable: true},
	POTATO_BLOCK:                  {ID: POTATO_BLOCK, Name: "Potato Block", LightFilter: 0, Transparent: true, Flowable: true},
	WOODEN_BUTTON:                 {ID: WOODEN_BUTTON, Name: "Wooden Button", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 0, Solid: true, Transparent: true},
	SKULL_BLOCK:                   {ID: SKULL_BLOCK, Name: "Skeleton Skull", Hardness: 1, BlastResistance: 5, LightFilter: 0, Transparent: true, ToolType: ToolTypePickaxe},
	ANVIL:                         {ID: ANVIL, Name: "Anvil", Hardness: 5, BlastResistance: 6000, LightFilter: 15, ToolType: ToolTypePickaxe},
	TRAPPED_CHEST:                 {ID: TRAPPED_CHEST, Name: "Trapped Chest", Hardness: 2.5, BlastResistance: 12.5, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FuelTime: 300},
	LIGHT_WEIGHTED_PRESSURE_PLATE: {ID: LIGHT_WEIGHTED_PRESSURE_PLATE, Name: "Light Weighted Pressure Plate", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 0, Solid: true, Transparent: true},
	HEAVY_WEIGHTED_PRESSURE_PLATE: {ID: HEAVY_WEIGHTED_PRESSURE_PLATE, Name: "Heavy Weighted Pressure Plate", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 0, Solid: true, Transparent: true},
	UNPOWERED_COMPARATOR_BLOCK:    {ID: UNPOWERED_COMPARATOR_BLOCK, Name: "Unpowered Comparator", LightFilter: 0, Transparent: true, Flowable: true},
	DAYLIGHT_SENSOR:               {ID: DAYLIGHT_SENSOR, Name: "Daylight Sensor", Hardness: 0.2, BlastResistance: 1, LightFilter: 0, Solid: true, Transparent: true, FuelTime: 300},
	REDSTONE_BLOCK:                {ID: REDSTONE_BLOCK, Name: "Redstone Block", Hardness: 5, BlastResistance: 25, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	NETHER_QUARTZ_ORE:             {ID: NETHER_QUARTZ_ORE, Name: "Nether Quartz Ore", Hardness: 3, BlastResistance: 15, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	HOPPER_BLOCK:                  {ID: HOPPER_BLOCK, Name: "Hopper", Hardness: 3, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe},
	QUARTZ_BLOCK:                  {ID: QUARTZ_BLOCK, Name: "Quartz Block", Hardness: 0.8, BlastResistance: 4, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	QUARTZ_STAIRS:                 {ID: QUARTZ_STAIRS, Name: "Quartz Stairs", Hardness: 0.8, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe, FlammableChance: 5, BurnChance: 20},
	DOUBLE_WOOD_SLAB:              {ID: DOUBLE_WOOD_SLAB, Name: "Double Oak Wooden Slab", Hardness: 2, BlastResistance: 10, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe},
	WOOD_SLAB:                     {ID: WOOD_SLAB, Name: "Oak Wooden Slab", Hardness: 2, BlastResistance: 10, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe},
	STAINED_CLAY:                  {ID: STAINED_CLAY, Name: "White Stained Clay", Hardness: 1.25, BlastResistance: 6.25, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	LEAVES2:                       {ID: LEAVES2, Name: "Acacia Leaves", Hardness: 0.2, BlastResistance: 1, LightFilter: 0, Solid: true, Transparent: true, DiffusesSkyLight: true, ToolType: ToolTypeShears, FlammableChance: 30, BurnChance: 60},
	WOOD2:                         {ID: WOOD2, Name: "Acacia Wood", Hardness: 2, BlastResistance: 10, LightFilter: 15, Solid: true, ToolType: ToolTypeAxe, FlammableChance: 5, BurnChance: 10, FuelTime: 300},
	ACACIA_WOOD_STAIRS:            {ID: ACACIA_WOOD_STAIRS, Name: "Acacia Wood Stairs", Hardness: 2, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FlammableChance: 5, BurnChance: 20, FuelTime: 300},
	DARK_OAK_WOOD_STAIRS:          {ID: DARK_OAK_WOOD_STAIRS, Name: "Dark Oak Wood Stairs", Hardness: 2, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FlammableChance: 5, BurnChance: 20, FuelTime: 300},
	SLIME_BLOCK:                   {ID: SLIME_BLOCK, Name: "Slime Block", LightFilter: 15, Solid: true},
	IRON_TRAPDOOR:                 {ID: IRON_TRAPDOOR, Name: "Iron Trapdoor", Hardness: 5, BlastResistance: 25, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe},
	HAY_BALE:                      {ID: HAY_BALE, Name: "Hay Bale", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 15, Solid: true, FlammableChance: 60, BurnChance: 20},
	CARPET:                        {ID: CARPET, Name: "White Carpet", Hardness: 0.1, LightFilter: 0, Solid: true, Transparent: true, Flowable: true},
	HARDENED_CLAY:                 {ID: HARDENED_CLAY, Name: "Hardened Clay", Hardness: 1.25, BlastResistance: 6.25, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	COAL_BLOCK:                    {ID: COAL_BLOCK, Name: "Coal Block", Hardness: 5, BlastResistance: 25, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe, FlammableChance: 5, BurnChance: 5, FuelTime: 16000},
	PACKED_ICE:                    {ID: PACKED_ICE, Name: "Packed Ice", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	DOUBLE_PLANT:                  {ID: DOUBLE_PLANT, Name: "Sunflower", LightFilter: 0, Transparent: true, Replaceable: true, Flowable: true},
	DAYLIGHT_SENSOR_INVERTED:      {ID: DAYLIGHT_SENSOR_INVERTED, Name: "Daylight Sensor", Hardness: 0.2, BlastResistance: 1, LightFilter: 0, Solid: true, Transparent: true, FuelTime: 300},
	RED_SANDSTONE:                 {ID: RED_SANDSTONE, Name: "Red Sandstone", Hardness: 0.8, BlastResistance: 4, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	RED_SANDSTONE_STAIRS:          {ID: RED_SANDSTONE_STAIRS, Name: "Red Sandstone Stairs", Hardness: 0.8, BlastResistance: 15, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe, FlammableChance: 5, BurnChance: 20},
	DOUBLE_RED_SANDSTONE_SLAB:     {ID: DOUBLE_RED_SANDSTONE_SLAB, Name: "Double Red Sandstone Slab", Hardness: 2, BlastResistance: 10, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	RED_SANDSTONE_SLAB:            {ID: RED_SANDSTONE_SLAB, Name: "Red Sandstone Slab", Hardness: 2, BlastResistance: 10, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypePickaxe},
	FENCE_GATE_SPRUCE:             {ID: FENCE_GATE_SPRUCE, Name: "Spruce Fence Gate", Hardness: 2, BlastResistance: 10, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FuelTime: 300},
	FENCE_GATE_BIRCH:              {ID: FENCE_GATE_BIRCH, Name: "Birch Fence Gate", Hardness: 2, BlastResistance: 10, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FuelTime: 300},
	FENCE_GATE_JUNGLE:             {ID: FENCE_GATE_JUNGLE, Name: "Jungle Fence Gate", Hardness: 2, BlastResistance: 10, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FuelTime: 300},
	FENCE_GATE_DARK_OAK:           {ID: FENCE_GATE_DARK_OAK, Name: "Dark Oak Fence Gate", Hardness: 2, BlastResistance: 10, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FuelTime: 300},
	FENCE_GATE_ACACIA:             {ID: FENCE_GATE_ACACIA, Name: "Acacia Fence Gate", Hardness: 2, BlastResistance: 10, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeAxe, FuelTime: 300},
	SPRUCE_DOOR_BLOCK:             {ID: SPRUCE_DOOR_BLOCK, Name: "Spruce Door Block", Hardness: 3, BlastResistance: 15, LightFilter: 0, Transparent: true, ToolType: ToolTypeAxe},
	BIRCH_DOOR_BLOCK:              {ID: BIRCH_DOOR_BLOCK, Name: "Birch Door Block", Hardness: 3, BlastResistance: 15, LightFilter: 0, Transparent: true, ToolType: ToolTypeAxe},
	JUNGLE_DOOR_BLOCK:             {ID: JUNGLE_DOOR_BLOCK, Name: "Jungle Door Block", Hardness: 3, BlastResistance: 15, LightFilter: 0, Transparent: true, ToolType: ToolTypeAxe},
	ACACIA_DOOR_BLOCK:             {ID: ACACIA_DOOR_BLOCK, Name: "Acacia Door Block", Hardness: 3, BlastResistance: 15, LightFilter: 0, Transparent: true, ToolType: ToolTypeAxe},
	DARK_OAK_DOOR_BLOCK:           {ID: DARK_OAK_DOOR_BLOCK, Name: "Dark Oak Door Block", Hardness: 3, BlastResistance: 15, LightFilter: 0, Transparent: true, ToolType: ToolTypeAxe},
	GRASS_PATH:                    {ID: GRASS_PATH, Name: "Grass Path", Hardness: 0.6, BlastResistance: 3, LightFilter: 0, Solid: true, Transparent: true, ToolType: ToolTypeShovel},
	ITEM_FRAME_BLOCK:              {ID: ITEM_FRAME_BLOCK, Name: "Item Frame", Hardness: 0.25, BlastResistance: 1.25, LightFilter: 0, Solid: true, Transparent: true},
	PODZOL:                        {ID: PODZOL, Name: "Podzol", Hardness: 0.5, BlastResistance: 2.5, LightFilter: 15, Solid: true, ToolType: ToolTypeShovel},
	BEETROOT_BLOCK:                {ID: BEETROOT_BLOCK, Name: "Beetroot Block", LightFilter: 0, Transparent: true, Flowable: true},
	STONECUTTER:                   {ID: STONECUTTER, Name: "Stonecutter", Hardness: 2.5, BlastResistance: 12.5, LightFilter: 15, Solid: true, ToolType: ToolTypePickaxe},
	GLOWING_OBSIDIAN:              {ID: GLOWING_OBSIDIAN, Name: "Glowing Obsidian", Hardness: 10, BlastResistance: 50, LightLevel: 12, LightFilter: 15, Solid: true},
	NETHER_REACTOR:                {ID: NETHER_REACTOR, Name: "Nether Reactor", Hardness: 3, BlastResistance: 15, LightFilter: 15, Solid: true},
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
