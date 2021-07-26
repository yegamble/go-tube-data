package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yegamble/go-tube-api/modules/api/auth"
	"os"
	"time"
)

type Session struct {
	ID           uint64
	SessionToken string `json:"session_token"`
	UserID       uint64
	User         User   `json:"user_id" form:"user_id" gorm:"foreignKey:UserID;references:ID"`
	Fingerprint  string `json:"fingerprint"`
}

func SaveSession(user User, c *fiber.Ctx) error {

	var session Session

	cookie := new(fiber.Cookie)
	cookie.Name = "session_token"
	cookie.Value = auth.GenerateSessionToken(64)
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Domain = os.Getenv("APP_URL")
	cookie.Path = "/"
	cookie.SameSite = "lax"
	cookie.Secure = true
	c.Cookie(cookie)

	session.SessionToken = cookie.Value
	session.UserID = user.ID
	session.Fingerprint = c.Get("User-Agent")

	err := db.Create(&session).Error
	if err != nil {
		c.ClearCookie("session_token")
		return err
	}

	InsertUserIPLog("logged in", user, c)

	return nil
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
