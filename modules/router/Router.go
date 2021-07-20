package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/yegamble/go-tube-api/modules/api/handler"
	"github.com/yegamble/go-tube-api/modules/api/user"
	"os"
)

var (
	res []handler.ErrorResponse
)

func SetRoutes() {
	app := fiber.New()

	routeHandler := app.Group("/", logger.New())

	routeHandler.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON("Welcome to " + os.Getenv("APP_NAME"))
	})

	routeHandler.Post("/user", func(c *fiber.Ctx) error {
		formResponse, error := user.RegisterUserFormParser(c)
		if error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(error)
		}
		return c.Status(fiber.StatusOK).JSON(formResponse)
	})

	err := app.Listen("localhost:3000")
	if err != nil {
		panic(err)
	}
}
