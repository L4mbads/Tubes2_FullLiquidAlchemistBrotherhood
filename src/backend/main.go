package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func getElementType(index int) string {
	switch index {
	case 1:
		return "Starting"
	case 2:
		// Table[2] is special, we skip it (Ruins/Archeologist)
		return ""
	case 3:
		return "Tier1"
	case 4:
		return "Tier2"
	case 5:
		return "Tier3"
	case 6:
		return "Tier4"
	case 7:
		return "Tier5"
	case 8:
		return "Tier6"
	case 9:
		return "Tier7"
	case 10:
		return "Tier8"
	case 11:
		return "Tier9"
	case 12:
		return "Tier10"
	case 13:
		return "Tier11"
	case 14:
		return "Tier12"
	case 15:
		return "Tier13"
	case 16:
		return "Tier14"
	case 17:
		return "Tier15"
	default:
		return ""
	}
}

func main() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

	// connect to database
	// db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS recipes (id SMALLINT PRIMARY KEY, element TEXT, image1 TEXT, image2 TEXT, ingredient1 TEXT, ingredient2 TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS recipes")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS elements")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS elements (name TEXT, image_url TEXT, type TEXT)")
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
		if elementType == "" {
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
