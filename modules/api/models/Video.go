package models

import (
	"errors"
	"github.com/dchest/uniuri"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
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
	IsConverted   bool      `json:"is_converted" form:"is_converted" gorm:"type:bool"`
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

type ConversionQueue struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	VideoID   uuid.UUID `json:"video_id" form:"video_id"`
	Video     Video     `gorm:"foreignKey:VideoID;references:ID;not null"`
	CreatedAt time.Time
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

	video.ShortID = uniuri.NewLenChars(10, StdChars)

	createVideo(file)

	return nil
}

func createVideo(file *multipart.FileHeader) error {
	dir := "uploads/video/" + user.Username + "/"

	filename, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	tempDst, err := os.Create(filepath.Join(dir+"temp/", filepath.Base(strings.Replace(filename.String()+os.Getenv("APP_VIDEO_EXTENSION"), "-", "_", -1))))
	if err != nil {
		return err
	}

	defer tempDst.Close()

	if _, err = io.Copy(tempDst, src); err != nil {
		return err
	}

	err = convertVideo(tempDst.Name())
	if err != nil {
		return err
	}

	return nil
}

func convertVideo(videoDir string) error {

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
