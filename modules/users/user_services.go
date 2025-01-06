package users

import (
	"database/sql"
	"fmt"
	"my-pp/share/variables"

	_ "github.com/go-sql-driver/mysql"
)

func GetUserLogin(userName, password string) (User, error) {
	query := "SELECT id, name, user_name, password FROM user WHERE user_name = ? AND password = ?"
	row := variables.DB.QueryRow(query, userName, password)

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

func GetUserById(userId string) (User, error) {
	query := `SELECT * FROM golang_crud.user WHERE id = ?`
	row := variables.DB.QueryRow(query, userId)

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

func RegisterUser(id, name, userName, password string) (User, error) {
	query := "INSERT INTO golang_crud.user (id, name, user_name, password) VALUES (?, ?, ?, ?)"
	_, err := variables.DB.Exec(query, id, name, userName, password)
	if err != nil {
		return User{}, err
	}

	return GetUserById(id)
}

func UpdateUser(id, name, password string) (User, error) {
	query := "UPDATE golang_crud.user SET name = ?, password = ? WHERE id = ?"
	_, err := variables.DB.Exec(query, name, password, id)
	if err != nil {
		return User{}, err
	}

	return GetUserById(id)
}
