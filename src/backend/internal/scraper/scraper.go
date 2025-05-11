package scraper

import (
	"database/sql"
	"flab/internal/db"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func getElementType(index int) int {
	switch index {
	case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17:
		return index - 1
	default:
		return -1
	}
}

func ScrapeElements(dbConn *sql.DB) {
	url := "little-alchemy.fandom.com"
	c := colly.NewCollector(colly.AllowedDomains(url))

	tableIndex := 0
	elementCounter := 0
	recipeCounter := 0
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

			aTags := h.DOM.Find("td:nth-of-type(1) a")
			imgUrl, _ := aTags.Eq(0).Find("img").Attr("data-src")

			err := db.InsertElement(dbConn, element, imgUrl, elementType)
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
				err := db.InsertRecipe(dbConn, element, ingredient1, ingredient2)
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

	c.Visit(url)
}
