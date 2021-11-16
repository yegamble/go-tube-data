package models

import "github.com/gofiber/fiber/v2"

type Config struct {
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

func addBannedIPAddress(c *fiber.Ctx) error {
	ipAddress := c.FormValue("ip_address")
	err := BanIPAddress(ipAddress)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "ip address banned successfully",
	})
}
