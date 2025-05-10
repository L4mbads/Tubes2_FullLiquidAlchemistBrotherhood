package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type RecipeType struct {
	Element     string
	Ingredient1 string
	Ingredient2 string
}

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
	// fmt.Println("1")
	if recipe.Ingredient1 == nil || recipe.Ingredient2 == nil {
		return
	}
	if recipe.Ingredient1.IsValid && recipe.Ingredient2.IsValid {
		// fmt.Println("2")
		for len(node.ValidRecipeIdx) <= recipeIdx {
			// fmt.Println("3")
			node.ValidRecipeIdx = append(node.ValidRecipeIdx, 0)
		}
		// fmt.Println("4")
		node.ValidRecipeIdx[recipeIdx] = sumSlice(node.Recipes[recipeIdx].Ingredient1.ValidRecipeIdx) * sumSlice(node.Recipes[recipeIdx].Ingredient2.ValidRecipeIdx)
		// node.ValidRecipeIdx[recipeIdx]++
		fmt.Printf("[%d]%d anjay\n", recipeIdx, node.ValidRecipeIdx[recipeIdx])
		// fmt.Println("5")
		node.setValid()
	}
}

func (node *ElementNode) setValid() {
	node.IsValid = true
	if node.Parent != nil {
		node.Parent.checkValidRecipe(node.Index)
	}
}

type RecipeNode struct {
	Ingredient1 *ElementNode `json:"ingredient1"`
	Ingredient2 *ElementNode `json:"ingredient2"`
}

func getElementType(index int) int {
	switch index {
	case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17:
		return index - 1
	default:
		return -1
	}
}

func main() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

	// connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Configure the connection pool
	db.SetMaxOpenConns(10) // Adjust based on your server's capacity
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	// create table if it doesn't exist
	_, err = db.Exec("DROP TABLE IF EXISTS recipes")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS elements")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS elements (name TEXT, image_url TEXT, type SMALLINT)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS recipes (element TEXT, ingredient1 TEXT, ingredient2 TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// create router
	router := mux.NewRouter()
	router.HandleFunc("/api/go/elements", getElements(db)).Methods("GET")
	router.HandleFunc("/api/go/recipes", getRecipes(db)).Methods("GET")
	router.HandleFunc("/api/go/recipe", getRecipe(db)).Methods("GET")

	// wrap the router with CORS and JSON content type middlewares
	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))

	c := colly.NewCollector(colly.AllowedDomains("little-alchemy.fandom.com"))
	tableIndex := 0
	elementCounter := 0
	recipeCounter := 0

	// each table (starting and tiers)
	c.OnHTML("table.list-table", func(table *colly.HTMLElement) {
		tableIndex++
		elementType := getElementType(tableIndex)
		if elementType == -1 {
			return
		}

		// each element generated
		table.ForEach("tbody tr", func(_ int, h *colly.HTMLElement) {
			element := strings.TrimSpace(h.ChildText("td:first-of-type a"))
			if element == "" || element == "Time" || element == "Ruins" || element == "Archeologist" {
				return
			}

			elementCounter++
			// fmt.Printf("\nElement[%v]: %-10s | %s\n", elementCounter, element, elementType)

			aTags := h.DOM.Find("td:nth-of-type(1) a")
			imgUrl, _ := aTags.Eq(0).Find("img").Attr("data-src")

			sqlStatement := `
			INSERT INTO elements (name, image_url, type)
			VALUES ($1, $2, $3)`
			_, err = db.Exec(sqlStatement, element, imgUrl, elementType)
			if err != nil {
				panic(err)
			}

			h.ForEach("td:nth-of-type(2) li", func(_ int, li *colly.HTMLElement) {
				recipeCounter++
				aTags := li.DOM.Find("a")

				if aTags.Length() < 2 {
					return
				}

				// imgUrl1, _ := aTags.Eq(0).Find("img").Attr("data-src")
				// imgUrl2, _ := aTags.Eq(2).Find("img").Attr("data-src")
				ingredient1 := strings.TrimSpace(aTags.Eq(1).Text())
				ingredient2 := strings.TrimSpace(aTags.Eq(3).Text())

				if ingredient1 == "Time" || ingredient2 == "Time" || ingredient1 == "Ruins" || ingredient2 == "Ruins" || ingredient1 == "Archeologist" || ingredient2 == "Archeologist" {
					return
				}

				// Insert into recipes table
				sqlStatement := `
				INSERT INTO recipes (element, ingredient1, ingredient2)
				VALUES ($1, $2, $3)`
				_, err = db.Exec(sqlStatement, element, ingredient1, ingredient2)
				if err != nil {
					log.Printf("Error inserting recipe for element '%s': %v", element, err)
					return
				}

			})
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Print("Visiting ", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Print(e.Error())
	})

	err = c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	// start server
	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))
}

