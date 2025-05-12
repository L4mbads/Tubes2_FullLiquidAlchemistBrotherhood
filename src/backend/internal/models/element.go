package models

type ElementType struct {
	Name     string
	ImageUrl string
	Type     string
}

type ElementNode struct {
	Index          int           `json:"index"`
	Name           string        `json:"name"`
	Parent         *ElementNode  `json:"-"`
	IsValid        bool          `json:"-"`
	Depth          int           `json:"-"`
	ValidRecipeIdx []int         `json:"validRecipeIdx"`
	Recipes        []*RecipeNode `json:"recipes"`
}

func (node *ElementNode) checkValidRecipe(recipeIdx int) {
	recipe := node.Recipes[recipeIdx]
	if recipe.Ingredient1 == nil || recipe.Ingredient2 == nil {
		return
	}
	if recipe.Ingredient1.IsValid && recipe.Ingredient2.IsValid {
		for len(node.ValidRecipeIdx) <= recipeIdx {
			node.ValidRecipeIdx = append(node.ValidRecipeIdx, 0)
		}
		node.ValidRecipeIdx[recipeIdx] = sumSlice(node.Recipes[recipeIdx].Ingredient1.ValidRecipeIdx) * sumSlice(node.Recipes[recipeIdx].Ingredient2.ValidRecipeIdx)
		// node.ValidRecipeIdx[recipeIdx]++
		node.setValid()
	}
}

func (node *ElementNode) setValid() {
	node.IsValid = true
	if node.Parent != nil {
		node.Parent.checkValidRecipe(node.Index)
	}
}
