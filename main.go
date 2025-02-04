package main

import (
	"log"
	"my-pp/modules/databases"
	"my-pp/modules/routes"
	"my-pp/share/variables"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//Connect to database
	db, err := databases.OpenDatabase()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	variables.DB = db

	err = databases.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	//Add cronjob backup db
	s := gocron.NewScheduler(time.Local)
	s.Every(1).Day().At("07:00").Do(databases.BackupDatabase)
	s.StartAsync()

	// Create a new router
	r := routes.SetupRoutes()

	// Start the HTTP server on port 3000
	log.Println("Server listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
