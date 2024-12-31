package src

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"my-pp/share/types"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	dbDriver = "mysql"
	dbUser   = "root"
	dbPass   = "letmein"
	dbName   = "golang_crud"
)

func CreateItem(db *sql.DB, id, title, text, userId string, updateAt time.Time) (types.ItemData, error) {
	query := "INSERT INTO golang_crud.list (id, title, text, updated_at, user_id) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, id, title, text, updateAt, userId)
	if err != nil {
		return types.ItemData{}, err
	}

	return GetItemById(db, id)
}

func GetItems(db *sql.DB, start, limit int) ([]types.ItemData, error) {
	query := `SELECT 
				l.id, l.title, l.text, l.updated_at, l.user_id,
				u.id, u.name, u.user_name, u.password
			  FROM golang_crud.list l
			  LEFT JOIN golang_crud.user u ON l.user_id = u.id
			  ORDER BY l.updated_at desc
			  LIMIT ? OFFSET ?`

	rows, err := db.Query(query, limit, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []types.ItemData{}

	for rows.Next() {
		var item types.ItemData
		var updatedAt string

		item.Body.User = &types.User{}

		err := rows.Scan(
			&item.ID,
			&item.Body.Title,
			&item.Body.Text,
			&updatedAt,
			&item.Body.UserId,
			&item.Body.User.ID,
			&item.Body.User.Body.Name,
			&item.Body.User.Body.UserName,
			&item.Body.User.Body.Password,
		)
		if err != nil {
			return nil, err
		}

		item.Body.UpdateAt, err = time.Parse("2006-01-02 15:04:05", updatedAt)
		if err != nil {
			return nil, err
		}

		if item.Body.User.ID == "" {
			item.Body.User = nil
		}

		data = append(data, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}
func GetItemById(db *sql.DB, itemId string) (types.ItemData, error) {
	query := `SELECT l.id, l.title, l.text, l.updated_at, l.user_id,
				u.id, u.name, u.user_name, u.password 
			  FROM golang_crud.list l
			  LEFT JOIN golang_crud.user u ON l.user_id = u.id
			  WHERE l.id = ?`
	row := db.QueryRow(query, itemId)

	var item types.ItemData
	var updatedAt string
	item.Body.User = &types.User{}

	err := row.Scan(
		&item.ID,
		&item.Body.Title,
		&item.Body.Text,
		&updatedAt,
		&item.Body.UserId,
		&item.Body.User.ID,
		&item.Body.User.Body.Name,
		&item.Body.User.Body.UserName,
		&item.Body.User.Body.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.ItemData{}, fmt.Errorf("item not found")
		}
		return types.ItemData{}, err
	}

	item.Body.UpdateAt, err = time.Parse("2006-01-02 15:04:05", updatedAt)
	if err != nil {
		return types.ItemData{}, err
	}

	return item, nil
}

func UpdateItem(db *sql.DB, id, title, text string) (types.ItemData, error) {
	query := "UPDATE golang_crud.list SET title = ?, text = ? WHERE id = ?"
	_, err := db.Exec(query, title, text, id)
	if err != nil {
		return types.ItemData{}, err
	}

	return GetItemById(db, id)
}

func DeleteItem(db *sql.DB, id string) (types.ItemData, error) {
	item, err := GetItemById(db, id)
	if err != nil {
		return types.ItemData{}, err
	}

	query := "DELETE FROM golang_crud.list WHERE id = ?"
	_, err = db.Exec(query, id)
	if err != nil {
		return types.ItemData{}, err
	}

	return item, nil
}

func GetUserLogin(db *sql.DB, userName, password string) (types.User, error) {
	query := "SELECT id, name, user_name, password FROM user WHERE user_name = ? AND password = ?"
	row := db.QueryRow(query, userName, password)

	var user types.User
	err := row.Scan(&user.ID, &user.Body.Name, &user.Body.UserName, &user.Body.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.User{}, fmt.Errorf("user not found")
		}
		return types.User{}, err
	}
	return user, nil
}

func createItemHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	id := uuid.New()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var item types.ItemData
	err = json.NewDecoder(r.Body).Decode(&item)

	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	createdItem, err := CreateItem(db, id.String(), item.Body.Title, item.Body.Text, item.Body.UserId, item.Body.UpdateAt)
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
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

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

	items, err := GetItems(db, start, limit)
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
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	itemId := vars["id"]

	item, err := GetItemById(db, itemId)
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
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	var item types.ItemData
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	updatedItem, err := UpdateItem(db, idStr, item.Body.Title, item.Body.Text)
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
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	item, err := DeleteItem(db, idStr)
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

func getUserLoginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	userName := r.URL.Query().Get("userName")
	password := r.URL.Query().Get("password")

	user, err := GetUserLogin(db, userName, password)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