func getElements(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM elements")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		elements := []ElementType{}
		for rows.Next() {
			var u ElementType
			if err := rows.Scan(&u.Name, &u.ImageUrl, &u.Type); err != nil {
				log.Fatal(err)
			}
			elements = append(elements, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(elements)

	}

}

func getRecipes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM recipes")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		recipes := []RecipeType{}
		for rows.Next() {
			var u RecipeType
			if err := rows.Scan(&u.Element, &u.Ingredient1, &u.Ingredient2); err != nil {
				log.Fatal(err)
			}
			recipes = append(recipes, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(recipes)

	}

}

func getRecipe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		element := r.URL.Query().Get("element")
		strategy := strings.ToLower(r.URL.Query().Get("strategy")) // "bfs" or "dfs"
		count, err := strconv.Atoi(r.URL.Query().Get("count"))
		if element == "" {
			http.Error(w, "element parameter required", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "target count parameter error", http.StatusBadRequest)
			return
		}

		if strategy == "dfs" {
			root, _, err := shortestDFS(db, nil, element, 0, 999999, count)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(root)
		} else {
			root, _, err := shortestBFS(db, nil, element, 0, 999999, count)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(root)
		}
	}
}

func isBasicElement(element string) bool {

	if element == "Water" || element == "Air" || element == "Fire" || element == "Earth" {
		return true
	}
	return false
}

func isElementNotLooping(node *ElementNode, element string) bool {
	if node == nil {
		return true
	}
	if node.Name == element {
		return false
	}

	if node.Parent != nil {
		return isElementNotLooping(node.Parent, element)
	}
	return true
}

type Recipe struct {
	Ingredient1 string
	Ingredient2 string
}

type RecipeQueue []*RecipeNode

func (queue *RecipeQueue) enqueue(recipe *RecipeNode) {
	*queue = append(*queue, recipe) // Simply append to enqueue.
	// fmt.Println("Enqueued:", *recipe)
}

func (queue *RecipeQueue) dequeue() *RecipeNode {
	element := (*queue)[0] // The first element is the one to be dequeued.
	*queue = (*queue)[1:]
	return element
}

func (queue *RecipeQueue) isEmpty() bool {
	return len(*queue) == 0
}

