package models

type RecipeType struct {
	Element     string
	Ingredient1 string
	Ingredient2 string
}

type RecipeNode struct {
	Ingredient1 *ElementNode `json:"ingredient1"`
	Ingredient2 *ElementNode `json:"ingredient2"`
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
