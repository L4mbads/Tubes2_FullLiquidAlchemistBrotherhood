package models

import (
	"database/sql"
	"fmt"
	"slices"
	"sync"
)

type Recipe struct {
	Ingredient1 string
	Ingredient2 string
}

// func DFS(db *sql.DB, parentNode *ElementNode, element string, targetCount int) (*ElementNode, error) {
// 	node := &ElementNode{Name: element, Parent: parentNode, IsValid: false}

// 	fmt.Printf("Processing %s \n", element)

// 	if isBasicElement(element) {
// 		fmt.Printf("%s leaf\n", element)
// 		node.ValidRecipeIdx = append(node.ValidRecipeIdx, 1)
// 		node.setValid()
// 		return node, nil
// 	}

// 	typeQuery := "SELECT type FROM elements WHERE name = $1"
// 	row := db.QueryRow(typeQuery, element)
// 	var elementType int
// 	err := row.Scan(&elementType)
// 	if err != nil {
// 		return nil, err
// 	}

// 	query := "SELECT ingredient1, ingredient2 FROM recipes WHERE element = $1"
// 	rows, err := db.Query(query, element)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var recipes []Recipe
// 	for rows.Next() {
// 		var ing1, ing2 string
// 		if err := rows.Scan(&ing1, &ing2); err != nil {
// 			continue
// 		}

// 		var type1, type2 int
// 		row = db.QueryRow(typeQuery, ing1)
// 		if err := row.Scan(&type1); err != nil || type1 >= elementType {
// 			continue
// 		}

// 		row = db.QueryRow(typeQuery, ing2)
// 		if err := row.Scan(&type2); err != nil || type2 >= elementType {
// 			continue
// 		}

// 		recipes = append(recipes, Recipe{Ingredient1: ing1, Ingredient2: ing2})
// 	}

// 	rows.Close()
// 	if len(recipes) == 0 {
// 		return nil, fmt.Errorf("no valid recipe found for %s", element)
// 	}
// 	fmt.Println("recipe found")

// 	var mu sync.Mutex
// 	var wg sync.WaitGroup
// 	root := node
// 	for root.Parent != nil {
// 		root = root.Parent
// 	}

// 	validRecipeCount := 0
// 	for _, recipe := range recipes {
// 		if validRecipeCount >= targetCount {
// 			break
// 		}
// 		if sumSlice(root.ValidRecipeIdx) >= targetCount {
// 			break
// 		}

// 		wg.Add(1)
// 		go func(idx int, r Recipe) {
// 			defer wg.Done()

// 			recipeNode := &RecipeNode{}
// 			childNode := &ElementNode{Name: element, Parent: parentNode}
// 			childNode.Index = idx
// 			childNode.Recipes = []*RecipeNode{recipeNode}

// 			var child1, child2 *ElementNode
// 			var err1, err2 error

// 			var innerWg sync.WaitGroup
// 			innerWg.Add(2)

// 			go func() {
// 				defer innerWg.Done()
// 				child1, err1 = DFS(db, childNode, r.Ingredient1, targetCount)
// 			}()
// 			go func() {
// 				defer innerWg.Done()
// 				child2, err2 = DFS(db, childNode, r.Ingredient2, targetCount)
// 			}()
// 			innerWg.Wait()

// 			if err1 == nil && err2 == nil && child1 != nil && child2 != nil {
// 				child1.Index = idx
// 				child2.Index = idx
// 				recipeNode.Ingredient1 = child1
// 				recipeNode.Ingredient2 = child2
// 				child1.setValid()
// 				child2.setValid()

// 				mu.Lock()
// 				node.Recipes = append(node.Recipes, recipeNode)
// 				node.ValidRecipeIdx = append(node.ValidRecipeIdx, 1)
// 				validRecipeCount++
// 				mu.Unlock()
// 			}
// 		}(validRecipeCount, recipe)
// 	}

// 	wg.Wait()
// 	cutTree(node)
// 	return node, nil
// }

// func DFS(ctx context.Context, db *sql.DB, parentNode *ElementNode, element string, targetCount int, sem chan struct{}) (*ElementNode, error) {
// 	node := &ElementNode{Name: element, Parent: parentNode, IsValid: false}
// 	if isBasicElement(element) {
// 		fmt.Printf("leaf %s\n", element)
// 		node.ValidRecipeIdx = append(node.ValidRecipeIdx, 1)
// 		node.setValid()
// 		return node, nil
// 	} else {
// 		fmt.Printf("node %s\n", element)
// 	}

