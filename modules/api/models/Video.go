package models

import (
	"errors"
	"fmt"
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
	"math/rand"
	"mime/multipart"
	"net"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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
	ld144           *string           `json:"144p" gorm:"type:varchar(255)"`
	ld240           *string           `json:"240p" gorm:"type:varchar(255)"`
	SD360           *string           `json:"360p" gorm:"type:varchar(255)"`
	SD480           *string           `json:"480p" gorm:"type:varchar(255)"`
	HD              *string           `json:"720p" gorm:"type:varchar(255)"`
	FHD             *string           `json:"1080p" gorm:"type:varchar(255)"`
	QHD             *string           `json:"1920p" gorm:"type:varchar(255)"`
	HD2K            *string           `json:"2048p" gorm:"type:varchar(255)"`
	UHD             *string           `json:"3840p" gorm:"type:varchar(255)"`
	FUHD            *string           `json:"7680p" gorm:"type:varchar(255)"`
	IsConverted     bool              `json:"is_converted" form:"is_converted" gorm:"type:bool"`
	ConversionQueue []ConversionQueue `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;type:varchar(255);"`
	Permission      int               `json:"permission"  gorm:"type:int;default:0"`
	PublishedAt     time.Time         `json:"published_at" gorm:"autoCreateTime"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}

type VidRes struct {
	ID         uint64
	Resolution string
}

type View struct {
	ID        int64
	UserID    uint64 `json:"user_id" form:"user_id"`
	User      User   `gorm:"foreignKey:UserID;references:ID;OnUpdate:CASCADE,OnDelete:SET NULL;"`
	VideoID   uint64 `json:"video_id" form:"video_id"`
	Video     Video  `gorm:"foreignKey:VideoID;references:ID;OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time
}

type ConversionQueue struct {
	ID         uint64    `json:"id" gorm:"primary_key"`
	VideoID    uuid.UUID `json:"video_id" form:"video_id"`
	Resolution string    `json:"resolution" gorm:"type:varchar(100)"`
	TempFile   string    `json:"temp_file" gorm:"type:varchar(255)"`
	Status     string    `json:"status" gorm:"type:varchar(100)"`
	CreatedAt  time.Time
}

type WatchLaterQueue struct {
	ID        uuid.UUID
	UserID    uint64
	User      User    `json:"user_id" form:"user_id" gorm:"foreignKey:UserID;references:ID"`
	Videos    *string `json:"videos,omitempty"`
	CreatedAt time.Time
}

var (
	queue         []ConversionQueue
	video         Video
	videos        []Video
	StdChars      = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
	acceptedMimes = map[string]string{
		"video/mp4":       "mp4",
		"video/quicktime": "mov",
		"video/x-ms-wmv":  "wmv",
	}
	resolutions = map[int64]ffmpeg.KwArgs{
		7680: {"filter:v": "scale=7680:-2", "b:v": "80M"},
		3840: {"filter:v": "scale=3840:-2", "b:v": "50M"},
		2048: {"filter:v": "scale=2048:-2", "b:v": "20M"},
		1920: {"filter:v": "scale=1920:-2", "b:v": "10M"},
		1080: {"filter:v": "scale=1920:-2", "b:v": "10M"},
		720:  {"filter:v": "scale=1280:-2", "b:v": "5M"},
		480:  {"filter:v": "scale=854:-2", "b:v": "2.5M"},
		360:  {"filter:v": "scale=640:-2", "b:v": "1M"},
		240:  {"filter:v": "scale=240:-2", "b:v": "0.5M"},
		144:  {"filter:v": "scale=144:-2", "b:v": "0.25M"},
	}
	baseArgs = ffmpeg.KwArgs{"c:v": "libx264",
		"preset":    "slow",
		"profile:v": "high",
		"crf":       "18",
		"coder":     "1",
		"pix_fmt":   "yuv420p",
		"movflags":  "+faststart",
		"g":         "30",
		"bf":        "2",
		"c:a":       "aac",
		"b:a":       "384k",
		"profile:a": "aac_low"}
)

func GetTrendingVideos(maxVideoResults int) (*[]Video, error) {

	offset := (page - 1) * config.GetResultsLimit()

	db.Model(&View{}).Offset(offset).Limit(config.VideoResultsLimit).Select("count(distinct(video_id))")

	return &videos, nil
}

