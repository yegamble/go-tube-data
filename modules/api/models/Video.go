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
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type Video struct {
	ID              uuid.UUID          `json:"uid" gorm:"primary_key;unique;required;type:varchar(255);"`
	Slug            *string            `json:"slug" gorm:"unique"`
	ShortID         *string            `json:"short_id" gorm:"unique;required"`
	Title           *string            `json:"title" gorm:"required;not null" validate:"min=1,max=255"`
	UserID          uuid.UUID          `json:"user_id" form:"user_id" gorm:"type:varchar(255);"`
	Description     *string            `json:"description" gorm:"type:string"`
	Tags            []*Tag             `json:"tags,omitempty" gorm:"many2many:video_tags;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Thumbnail       *string            `json:"thumbnail" gorm:"type:varchar(100)"`
	Duration        float64            `json:"duration" gorm:"type:float;default:0"`
	IsConverted     bool               `json:"is_converted" form:"is_converted" gorm:"type:bool"`
	IsPrivate       bool               `json:"is_private" form:"is_private" gorm:"type:bool"`
	IsUnlisted      bool               `json:"is_unlisted" form:"is_unlisted" gorm:"type:bool"`
	CategoryID      uuid.UUID          `json:"category_id" form:"category_id" gorm:"type:varchar(255);constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Category        Category           `gorm:"type:varchar(255)"`
	ConversionQueue []*ConversionQueue `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;type:varchar(255);"`
	Views           []*User            `gorm:"many2many:user_views;References:ID;joinReferences:UserID"`
	Likes           []*Like            `json:"likes"`
	Permission      int                `json:"permission"  gorm:"type:int;default:0"`
	PublishedAt     time.Time          `json:"published_at" gorm:"autoCreateTime"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}

type Like struct {
	gorm.Model
	ID      uuid.UUID `json:"id" gorm:"primary_key"`
	VideoID uuid.UUID `json:"video_id" gorm:"type:varchar(255);"`
	UserID  uuid.UUID `json:"user_id" gorm:"type:varchar(255);"`
	User    User
	Rating  bool `json:"rating" gorm:"type:int"`
}

type VideoFile struct {
	ID         uint64    `json:"id" gorm:"primary_key"`
	VideoID    uuid.UUID `json:"video_id"`
	Resolution *string   `json:"resolution" gorm:"type:varchar(255);"`
	FileName   *string   `json:"file_name" gorm:"type:varchar(255);"`
	FileSize   *int64    `json:"file_size"`
	FileType   *string   `json:"file_type" gorm:"type:varchar(255);"`
	FilePath   *string   `json:"file_path" gorm:"type:varchar(255);"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

type Category struct {
	ID          uuid.UUID `json:"id" gorm:"unique;required;type:varchar(255);primary_key"`
	Name        *string   `json:"name" gorm:"required;not null" validate:"min=1,max=255"`
	Description *string   `json:"description" gorm:"type:string"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type UserView struct {
	ID        uuid.UUID
	UserID    uuid.UUID `json:"user_id" form:"user_id"`
	VideoID   uuid.UUID `json:"video_id" form:"video_id"`
	CreatedAt time.Time
}

type WatchLaterQueue struct {
	ID        uuid.UUID
	UserID    uuid.UUID `json:"user_id" form:"user_id" gorm:"type:varchar(255);"`
	Videos    *string   `json:"videos,omitempty"`
	CreatedAt time.Time
}

var (
	queue    []ConversionQueue
	video    Video
	videos   []Video
	StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
)

//func GetTrendingVideos(maxVideoResults int) (*[]Video, error) {
//
//	offset := (page - 1) * config.GetResultsLimit()
//
//	db.Model(&View{}).Offset(offset).Limit(config.VideoResultsLimit).Select("count(distinct(video_id))")
//
//	return &videos, nil
//}

func (video *Video) CreateUserView(user *User) error {
	tx := db.Begin()
	video.Views = append(video.Views, user)
	err := tx.Save(&user).Error
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (queue *ConversionQueue) createVideoFile(destination *string, fileName string, tx *gorm.DB) error {
	videoFile := VideoFile{}
	videoFile.VideoID = queue.VideoID
	videoFile.Resolution = queue.Resolution
	videoFile.FileName = &fileName
	videoFile.FilePath = destination

	err := tx.Create(&videoFile).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
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
	db.Where("private = ?", false)
	db.Where("unlisted = ?", false)
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
	video.ID = filename

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

	return video.ID, nil
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
	tx.Model(&video).Where("uid = ?", video.ID)

	for _, resolutionColumn := range resolutionKey {

		var result = map[string]interface{}{}
		result[resolutionColumn] = ""

		tx.Model(&video).First(&result, "uid = ?", video.ID)

		log.Println(result[resolutionColumn])
		videoFile, err := os.Open("")
		if os.IsNotExist(err) {
			log.Println(err.Error())
			tx.Model(&video).Where("uid = ?", video.ID).Update(resolutionColumn, nil)
		} else {
			log.Println(videoFile.Name())
			err = os.Remove(videoFile.Name())
			if err != nil {
				tx.Rollback()
				return err
			}
			tx.Model(&video).Update(resolutionColumn, nil)
		}
	}

	err := tx.Where("uid = ?", video.ID).Delete(&video).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (user *User) GetUserVideos() (err error) {
	return err
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