// 	typeQuery := "SELECT type FROM elements WHERE name = $1"
// 	row := db.QueryRowContext(ctx, typeQuery, element)
// 	var elementType int
// 	if err := row.Scan(&elementType); err != nil {
// 		return nil, err
// 	}

// 	query := "SELECT ingredient1, ingredient2 FROM recipes WHERE element = $1"
// 	rows, err := db.QueryContext(ctx, query, element)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var recipes []Recipe
// 	for rows.Next() {
// 		var ing1, ing2 string
// 		if err := rows.Scan(&ing1, &ing2); err != nil {
// 			continue
// 		}

// 		var type1, type2 int
// 		row = db.QueryRowContext(ctx, typeQuery, ing1)
// 		if err := row.Scan(&type1); err != nil || type1 >= elementType {
// 			continue
// 		}

// 		row = db.QueryRowContext(ctx, typeQuery, ing2)
// 		if err := row.Scan(&type2); err != nil || type2 >= elementType {
// 			continue
// 		}

// 		recipes = append(recipes, Recipe{Ingredient1: ing1, Ingredient2: ing2})
// 	}

// 	if len(recipes) == 0 {
// 		return nil, fmt.Errorf("no valid recipe found for %s", element)
// 	}

// 	var mu sync.Mutex
// 	var wg sync.WaitGroup
// 	root := node
// 	for root.Parent != nil {
// 		root = root.Parent
// 	}

