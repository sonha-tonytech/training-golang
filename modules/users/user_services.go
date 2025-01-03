package users

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func GetUserLogin(db *sql.DB, userName, password string) (User, error) {
	query := "SELECT id, name, user_name, password FROM user WHERE user_name = ? AND password = ?"
	row := db.QueryRow(query, userName, password)

	var user User
	err := row.Scan(&user.ID, &user.Body.Name, &user.Body.UserName, &user.Body.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, err
	}
	return user, nil
}

func GetUserById(db *sql.DB, userId string) (User, error) {
	query := `SELECT * FROM golang_crud.user WHERE id = ?`
	row := db.QueryRow(query, userId)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Body.Name,
		&user.Body.UserName,
		&user.Body.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, err
	}

	return user, nil
}

func RegisterUser(db *sql.DB, id, name, userName, password string) (User, error) {
	query := "INSERT INTO golang_crud.user (id, name, user_name, password) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, id, name, userName, password)
	if err != nil {
		return User{}, err
	}

	return GetUserById(db, id)
}

func UpdateUser(db *sql.DB, id, name, password string) (User, error) {
	query := "UPDATE golang_crud.user SET name = ?, password = ? WHERE id = ?"
	_, err := db.Exec(query, name, password, id)
	if err != nil {
		return User{}, err
	}

	return GetUserById(db, id)
}
