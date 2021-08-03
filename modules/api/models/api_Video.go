package models

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

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

	return c.Status(fiber.StatusUnsupportedMediaType).JSON(video.ShortID)
}
