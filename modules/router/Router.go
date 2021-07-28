package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/yegamble/go-tube-api/modules/api/models"
	"os"
	"time"
)

func SetRoutes() {

	app := fiber.New()

	userHandler := app.Group("/user", logger.New())

	//create user
	userHandler.Post("/create", func(c *fiber.Ctx) error {
		return models.RegisterUser(c)
	})

	//edit user
	userHandler.Patch("/edit/:id", func(c *fiber.Ctx) error {
		return models.EditUser(c)
	})

	//login user
	userHandler.Post("/login", func(c *fiber.Ctx) error {
		return models.Login(c)
	})

	//logout user
	userHandler.Post("/logout", func(c *fiber.Ctx) error {
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
	userHandler.Get("/search/*", func(c *fiber.Ctx) error {
		return models.SearchUsersByUsername(c)
	})

	//get user profile
	userHandler.Get("/:username", func(c *fiber.Ctx) error {
		return models.GetUserByUsername(c)
	})

	//delete user
	userHandler.Delete("/:uid", func(c *fiber.Ctx) error {
		return models.DeleteUser(c)
	})

	//get user by id
	userHandler.Get("/id/:id", func(c *fiber.Ctx) error {
		return models.GetUserByID(c)
	})

	//get user by uid
	userHandler.Get("/uid/:uid", func(c *fiber.Ctx) error {
		return models.SearchUserByUID(c)
	})

	//user upload
	uploadHandler := userHandler.Group("/upload")

	uploadHandler.Post("/profile-photo/:id", func(c *fiber.Ctx) error {
		return models.UploadUserPhoto(c, "profile_photo")
	})

	uploadHandler.Post("/header-photo/:id", func(c *fiber.Ctx) error {
		return models.UploadUserPhoto(c, "header_photo")
	})

	uploadHandler.Delete("/profile-photo/:id", func(c *fiber.Ctx) error {
		return models.DeleteUserPhoto(c, "profile_photo")
	})

	uploadHandler.Delete("/header-photo/:id", func(c *fiber.Ctx) error {
		return models.DeleteUserPhoto(c, "header_photo")
	})

	err := app.Listen(os.Getenv("APP_URL") + ":" + os.Getenv("APP_PORT"))
	if err != nil {
		panic(err)
	}

}
