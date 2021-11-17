package models

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/modules/api/config"
)

func GetAllVideos(c *fiber.Ctx) error {

	offset := (page - 1) * config.GetResultsLimit()

	db.Offset(offset).Limit(config.VideoResultsLimit).Find(&videos)
	if len(videos) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "false",
			"message": "Videos Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(videos)
}

func FetchVideoByID(c *fiber.Ctx) error {
	video, err := GetVideoByID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(video)
}

func FetchVideoByUID(c *fiber.Ctx) error {
	video, err := GetVideoByUID(c.Params("uid"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(video)
}

func TriggerConversionByVideoUID(c *fiber.Ctx) error {
	videoUID, err := uuid.Parse(c.FormValue("uid"))

	err = convertQueueByVideo(videoUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "true",
		"message": "Video Conversion Complete",
	})
}

func TriggerConversionByQueueID(c *fiber.Ctx) error {
	conversionQueueID, err := uuid.Parse(c.Params("id"))

	err = convertQueueByVideo(conversionQueueID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "true",
		"message": "Video Conversion Complete",
	})
}

func UploadVideo(c *fiber.Ctx) error {

	var body Video
	var user *User

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	file, err := c.FormFile("video")
	if err != nil {
		return err
	}

	user, err = GetUserByUUID(uuid.MustParse(c.FormValue("user_id")))
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	contentType := file.Header.Get("content-type")
	_, exists := acceptedMimes[contentType]
	if !exists {
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(errors.New("unsupported video format").Error())
	}

	video.UserID = user.ID
	videoUUID, err := createVideo(&body, user, file)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusUnsupportedMediaType).JSON(videoUUID)
}

func DeleteVideo(c *fiber.Ctx) error {

	video = Video{}
	video.ID = uuid.MustParse(c.FormValue("uid"))

	err := video.DeleteVideo()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
		"status":  "false",
		"message": "video deleted",
	})
}
