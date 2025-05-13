package models

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

type Recipe struct {
	Ingredient1 string
	Ingredient2 string
}

type RecipeData struct {
	Tree        *ElementNode `json:"tree"`
	NodeCount   int          `json:"nodeCount"`
	RecipeCount int          `json:"recipeCount"`
}

var nodeCount atomic.Uint64

func DFSLive(db *sql.DB, element string, targetCount int, emit func(*ElementNode)) (RecipeData, error) {
	nodeCount.Store(0)
	ctx := context.Background()
	sem := make(chan struct{}, runtime.NumCPU())
	root, err := DFSRecursiveLive(ctx, db, nil, element, targetCount, sem, emit, 0, true)
	CutTree(root)
	return RecipeData{root, int(nodeCount.Load()), sumSlice(root.ValidRecipeIdx)}, err
}

func DFSRecursiveLive(ctx context.Context, db *sql.DB, parentNode *ElementNode, element string, targetCount int, sem chan struct{}, emit func(*ElementNode), idx int, isLeft bool) (*ElementNode, error) {
	node := &ElementNode{Name: element, Parent: parentNode, IsValid: false}
	nodeCount.Add(1)
	if parentNode != nil {
		if isLeft {
			parentNode.Recipes[idx].Ingredient1 = node
		} else {
			parentNode.Recipes[idx].Ingredient2 = node
		}
	}

	root := node
	for root.Parent != nil {
		root = root.Parent
	}

	emit(root)
	time.Sleep(1000 * time.Millisecond)

	if isBasicElement(element) {
		// fmt.Printf("%s leaf\n", element)
		node.ValidRecipeIdx = append(node.ValidRecipeIdx, 1)
		node.setValid()
		return node, nil
	}
	// fmt.Printf("%s node\n", element)

	typeQuery := "SELECT type FROM elements WHERE name = $1"
	row := db.QueryRowContext(ctx, typeQuery, element)
	var elementType int
	if err := row.Scan(&elementType); err != nil {
		return nil, err
	}

	query := "SELECT ingredient1, ingredient2 FROM recipes WHERE element = $1"
	rows, err := db.QueryContext(ctx, query, element)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var ing1, ing2 string
		if err := rows.Scan(&ing1, &ing2); err != nil {
			continue
		}

		var type1, type2 int
		row = db.QueryRowContext(ctx, typeQuery, ing1)
		if err := row.Scan(&type1); err != nil || type1 >= elementType {
			continue
		}
		row = db.QueryRowContext(ctx, typeQuery, ing2)
		if err := row.Scan(&type2); err != nil || type2 >= elementType {
			continue
		}

		recipes = append(recipes, Recipe{Ingredient1: ing1, Ingredient2: ing2})
	}

	if len(recipes) == 0 {
		return nil, fmt.Errorf("no valid recipe found for %s", element)
	}

	// tmp := targetCount / len(recipes)
	tmp := targetCount
	targetCount = int(math.Ceil(float64(targetCount) / float64(len(recipes))))

	i := 0

	tryAcquire := func() bool {
		// return true
		select {
		case sem <- struct{}{}:
			return true
		default:
			return false
		}
	}

	release := func() {
		<-sem
	}

	for _, recipe := range recipes {
		if i > 0 && tmp < 1 {
			break
		}

		var err1, err2 error
		var wg sync.WaitGroup
		if ctx.Err() != nil {
			break
		}
		recipeNode := &RecipeNode{Ingredient1: &ElementNode{Name: recipe.Ingredient1, Parent: node, Index: i}, Ingredient2: &ElementNode{Name: recipe.Ingredient2, Parent: node, Index: i}}
		// recipeNode := &RecipeNode{}
		node.Recipes = append(node.Recipes, recipeNode)

		// wg.Add(1)

		if tryAcquire() {
			wg.Add(1)
			go func() {
				defer release()
				defer wg.Done()
				recipeNode.Ingredient1, err1 = DFSRecursive(ctx, db, node, recipe.Ingredient1, targetCount, sem)
			}()
		} else {
			_, err1 = DFSRecursiveLive(ctx, db, node, recipe.Ingredient1, targetCount, sem, emit, i, true)
		}

		// if tryAcquire() {
		// 	wg.Add(1)
		// 	go func() {
		// 		defer release()
		// 		defer wg.Done()
		// 		recipeNode.Ingredient2, err2 = DFS(ctx, db, node, recipe.Ingredient2, targetCount, sem)
		// 	}()
		// } else {
		// 	recipeNode.Ingredient2, err2 = DFS(ctx, db, node, recipe.Ingredient2, targetCount, sem)
		// }

		// ss := sumSlice(recipeNode.Ingredient1.ValidRecipeIdx)
		// fmt.Printf("%s = %d resep\n", recipeNode.Ingredient1.Name, ss)
		_, err2 = DFSRecursiveLive(ctx, db, node, recipe.Ingredient2, targetCount, sem, emit, i, false)
		wg.Wait()

		hasValid := true
		if err1 != nil {
			hasValid = false
		}
		if err2 != nil {
			hasValid = false
		}

		if hasValid {
			recipeNode.Ingredient1.Index = i
			recipeNode.Ingredient2.Index = i
			// recipeNode.Ingredient1 = child1
			// recipeNode.Ingredient2 = child2

			// Perform setValid and parent updates in main thread
			recipeNode.Ingredient1.setValid()
			recipeNode.Ingredient2.setValid()

			i++
			x := sumSlice(node.ValidRecipeIdx)
			// fmt.Printf("SDFJSFH %s %d\n", element, x)
			tmp -= x
		}
		emit(root)
		time.Sleep(1000 * time.Millisecond)
	}

	// CutTree(node)
	if i == 0 {
		return node, fmt.Errorf("no valid sub-recipes for %s", element)
	} else {
		return node, nil
	}
}

