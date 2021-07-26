package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Video struct {
	ID            uint64    `json:"id" gorm:"primary_key"`
	UID           uuid.UUID `json:"uid"`
	ShortID       string    `json:"short_id" gorm:"unique"`
	Title         string    `json:"title" gorm:"required"`
	UserID        uint64    `json:"user_id`
	Description   string    `json:"description" gorm:"type:string"`
	Thumbnail     string    `json:"thumbnail" gorm:"type:varchar(100)"`
	Resolutions   string    `json:"resolutions" gorm:"required"`
	MaxResolution string    `json:"max_resolution"`
	Private       bool      `json:"private"  gorm:"type:bool;default:false"`
	Unlisted      bool      `json:"unlisted" gorm:"type:bool;default:false"`
	PublishedAt   time.Time `json:"published_at"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

type VidRes struct {
	ID         uint64
	Resolution string
}

//getVideosByUserID(userID uint){
//
//}
