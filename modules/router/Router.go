package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/yegamble/go-tube-api/modules/api/models"
	"os"
	"time"
)

func SetRoutes() {

	app := fiber.New(fiber.Config{BodyLimit: 10000 * 1024 * 1024})

	userHandler := app.Group("/user", logger.New())

	//create user
	userHandler.Post("/create", func(c *fiber.Ctx) error {
		return models.RegisterUser(c)
	})

	//edit user
	userHandler.Patch("/edit/:id", models.AuthRequired(), func(c *fiber.Ctx) error {
		return models.EditUserRequest(c)
	})

	//login user
	userHandler.Post("/login", func(c *fiber.Ctx) error {
		return models.Login(c)
	})

	//refresh user token
	userHandler.Post("/refresh/token", models.AuthRequired(), func(c *fiber.Ctx) error {
		return models.RefreshAuthorisation(c)
	})

	//logout user
	userHandler.Post("/logout", models.AuthRequired(), func(c *fiber.Ctx) error {
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
		return models.FetchUserByUsername(c)
	})

	//delete user
	userHandler.Delete("/:uid", models.AuthRequired(), func(c *fiber.Ctx) error {
		return models.DeleteUser(c)
	})

	//get user by id
	userHandler.Get("/id/:id", models.AuthRequired(), func(c *fiber.Ctx) error {
		return models.FetchUserByID(c)
	})

	//get user by uid
	userHandler.Get("/uid/:uid", models.AuthRequired(), func(c *fiber.Ctx) error {
		return models.FetchUserByUID(c)
	})

	//user upload
	uploadHandler := userHandler.Group("/upload")

	uploadHandler.Post("/profile-photo/:id", models.AuthRequired(), func(c *fiber.Ctx) error {
		return models.UploadUserPhoto(c, "profile_photo")
	})

	uploadHandler.Post("/header-photo/:id", models.AuthRequired(), func(c *fiber.Ctx) error {
		return models.UploadUserPhoto(c, "header_photo")
	})

	uploadHandler.Delete("/profile-photo/:id", models.AuthRequired(), func(c *fiber.Ctx) error {
		return models.DeleteUserPhoto(c, "profile_photo")
	})

	uploadHandler.Delete("/header-photo/:id", models.AuthRequired(), func(c *fiber.Ctx) error {
		return models.DeleteUserPhoto(c, "header_photo")
	})

	videoHandler := app.Group("api/video", logger.New())

	videoHandler.Get("/", func(c *fiber.Ctx) error {
		return models.GetAllVideos(c)
	})

	videoHandler.Get("/trending", func(c *fiber.Ctx) error {
		return nil
	})

	videoHandler.Get("/top", func(c *fiber.Ctx) error {
		return nil
	})

	videoHandler.Get("/id/:id", func(c *fiber.Ctx) error {
		return models.FetchVideoByID(c)
	})

	videoHandler.Get("/uid/:id", func(c *fiber.Ctx) error {
		return models.FetchVideoByUID(c)
	})

	videoHandler.Post("/upload", func(c *fiber.Ctx) error {
		return models.UploadVideo(c)
	})

	videoHandler.Post("/convert", func(c *fiber.Ctx) error {
		return models.TriggerConversionByVideoUID(c)
	})

	err := app.Listen(os.Getenv("APP_URL") + ":" + os.Getenv("APP_PORT"))
	if err != nil {
		panic(err)
	}

}
