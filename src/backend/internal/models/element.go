package models

import "fmt"

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
	fmt.Printf("check %s %d, %d\n", node.Name, recipeIdx, len(node.ValidRecipeIdx))
	recipe := node.Recipes[recipeIdx]
	fmt.Println("1")
	if recipe.Ingredient1 == nil || recipe.Ingredient2 == nil {
		return
	}
	if recipe.Ingredient1.IsValid && recipe.Ingredient2.IsValid {
		fmt.Println("2")
		for len(node.ValidRecipeIdx) <= recipeIdx {
			fmt.Println("3")
			node.ValidRecipeIdx = append(node.ValidRecipeIdx, 0)
		}
		fmt.Println("4")
		node.ValidRecipeIdx[recipeIdx] = sumSlice(node.Recipes[recipeIdx].Ingredient1.ValidRecipeIdx) * sumSlice(node.Recipes[recipeIdx].Ingredient2.ValidRecipeIdx)
		// node.ValidRecipeIdx[recipeIdx]++
		fmt.Printf("[%d]%d anjay\n", recipeIdx, node.ValidRecipeIdx[recipeIdx])
		fmt.Println("5")
		node.setValid()
	}
}

func (node *ElementNode) setValid() {
	node.IsValid = true
	if node.Parent != nil {
		node.Parent.checkValidRecipe(node.Index)
	}
}