func DFS(db *sql.DB, element string, targetCount int) (RecipeData, error) {
	nodeCount.Store(0)
	ctx := context.Background()
	sem := make(chan struct{}, runtime.NumCPU())
	root, err := DFSRecursive(ctx, db, nil, element, targetCount, sem)
	CutTree(root)
	return RecipeData{root, int(nodeCount.Load()), sumSlice(root.ValidRecipeIdx)}, err
}

func DFSRecursive(ctx context.Context, db *sql.DB, parentNode *ElementNode, element string, targetCount int, sem chan struct{}) (*ElementNode, error) {
	node := &ElementNode{Name: element, Parent: parentNode, IsValid: false}
	nodeCount.Add(1)

	if isBasicElement(element) {
		// fmt.Printf("%s leaf\n", element)
		node.ValidRecipeIdx = append(node.ValidRecipeIdx, 1)
		node.setValid()
		return node, nil
	}
	// fmt.Printf("%s node\n", element)

	typeQuery := "SELECT type FROM elements WHERE name = $1"
	row := db.QueryRowContext(ctx, typeQuery, element)
	var elementType int
	if err := row.Scan(&elementType); err != nil {
		return nil, err
	}

	query := "SELECT ingredient1, ingredient2 FROM recipes WHERE element = $1"
	rows, err := db.QueryContext(ctx, query, element)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var ing1, ing2 string
		if err := rows.Scan(&ing1, &ing2); err != nil {
			continue
		}

		var type1, type2 int
		row = db.QueryRowContext(ctx, typeQuery, ing1)
		if err := row.Scan(&type1); err != nil || type1 >= elementType {
			continue
		}
		row = db.QueryRowContext(ctx, typeQuery, ing2)
		if err := row.Scan(&type2); err != nil || type2 >= elementType {
			continue
		}

		recipes = append(recipes, Recipe{Ingredient1: ing1, Ingredient2: ing2})
	}

	if len(recipes) == 0 {
		return nil, fmt.Errorf("no valid recipe found for %s", element)
	}

	// tmp := targetCount / len(recipes)
	tmp := targetCount
	targetCount = int(math.Ceil(float64(targetCount) / float64(len(recipes))))

	i := 0
	// root := node
	// for root.Parent != nil {
	// 	root = root.Parent
	// }

	tryAcquire := func() bool {
		// return true
		select {
		case sem <- struct{}{}:
			return true
		default:
			return false
		}
	}

	release := func() {
		<-sem
	}

	for _, recipe := range recipes {
		if i > 0 && tmp < 1 {
			break
		}

		recipeNode := &RecipeNode{}
		node.Recipes = append(node.Recipes, recipeNode)
		var err1, err2 error
		var wg sync.WaitGroup
		if ctx.Err() != nil {
			break
		}

		// wg.Add(1)

		if tryAcquire() {
			wg.Add(1)
			go func() {
				defer release()
				defer wg.Done()
				recipeNode.Ingredient1, err1 = DFSRecursive(ctx, db, node, recipe.Ingredient1, targetCount, sem)
			}()
		} else {
			recipeNode.Ingredient1, err1 = DFSRecursive(ctx, db, node, recipe.Ingredient1, targetCount, sem)
		}

		// if tryAcquire() {
		// 	wg.Add(1)
		// 	go func() {
		// 		defer release()
		// 		defer wg.Done()
		// 		recipeNode.Ingredient2, err2 = DFS(ctx, db, node, recipe.Ingredient2, targetCount, sem)
		// 	}()
		// } else {
		// 	recipeNode.Ingredient2, err2 = DFS(ctx, db, node, recipe.Ingredient2, targetCount, sem)
		// }

		// ss := sumSlice(recipeNode.Ingredient1.ValidRecipeIdx)
		// fmt.Printf("%s = %d resep\n", recipeNode.Ingredient1.Name, ss)
		recipeNode.Ingredient2, err2 = DFSRecursive(ctx, db, node, recipe.Ingredient2, targetCount, sem)
		wg.Wait()

		hasValid := true
		if err1 != nil {
			hasValid = false
		}
		if err2 != nil {
			hasValid = false
		}

		if hasValid {
			recipeNode.Ingredient1.Index = i
			recipeNode.Ingredient2.Index = i
			// recipeNode.Ingredient1 = child1
			// recipeNode.Ingredient2 = child2

			// Perform setValid and parent updates in main thread
			recipeNode.Ingredient1.setValid()
			recipeNode.Ingredient2.setValid()

			i++
			x := sumSlice(node.ValidRecipeIdx)
			// fmt.Printf("SDFJSFH %s %d\n", element, x)
			tmp -= x
		}
	}

	// CutTree(node)
	if i == 0 {
		return node, fmt.Errorf("no valid sub-recipes for %s", element)
	} else {
		return node, nil
	}
}

