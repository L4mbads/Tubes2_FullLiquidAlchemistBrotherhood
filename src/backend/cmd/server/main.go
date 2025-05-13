package main

import (
	"flab/internal/api"
	"flab/internal/db"
	"log"
	"net/http"
)

func main() {
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	// set up router
	router := api.NewRouter(dbConn)

	// cors
	enhancedRouter := api.EnableCORS(api.JSONContentTypeMiddleware(router))

	// scraper.ScrapeElements(dbConn)

	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))
}
