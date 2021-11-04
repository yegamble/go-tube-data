package models

import (
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"os"
	"strconv"
	"strings"
	"time"
)

type TokenDetails struct {
	AccessToken    string
	RefreshToken   string
	AccessUuid     string
	RefreshUuid    string
	AtExpires      int64
	RtExpires      int64
	CookieHTTPOnly bool
	CookieSameSite string
	KeyGenerator   func() string
}

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
	IsBanned   bool
	IsAdmin    bool
}

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

func CreateJWTToken(userid uuid.UUID, isAdmin bool) (*TokenDetails, error) {
	td := &TokenDetails{}

	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewString()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewString()

	var err error

	//Creating Access Token
	os.Setenv("ACCESS_SECRET", os.Getenv("ACCESS_SECRET"))
	atClaims := jwtgo.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["is_admin"] = isAdmin
	rtClaims["user_id"] = userid
	atClaims["is_admin"] = isAdmin
	rtClaims["exp"] = td.RtExpires
	rt := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func ExtractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")

	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractTokenMetadata(c *fiber.Ctx) (*AccessDetails, error) {
	token, err := VerifyToken(c)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}

		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}
