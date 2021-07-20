package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/yegamble/go-tube-api/database"
	"github.com/yegamble/go-tube-api/modules/api/user"
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
	initDatabase()
	router.SetRoutes()
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open(os.Getenv("DB_NAME")))
	if err != nil {
		panic("Failed to Connect to Database")
	}
	fmt.Println("Database connection successfully opened")

	database.DBConn.AutoMigrate(&user.User{})
}
