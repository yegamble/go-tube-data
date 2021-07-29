package models

import (
	"errors"
	"github.com/dchest/uniuri"
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

var (
	video         Video
	StdChars      = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
	acceptedMimes = map[string]string{
		"video/mp4": "mp4",
	}
)

func UploadVideo(c *fiber.Ctx) error {

	file, err := c.FormFile("video")
	if err != nil {
		return err
	}

	contentType := file.Header.Get("content-type")
	_, exists := acceptedMimes[contentType]
	if !exists {
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(errors.New("unsupported video format").Error())
	}

	return nil
}

func createVideo(c *fiber.Ctx) error {

	video.ShortID = uniuri.NewLenChars(10, StdChars)
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
