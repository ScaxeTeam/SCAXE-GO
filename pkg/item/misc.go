package item
type MiscItemInfo struct {
	ID           int
	Name         string
	MaxStack     int
	FuelTime     int
	IsPlaceable  bool
	PlaceBlockID int
}
var miscItems = map[int]MiscItemInfo{
	BUCKET:          {ID: BUCKET, Name: "Bucket", MaxStack: 16},
	FLINT_AND_STEEL: {ID: FLINT_AND_STEEL, Name: "Flint and Steel", MaxStack: 1},
	COMPASS:         {ID: COMPASS, Name: "Compass", MaxStack: 64},
	CLOCK:           {ID: CLOCK, Name: "Clock", MaxStack: 64},
	SHEARS:          {ID: SHEARS, Name: "Shears", MaxStack: 1},
	FISHING_ROD:     {ID: FISHING_ROD, Name: "Fishing Rod", MaxStack: 1},
	SADDLE:          {ID: SADDLE, Name: "Saddle", MaxStack: 1},
	FILLED_MAP:      {ID: FILLED_MAP, Name: "Map", MaxStack: 64},
	MAP:             {ID: MAP, Name: "Empty Map", MaxStack: 64},
	PAINTING:   {ID: PAINTING, Name: "Painting", MaxStack: 64, IsPlaceable: true},
	SIGN:       {ID: SIGN, Name: "Sign", MaxStack: 16, IsPlaceable: true, FuelTime: 200},
	ITEM_FRAME: {ID: ITEM_FRAME, Name: "Item Frame", MaxStack: 64, IsPlaceable: true},
	FLOWER_POT: {ID: FLOWER_POT, Name: "Flower Pot", MaxStack: 64, IsPlaceable: true},
	BED:        {ID: BED, Name: "Bed", MaxStack: 1, IsPlaceable: true},
	CAKE:       {ID: CAKE, Name: "Cake", MaxStack: 1, IsPlaceable: true},
	BOOK:           {ID: BOOK, Name: "Book", MaxStack: 64},
	ENCHANTED_BOOK: {ID: ENCHANTED_BOOK, Name: "Enchanted Book", MaxStack: 1},
	MINECART:        {ID: MINECART, Name: "Minecart", MaxStack: 1},
	CHEST_MINECART:  {ID: CHEST_MINECART, Name: "Chest Minecart", MaxStack: 1},
	TNT_MINECART:    {ID: TNT_MINECART, Name: "TNT Minecart", MaxStack: 1},
	HOPPER_MINECART: {ID: HOPPER_MINECART, Name: "Hopper Minecart", MaxStack: 1},
	BOAT:            {ID: BOAT, Name: "Oak Boat", MaxStack: 1},
	SPRUCE_BOAT:     {ID: SPRUCE_BOAT, Name: "Spruce Boat", MaxStack: 1},
	BIRCH_BOAT:      {ID: BIRCH_BOAT, Name: "Birch Boat", MaxStack: 1},
	JUNGLE_BOAT:     {ID: JUNGLE_BOAT, Name: "Jungle Boat", MaxStack: 1},
	ACACIA_BOAT:     {ID: ACACIA_BOAT, Name: "Acacia Boat", MaxStack: 1},
	DARK_OAK_BOAT:   {ID: DARK_OAK_BOAT, Name: "Dark Oak Boat", MaxStack: 1},
	COAL:          {ID: COAL, Name: "Coal", MaxStack: 64, FuelTime: 1600},
	DIAMOND:       {ID: DIAMOND, Name: "Diamond", MaxStack: 64},
	IRON_INGOT:    {ID: IRON_INGOT, Name: "Iron Ingot", MaxStack: 64},
	GOLD_INGOT:    {ID: GOLD_INGOT, Name: "Gold Ingot", MaxStack: 64},
	EMERALD:       {ID: EMERALD, Name: "Emerald", MaxStack: 64},
	STICK:         {ID: STICK, Name: "Stick", MaxStack: 64, FuelTime: 100},
	STRING:        {ID: STRING, Name: "String", MaxStack: 64},
	FEATHER:       {ID: FEATHER, Name: "Feather", MaxStack: 64},
	FLINT:         {ID: FLINT, Name: "Flint", MaxStack: 64},
	GUNPOWDER:     {ID: GUNPOWDER, Name: "Gunpowder", MaxStack: 64},
	LEATHER:       {ID: LEATHER, Name: "Leather", MaxStack: 64},
	BRICK:         {ID: BRICK, Name: "Brick", MaxStack: 64},
	CLAY:          {ID: CLAY, Name: "Clay Ball", MaxStack: 64},
	PAPER:         {ID: PAPER, Name: "Paper", MaxStack: 64},
	SUGAR:         {ID: SUGAR, Name: "Sugar", MaxStack: 64},
	SLIMEBALL:     {ID: SLIMEBALL, Name: "Slimeball", MaxStack: 64},
	BONE:          {ID: BONE, Name: "Bone", MaxStack: 64},
	GOLD_NUGGET:   {ID: GOLD_NUGGET, Name: "Gold Nugget", MaxStack: 64},
	NETHER_QUARTZ: {ID: NETHER_QUARTZ, Name: "Nether Quartz", MaxStack: 64},
	NETHER_BRICK:  {ID: NETHER_BRICK, Name: "Nether Brick", MaxStack: 64},
	REDSTONE:   {ID: REDSTONE, Name: "Redstone Dust", MaxStack: 64},
	REPEATER:   {ID: REPEATER, Name: "Repeater", MaxStack: 64, IsPlaceable: true},
	COMPARATOR: {ID: COMPARATOR, Name: "Comparator", MaxStack: 64, IsPlaceable: true},
	WOODEN_DOOR:   {ID: WOODEN_DOOR, Name: "Oak Door", MaxStack: 64, IsPlaceable: true},
	IRON_DOOR:     {ID: IRON_DOOR, Name: "Iron Door", MaxStack: 64, IsPlaceable: true},
	SPRUCE_DOOR:   {ID: SPRUCE_DOOR, Name: "Spruce Door", MaxStack: 64, IsPlaceable: true},
	BIRCH_DOOR:    {ID: BIRCH_DOOR, Name: "Birch Door", MaxStack: 64, IsPlaceable: true},
	JUNGLE_DOOR:   {ID: JUNGLE_DOOR, Name: "Jungle Door", MaxStack: 64, IsPlaceable: true},
	ACACIA_DOOR:   {ID: ACACIA_DOOR, Name: "Acacia Door", MaxStack: 64, IsPlaceable: true},
	DARK_OAK_DOOR: {ID: DARK_OAK_DOOR, Name: "Dark Oak Door", MaxStack: 64, IsPlaceable: true},
	GLASS_BOTTLE:         {ID: GLASS_BOTTLE, Name: "Glass Bottle", MaxStack: 64},
	POTION:               {ID: POTION, Name: "Potion", MaxStack: 1},
	SPLASH_POTION:        {ID: SPLASH_POTION, Name: "Splash Potion", MaxStack: 1},
	BOTTLE_O_ENCHANTING:  {ID: BOTTLE_O_ENCHANTING, Name: "Bottle o' Enchanting", MaxStack: 64},
	BLAZE_ROD:            {ID: BLAZE_ROD, Name: "Blaze Rod", MaxStack: 64, FuelTime: 2400},
	BLAZE_POWDER:         {ID: BLAZE_POWDER, Name: "Blaze Powder", MaxStack: 64},
	GHAST_TEAR:           {ID: GHAST_TEAR, Name: "Ghast Tear", MaxStack: 64},
	MAGMA_CREAM:          {ID: MAGMA_CREAM, Name: "Magma Cream", MaxStack: 64},
	FERMENTED_SPIDER_EYE: {ID: FERMENTED_SPIDER_EYE, Name: "Fermented Spider Eye", MaxStack: 64},
	GLOWSTONE_DUST:       {ID: GLOWSTONE_DUST, Name: "Glowstone Dust", MaxStack: 64},
	GLISTERING_MELON:     {ID: GLISTERING_MELON, Name: "Glistering Melon", MaxStack: 64},
	NETHER_WART:          {ID: NETHER_WART, Name: "Nether Wart", MaxStack: 64, IsPlaceable: true},
	SEEDS:          {ID: SEEDS, Name: "Wheat Seeds", MaxStack: 64, IsPlaceable: true},
	PUMPKIN_SEEDS:  {ID: PUMPKIN_SEEDS, Name: "Pumpkin Seeds", MaxStack: 64, IsPlaceable: true},
	MELON_SEEDS:    {ID: MELON_SEEDS, Name: "Melon Seeds", MaxStack: 64, IsPlaceable: true},
	BEETROOT_SEEDS: {ID: BEETROOT_SEEDS, Name: "Beetroot Seeds", MaxStack: 64, IsPlaceable: true},
	SPAWN_EGG: {ID: SPAWN_EGG, Name: "Spawn Egg", MaxStack: 64},
	EGG:       {ID: EGG, Name: "Egg", MaxStack: 16},
	SNOWBALL:  {ID: SNOWBALL, Name: "Snowball", MaxStack: 16},
	SUGARCANE: {ID: SUGARCANE, Name: "Sugar Cane", MaxStack: 64, IsPlaceable: true},
	DYE:       {ID: DYE, Name: "Dye", MaxStack: 64},
	MOB_HEAD:  {ID: MOB_HEAD, Name: "Mob Head", MaxStack: 64, IsPlaceable: true},
	RABBIT_FOOT: {ID: RABBIT_FOOT, Name: "Rabbit's Foot", MaxStack: 64},
	RABBIT_HIDE: {ID: RABBIT_HIDE, Name: "Rabbit Hide", MaxStack: 64},
	BREWING_STAND: {ID: BREWING_STAND, Name: "Brewing Stand", MaxStack: 64, IsPlaceable: true},
	CAULDRON:      {ID: CAULDRON, Name: "Cauldron", MaxStack: 64, IsPlaceable: true},
	HOPPER:        {ID: HOPPER, Name: "Hopper", MaxStack: 64, IsPlaceable: true},
}
func GetMiscItemInfo(id int) *MiscItemInfo {
	info, ok := miscItems[id]
	if !ok {
		return nil
	}
	return &info
}
func GetFuelTime(id int) int {
	if info, ok := miscItems[id]; ok && info.FuelTime > 0 {
		return info.FuelTime
	}
	return 0
}
func IsMiscPlaceable(id int) bool {
	if info, ok := miscItems[id]; ok {
		return info.IsPlaceable
	}
	return false
}

const (
	BucketEmpty = 0
	BucketMilk  = 1
	BucketWater = 8
	BucketLava  = 10
)
func BucketTypeToBlock(meta int) int {
	switch meta {
	case BucketWater:
		return 8
	case BucketLava:
		return 10
	default:
		return 0
	}
}
