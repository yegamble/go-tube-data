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

	var dsn string
	if len(os.Getenv("ENV_VAR")) > 0 {
		dsn = os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	} else {
		dsn = "gotube_admin:btQj49JylBTeweuP@tcp(localhost:3306)/gotube_db?charset=utf8mb4&parseTime=True&loc=Local"
	}
	print(dsn)
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
		dsnRedis = os.Getenv("REDIS_USER") + ":" + os.Getenv("REDIS_PASSWORD") + "@localhost:" + os.Getenv("DB_USER")
	}
	redisDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", //redis port
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
		&VideoFile{},
		&View{},
		&UserBlock{},
		&Subscription{},
		&IPLog{},
		&BannedIPLog{},
		&ConversionQueue{},
		&UserSettings{},
		&Category{},
		&config.Config{},
	)
}
