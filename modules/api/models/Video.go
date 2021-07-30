package models

import (
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
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

	//contentType := file.Header.Get("content-type")
	//_, exists := acceptedMimes[contentType]
	//if !exists {
	//	return c.Status(fiber.StatusUnsupportedMediaType).JSON(errors.New("unsupported video format").Error())
	//}

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

	err = convertVideo(tempDst.Name(), dir+filename.String())
	if err != nil {
		return err
	}

	return nil
}

func convertVideo(videoDir string, dstDir string) error {

	//command := []string{"expr:gte(t,n_forced/2)"}
	a, err := ffmpeg.Probe(videoDir)
	if err != nil {
		panic(err)
	}
	totalDuration := gjson.Get(a, "format.duration").Float()

	err = ffmpeg.Input(videoDir, nil).Output(dstDir+os.Getenv("APP_VIDEO_EXTENSION"), ffmpeg.KwArgs{"c:v": "libx264", "b:v": "15M", "preset": "slow", "profile:v": "high", "crf": "18", "coder": "1", "pix_fmt": "yuv420p", "movflags": "+faststart", "g": "30", "bf": "2", "c:a": "aac", "b:a": "384k", "profile:a": "aac_low"}).GlobalArgs("-progress", "unix://"+TempSock(totalDuration)).OverWriteOutput().Run()
	//Output(dstDir+os.Getenv("APP_VIDEO_EXTENSION"), ffmpeg.KwArgs{"vf": "yadif,format=yuv422p","force_key_frames":strings.Join(command, "', '"),"c:v": "libx264","b:v":"15","bf":"2","c:a":"aac","crf":"18","ac":"2","ar":"44100","use_editlist":"0","movflags":"+faststart"}).OverWriteOutput().Run()
	//"-vf": "yadif, format=yuv422p","force_key_frames": "expr:gte(t\\,n_forced/2)", "c:v":"libx264", "b:v": "<60M for 1080, 50M for 720, 15M for SD>", "bf": "2", "c:a": "flac", "ac": "2", "ar": "44100", "use_editlist": "0", "movflags": "+faststart"

	if err != nil {
		return err
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
