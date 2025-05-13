package api

import (
	"database/sql"

	"github.com/gorilla/mux"
)

func NewRouter(dbConn *sql.DB) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/go/elements", GetElementsHandler(dbConn)).Methods("GET")
	router.HandleFunc("/api/go/recipes", GetRecipesHandler(dbConn)).Methods("GET")

	router.HandleFunc("/api/go/element/{element}", GetElementHandler(dbConn)).Methods("GET")
	router.HandleFunc("/api/go/recipe", GetRecipeHandler(dbConn)).Methods("GET")
	router.HandleFunc("/api/go/liverecipe", GetLiveRecipeHandler(dbConn)).Methods("GET")
	return router
}
