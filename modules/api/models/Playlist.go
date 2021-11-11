package models

import "github.com/google/uuid"

type UserPlaylist struct {
	UUID        uuid.UUID           `json:"id" gorm:"primary_key"`
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	UserUUID    uuid.UUID           `json:"user_uuid" gorm:"type:varchar(255);"`
	Queue       []UserPlaylistQueue `json:"queue" gorm:"foreignKey:UUID;references:UUID;type:varchar(255);"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
}

type UserPlaylistQueue struct {
	UUID      uuid.UUID `json:"uuid"`
	VideoUUID uuid.UUID `json:"video_uuid"`
	Video     Video     `json:"video" gorm:"foreignKey:VideoUUID;references:uuid;type:varchar(255);"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
