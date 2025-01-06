package main

import (
	"log"
	"my-pp/share/utils"
	"my-pp/share/variables"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
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

	err = utils.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	//Add cronjob backup db
	s := gocron.NewScheduler(time.Local)
	s.Every(1).Day().At("07:00").Do(utils.BackupDatabase)
	s.StartAsync()

	// Create a new router
	r := utils.SetupRoutes()

	// Start the HTTP server on port 3000
	log.Println("Server listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
