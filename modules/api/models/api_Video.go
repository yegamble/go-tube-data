package models

import (
	"errors"
	"github.com/gofiber/fiber/v2"
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
	err = createVideo(&body, user, file)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusUnsupportedMediaType).JSON(video.Slug)
}
