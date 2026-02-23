package item

type FoodInfo struct {
	ID             int
	Name           string
	FoodRestore    int
	Saturation     float32
	ResidueID      int
	ResidueMeta    int
	EffectChance   float32
	EffectID       int
	EffectDuration int
	EffectLevel    int
	AlwaysEdible   bool
}

var foods = map[int]FoodInfo{

	APPLE:        {ID: APPLE, Name: "Apple", FoodRestore: 4, Saturation: 2.4},
	BREAD:        {ID: BREAD, Name: "Bread", FoodRestore: 5, Saturation: 6.0},
	COOKIE:       {ID: COOKIE, Name: "Cookie", FoodRestore: 2, Saturation: 0.4},
	MELON:        {ID: MELON, Name: "Melon Slice", FoodRestore: 2, Saturation: 1.2},
	CARROT:       {ID: CARROT, Name: "Carrot", FoodRestore: 3, Saturation: 3.6},
	POTATO:       {ID: POTATO, Name: "Potato", FoodRestore: 1, Saturation: 0.6},
	BAKED_POTATO: {ID: BAKED_POTATO, Name: "Baked Potato", FoodRestore: 5, Saturation: 6.0},
	BEETROOT:     {ID: BEETROOT, Name: "Beetroot", FoodRestore: 1, Saturation: 1.2},
	PUMPKIN_PIE:  {ID: PUMPKIN_PIE, Name: "Pumpkin Pie", FoodRestore: 8, Saturation: 4.8},

	RAW_PORKCHOP: {ID: RAW_PORKCHOP, Name: "Raw Porkchop", FoodRestore: 3, Saturation: 1.8},
	RAW_BEEF:     {ID: RAW_BEEF, Name: "Raw Beef", FoodRestore: 3, Saturation: 1.8},
	RAW_CHICKEN: {ID: RAW_CHICKEN, Name: "Raw Chicken", FoodRestore: 2, Saturation: 1.2,
		EffectChance: 0.3, EffectID: 17, EffectDuration: 600, EffectLevel: 0},
	RAW_RABBIT: {ID: RAW_RABBIT, Name: "Raw Rabbit", FoodRestore: 3, Saturation: 1.8},
	RAW_FISH:   {ID: RAW_FISH, Name: "Raw Fish", FoodRestore: 2, Saturation: 0.4},

	COOKED_PORKCHOP: {ID: COOKED_PORKCHOP, Name: "Cooked Porkchop", FoodRestore: 8, Saturation: 12.8},
	STEAK:           {ID: STEAK, Name: "Steak", FoodRestore: 8, Saturation: 12.8},
	COOKED_CHICKEN:  {ID: COOKED_CHICKEN, Name: "Cooked Chicken", FoodRestore: 6, Saturation: 7.2},
	COOKED_RABBIT:   {ID: COOKED_RABBIT, Name: "Cooked Rabbit", FoodRestore: 5, Saturation: 6.0},
	COOKED_FISH:     {ID: COOKED_FISH, Name: "Cooked Fish", FoodRestore: 5, Saturation: 6.0},

	GOLDEN_APPLE:  {ID: GOLDEN_APPLE, Name: "Golden Apple", FoodRestore: 4, Saturation: 9.6, AlwaysEdible: true},
	GOLDEN_CARROT: {ID: GOLDEN_CARROT, Name: "Golden Carrot", FoodRestore: 6, Saturation: 14.4},
	ROTTEN_FLESH: {ID: ROTTEN_FLESH, Name: "Rotten Flesh", FoodRestore: 4, Saturation: 0.8,
		EffectChance: 0.8, EffectID: 17, EffectDuration: 600, EffectLevel: 0},
	POISONOUS_POTATO: {ID: POISONOUS_POTATO, Name: "Poisonous Potato", FoodRestore: 2, Saturation: 1.2,
		EffectChance: 0.6, EffectID: 19, EffectDuration: 100, EffectLevel: 0},

	MUSHROOM_STEW: {ID: MUSHROOM_STEW, Name: "Mushroom Stew", FoodRestore: 6, Saturation: 7.2,
		ResidueID: BOWL},
	RABBIT_STEW: {ID: RABBIT_STEW, Name: "Rabbit Stew", FoodRestore: 10, Saturation: 12.0,
		ResidueID: BOWL},
	BEETROOT_SOUP: {ID: BEETROOT_SOUP, Name: "Beetroot Soup", FoodRestore: 6, Saturation: 7.2,
		ResidueID: BOWL},

	RAW_SALMON:    {ID: RAW_SALMON, Name: "Raw Salmon", FoodRestore: 2, Saturation: 0.4},
	COOKED_SALMON: {ID: COOKED_SALMON, Name: "Cooked Salmon", FoodRestore: 6, Saturation: 9.6},
	CLOWN_FISH:    {ID: CLOWN_FISH, Name: "Clownfish", FoodRestore: 1, Saturation: 0.2},
	PUFFER_FISH: {ID: PUFFER_FISH, Name: "Pufferfish", FoodRestore: 1, Saturation: 0.2,
		EffectChance: 1.0, EffectID: 19, EffectDuration: 1200, EffectLevel: 3},

	SPIDER_EYE: {ID: SPIDER_EYE, Name: "Spider Eye", FoodRestore: 2, Saturation: 3.2,
		EffectChance: 1.0, EffectID: 19, EffectDuration: 100, EffectLevel: 0},
	ENCHANTED_GOLDEN_APPLE: {ID: ENCHANTED_GOLDEN_APPLE, Name: "Enchanted Golden Apple",
		FoodRestore: 4, Saturation: 9.6, AlwaysEdible: true},
}

