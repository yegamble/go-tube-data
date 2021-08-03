package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

func CreateAuthRecord(userid uint64, td *TokenDetails) error {

	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(reflect.ValueOf((*td).AccessUuid).String(), strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(reflect.ValueOf((*td).RefreshUuid).String(), strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func CheckAuthorisationIsValid(c *fiber.Ctx) (*User, error) {
	userToken, err := VerifyToken(c)

	tokenAuth, err := ExtractTokenMetadata(c)
	if err != nil {
		return nil, c.Status(fiber.StatusUnauthorized).JSON("unauthorized")
	}

	userId, err := FetchAccessDetailsFromDB(tokenAuth)
	if err != nil {
		return nil, c.Status(fiber.StatusUnauthorized).JSON("unauthorized")
	}

	user.ID = userId

	claims := userToken.Claims.(jwt.MapClaims)
	log.Println(claims["user_id"])

	return &user, err
}

func FetchAccessDetailsFromDB(authDetails *AccessDetails) (uint64, error) {
	userid, err := client.Get(authDetails.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func SaveUserCookies(accessToken string, refreshToken string, c *fiber.Ctx) {
	cookie := new(fiber.Cookie)
	cookie.Name = "access_token"
	cookie.Value = accessToken
	cookie.Expires = time.Now().Add((24 * time.Hour) * 7)
	cookie.Domain = os.Getenv("APP_URL")
	cookie.Path = "/"
	cookie.SameSite = "lax"
	cookie.HTTPOnly = true
	//cookie2.Secure = true
	c.Cookie(cookie)

	cookie2 := new(fiber.Cookie)
	cookie2.Name = "refresh_token"
	cookie2.Value = refreshToken
	cookie2.Expires = time.Now().Add((24 * time.Hour) * 7)
	cookie2.Domain = os.Getenv("APP_URL")
	cookie2.Path = "/"
	cookie2.SameSite = "lax"
	cookie2.HTTPOnly = true
	//cookie2.Secure = true
	c.Cookie(cookie2)
}

func DeleteAuth(refreshUuid string) (int64, error) {

	deleted := client.Del(refreshUuid)
	if deleted.Err() != nil {
		return 0, deleted.Err()
	}

	return deleted.Val(), nil
}
