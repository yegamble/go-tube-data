package main

import (
	"github.com/joho/godotenv"
	"github.com/yegamble/go-tube-api/modules/api/models"
	"github.com/yegamble/go-tube-api/modules/router"
	"log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	err = models.SyncModels()
	if err != nil {
		log.Fatalf("Error Syncing Models to Database")
	}
	router.SetRoutes()
	models.ScheduleCleanup()
}
