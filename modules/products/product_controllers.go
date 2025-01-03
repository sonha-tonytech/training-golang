package products

import (
	"encoding/json"
	"my-pp/share/variables"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func createItemHandler(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()
	var item ItemData

	err := json.NewDecoder(r.Body).Decode(&item)

	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	createdItem, err := CreateItem(variables.DB, id.String(), item.Body.Title, item.Body.Text, item.Body.UserId, item.Body.UpdateAt)
	if err != nil {
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdItem)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func getItemsHandler(w http.ResponseWriter, r *http.Request) {
	_start := r.URL.Query().Get("_start")
	_limit := r.URL.Query().Get("_limit")

	start, err := strconv.Atoi(_start)
	if err != nil {
		http.Error(w, "Invalid start value", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(_limit)
	if err != nil {
		http.Error(w, "Invalid limit value", http.StatusBadRequest)
		return
	}

	items, err := GetItems(variables.DB, start, limit)
	if err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
	}
}

func getItemByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemId := vars["id"]

	item, err := GetItemById(variables.DB, itemId)
	if err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
	}
}

func updateItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	var item ItemData
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	updatedItem, err := UpdateItem(variables.DB, idStr, item.Body.Title, item.Body.Text)
	if err != nil {
		http.Error(w, "Item not found or failed to update", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(updatedItem)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	item, err := DeleteItem(variables.DB, idStr)
	if err != nil {
		http.Error(w, "Item not found or could not be deleted", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
