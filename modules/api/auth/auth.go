package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/yegamble/go-tube-api/modules/api/user"
	"os"
	"time"
)

func Login(c *fiber.Ctx) error {
	var u user.User
	if err := c.BodyParser(&u); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON("Invalid json provided")
	}
	//compare the user from the request, with the one we defined:
	if c.FormValue("username") != u.Username || c.FormValue("password") != u.Password {
		return c.Status(fiber.StatusUnauthorized).JSON("Invalid Login Details")
	}

	match, err := ComparePasswordAndHash(c.FormValue("password"))
	if err != nil {
		return err
	} else if !match {
		return errors.New("invalid login details")
	}

	token, err := CreateToken(u.ID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(token)
}

func CreateToken(userid uint64) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", os.Getenv("ACCESS_SECRET")) //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
