package video

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Video struct {
	gorm.Model
	ID            uuid.UUID `json:"id" gorm:"primary_key"`
	ShortID       string    `json:"short_id" gorm:"unique"`
	Title         string    `json:"title" gorm:"required"`
	UserID        uuid.UUID `json:"user_id`
	Description   string    `json:"description" gorm:"type:string"`
	Thumbnail     string    `json:"thumbnail""`
	Resolutions   string    `json:"resolutions" gorm:"required"`
	MaxResolution string    `json:"max_resolution"`
	PublishedAt   time.Time `json:"published_at"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

type VidRes struct {
	ID         int64
	Resolution string
}
