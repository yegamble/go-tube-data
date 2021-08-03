package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yegamble/go-tube-api/modules/api/auth"
	"os"
	"strconv"
	"time"
)

type Session struct {
	ID           uint64
	SessionToken string `json:"session_token"`
	UserID       uint64
	User         User   `json:"user_id" form:"user_id" gorm:"foreignKey:UserID;references:ID"`
	Fingerprint  string `json:"fingerprint"`
}

func SaveSession(userID uint64, accessToken string, refreshToken string, c *fiber.Ctx) error {

	var session Session

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

	session.SessionToken = cookie.Value
	session.UserID = userID
	session.Fingerprint = c.Get("User-Agent")

	err := db.Create(&session).Error
	if err != nil {
		c.ClearCookie("session_token")
		return err
	}

	CreateUserLog("new session", userID, c)

	return nil
}

func FetchAuth(authD *auth.AccessDetails) (uint64, error) {
	userid, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

//func DeleteSession(user User, c *fiber.Ctx) error{
//
//	err := db.Where("session_token = ?", c.Cookies("session_token")).Delete(&Session{}).Error
//	if err != nil {
//		return err
//	}
//
//	c.ClearCookie()
//
//}
