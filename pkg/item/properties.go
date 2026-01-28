package item

type ItemProperty struct {
	ID            int
	Name          string
	MaxStackSize  int
	MaxDurability int
	FoodRestore   int
	Saturation    float32
}

var defaultProperty = ItemProperty{
	Name:         "Unknown",
	MaxStackSize: 64,
}

func GetItemName(id int) string {
	if prop, ok := itemProperties[id]; ok {
		return prop.Name
	}

	if id < 256 {
		return "Block"
	}
	return "Unknown Item"
}

func GetMaxStackSizeFor(id int) int {
	if prop, ok := itemProperties[id]; ok && prop.MaxStackSize > 0 {
		return prop.MaxStackSize
	}
	if IsTool(id) || IsArmor(id) {
		return 1
	}
	return 64
}

func IsFood(id int) bool {
	if prop, ok := itemProperties[id]; ok {
		return prop.FoodRestore > 0
	}
	return false
}

func GetFoodRestore(id int) int {
	if prop, ok := itemProperties[id]; ok {
		return prop.FoodRestore
	}
	return 0
}

func GetSaturation(id int) float32 {
	if prop, ok := itemProperties[id]; ok {
		return prop.Saturation
	}
	return 0
}

func IsArmor(id int) bool {
	return id >= LEATHER_CAP && id <= GOLD_BOOTS
}

func GetArmorType(id int) int {
	if !IsArmor(id) {
		return -1
	}
	switch id {
	case LEATHER_CAP, CHAIN_HELMET, IRON_HELMET, DIAMOND_HELMET, GOLD_HELMET:
		return 0
	case LEATHER_TUNIC, CHAIN_CHESTPLATE, IRON_CHESTPLATE, DIAMOND_CHESTPLATE, GOLD_CHESTPLATE:
		return 1
	case LEATHER_PANTS, CHAIN_LEGGINGS, IRON_LEGGINGS, DIAMOND_LEGGINGS, GOLD_LEGGINGS:
		return 2
	case LEATHER_BOOTS, CHAIN_BOOTS, IRON_BOOTS, DIAMOND_BOOTS, GOLD_BOOTS:
		return 3
	}
	return -1
}

func GetArmorDefense(id int) int {
	switch id {

	case LEATHER_CAP:
		return 1
	case LEATHER_TUNIC:
		return 3
	case LEATHER_PANTS:
		return 2
	case LEATHER_BOOTS:
		return 1

	case CHAIN_HELMET:
		return 2
	case CHAIN_CHESTPLATE:
		return 5
	case CHAIN_LEGGINGS:
		return 4
	case CHAIN_BOOTS:
		return 1

	case IRON_HELMET:
		return 2
	case IRON_CHESTPLATE:
		return 6
	case IRON_LEGGINGS:
		return 5
	case IRON_BOOTS:
		return 2

	case DIAMOND_HELMET:
		return 3
	case DIAMOND_CHESTPLATE:
		return 8
	case DIAMOND_LEGGINGS:
		return 6
	case DIAMOND_BOOTS:
		return 3

	case GOLD_HELMET:
		return 2
	case GOLD_CHESTPLATE:
		return 5
	case GOLD_LEGGINGS:
		return 3
	case GOLD_BOOTS:
		return 1
	}
	return 0
}