func countVideoView(*Video, error) {

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

	defer src.Close()

	os.MkdirAll(dir+"tmp/", 0777)
	if err != nil {
		return uuid.Nil, err
	}

	tempDst, err := os.Create(filepath.Join(dir+"tmp/",
		filepath.Base(filename.String())))
	if err != nil {
		return uuid.Nil, err
	}

	defer tempDst.Close()

	if _, err = io.Copy(tempDst, src); err != nil {
		return uuid.Nil, err
	}

	video.UserID = user.ID
	shortID := uniuri.NewLenChars(10, StdChars)
	video.ShortID = &shortID
	video.UID = filename

	//createVideoQueue
	video.createConversionQueue(tempDst.Name())

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

func getVideoWidth(videoDirectory string) (int64, error) {
	a, err := ffmpeg.Probe(videoDirectory)
	if err != nil {
		return 0, err
	}

	streamsArray := gjson.Get(a, "streams").Array()
	for _, stream := range streamsArray {
		if stream.Get("width").Int() > 0 {
			return stream.Get("width").Int(), nil
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

func (video *Video) createConversionQueue(temporaryVideoDirectory string) error {

	tx := db.Begin()

	videoWidth, err := getVideoWidth(temporaryVideoDirectory)
	if err != nil {
		return err
	}

	for resolution, _ := range resolutions {
		if resolution <= videoWidth {
			queue := ConversionQueue{
				VideoID:    video.UID,
				Resolution: strconv.FormatInt(resolution, 10),
				Status:     "pending",
			}

			err := tx.Create(&queue).Error
			if err != nil {
				tx.Rollback()
				return err
			}

		}
	}

	tx.Commit()
	return nil
}

func getConversionQueue(videoUID uuid.UUID) ([]ConversionQueue, error) {

	tx := db.Begin()

	if videoUID != uuid.Nil {
		tx.Where("video_id = ?", videoUID).Find(&queue)
	} else {
		tx.Find(&queue, "status = ?", "pending")
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return queue, nil
}

func convertQueueByVideo(uid uuid.UUID) error {
	if uid == uuid.Nil {
		return errors.New("video uid cannot be nil")
	}

	video, err := getVideoByUID(uid)
	if err != nil {
		return err
	}

	queue, err := getConversionQueue(video.UID)
	if err != nil {
		return err
	}

	for _, q := range queue {
		err := q.ConvertVideoFromQueue(db)
		if err != nil {
			return err
		}
	}

	return nil
}

func getVideoByUID(uid uuid.UUID) (Video, error) {
	tx := db.Begin()
	video := Video{}
	err := tx.Where("uid = ?", uid).First(&video).Error
	if err != nil {
		tx.Rollback()
		return video, err
	}

	tx.Commit()
	return video, nil
}

func (queue *ConversionQueue) ConvertVideoFromQueue(tx *gorm.DB) error {
	tx.Begin()
	tx.Model(&queue).Statement.SetColumn("status", "processing")
	tmpFile := os.Getenv("VIDEO_DIR_TMP") + queue.VideoID.String()
	fmt.Println(tmpFile)

	a, err := ffmpeg.Probe(tmpFile)
	if err != nil {
		log.Println(err.Error())
		return errors.New("error with temporary video file")
	}

	totalDuration := gjson.Get(a, "format.duration").Float()

	input := ffmpeg.Input(tmpFile, nil)

	vidResolution, err := strconv.ParseInt(queue.Resolution, 10, 64)
	if err != nil {
		return err
	}

	filename := os.Getenv("VIDEO_DIR") + queue.VideoID.String() + os.Getenv("APP_VIDEO_EXTENSION")

	if resolutions[vidResolution] != nil {

		err = input.Output(filename, baseArgs, resolutions[vidResolution]).
			GlobalArgs("-progress", "unix://"+TempSock(totalDuration)).
			OverWriteOutput().
			Run()
		if err != nil {
			return err
		}
	}

	return err
}

func saveConvertedVideo(queue *ConversionQueue) error {
	tx := db.Begin()

	db.Model(&User{}).Where("active = ?", true).Update(queue.Resolution, "hello")

	tx.Model(&queue).Statement.SetColumn("status", "done")
	tx.Commit()
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

func TempSock(totalDuration float64) string {
	// serve

	rand.Seed(time.Now().Unix())
	sockFileName := path.Join(os.TempDir(), fmt.Sprintf("%d_sock", rand.Int()))
	l, err := net.Listen("unix", sockFileName)
	if err != nil {
		panic(err)
	}

	go func() {
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		buf := make([]byte, 16)
		data := ""
		progress := ""
		for {
			_, err := fd.Read(buf)
			if err != nil {
				return
			}
			data += string(buf)
			a := re.FindAllStringSubmatch(data, -1)
			cp := ""
			if len(a) > 0 && len(a[len(a)-1]) > 0 {
				c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
				cp = fmt.Sprintf("%.2f", float64(c)/totalDuration/1000000)
			}
			if strings.Contains(data, "progress=end") {
				cp = "done"
			}
			if cp == "" {
				cp = ".0"
			}
			if cp != progress {
				progress = cp
				fmt.Println("progress: ", progress)
			}
		}
	}()

	return sockFileName
}