// func DFS(db *sql.DB, parentNode *ElementNode, element string, targetCount int) (*ElementNode, error) {
// 	node := &ElementNode{Name: element, Parent: parentNode, IsValid: false}

// 	if isBasicElement(element) {
// 		node.ValidRecipeIdx = append(node.ValidRecipeIdx, 1)
// 		node.setValid()
// 		return node, nil
// 	}

// 	typeQuery := "SELECT type FROM elements WHERE name = $1"
// 	row := db.QueryRow(typeQuery, element)
// 	var elementType int
// 	err := row.Scan(&elementType)
// 	if err == sql.ErrNoRows {
// 		return nil, err
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	query := "SELECT ingredient1, ingredient2 FROM recipes WHERE element = $1"
// 	rows, err := db.Query(query, element)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var recipes []Recipe
// 	// create a list/array to store ingredients
// 	for rows.Next() {
// 		var ingredient1, ingredient2 string
// 		if err := rows.Scan(&ingredient1, &ingredient2); err != nil {
// 			return nil, err
// 		}

// 		// do not continue path if recipes are higher type
// 		query := "SELECT type FROM elements WHERE name = $1"
// 		row := db.QueryRow(query, ingredient1)
// 		var elementType1 int
// 		err := row.Scan(&elementType1)
// 		if err == sql.ErrNoRows {
// 			continue
// 		} else if err != nil {
// 			return nil, err
// 		}

// 		if elementType1 >= elementType {
// 			continue
// 		}

// 		query = "SELECT type FROM elements WHERE name = $1"
// 		row = db.QueryRow(query, ingredient2)
// 		var elementType2 int
// 		err = row.Scan(&elementType2)
// 		if err == sql.ErrNoRows {
// 			continue
// 		} else if err != nil {
// 			return nil, err
// 		}

// 		if elementType2 >= elementType {
// 			continue
// 		}

// 		recipe := Recipe{
// 			Ingredient1: ingredient1,
// 			Ingredient2: ingredient2,
// 		}

// 		recipes = append(recipes, recipe)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	rows.Close()

// 	if len(recipes) == 0 {
// 		return nil, fmt.Errorf("no valid recipe found for %s", element)
// 	}

// 	i := 0
// 	root := node
// 	for root.Parent != nil {
// 		root = root.Parent
// 	}

// 	for _, recipe := range recipes {
// 		if i >= targetCount {
// 			break
// 		}

