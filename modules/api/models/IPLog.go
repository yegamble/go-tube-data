package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yegamble/go-tube-api/database"
	"time"
)

type IPLog struct {
	ID        uint64 `json:"id" gorm:"primary_key"`
	UserID    uint64 `json:"user_id" form:"user_id"`
	User      User   `gorm:"foreignKey:UserID;references:ID"`
	IPAddress string `json:"ip_address" gorm:"type:text"`
	Activity  string `json:"activity" gorm:"type:text"`
	CreatedAt time.Time
}

type BannedIPLog struct {
	IPAddress string
}

func InsertUserIPLog(activity string, u User, ctx *fiber.Ctx) (uint64, error) {

	var log IPLog
	log.UserID = u.ID
	log.IPAddress = ctx.IP()
	log.Activity = activity

	result := db.Create(&log)

	return log.ID, result.Error
}

func clearIPLogs(clearAll bool) error {
	db := database.DBConn
	result := db.Where("created_at < NOW() - INTERVAL 4 WEEK").Delete(IPLog{})
	return result.Error
}
