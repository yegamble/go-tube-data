package models

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
	"github.com/yegamble/go-tube-api/modules/api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var (
	db      *gorm.DB
	redisDB *redis.Client
)

func init() {
	err := godotenv.Load(".env")

	dsn := "root@tcp(127.0.0.1:3306)/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connection successfully opened")

	StartRedis()
}

func StartRedis() {
	//Initializing redis
	dsnRedis := os.Getenv("REDIS_DSN")
	if len(dsnRedis) == 0 {
		dsnRedis = "localhost:6379"
	}
	redisDB = redis.NewClient(&redis.Options{
		Addr: dsnRedis, //redis port
	})
	_, err := redisDB.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func SyncModels() {

	db.AutoMigrate(
		&User{},
		&Session{},
		&WatchLaterQueue{},
		&Video{},
		&UserBlock{},
		&Subscription{},
		&IPLog{},
		&BannedIPLog{},
		&ConversionQueue{},
		&config.Config{},
	)
}
