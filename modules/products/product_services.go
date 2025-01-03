package products

import (
	"database/sql"
	"fmt"
	"my-pp/modules/users"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func CreateItem(db *sql.DB, id, title, text, userId string, updateAt time.Time) (ItemData, error) {
	query := "INSERT INTO golang_crud.list (id, title, text, updated_at, user_id) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, id, title, text, updateAt, userId)
	if err != nil {
		return ItemData{}, err
	}

	return GetItemById(db, id)
}

func GetItems(db *sql.DB, start, limit int) ([]ItemData, error) {
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

	data := []ItemData{}

	for rows.Next() {
		var item ItemData
		var updatedAt string

		item.Body.User = &users.User{}

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
func GetItemById(db *sql.DB, itemId string) (ItemData, error) {
	query := `SELECT l.id, l.title, l.text, l.updated_at, l.user_id,
				u.id, u.name, u.user_name, u.password 
			  FROM golang_crud.list l
			  LEFT JOIN golang_crud.user u ON l.user_id = u.id
			  WHERE l.id = ?`
	row := db.QueryRow(query, itemId)

	var item ItemData
	var updatedAt string
	item.Body.User = &users.User{}

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
			return ItemData{}, fmt.Errorf("item not found")
		}
		return ItemData{}, err
	}

	item.Body.UpdateAt, err = time.Parse("2006-01-02 15:04:05", updatedAt)
	if err != nil {
		return ItemData{}, err
	}

	return item, nil
}

func UpdateItem(db *sql.DB, id, title, text string) (ItemData, error) {
	query := "UPDATE golang_crud.list SET title = ?, text = ? WHERE id = ?"
	_, err := db.Exec(query, title, text, id)
	if err != nil {
		return ItemData{}, err
	}

	return GetItemById(db, id)
}

func DeleteItem(db *sql.DB, id string) (ItemData, error) {
	item, err := GetItemById(db, id)
	if err != nil {
		return ItemData{}, err
	}

	query := "DELETE FROM golang_crud.list WHERE id = ?"
	_, err = db.Exec(query, id)
	if err != nil {
		return ItemData{}, err
	}

	return item, nil
}
