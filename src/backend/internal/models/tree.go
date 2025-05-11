package models

import (
	"database/sql"
	"fmt"
	"slices"
)

type Recipe struct {
	Ingredient1 string
	Ingredient2 string
}

func ShortestDFS(db *sql.DB, parentNode *ElementNode, element string, depth int, limit int, targetCount int) (*ElementNode, int, error) {
	node := &ElementNode{Name: element, Parent: parentNode, IsValid: false}

	fmt.Printf("babik now %s\n", element)
	if isBasicElement(element) {
		node.ValidRecipeIdx = append(node.ValidRecipeIdx, 1)
		fmt.Println("single done")
		node.setValid()
		return node, 0, nil
	}

	typeQuery := "SELECT type FROM elements WHERE name = $1"
	row := db.QueryRow(typeQuery, element)
	var elementType int
	err := row.Scan(&elementType)
	if err == sql.ErrNoRows {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, nil
	}

	query := "SELECT ingredient1, ingredient2 FROM recipes WHERE element = $1"
	rows, err := db.Query(query, element)
	if err != nil {
		return nil, 0, err
	}

	var recipes []Recipe
	// create a list/array to store ingredients
	for rows.Next() {
		var ingredient1, ingredient2 string
		if err := rows.Scan(&ingredient1, &ingredient2); err != nil {
			return nil, 0, err
		}
		// if !isElementNotLooping(parentNode, ingredient1) || !isElementNotLooping(parentNode, ingredient2) {
		// 	continue
		// }

		// do not continue path if recipes are higher type
		query := "SELECT type FROM elements WHERE name = $1"
		row := db.QueryRow(query, ingredient1)
		var elementType1 int
		err := row.Scan(&elementType1)
		if err == sql.ErrNoRows {
			continue
		} else if err != nil {
			return nil, 0, nil
		}

		if elementType1 >= elementType {
			continue
		}

		query = "SELECT type FROM elements WHERE name = $1"
		row = db.QueryRow(query, ingredient2)
		var elementType2 int
		err = row.Scan(&elementType2)
		if err == sql.ErrNoRows {
			continue
		} else if err != nil {
			return nil, 0, nil
		}

		if elementType2 >= elementType {
			continue
		}

		recipe := Recipe{
			Ingredient1: ingredient1,
			Ingredient2: ingredient2,
		}

		fmt.Printf("%s %s\n", ingredient1, ingredient2)

		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	rows.Close()

	if len(recipes) == 0 {
		fmt.Printf("error kosong\n")
		return nil, 0, err
	}

	minDepth := 999999
	i := 0
	validPath := 0
	root := node
	for root.Parent != nil {
		root = root.Parent
	}
	// fmt.Printf("ROOTNYA BENER GAK %s %d %d\n", root.Name, sumSlice(root.ValidRecipeIdx), targetCount)
	for _, recipe := range recipes {
		if validPath >= targetCount {
			break
		}
		x := sumSlice(root.ValidRecipeIdx)
		fmt.Printf("AAAAAAAA %d >= %d\n", x, targetCount)
		if x >= targetCount {
			fmt.Printf("%d / %d VALID RECIPE IDX\n", targetCount, x)
			break
		}

		ingredient1 := recipe.Ingredient1
		ingredient2 := recipe.Ingredient2

		hasValidRecipe := true
		// maxDepth := -1
		currentRecipeNode := &RecipeNode{Ingredient1: nil, Ingredient2: nil}
		node.Recipes = append(node.Recipes, currentRecipeNode)
		if ingredient1 != "" {
			fmt.Printf("entering %s (%d, %d)\n", ingredient1, depth, minDepth)
			child1, _, err := ShortestDFS(db, node, ingredient1, depth+1, limit, targetCount)
			if err != nil {
				return nil, 0, err
			}
			if child1 != nil {
				child1.Index = i
				currentRecipeNode.Ingredient1 = child1
				child1.setValid()
				// maxDepth = depthIngredient1
				// fmt.Printf("%s: maxDepth1(%s): %d\n", element, ingredient1, depthIngredient1)
			} else {
				hasValidRecipe = false
			}
		}
		x = sumSlice(root.ValidRecipeIdx)
		fmt.Printf("AAAAAAAA %d >= %d\n", x, targetCount)
		if x >= targetCount {
			fmt.Printf("%d / %d VALID RECIPE IDX\n", targetCount, x)
			break
		}
		if hasValidRecipe && ingredient2 != "" {
			fmt.Printf("entering %s (%d, %d)\n", ingredient2, depth, minDepth)
			child2, _, err := ShortestDFS(db, node, ingredient2, depth+1, limit, targetCount)
			if err != nil {
				return nil, 0, err
			}
			if child2 != nil {
				child2.Index = i
				currentRecipeNode.Ingredient2 = child2
				child2.setValid()
				// fmt.Printf("%s: maxDepth vs depth2: %d v %d\n", element, maxDepth, depthIngredient2)
				// if depthIngredient2 > maxDepth {
				// 	fmt.Printf("%s: maxDepth2(%s): %d\n", element, ingredient2, depthIngredient2)
				// 	maxDepth = depthIngredient2
				// }
			} else {
				hasValidRecipe = false
			}
		}

		if hasValidRecipe {
			// node.setValid()
			// 	// recipeNode = currentRecipeNode
			// 	minDepth = maxDepth
			// 	limit = minDepth + depth
			// 	fmt.Printf("%s %s valid (%d)\n", ingredient1, ingredient2, depth)
			validPath++
		}
		i++
	}
	// if minDepth != -1 {
	// 	node.Recipes = append(node.Recipes, recipeNode)
	// 	fmt.Println("ada yg elesai")
	// }
	cutTree(node)
	return node, minDepth + 1, nil
}

func ShortestBFS(db *sql.DB, parentNode *ElementNode, element string, depth int, limit int, targetCount int) (*ElementNode, int, error) {
	// initialize root node
	root := &ElementNode{Name: element, Parent: parentNode, IsValid: false, Depth: depth}

	// use queue for BFS
	queue := RecipeQueue{}
	queue.enqueue(&RecipeNode{Ingredient1: root, Ingredient2: nil})

	var shortestRecipeNode *RecipeNode
	var minDepth = limit

	for !queue.isEmpty() {
		recipeCount := sumSlice(root.ValidRecipeIdx)
		fmt.Printf("RECIPE COUNT: %d\n", recipeCount)
		if recipeCount >= targetCount {
			fmt.Printf("%d VALID RECIPE IDX\n", targetCount)
			break
		}
		currentRecipe := queue.dequeue()

		currentNode1 := currentRecipe.Ingredient1
		currentNode2 := currentRecipe.Ingredient2

		// both basic element
		if currentNode1 != nil && currentNode2 != nil &&
			isBasicElement(currentNode1.Name) && isBasicElement(currentNode2.Name) {
			currentNode1.ValidRecipeIdx = append(currentNode1.ValidRecipeIdx, 1)
			currentNode2.ValidRecipeIdx = append(currentNode2.ValidRecipeIdx, 1)
			currentNode1.setValid()
			currentNode2.setValid()
			shortestRecipeNode = currentRecipe
			minDepth = currentNode1.Depth
			continue
		}
		// recipeCount = sumSlice(root.ValidRecipeIdx)
		// fmt.Printf("RECIPE COUNT: %d\n", recipeCount)
		// if recipeCount >= targetCount {
		// 	fmt.Printf("%d VALID RECIPE IDX\n", targetCount)
		// 	break
		// }

		if currentNode1 != nil {
			if isBasicElement(currentNode1.Name) {
				// ingredient 1 basic element
				currentNode1.ValidRecipeIdx = append(currentNode1.ValidRecipeIdx, 1)
				currentNode1.setValid()
			} else {
				typeQuery := "SELECT type FROM elements WHERE name = $1"
				row := db.QueryRow(typeQuery, currentNode1.Name)
				var elementType int
				err := row.Scan(&elementType)
				if err == sql.ErrNoRows {
					continue
				} else if err != nil {
					return nil, 0, nil
				}

				query := "SELECT ingredient1, ingredient2 FROM recipes WHERE element = $1"
				rows, err := db.Query(query, currentNode1.Name)
				if err != nil {
					return nil, 0, err
				}

				i := 0
				for rows.Next() {
					// if i >= targetCount {
					// 	break
					// }
					var ingredient1, ingredient2 string
					if err := rows.Scan(&ingredient1, &ingredient2); err != nil {
						rows.Close()
						return nil, 0, err
					}
					// do not continue path if recipes are higher type
					query := "SELECT type FROM elements WHERE name = $1"
					row := db.QueryRow(query, ingredient1)
					var elementType1 int
					err := row.Scan(&elementType1)
					if err == sql.ErrNoRows {
						continue
					} else if err != nil {
						return nil, 0, nil
					}
					if elementType1 >= elementType {
						continue
					}

					query = "SELECT type FROM elements WHERE name = $1"
					row = db.QueryRow(query, ingredient2)
					var elementType2 int
					err = row.Scan(&elementType2)
					if err == sql.ErrNoRows {
						return nil, 0, nil
					} else if err != nil {
						return nil, 0, nil
					}

					if elementType2 >= elementType {
						continue
					}

					// if !isElementNotLooping(currentNode1.Parent, ingredient1) || !isElementNotLooping(currentNode1.Parent, ingredient2) {
					// 	continue
					// }

					child1 := &ElementNode{Name: ingredient1, Parent: currentNode1, Depth: currentNode1.Depth + 1, Index: i}
					child2 := &ElementNode{Name: ingredient2, Parent: currentNode1, Depth: currentNode1.Depth + 1, Index: i}

					newRecipe := &RecipeNode{Ingredient1: child1, Ingredient2: child2}
					queue.enqueue(newRecipe)

					currentNode1.Recipes = append(currentNode1.Recipes, newRecipe)
					i++
				}

				rows.Close()
			}
		}

		// recipeCount = sumSlice(root.ValidRecipeIdx)
		// fmt.Printf("RECIPE COUNT: %d\n", recipeCount)
		// if recipeCount >= targetCount {
		// 	fmt.Printf("%d VALID RECIPE IDX\n", targetCount)
		// 	break
		// }
		if currentNode2 != nil {
			if isBasicElement(currentNode2.Name) {
				// ingredient 2 basic element
				currentNode2.ValidRecipeIdx = append(currentNode2.ValidRecipeIdx, 1)
				currentNode2.setValid()
			} else {
				typeQuery := "SELECT type FROM elements WHERE name = $1"
				row := db.QueryRow(typeQuery, currentNode2.Name)
				var elementType int
				err := row.Scan(&elementType)
				if err == sql.ErrNoRows {
					return nil, 0, nil
				} else if err != nil {
					return nil, 0, nil
				}

				query := "SELECT ingredient1, ingredient2 FROM recipes WHERE element = $1"
				rows, err := db.Query(query, currentNode2.Name)
				if err == sql.ErrNoRows {
					continue
				} else if err != nil {
					return nil, 0, nil
				}

				i := 0
				for rows.Next() {
					// if i >= targetCount {
					// 	break
					// }
					var ingredient1, ingredient2 string
					if err := rows.Scan(&ingredient1, &ingredient2); err != nil {
						rows.Close()
						return nil, 0, err
					}

					// Do not continue path if recipes are higher type
					query := "SELECT type FROM elements WHERE name = $1"
					row := db.QueryRow(query, ingredient1)
					var elementType1 int
					err := row.Scan(&elementType1)
					if err == sql.ErrNoRows {
						continue
					} else if err != nil {
						return nil, 0, nil
					}

					if elementType1 >= elementType {
						continue
					}

					query = "SELECT type FROM elements WHERE name = $1"
					row = db.QueryRow(query, ingredient2)
					var elementType2 int
					err = row.Scan(&elementType2)
					if err == sql.ErrNoRows {
						continue
					} else if err != nil {
						return nil, 0, nil
					}

					if elementType2 >= elementType {
						continue
					}

					// if !isElementNotLooping(currentNode2.Parent, ingredient1) || !isElementNotLooping(currentNode2.Parent, ingredient2) {
					// 	continue
					// }

					child1 := &ElementNode{Name: ingredient1, Parent: currentNode2, Depth: currentNode2.Depth + 1, Index: i}
					child2 := &ElementNode{Name: ingredient2, Parent: currentNode2, Depth: currentNode2.Depth + 1, Index: i}

					newRecipe := &RecipeNode{Ingredient1: child1, Ingredient2: child2}
					queue.enqueue(newRecipe)

					currentNode2.Recipes = append(currentNode2.Recipes, newRecipe)
					i++
				}

				rows.Close()
			}
		}
	}

	if shortestRecipeNode == nil {
		return nil, 0, fmt.Errorf("no valid recipe found for element: %s", element)
	}

	// fmt.Printf("%d anjer", targetCount)
	cutTree(root)
	return root, minDepth, nil
}

func cutTree(node *ElementNode) {
	if node == nil {
		return
	}

	var indicesToDelete []int
	// fmt.Printf("cutting %s\n", node.Name)
	for i, child := range node.Recipes {
		// fmt.Printf("cutting in %s %s\n", child.Ingredient1.Name, child.Ingredient2.Name)
		if !child.Ingredient1.IsValid || !child.Ingredient2.IsValid {
			// fmt.Printf("invalid in %t %s %t %s\n", child.Ingredient1.IsValid, child.Ingredient1.Name, child.Ingredient2.IsValid, child.Ingredient2.Name)
			indicesToDelete = append(indicesToDelete, i)
		}
		// cutTree(child.Ingredient1)
		// cutTree(child.Ingredient2)
	}

	for j := len(indicesToDelete) - 1; j >= 0; j-- {
		idx := indicesToDelete[j]
		node.Recipes = slices.Delete(node.Recipes, idx, idx+1)
	}
	for _, child := range node.Recipes {
		cutTree(child.Ingredient1)
		cutTree(child.Ingredient2)
	}
}

func isBasicElement(element string) bool {
	return element == "Water" || element == "Air" || element == "Fire" || element == "Earth"
}

func sumSlice(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}