// 		x := sumSlice(root.ValidRecipeIdx)
// 		if x >= targetCount {
// 			fmt.Printf("%d / %d VALID RECIPE IDX\n", targetCount, x)
// 			break
// 		}

// 		ingredient1 := recipe.Ingredient1
// 		ingredient2 := recipe.Ingredient2

// 		hasValidRecipe := true
// 		currentRecipeNode := &RecipeNode{Ingredient1: nil, Ingredient2: nil}
// 		node.Recipes = append(node.Recipes, currentRecipeNode)
// 		// node.ValidRecipeIdx = append(node.ValidRecipeIdx, 1)
// 		if ingredient1 != "" {
// 			// fmt.Printf("entering %s (%d, %d)\n", ingredient1, depth, minDepth)
// 			child1, err := DFS(db, node, ingredient1, targetCount)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if child1 != nil {
// 				child1.Index = i
// 				currentRecipeNode.Ingredient1 = child1
// 				child1.setValid()
// 			} else {
// 				hasValidRecipe = false
// 			}
// 		}

// 		if hasValidRecipe && ingredient2 != "" {
// 			// fmt.Printf("entering %s (%d, %d)\n", ingredient2, depth, minDepth)
// 			// fmt.Printf("slice %s = %d\n", ingredient1, sumSlice(currentRecipeNode.Ingredient1.ValidRecipeIdx))
// 			child2, err := DFS(db, node, ingredient2, targetCount)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if child2 != nil {
// 				child2.Index = i
// 				currentRecipeNode.Ingredient2 = child2
// 				child2.setValid()
// 			} else {
// 				hasValidRecipe = false
// 			}
// 		}

// 		if hasValidRecipe {
// 			i++
// 		} else {
// 			// remove invalid recipe
// 			node.Recipes = node.Recipes[:len(node.Recipes)-1]
// 			// node.ValidRecipeIdx = node.ValidRecipeIdx[:len(node.ValidRecipeIdx)-1]
// 		}
// 	}

//		CutTree(node)
//		return node, nil
//	}
func BFSLive(db *sql.DB, element string, targetCount int, emit func(*ElementNode)) (RecipeData, error) {
	nodeCount.Store(0)
	// initialize root node
	root := &ElementNode{Name: element, Parent: nil, IsValid: false}
	nodeCount.Add(1)

	// use queue for BFS
	queue := RecipeQueue{}
	queue.enqueue(&RecipeNode{Ingredient1: root, Ingredient2: nil})

	var wg sync.WaitGroup

	for !queue.isEmpty() {

		currentRecipe := queue.dequeue()
		currentNode1 := currentRecipe.Ingredient1
		currentNode2 := currentRecipe.Ingredient2
		nodeCount.Add(2)

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
			nodeCount.Add(2)
			emit(root)
			time.Sleep(1000 * time.Millisecond)
			continue
		}

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
		emit(root)
		time.Sleep(1000 * time.Millisecond)
	}

	// cut invalid subtrees
	CutTree(root)
	emit(root)
	return RecipeData{root, int(nodeCount.Load()), sumSlice(root.ValidRecipeIdx)}, nil
}

func BFS(db *sql.DB, element string, targetCount int) (RecipeData, error) {
	nodeCount.Store(0)
	// initialize root node
	root := &ElementNode{Name: element, Parent: nil, IsValid: false}
	nodeCount.Add(1)

	// use queue for BFS
	queue := RecipeQueue{}
	queue.enqueue(&RecipeNode{Ingredient1: root, Ingredient2: nil})

	var wg sync.WaitGroup

	for !queue.isEmpty() {
		currentRecipe := queue.dequeue()

		currentNode1 := currentRecipe.Ingredient1
		currentNode2 := currentRecipe.Ingredient2
		nodeCount.Add(2)

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
			nodeCount.Add(2)
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
	CutTree(root)
	return RecipeData{root, int(nodeCount.Load()), sumSlice(root.ValidRecipeIdx)}, nil
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

		i++
	}
}

func CutTree(node *ElementNode) {
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
		// CutTree(child.Ingredient1)
		// CutTree(child.Ingredient2)
	}

	for j := len(indicesToDelete) - 1; j >= 0; j-- {
		idx := indicesToDelete[j]
		node.Recipes = slices.Delete(node.Recipes, idx, idx+1)
	}
	for _, child := range node.Recipes {
		CutTree(child.Ingredient1)
		CutTree(child.Ingredient2)
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
