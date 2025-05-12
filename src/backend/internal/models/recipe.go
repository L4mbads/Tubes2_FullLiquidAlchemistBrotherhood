package models

import "sync"

type RecipeType struct {
	Element     string
	Ingredient1 string
	Ingredient2 string
}

type RecipeNode struct {
	Ingredient1 *ElementNode `json:"ingredient1"`
	Ingredient2 *ElementNode `json:"ingredient2"`
}

type RecipeQueue struct {
	recipe []*RecipeNode
	mu     sync.Mutex
}

func (queue *RecipeQueue) enqueue(recipe *RecipeNode) {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	queue.recipe = append(queue.recipe, recipe) // Simply append to enqueue.
	// fmt.Println("Enqueued:", *recipe)
}

func (queue *RecipeQueue) dequeue() *RecipeNode {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	if len(queue.recipe) == 0 {
		return nil
	}
	element := (queue.recipe)[0] // The first element is the one to be dequeued.
	queue.recipe = (queue.recipe)[1:]
	return element
}

func (queue *RecipeQueue) isEmpty() bool {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	return len(queue.recipe) == 0
}