func sumSlice(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

func cutTree(node *ElementNode) {
	fmt.Printf("cutting %s\n", node.Name)
	// if node.Parent != nil {
	// 	if node.Index < len(node.Parent.ValidRecipeIdx) && node.Parent.ValidRecipeIdx[node.Index] == 0 {
	// 		fmt.Printf("deleting index %d len %d", node.Index, len(node.Parent.Recipes))
	// 		node.Parent.Recipes = slices.Delete(node.Parent.Recipes, node.Index, node.Index+1)
	// 	}
	// }
	for idx, child := range node.Recipes {
		if !child.Ingredient1.IsValid || !child.Ingredient2.IsValid {
			node.Parent.Recipes = slices.Delete(node.Parent.Recipes, idx, idx+1)
		}
		cutTree(child.Ingredient1)
		cutTree(child.Ingredient2)
	}
}

func shortestBFS(db *sql.DB, parentNode *ElementNode, element string, depth int, limit int, targetCount int) (*ElementNode, int, error) {
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

		if currentNode1 != nil {
			if isBasicElement(currentNode1.Name) {
				// ingredient 1 basic element
				currentNode1.ValidRecipeIdx = append(currentNode1.ValidRecipeIdx, 1)
				currentNode1.setValid()
			} else {
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
				rows, err := db.Query(query, currentNode1.Name)
				if err != nil {
					return nil, 0, err
				}

				i := 0
				for rows.Next() {
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
						return nil, 0, nil
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

					if !isElementNotLooping(currentNode1.Parent, ingredient1) || !isElementNotLooping(currentNode1.Parent, ingredient2) {
						continue
					}

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

		if currentNode2 != nil {
			if isBasicElement(currentNode2.Name) {
				// ingredient 2 basic element
				currentNode2.ValidRecipeIdx = append(currentNode2.ValidRecipeIdx, 1)
				currentNode2.setValid()
			} else {
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
				rows, err := db.Query(query, currentNode2.Name)
				if err != nil {
					return nil, 0, err
				}

				i := 0
				for rows.Next() {
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
						return nil, 0, nil
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

					if !isElementNotLooping(currentNode2.Parent, ingredient1) || !isElementNotLooping(currentNode2.Parent, ingredient2) {
						continue
					}

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
	// cutTree(root)
	return root, minDepth, nil
}

func shortestDFS(db *sql.DB, parentNode *ElementNode, element string, depth int, limit int, targetCount int) (*ElementNode, int, error) {
	node := &ElementNode{Name: element, Parent: parentNode, IsValid: false}

	for i := 0; i < depth; i++ {
		fmt.Print("-")
	}
	fmt.Printf("%d %s\n", depth, element)

	if isBasicElement(element) {
		node.ValidRecipeIdx = append(node.ValidRecipeIdx, 1)
		fmt.Println("single done")
		// node.setValid()
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
	// leafRecipeIndex := -1
	// singleLeafRecipeIndex := -1
	// i := 0
	for rows.Next() {
		var ingredient1, ingredient2 string
		if err := rows.Scan(&ingredient1, &ingredient2); err != nil {
			return nil, 0, err
		}
		if !isElementNotLooping(parentNode, ingredient1) || !isElementNotLooping(parentNode, ingredient2) {
			continue
		}

		// do not continue path if recipes are higher type
		query := "SELECT type FROM elements WHERE name = $1"
		row := db.QueryRow(query, ingredient1)
		var elementType1 int
		err := row.Scan(&elementType1)
		if err == sql.ErrNoRows {
			return nil, 0, nil
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

		recipe := Recipe{
			Ingredient1: ingredient1,
			Ingredient2: ingredient2,
		}

		fmt.Printf("%s %s\n", ingredient1, ingredient2)

		recipes = append(recipes, recipe)

		// // prioritize leaf nodes
		// if isBasicElement(ingredient1) && isBasicElement(ingredient2) {
		// 	leafRecipeIndex = i
		// } else if singleLeafRecipeIndex == -1 && (isBasicElement(ingredient1) || isBasicElement(ingredient2)) {
		// 	singleLeafRecipeIndex = i
		// }
		// i++
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	rows.Close()

	if len(recipes) == 0 {
		return nil, 0, err
	}

	// if leafRecipeIndex != -1 {
	// 	var recipe = recipes[leafRecipeIndex]
	// 	ingredient1Node := &ElementNode{Name: recipe.Ingredient1, IsValid: true}
	// 	ingredient2Node := &ElementNode{Name: recipe.Ingredient2, IsValid: true}
	// 	recipeNode := &RecipeNode{Ingredient1: ingredient1Node, Ingredient2: ingredient2Node}

	// 	node.Recipes = append(node.Recipes, recipeNode)
	// 	node.IsValid = true
	// 	fmt.Println("done")
	// 	return node, 1, nil
	// } else if limit <= depth+1 {
	// 	return nil, 0, nil
	// } else if singleLeafRecipeIndex != -1 {
	// 	fmt.Printf("swapping 0 with %d\n", singleLeafRecipeIndex)
	// 	recipes[singleLeafRecipeIndex], recipes[0] = recipes[0], recipes[singleLeafRecipeIndex]
	// }

	// recipeNode := &RecipeNode{Ingredient1: nil, Ingredient2: nil}
	// node.Recipes = append(node.Recipes, recipeNode)
	minDepth := 999999
	i := 0
	root := node
	for root.Parent != nil {
		root = root.Parent
	}
	fmt.Printf("ROOTNYA BENER GAK %s %d %d\n", root.Name, sumSlice(root.ValidRecipeIdx), targetCount)
	for _, recipe := range recipes {
		if sumSlice(root.ValidRecipeIdx) >= targetCount {
			fmt.Printf("%d VALID RECIPE IDX\n", targetCount)
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
			child1, _, err := shortestDFS(db, node, ingredient1, depth+1, limit, targetCount)
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

		if sumSlice(root.ValidRecipeIdx) >= targetCount {
			fmt.Printf("%d VALID RECIPE IDX\n", targetCount)
			break
		}
		if hasValidRecipe && ingredient2 != "" {
			fmt.Printf("entering %s (%d, %d)\n", ingredient2, depth, minDepth)
			child2, _, err := shortestDFS(db, node, ingredient2, depth+1, limit, targetCount)
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

		// if hasValidRecipe {
		// node.setValid()
		// 	// recipeNode = currentRecipeNode
		// 	minDepth = maxDepth
		// 	limit = minDepth + depth
		// 	fmt.Printf("%s %s valid (%d)\n", ingredient1, ingredient2, depth)

		// 	break
		// }
		i++
	}
	// if minDepth != -1 {
	// 	node.Recipes = append(node.Recipes, recipeNode)
	// 	fmt.Println("ada yg elesai")
	// }
	return node, minDepth + 1, nil
}

func buildRecipeTree(db *sql.DB, element string, visited map[string]bool, depth int) (*ElementNode, error) {
	// Check if the element is already visited to avoid cycles
	// if visited[element] {
	// 	return nil, nil
	// }
	// visited[element] = true
	// if depth == 4 {
	// 	return nil, nil
	// }

	fmt.Print(depth)
	fmt.Println(element)

	// Create a node for the current element
	node := &ElementNode{Name: element}

	// If the element is a basic element, return it as a leaf node
	if isBasicElement(element) {
		return node, nil
	}

	// Query the database for all recipes of the current element
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

	for _, recipe := range recipes {
		ingredient1 := recipe.Ingredient1
		ingredient2 := recipe.Ingredient2
		// if err := rows.Scan(&ingredient1, &ingredient2); err != nil {
		// 	return nil, err
		// }

		if !isElementNotLooping(node, ingredient1) || !isElementNotLooping(node, ingredient2) {
			continue
		}

		recipeNode := &RecipeNode{Ingredient1: nil, Ingredient2: nil}

		hasValidRecipe := true
		// Recursively build the tree for each ingredient
		if ingredient1 != "" {
			child1, err := buildRecipeTree(db, ingredient1, visited, depth+1)
			if err != nil {
				return nil, err
			}
			if child1 != nil {
				child1.Parent = node
				recipeNode.Ingredient1 = child1
			} else {
				hasValidRecipe = false
			}
		}
		if ingredient2 != "" {
			child2, err := buildRecipeTree(db, ingredient2, visited, depth+1)
			if err != nil {
				return nil, err
			}
			if child2 != nil {
				child2.Parent = node
				recipeNode.Ingredient2 = child2
			} else {
				hasValidRecipe = false
			}
		}
		// Add the recipe node as a child of the current element node
		if hasValidRecipe { // && isElementNotLooping(node, ingredient1) && isElementNotLooping(node, ingredient2) {
			node.Recipes = append(node.Recipes, recipeNode)
		}
		if isBasicElement(ingredient1) && isBasicElement(ingredient2) {
			break
		}
	}

	return node, nil
}

// Recursive function to build the recipe tree
// func buildRecipeTree(db *sql.DB, element string, visited map[string]bool) (*RecipeNode, error) {
// 	// Check if the element is already visited to avoid cycles
// 	if visited[element] {
// 		return nil, nil
// 	}
// 	visited[element] = true

// 	// Create a node for the current element
// 	node := &RecipeNode{Name: element}

// 	if isBasicElement(element) {
// 		return node, nil
// 	}

// 	// Query the database for the recipe of the current element
// 	query := "SELECT ingredient1, ingredient2 FROM recipes WHERE element = $1"
// 	row := db.QueryRow(query, element)

// 	var ingredient1, ingredient2 string
// 	err := row.Scan(&ingredient1, &ingredient2)
// 	if err == sql.ErrNoRows {
// 		// No recipe found for this element, return the node as a leaf
// 		return node, nil
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	// Recursively build the tree for each ingredient
// 	if ingredient1 != "" {
// 		child1, err := buildRecipeTree(db, ingredient1, visited)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if child1 != nil {
// 			node.Recipes = append(node.Recipes, child1)
// 		}
// 	}
// 	if ingredient2 != "" {
// 		child2, err := buildRecipeTree(db, ingredient2, visited)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if child2 != nil {
// 			node.Recipes = append(node.Recipes, child2)
// 		}
// 	}

// 	return node, nil
// }

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Check if the request is for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set JSON Content-Type
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
