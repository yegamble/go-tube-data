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
	ID            uint64    `json:"id" gorm:"primary_key"`
	UID           uuid.UUID `json:"uid" gorm:"unique;required"`
	Slug          string    `json:"slug" gorm:"unique"`
	ShortID       string    `json:"short_id" gorm:"unique;required"`
	Title         string    `json:"title" gorm:"required;not null" validate:"min=1,max=255"`
	UserID        uint64    `json:"user_id" form:"user_id"`
	User          User      `gorm:"foreignKey:UserID;references:ID"`
	Description   string    `json:"description" gorm:"type:string"`
	Tags          []string  `json:"tags" gorm:"type:string"`
	Thumbnail     string    `json:"thumbnail" gorm:"type:varchar(100)"`
	Duration      uint64    `json:"duration" gorm:"type:int;default:0"`
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

type View struct {
	ID        int64
	UserID    uint64 `json:"user_id" form:"user_id"`
	User      User   `gorm:"foreignKey:UserID;references:ID"`
	VideoID   uint64 `json:"video_id" form:"video_id"`
	Video     Video  `gorm:"foreignKey:VideoID;references:ID"`
	CreatedAt time.Time
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
	videos        []Video
	StdChars      = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
	acceptedMimes = map[string]string{
		"video/mp4":       "mp4",
		"video/quicktime": "mov",
		"video/x-ms-wmv":  "wmv",
	}
	scale360Args  = ffmpeg.KwArgs{"filter:v": "scale=640:-2", "b:v": "1M"}
	scale480Args  = ffmpeg.KwArgs{"filter:v": "scale=854:-2", "b:v": "2.5M"}
	scale720Args  = ffmpeg.KwArgs{"filter:v": "scale=1280:-2", "b:v": "5M"}
	scale1080Args = ffmpeg.KwArgs{"filter:v": "scale=1920:-2", "b:v": "10M"}
	scale2kArgs   = ffmpeg.KwArgs{"filter:v": "scale=2048:-2", "b:v": "20M"}
	scale4kArgs   = ffmpeg.KwArgs{"filter:v": "scale=3840:-2", "b:v": "50M"}
	scale8kArgs   = ffmpeg.KwArgs{"filter:v": "scale=7680:-2", "b:v": "80M"}
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

func createVideo(video *Video, user *User, file *multipart.FileHeader) error {

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

	tempDst, err := os.Create(filepath.Join(dir+"temp/",
		filepath.Base(strings.Replace(filename.String()+os.Getenv("APP_VIDEO_EXTENSION"),
			"-", "_", -1))))
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

	err = convertVideo(tempDst.Name(), dir, filename.String())
	if err != nil {
		return err
	}

	return nil
}

func convertVideo(videoDir string, dstDir string, filename string) error {

	baseArgs := ffmpeg.KwArgs{"c:v": "libx264",
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

	a, err := ffmpeg.Probe(videoDir)
	if err != nil {
		return err
	}

	totalDuration := gjson.Get(a, "format.duration").Float()

	input := ffmpeg.Input(videoDir, nil)

	vidWidth := gjson.Get(a, "streams.0.width").Int()

	if vidWidth > 0 {
		err = input.Output(dstDir+filename+"_360p"+os.Getenv("APP_VIDEO_EXTENSION"), baseArgs, scale360Args).
			GlobalArgs("-progress", "unix://"+TempSock(totalDuration)).
			OverWriteOutput().
			Run()
		if err != nil {
			return err
		}
	}

	if vidWidth >= 854 {
		err = input.Output(dstDir+filename+"_480p"+os.Getenv("APP_VIDEO_EXTENSION"), baseArgs, scale480Args).
			GlobalArgs("-progress", "unix://"+TempSock(totalDuration)).
			OverWriteOutput().
			Run()
		if err != nil {
			return err
		}
	}

	if vidWidth >= 1280 {
		err = input.Output(dstDir+filename+"_720p"+os.Getenv("APP_VIDEO_EXTENSION"), baseArgs, scale720Args).
			GlobalArgs("-progress", "unix://"+TempSock(totalDuration)).
			OverWriteOutput().
			Run()
		if err != nil {
			return err
		}
	}

	if vidWidth >= 1920 {
		err = input.Output(dstDir+filename+"_1080p"+os.Getenv("APP_VIDEO_EXTENSION"), baseArgs, scale1080Args).
			GlobalArgs("-progress", "unix://"+TempSock(totalDuration)).
			OverWriteOutput().
			Run()
		if err != nil {
			return err
		}
	}

	if vidWidth >= 2048 {
		err = input.Output(dstDir+filename+"_480p"+os.Getenv("APP_VIDEO_EXTENSION"), baseArgs, scale2kArgs).
			GlobalArgs("-progress", "unix://"+TempSock(totalDuration)).
			OverWriteOutput().
			Run()
		if err != nil {
			return err
		}
	}

	if vidWidth >= 3840 {
		err = input.Output(dstDir+filename+"_480p"+os.Getenv("APP_VIDEO_EXTENSION"), baseArgs, scale4kArgs).
			GlobalArgs("-progress", "unix://"+TempSock(totalDuration)).
			OverWriteOutput().
			Run()
		if err != nil {
			return err
		}
	}

	if vidWidth >= 7680 {
		err = input.Output(dstDir+filename+"_480p"+os.Getenv("APP_VIDEO_EXTENSION"), baseArgs, scale8kArgs).
			GlobalArgs("-progress", "unix://"+TempSock(totalDuration)).
			OverWriteOutput().
			Run()
		if err != nil {
			return err
		}
	}

	return err
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
