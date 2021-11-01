package models

import (
	"github.com/gofiber/fiber/v2"
)

type Session struct {
	ID          uint64
	AccessToken string `json:"access_token"`
	UserID      uint64
	User        User   `json:"user_id" form:"user_id" gorm:"foreignKey:UserID;references:ID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Fingerprint string `json:"fingerprint"`
}

func SaveSession(userID uint64, cookieValue string, c *fiber.Ctx) error {
	var session Session

	session.AccessToken = cookieValue
	session.UserID = userID
	session.Fingerprint = c.Get("User-Agent")

	err := db.Create(&session).Error
	if err != nil {
		c.ClearCookie()
		return err
	}

	CreateUserLog("new session", userID, c)

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
