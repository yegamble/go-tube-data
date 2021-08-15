package config

import (
	"os"
	"strconv"
)

type Config struct {
	ID    uint64 `json:"id" gorm:"primary_key"`
	Name  string `json:"name" gorm:"type:text"`
	Value string `json:"value" gorm:"type:text"`
}

var (
	UserResultsLimit  = 50
	VideoResultsLimit = 25
)

func GetResultsLimit() int {
	res, _ := strconv.Atoi(os.Getenv("RESULTS_LIMIT"))

	return res
}
