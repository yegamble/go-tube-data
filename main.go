package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/yegamble/go-tube-api/database"
	"github.com/yegamble/go-tube-api/modules/api/config"
	"github.com/yegamble/go-tube-api/modules/api/user"
	"github.com/yegamble/go-tube-api/modules/api/video"
	"github.com/yegamble/go-tube-api/modules/router"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {

	//user.UserFormParser()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//Initialise Database
	err = initDatabase()
	if err != nil {
		log.Panic(err)
	}
	router.SetRoutes()
}

func initDatabase() error {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open(os.Getenv("DB_NAME")))
	if err != nil {
		panic("Failed to Connect to Database")
	}
	fmt.Println("Database connection successfully opened")

	result := database.DBConn.AutoMigrate(&user.User{},
		&user.WatchLaterVideo{}, &video.Video{},
		&user.UserBlock{}, &user.Subscription{},
		&user.IPLog{}, &config.Config{})

	return result
}
