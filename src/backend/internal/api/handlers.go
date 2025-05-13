package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"flab/internal/db"
	"flab/internal/models"

	"github.com/gorilla/mux"
)

func GetElementsHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		elements, err := db.GetElements(dbConn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(elements)
	}
}

func GetRecipesHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipes, err := db.GetRecipes(dbConn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(recipes)
	}
}

func GetElementHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		elementToQuery := vars["element"]

		element, err := db.GetElement(dbConn, elementToQuery)

		if err != nil {
			http.Error(w, createErrorResponse("Element not found", http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(element)

	}

}

func GetRecipeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		element := r.URL.Query().Get("element")
		strategy := strings.ToLower(r.URL.Query().Get("strategy"))
		count, err := strconv.Atoi(r.URL.Query().Get("count"))

		if element == "" {
			http.Error(w, createErrorResponse("Element parameter required", http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, createErrorResponse("Target count parameter error", http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if strategy == "dfs" {
			// root, err := models.DFS(db, nil, element, count)
			ctx := context.Background()
			sem := make(chan struct{}, runtime.NumCPU()*4)
			root, err := models.DFS(ctx, db, nil, element, count, sem)
			if err != nil {
				http.Error(w, createErrorResponse(err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			models.CutTree(root)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(root)

		} else {
			root, err := models.BFS(db, element, count)
			if err != nil {
				http.Error(w, createErrorResponse(err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(root)
		}
		fmt.Println("done")

	}
}

func createErrorResponse(message string, statusCode int) string {
	errorResponse := map[string]interface{}{
		"error":      true,
		"message":    message,
		"statusCode": statusCode,
	}
	errorJSON, _ := json.Marshal(errorResponse)
	return string(errorJSON)
}
