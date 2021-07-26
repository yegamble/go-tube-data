package router

import (
	"github.com/alexedwards/scs/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/yegamble/go-tube-api/modules/api/handler"
	"github.com/yegamble/go-tube-api/modules/api/models"
	"log"
	"os"
	"time"
)

var (
	res []handler.ErrorResponse
)

func SetRoutes() {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	app := fiber.New()

	routeHandler := app.Group("/user", logger.New())

	//create user
	routeHandler.Post("/create", func(c *fiber.Ctx) error {
		return models.RegisterUser(c)
	})

	////edit user
	//routeHandler.Post("/edit/:uid", func(c *fiber.Ctx) error {
	//	return models.EditUser(c)
	//})

	//login user
	routeHandler.Post("/login", func(c *fiber.Ctx) error {
		store := session.New()
		log.Print(store)
		return models.Login(c)
	})

	//search user
	routeHandler.Get("/search/*", func(c *fiber.Ctx) error {
		return models.SearchUserByUsername(c)
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
		return models.GetUserByUID(c)
	})

	err := app.Listen(os.Getenv("APP_URL") + ":" + os.Getenv("APP_PORT"))
	if err != nil {
		panic(err)
	}

}
