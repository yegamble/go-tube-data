package main

import (
	"fmt"
	"github.com/yegamble/go-tube-api/database"
	"github.com/yegamble/go-tube-api/modules/api/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Hello world, the web server
	initDatabase()
	//user.UserFormParser()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("goTube.db"))
	if err != nil {
		panic("Failed to Connect to Database")
	}
	fmt.Println("Database connection successfully opened")

	database.DBConn.AutoMigrate(&user.User{})
}