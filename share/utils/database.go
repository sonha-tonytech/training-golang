package utils

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func OpenDatabase() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var (
		dbDriver = os.Getenv("DB_DRIVER")
		dbUser   = os.Getenv("DB_USER")
		dbPass   = os.Getenv("DB_PASS")
		dbName   = os.Getenv("DB_NAME")
	)

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	return db, nil
}
