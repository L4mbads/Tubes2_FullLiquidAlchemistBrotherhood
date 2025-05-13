package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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

		if count <= 0 {
			http.Error(w, createErrorResponse("Target count must be positive integer", http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if strategy != "bfs" && strategy != "dfs" {
			http.Error(w, createErrorResponse("Strategy parameter error", http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if strategy == "dfs" {
			fmt.Println("DFS")

			root, err := models.DFS(db, element, count)
			if err != nil {
				http.Error(w, createErrorResponse(err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(root)

		} else {
			fmt.Println("BFS")

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

func GetLiveRecipeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

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

		if count <= 0 {
			http.Error(w, createErrorResponse("Target count must be positive integer", http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if strategy != "bfs" && strategy != "dfs" {
			http.Error(w, createErrorResponse("Strategy parameter error", http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		emit := func(node *models.ElementNode) {
			data, _ := json.Marshal(node)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}

		if strategy == "dfs" {
			models.DFSLive(db, element, count, emit)
		} else {
			models.BFSLive(db, element, count, emit)
		}
		fmt.Fprintf(w, "event: done\ndata: {}\n\n")
		flusher.Flush()
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
