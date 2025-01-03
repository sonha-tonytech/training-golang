package main

import (
	"log"
	"my-pp/share/utils"
	"my-pp/share/variables"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//Connect to database
	db, err := utils.OpenDatabase()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	variables.DB = db

	// Create a new router
	r := SetupRoutes()

	// Start the HTTP server on port 3000
	log.Println("Server listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
