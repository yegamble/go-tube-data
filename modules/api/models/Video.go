package models

import (
	"errors"
	"github.com/dchest/uniuri"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/yegamble/go-tube-api/modules/api/config"
	"github.com/yegamble/go-tube-api/modules/api/handler"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type Video struct {
	ID              uint64            `json:"id" gorm:"primary_key"`
	UID             uuid.UUID         `json:"uid" gorm:"unique;required;type:varchar(255);"`
	Slug            *string           `json:"slug" gorm:"unique"`
	ShortID         *string           `json:"short_id" gorm:"unique;required"`
	Title           *string           `json:"title" gorm:"required;not null" validate:"min=1,max=255"`
	UserID          uint64            `json:"user_id" form:"user_id"`
	User            User              `gorm:"foreignKey:UserID;references:ID;OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Description     *string           `json:"description" gorm:"type:string"`
	Tags            []string          `json:"tags" gorm:"type:string"`
	Thumbnail       *string           `json:"thumbnail" gorm:"type:varchar(100)"`
	Duration        float64           `json:"duration" gorm:"type:float;default:0"`
	Ld144           *string           `json:"144p" gorm:"type:varchar(255)"`
	Ld240           *string           `json:"240p" gorm:"type:varchar(255)"`
	SD360           *string           `json:"360p" gorm:"type:varchar(255)"`
	SD480           *string           `json:"480p" gorm:"type:varchar(255)"`
	HD              *string           `json:"720p" gorm:"type:varchar(255)"`
	FHD             *string           `json:"1080p" gorm:"type:varchar(255)"`
	QHD             *string           `json:"1920p" gorm:"type:varchar(255)"`
	HD2K            *string           `json:"2048p" gorm:"column:hdk2;type:varchar(255)"`
	UHD             *string           `json:"3840p" gorm:"type:varchar(255)"`
	FUHD            *string           `json:"7680p" gorm:"type:varchar(255)"`
	IsConverted     bool              `json:"is_converted" form:"is_converted" gorm:"type:bool"`
	ConversionQueue []ConversionQueue `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;type:varchar(255);"`
	Views           []View            `gorm:"type:varchar(255);"`
	Permission      int               `json:"permission"  gorm:"type:int;default:0"`
	PublishedAt     time.Time         `json:"published_at" gorm:"autoCreateTime"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}

type View struct {
	ID        int64
	UserID    uuid.UUID `json:"user_id" form:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:UID;OnUpdate:CASCADE,OnDelete:SET NULL;"`
	VideoID   uuid.UUID `json:"video_id" form:"video_id"`
	Video     Video     `gorm:"foreignKey:VideoID;references:UID;OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time
}

type WatchLaterQueue struct {
	ID        uuid.UUID
	UserID    uint64
	User      User    `json:"user_id" form:"user_id" gorm:"foreignKey:UserID;references:ID"`
	Videos    *string `json:"videos,omitempty"`
	CreatedAt time.Time
}

var (
	queue    []ConversionQueue
	video    Video
	videos   []Video
	StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
)

func GetTrendingVideos(maxVideoResults int) (*[]Video, error) {

	offset := (page - 1) * config.GetResultsLimit()

	db.Model(&View{}).Offset(offset).Limit(config.VideoResultsLimit).Select("count(distinct(video_id))")

	return &videos, nil
}

func (video *Video) countVideoView(*Video, error) {

}

func GetVideoByID(id string) (*Video, error) {
	tx := db.Begin()
	err := tx.First(&video, "id = ?", id).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &video, nil
}

func GetVideoByUID(uid string) (*Video, error) {
	tx := db.Begin()
	err := tx.First(&video, "uid = ?", uid).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &video, nil
}

func SearchVideo(searchTerm string, limit int, page int) (*[]User, error) {

	if page < 0 {
		return nil, errors.New("page cannot be negative")
	}

	offset := (page - 1) * config.GetResultsLimit()

	db.Offset(offset).Limit(config.UserResultsLimit)
	db.Where("title LIKE ?", "%"+searchTerm+"%")
	db.Where("description LIKE ?", "%"+searchTerm+"%")
	err := db.Find(&video).Error
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func createVideo(video *Video, user *User, file *multipart.FileHeader) (uuid.UUID, error) {

	tx := db.Begin()

	dir := os.Getenv("VIDEO_DIR")

	filename, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	src, err := file.Open()
	if err != nil {
		return uuid.Nil, err
	}

	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)

	err = os.MkdirAll(dir+"tmp/", 0777)
	if err != nil {
		return [16]byte{}, err
	}
	if err != nil {
		return uuid.Nil, err
	}

	tempDst, err := os.Create(filepath.Join(dir+"tmp/",
		filepath.Base(filename.String())))
	if err != nil {
		return uuid.Nil, err
	}

	defer func(tempDst *os.File) {
		_ = tempDst.Close()
	}(tempDst)

	if _, err = io.Copy(tempDst, src); err != nil {
		return uuid.Nil, err
	}

	video.UserID = user.ID
	shortID := uniuri.NewLenChars(10, StdChars)
	video.ShortID = &shortID
	video.UID = filename

	//createVideoQueue
	video.createConversionQueue(tempDst.Name())
	video.Duration, err = getVideoDuration(tempDst.Name())

	err = db.Create(&video).Error
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	return video.UID, nil
}

func getVideoWidth(videoDirectory string) (int, error) {
	a, err := ffmpeg.Probe(videoDirectory)
	if err != nil {
		return 0, err
	}

	streamsArray := gjson.Get(a, "streams").Array()
	for _, stream := range streamsArray {
		if stream.Get("width").Int() > 0 {
			return int(stream.Get("width").Int()), nil
		}
	}

	return 0, errors.New("video streams not found")
}

func getVideoDuration(videoDirectory string) (float64, error) {
	a, err := ffmpeg.Probe(videoDirectory)
	if err != nil {
		return 0, err
	}

	totalDuration := gjson.Get(a, "format.duration").Float()

	if totalDuration == 0 {
		return 0, errors.New("video streams not found")
	}

	return totalDuration, nil
}

func getVideoByUID(uid uuid.UUID) (Video, error) {
	tx := db.Begin()
	video := Video{}
	err := tx.Where("uid = ?", uid).First(&video).Error
	if err != nil {
		return video, err
	}

	tx.Commit()
	return video, nil
}

func (video *Video) setVideoAsConverted() error {
	tx := db.Begin()

	db.Model(&video).Update("is_converted", true)
	tx.Commit()
	return nil
}

func EditVideo() error {
	return nil
}

func (video *Video) DeleteVideo() error {
	tx := db.Begin()

	for _, resolutionColumn := range resolutionKey {

		tx2 := db.Begin()
		tempStruct := struct{ a string }{}
		tx2.Select(resolutionColumn).First(&video, "uid = ?", video.UID).Scan(tempStruct)

		videoFile, err := os.Open(tempStruct.a)
		if os.IsNotExist(err) {
			tx2.Model(&video).Update(resolutionColumn, nil)
		} else {
			err = os.Remove(videoFile.Name())
			if err != nil {
				tx2.Rollback()
				return err
			}
			tx2.Model(&video).Update(resolutionColumn, nil)
		}
	}

	err := tx.Delete(&video).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func GetUserVideos() error {
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
