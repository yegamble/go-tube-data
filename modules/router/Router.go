package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/yegamble/go-tube-api/modules/api/handler"
	"github.com/yegamble/go-tube-api/modules/api/models"
	"os"
	"time"
)

var (
	res []handler.ErrorResponse
)

func SetRoutes() {

	app := fiber.New()

	routeHandler := app.Group("/user", logger.New())

	//create user
	routeHandler.Post("/create", func(c *fiber.Ctx) error {
		return models.RegisterUser(c)
	})

	//edit user
	routeHandler.Patch("/edit/:id", func(c *fiber.Ctx) error {
		return models.EditUser(c)
	})

	//login user
	routeHandler.Post("/login", func(c *fiber.Ctx) error {
		return models.Login(c)
	})

	//logout user
	routeHandler.Post("/logout", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name: "session_token",
			// Set expiry date to the past
			Expires:  time.Now().Add(-(time.Hour * 2)),
			HTTPOnly: true,
			Domain:   os.Getenv("APP_URL"),
			Path:     "/",
			SameSite: "lax",
		})
		return models.Logout(c)
	})

	//search user
	routeHandler.Get("/search/*", func(c *fiber.Ctx) error {
		return models.SearchUsersByUsername(c)
	})

	//get user profile
	routeHandler.Get("/:username", func(c *fiber.Ctx) error {
		return models.GetUserByUsername(c)
	})

	//delete user
	routeHandler.Delete("/:uid", func(c *fiber.Ctx) error {
		return models.DeleteUser(c)
	})

	//get user by id
	routeHandler.Get("/id/:id", func(c *fiber.Ctx) error {
		return models.GetUserByID(c)
	})

	//get user by uid
	routeHandler.Get("/uid/:uid", func(c *fiber.Ctx) error {
		return models.SearchUserByUID(c)
	})

	err := app.Listen(os.Getenv("APP_URL") + ":" + os.Getenv("APP_PORT"))
	if err != nil {
		panic(err)
	}

}
