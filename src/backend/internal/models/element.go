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
	ValidRecipeIdx []int         `json:"validRecipeIdx"`
	Recipes        []*RecipeNode `json:"recipes"`
}

func (node *ElementNode) checkValidRecipe(recipeIdx int) {
	// fmt.Printf("checking valid recipe of%s, ", node.Name)
	recipe := node.Recipes[recipeIdx]
	if recipe.Ingredient1 == nil || recipe.Ingredient2 == nil {
		// fmt.Printf("ingreidents nil%s\n", node.Name)
		return
	}
	if recipe.Ingredient1.IsValid && recipe.Ingredient2.IsValid {
		for len(node.ValidRecipeIdx) <= recipeIdx {
			node.ValidRecipeIdx = append(node.ValidRecipeIdx, 0)
		}
		// fmt.Printf("ingreidents valid%s\n", node.Name)
		node.ValidRecipeIdx[recipeIdx] = sumSlice(node.Recipes[recipeIdx].Ingredient1.ValidRecipeIdx) * sumSlice(node.Recipes[recipeIdx].Ingredient2.ValidRecipeIdx)
		// node.ValidRecipeIdx[recipeIdx]++
		node.setValid()
	} else {
		// fmt.Printf("whats valid?%s %s(%t) %s(%t)\n", node.Name, recipe.Ingredient1.Name, recipe.Ingredient1.IsValid, recipe.Ingredient2.Name, recipe.Ingredient2.IsValid)
	}
}

func (node *ElementNode) setValid() {
	node.IsValid = true
	// fmt.Printf("validating %s\n", node.Name)
	if node.Parent != nil {
		node.Parent.checkValidRecipe(node.Index)
	}
}
