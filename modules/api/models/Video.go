package models

import (
	"errors"
	"github.com/dchest/uniuri"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/modules/api/handler"
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
	UID           uuid.UUID `json:"uid" gorm:"unique;required"`
	ShortID       string    `json:"short_id" gorm:"unique;required"`
	Title         string    `json:"title" gorm:"required;not null" validate:"min=1,max=255"`
	UserID        uint64    `json:"user_id" form:"user_id"`
	User          User      `gorm:"foreignKey:UserID;references:ID"`
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
	UserID    uint64    `json:"user_id" form:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
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

	var body Video

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	file, err := c.FormFile("video")
	if err != nil {
		return err
	}

	user, _ = GetUserByID(c.FormValue("user_id"))
	if err != nil {
		return err
	}

	contentType := file.Header.Get("content-type")
	_, exists := acceptedMimes[contentType]
	if !exists {
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(errors.New("unsupported video format").Error())
	}

	video.UserID = user.ID

	createVideo(&body, user, file)

	return nil
}

func createVideo(video *Video, user User, file *multipart.FileHeader) error {

	dir := "uploads/videos/" + user.Username + "/"

	filename, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	os.MkdirAll(dir+"temp/", 0777)
	if err != nil {
		return err
	}

	tempDst, err := os.Create(filepath.Join(dir+"temp/", filepath.Base(strings.Replace(filename.String()+os.Getenv("APP_VIDEO_EXTENSION"), "-", "_", -1))))
	if err != nil {
		return err
	}

	defer tempDst.Close()

	if _, err = io.Copy(tempDst, src); err != nil {
		return err
	}

	video.UserID = user.ID
	video.ShortID = uniuri.NewLenChars(10, StdChars)
	video.UID = filename
	err = db.Create(&video).Error
	if err != nil {
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

func ValidateVideoStruct(video *Video) []*handler.ErrorResponse {
	var errors []*handler.ErrorResponse
	var element handler.ErrorResponse
	validate := validator.New()

	err := validate.Struct(video)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}
