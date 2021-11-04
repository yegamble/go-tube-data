package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Session struct {
	ID          uint64
	AccessToken string `json:"access_token"`
	UserID      uuid.UUID
	User        User   `json:"user_id" form:"user_id" gorm:"foreignKey:UserID;references:ID;OnUpdate:CASCADE,OnDelete:CASCADE;type:varchar(255);"`
	Fingerprint string `json:"fingerprint"`
}

func SaveSession(userID uuid.UUID, cookieValue string, c *fiber.Ctx) error {
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
