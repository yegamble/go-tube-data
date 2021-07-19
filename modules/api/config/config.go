package config

import "github.com/google/uuid"

type Config struct {
	ID    uuid.UUID `json:"id" gorm:"primary_key"`
	Name  string    `json:"name" gorm:"type:text"`
	Value string    `json:"value" gorm:"type:text"`
}
