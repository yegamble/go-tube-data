package gtredis

import (
	"github.com/go-redis/redis/v7"
	"os"
)

var (
	Client *redis.Client
)

func StartRedis() {
	//Initializing redis
	dsnRedis := os.Getenv("REDIS_DSN")
	if len(dsnRedis) == 0 {
		dsnRedis = "localhost:6379"
	}
	Client = redis.NewClient(&redis.Options{
		Addr: dsnRedis, //redis port
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