// 	validCount := 0
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	for _, recipe := range recipes {
// 		select {
// 		case <-ctx.Done():
// 			break
// 		default:
// 		}

// 		mu.Lock()
// 		if validCount >= targetCount || sumSlice(root.ValidRecipeIdx) >= targetCount {
// 			mu.Unlock()
// 			break
// 		}
// 		mu.Unlock()

// 		wg.Add(1)
// 		go func(idx *int, r Recipe) {
// 			defer wg.Done()

// 			// Limit concurrency
// 			select {
// 			case sem <- struct{}{}:
// 			case <-ctx.Done():
// 				return
// 			}
// 			defer func() { <-sem }()

// 			ingredient1 := recipe.Ingredient1
// 			ingredient2 := recipe.Ingredient2

// 			hasValidRecipe := true
// 			currentRecipeNode := &RecipeNode{Ingredient1: nil, Ingredient2: nil}
// 			mu.Lock()
// 			node.Recipes = append(node.Recipes, currentRecipeNode)
// 			mu.Unlock()
// 			if ingredient1 != "" {
// 				// fmt.Printf("entering %s (%d, %d)\n", ingredient1, depth, minDepth)
// 				child1, err := DFS(ctx, db, node, ingredient1, targetCount, sem)
// 				if err == nil && child1 != nil {
// 					child1.Index = validCount
// 					currentRecipeNode.Ingredient1 = child1
// 					child1.setValid()
// 				} else {
// 					hasValidRecipe = false
// 				}
// 			}

// 			if hasValidRecipe && ingredient2 != "" {
// 				// fmt.Printf("entering %s (%d, %d)\n", ingredient2, depth, minDepth)
// 				child2, err := DFS(ctx, db, node, ingredient2, targetCount, sem)
// 				if err == nil && child2 != nil {
// 					child2.Index = validCount
// 					currentRecipeNode.Ingredient2 = child2
// 					child2.setValid()
// 				} else {
// 					hasValidRecipe = false
// 				}
// 			}

// 			if hasValidRecipe {
// 				validCount++
// 			} else {
// 				// remove invalid recipe
// 				node.Recipes = node.Recipes[:len(node.Recipes)-1]
// 			}

// 			// recipeNode := &RecipeNode{}
// 			// childNode := &ElementNode{Name: element, Parent: parentNode}
// 			// childNode.Index = *idx
// 			// childNode.Recipes = []*RecipeNode{recipeNode}
// 			// var child1, child2 *ElementNode
// 			// mu.Lock()
// 			// node.Recipes = append(node.Recipes, recipeNode)
// 			// mu.Unlock()

// 			// var err1, err2 error
// 			// var innerWg sync.WaitGroup
// 			// innerWg.Add(2)

// 			// go func() {
// 			// 	defer innerWg.Done()
// 			// 	child1, err1 = DFS(ctx, db, childNode, r.Ingredient1, targetCount, sem)
// 			// }()
// 			// go func() {
// 			// 	defer innerWg.Done()
// 			// 	child2, err2 = DFS(ctx, db, childNode, r.Ingredient2, targetCount, sem)
// 			// }()
// 			// innerWg.Wait()

// 			// if err1 == nil && err2 == nil && child1 != nil && child2 != nil {
// 			// 	child1.Index = *idx
// 			// 	child2.Index = *idx

// 			// 	mu.Lock()
// 			// 	recipeNode.Ingredient1 = child1
// 			// 	recipeNode.Ingredient2 = child2
// 			// 	// child1.setValid()
// 			// 	// child2.setValid()
// 			// 	// node.Recipes = append(node.Recipes, recipeNode)
// 			// 	// node.ValidRecipeIdx = append(node.ValidRecipeIdx, 1)
// 			// 	for len(node.ValidRecipeIdx) <= validCount {
// 			// 		node.ValidRecipeIdx = append(node.ValidRecipeIdx, 0)
// 			// 	}
// 			// 	// node.ValidRecipeIdx[validCount] = sumSlice(node.Recipes[validCount].Ingredient1.ValidRecipeIdx) * sumSlice(node.Recipes[validCount].Ingredient2.ValidRecipeIdx)
// 			// 	validCount++
// 			// 	if validCount >= targetCount {
// 			// 		cancel() // early cancel!
// 			// 	}
// 			// 	mu.Unlock()
// 			// }
// 		}(&validCount, recipe)
// 	}

// 	wg.Wait()
// 	cutTree(node)
// 	return node, nil
// }

func DFS(db *sql.DB, parentNode *ElementNode, element string, targetCount int) (*ElementNode, error) {
	node := &ElementNode{Name: element, Parent: parentNode, IsValid: false}

	if isBasicElement(element) {
		node.ValidRecipeIdx = append(node.ValidRecipeIdx, 1)
		node.setValid()
		return node, nil
	}

	typeQuery := "SELECT type FROM elements WHERE name = $1"
	row := db.QueryRow(typeQuery, element)
	var elementType int
	err := row.Scan(&elementType)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	query := "SELECT ingredient1, ingredient2 FROM recipes WHERE element = $1"
	rows, err := db.Query(query, element)
	if err != nil {
		return nil, err
	}

	var recipes []Recipe
	// create a list/array to store ingredients
	for rows.Next() {
		var ingredient1, ingredient2 string
		if err := rows.Scan(&ingredient1, &ingredient2); err != nil {
			return nil, err
		}

		// do not continue path if recipes are higher type
		query := "SELECT type FROM elements WHERE name = $1"
		row := db.QueryRow(query, ingredient1)
		var elementType1 int
		err := row.Scan(&elementType1)
		if err == sql.ErrNoRows {
			continue
		} else if err != nil {
			return nil, err
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
			return nil, err
		}

		if elementType2 >= elementType {
			continue
		}

		recipe := Recipe{
			Ingredient1: ingredient1,
			Ingredient2: ingredient2,
		}

		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	rows.Close()

	if len(recipes) == 0 {
		return nil, fmt.Errorf("no valid recipe found for %s", element)
	}

	i := 0
	root := node
	for root.Parent != nil {
		root = root.Parent
	}

	for _, recipe := range recipes {
		if i >= targetCount {
			break
		}

		x := sumSlice(root.ValidRecipeIdx)
		if x >= targetCount {
			fmt.Printf("%d / %d VALID RECIPE IDX\n", targetCount, x)
			break
		}

		ingredient1 := recipe.Ingredient1
		ingredient2 := recipe.Ingredient2

		hasValidRecipe := true
		currentRecipeNode := &RecipeNode{Ingredient1: nil, Ingredient2: nil}
		node.Recipes = append(node.Recipes, currentRecipeNode)
		if ingredient1 != "" {
			// fmt.Printf("entering %s (%d, %d)\n", ingredient1, depth, minDepth)
			child1, err := DFS(db, node, ingredient1, targetCount)
			if err != nil {
				return nil, err
			}
			if child1 != nil {
				child1.Index = i
				currentRecipeNode.Ingredient1 = child1
				child1.setValid()
			} else {
				hasValidRecipe = false
			}
		}

		if hasValidRecipe && ingredient2 != "" {
			// fmt.Printf("entering %s (%d, %d)\n", ingredient2, depth, minDepth)
			child2, err := DFS(db, node, ingredient2, targetCount)
			if err != nil {
				return nil, err
			}
			if child2 != nil {
				child2.Index = i
				currentRecipeNode.Ingredient2 = child2
				child2.setValid()
			} else {
				hasValidRecipe = false
			}
		}

		if hasValidRecipe {
			i++
		} else {
			// remove invalid recipe
			node.Recipes = node.Recipes[:len(node.Recipes)-1]
		}
	}

	cutTree(node)
	return node, nil
}

func BFS(db *sql.DB, element string, targetCount int) (*ElementNode, error) {
	// initialize root node
	root := &ElementNode{Name: element, Parent: nil, IsValid: false}

	// use queue for BFS
	queue := RecipeQueue{}
	queue.enqueue(&RecipeNode{Ingredient1: root, Ingredient2: nil})

	var wg sync.WaitGroup

	for !queue.isEmpty() {

		currentRecipe := queue.dequeue()

		currentNode1 := currentRecipe.Ingredient1
		currentNode2 := currentRecipe.Ingredient2

		branchHitTarget := false
		nodeptr := currentNode1
		for nodeptr != nil {
			recipeCount := sumSlice(nodeptr.ValidRecipeIdx)
			if recipeCount >= targetCount {
				fmt.Printf("%d %s VALID RECIPE IDX\n", targetCount, nodeptr.Name)
				// early stop
				branchHitTarget = true
				break
			}
			nodeptr = nodeptr.Parent
		}
		if branchHitTarget {
			continue
		}

		// both basic element
		if currentNode1 != nil && currentNode2 != nil &&
			isBasicElement(currentNode1.Name) && isBasicElement(currentNode2.Name) {
			currentNode1.ValidRecipeIdx = append(currentNode1.ValidRecipeIdx, 1)
			currentNode2.ValidRecipeIdx = append(currentNode2.ValidRecipeIdx, 1)
			currentNode1.setValid()
			currentNode2.setValid()
			continue
		}
		// processNodeBFS(db, currentNode1, &queue)
		// processNodeBFS(db, currentNode2, &queue)

		wg.Add(2)
		go func(node *ElementNode) {
			defer wg.Done()
			if node == nil {
				return
			}
			processNodeBFS(db, node, &queue)
		}(currentNode1)

		go func(node *ElementNode) {
			defer wg.Done()
			if node == nil {
				return
			}
			processNodeBFS(db, node, &queue)
		}(currentNode2)

		wg.Wait()
	}

	// cut invalid subtrees
	cutTree(root)
	return root, nil
}

func processNodeBFS(db *sql.DB, node *ElementNode, queue *RecipeQueue) {
	if node == nil {
		return
	}

	if isBasicElement(node.Name) {
		node.ValidRecipeIdx = append(node.ValidRecipeIdx, 1)
		node.setValid()
		return
	}

	typeQuery := "SELECT type FROM elements WHERE name = $1"
	row := db.QueryRow(typeQuery, node.Name)
	var elementType int
	err := row.Scan(&elementType)
	if err == sql.ErrNoRows || err != nil {
		return
	}

	query := "SELECT ingredient1, ingredient2 FROM recipes WHERE element = $1"
	rows, err := db.Query(query, node.Name)
	if err != nil {
		return
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var ingredient1, ingredient2 string
		if err := rows.Scan(&ingredient1, &ingredient2); err != nil {
			continue
		}

		// Do not continue path if recipes are higher type
		query := "SELECT type FROM elements WHERE name = $1"
		row := db.QueryRow(query, ingredient1)
		var elementType1 int
		err := row.Scan(&elementType1)
		if err == sql.ErrNoRows {
			continue
		} else if err != nil {
			return
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
			return
		}

		if elementType2 >= elementType {
			continue
		}

		child1 := &ElementNode{Name: ingredient1, Parent: node, Index: i}
		child2 := &ElementNode{Name: ingredient2, Parent: node, Index: i}

		newRecipe := &RecipeNode{Ingredient1: child1, Ingredient2: child2}

		queue.enqueue(newRecipe)
		node.Recipes = append(node.Recipes, newRecipe)
		break

		i++
	}
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
