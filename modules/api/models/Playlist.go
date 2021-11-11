package models

import "github.com/google/uuid"

type UserPlaylist struct {
	Id          uuid.UUID `json:"id" gorm:"primary_key"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	UserId      uuid.UUID `json:"user_uuid"`
	Videos      []UserPlaylistQueue
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UserPlaylistQueue struct {
	Id        uuid.UUID `json:"id" gorm:"primary_key"`
	UserId    uuid.UUID `json:"user_id" gorm:"foreignKey:UserUUID;references:UUID;OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User      User      `json:"user"`
	VideoUID  uuid.UUID `json:"video_uid"`
	Videos    Video
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
