package models

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/yegamble/go-tube-api/modules/api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var (
	db *gorm.DB
)

func init() {
	err := godotenv.Load(".env")

	dsn := "root@tcp(127.0.0.1:3306)/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connection successfully opened")
}

func SyncModels() {

	db.AutoMigrate(
		&User{},
		&WatchLaterVideo{},
		&Video{},
		&UserBlock{},
		&Subscription{},
		&IPLog{},
		&config.Config{},
	)
}
