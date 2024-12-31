package main

import (
	"log"
	"my-pp/src"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Create a new router
	r := src.SetupRoutes()

	// Start the HTTP server on port 3001
	log.Println("Server listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
