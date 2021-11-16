package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var (
	UserResultsLimit  = 50
	VideoResultsLimit = 25
)

func GetResultsLimit() int {
	res, _ := strconv.Atoi(os.Getenv("RESULTS_LIMIT"))

	return res
}

// Config func to get env value
func Config(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	// Return the value of the variable
	return os.Getenv(key)
}
