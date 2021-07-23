package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/yegamble/go-tube-api/modules/api/auth"
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

	//create user
	routeHandler.Post("/user", func(c *fiber.Ctx) error {
		return user.RegisterUser(c)
	})

	//Login user
	routeHandler.Post("/login", func(c *fiber.Ctx) error {
		return auth.Login(c)
	})

	//get user
	routeHandler.Get("/user", func(c *fiber.Ctx) error {
		return user.GetUserByName(c)
	})

	err := app.Listen(os.Getenv("APP_URL") + ":" + os.Getenv("APP_PORT"))
	if err != nil {
		panic(err)
	}
}
