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

	routeHandler := app.Group("/user", logger.New())
	//create user
	routeHandler.Post("/", func(c *fiber.Ctx) error {
		return user.RegisterUser(c)
	})

	//Login user
	routeHandler.Post("/login", func(c *fiber.Ctx) error {
		return user.Login(c)
	})

	//get user
	routeHandler.Get("/", func(c *fiber.Ctx) error {
		return user.GetUserByUsername(c)
	})

	err := app.Listen(os.Getenv("APP_URL") + ":" + os.Getenv("APP_PORT"))
	if err != nil {
		panic(err)
	}
}
