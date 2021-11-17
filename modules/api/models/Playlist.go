package models

import "github.com/google/uuid"

type UserPlaylist struct {
	UUID        uuid.UUID           `json:"id" gorm:"primary_key"`
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	UserID      uuid.UUID           `json:"user_id" gorm:"type:varchar(255);"`
	Queue       []UserPlaylistQueue `json:"queue" gorm:"references:ID;type:varchar(255);"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
}

type UserPlaylistQueue struct {
	ID        uuid.UUID `json:"uuid"`
	VideoID   uuid.UUID `json:"video_id"`
	Video     Video     `json:"video" gorm:"type:varchar(255);"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
