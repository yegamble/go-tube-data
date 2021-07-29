package models

import (
	"github.com/gofiber/fiber/v2"
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
	Permission    int       `json:"permission"  gorm:"type:int;default:0"`
	PublishedAt   time.Time `json:"published_at" gorm:"autoCreateTime"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

type VidRes struct {
	ID         uint64
	Resolution string
}

func createVideo(c *fiber.Ctx) error {
	return nil
}

func convertVideo(video *Video) error {
	return nil
}

func EditVideo() error {
	return nil
}

func DeleteVideo() error {
	return nil
}

func GetUserVideos() error {
	return nil
}

func GetAllVideosPublic() error {
	return nil
}