func GetFoodInfo(id int) *FoodInfo {
	info, ok := foods[id]
	if !ok {
		return nil
	}
	return &info
}

func IsFoodItem(id int) bool {
	_, ok := foods[id]
	return ok
}

type EatResult struct {
	FoodRestore    int
	Saturation     float32
	ResidueItem    *Item
	HasEffect      bool
	EffectID       int
	EffectDuration int
	EffectLevel    int
}

func CanEat(itemID int, currentFood int, maxFood int) bool {
	info := GetFoodInfo(itemID)
	if info == nil {
		return false
	}
	if info.AlwaysEdible {
		return true
	}
	return currentFood < maxFood
}

func Eat(item Item) *EatResult {
	info := GetFoodInfo(item.ID)
	if info == nil {
		return nil
	}

	result := &EatResult{
		FoodRestore: info.FoodRestore,
		Saturation:  info.Saturation,
	}

	if info.ResidueID > 0 {
		residue := NewItem(info.ResidueID, info.ResidueMeta, 1)
		result.ResidueItem = &residue
	}

	if info.EffectChance > 0 && info.EffectID > 0 {

		result.HasEffect = true
		result.EffectID = info.EffectID
		result.EffectDuration = info.EffectDuration
		result.EffectLevel = info.EffectLevel
	}

	return result
}

const (
	MaxFood       = 20
	MaxSaturation = 20.0
	FoodTickTimer = 80

	ExhaustionActionAttack = 0.1
	ExhaustionActionDamage = 0.1
	ExhaustionActionMine   = 0.005
	ExhaustionActionSprint = 0.1
	ExhaustionActionJump   = 0.05
	ExhaustionActionSwim   = 0.01
	ExhaustionActionWalk   = 0.0
	ExhaustionActionRegen  = 6.0
)

type HungerResult struct {
	ShouldHeal   bool
	ShouldDamage bool
}

func CalcHungerTick(currentFood int, saturation float32) HungerResult {
	return HungerResult{
		ShouldHeal:   currentFood >= 18 && saturation > 0,
		ShouldDamage: currentFood <= 0,
	}
}

func AllFoodIDs() []int {
	ids := make([]int, 0, len(foods))
	for id := range foods {
		ids = append(ids, id)
	}
	return ids
}