var itemProperties = map[int]ItemProperty{

	IRON_SHOVEL:     {ID: IRON_SHOVEL, Name: "Iron Shovel", MaxStackSize: 1, MaxDurability: 251},
	IRON_PICKAXE:    {ID: IRON_PICKAXE, Name: "Iron Pickaxe", MaxStackSize: 1, MaxDurability: 251},
	IRON_AXE:        {ID: IRON_AXE, Name: "Iron Axe", MaxStackSize: 1, MaxDurability: 251},
	IRON_SWORD:      {ID: IRON_SWORD, Name: "Iron Sword", MaxStackSize: 1, MaxDurability: 251},
	IRON_HOE:        {ID: IRON_HOE, Name: "Iron Hoe", MaxStackSize: 1, MaxDurability: 251},
	WOODEN_SWORD:    {ID: WOODEN_SWORD, Name: "Wooden Sword", MaxStackSize: 1, MaxDurability: 60},
	WOODEN_SHOVEL:   {ID: WOODEN_SHOVEL, Name: "Wooden Shovel", MaxStackSize: 1, MaxDurability: 60},
	WOODEN_PICKAXE:  {ID: WOODEN_PICKAXE, Name: "Wooden Pickaxe", MaxStackSize: 1, MaxDurability: 60},
	WOODEN_AXE:      {ID: WOODEN_AXE, Name: "Wooden Axe", MaxStackSize: 1, MaxDurability: 60},
	WOODEN_HOE:      {ID: WOODEN_HOE, Name: "Wooden Hoe", MaxStackSize: 1, MaxDurability: 60},
	STONE_SWORD:     {ID: STONE_SWORD, Name: "Stone Sword", MaxStackSize: 1, MaxDurability: 132},
	STONE_SHOVEL:    {ID: STONE_SHOVEL, Name: "Stone Shovel", MaxStackSize: 1, MaxDurability: 132},
	STONE_PICKAXE:   {ID: STONE_PICKAXE, Name: "Stone Pickaxe", MaxStackSize: 1, MaxDurability: 132},
	STONE_AXE:       {ID: STONE_AXE, Name: "Stone Axe", MaxStackSize: 1, MaxDurability: 132},
	STONE_HOE:       {ID: STONE_HOE, Name: "Stone Hoe", MaxStackSize: 1, MaxDurability: 132},
	DIAMOND_SWORD:   {ID: DIAMOND_SWORD, Name: "Diamond Sword", MaxStackSize: 1, MaxDurability: 1562},
	DIAMOND_SHOVEL:  {ID: DIAMOND_SHOVEL, Name: "Diamond Shovel", MaxStackSize: 1, MaxDurability: 1562},
	DIAMOND_PICKAXE: {ID: DIAMOND_PICKAXE, Name: "Diamond Pickaxe", MaxStackSize: 1, MaxDurability: 1562},
	DIAMOND_AXE:     {ID: DIAMOND_AXE, Name: "Diamond Axe", MaxStackSize: 1, MaxDurability: 1562},
	DIAMOND_HOE:     {ID: DIAMOND_HOE, Name: "Diamond Hoe", MaxStackSize: 1, MaxDurability: 1562},
	GOLD_SWORD:      {ID: GOLD_SWORD, Name: "Golden Sword", MaxStackSize: 1, MaxDurability: 33},
	GOLD_SHOVEL:     {ID: GOLD_SHOVEL, Name: "Golden Shovel", MaxStackSize: 1, MaxDurability: 33},
	GOLD_PICKAXE:    {ID: GOLD_PICKAXE, Name: "Golden Pickaxe", MaxStackSize: 1, MaxDurability: 33},
	GOLD_AXE:        {ID: GOLD_AXE, Name: "Golden Axe", MaxStackSize: 1, MaxDurability: 33},
	GOLD_HOE:        {ID: GOLD_HOE, Name: "Golden Hoe", MaxStackSize: 1, MaxDurability: 33},

	FLINT_AND_STEEL: {ID: FLINT_AND_STEEL, Name: "Flint and Steel", MaxStackSize: 1, MaxDurability: 65},
	BOW:             {ID: BOW, Name: "Bow", MaxStackSize: 1, MaxDurability: 385},
	SHEARS:          {ID: SHEARS, Name: "Shears", MaxStackSize: 1, MaxDurability: 239},
	FISHING_ROD:     {ID: FISHING_ROD, Name: "Fishing Rod", MaxStackSize: 1, MaxDurability: 65},

	LEATHER_CAP:        {ID: LEATHER_CAP, Name: "Leather Cap", MaxStackSize: 1, MaxDurability: 56},
	LEATHER_TUNIC:      {ID: LEATHER_TUNIC, Name: "Leather Tunic", MaxStackSize: 1, MaxDurability: 81},
	LEATHER_PANTS:      {ID: LEATHER_PANTS, Name: "Leather Pants", MaxStackSize: 1, MaxDurability: 76},
	LEATHER_BOOTS:      {ID: LEATHER_BOOTS, Name: "Leather Boots", MaxStackSize: 1, MaxDurability: 66},
	CHAIN_HELMET:       {ID: CHAIN_HELMET, Name: "Chain Helmet", MaxStackSize: 1, MaxDurability: 166},
	CHAIN_CHESTPLATE:   {ID: CHAIN_CHESTPLATE, Name: "Chain Chestplate", MaxStackSize: 1, MaxDurability: 241},
	CHAIN_LEGGINGS:     {ID: CHAIN_LEGGINGS, Name: "Chain Leggings", MaxStackSize: 1, MaxDurability: 226},
	CHAIN_BOOTS:        {ID: CHAIN_BOOTS, Name: "Chain Boots", MaxStackSize: 1, MaxDurability: 196},
	IRON_HELMET:        {ID: IRON_HELMET, Name: "Iron Helmet", MaxStackSize: 1, MaxDurability: 166},
	IRON_CHESTPLATE:    {ID: IRON_CHESTPLATE, Name: "Iron Chestplate", MaxStackSize: 1, MaxDurability: 241},
	IRON_LEGGINGS:      {ID: IRON_LEGGINGS, Name: "Iron Leggings", MaxStackSize: 1, MaxDurability: 226},
	IRON_BOOTS:         {ID: IRON_BOOTS, Name: "Iron Boots", MaxStackSize: 1, MaxDurability: 196},
	DIAMOND_HELMET:     {ID: DIAMOND_HELMET, Name: "Diamond Helmet", MaxStackSize: 1, MaxDurability: 364},
	DIAMOND_CHESTPLATE: {ID: DIAMOND_CHESTPLATE, Name: "Diamond Chestplate", MaxStackSize: 1, MaxDurability: 529},
	DIAMOND_LEGGINGS:   {ID: DIAMOND_LEGGINGS, Name: "Diamond Leggings", MaxStackSize: 1, MaxDurability: 496},
	DIAMOND_BOOTS:      {ID: DIAMOND_BOOTS, Name: "Diamond Boots", MaxStackSize: 1, MaxDurability: 430},
	GOLD_HELMET:        {ID: GOLD_HELMET, Name: "Golden Helmet", MaxStackSize: 1, MaxDurability: 78},
	GOLD_CHESTPLATE:    {ID: GOLD_CHESTPLATE, Name: "Golden Chestplate", MaxStackSize: 1, MaxDurability: 113},
	GOLD_LEGGINGS:      {ID: GOLD_LEGGINGS, Name: "Golden Leggings", MaxStackSize: 1, MaxDurability: 106},
	GOLD_BOOTS:         {ID: GOLD_BOOTS, Name: "Golden Boots", MaxStackSize: 1, MaxDurability: 92},

	COAL:           {ID: COAL, Name: "Coal"},
	DIAMOND:        {ID: DIAMOND, Name: "Diamond"},
	IRON_INGOT:     {ID: IRON_INGOT, Name: "Iron Ingot"},
	GOLD_INGOT:     {ID: GOLD_INGOT, Name: "Gold Ingot"},
	STICK:          {ID: STICK, Name: "Stick"},
	STRING:         {ID: STRING, Name: "String"},
	FEATHER:        {ID: FEATHER, Name: "Feather"},
	GUNPOWDER:      {ID: GUNPOWDER, Name: "Gunpowder"},
	FLINT:          {ID: FLINT, Name: "Flint"},
	LEATHER:        {ID: LEATHER, Name: "Leather"},
	BRICK:          {ID: BRICK, Name: "Brick"},
	CLAY:           {ID: CLAY, Name: "Clay Ball"},
	PAPER:          {ID: PAPER, Name: "Paper"},
	BOOK:           {ID: BOOK, Name: "Book"},
	SLIMEBALL:      {ID: SLIMEBALL, Name: "Slimeball"},
	BONE:           {ID: BONE, Name: "Bone"},
	REDSTONE:       {ID: REDSTONE, Name: "Redstone"},
	GLOWSTONE_DUST: {ID: GLOWSTONE_DUST, Name: "Glowstone Dust"},
	EMERALD:        {ID: EMERALD, Name: "Emerald"},
	GOLD_NUGGET:    {ID: GOLD_NUGGET, Name: "Gold Nugget"},
	QUARTZ:         {ID: QUARTZ, Name: "Nether Quartz"},

	APPLE:            {ID: APPLE, Name: "Apple", FoodRestore: 4, Saturation: 2.4},
	BREAD:            {ID: BREAD, Name: "Bread", FoodRestore: 5, Saturation: 6.0},
	RAW_PORKCHOP:     {ID: RAW_PORKCHOP, Name: "Raw Porkchop", FoodRestore: 3, Saturation: 1.8},
	COOKED_PORKCHOP:  {ID: COOKED_PORKCHOP, Name: "Cooked Porkchop", FoodRestore: 8, Saturation: 12.8},
	GOLDEN_APPLE:     {ID: GOLDEN_APPLE, Name: "Golden Apple", FoodRestore: 4, Saturation: 9.6},
	RAW_FISH:         {ID: RAW_FISH, Name: "Raw Fish", FoodRestore: 2, Saturation: 0.4},
	COOKED_FISH:      {ID: COOKED_FISH, Name: "Cooked Fish", FoodRestore: 5, Saturation: 6.0},
	COOKIE:           {ID: COOKIE, Name: "Cookie", FoodRestore: 2, Saturation: 0.4},
	MELON:            {ID: MELON, Name: "Melon Slice", FoodRestore: 2, Saturation: 1.2},
	RAW_BEEF:         {ID: RAW_BEEF, Name: "Raw Beef", FoodRestore: 3, Saturation: 1.8},
	STEAK:            {ID: STEAK, Name: "Steak", FoodRestore: 8, Saturation: 12.8},
	RAW_CHICKEN:      {ID: RAW_CHICKEN, Name: "Raw Chicken", FoodRestore: 2, Saturation: 1.2},
	COOKED_CHICKEN:   {ID: COOKED_CHICKEN, Name: "Cooked Chicken", FoodRestore: 6, Saturation: 7.2},
	ROTTEN_FLESH:     {ID: ROTTEN_FLESH, Name: "Rotten Flesh", FoodRestore: 4, Saturation: 0.8},
	CARROT:           {ID: CARROT, Name: "Carrot", FoodRestore: 3, Saturation: 3.6},
	POTATO:           {ID: POTATO, Name: "Potato", FoodRestore: 1, Saturation: 0.6},
	BAKED_POTATO:     {ID: BAKED_POTATO, Name: "Baked Potato", FoodRestore: 5, Saturation: 6.0},
	POISONOUS_POTATO: {ID: POISONOUS_POTATO, Name: "Poisonous Potato", FoodRestore: 2, Saturation: 1.2},
	GOLDEN_CARROT:    {ID: GOLDEN_CARROT, Name: "Golden Carrot", FoodRestore: 6, Saturation: 14.4},
	PUMPKIN_PIE:      {ID: PUMPKIN_PIE, Name: "Pumpkin Pie", FoodRestore: 8, Saturation: 4.8},
	RAW_RABBIT:       {ID: RAW_RABBIT, Name: "Raw Rabbit", FoodRestore: 3, Saturation: 1.8},
	COOKED_RABBIT:    {ID: COOKED_RABBIT, Name: "Cooked Rabbit", FoodRestore: 5, Saturation: 6.0},
	RABBIT_STEW:      {ID: RABBIT_STEW, Name: "Rabbit Stew", MaxStackSize: 1, FoodRestore: 10, Saturation: 12.0},
	BEETROOT:         {ID: BEETROOT, Name: "Beetroot", FoodRestore: 1, Saturation: 1.2},
	BEETROOT_SOUP:    {ID: BEETROOT_SOUP, Name: "Beetroot Soup", MaxStackSize: 1, FoodRestore: 6, Saturation: 7.2},
	MUSHROOM_STEW:    {ID: MUSHROOM_STEW, Name: "Mushroom Stew", MaxStackSize: 1, FoodRestore: 6, Saturation: 7.2},

	ARROW:        {ID: ARROW, Name: "Arrow"},
	BOWL:         {ID: BOWL, Name: "Bowl"},
	EGG:          {ID: EGG, Name: "Egg", MaxStackSize: 16},
	SNOWBALL:     {ID: SNOWBALL, Name: "Snowball", MaxStackSize: 16},
	BUCKET:       {ID: BUCKET, Name: "Bucket", MaxStackSize: 1},
	SIGN:         {ID: SIGN, Name: "Sign", MaxStackSize: 16},
	BED:          {ID: BED, Name: "Bed", MaxStackSize: 1},
	CAKE:         {ID: CAKE, Name: "Cake", MaxStackSize: 1},
	COMPASS:      {ID: COMPASS, Name: "Compass"},
	CLOCK:        {ID: CLOCK, Name: "Clock"},
	DYE:          {ID: DYE, Name: "Dye"},
	SUGAR:        {ID: SUGAR, Name: "Sugar"},
	SEEDS:        {ID: SEEDS, Name: "Seeds"},
	POTION:       {ID: POTION, Name: "Potion", MaxStackSize: 1},
	GLASS_BOTTLE: {ID: GLASS_BOTTLE, Name: "Glass Bottle"},
	SPAWN_EGG:    {ID: SPAWN_EGG, Name: "Spawn Egg"},
	MAP:          {ID: MAP, Name: "Empty Map"},
	FILLED_MAP:   {ID: FILLED_MAP, Name: "Map"},
}
