package db

import (
	"database/sql"
	"flab/internal/models"
)

func GetElements(db *sql.DB) ([]models.ElementType, error) {
	rows, err := db.Query("SELECT * FROM elements")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var elements []models.ElementType
	for rows.Next() {
		var element models.ElementType
		if err := rows.Scan(&element.Name, &element.ImageUrl, &element.Type); err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}

	return elements, rows.Err()
}

func GetElement(db *sql.DB, elementToQuery string) (models.ElementType, error) {
	query := "SELECT * FROM elements WHERE name = $1"

	row := db.QueryRow(query, elementToQuery)

	var element models.ElementType
	err := row.Scan(&element.Name, &element.ImageUrl, &element.Type)

	return element, err
}

func GetRecipes(db *sql.DB) ([]models.RecipeType, error) {
	rows, err := db.Query("SELECT * FROM recipes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []models.RecipeType
	for rows.Next() {
		var recipe models.RecipeType
		if err := rows.Scan(&recipe.Element, &recipe.Ingredient1, &recipe.Ingredient2); err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	return recipes, rows.Err()
}

func InsertElement(db *sql.DB, element string, imgUrl string, elementType int) error {
	sqlStatement := `
		INSERT INTO elements (name, image_url, type)
		VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, element, imgUrl, elementType)
	return err
}

func InsertRecipe(db *sql.DB, element string, ingredient1 string, ingredient2 string) error {
	sqlStatement := `
	INSERT INTO recipes (element, ingredient1, ingredient2)
	VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, element, ingredient1, ingredient2)
	return err
}
