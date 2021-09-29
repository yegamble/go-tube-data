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

	dsn := os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
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
		dsnRedis = os.Getenv("REDIS_USER")+":"+os.Getenv("REDIS_PASSWORD")+"@localhost:"+os.Getenv("DB_USER")+":"+os.Getenv("DB_PORT")
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
		&View{},
		&UserBlock{},
		&Subscription{},
		&IPLog{},
		&BannedIPLog{},
		&ConversionQueue{},
		&UserSettings{},
		&config.Config{},
	)
}
