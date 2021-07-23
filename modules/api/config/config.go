package config

import (
	"github.com/google/uuid"
	"os"
	"strconv"
)

type Config struct {
	ID    uuid.UUID `json:"id" gorm:"primary_key"`
	Name  string    `json:"name" gorm:"type:text"`
	Value string    `json:"value" gorm:"type:text"`
}

var (
	UserResultsLimit  = 50
	videoResultsLimit = 25
)

func GetResultsLimit() int {
	res, _ := strconv.Atoi(os.Getenv("RESULTS_LIMIT"))

	return res
}
