package crafting

type Item struct {
	ID    int
	Meta  int
	Count int
}

func NewItem(id, meta, count int) *Item {
	return &Item{ID: id, Meta: meta, Count: count}
}

func (i *Item) Matches(other *Item, checkMeta bool) bool {
	if i.ID != other.ID {
		return false
	}
	if checkMeta && i.Meta != other.Meta && i.Meta != -1 && other.Meta != -1 {
		return false
	}
	return true
}

type Recipe interface {
	GetResult() *Item
	GetIngredients() []*Item
	Matches(items []*Item) bool
}

type ShapelessRecipe struct {
	result      *Item
	ingredients []*Item
}

func NewShapelessRecipe(result *Item, ingredients ...*Item) *ShapelessRecipe {
	return &ShapelessRecipe{
		result:      result,
		ingredients: ingredients,
	}
}

func (r *ShapelessRecipe) GetResult() *Item {
	return r.result
}

func (r *ShapelessRecipe) GetIngredients() []*Item {
	return r.ingredients
}

func (r *ShapelessRecipe) Matches(items []*Item) bool {
	need := make([]*Item, len(r.ingredients))
	copy(need, r.ingredients)

	for _, item := range items {
		if item == nil || item.ID == 0 {
			continue
		}
		found := false
		for i, ing := range need {
			if ing != nil && ing.Matches(item, true) {
				need[i] = nil
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	for _, ing := range need {
		if ing != nil {
			return false
		}
	}
	return true
}

type ShapedRecipe struct {
	result *Item
	shape  []string
	keys   map[rune]*Item
	width  int
	height int
}

func NewShapedRecipe(result *Item, shape ...string) *ShapedRecipe {
	height := len(shape)
	width := 0
	for _, row := range shape {
		if len(row) > width {
			width = len(row)
		}
	}
	return &ShapedRecipe{
		result: result,
		shape:  shape,
		keys:   make(map[rune]*Item),
		width:  width,
		height: height,
	}
}

func (r *ShapedRecipe) SetIngredient(char rune, item *Item) *ShapedRecipe {
	r.keys[char] = item
	return r
}

func (r *ShapedRecipe) GetResult() *Item {
	return r.result
}

func (r *ShapedRecipe) GetIngredients() []*Item {
	ingredients := make([]*Item, 0)
	for _, item := range r.keys {
		if item != nil && item.ID != 0 {
			ingredients = append(ingredients, item)
		}
	}
	return ingredients
}

func (r *ShapedRecipe) GetWidth() int {
	return r.width
}

func (r *ShapedRecipe) GetHeight() int {
	return r.height
}

func (r *ShapedRecipe) Matches(items []*Item) bool {

	return len(items) >= len(r.GetIngredients())
}

type FurnaceRecipe struct {
	result *Item
	input  *Item
}

func NewFurnaceRecipe(result, input *Item) *FurnaceRecipe {
	return &FurnaceRecipe{
		result: result,
		input:  input,
	}
}

func (r *FurnaceRecipe) GetResult() *Item {
	return r.result
}

func (r *FurnaceRecipe) GetInput() *Item {
	return r.input
}

func (r *FurnaceRecipe) GetIngredients() []*Item {
	return []*Item{r.input}
}

func (r *FurnaceRecipe) Matches(items []*Item) bool {
	if len(items) != 1 {
		return false
	}
	return r.input.Matches(items[0], true)
}

type CraftingManager struct {
	recipes        []Recipe
	furnaceRecipes []*FurnaceRecipe
	brewingRecipes []*BrewingRecipe
}

func NewCraftingManager() *CraftingManager {
	cm := &CraftingManager{
		recipes:        make([]Recipe, 0),
		furnaceRecipes: make([]*FurnaceRecipe, 0),
		brewingRecipes: make([]*BrewingRecipe, 0),
	}
	cm.registerDefaultRecipes()
	cm.RegisterDefaultBrewingRecipes()
	return cm
}

var globalCraftingManager *CraftingManager

func GetCraftingManager() *CraftingManager {
	if globalCraftingManager == nil {
		globalCraftingManager = NewCraftingManager()
	}
	return globalCraftingManager
}

func (m *CraftingManager) RegisterRecipe(recipe Recipe) {
	m.recipes = append(m.recipes, recipe)
}

func (m *CraftingManager) RegisterFurnaceRecipe(recipe *FurnaceRecipe) {
	m.furnaceRecipes = append(m.furnaceRecipes, recipe)
}

func (m *CraftingManager) GetRecipes() []Recipe {
	return m.recipes
}

func (m *CraftingManager) FindRecipe(items []*Item) Recipe {
	for _, recipe := range m.recipes {
		if recipe.Matches(items) {
			return recipe
		}
	}
	return nil
}

func (m *CraftingManager) FindFurnaceRecipe(input *Item) *FurnaceRecipe {
	for _, recipe := range m.furnaceRecipes {
		if recipe.input.Matches(input, true) {
			return recipe
		}
	}
	return nil
}

func (m *CraftingManager) registerDefaultRecipes() {

	m.RegisterRecipe(NewShapedRecipe(NewItem(58, 0, 1), "XX", "XX").
		SetIngredient('X', NewItem(5, -1, 1)))

	m.RegisterRecipe(NewShapedRecipe(NewItem(280, 0, 4), "X", "X").
		SetIngredient('X', NewItem(5, -1, 1)))

	m.RegisterRecipe(NewShapedRecipe(NewItem(61, 0, 1), "XXX", "X X", "XXX").
		SetIngredient('X', NewItem(4, 0, 1)))

	m.RegisterRecipe(NewShapedRecipe(NewItem(54, 0, 1), "XXX", "X X", "XXX").
		SetIngredient('X', NewItem(5, -1, 1)))

	m.RegisterRecipe(NewShapedRecipe(NewItem(50, 0, 4), "C", "S").
		SetIngredient('C', NewItem(263, 0, 1)).
		SetIngredient('S', NewItem(280, 0, 1)))

	m.RegisterRecipe(NewShapelessRecipe(NewItem(5, 0, 4), NewItem(17, 0, 1)))
	m.RegisterRecipe(NewShapelessRecipe(NewItem(5, 1, 4), NewItem(17, 1, 1)))
	m.RegisterRecipe(NewShapelessRecipe(NewItem(5, 2, 4), NewItem(17, 2, 1)))
	m.RegisterRecipe(NewShapelessRecipe(NewItem(5, 3, 4), NewItem(17, 3, 1)))

	m.RegisterFurnaceRecipe(NewFurnaceRecipe(NewItem(266, 0, 1), NewItem(14, 0, 1)))
	m.RegisterFurnaceRecipe(NewFurnaceRecipe(NewItem(265, 0, 1), NewItem(15, 0, 1)))
	m.RegisterFurnaceRecipe(NewFurnaceRecipe(NewItem(263, 0, 1), NewItem(17, -1, 1)))
	m.RegisterFurnaceRecipe(NewFurnaceRecipe(NewItem(1, 0, 1), NewItem(4, 0, 1)))
	m.RegisterFurnaceRecipe(NewFurnaceRecipe(NewItem(20, 0, 1), NewItem(12, 0, 1)))
	m.RegisterFurnaceRecipe(NewFurnaceRecipe(NewItem(336, 0, 1), NewItem(337, 0, 1)))
	m.RegisterFurnaceRecipe(NewFurnaceRecipe(NewItem(351, 4, 1), NewItem(21, 0, 1)))
	m.RegisterFurnaceRecipe(NewFurnaceRecipe(NewItem(264, 0, 1), NewItem(56, 0, 1)))
	m.RegisterFurnaceRecipe(NewFurnaceRecipe(NewItem(331, 0, 1), NewItem(73, 0, 1)))
	m.RegisterFurnaceRecipe(NewFurnaceRecipe(NewItem(388, 0, 1), NewItem(129, 0, 1)))
	m.RegisterFurnaceRecipe(NewFurnaceRecipe(NewItem(405, 0, 1), NewItem(87, 0, 1)))
}
