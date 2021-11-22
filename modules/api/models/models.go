package models

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	db      *gorm.DB
	redisDB *redis.Client
)

func init() {

	err := godotenv.Load(".env")

	var dsn string
	if len(os.Getenv("ENV_VAR")) == 0 {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	} else {
		panic("environment variables not set")
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

	redisDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:" + os.Getenv("REDIS_PORT"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})

	result, err := redisDB.Ping().Result()
	if err != nil {
		log.Println(result)
		panic(err)
	}
}

func SyncModels() error {

	err := db.AutoMigrate(
		&User{},
		&Session{},
		&MessageActivityLog{},
		&MessageThread{},
		&MessageThreadParticipant{},
		&WatchLaterQueue{},
		&Video{},
		&VideoFile{},
		&BlockedUserRecord{},
		&Log{},
		&UserView{},
		&BannedIP{},
		&ConversionQueue{},
		&UserSettings{},
		&Category{},
		&Config{},
	)
	if err != nil {
		return err
	}

	return nil
}
