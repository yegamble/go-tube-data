package models

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"math/rand"
	"net"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ConversionQueue struct {
	ID         uint64    `json:"id" gorm:"primary_key"`
	UserUID    uuid.UUID `json:"user_uid" form:"user_uid" gorm:"type:varchar(255);"`
	VideoUID   uuid.UUID `json:"video_id" form:"video_id"`
	Resolution *string   `json:"resolution" gorm:"type:varchar(100)"`
	TempFile   string    `json:"temp_file" gorm:"type:varchar(255)"`
	Status     string    `json:"status" gorm:"type:varchar(100)"`
	CreatedAt  time.Time
}

var (
	acceptedMimes = map[string]string{
		"video/mp4":       "mp4",
		"video/quicktime": "mov",
		"video/x-ms-wmv":  "wmv",
	}

	resolutionKey = map[int]string{
		144:  "Ld144",
		240:  "Ld240",
		360:  "SD360",
		480:  "SD480",
		720:  "HD",
		1080: "FHD",
		1920: "QHD",
		2048: "HD2K",
		3840: "UHD",
		7680: "FUHD",
	}

	resolutions = map[string]ffmpeg.KwArgs{
		"FUHD":  {"filter:v": "scale=7680:-2", "b:v": "80M", "b:a": "512k", "crf": "18"},
		"UHD":   {"filter:v": "scale=3840:-2", "b:v": "50M", "b:a": "512k", "crf": "18"},
		"HD2K":  {"filter:v": "scale=2048:-2", "b:v": "20M", "b:a": "512k", "crf": "18"},
		"QHD":   {"filter:v": "scale=1920:-2", "b:v": "10M", "b:a": "512k", "crf": "18"},
		"FHD":   {"filter:v": "scale=1920:-2", "b:v": "10M", "b:a": "512k", "crf": "18"},
		"HD":    {"filter:v": "scale=1280:-2", "b:v": "5M", "b:a": "512k", "crf": "18"},
		"SD480": {"filter:v": "scale=854:-2", "b:v": "2.5M", "b:a": "384k", "crf": "23"},
		"SD360": {"filter:v": "scale=640:-2", "b:v": "1.5M", "b:a": "384k", "crf": "25"},
		"Ld240": {"filter:v": "scale=426:-2", "b:v": "0.5M", "b:a": "128k", "crf": "28"},
		"Ld144": {"filter:v": "scale=256:-2", "b:v": "0.25M", "b:a": "128k", "crf": "28"},
	}
	baseArgs = ffmpeg.KwArgs{"c:v": "libx264",
		"preset":    "slow",
		"profile:v": "high",
		"coder":     "1",
		"pix_fmt":   "yuv420p",
		"movflags":  "+faststart",
		"g":         "30",
		"bf":        "2",
		"c:a":       "aac",
		"profile:a": "aac_low"}
)

func (video *Video) createConversionQueue(temporaryVideoDirectory string) error {

	tx := db.Begin()

	videoWidth, err := getVideoWidth(temporaryVideoDirectory)
	if err != nil {
		return err
	}

	for key, resolution := range resolutionKey {
		fmt.Println(key)
		if videoWidth >= key {
			queue := ConversionQueue{
				UserUID:    video.UserID,
				VideoUID:   video.UID,
				Resolution: &resolution,
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
		err := q.convertVideo()
		if err != nil {
			return err
		}
	}

	err = video.setVideoAsConverted()
	if err != nil {
		return err
	}

	return nil
}

func (queue ConversionQueue) convertVideo() error {
	tx := db.Begin()
	video = Video{}

	tx.Model(&queue).Update("status", "processing")
	tmpFile := os.Getenv("VIDEO_DIR_TMP") + queue.VideoUID.String()

	a, err := ffmpeg.Probe(tmpFile)
	if err != nil {
		tx.Model(&queue).Update("status", "error")
		tx.Commit()
		return errors.New("error with video file")
	}

	tx.Commit()

	tx2 := db.Begin()

	totalDuration := gjson.Get(a, "format.duration").Float()

	input := ffmpeg.Input(tmpFile, nil)
	filename := uuid.NewString()
	filelocation := os.Getenv("VIDEO_DIR") + uuid.NewString() + os.Getenv("APP_VIDEO_EXTENSION")

	if resolutions[*queue.Resolution] != nil {
		//err = tx2.Model(&video).Where("uid = ?", queue.VideoUID).Update(queue.Resolution, filelocation).Error
		err = queue.createVideoFile(&filelocation, filename, tx2)
		if err != nil {
			return err
		}

		err = input.Output(filelocation, baseArgs, resolutions[*queue.Resolution]).
			GlobalArgs("-progress", "unix://"+conversionProgressSock(totalDuration)).
			OverWriteOutput().
			Run()
		if err != nil {
			tx2.Model(&queue).Update("status", "error")
			tx2.Commit()
			return err
		}
	}

	tx2.Model(&queue).Update("status", "complete")

	err = tx2.Commit().Error
	if err != nil {
		tx2.Rollback()
		return err
	}

	return err
}

func conversionProgressSock(totalDuration float64) string {
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
