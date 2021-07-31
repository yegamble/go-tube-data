package auth

import (
	"crypto/rand"
	"encoding/hex"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	jwt "github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"os"
	"time"
)

func AuthRequired() fiber.Handler {
	return jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
		SigningKey: []byte(os.Getenv("ACCESS_SECRET")),
	})
}

func CreateJWTToken(userid uuid.UUID) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", os.Getenv("ACCESS_SECRET")) //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func GenerateSessionToken(length int) string {

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
