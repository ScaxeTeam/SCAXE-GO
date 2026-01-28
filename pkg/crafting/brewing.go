package crafting

type BrewingRecipe struct {
	result     *Item
	ingredient *Item
	potion     *Item
}

func NewBrewingRecipe(result, ingredient, potion *Item) *BrewingRecipe {
	return &BrewingRecipe{
		result:     result,
		ingredient: ingredient,
		potion:     potion,
	}
}

func (r *BrewingRecipe) GetResult() *Item {
	return r.result
}

func (r *BrewingRecipe) GetIngredient() *Item {
	return r.ingredient
}

func (r *BrewingRecipe) GetPotion() *Item {
	return r.potion
}

func (r *BrewingRecipe) GetIngredients() []*Item {
	return []*Item{r.ingredient, r.potion}
}

func (r *BrewingRecipe) Matches(items []*Item) bool {
	if len(items) != 2 {
		return false
	}
	hasIngredient := false
	hasPotion := false
	for _, item := range items {
		if r.ingredient.Matches(item, true) {
			hasIngredient = true
		}
		if r.potion.Matches(item, true) {
			hasPotion = true
		}
	}
	return hasIngredient && hasPotion
}

const (
	PotionWater              = 0
	PotionMundane            = 1
	PotionMundaneExtended    = 2
	PotionThick              = 3
	PotionAwkward            = 4
	PotionNightVision        = 5
	PotionNightVisionLong    = 6
	PotionInvisibility       = 7
	PotionInvisibilityLong   = 8
	PotionLeaping            = 9
	PotionLeapingLong        = 10
	PotionLeapingStrong      = 11
	PotionFireResistance     = 12
	PotionFireResistanceLong = 13
	PotionSwiftness          = 14
	PotionSwiftnessLong      = 15
	PotionSwiftnessStrong    = 16
	PotionSlowness           = 17
	PotionSlownessLong       = 18
	PotionWaterBreathing     = 19
	PotionWaterBreathingLong = 20
	PotionHealing            = 21
	PotionHealingStrong      = 22
	PotionHarming            = 23
	PotionHarmingStrong      = 24
	PotionPoison             = 25
	PotionPoisonLong         = 26
	PotionPoisonStrong       = 27
	PotionRegeneration       = 28
	PotionRegenerationLong   = 29
	PotionRegenerationStrong = 30
	PotionStrength           = 31
	PotionStrengthLong       = 32
	PotionStrengthStrong     = 33
	PotionWeakness           = 34
	PotionWeaknessLong       = 35
)

const (
	ItemPotion        = 373
	ItemSplashPotion  = 438
	ItemNetherWart    = 372
	ItemRedstoneDust  = 331
	ItemGlowstoneDust = 348
	ItemFermentedEye  = 376
	ItemGunpowder     = 289
	ItemMagmaCream    = 378
	ItemBlazeRod      = 369
	ItemGhastTear     = 370
	ItemSpiderEye     = 375
	ItemGoldenCarrot  = 396
	ItemSugar         = 353
	ItemRabbitFoot    = 414
	ItemPufferfish    = 462
)

func (m *CraftingManager) RegisterDefaultBrewingRecipes() {

	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionAwkward, 1),
		NewItem(ItemNetherWart, 0, 1),
		NewItem(ItemPotion, PotionWater, 1),
	))

	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionNightVision, 1),
		NewItem(ItemGoldenCarrot, 0, 1),
		NewItem(ItemPotion, PotionAwkward, 1),
	))
	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionNightVisionLong, 1),
		NewItem(ItemRedstoneDust, 0, 1),
		NewItem(ItemPotion, PotionNightVision, 1),
	))

	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionInvisibility, 1),
		NewItem(ItemFermentedEye, 0, 1),
		NewItem(ItemPotion, PotionNightVision, 1),
	))
	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionInvisibilityLong, 1),
		NewItem(ItemRedstoneDust, 0, 1),
		NewItem(ItemPotion, PotionInvisibility, 1),
	))

	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionFireResistance, 1),
		NewItem(ItemMagmaCream, 0, 1),
		NewItem(ItemPotion, PotionAwkward, 1),
	))
	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionFireResistanceLong, 1),
		NewItem(ItemRedstoneDust, 0, 1),
		NewItem(ItemPotion, PotionFireResistance, 1),
	))

	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionSwiftness, 1),
		NewItem(ItemSugar, 0, 1),
		NewItem(ItemPotion, PotionAwkward, 1),
	))
	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionSwiftnessLong, 1),
		NewItem(ItemRedstoneDust, 0, 1),
		NewItem(ItemPotion, PotionSwiftness, 1),
	))
	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionSwiftnessStrong, 1),
		NewItem(ItemGlowstoneDust, 0, 1),
		NewItem(ItemPotion, PotionSwiftness, 1),
	))

	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionHealing, 1),
		NewItem(ItemGhastTear, 0, 1),
		NewItem(ItemPotion, PotionAwkward, 1),
	))
	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionHealingStrong, 1),
		NewItem(ItemGlowstoneDust, 0, 1),
		NewItem(ItemPotion, PotionHealing, 1),
	))

	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionPoison, 1),
		NewItem(ItemSpiderEye, 0, 1),
		NewItem(ItemPotion, PotionAwkward, 1),
	))
	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionPoisonLong, 1),
		NewItem(ItemRedstoneDust, 0, 1),
		NewItem(ItemPotion, PotionPoison, 1),
	))

	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionRegeneration, 1),
		NewItem(ItemGhastTear, 0, 1),
		NewItem(ItemPotion, PotionAwkward, 1),
	))

	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionStrength, 1),
		NewItem(ItemBlazeRod, 0, 1),
		NewItem(ItemPotion, PotionAwkward, 1),
	))

	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionWeakness, 1),
		NewItem(ItemFermentedEye, 0, 1),
		NewItem(ItemPotion, PotionWater, 1),
	))

	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionWaterBreathing, 1),
		NewItem(ItemPufferfish, 0, 1),
		NewItem(ItemPotion, PotionAwkward, 1),
	))

	m.RegisterBrewingRecipe(NewBrewingRecipe(
		NewItem(ItemPotion, PotionLeaping, 1),
		NewItem(ItemRabbitFoot, 0, 1),
		NewItem(ItemPotion, PotionAwkward, 1),
	))
}

func (m *CraftingManager) GetBrewingRecipes() []*BrewingRecipe {
	return m.brewingRecipes
}

func (m *CraftingManager) RegisterBrewingRecipe(recipe *BrewingRecipe) {
	if m.brewingRecipes == nil {
		m.brewingRecipes = make([]*BrewingRecipe, 0)
	}
	m.brewingRecipes = append(m.brewingRecipes, recipe)
}

func (m *CraftingManager) FindBrewingRecipe(ingredient, potion *Item) *BrewingRecipe {
	for _, recipe := range m.brewingRecipes {
		if recipe.ingredient.Matches(ingredient, true) && recipe.potion.Matches(potion, true) {
			return recipe
		}
	}
	return nil
}
